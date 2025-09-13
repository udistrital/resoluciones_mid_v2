package models

import "time"

type ContratoPreliquidacion struct {
	NumeroContrato           string
	Vigencia                 int
	NombreCompleto           string
	Documento                string
	PersonaId                int
	TipoNominaId             int
	FechaInicio              time.Time
	FechaFin                 time.Time
	ValorContrato            float64
	DependenciaId            int
	Cdp                      int
	Rp                       int
	Activo                   bool
	NumeroSemanas            int
	ResolucionId             int
	Resolucion               string
	PorcentajesDesagregadoId int
	Desagregado              *map[string]float64
}
