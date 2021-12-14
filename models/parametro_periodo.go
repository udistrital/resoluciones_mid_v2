package models

type ParametroPeriodo struct {
	Id          int
	Valor       string
	Activo      bool
	ParametroId *Parametro
}
