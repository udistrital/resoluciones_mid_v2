package services_test

import (
	"errors"
	"testing"

	servicepkg "github.com/udistrital/resoluciones_mid_v2/services"
)

func TestCoerceIDString(t *testing.T) {
	tests := []struct {
		name  string
		value interface{}
		want  string
		ok    bool
	}{
		{name: "float64", value: float64(180), want: "180", ok: true},
		{name: "int", value: 180, want: "180", ok: true},
		{name: "string", value: " 180 ", want: "180", ok: true},
		{name: "empty string", value: " ", want: "", ok: false},
		{name: "nil", value: nil, want: "", ok: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := servicepkg.TestHookCoerceIDString(tt.value)
			if got != tt.want || ok != tt.ok {
				t.Fatalf("got (%q, %t), want (%q, %t)", got, ok, tt.want, tt.ok)
			}
		})
	}
}

func TestSafePanicToError(t *testing.T) {
	err := servicepkg.TestHookSafePanicToError(func() {
		panic(errors.New("crud no disponible"))
	})
	if err == nil || err.Error() != "crud no disponible" {
		t.Fatalf("error incorrecto: %v", err)
	}

	err = servicepkg.TestHookSafePanicToError(func() {
		panic("respuesta inválida")
	})
	if err == nil || err.Error() != "respuesta inválida" {
		t.Fatalf("error incorrecto: %v", err)
	}
}

func TestConstruirRegistrosRpIgnoraFilasVacias(t *testing.T) {
	headers := map[string]int{
		"cod_resolucion": 0,
		"cod_facultad":   1,
		"documento":      2,
		"cod_proyecto":   3,
		"crp":            4,
	}
	rows := [][]string{
		{"'0180'", "65", "79561412", "37", "5691"},
		{"", "", "", "", ""},
	}

	registros := servicepkg.TestHookConstruirRegistrosRp(rows, headers)
	if len(registros) != 1 {
		t.Fatalf("se esperaba una sola fila útil, got %+v", registros)
	}
	if registros[0].FilaExcel != 2 {
		t.Fatalf("fila excel incorrecta: %+v", registros[0])
	}
}
