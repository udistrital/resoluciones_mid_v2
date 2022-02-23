package models

import (
	"time"
)

type RelacionParametro struct {
	Id             int
	Descripcion    string
	EstadoRegistro bool
	FechaRegistro  time.Time
}
