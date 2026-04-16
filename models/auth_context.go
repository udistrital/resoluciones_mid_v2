package models

type RequestAuthContext struct {
	NumeroDocumento string
	Roles           []string
	Source          string
	Trusted         bool
}
