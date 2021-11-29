package models

type Parametro struct {
	Id                int
	Nombre            string
	Descripcion       string
	CodigoAbreviacion string
	Activo            bool
	NumeroOrden       float64
	TipoParametroId   *TipoParametro
	ParametroPadreId  *Parametro
}
