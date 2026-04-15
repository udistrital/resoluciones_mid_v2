package services_test

import (
	"testing"

	"github.com/udistrital/resoluciones_mid_v2/models"
	servicepkg "github.com/udistrital/resoluciones_mid_v2/services"
)

func TestNormalizeRol(t *testing.T) {
	if got := servicepkg.TestHookNormalizeRol(" decano "); got != "DECANO" {
		t.Fatalf("rol normalizado incorrecto: got %q want %q", got, "DECANO")
	}
}

func TestJoinWSO2URL(t *testing.T) {
	got := servicepkg.TestHookJoinWSO2URL("https://", "/pruebasapi.intranetoas.udistrital.edu.co:8104/", "/academica_crud_api/", "/decano/1023")
	want := "https://pruebasapi.intranetoas.udistrital.edu.co:8104/academica_crud_api/decano/1023"

	if got != want {
		t.Fatalf("url construida incorrecta: got %q want %q", got, want)
	}
}

func TestGetHighestPriorityRol(t *testing.T) {
	got := servicepkg.TestHookGetHighestPriorityRol([]string{"docente", "decano", "ADMINISTRADOR_RESOLUCIONES"})
	if got != "ADMINISTRADOR_RESOLUCIONES" {
		t.Fatalf("rol prioritario incorrecto: got %q want %q", got, "ADMINISTRADOR_RESOLUCIONES")
	}
}

func TestIsGlobalRol(t *testing.T) {
	if !servicepkg.TestHookIsGlobalRol("asis_financiera") {
		t.Fatal("ASIS_FINANCIERA debía ser global")
	}
	if servicepkg.TestHookIsGlobalRol("decano") {
		t.Fatal("DECANO no debía ser global")
	}
}

func TestDeduplicateDependencias(t *testing.T) {
	items := []models.DependenciaUsuario{
		{CodigoDependencia: 10, IdOikos: 101, Nombre: "A"},
		{CodigoDependencia: 10, IdOikos: 101, Nombre: "A repetida"},
		{CodigoDependencia: 20, IdOikos: 202, Nombre: "B"},
	}

	resultado := servicepkg.TestHookDeduplicateDependencias(items)
	if len(resultado) != 2 {
		t.Fatalf("cantidad de dependencias deduplicadas incorrecta: got %d want %d", len(resultado), 2)
	}
	if resultado[0].IdOikos != 101 || resultado[1].IdOikos != 202 {
		t.Fatalf("orden o ids inesperados tras deduplicar: %+v", resultado)
	}
}
