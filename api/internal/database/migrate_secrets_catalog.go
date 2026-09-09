package database

import (
	"fmt"
	"log"
	"sort"
	"time"

	"github.com/envm-vault/envm-vault/api/internal/model"
	"github.com/envm-vault/envm-vault/api/internal/secretcrypto"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// migrateSecretsToGlobalCatalog converts per-project secrets into a global catalog
// with project_secrets M2M. Idempotent: no-op when secrets.project_id is absent.
func migrateSecretsToGlobalCatalog(db *gorm.DB, crypto secretcrypto.Crypter) error {
	if !db.Migrator().HasTable("secrets") {
		return nil
	}
	if !db.Migrator().HasColumn(&legacySecret{}, "project_id") {
		return nil
	}
	if crypto == nil {
		return fmt.Errorf("migrate secrets catalog: encryption required to merge duplicate names")
	}

	log.Println("migrating secrets to global catalog…")

	if err := db.Exec(`
		CREATE TABLE IF NOT EXISTS project_secrets (
			project_id uuid NOT NULL,
			secret_id uuid NOT NULL,
			PRIMARY KEY (project_id, secret_id)
		)`).Error; err != nil {
		return fmt.Errorf("create project_secrets: %w", err)
	}

	var rows []legacySecret
	if err := db.Order("created_at ASC").Find(&rows).Error; err != nil {
		return fmt.Errorf("load legacy secrets: %w", err)
	}

	byName := make(map[string][]legacySecret)
	for _, row := range rows {
		byName[row.Name] = append(byName[row.Name], row)
	}

	for name, group := range byName {
		sort.Slice(group, func(i, j int) bool {
			if group[i].CreatedAt.Equal(group[j].CreatedAt) {
				return group[i].ID.String() < group[j].ID.String()
			}
			return group[i].CreatedAt.Before(group[j].CreatedAt)
		})
		canonical := group[0]
		if err := attachProjectSecret(db, canonical.ProjectID, canonical.ID); err != nil {
			return err
		}

		for _, dup := range group[1:] {
			if err := attachProjectSecret(db, dup.ProjectID, canonical.ID); err != nil {
				return err
			}
			if err := mergeSecretValues(db, crypto, canonical.ID, dup.ID); err != nil {
				return fmt.Errorf("merge values for %s: %w", name, err)
			}
			if err := db.Where("id = ?", dup.ID).Delete(&model.Secret{}).Error; err != nil {
				return fmt.Errorf("delete duplicate secret %s: %w", dup.ID, err)
			}
		}
	}

	if err := db.Exec(`DROP INDEX IF EXISTS idx_secrets_project_name`).Error; err != nil {
		return fmt.Errorf("drop idx_secrets_project_name: %w", err)
	}
	if db.Migrator().HasColumn(&legacySecret{}, "project_id") {
		if err := db.Migrator().DropColumn(&legacySecret{}, "project_id"); err != nil {
			return fmt.Errorf("drop secrets.project_id: %w", err)
		}
	}

	log.Println("secrets catalog migration complete")
	return nil
}

type legacySecret struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey"`
	ProjectID uuid.UUID `gorm:"type:uuid;column:project_id"`
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (legacySecret) TableName() string { return "secrets" }

func attachProjectSecret(db *gorm.DB, projectID, secretID uuid.UUID) error {
	err := db.Clauses(clause.OnConflict{DoNothing: true}).Create(&model.ProjectSecret{
		ProjectID: projectID,
		SecretID:  secretID,
	}).Error
	if err != nil {
		return fmt.Errorf("attach project_secret: %w", err)
	}
	return nil
}

func mergeSecretValues(db *gorm.DB, crypto secretcrypto.Crypter, canonicalID, dupID uuid.UUID) error {
	var dupValues []model.SecretValue
	if err := db.Where("secret_id = ?", dupID).Find(&dupValues).Error; err != nil {
		return fmt.Errorf("load duplicate values: %w", err)
	}

	for _, dv := range dupValues {
		var existing model.SecretValue
		err := db.Where("secret_id = ? AND environment_id = ?", canonicalID, dv.EnvironmentID).
			First(&existing).Error
		if err != nil && err != gorm.ErrRecordNotFound {
			return fmt.Errorf("lookup canonical value: %w", err)
		}

		plaintext, decErr := crypto.Decrypt(
			dv.Ciphertext,
			secretcrypto.AAD(dupID, dv.EnvironmentID),
			dv.KeyVersion,
		)
		if decErr != nil {
			return fmt.Errorf("decrypt duplicate value: %w", decErr)
		}

		if err == gorm.ErrRecordNotFound {
			blob, version, encErr := crypto.Encrypt(plaintext, secretcrypto.AAD(canonicalID, dv.EnvironmentID))
			if encErr != nil {
				return fmt.Errorf("re-encrypt value: %w", encErr)
			}
			if createErr := db.Create(&model.SecretValue{
				SecretID:      canonicalID,
				EnvironmentID: dv.EnvironmentID,
				Ciphertext:    blob,
				KeyVersion:    version,
			}).Error; createErr != nil {
				return fmt.Errorf("create merged value: %w", createErr)
			}
			continue
		}

		// Both have a value: keep the more recently updated ciphertext (re-keyed to canonical AAD).
		keepDup := !dv.UpdatedAt.Before(existing.UpdatedAt)
		if !keepDup {
			continue
		}
		blob, version, encErr := crypto.Encrypt(plaintext, secretcrypto.AAD(canonicalID, dv.EnvironmentID))
		if encErr != nil {
			return fmt.Errorf("re-encrypt newer value: %w", encErr)
		}
		existing.Ciphertext = blob
		existing.KeyVersion = version
		existing.UpdatedAt = dv.UpdatedAt
		if saveErr := db.Save(&existing).Error; saveErr != nil {
			return fmt.Errorf("update merged value: %w", saveErr)
		}
	}

	if err := db.Where("secret_id = ?", dupID).Delete(&model.SecretValue{}).Error; err != nil {
		return fmt.Errorf("delete duplicate values: %w", err)
	}
	return nil
}

func ensureSecretSearchIndexes(db *gorm.DB) error {
	if err := db.Exec(`CREATE EXTENSION IF NOT EXISTS pg_trgm`).Error; err != nil {
		return fmt.Errorf("enable pg_trgm: %w", err)
	}
	if err := db.Exec(`
		CREATE INDEX IF NOT EXISTS idx_secrets_name_trgm
		ON secrets USING GIN (name gin_trgm_ops)`).Error; err != nil {
		return fmt.Errorf("create secrets name trigram index: %w", err)
	}
	return nil
}
