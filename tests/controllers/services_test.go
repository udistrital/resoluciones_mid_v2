package controllers

import (
	"net/http"
	"testing"
)

func TestDesagregadoPlaneacion(t *testing.T) {
	if response, err := http.Post("http://localhost:8529/v1/services/desagregado_planeacion/", "", nil); err == nil {
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
