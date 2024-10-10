package models

import "time"

type VinculacionDocente struct {
	Id                             int
	NumeroContrato                 *string
	Vigencia                       int
	PersonaId                      float64
	NumeroHorasSemanales           int
	NumeroSemanas                  int
	ValorPuntoSalarial             float64
	SalarioMinimoId                int
	ResolucionVinculacionDocenteId *ResolucionVinculacionDocente
	DedicacionId                   int
	ProyectoCurricularId           int
	ValorContrato                  float64
	Categoria                      string
	Emerito                        bool
	DependenciaAcademica           int
	NumeroRp                       float64
	VigenciaRp                     float64
	FechaInicio                    time.Time
	Activo                         bool
	FechaCreacion                  string
	FechaModificacion              string
	NumeroHorasTrabajadas          int
}
