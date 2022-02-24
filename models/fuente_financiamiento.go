package models

type FuenteFinanciamiento struct {
	Id                       int
	Descripcion              string
	Nombre                   string
	Codigo                   string
	TipoFuenteFinanciamiento *TipoFuenteFinanciamiento
}
