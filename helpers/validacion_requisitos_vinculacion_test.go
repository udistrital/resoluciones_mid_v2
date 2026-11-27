package helpers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	"github.com/astaxie/beego"
	"github.com/udistrital/resoluciones_mid_v2/models"
)

func TestConsultarRequisitosVinculacionContrato(t *testing.T) {
	configurarPrevinculacionParaPrueba(t, "", "")
	datos := models.ObjetoPrevinculaciones{Vigencia: 2026, Docentes: []models.CargaLectiva{
		{DocDocente: "1015402605", Periodo: " 3 "},
		{DocDocente: "79708124", Periodo: "3"},
		{DocDocente: "1015402605", Periodo: "3"},
		{DocDocente: " 79777053 ", Periodo: "3"},
	}}
	calls := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls++
		if r.Method != http.MethodPost || r.URL.RequestURI() != "/busservicios/urano_produccion/previnculacion" {
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
	configurarPrevinculacionParaPrueba(t, server.URL+"/busservicios", "/urano_produccion")
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

	_ = beego.AppConfig.Set("UrlcrudWSO2", server.URL+"/busservicios/")
	_ = beego.AppConfig.Set("UrlPrevinculaciones", "/urano_produccion/")
	if _, errMap := consultarRequisitosVinculacion(context.Background(), datos); errMap != nil {
		t.Fatalf("URL base con barra final produjo error: %v", errMap)
	}
	if calls != 2 {
		t.Fatalf("se esperaba un POST adicional con barra final; llamadas: %d", calls)
	}

	datos.Docentes[0].Periodo = "inválido"
	if _, errMap = consultarRequisitosVinculacion(context.Background(), datos); errMap == nil || errMap["status"] != "400" {
		t.Fatalf("se esperaba error de período: %v", errMap)
	}
	if calls != 2 {
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

func configurarPrevinculacionParaPrueba(t *testing.T, baseURL, namespace string) {
	t.Helper()
	for key, value := range map[string]string{"UrlcrudWSO2": baseURL, "UrlPrevinculaciones": namespace} {
		original := beego.AppConfig.String(key)
		t.Cleanup(func() { _ = beego.AppConfig.Set(key, original) })
		if err := beego.AppConfig.Set(key, value); err != nil {
			t.Fatal(err)
		}
	}
}

func TestConstruirURLPrevinculacionBarras(t *testing.T) {
	for _, scheme := range []string{"http", "https"} {
		for _, baseSuffix := range []string{"", "/", "///"} {
			for _, namespace := range []string{"urano_produccion", "/urano_produccion", "urano_produccion/", "/urano_produccion/", "///urano_produccion///"} {
				baseURL := scheme + "://api.example.com/busservicios" + baseSuffix
				t.Run(baseURL+"|"+namespace, func(t *testing.T) {
					got, err := construirURLPrevinculacion(baseURL, namespace)
					want := scheme + "://api.example.com/busservicios/urano_produccion/previnculacion"
					if err != nil || got != want {
						t.Fatalf("URL incorrecta: got=%q, want=%q, err=%v", got, want, err)
					}
				})
			}
		}
	}
	for _, tc := range []struct{ name, baseURL, namespace, want string }{
		{"sin ruta base", "https://api.example.com", "/urano_produccion", "https://api.example.com/urano_produccion/previnculacion"},
		{"espacios externos", " https://api.example.com/busservicios/ ", " /urano_produccion/ ", "https://api.example.com/busservicios/urano_produccion/previnculacion"},
		{"puerto local", "http://localhost:8080/busservicios", "/urano_produccion", "http://localhost:8080/busservicios/urano_produccion/previnculacion"},
		{"IPv6", "http://[::1]:8080/busservicios", "/urano_produccion", "http://[::1]:8080/busservicios/urano_produccion/previnculacion"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got, err := construirURLPrevinculacion(tc.baseURL, tc.namespace)
			if err != nil || got != tc.want {
				t.Fatalf("URL incorrecta: got=%q, want=%q, err=%v", got, tc.want, err)
			}
		})
	}
}

func TestConsultarRequisitosVinculacionConfiguracionInvalida(t *testing.T) {
	calls := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls++
		_, _ = w.Write([]byte(`[]`))
	}))
	defer server.Close()
	datos := models.ObjetoPrevinculaciones{Vigencia: 2026, Docentes: []models.CargaLectiva{{DocDocente: "1015402605", Periodo: "3"}}}
	for _, tc := range []struct{ name, baseURL, namespace, key string }{
		{"base vacía", "", "/urano_produccion", "RESOLUCIONES_MID_V2_WSO2_URL"},
		{"base en blanco", "  ", "/urano_produccion", "RESOLUCIONES_MID_V2_WSO2_URL"},
		{"namespace vacío", server.URL, "", "RESOLUCIONES_MID_V2_PREVINCULACIONES"},
		{"namespace en blanco", server.URL, "  ", "RESOLUCIONES_MID_V2_PREVINCULACIONES"},
		{"namespace sin segmento", server.URL, "///", "RESOLUCIONES_MID_V2_PREVINCULACIONES"},
		{"ambas vacías", "", "", "RESOLUCIONES_MID_V2_WSO2_URL"},
		{"base relativa", "/busservicios", "/urano_produccion", "RESOLUCIONES_MID_V2_WSO2_URL"},
		{"sin protocolo", "api.example.com/busservicios", "/urano_produccion", "RESOLUCIONES_MID_V2_WSO2_URL"},
		{"protocolo relativo", "//api.example.com/busservicios", "/urano_produccion", "RESOLUCIONES_MID_V2_WSO2_URL"},
		{"protocolo no HTTP", "ftp://api.example.com/busservicios", "/urano_produccion", "RESOLUCIONES_MID_V2_WSO2_URL"},
		{"sin dominio", "https:///busservicios", "/urano_produccion", "RESOLUCIONES_MID_V2_WSO2_URL"},
		{"dominio vacío con puerto", "http://:8080/busservicios", "/urano_produccion", "RESOLUCIONES_MID_V2_WSO2_URL"},
		{"URL opaca", "https:busservicios", "/urano_produccion", "RESOLUCIONES_MID_V2_WSO2_URL"},
		{"dominio con espacios", "https://api example.com", "/urano_produccion", "RESOLUCIONES_MID_V2_WSO2_URL"},
		{"puerto inválido", "https://api.example.com:abc", "/urano_produccion", "RESOLUCIONES_MID_V2_WSO2_URL"},
		{"puerto fuera de rango", "https://api.example.com:65536", "/urano_produccion", "RESOLUCIONES_MID_V2_WSO2_URL"},
		{"puerto cero", "https://api.example.com:0", "/urano_produccion", "RESOLUCIONES_MID_V2_WSO2_URL"},
		{"escape inválido", "https://api.example.com/%zz", "/urano_produccion", "RESOLUCIONES_MID_V2_WSO2_URL"},
		{"consulta en base", server.URL + "?q=1", "/urano_produccion", "RESOLUCIONES_MID_V2_WSO2_URL"},
		{"fragmento en base", server.URL + "#ruta", "/urano_produccion", "RESOLUCIONES_MID_V2_WSO2_URL"},
		{"credenciales en base", "https://user:pass@api.example.com", "/urano_produccion", "RESOLUCIONES_MID_V2_WSO2_URL"},
		{"namespace absoluto", server.URL, "https://api.example.com/urano_produccion", "RESOLUCIONES_MID_V2_PREVINCULACIONES"},
		{"namespace con consulta", server.URL, "/urano_produccion?q=1", "RESOLUCIONES_MID_V2_PREVINCULACIONES"},
		{"namespace con fragmento", server.URL, "/urano_produccion#ruta", "RESOLUCIONES_MID_V2_PREVINCULACIONES"},
		{"namespace con escape inválido", server.URL, "/urano_%zz", "RESOLUCIONES_MID_V2_PREVINCULACIONES"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			configurarPrevinculacionParaPrueba(t, tc.baseURL, tc.namespace)
			respuesta, errMap := consultarRequisitosVinculacion(context.Background(), datos)
			if respuesta != nil || errMap == nil || errMap["status"] != "500" {
				t.Fatalf("se esperaba error 500 de configuración: respuesta=%+v, error=%v", respuesta, errMap)
			}
			if !strings.Contains(fmt.Sprint(errMap["err"]), tc.key) {
				t.Fatalf("error debe identificar configuración inválida %s: %v", tc.key, errMap)
			}
			if calls != 0 {
				t.Fatal("configuración inválida produjo petición HTTP")
			}
		})
	}
}

func TestConsultarRequisitosVinculacionErroresComunicacion(t *testing.T) {
	datos := models.ObjetoPrevinculaciones{Vigencia: 2026, Docentes: []models.CargaLectiva{{DocDocente: "1015402605", Periodo: "3"}}}
	t.Run("servicio no disponible", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusServiceUnavailable)
		}))
		defer server.Close()
		configurarPrevinculacionParaPrueba(t, server.URL, "/urano_produccion")
		if _, errMap := consultarRequisitosVinculacion(context.Background(), datos); errMap == nil || errMap["status"] != "502" {
			t.Fatalf("se esperaba 502 del servicio externo: %v", errMap)
		}
	})
	t.Run("conexión rechazada", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
		server.Close()
		configurarPrevinculacionParaPrueba(t, server.URL, "/urano_produccion")
		if _, errMap := consultarRequisitosVinculacion(context.Background(), datos); errMap == nil || errMap["status"] != "502" {
			t.Fatalf("se esperaba 502 de conexión: %v", errMap)
		}
	})
}
