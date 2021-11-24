package models

type ModificacionResolucion struct {
	Id                   int
	ResolucionNuevaId    *Resolucion
	ResolucionAnteriorId *Resolucion
	Activo               bool
	FechaCreacion        string
	FechaModificacion    string
}
