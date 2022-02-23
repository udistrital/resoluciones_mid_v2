package models

import (
	"time"
)

type ContratoEstado struct {
	NumeroContrato string
	Vigencia       int
	FechaRegistro  time.Time
	Id             int
	Estado         *EstadoContrato
	Usuario        string
}
