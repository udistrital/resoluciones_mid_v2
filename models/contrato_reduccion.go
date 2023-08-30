package models

import "time"

type Reduccion struct {
	Vigencia            int
	Documento           string
	FechaReduccion      time.Time
	NivelAcademico      string
	Semanas             int
	SemanasAnteriores   int
	ContratosOriginales []ContratoReducir
	ContratoNuevo       *ContratoReducido
}

type ContratoReducir struct {
	NumeroContratoOriginal string
	ValorContratoReducido  float64
	DesagregadoOriginal    *map[string]float64
}

type ContratoReducido struct {
	NumeroContratoReduccion string
	ValorContratoReduccion  float64
	DesagregadoReduccion    *map[string]float64
}
