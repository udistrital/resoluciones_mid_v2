package models

type Vinculaciones struct {
	Id                   int
	Nombre               string
	PersonaId            float64
	TipoDocumento        string
	ExpedicionDocumento  string
	NumeroContrato       string
	Vigencia             int
	Categoria            string
	Dedicacion           string
	NumeroHorasSemanales int
	NumeroSemanas        int
	Disponibilidad       int
	RegistroPresupuestal int
	ValorContratoFormato string
	ProyectoCurricularId int
}
