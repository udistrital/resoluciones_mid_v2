package models

type AuthenticatedContext struct {
	NumeroDocumento string
	Roles           []string
	Source          string
}
