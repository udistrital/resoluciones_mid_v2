package models

type ResumenDashboardResolucion struct {
	ResolucionId              int     `json:"resolucion_id"`
	NumeroResolucion          string  `json:"numero_resolucion"`
	Vigencia                  int     `json:"vigencia"`
	DependenciaId             int     `json:"dependencia_id"`
	DependenciaNombre         string  `json:"dependencia_nombre"`
	Total                     int     `json:"total"`
	TotalConRp                int     `json:"total_con_rp"`
	Completas                 int     `json:"completas"`
	PendientesTitan           int     `json:"pendientes_titan"`
	SinRp                     int     `json:"sin_rp"`
	PorcentajeCompletas       float64 `json:"porcentaje_completas"`
	PorcentajePendientesTitan float64 `json:"porcentaje_pendientes_titan"`
	PorcentajeSinRp           float64 `json:"porcentaje_sin_rp"`
	EstadoGeneralCodigo       string  `json:"estado_general_codigo"`
	EstadoGeneralNombre       string  `json:"estado_general_nombre"`
}
