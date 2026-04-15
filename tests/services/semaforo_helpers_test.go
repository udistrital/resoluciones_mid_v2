package services_test

import (
	"testing"

	"github.com/udistrital/resoluciones_mid_v2/models"
	servicepkg "github.com/udistrital/resoluciones_mid_v2/services"
)

func TestConstruirMapaContratosTitan(t *testing.T) {
	contratos := []models.ContratoTitan{
		{NumeroContrato: " CPS-001 ", Activo: true},
		{NumeroContrato: "CPS-001", Activo: true},
		{NumeroContrato: "CPS-002", Activo: false},
		{NumeroContrato: "", Activo: true},
	}

	mapa := servicepkg.TestHookConstruirMapaContratosTitan(contratos)
	if len(mapa) != 1 {
		t.Fatalf("mapa de Titan incorrecto: %+v", mapa)
	}
	if !mapa["CPS-001"] {
		t.Fatalf("se esperaba contrato activo CPS-001: %+v", mapa)
	}
}

func TestResumirEstadoRp(t *testing.T) {
	numeroContrato := " CPS-777 "
	vinculacion := models.VinculacionDocente{
		NumeroContrato: &numeroContrato,
		NumeroRp:       123,
		VigenciaRp:     2026,
	}

	resumen := servicepkg.TestHookResumirEstadoRp(vinculacion)
	if resumen.NumeroContrato != "CPS-777" {
		t.Fatalf("numero de contrato incorrecto: got %q want %q", resumen.NumeroContrato, "CPS-777")
	}
	if !resumen.TieneRpResoluciones {
		t.Fatal("debía marcar RP cargado en resoluciones")
	}
}

func TestClasificarEstadoSemaforoVinculacion(t *testing.T) {
	numeroContrato := "CPS-888"
	vinculacion := models.VinculacionDocente{
		Id:             10,
		PersonaId:      1023,
		Vigencia:       2026,
		NumeroContrato: &numeroContrato,
		NumeroRp:       500,
		VigenciaRp:     2026,
	}

	estadoPendiente := servicepkg.TestHookClasificarEstadoSemaforoVinculacion(vinculacion, map[string]bool{})
	if estadoPendiente.EstadoCodigo != "PENDIENTE_TITAN" {
		t.Fatalf("estado pendiente incorrecto: %+v", estadoPendiente)
	}

	estadoCompleto := servicepkg.TestHookClasificarEstadoSemaforoVinculacion(vinculacion, map[string]bool{"CPS-888": true})
	if estadoCompleto.EstadoCodigo != "COMPLETO" || !estadoCompleto.TieneRpTitan {
		t.Fatalf("estado completo incorrecto: %+v", estadoCompleto)
	}
}
