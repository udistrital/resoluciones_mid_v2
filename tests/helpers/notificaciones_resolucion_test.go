package helpers_test

import (
	"testing"

	helperspkg "github.com/udistrital/resoluciones_mid_v2/helpers"
)

func TestPlantillaNotificacionResolucion(t *testing.T) {
	casos := map[string]string{
		"RVIN": "RESOLUCIONES_VINCULACION_PLANTILLA",
		"RCAN": "RESOLUCIONES_CANCELACION_PLANTILLA",
		"RRED": "RESOLUCIONES_REDUCCION_PLANTILLA",
		"RADD": "RESOLUCIONES_ADICION_PLANTILLA",
		"OTRO": "",
	}

	for tipo, esperada := range casos {
		if actual := helperspkg.TestHookPlantillaNotificacionResolucion(tipo); actual != esperada {
			t.Fatalf("plantilla incorrecta para %s: got %q want %q", tipo, actual, esperada)
		}
	}
}
