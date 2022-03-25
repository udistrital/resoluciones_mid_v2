package models

import "time"

type ObjetoModificaciones struct {
	CambiosVinculacion       *CambioVinculacion
	ResolucionNuevaId        *ResolucionVinculacionDocente
	ModificacionResolucionId int
}

type ObjetoCancelaciones struct {
	CambiosVinculacion       []CambioVinculacion
	ResolucionNuevaId        *ResolucionVinculacionDocente
	ModificacionResolucionId int
}

type CambioVinculacion struct {
	NumeroHorasSemanales int
	NumeroSemanas        int
	FechaInicio          time.Time
	DocPresupuestal      *DocumentoPresupuestal
	VinculacionOriginal  *Vinculaciones
}
