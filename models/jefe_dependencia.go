package models

import (
	"time"
)

type JefeDependencia struct {
	Id             int
	FechaInicio    time.Time
	FechaFin       time.Time
	TerceroId      int
	DependenciaId  int
	ActaAprobacion string
}
