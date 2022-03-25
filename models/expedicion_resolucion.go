package models

import "time"

type ExpedicionResolucion struct {
	Vinculaciones   *[]ContratoVinculacion
	IdResolucion    int
	FechaExpedicion time.Time
}
