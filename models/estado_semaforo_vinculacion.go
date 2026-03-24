package models

type EstadoSemaforoVinculacion struct {
	VinculacionId       int    `json:"vinculacion_id"`
	NumeroDocumento     int    `json:"numero_documento"`
	NumeroContrato      string `json:"numero_contrato"`
	Vigencia            int    `json:"vigencia"`
	NumeroRp            int    `json:"numero_rp"`
	VigenciaRp          int    `json:"vigencia_rp"`
	TieneRpResoluciones bool   `json:"tiene_rp_resoluciones"`
	TieneRpTitan        bool   `json:"tiene_rp_titan"`
	EstadoCodigo        string `json:"estado_codigo"`
	EstadoNombre        string `json:"estado_nombre"`
	Prioridad           int    `json:"prioridad"`
}
