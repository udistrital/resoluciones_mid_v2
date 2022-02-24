package models

import (
	"time"
)

type EstadoContrato struct {
	NombreEstado  string
	FechaRegistro time.Time
	Id            int
}
