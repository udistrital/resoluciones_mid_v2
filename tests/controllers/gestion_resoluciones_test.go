package controllers

import (
	"net/http"
	"testing"
)

func TestResolucionesPost(t *testing.T) {
	if response, err := http.Post("http://localhost:8529/v1/gestion_resoluciones/", "", nil); err == nil {
		if response.StatusCode != 400 {
			t.Error("Error TestResolucionesPost: Se esperaba 400 y se obtuvo", response.StatusCode)
			t.Fail()
		} else {
			t.Log("TestResolucionesPost  Finalizado Correctamente (OK)")
		}
	} else {
		t.Error("Error TestResolucionesPost:", err.Error())
		t.Fail()
	}
}

func TestResolucionesGetOne(t *testing.T) {
	if response, err := http.Get("http://localhost:8529/v1/gestion_resoluciones/219"); err == nil {
		if response.StatusCode != 200 {
			t.Error("Error TestResolucionesGetOne: Se esperaba 200 y se obtuvo", response.StatusCode)
			t.Fail()
		} else {
			t.Log("TestResolucionesGetOne Finalizado Correctamente (OK)")
		}
	} else {
		t.Error("Error TestResolucionesGetOne:", err.Error())
		t.Fail()
	}
}

func TestResolucionesGetAll(t *testing.T) {
	if response, err := http.Get("http://localhost:8529/v1/gestion_resoluciones?limit=10&offset=0"); err == nil {
		if response.StatusCode != 200 {
			t.Error("Error TestResolucionesGetAll: Se esperaba 200 y se obtuvo", response.StatusCode)
			t.Fail()
		} else {
			t.Log("TestResolucionesGetAll Finalizado Correctamente (OK)")
		}
	} else {
		t.Error("Error TestResolucionesGetAll:", err.Error())
		t.Fail()
	}
}

func TestResolucionesConsultaDocente(t *testing.T) {
	if response, err := http.Get("http://localhost:8529/v1/gestion_resoluciones/consultar_docente/79777053"); err == nil {
		if response.StatusCode != 200 {
			t.Error("Error TestResolucionesConsultaDocente: Se esperaba 200 y se obtuvo", response.StatusCode)
			t.Fail()
		} else {
			t.Log("TestResolucionesConsultaDocente Finalizado Correctamente (OK)")
		}
	} else {
		t.Error("Error TestResolucionesConsultaDocente:", err.Error())
		t.Fail()
	}
}

func TestGetResolucionesInscritas(t *testing.T) {
	if response, err := http.Get("http://localhost:8529/v1/gestion_resoluciones/resoluciones_inscritas?limit=10&offset=0"); err == nil {
		if response.StatusCode != 200 {
			t.Error("Error TestGetResolucionesInscritas: Se esperaba 200 y se obtuvo", response.StatusCode)
			t.Fail()
		} else {
			t.Log("TestGetResolucionesInscritas Finalizado Correctamente (OK)")
		}
	} else {
		t.Error("Error TestGetResolucionesInscritas:", err.Error())
		t.Fail()
	}
}

func TestGetResolucionesAprobadas(t *testing.T) {
	if response, err := http.Get("http://localhost:8529/v1/gestion_resoluciones/resoluciones_aprobadas?limit=10&offset=0"); err == nil {
		if response.StatusCode != 200 {
			t.Error("Error TestGetResolucionesAprobadas: Se esperaba 200 y se obtuvo", response.StatusCode)
			t.Fail()
		} else {
			t.Log("TestGetResolucionesAprobadas Finalizado Correctamente (OK)")
		}
	} else {
		t.Error("Error TestGetResolucionesAprobadas:", err.Error())
		t.Fail()
	}
}

func TestResolucionesGetExpedidas(t *testing.T) {
	if response, err := http.Get("http://localhost:8529/v1/gestion_resoluciones/resoluciones_expedidas?limit=10&offset=0"); err == nil {
		if response.StatusCode != 200 {
			t.Error("Error TestResolucionesGetExpedidas: Se esperaba 200 y se obtuvo", response.StatusCode)
			t.Fail()
		} else {
			t.Log("TestResolucionesGetExpedidas Finalizado Correctamente (OK)")
		}
	} else {
		t.Error("Error TestResolucionesGetExpedidas:", err.Error())
		t.Fail()
	}
}

func TestResolucionesGenerarResolucion(t *testing.T) {
	if response, err := http.Get("http://localhost:8529/v1/gestion_resoluciones/generar_resolucion/219"); err == nil {
		if response.StatusCode != 200 {
			t.Error("Error TestResolucionesGenerarResolucion: Se esperaba 200 y se obtuvo", response.StatusCode)
			t.Fail()
		} else {
			t.Log("TestResolucionesGenerarResolucion Finalizado Correctamente (OK)")
		}
	} else {
		t.Error("Error TestResolucionesGenerarResolucion:", err.Error())
		t.Fail()
	}
}
