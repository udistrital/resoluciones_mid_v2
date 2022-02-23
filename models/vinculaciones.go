package models

type Vinculaciones struct {
	Id                   int
	Nombre               string
	PersonaId            float64
	NumeroContrato       string
	Vigencia             int
	Categoria            string
	Dedicacion           string
	NumeroHorasSemanales int
	NumeroSemanas        int
	Disponibilidad       int
	ValorContratoFormato string
	ProyectoCurricularId int
}
