package models

type ResumenSemaforoResolucion struct {
	ResolucionId              int     `json:"resolucion_id"`
	Total                     int     `json:"total"`
	TotalConRp                int     `json:"total_con_rp"`
	Completas                 int     `json:"completas"`
	PendientesTitan           int     `json:"pendientes_titan"`
	SinRp                     int     `json:"sin_rp"`
	PorcentajeCompletas       float64 `json:"porcentaje_completas"`
	PorcentajePendientesTitan float64 `json:"porcentaje_pendientes_titan"`
	PorcentajeSinRp           float64 `json:"porcentaje_sin_rp"`
}
