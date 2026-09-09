package maps

import "errors"

// Mapped is an HTTP status + public error message for a known sentinel.
type Mapped struct {
	Status  int
	Message string
}

type mapperFunc func(err error) *Mapped

// Lookup walks context mappers until one recognizes err.
func Lookup(err error) *Mapped {
	for _, mapFn := range contextMappers {
		if mapped := mapFn(err); mapped != nil {
			return mapped
		}
	}
	return nil
}

var contextMappers = []mapperFunc{
	mapAuth,
	mapUser,
	mapToken,
	mapTeam,
	mapTag,
	mapEnvironment,
	mapProject,
	mapSecret,
}

func match(err, sentinel error, status int, message string) *Mapped {
	if errors.Is(err, sentinel) {
		return &Mapped{Status: status, Message: message}
	}
	return nil
}

func firstMatch(try ...*Mapped) *Mapped {
	for _, m := range try {
		if m != nil {
			return m
		}
	}
	return nil
}
