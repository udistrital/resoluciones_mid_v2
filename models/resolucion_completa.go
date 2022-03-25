package models

type ResolucionCompleta struct {
	Vinculacion             ResolucionVinculacionDocente
	Consideracion           string
	Preambulo               string
	Vigencia                int
	Numero                  string
	Id                      int
	Articulos               []Articulo
	OrdenadorGasto          OrdenadorGasto
	Titulo                  string
	CuadroResponsabilidades string
}
