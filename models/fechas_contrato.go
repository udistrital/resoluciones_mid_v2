package models

import "time"

type FechasContrato struct {
	FechaInicioReal time.Time // Fecha de inicio real del contrato (ajustada si cae en 31)
	FechaFinReal    time.Time // Fecha de fin real del contrato (ajustada si cae en 31)
	SemanasReales   float64   // Número de semanas reales completas

	FechaInicioPago   time.Time // Fecha de inicio de pago (ajustada)
	FechaFinPago      time.Time // Fecha de fin de pago (ajustada a titan)
	SemanasPagoReales float64   // Número de semanas de pago reales segun fechas de titan
	SemanasPagoDve    float64   // Número de semanas de pago segun liquidador
}
