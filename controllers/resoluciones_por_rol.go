package controllers

import (
	"strings"

	"github.com/astaxie/beego"
	"github.com/udistrital/resoluciones_mid_v2/helpers"
	"github.com/udistrital/resoluciones_mid_v2/services"
)

type ResolucionesPorRolController struct {
	beego.Controller
}

func (c *ResolucionesPorRolController) URLMapping() {
	c.Mapping("GetByRol", c.GetByRol)
}

// @router / [get]
func (c *ResolucionesPorRolController) GetByRol() {
	defer helpers.ErrorController(c.Controller, "ResolucionesPorRolController")

	numeroDocumento := strings.TrimSpace(c.GetString("numero_documento"))
	rol := strings.ToUpper(strings.TrimSpace(c.GetString("rol")))
	vigencia, errVigencia := c.GetInt("vigencia")

	if numeroDocumento == "" {
		panic(map[string]interface{}{
			"funcion": "GetByRol",
			"err":     "numero_documento es requerido",
			"status":  "400",
		})
	}

	if rol == "" {
		panic(map[string]interface{}{
			"funcion": "GetByRol",
			"err":     "rol es requerido",
			"status":  "400",
		})
	}

	if errVigencia != nil || vigencia <= 0 {
		panic(map[string]interface{}{
			"funcion": "GetByRol",
			"err":     "vigencia es requerida y debe ser válida",
			"status":  "400",
		})
	}

	switch rol {
	case "DECANO", "ASISTENTE_DECANATURA":
	default:
		panic(map[string]interface{}{
			"funcion": "GetByRol",
			"err":     "rol no soportado",
			"status":  "400",
		})
	}

	idsOikos, errMap := services.ResolveOikosByRol(numeroDocumento, rol)
	if errMap != nil {
		panic(errMap)
	}

	res, errMap := services.GetResolucionesByDependenciaIdsAll(idsOikos, vigencia)
	if errMap != nil {
		panic(errMap)
	}

	c.Data["json"] = map[string]interface{}{
		"Success": true,
		"Status":  200,
		"Message": "Resoluciones por rol",
		"Data": map[string]interface{}{
			"rol":              rol,
			"numero_documento": numeroDocumento,
			"vigencia":         vigencia,
			"id_oikos":         idsOikos,
			"filtro": map[string]interface{}{
				"campo": "DependenciaId, Vigencia",
				"valor": map[string]interface{}{
					"id_oikos": idsOikos,
					"vigencia": vigencia,
				},
			},
			"resoluciones": res,
		},
	}
	c.ServeJSON()
}
