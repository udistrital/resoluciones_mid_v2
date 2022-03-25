package helpers

import (
	"github.com/udistrital/resoluciones_mid_v2/models"
)

func HomologarDedicacionNombre(dedicacion, tipo string) (dedicacion_id int) {
	HCH := map[string]int{"old": 5, "new": 299}
	HCP := map[string]int{"old": 4, "new": 297}
	MTO := map[string]int{"old": 3, "new": 298}
	TCO := map[string]int{"old": 2, "new": 296}
	dedicaciones := map[string]map[string]int{
		"HCH": HCH,
		"HCP": HCP,
		"TCO": TCO,
		"MTO": MTO,
	}

	return dedicaciones[dedicacion][tipo]
}

func HomologarFacultad(tipo, facultad string) (facultadId string, outputError map[string]interface{}) {
	var endpoint string
	var respuesta models.ObjetoFacultad

	if tipo == "new" {
		endpoint = "facultad_gedep_oikos"
	} else {
		endpoint = "facultad_oikos_gedep"
	}

	url := endpoint + "/" + facultad
	if err := GetRequestWSO2("NscrudHomologacion", url, &respuesta); err != nil {
		outputError = map[string]interface{}{"funcion": "/HomologarFacultad", "err": err.Error(), "status": "500"}
		return facultadId, outputError
	}

	if tipo == "new" {
		facultadId = respuesta.Homologacion.IdGeDep
	} else {
		facultadId = respuesta.Homologacion.IdOikos
	}

	return
}

func HomologarProyectoCurricular(proyectoOld string) (proyectoId string, outputError map[string]interface{}) {
	var respuesta models.ObjetoProyectoCurricular

	url := "proyecto_curricular_cod_proyecto/" + proyectoOld
	if err := GetRequestWSO2("NscrudHomologacion", url, &respuesta); err != nil {
		outputError = map[string]interface{}{"funcion": "/HomologarProyectoCurricular", "err": err.Error(), "status": "500"}
		return proyectoId, outputError
	}
	proyectoId = respuesta.Homologacion.IDOikos

	return
}
