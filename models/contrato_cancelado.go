package models

import "time"

type ContratoCancelado struct {
	Id                int
	NumeroContrato    string
	Vigencia          int
	FechaCancelacion  time.Time
	MotivoCancelacion string
	Usuario           string
	FechaRegistro     time.Time
	Estado            bool
}
