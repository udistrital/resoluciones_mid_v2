package models

type ObjetoPrevinculaciones struct {
	Docentes       []CargaLectiva
	ResolucionData *ResolucionVinculacionDocente
	NumeroSemanas  int
	Vigencia       int
	Disponibilidad []DocumentoPresupuestal
}
