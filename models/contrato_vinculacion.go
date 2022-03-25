package models

type ContratoVinculacion struct {
	ContratoGeneral    *ContratoGeneral
	VinculacionDocente *VinculacionDocente
	ActaInicio         *ActaInicio
	Cdp                *ContratoDisponibilidad
}
