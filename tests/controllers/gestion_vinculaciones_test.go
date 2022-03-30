package controllers

import (
	"net/http"
	"testing"
)

func TestVinculacionesPost(t *testing.T) {
	if response, err := http.Post("http://localhost:8529/v1/gestion_vinculaciones/", "", nil); err == nil {
		if response.StatusCode != 400 {
			t.Error("Error TestVinculacionesPost: Se esperaba 400 y se obtuvo", response.StatusCode)
			t.Fail()
		} else {
			t.Log("TestVinculacionesPost  Finalizado Correctamente (OK)")
		}
	} else {
		t.Error("Error TestVinculacionesPost:", err.Error())
		t.Fail()
	}
}

func TestModificarVinculacionPost(t *testing.T) {
	if response, err := http.Post("http://localhost:8529/v1/gestion_vinculaciones/modificar_vinculacion", "", nil); err == nil {
		if response.StatusCode != 400 {
			t.Error("Error TestModificarVinculacionPost: Se esperaba 400 y se obtuvo", response.StatusCode)
			t.Fail()
		} else {
			t.Log("TestModificarVinculacionPost  Finalizado Correctamente (OK)")
		}
	} else {
		t.Error("Error TestModificarVinculacionPost:", err.Error())
		t.Fail()
	}
}

func TestDesvincularPost(t *testing.T) {
	if response, err := http.Post("http://localhost:8529/v1/gestion_vinculaciones/desvincular", "", nil); err == nil {
		if response.StatusCode != 400 {
			t.Error("Error TestDesvincularPost: Se esperaba 400 y se obtuvo", response.StatusCode)
			t.Fail()
		} else {
			t.Log("TestDesvincularPost  Finalizado Correctamente (OK)")
		}
	} else {
		t.Error("Error TestDesvincularPost:", err.Error())
		t.Fail()
	}
}

func TestDocentesPrevinculadosGet(t *testing.T) {
	if response, err := http.Get("http://localhost:8529/v1/gestion_vinculaciones/219"); err == nil {
		if response.StatusCode != 200 {
			t.Error("Error TestDocentesPrevinculadosGet: Se esperaba 200 y se obtuvo", response.StatusCode)
			t.Fail()
		} else {
			t.Log("TestDocentesPrevinculadosGet Finalizado Correctamente (OK)")
		}
	} else {
		t.Error("Error TestDocentesPrevinculadosGet:", err.Error())
		t.Fail()
	}
}

func TestDocentesCargaHorariaGet(t *testing.T) {
	if response, err := http.Get("http://localhost:8529/v1/gestion_vinculaciones/docentes_carga_horaria/2021/1/HCH/33/PREGRADO"); err == nil {
		if response.StatusCode != 200 {
			t.Error("Error TestDocentesCargaHorariaGet: Se esperaba 200 y se obtuvo", response.StatusCode)
			t.Fail()
		} else {
			t.Log("TestDocentesCargaHorariaGet Finalizado Correctamente (OK)")
		}
	} else {
		t.Error("Error TestDocentesCargaHorariaGet:", err.Error())
		t.Fail()
	}
}

func TestInformeVinculacionesPost(t *testing.T) {
	if response, err := http.Post("http://localhost:8529/v1/gestion_vinculaciones/informe_vinculaciones", "", nil); err == nil {
		if response.StatusCode != 400 {
			t.Error("Error TestInformeVinculacionesPost: Se esperaba 400 y se obtuvo", response.StatusCode)
			t.Fail()
		} else {
			t.Log("TestInformeVinculacionesPost  Finalizado Correctamente (OK)")
		}
	} else {
		t.Error("Error TestInformeVinculacionesPost:", err.Error())
		t.Fail()
	}
}

func TestDesvincularDocentesPost(t *testing.T) {
	if response, err := http.Post("http://localhost:8529/v1/gestion_vinculaciones/desvincular_docentes", "", nil); err == nil {
		if response.StatusCode != 400 {
			t.Error("Error TestDesvincularDocentesPost: Se esperaba 400 y se obtuvo", response.StatusCode)
			t.Fail()
		} else {
			t.Log("TestDesvincularDocentesPost  Finalizado Correctamente (OK)")
		}
	} else {
		t.Error("Error TestDesvincularDocentesPost:", err.Error())
		t.Fail()
	}
}

func TestCalcularValorContratosSeleccionadosPost(t *testing.T) {
	if response, err := http.Post("http://localhost:8529/v1/gestion_vinculaciones/calcular_valor_contratos_seleccionados", "", nil); err == nil {
		if response.StatusCode != 400 {
			t.Error("Error TestCalcularValorContratosSeleccionadosPost: Se esperaba 400 y se obtuvo", response.StatusCode)
			t.Fail()
		} else {
			t.Log("TestCalcularValorContratosSeleccionadosPost  Finalizado Correctamente (OK)")
		}
	} else {
		t.Error("Error TestCalcularValorContratosSeleccionadosPost:", err.Error())
		t.Fail()
	}
}
