package helpers

import (
	"time"

	"github.com/udistrital/resoluciones_mid_v2/models"
)

func calcularSemanasContratoDVE(fechaInicio time.Time, fechaFin time.Time) (semanas float64) {
	var a, m, d int
	var mesesContrato float64
	if fechaFin.IsZero() {
		fechaFinActual := time.Now()
		a, m, d = diff(fechaInicio, fechaFinActual)
		mesesContrato = (float64(a * 12)) + float64(m) + (float64(d) / 30)
	} else {
		a, m, d = diff(fechaInicio, fechaFin)
		d += 1
		if d == 22 {
			d += 1
		}
		mesesContrato = (float64(a * 12)) + float64(m) + (float64(d) / 30)
	}

	if mesesContrato/float64(int(mesesContrato)) != 1 {
		return (mesesContrato * 4) + 1
	}

	return mesesContrato * 4
}

// CalcularFechasContrato calcula la fecha de fin de un contrato a partir de la fecha de inicio y el numero de semanas.
// Calcula las fechas reales del contrato y las fechas ajustadas para TITAN.
// Ajusta la fecha inicio y fin para que inicie un lunes.
func CalcularFechasContrato(fechaInicio time.Time, numeroSemanas int) models.FechasContrato {
	var resultado models.FechasContrato

	resultado.FechaInicioReal = fechaInicio
	dias := numeroSemanas * 7
	resultado.FechaFinReal = fechaInicio.AddDate(0, 0, dias)
	resultado.SemanasReales = float64(dias) / 7
	resultado.FechaInicioPago = fechaInicio

	if resultado.FechaInicioPago.Weekday() != time.Monday {
		diasHastaLunes := int(time.Monday - resultado.FechaInicioPago.Weekday())
		if diasHastaLunes <= 0 {
			diasHastaLunes += 7
		}
		resultado.FechaInicioPago = resultado.FechaInicioPago.AddDate(0, 0, diasHastaLunes)
	}

	diasPago := numeroSemanas * 7
	if diasPago == 0 {
		resultado.FechaFinPago = resultado.FechaInicioPago
	} else {
		resultado.FechaFinPago = resultado.FechaInicioPago.AddDate(0, 0, diasPago-1)
	}

	if resultado.FechaFinPago.Day() == 31 {
		resultado.FechaFinPago = resultado.FechaFinPago.AddDate(0, 0, -1)
	}

	diasDiferencia := resultado.FechaFinPago.Sub(resultado.FechaInicioPago).Hours() / 24
	resultado.SemanasPagoReales = (diasDiferencia + 1) / 7
	resultado.SemanasPagoDve = calcularSemanasContratoDVE(resultado.FechaInicioPago, resultado.FechaFinPago)

	return resultado
}
