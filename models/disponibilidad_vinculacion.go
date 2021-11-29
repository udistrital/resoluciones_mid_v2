package models

type DisponibilidadVinculacion struct {
	Id                   int
	Disponibilidad       int
	Rubro                string
	Valor                float64
	VinculacionDocenteId *VinculacionDocente
	Activo               bool
	FechaCreacion        string
	FechaModificacion    string
}
