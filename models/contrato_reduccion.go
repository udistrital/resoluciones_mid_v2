package models

import "time"

type Reduccion struct {
	Vigencia            int
	Documento           string
	FechaReduccion      time.Time
	ContratosOriginales []ContratoReducir
	ContratoNuevo       *ContratoReducido
}

type ContratoReducir struct {
	NumeroContratoOriginal string
	DesagregadoOriginal    *map[string]float64
}

type ContratoReducido struct {
	NumeroContratoReduccion string
	ValorContratoReduccion  float64
	DesagregadoReduccion    *map[string]float64
}
