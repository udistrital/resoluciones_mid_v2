package models

import (
	"time"
)

type ContratoGeneral struct {
	Id                           string
	VigenciaContrato             int
	ObjetoContrato               string
	PlazoEjecucion               int
	FormaPago                    *Parametros
	OrdenadorGasto               int
	ClausulaRegistroPresupuestal bool
	SedeSolicitante              string
	DependenciaSolicitante       string
	Contratista                  int
	ValorContrato                float64
	Justificacion                string
	DescripcionFormaPago         string
	Condiciones                  string
	FechaRegistro                time.Time
	TipologiaContrato            int
	TipoCompromiso               int
	ModalidadSeleccion           int
	Procedimiento                int
	RegimenContratacion          int
	TipoGasto                    int
	TemaGastoInversion           int
	OrigenPresupueso             int
	OrigenRecursos               int
	TipoMoneda                   int
	ValorContratoMe              float64
	ValorTasaCambio              float64
	TipoControl                  int
	Observaciones                string
	Supervisor                   *SupervisorContrato
	ClaseContratista             int
	Convenio                     string
	NumeroConstancia             int
	Estado                       bool
	TipoContrato                 *TipoContrato
	LugarEjecucion               *LugarEjecucion
	UnidadEjecucion              *Parametros
	UnidadEjecutora              int
}
