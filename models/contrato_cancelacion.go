package models

import "time"

type ContratoCancelacion struct {
	NumeroContrato string
	Vigencia       int
	Documento      string
	FechaAnulacion time.Time
	Desagregado    *map[string]float64
}
