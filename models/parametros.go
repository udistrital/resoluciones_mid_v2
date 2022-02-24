package models

import (
	"time"
)

type Parametros struct {
	Id                int
	Descripcion       string
	CodigoContraloria string
	RelParametro      *RelacionParametro
	EstadoRegistro    bool
	FechaRegistro     time.Time
}
