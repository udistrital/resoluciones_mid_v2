package controllers

import (
	"strings"

	"github.com/astaxie/beego"
	"github.com/udistrital/resoluciones_mid_v2/helpers"
	"github.com/udistrital/resoluciones_mid_v2/services"
)

// ResolucionesPorRolController operations for ResolucionesPorRol
type ResolucionesPorRolController struct {
	beego.Controller
}

// URLMapping ...
func (c *ResolucionesPorRolController) URLMapping() {
	c.Mapping("GetDependenciasByRol", c.GetDependenciasByRol)
	c.Mapping("GetResolucionesByDependencia", c.GetResolucionesByDependencia)
}

// GetDependenciasByRol ...
// @Title GetDependenciasByRol
// @Description Obtiene las dependencias asociadas a un usuario según su rol
// @Param	numero_documento	query	string	true	"Número de documento del usuario"
// @Param	rol				query	string	true	"Rol del usuario (DECANO o ASISTENTE_DECANATURA)"
// @Success 200 {object} map[string]interface{}
// @Failure 400 bad request
// @Failure 500 Internal server error
// @router /dependencias [get]
func (c *ResolucionesPorRolController) GetDependenciasByRol() {
	defer helpers.ErrorController(c.Controller, "ResolucionesPorRolController")

	numeroDocumento := strings.TrimSpace(c.GetString("numero_documento"))
	rol := strings.ToUpper(strings.TrimSpace(c.GetString("rol")))

	if numeroDocumento == "" {
		panic(map[string]interface{}{"funcion": "GetDependenciasByRol", "err": "numero_documento es requerido", "status": "400"})
	}

	if rol == "" {
		panic(map[string]interface{}{"funcion": "GetDependenciasByRol", "err": "rol es requerido", "status": "400"})
	}

	switch rol {
	case "DECANO", "ASISTENTE_DECANATURA":
	default:
		panic(map[string]interface{}{"funcion": "GetDependenciasByRol", "err": "rol no soportado", "status": "400"})
	}

	dependencias, errMap := services.ResolveDependenciasByRol(numeroDocumento, rol)
	if errMap != nil {
		panic(errMap)
	}

	c.Ctx.Output.SetStatus(200)
	c.Data["json"] = map[string]interface{}{
		"Success": true,
		"Status":  200,
		"Message": "Dependencias obtenidas con éxito",
		"Data":    dependencias,
	}
	c.ServeJSON()
}

// GetResolucionesByDependencia ...
// @Title GetResolucionesByDependencia
// @Description Obtiene las resoluciones asociadas a una dependencia y vigencia
// @Param	id_oikos	query	int	true	"ID de la dependencia (OIKOS)"
// @Param	vigencia	query	int	true	"Vigencia de las resoluciones"
// @Success 200 {object} map[string]interface{}
// @Failure 400 bad request
// @Failure 500 Internal server error
// @router /consulta [get]
func (c *ResolucionesPorRolController) GetResolucionesByDependencia() {
	defer helpers.ErrorController(c.Controller, "ResolucionesPorRolController")

	idOikos, errIdOikos := c.GetInt("id_oikos")
	vigencia, errVigencia := c.GetInt("vigencia")

	if errIdOikos != nil || idOikos <= 0 {
		panic(map[string]interface{}{"funcion": "GetResolucionesByDependencia", "err": "id_oikos es requerido y debe ser válido", "status": "400"})
	}

	if errVigencia != nil || vigencia <= 0 {
		panic(map[string]interface{}{"funcion": "GetResolucionesByDependencia", "err": "vigencia es requerida y debe ser válida", "status": "400"})
	}

	resoluciones, errMap := services.GetResolucionesByDependenciaIdAndVigencia(idOikos, vigencia)
	if errMap != nil {
		panic(errMap)
	}

	c.Ctx.Output.SetStatus(200)
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
