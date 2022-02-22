package models

type ObjetoCategoriaDocente struct {
	CategoriaDocente struct {
		Anio           string `json:"anio" xml:"anio"`
		Categoria      string `json:"categoria" xml:"categoria"`
		Identificacion string `json:"identificacion" xml:"identificacion"`
		IDCategoria    string `json:"id_categoria" xml:"id_categoria"`
		Periodo        string `json:"periodo" xml:"periodo"`
	} `json:"categoria_docente" xml:"categoria_docente"`
}
