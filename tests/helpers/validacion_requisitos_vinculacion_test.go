package helpers_test

import (
	"reflect"
	"testing"

	helperspkg "github.com/udistrital/resoluciones_mid_v2/helpers"
	"github.com/udistrital/resoluciones_mid_v2/models"
)

func TestConstruirIdUsuariosOdin(t *testing.T) {
	docentes := []models.CargaLectiva{
		{DocDocente: "1015402605"},
		{DocDocente: "79708124"},
		{DocDocente: "1015402605"},
		{DocDocente: " 79777053 "},
	}

	got := helperspkg.TestHookConstruirIdUsuariosOdin(docentes)
	want := "'1015402605','79708124','79777053'"

	if got != want {
		t.Fatalf("id_usuario incorrecto. got=%q want=%q", got, want)
	}
}

func TestConstruirMotivosNoVinculable(t *testing.T) {
	requisito := models.OdinRequisitoVinculacion{
		ApruebaSoporte:      "N   ",
		RegistroTercero:     "S   ",
		RegistroProveedor:   "N   ",
		RegistroTipoVinculo: "N   ",
		RegistroNoCruce:     "S   ",
		RegistroPreCarga:    "N   ",
		Observacion:         "Pendiente actualización en proveedor",
	}

	got := helperspkg.TestHookConstruirMotivosNoVinculable(requisito)
	want := []string{
		"soporte no aprobado",
		"sin registro de proveedor",
		"sin registro de tipo de vínculo",
		"sin precarga",
		"Pendiente actualización en proveedor",
	}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("motivos incorrectos. got=%v want=%v", got, want)
	}
}

func TestConstruirDocentesNoVinculables(t *testing.T) {
	docentes := []models.CargaLectiva{
		{DocDocente: "1015402605", DocenteApellido: "Laguna", DocenteNombre: "Cristian"},
		{DocDocente: "79708124", DocenteApellido: "Calvo", DocenteNombre: "Martha"},
		{DocDocente: "555", DocenteApellido: "Sin", DocenteNombre: "Respuesta"},
	}
	respuesta := []models.OdinRequisitoVinculacion{
		{
			IdUsuario:            "1015402605              ",
			ApruebaSoporte:       "N   ",
			RegistroTercero:      "S",
			RegistroProveedor:    "S",
			RegistroTipoVinculo:  "S",
			RegistroNoCruce:      "S",
			RegistroPreCarga:     "S",
			ResolucionVinculable: "N",
		},
		{
			IdUsuario:            "79708124",
			ResolucionVinculable: "S",
		},
	}

	got := helperspkg.TestHookConstruirDocentesNoVinculables(docentes, respuesta)
	if len(got) != 2 {
		t.Fatalf("se esperaban 2 docentes no vinculables y llegaron %d: %+v", len(got), got)
	}
	if got[0].Documento != "1015402605" || got[0].Nombre != "Laguna Cristian" {
		t.Fatalf("primer docente incorrecto: %+v", got[0])
	}
	if len(got[0].Motivos) != 1 || got[0].Motivos[0] != "soporte no aprobado" {
		t.Fatalf("motivos primer docente incorrectos: %+v", got[0])
	}
	if got[1].Documento != "555" || got[1].Nombre != "Sin Respuesta" {
		t.Fatalf("segundo docente incorrecto: %+v", got[1])
	}
	if len(got[1].Motivos) != 1 || got[1].Motivos[0] != "No se encontraron datos en el proceso Cargue de Soportes Previnculación para el docente" {
		t.Fatalf("motivos segundo docente incorrectos: %+v", got[1])
	}
}

func TestNormalizarRespuestaRequisitosODINSinDatos(t *testing.T) {
	raw := map[string]interface{}{
		"Error": "No se encontrarón datos asociados al reporte",
	}

	got, errMap := helperspkg.TestHookNormalizarRespuestaRequisitosODIN(raw)
	if errMap != nil {
		t.Fatalf("no se esperaba error y llegó: %+v", errMap)
	}
	if len(got) != 0 {
		t.Fatalf("se esperaba respuesta vacía y llegó: %+v", got)
	}
}
