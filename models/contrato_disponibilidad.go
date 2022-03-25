package models

import "time"

type ContratoDisponibilidad struct {
	Id             int
	NumeroCdp      int
	NumeroContrato string
	Vigencia       int
	Estado         bool
	FechaRegistro  time.Time
	VigenciaCdp    int
}
