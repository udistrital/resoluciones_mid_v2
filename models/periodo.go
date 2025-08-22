package models

type Periodo struct {
	Id                int
	Nombre            string
	Descripcion       string
	Year              int
	Ciclo             string
	CodigoAbreviacion string
	Activo            bool
	AplicacionId      int
	InicioVigencia    string
	FinVigencia       string
	FechaCreacion     string
	FechaModificacion string
}
