package models

import "time"

type InformacionPersonaNatural struct {
	TipoDocumento               *ParametroEstandar
	Id                          string
	DigitoVerificacion          float64
	PrimerApellido              string
	SegundoApellido             string
	PrimerNombre                string
	SegundoNombre               string
	Cargo                       string
	IdPaisNacimiento            float64
	Perfil                      *ParametroEstandar
	Profesion                   string
	Especialidad                string
	MontoCapitalAutorizado      float64
	Genero                      string
	FechaExpedicionDocumento    time.Time
	IdCiudadExpedicionDocumento float64
	NomProveedor                string
	CiudadExpedicionDocumento   string
}
