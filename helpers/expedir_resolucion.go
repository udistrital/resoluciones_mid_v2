package helpers

import (
	"fmt"
	"strconv"
	"time"

	"github.com/udistrital/resoluciones_mid_v2/models"
)

func cargarJefeDependenciaActual(dependenciaId int, fecha string) (models.JefeDependencia, error) {
	var jefes []models.JefeDependencia
	url := "jefe_dependencia?query=DependenciaId:" + strconv.Itoa(dependenciaId) + ",FechaFin__gte:" + fecha + ",FechaInicio__lte:" + fecha
	if err := GetRequestLegacy("UrlcrudCore", url, &jefes); err != nil {
		return models.JefeDependencia{}, err
	}
	if len(jefes) == 0 {
		return models.JefeDependencia{}, fmt.Errorf("no se encontró jefe de dependencia vigente para dependencia %d", dependenciaId)
	}

	return jefes[0], nil
}

func cargarSupervisorContratoActual(documento int, fecha string) (models.SupervisorContrato, error) {
	var supervisores []models.SupervisorContrato
	url := "supervisor_contrato?order=desc&sortby=Id&query=Documento:" + strconv.Itoa(documento) + ",FechaFin__gte:" + fecha + ",FechaInicio__lte:" + fecha + "&CargoId.Cargo__startswith:DECANO|VICE"
	if err := GetRequestLegacy("UrlcrudAgora", url, &supervisores); err != nil {
		return models.SupervisorContrato{}, err
	}
	if len(supervisores) == 0 {
		return models.SupervisorContrato{}, fmt.Errorf("no se encontró supervisor vigente para documento %d", documento)
	}

	return supervisores[0], nil
}

func SupervisorActual(dependenciaId int) (supervisorActual models.SupervisorContrato, outputError map[string]interface{}) {
	var fecha = time.Now().Format("2006-01-02") // -- Se debe dejar este una vez se suba
	// var fecha = "2021-01-01"

	jefe, err := cargarJefeDependenciaActual(dependenciaId, fecha)
	if err != nil {
		outputError = map[string]interface{}{"funcion": "/SupervisorActual2", "err": err.Error(), "status": "404"}
		return supervisorActual, outputError
	}

	supervisorActual, err = cargarSupervisorContratoActual(jefe.TerceroId, fecha)
	if err != nil {
		outputError = map[string]interface{}{"funcion": "/SupervisorActual3", "err": err.Error(), "status": "404"}
		return supervisorActual, outputError
	}

	return supervisorActual, nil
}
