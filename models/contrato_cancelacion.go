package models

import "time"

type ContratoCancelacion struct {
	NumeroContrato string
	Vigencia       int
	Documento      string
	ValorContrato  float64
	FechaAnulacion time.Time
	Desagregado    *map[string]float64
}
