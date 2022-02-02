package models

type Vinculaciones struct {
	Id                   int
	Nombre               string
	PersonaId            float64
	Categoria            string
	Dedicacion           string
	NumeroHorasSemanales int
	NumeroSemanas        int
	Disponibilidad       int
	ValorContratoFormato string
	ProyectoCurricularId int16
}
