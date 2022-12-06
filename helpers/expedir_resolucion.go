package helpers

import (
	"fmt"
	"strconv"
	"time"

	"github.com/udistrital/resoluciones_mid_v2/models"
)

func SupervisorActual(resolucionId int) (supervisorActual models.SupervisorContrato, outputError map[string]interface{}) {
	var r models.Resolucion
	var j []models.JefeDependencia
	var s []models.SupervisorContrato
	var fecha = time.Now().Format("2006-01-02") // -- Se debe dejar este una vez se suba
	// var fecha = "2021-01-01"
	//If Resolucion (GET)
	url := "resolucion/" + strconv.Itoa(resolucionId)
	if err := GetRequestNew("UrlCrudResoluciones", url, &r); err == nil {
		//If Jefe_dependencia (GET)
		url = "jefe_dependencia?query=DependenciaId:" + strconv.Itoa(r.DependenciaId) + ",FechaFin__gte:" + fecha + ",FechaInicio__lte:" + fecha
		if err := GetRequestLegacy("UrlcrudCore", url, &j); err == nil && len(j) > 0 {
			//If Supervisor (GET)
			url = "supervisor_contrato?order=desc&sortby=Id&query=Documento:" + strconv.Itoa(j[0].TerceroId) + ",FechaFin__gte:" + fecha + ",FechaInicio__lte:" + fecha + "&CargoId.Cargo__startswith:DECANO|VICE"
			if err := GetRequestLegacy("UrlcrudAgora", url, &s); err == nil && len(s) > 0 {
				fmt.Println(s)
				return s[0], nil
			} else { //If Jefe_dependencia (GET)
				fmt.Println("He fallado un poquito en If Supervisor 1 (GET) en el método SupervisorActual, solucioname!!! ", err)
				outputError = map[string]interface{}{"funcion": "/SupervisorActual3", "err": err.Error(), "status": "404"}
				return supervisorActual, outputError
			}
		} else { //If Jefe_dependencia (GET)
			fmt.Println("He fallado un poquito en If Jefe_dependencia 2 (GET) en el método SupervisorActual, solucioname!!! ", err)
			outputError = map[string]interface{}{"funcion": "/SupervisorActua2", "err": err.Error(), "status": "404"}
			return supervisorActual, outputError
		}
	} else { //If Resolucion (GET)
		fmt.Println("He fallado un poquito en If Resolucion 3 (GET) en el método SupervisorActual, solucioname!!! ", err)
		outputError = map[string]interface{}{"funcion": "/SupervisorActual", "err": err.Error(), "status": "404"}
		return supervisorActual, outputError
	}
}

// Calcula la fecha de fin de un contrato a partir de la fecha de inicio y el numero de semanas
func CalcularFechaFin(fechaInicio time.Time, numeroSemanas int) (fechaFin time.Time) {
	var mesEntero, dias int
	var decimal, meses float32
	if numeroSemanas%4 == 0 {
		// Meses de 4 semanas
		meses = float32(numeroSemanas) / 4
	} else {
		dias = numeroSemanas * 7
		meses = float32(dias) / 30
	}
	mesEntero = int(meses)
	decimal = meses - float32(mesEntero)
	numeroDias := decimal * 30
	f := fechaInicio
	after := f.AddDate(0, mesEntero, int(numeroDias-1))
	return after
}
