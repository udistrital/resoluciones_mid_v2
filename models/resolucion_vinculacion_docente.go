package models

type ResolucionVinculacionDocente struct {
	Id                int
	FacultadId        int
	Dedicacion        string
	NivelAcademico    string
	Activo            bool
	FechaCreacion     string
	FechaModificacion string
}
