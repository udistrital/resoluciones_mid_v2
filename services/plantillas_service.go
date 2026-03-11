package services

import (
	"strconv"

	"github.com/udistrital/resoluciones_mid_v2/helpers"
	"github.com/udistrital/resoluciones_mid_v2/models"
)

func GetPlantillasByAlcance(numeroDocumento string, roles []string, dependenciaFiltro *int) (res []models.Plantilla, outputError map[string]interface{}) {
	alcance, err := ResolveAlcanceUsuario(numeroDocumento, roles)
	if err != nil {
		return nil, err
	}

	if alcance.EsGlobal {
		if dependenciaFiltro != nil && *dependenciaFiltro > 0 {
			return helpers.ListarPlantillas(strconv.Itoa(*dependenciaFiltro))
		}
		return helpers.ListarPlantillas("")
	}

	if len(alcance.Dependencias) == 0 {
		return nil, map[string]interface{}{
			"funcion": "GetPlantillasByAlcance",
			"err":     "el usuario no tiene dependencias asociadas para consultar plantillas",
			"status":  "404",
		}
	}

	if dependenciaFiltro == nil || *dependenciaFiltro <= 0 {
		return nil, map[string]interface{}{
			"funcion": "GetPlantillasByAlcance",
			"err":     "debe seleccionar una dependencia para realizar la consulta",
			"status":  "400",
		}
	}

	if !DependenciaPermitida(*dependenciaFiltro, alcance.Dependencias) {
		return nil, map[string]interface{}{
			"funcion": "GetPlantillasByAlcance",
			"err":     "la dependencia consultada no está autorizada para el usuario",
			"status":  "403",
		}
	}

	return helpers.ListarPlantillas(strconv.Itoa(*dependenciaFiltro))
}
