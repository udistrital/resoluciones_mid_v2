package services

import (
	"fmt"
	"sort"
	"strconv"

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
		if !DependenciaPermitida(*dependenciaFiltro, alcance.Dependencias) {
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

func GetResolucionesTablaByAlcance(numeroDocumento string, roles []string, filtro models.Filtro, dependenciaFiltro *int) (res []models.Resoluciones, total int, outputError map[string]interface{}) {
	alcance, err := ResolveAlcanceUsuario(numeroDocumento, roles)
	if err != nil {
		return nil, 0, err
	}

	if alcance.EsGlobal {
		if dependenciaFiltro != nil && *dependenciaFiltro > 0 {
			filtro.FacultadId = strconv.Itoa(*dependenciaFiltro)
		}
	} else {
		if len(alcance.Dependencias) == 0 {
			return nil, 0, map[string]interface{}{
				"funcion": "GetResolucionesTablaByAlcance",
				"err":     "el usuario no tiene dependencias asociadas para consultar resoluciones",
				"status":  "404",
			}
		}

		if dependenciaFiltro != nil && *dependenciaFiltro > 0 {
			if !DependenciaPermitida(*dependenciaFiltro, alcance.Dependencias) {
				return nil, 0, map[string]interface{}{
					"funcion": "GetResolucionesTablaByAlcance",
					"err":     "la dependencia consultada no está autorizada para el usuario",
					"status":  "403",
				}
			}

			filtro.FacultadId = strconv.Itoa(*dependenciaFiltro)
		} else {
			if len(alcance.Dependencias) == 1 {
				filtro.FacultadId = strconv.Itoa(alcance.Dependencias[0].IdOikos)
			} else {
				return nil, 0, map[string]interface{}{
					"funcion": "GetResolucionesTablaByAlcance",
					"err":     "debe seleccionar una dependencia para realizar la consulta",
					"status":  "400",
				}
			}
		}
	}

	listado, total, err2 := helpers.ListarResolucionesFiltradas(filtro)
	if err2 != nil {
		err2["funcion"] = "GetResolucionesTablaByAlcance"
		return nil, 0, err2
	}

	return listado, total, nil
}

func UsuarioPuedeConsultarResolucion(resolucionId int, numeroDocumento string, roles []string) (bool, map[string]interface{}) {
	var resolucion models.Resolucion

	route := fmt.Sprintf("resolucion/%d", resolucionId)
	if err := helpers.GetRequestNew("UrlCrudResoluciones", route, &resolucion); err != nil {
		return false, map[string]interface{}{
			"funcion": "UsuarioPuedeConsultarResolucion",
			"err":     err.Error(),
			"status":  "502",
		}
	}

	alcance, errMap := ResolveAlcanceUsuario(numeroDocumento, roles)
	if errMap != nil {
		return false, errMap
	}

	if alcance.EsGlobal {
		return true, nil
	}

	if !DependenciaPermitida(resolucion.DependenciaId, alcance.Dependencias) {
		return false, map[string]interface{}{
			"funcion": "UsuarioPuedeConsultarResolucion",
			"err":     "La resolución no está autorizada para el usuario",
			"status":  "403",
		}
	}

	return true, nil
}
