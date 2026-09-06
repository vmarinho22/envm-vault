package audit

import "github.com/google/uuid"

// Event is a domain audit event. Origin fields are filled from context by Record.
type Event struct {
	Action       string
	ActorUserID  *uuid.UUID
	ActorEmail   string
	ActorRole    string
	TargetUserID *uuid.UUID
}
