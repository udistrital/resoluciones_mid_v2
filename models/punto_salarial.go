package models

import (
	"time"
)

type PuntoSalarial struct {
	Decreto       string
	ValorPunto    int
	FechaRegistro time.Time
	Vigencia      int
	Id            int
}
