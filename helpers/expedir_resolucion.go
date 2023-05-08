package helpers

import (
	"fmt"
	"strconv"
	"time"

	"github.com/udistrital/resoluciones_mid_v2/models"
)

func SupervisorActual(dependenciaId int) (supervisorActual models.SupervisorContrato, outputError map[string]interface{}) {
	var j []models.JefeDependencia
	var s []models.SupervisorContrato
	var fecha = time.Now().Format("2006-01-02") // -- Se debe dejar este una vez se suba
	// var fecha = "2021-01-01"
	//If Jefe_dependencia (GET)
	url := "jefe_dependencia?query=DependenciaId:" + strconv.Itoa(dependenciaId) + ",FechaFin__gte:" + fecha + ",FechaInicio__lte:" + fecha
	if err := GetRequestLegacy("UrlcrudCore", url, &j); err == nil && len(j) > 0 {
		//If Supervisor (GET)
		url = "supervisor_contrato?order=desc&sortby=Id&query=Documento:" + strconv.Itoa(j[0].TerceroId) + ",FechaFin__gte:" + fecha + ",FechaInicio__lte:" + fecha + "&CargoId.Cargo__startswith:DECANO|VICE"
		if err := GetRequestLegacy("UrlcrudAgora", url, &s); err == nil && len(s) > 0 {
			fmt.Println(s)
			return s[0], nil
		} else { //If Jefe_dependencia (GET)
			fmt.Println("No se ha encontrado supervisor activo en la fecha actual!!!", err)
			outputError = map[string]interface{}{"funcion": "/SupervisorActual3", "err": err.Error(), "status": "404"}
			return supervisorActual, outputError
		}
	} else { //If Jefe_dependencia (GET)
		fmt.Println("No se ha encontrado jefe de dependencia activo en la fecha actual!!! ", err)
		outputError = map[string]interface{}{"funcion": "/SupervisorActual2", "err": err.Error(), "status": "404"}
		return supervisorActual, outputError
	}

}

// Calcula la fecha de fin de un contrato a partir de la fecha de inicio y el numero de semanas
func CalcularFechaFin(fechaInicio time.Time, numeroSemanas int) (fechaFin time.Time) {

	// Versión original con meses de 4 semanas
	// meses = float32(numeroSemanas) / 4
	// entero = int(meses)
	// decimal = meses - float32(entero)
	// numero_dias := ((decimal * 4) * 7)
	// f := fecha_inicio
	// after := f.AddDate(0, entero, int(numero_dias))

	// Primera modificación con meses de 30 dias estrictos
	// var mesEntero, dias int
	// var decimal, meses float32
	// dias = numeroSemanas * 7
	// meses = float32(dias) / 30
	// mesEntero = int(meses)
	// decimal = meses - float32(mesEntero)
	// numeroDias := decimal * 30
	// after := fechaInicio.AddDate(0, mesEntero, int(numeroDias)-1)

	// Segunda modificación, estrictamente por dias o semanas de 7 dias ajustando
	// los dias sobrantes cuando el calendario academico inicia a mitad de semana
	// de manera que la fecha de fin resultante sea a final de semana

	dias := numeroSemanas * 7
	if fechaInicio.Weekday() >= 3 {
		dias += (7 - int(fechaInicio.Weekday()))
	} else {
		dias -= int(fechaInicio.Weekday())
	}
	after := fechaInicio.AddDate(0, 0, dias-1)
	return after
}
