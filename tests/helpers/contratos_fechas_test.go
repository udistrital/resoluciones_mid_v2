package helpers_test

import (
	"testing"
	"time"

	helperspkg "github.com/udistrital/resoluciones_mid_v2/helpers"
)

func TestCalcularFechasContratoAjustaInicioAlSiguienteLunes(t *testing.T) {
	fechaInicio := time.Date(2026, time.April, 15, 0, 0, 0, 0, time.UTC)

	resultado := helperspkg.CalcularFechasContrato(fechaInicio, 2)

	esperadaInicioPago := time.Date(2026, time.April, 20, 0, 0, 0, 0, time.UTC)
	esperadaFinPago := time.Date(2026, time.May, 3, 0, 0, 0, 0, time.UTC)

	if !resultado.FechaInicioPago.Equal(esperadaInicioPago) {
		t.Fatalf("fecha inicio pago incorrecta: got %v want %v", resultado.FechaInicioPago, esperadaInicioPago)
	}
	if !resultado.FechaFinPago.Equal(esperadaFinPago) {
		t.Fatalf("fecha fin pago incorrecta: got %v want %v", resultado.FechaFinPago, esperadaFinPago)
	}
	if resultado.SemanasPagoReales != 2 {
		t.Fatalf("semanas pago reales incorrectas: got %v want 2", resultado.SemanasPagoReales)
	}
}

func TestCalcularFechasContratoAjustaFinSiCaeEn31(t *testing.T) {
	fechaInicio := time.Date(2026, time.March, 16, 0, 0, 0, 0, time.UTC)

	resultado := helperspkg.CalcularFechasContrato(fechaInicio, 3)

	esperadaFinPago := time.Date(2026, time.April, 5, 0, 0, 0, 0, time.UTC)
	if !resultado.FechaFinPago.Equal(esperadaFinPago) {
		t.Fatalf("fecha fin pago incorrecta: got %v want %v", resultado.FechaFinPago, esperadaFinPago)
	}
}
