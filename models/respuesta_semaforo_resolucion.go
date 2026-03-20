package models

type RespuestaSemaforoResolucion struct {
	Resumen ResumenSemaforoResolucion   `json:"resumen"`
	Detalle []EstadoSemaforoVinculacion `json:"detalle"`
}
