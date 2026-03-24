package models

type ResumenGlobalDashboardResoluciones struct {
	TotalResoluciones               int     `json:"total_resoluciones"`
	ResolucionesCompletas           int     `json:"resoluciones_completas"`
	ResolucionesConPendientesTitan  int     `json:"resoluciones_con_pendientes_titan"`
	ResolucionesConSinRp            int     `json:"resoluciones_con_sin_rp"`
	PorcentajeResolucionesCompletas float64 `json:"porcentaje_resoluciones_completas"`
}

type RespuestaDashboardResoluciones struct {
	ResumenGlobal ResumenGlobalDashboardResoluciones `json:"resumen_global"`
	Resoluciones  []ResumenDashboardResolucion       `json:"resoluciones"`
	Total         int                                `json:"total"`
	Limit         int                                `json:"limit"`
	Offset        int                                `json:"offset"`
}
