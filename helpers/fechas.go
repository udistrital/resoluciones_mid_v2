package helpers

import "time"

// NormalizarFechaTimezone ajusta la fecha para compensar diferencias de zona horaria
// entre APIs y bases de datos
func NormalizarFechaTimezone(fecha *time.Time) *time.Time {
	if fecha == nil {
		return nil
	}
	// Crear una nueva fecha con la misma fecha pero a las 00:00:00
	t := time.Date(
		fecha.Year(),
		fecha.Month(),
		fecha.Day(),
		0,        // hora
		0,        // minutos
		0,        // segundos
		0,        // nanosegundos
		time.UTC, // usar zona horaria UTC
	)
	return &t
}

// NormalizarFechasResolucion normaliza tanto la fecha de inicio como la fecha fin
// de una resolución para compensar diferencias de zona horaria
func NormalizarFechasResolucion(fechaInicio, fechaFin *time.Time) (*time.Time, *time.Time) {
	return NormalizarFechaTimezone(fechaInicio), NormalizarFechaTimezone(fechaFin)
}
