package helpers

import (
	"fmt"
	"strconv"
	"time"

	"github.com/udistrital/resoluciones_mid_v2/models"
)

func SupervisorActual(id_resolucion int) (supervisor_actual models.SupervisorContrato, outputError map[string]interface{}) {
	var r models.Resolucion
	var j []models.JefeDependencia
	var s []models.SupervisorContrato
	//var fecha = time.Now().Format("2006-01-02")   -- Se debe dejar este una vez se suba
	var fecha = "2018-01-01"
	//If Resolucion (GET)
	url := "resolucion/" + strconv.Itoa(id_resolucion)
	if err := GetRequestNew("UrlCrudResoluciones", url, &r); err == nil {
		//If Jefe_dependencia (GET)
		url = "jefe_dependencia?query=DependenciaId:" + strconv.Itoa(r.DependenciaId) + ",FechaFin__gte:" + fecha + ",FechaInicio__lte:" + fecha
		if err := GetRequestLegacy("UrlcrudCore", url, &j); err == nil {
			//If Supervisor (GET)
			url = "supervisor_contrato?query=Documento:" + strconv.Itoa(j[0].TerceroId) + ",FechaFin__gte:" + fecha + ",FechaInicio__lte:" + fecha + "&CargoId.Cargo__startswith:DECANO|VICE"
			if err := GetRequestLegacy("UrlcrudAgora", url, &s); err == nil {
				fmt.Println(s[0])
				return s[0], nil
			} else { //If Jefe_dependencia (GET)
				fmt.Println("He fallado un poquito en If Supervisor 1 (GET) en el método SupervisorActual, solucioname!!! ", err)
				outputError = map[string]interface{}{"funcion": "/SupervisorActual3", "err": err.Error(), "status": "404"}
				return s[0], outputError
			}
		} else { //If Jefe_dependencia (GET)
			fmt.Println("He fallado un poquito en If Jefe_dependencia 2 (GET) en el método SupervisorActual, solucioname!!! ", err)
			outputError = map[string]interface{}{"funcion": "/SupervisorActua2", "err": err.Error(), "status": "404"}
			return s[0], outputError
		}
	} else { //If Resolucion (GET)
		fmt.Println("He fallado un poquito en If Resolucion 3 (GET) en el método SupervisorActual, solucioname!!! ", err)
		outputError = map[string]interface{}{"funcion": "/SupervisorActual", "err": err.Error(), "status": "404"}
		return s[0], outputError
	}
}

func CalcularFechaFin(fecha_inicio time.Time, numero_semanas int) (fecha_fin time.Time) {
	var entero int
	var decimal float32
	meses := float32(numero_semanas) / 4
	entero = int(meses)
	decimal = meses - float32(entero)
	numero_dias := ((decimal * 4) * 7)
	f_i := fecha_inicio
	after := f_i.AddDate(0, entero, int(numero_dias))
	return after
}
