package models

type ObjetoCargaLectiva struct {
	CargasLectivas struct {
		CargaLectiva []CargaLectiva `json:"carga_lectiva"`
	} `json:"cargas_lectivas"`
}

type CargaLectiva struct {
	Anio                  string `json:"anio"`
	HorasLectivas         string `json:"horas_lectivas"`
	DocDocente            string `json:"docente_documento"`
	IDFacultad            string `json:"id_facultad"`
	IDProyecto            string `json:"id_proyecto"`
	DependenciaAcademica  int    //`json:"id_proyecto_condor"`
	IDTipoVinculacion     string `json:"id_tipo_vinculacion"`
	NombreFacultad        string `json:"facultad_nombre"`
	NombreProyecto        string `json:"proyecto_nombre"`
	NombreTipoVinculacion string `json:"tipo_vinculacion_nombre"`
	Periodo               string `json:"periodo"`
	DocenteApellido       string `json:"docente_apellido"`
	DocenteNombre         string `json:"docente_nombre"`
	CategoriaNombre       string
	IDCategoria           string
	IdProveedor           string
	Emerito               bool
}
