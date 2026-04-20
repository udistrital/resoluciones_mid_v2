package controllers_test

import (
	"testing"

	controllerpkg "github.com/udistrital/resoluciones_mid_v2/controllers"
	"github.com/udistrital/resoluciones_mid_v2/models"
)

func TestBuildFiltroConsultaGestionResoluciones(t *testing.T) {
	filtro := controllerpkg.TestHookBuildFiltroConsulta(map[string]string{
		"limit":            "10",
		"offset":           "1",
		"NumeroResolucion": "084",
		"Vigencia":         "2026",
		"Periodo":          "1",
		"Semanas":          "16",
		"Facultad":         "FACULTAD TECNOLOGICA",
		"NivelAcademico":   "POSGRADO",
		"Dedicacion":       "HCP",
		"Estado":           "Expedida|Aprobada",
		"TipoResolucion":   "Resolución de Vinculación",
		"ExcluirTipo":      "RVE",
	}, "Vigencia")

	if filtro.Limit != "10" || filtro.Offset != "1" {
		t.Fatalf("paginación incorrecta: %+v", filtro)
	}
	if filtro.NumeroResolucion != "084" || filtro.Vigencia != "2026" {
		t.Fatalf("filtros básicos incorrectos: %+v", filtro)
	}
	if filtro.FacultadId != "FACULTAD TECNOLOGICA" {
		t.Fatalf("facultad incorrecta: %+v", filtro)
	}
	if filtro.TipoResolucion != "Resolución de Vinculación" || filtro.ExcluirTipo != "RVE" {
		t.Fatalf("tipo de resolución incorrecto: %+v", filtro)
	}
}

func TestBuildFiltroConsultaResolucionesPorRolPriorizaAliases(t *testing.T) {
	filtro := controllerpkg.TestHookBuildFiltroConsulta(map[string]string{
		"limit":             "10",
		"offset":            "2",
		"numero_resolucion": "001",
		"vigencia":          "2026",
		"id_oikos":          "12",
	}, "vigencia")

	if filtro.NumeroResolucion != "001" {
		t.Fatalf("numero_resolucion no fue tomado: %+v", filtro)
	}
	if filtro.Vigencia != "2026" {
		t.Fatalf("vigencia no fue tomada: %+v", filtro)
	}
	if filtro.FacultadId != "12" {
		t.Fatalf("id_oikos no fue tomado como facultad: %+v", filtro)
	}
}

func TestBuildFiltroConsultaPrefierePrimerValorNoVacio(t *testing.T) {
	filtro := controllerpkg.TestHookBuildFiltroConsulta(map[string]string{
		"NumeroResolucion":  "094",
		"numero_resolucion": "001",
		"Vigencia":          "",
		"vigencia":          "2027",
		"Facultad":          "FACULTAD DE INGENIERIA",
		"id_oikos":          "33",
	}, "Vigencia", "vigencia")

	if filtro.NumeroResolucion != "094" {
		t.Fatalf("debía priorizar NumeroResolucion: %+v", filtro)
	}
	if filtro.Vigencia != "2027" {
		t.Fatalf("debía tomar vigencia en minúscula: %+v", filtro)
	}
	if filtro.FacultadId != "FACULTAD DE INGENIERIA" {
		t.Fatalf("debía priorizar Facultad textual: %+v", filtro)
	}
}

func TestValidateFiltroConsultaAceptaValoresValidos(t *testing.T) {
	filtro := models.Filtro{
		Limit:            "10",
		Offset:           "1",
		NumeroResolucion: "084",
		Vigencia:         "2026",
		Periodo:          "1",
		Semanas:          "16",
	}

	if err := controllerpkg.TestHookValidateFiltroConsulta(filtro); err != nil {
		t.Fatalf("se esperaba filtro válido y se obtuvo error: %v", err)
	}
}

func TestValidateFiltroConsultaRechazaValoresInvalidos(t *testing.T) {
	casos := []models.Filtro{
		{Limit: "0", Offset: "1"},
		{Limit: "10", Offset: "-1"},
		{Limit: "10", Offset: "1", NumeroResolucion: "ABC"},
		{Limit: "10", Offset: "1", Vigencia: "20A6"},
		{Limit: "10", Offset: "1", Periodo: "uno"},
		{Limit: "10", Offset: "1", Semanas: "dieciseis"},
	}

	for _, caso := range casos {
		if err := controllerpkg.TestHookValidateFiltroConsulta(caso); err == nil {
			t.Fatalf("se esperaba error para filtro inválido: %+v", caso)
		}
	}
}

func TestValidateRequiredHelpers(t *testing.T) {
	if err := controllerpkg.TestHookValidateRequiredText("1023026828", "numero_documento es requerido"); err != nil {
		t.Fatalf("no debía fallar validateRequiredText con dato válido: %v", err)
	}
	if err := controllerpkg.TestHookValidateRequiredRoles([]string{"DECANO"}); err != nil {
		t.Fatalf("no debía fallar validateRequiredRoles con dato válido: %v", err)
	}
	if err := controllerpkg.TestHookValidateNamedPositiveInt("2026", "vigencia"); err != nil {
		t.Fatalf("no debía fallar validateNamedPositiveInt con dato válido: %v", err)
	}
}

func TestValidateRequiredHelpersRechazaInvalidos(t *testing.T) {
	if err := controllerpkg.TestHookValidateRequiredText("", "numero_documento es requerido"); err == nil {
		t.Fatal("se esperaba error por numero_documento vacío")
	}
	if err := controllerpkg.TestHookValidateRequiredRoles(nil); err == nil {
		t.Fatal("se esperaba error por roles vacíos")
	}
	if err := controllerpkg.TestHookValidateNamedPositiveInt("0", "vigencia"); err == nil {
		t.Fatal("se esperaba error por vigencia inválida")
	}
}
