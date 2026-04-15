package helpers

import "time"

var meses = map[string][]string{
	"es": {"Enero", "Febrero", "Marzo", "Abril", "Mayo", "Junio", "Julio", "Agosto", "Septiembre", "Octubre", "Noviembre", "Diciembre"},
}

var periodo = map[string][]string{
	"es": {"Primer", "Segundo", "Tercer"},
}

func obtenerNombreMes(m time.Month, idioma string) string {
	nombresMeses, ok := meses[idioma]
	if !ok {
		return ""
	}

	mes := int(m) - 1
	if mes < 0 || mes >= len(nombresMeses) {
		return ""
	}

	return nombresMeses[mes]
}

func cambiarString(original string) (cambiado string) {
	switch {
	case original == "HCH":
		cambiado = "Hora Cátedra Honorarios"
	case original == "HCP":
		cambiado = "Hora Cátedra Salarios"
	case original == "TCO-MTO":
		cambiado = "Tiempo Completo Ocasional - Medio Tiempo Ocasional"
	default:
		cambiado = original
	}
	return cambiado
}

func obtenerPeriodo(periodoA int, idioma string) string {
	periodoV, ok := periodo[idioma]
	if !ok {
		return ""
	}

	periodoF := periodoA - 1
	if periodoF < 0 || periodoF >= len(periodoV) {
		return ""
	}

	return periodoV[periodoF]
}
