package helpers

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/astaxie/beego"
	"github.com/udistrital/resoluciones_mid_v2/models"
)

func TestConsultarRequisitosVinculacionContrato(t *testing.T) {
	original := beego.AppConfig.String("UrlPrevinculaciones")
	t.Cleanup(func() { _ = beego.AppConfig.Set("UrlPrevinculaciones", original) })
	datos := models.ObjetoPrevinculaciones{Vigencia: 2026, Docentes: []models.CargaLectiva{
		{DocDocente: "1015402605", Periodo: " 3 "},
		{DocDocente: "79708124", Periodo: "3"},
		{DocDocente: "1015402605", Periodo: "3"},
		{DocDocente: " 79777053 ", Periodo: "3"},
	}}
	calls := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls++
		if r.Method != http.MethodPost || r.URL.RequestURI() != "/previnculacion" {
			t.Errorf("petición inesperada: %s %s", r.Method, r.URL)
		}
		if r.Header.Get("Accept") != "application/json" || r.Header.Get("Content-Type") != "application/json" {
			t.Errorf("headers incorrectos: %v", r.Header)
		}
		if r.Header.Get("Authorization") != "" {
			t.Error("no se esperaba autenticación ODIN")
		}
		var body map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Error(err)
		}
		want := map[string]interface{}{"anio": float64(2026), "periodo": float64(3), "id_usuario": []interface{}{"1015402605", "79708124", "79777053"}}
		if !reflect.DeepEqual(body, want) {
			t.Errorf("payload incorrecto: %#v", body)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`[
 {"estado":"A","registro_tipo_vinculo":"S   ","aprueba_soporte":"S   ","periodo":3,"id_categoria":"2     ","categoria":"Asistente    ","registro_proveedor":"S   ","nocruce":1,"nombre_estado":"Activa","registro_tercero":"S   ","registro_no_cruce":"S   ","id_usuario":"1015402605              ","id":29946,"resolucion_vinculable":"S","anio":2026,"observacion":"","registro_pre_carga":"S   "},
 {"estado":"A","registro_tipo_vinculo":"S   ","aprueba_soporte":"S   ","periodo":3,"id_categoria":"2     ","categoria":"Asistente    ","registro_proveedor":"S   ","nocruce":1,"nombre_estado":"Activa","registro_tercero":"S   ","registro_no_cruce":"S   ","id_usuario":"79708124                ","id":29182,"resolucion_vinculable":"S","anio":2026,"observacion":"","registro_pre_carga":"S   "}
]`))
	}))
	defer server.Close()
	_ = beego.AppConfig.Set("UrlPrevinculaciones", server.URL+"/previnculacion")
	respuesta, errMap := consultarRequisitosVinculacion(context.Background(), datos)
	if errMap != nil || len(respuesta) != 2 || respuesta[0].ResolucionVinculable != "S" {
		t.Fatalf("respuesta incorrecta: %+v, %v", respuesta, errMap)
	}
	if noVinculables := construirDocentesNoVinculables(datos.Docentes, respuesta); len(noVinculables) != 1 || noVinculables[0].Documento != "79777053" {
		t.Fatalf("solo el docente ausente en la respuesta debe bloquearse: %+v", noVinculables)
	}
	for _, requisito := range respuesta {
		if requisito.Anio != 2026 || requisito.Periodo != 3 || requisito.NoCruce != 1 || !banderaAprobada(requisito.ApruebaSoporte) {
			t.Fatalf("campos de respuesta mal interpretados: %+v", requisito)
		}
	}
	if calls != 1 {
		t.Fatalf("se esperaba un único POST; llamadas: %d", calls)
	}

	datos.Docentes[0].Periodo = "inválido"
	if _, errMap = consultarRequisitosVinculacion(context.Background(), datos); errMap == nil || errMap["status"] != "400" {
		t.Fatalf("se esperaba error de período: %v", errMap)
	}
	if calls != 1 {
		t.Fatal("período inválido produjo petición")
	}
	datos.Docentes[0].Periodo = "3"
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if _, errMap = consultarRequisitosVinculacion(ctx, datos); errMap == nil || errMap["status"] != "502" {
		t.Fatalf("se esperaba error por cancelación: %v", errMap)
	}
	_ = beego.AppConfig.Set("UrlPrevinculaciones", "")
	if _, errMap = consultarRequisitosVinculacion(context.Background(), datos); errMap == nil || errMap["status"] != "500" {
		t.Fatalf("se esperaba error de configuración: %v", errMap)
	}
}
