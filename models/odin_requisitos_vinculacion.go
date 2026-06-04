package models

type OdinAuthRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Version  string `json:"version"`
}

type OdinAuthResponse struct {
	Token       string `json:"token"`
	AccessToken string `json:"access_token"`
}

type OdinRequisitosVinculacionRequest struct {
	Parametros OdinRequisitosVinculacionParametros `json:"parametros"`
}

type OdinRequisitosVinculacionParametros struct {
	Anio      string `json:"anio"`
	Periodo   string `json:"periodo"`
	IdUsuario string `json:"id_usuario"`
}

type OdinRequisitoVinculacion struct {
	Id                   int    `json:"id"`
	Anio                 int    `json:"anio"`
	Periodo              int    `json:"periodo"`
	IdUsuario            string `json:"id_usuario"`
	ApruebaSoporte       string `json:"aprueba_soporte"`
	IdCategoria          string `json:"id_categoria"`
	Categoria            string `json:"categoria"`
	Observacion          string `json:"observacion"`
	Estado               string `json:"estado"`
	RegistroTercero      string `json:"registro_tercero"`
	RegistroProveedor    string `json:"registro_proveedor"`
	RegistroTipoVinculo  string `json:"registro_tipo_vinculo"`
	RegistroNoCruce      string `json:"registro_no_cruce"`
	RegistroPreCarga     string `json:"registro_pre_carga"`
	NombreEstado         string `json:"nombre_estado"`
	NoCruce              int    `json:"nocruce"`
	ResolucionVinculable string `json:"resolucion_vinculable"`
}

type DocenteNoVinculable struct {
	Documento string   `json:"documento"`
	Nombre    string   `json:"nombre"`
	Motivos   []string `json:"motivos"`
}
