package helpers

import (
	"fmt"
	"strconv"
	"time"

	"github.com/astaxie/beego/logs"
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

func GetContenidoResolucion(id_resolucion string, id_facultad string) (contenidoResolucion models.ResolucionCompleta, outputError map[string]interface{}) {

	defer func() {
		if err := recover(); err != nil {
			outputError = map[string]interface{}{"funcion": "/GetContenidoResolucion", "err": err, "status": "502"}
			panic(outputError)
		}
	}()

	var ordenador_gasto []models.OrdenadorGasto
	var jefe_dependencia []models.JefeDependencia

	url := "contenido_resolucion/" + id_resolucion
	if err1 := GetRequestNew("UrlCrudResoluciones", url, &contenidoResolucion); err1 == nil {
		url = "ordenador_gasto?limit=-1&query=DependenciaId:" + id_facultad
		if err2 := GetRequestLegacy("UrlcrudCore", url, &ordenador_gasto); err2 == nil {
			fmt.Println(ordenador_gasto)
			if ordenador_gasto == nil || len(ordenador_gasto) == 0 {
				url = "ordenador_gasto?query=Id:1"
				if err3 := GetRequestLegacy("UrlcrudCore", url, &ordenador_gasto); err3 == nil {
					contenidoResolucion.OrdenadorGasto = ordenador_gasto[0]
				} else {
					logs.Error(err3)
					outputError = map[string]interface{}{"funcion": "/GetContenidoResolucion3", "err3": err3, "status": "502"}
					return contenidoResolucion, outputError
				}
			} else {
				contenidoResolucion.OrdenadorGasto = ordenador_gasto[0]
			}

		} else {
			logs.Error(err2)
			outputError = map[string]interface{}{"funcion": "/GetContenidoResolucion2", "err2": err2, "status": "502"}
			return contenidoResolucion, outputError
		}
	} else {
		logs.Error(err1)
		outputError = map[string]interface{}{"funcion": "/GetContenidoResolucion1", "err": err1, "status": "502"}
		return contenidoResolucion, outputError
	}

	fecha_actual := time.Now().Format("2006-01-02")
	var err5 map[string]interface{}
	url = "jefe_dependencia?query=DependenciaId:" + id_facultad + ",FechaFin__gte:" + fecha_actual + ",FechaInicio__lte:" + fecha_actual
	if err4 := GetRequestLegacy("UrlcrudCore", url, &jefe_dependencia); err4 == nil {
		contenidoResolucion.OrdenadorGasto.NombreOrdenador, err5 = BuscarNombreProveedor(jefe_dependencia[0].TerceroId)
		if err5 != nil {
			logs.Error(err4)
			//outputError = map[string]interface{}{"funcion": "/GetContenidoResolucion5", "err5": err5, "status": "502"}
			return contenidoResolucion, err5
		}
	} else {
		logs.Error(err4)
		outputError = map[string]interface{}{"funcion": "/GetContenidoResolucion4", "err4": err4, "status": "502"}
		return contenidoResolucion, outputError
	}

	return contenidoResolucion, nil
}
