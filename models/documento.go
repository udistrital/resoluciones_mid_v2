package models

type DocumentoContainer struct {
	Res    Documento `json:"res"`
	Status string
	Error  string
}

type Documento struct {
	Id                int
	Nombre            string
	Descripcion       string
	Enlace            string
	TipoDocumento     map[string]interface{}
	Metadatos         string
	Activo            bool
	FechaCreacion     string
	Fechamodificacion string
}
