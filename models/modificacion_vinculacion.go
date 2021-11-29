package models

type ModificacionVinculacion struct {
	Id                             int
	ModificacionResolucionId       *ModificacionResolucion
	VinculacionDocenteCanceladaId  *VinculacionDocente
	VinculacionDocenteRegistradaId *VinculacionDocente
	Horas                          float64
	Activo                         bool
	FechaCreacion                  string
	FechaModificacion              string
}
