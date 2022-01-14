package controllers

import (
	"net/http"
	"testing"
)

func TestPlantillasPost(t *testing.T) {
	if response, err := http.Post("http://localhost:8529/v1/gestion_plantillas/", "", nil); err == nil {
		if response.StatusCode != 400 {
			t.Error("Error TestPlantillasPost: Se esperaba 400 y se obtuvo", response.StatusCode)
			t.Fail()
		} else {
			t.Log("TestPlantillasPost Finalizado Correctamente (OK)")
		}
	} else {
		t.Error("Error TestPlantillasPost:", err.Error())
		t.Fail()
	}
}

func TestPlantillasGetOne(t *testing.T) {
	if response, err := http.Get("http://localhost:8529/v1/gestion_plantillas/1"); err == nil {
		if response.StatusCode != 200 {
			t.Error("Error TestPlantillasGetOne: Se esperaba 200 y se obtuvo", response.StatusCode)
			t.Fail()
		} else {
			t.Log("TestPlantillasGetOne Finalizado Correctamente (OK)")
		}
	} else {
		t.Error("Error TestPlantillasGetOne:", err.Error())
		t.Fail()
	}
}

func TestPlantillasGetAll(t *testing.T) {
	if response, err := http.Get("http://localhost:8529/v1/gestion_plantillas/"); err == nil {
		if response.StatusCode != 200 {
			t.Error("Error TestPlantillasGetAll: Se esperaba 400 y se obtuvo", response.StatusCode)
			t.Fail()
		} else {
			t.Log("TestPlantillasGetAll Finalizado Correctamente (OK)")
		}
	} else {
		t.Error("ErrorTestPlantillasGetAll:", err.Error())
		t.Fail()
	}
}
