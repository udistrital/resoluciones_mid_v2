package models

type ResolucionEstado struct {
	Id                 int
	Usuario            string
	EstadoResolucionId int
	ResolucionId       *Resolucion
	Activo             bool
	FechaCreacion      string
	FechaModificacion  string
}
