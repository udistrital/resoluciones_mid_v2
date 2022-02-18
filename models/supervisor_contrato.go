package models

import (
	"time"
)

type SupervisorContrato struct {
	Id                    int
	Nombre                string
	Documento             int
	Cargo                 string
	SedeSupervisor        string
	DependenciaSupervisor string
	Tipo                  int
	Estado                bool
	DigitoVerificacion    int
	FechaInicio           time.Time
	FechaFin              time.Time
}
