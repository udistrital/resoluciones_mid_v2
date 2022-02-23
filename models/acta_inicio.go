package models

import (
	"time"
)

type ActaInicio struct {
	Id             int
	NumeroContrato string
	Vigencia       int
	FechaInicio    time.Time
	FechaFin       time.Time
	Descripcion    string
	Usuario        string
	FechaRegistro  time.Time
}
