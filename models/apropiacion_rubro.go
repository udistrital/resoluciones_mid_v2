package models

//ApropiacionRubro relaciona la información de la apropiación con el rubro asociado
type ApropiacionRubro struct {
	Id       int
	Vigencia float64
	Rubro    *Rubro
	Valor    float64
	Estado   *ApropiacionRubroEstado
	Saldo    int
}
