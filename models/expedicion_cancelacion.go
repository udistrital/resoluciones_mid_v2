package models

import "time"

type ExpedicionCancelacion struct {
	Vinculaciones   []*CancelacionContrato
	IdResolucion    int
	FechaExpedicion time.Time
}
