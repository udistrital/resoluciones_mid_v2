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
