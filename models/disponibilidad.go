package models

import "time"

type Disponibilidad struct {
	Id                        int
	Vigencia                  float64
	NumeroDisponibilidad      float64
	FechaRegistro             time.Time
	Solicitud                 int
	DisponibilidadApropiacion []*DisponibilidadApropiacion
}
