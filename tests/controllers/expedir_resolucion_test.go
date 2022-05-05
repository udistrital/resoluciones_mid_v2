package controllers

import (
	"net/http"
	"testing"
)

func TestExpedirPost(t *testing.T) {
	if response, err := http.Post("http://localhost:8529/v1/expedir_resolucion/expedir", "", nil); err == nil {
		if response.StatusCode != 400 {
			t.Error("Error TestExpedirPost: Se esperaba 400 y se obtuvo", response.StatusCode)
			t.Fail()
		} else {
			t.Log("TestExpedirPost  Finalizado Correctamente (OK)")
		}
	} else {
		t.Error("Error TestExpedirPost:", err.Error())
		t.Fail()
	}
}

func TestValidarDatosExpedicionPost(t *testing.T) {
	if response, err := http.Post("http://localhost:8529/v1/expedir_resolucion/validar_datos_expedicion", "", nil); err == nil {
		if response.StatusCode != 400 {
			t.Error("Error TestValidarDatosExpedicionPost: Se esperaba 400 y se obtuvo", response.StatusCode)
			t.Fail()
		} else {
			t.Log("TestValidarDatosExpedicionPost  Finalizado Correctamente (OK)")
		}
	} else {
		t.Error("Error TestValidarDatosExpedicionPost:", err.Error())
		t.Fail()
	}
}

func TestExpedirModificacionPost(t *testing.T) {
	if response, err := http.Post("http://localhost:8529/v1/expedir_resolucion/expedirModificacion", "", nil); err == nil {
		if response.StatusCode != 400 {
			t.Error("Error TestExpedirModificacionPost: Se esperaba 400 y se obtuvo", response.StatusCode)
			t.Fail()
		} else {
			t.Log("TestExpedirModificacionPost  Finalizado Correctamente (OK)")
		}
	} else {
		t.Error("Error TestExpedirModificacionPost:", err.Error())
		t.Fail()
	}
}

func TestExpedirCancelacionPost(t *testing.T) {
	if response, err := http.Post("http://localhost:8529/v1/expedir_resolucion/cancelar", "", nil); err == nil {
		if response.StatusCode != 400 {
			t.Error("Error TestExpedirCancelacionPost: Se esperaba 400 y se obtuvo", response.StatusCode)
			t.Fail()
		} else {
			t.Log("TestExpedirCancelacionPost  Finalizado Correctamente (OK)")
		}
	} else {
		t.Error("Error TestExpedirCancelacionPost:", err.Error())
		t.Fail()
	}
}
