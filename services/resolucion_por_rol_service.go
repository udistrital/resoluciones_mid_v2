package services

import (
	"fmt"
	"sort"

	"github.com/udistrital/resoluciones_mid_v2/helpers"
	"github.com/udistrital/resoluciones_mid_v2/models"
)

func GetResolucionesByDependenciaIdAndVigencia(idOikos int, vigencia int) (res []models.Resolucion, outputError map[string]interface{}) {
	var temp []models.Resolucion

	route := fmt.Sprintf(
		"resolucion?query=DependenciaId:%d,Vigencia:%d&limit=0&sortby=Id&order=desc",
		idOikos,
		vigencia,
	)

	if err := helpers.GetRequestNew("UrlCrudResoluciones", route, &temp); err != nil {
		return nil, map[string]interface{}{
			"funcion": "GetResolucionesByDependenciaIdAndVigencia",
			"err":     err.Error(),
			"status":  "502",
		}
	}

	sort.Slice(temp, func(i, j int) bool {
		return temp[i].Id > temp[j].Id
	})

	return temp, nil
}

func GetResolucionesByVigencia(vigencia int) (res []models.Resolucion, outputError map[string]interface{}) {
	var temp []models.Resolucion

	route := fmt.Sprintf(
		"resolucion?query=Vigencia:%d&limit=0&sortby=Id&order=desc",
		vigencia,
	)

	if err := helpers.GetRequestNew("UrlCrudResoluciones", route, &temp); err != nil {
		return nil, map[string]interface{}{
			"funcion": "GetResolucionesByVigencia",
			"err":     err.Error(),
			"status":  "502",
		}
	}

	sort.Slice(temp, func(i, j int) bool {
		return temp[i].Id > temp[j].Id
	})

	return temp, nil
}

func dependenciaPermitida(idOikos int, dependencias []models.DependenciaUsuario) bool {
	for _, dep := range dependencias {
		if dep.IdOikos == idOikos {
			return true
		}
	}
	return false
}

func GetResolucionesByAlcance(numeroDocumento string, roles []string, vigencia int, dependenciaFiltro *int) (res []models.Resolucion, outputError map[string]interface{}) {
	alcance, err := ResolveAlcanceUsuario(numeroDocumento, roles)
	if err != nil {
		return nil, err
	}

	if alcance.EsGlobal {
		if dependenciaFiltro != nil && *dependenciaFiltro > 0 {
			return GetResolucionesByDependenciaIdAndVigencia(*dependenciaFiltro, vigencia)
		}
		return GetResolucionesByVigencia(vigencia)
	}

	if len(alcance.Dependencias) == 0 {
		return nil, map[string]interface{}{
			"funcion": "GetResolucionesByAlcance",
			"err":     "el usuario no tiene dependencias asociadas para consultar resoluciones",
			"status":  "404",
		}
	}

	if dependenciaFiltro != nil && *dependenciaFiltro > 0 {
		if !dependenciaPermitida(*dependenciaFiltro, alcance.Dependencias) {
			return nil, map[string]interface{}{
				"funcion": "GetResolucionesByAlcance",
				"err":     "la dependencia consultada no está autorizada para el usuario",
				"status":  "403",
			}
		}

		return GetResolucionesByDependenciaIdAndVigencia(*dependenciaFiltro, vigencia)
	}

	resultado := make([]models.Resolucion, 0)
	seen := make(map[int]bool)

	for _, dep := range alcance.Dependencias {
		resoluciones, errMap := GetResolucionesByDependenciaIdAndVigencia(dep.IdOikos, vigencia)
		if errMap != nil {
			return nil, errMap
		}

		for _, resolucion := range resoluciones {
			if !seen[resolucion.Id] {
				seen[resolucion.Id] = true
				resultado = append(resultado, resolucion)
			}
		}
	}

	sort.Slice(resultado, func(i, j int) bool {
		return resultado[i].Id > resultado[j].Id
	})

	return resultado, nil
}
