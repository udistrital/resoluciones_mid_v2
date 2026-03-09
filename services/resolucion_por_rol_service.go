package services

import (
	"fmt"
	"sort"

	"github.com/udistrital/resoluciones_mid_v2/helpers"
	"github.com/udistrital/resoluciones_mid_v2/models"
)

func GetResolucionesByDependenciaIdsAll(idsOikos []int, vigencia int) (res []models.Resolucion, outputError map[string]interface{}) {
	resMap := make(map[int]models.Resolucion)

	for _, idOikos := range idsOikos {
		var temp []models.Resolucion

		route := fmt.Sprintf(
			"resolucion?query=DependenciaId:%d,Vigencia:%d&limit=0&sortby=Id&order=desc",
			idOikos,
			vigencia,
		)

		if err := helpers.GetRequestNew("UrlCrudResoluciones", route, &temp); err != nil {
			return nil, map[string]interface{}{
				"funcion": "GetResolucionesByDependenciaIdsAll",
				"err":     err.Error(),
				"status":  "502",
			}
		}

		for _, r := range temp {
			resMap[r.Id] = r
		}
	}

	for _, r := range resMap {
		res = append(res, r)
	}

	sort.Slice(res, func(i, j int) bool {
		return res[i].Id > res[j].Id
	})

	return res, nil
}
