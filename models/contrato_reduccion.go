package models

import "time"

type Reduccion struct {
	NumeroContratoReduccion string
	Vigencia                int
	Documento               string
	ValorContratoReduccion  float64
	FechaReduccion          time.Time
	ContratosOriginales     []ContratoReducir
	DesagregadoReduccion    *map[string]float64
}

type ContratoReducir struct {
	NumeroContratoOriginal string
	DesagregadoOriginal    *map[string]float64
}
