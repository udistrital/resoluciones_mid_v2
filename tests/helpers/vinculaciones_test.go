package helpers_test

import (
	"testing"

	helperspkg "github.com/udistrital/resoluciones_mid_v2/helpers"
	"github.com/udistrital/resoluciones_mid_v2/models"
)

func TestConstruirResumenVinculacion(t *testing.T) {
	numeroContrato := "DVE123"
	previnculacion := models.VinculacionDocente{
		Id:                   10,
		NumeroContrato:       &numeroContrato,
		Vigencia:             2026,
		PersonaId:            123456789,
		NumeroHorasSemanales: 16,
		NumeroSemanas:        18,
		ValorContrato:        2500000,
		Categoria:            " ASOCIADO ",
		NumeroRp:             42,
		ProyectoCurricularId: 7,
		ResolucionVinculacionDocenteId: &models.ResolucionVinculacionDocente{
			Dedicacion: "HCP",
		},
	}
	persona := models.InformacionPersonaNatural{
		NomProveedor:  "Docente Prueba",
		TipoDocumento: &models.ParametroEstandar{ValorParametro: "CC"},
	}
	proyecto := models.Dependencia{Id: 7, Nombre: "INGENIERIA"}
	disponibilidad := models.DisponibilidadVinculacion{Disponibilidad: 999}

	resumen := helperspkg.TestHookConstruirResumenVinculacion(previnculacion, persona, "Bogota", proyecto, disponibilidad)

	if resumen.NumeroContrato != "DVE123" {
		t.Fatalf("numero contrato incorrecto: %+v", resumen)
	}
	if resumen.Categoria != "ASOCIADO" {
		t.Fatalf("categoria incorrecta: %+v", resumen)
	}
	if resumen.ProyectoCurricularNombre != "INGENIERIA" {
		t.Fatalf("proyecto curricular incorrecto: %+v", resumen)
	}
	if resumen.Disponibilidad != 999 {
		t.Fatalf("disponibilidad incorrecta: %+v", resumen)
	}
}

func TestConstruirResumenVinculacionInicializaNumeroContratoNil(t *testing.T) {
	previnculacion := models.VinculacionDocente{
		Id:                   11,
		Vigencia:             2026,
		PersonaId:            987654321,
		NumeroHorasSemanales: 12,
		NumeroSemanas:        16,
		ValorContrato:        1800000,
		Categoria:            "ASISTENTE",
		ProyectoCurricularId: 9,
		ResolucionVinculacionDocenteId: &models.ResolucionVinculacionDocente{
			Dedicacion: "HCH",
		},
	}
	persona := models.InformacionPersonaNatural{
		NomProveedor:  "Docente Sin Contrato",
		TipoDocumento: &models.ParametroEstandar{ValorParametro: "CC"},
	}
	proyecto := models.Dependencia{Id: 9, Nombre: "ARTES"}
	disponibilidad := models.DisponibilidadVinculacion{Disponibilidad: 123}

	resumen := helperspkg.TestHookConstruirResumenVinculacion(previnculacion, persona, "Bogota", proyecto, disponibilidad)

	if resumen.NumeroContrato != "" {
		t.Fatalf("se esperaba numero contrato vacío: %+v", resumen)
	}
}
