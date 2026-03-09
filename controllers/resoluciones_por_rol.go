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
	c.Mapping("GetDependenciasByRol", c.GetDependenciasByRol)
	c.Mapping("GetResolucionesByDependencia", c.GetResolucionesByDependencia)
}

// @router /dependencias [get]
func (c *ResolucionesPorRolController) GetDependenciasByRol() {
	defer helpers.ErrorController(c.Controller, "ResolucionesPorRolController")

	numeroDocumento := strings.TrimSpace(c.GetString("numero_documento"))
	rol := strings.ToUpper(strings.TrimSpace(c.GetString("rol")))

	if numeroDocumento == "" {
		panic(map[string]interface{}{
			"funcion": "GetDependenciasByRol",
			"err":     "numero_documento es requerido",
			"status":  "400",
		})
	}

	if rol == "" {
		panic(map[string]interface{}{
			"funcion": "GetDependenciasByRol",
			"err":     "rol es requerido",
			"status":  "400",
		})
	}

	switch rol {
	case "DECANO", "ASISTENTE_DECANATURA":
		// ok
	default:
		panic(map[string]interface{}{
			"funcion": "GetDependenciasByRol",
			"err":     "rol no soportado",
			"status":  "400",
		})
	}

	dependencias, errMap := services.ResolveDependenciasByRol(numeroDocumento, rol)
	if errMap != nil {
		panic(errMap)
	}

	c.Data["json"] = map[string]interface{}{
		"Success": true,
		"Status":  200,
		"Message": "Dependencias obtenidas con éxito",
		"Data":    dependencias,
	}
	c.ServeJSON()
}

// @router /consulta [get]
func (c *ResolucionesPorRolController) GetResolucionesByDependencia() {
	defer helpers.ErrorController(c.Controller, "ResolucionesPorRolController")

	idOikos, errIdOikos := c.GetInt("id_oikos")
	vigencia, errVigencia := c.GetInt("vigencia")

	if errIdOikos != nil || idOikos <= 0 {
		panic(map[string]interface{}{
			"funcion": "GetResolucionesByDependencia",
			"err":     "id_oikos es requerido y debe ser válido",
			"status":  "400",
		})
	}

	if errVigencia != nil || vigencia <= 0 {
		panic(map[string]interface{}{
			"funcion": "GetResolucionesByDependencia",
			"err":     "vigencia es requerida y debe ser válida",
			"status":  "400",
		})
	}

	resoluciones, errMap := services.GetResolucionesByDependenciaIdAndVigencia(idOikos, vigencia)
	if errMap != nil {
		panic(errMap)
	}

	c.Data["json"] = map[string]interface{}{
		"Success": true,
		"Status":  200,
		"Message": "Resoluciones obtenidas con éxito",
		"Data": map[string]interface{}{
			"id_oikos":     idOikos,
			"vigencia":     vigencia,
			"resoluciones": resoluciones,
		},
	}
	c.ServeJSON()
}
