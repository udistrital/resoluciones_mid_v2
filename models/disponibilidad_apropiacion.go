package models

type DisponibilidadApropiacion struct {
	Id                   int
	Disponibilidad       *Disponibilidad
	Apropiacion          *ApropiacionRubro
	Valor                float64
	FuenteFinanciamiento *FuenteFinanciamiento
}
