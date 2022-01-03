package models

type ContenidoResolucion struct {
	Resolucion  Resolucion
	Articulos   []Articulo
	Vinculacion ResolucionVinculacionDocente
	Usuario     string
}
