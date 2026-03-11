package models

type DependenciaUsuario struct {
	CodigoDependencia int    `json:"codigo_dependencia"`
	IdOikos           int    `json:"id_oikos"`
	Nombre            string `json:"nombre,omitempty"`
	Rol               string `json:"rol,omitempty"`
}

type AlcanceUsuario struct {
	RolPrincipal string               `json:"rol_principal"`
	EsGlobal     bool                 `json:"es_global"`
	Dependencias []DependenciaUsuario `json:"dependencias"`
}

type DecanoFacultadResponse struct {
	Facultad struct {
		Decano []struct {
			CodigoFacultad string `json:"codigo_facultad"`
			NombreFacultad string `json:"facultad"`
			NombreDecano   string `json:"nombre"`
			FechaDesde     string `json:"fecha_desde"`
			FechaHasta     string `json:"fecha_hasta"`
		} `json:"decano"`
	} `json:"facultad"`
}

type AsistenteFacultadResponse struct {
	Asistente struct {
		Facultad []struct {
			CodigoDependencia string `json:"codigo_dependecia"`
			TipoUsuario       string `json:"tipo_usuario"`
			NombreDependencia string `json:"nombre_dependencia"`
		} `json:"facultad"`
	} `json:"asistente"`
}

type HomologacionFacultadResponse struct {
	Homologacion struct {
		IdOikos string `json:"id_oikos"`
		IdGedep string `json:"id_gedep"`
	} `json:"homologacion"`
}
