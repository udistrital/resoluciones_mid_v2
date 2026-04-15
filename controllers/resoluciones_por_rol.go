package controllers

import (
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
// @Description Obtiene el alcance del usuario según sus roles
// @Param	numero_documento	query	string	true	"Número de documento del usuario"
// @Param	roles			query	string	true	"Roles del usuario separados por coma"
// @Success 200 {object} map[string]interface{}
// @Failure 400 bad request
// @Failure 500 Internal server error
// @router /dependencias [get]
func (c *ResolucionesPorRolController) GetDependenciasByRol() {
	defer helpers.ErrorController(c.Controller, "ResolucionesPorRolController")

	authContext := requireAuthenticatedContext(buildAuthenticatedContext(&c.Controller), "GetDependenciasByRol")

	alcance, errMap := services.ResolveAlcanceUsuario(authContext.NumeroDocumento, authContext.Roles)
	if errMap != nil {
		panic(errMap)
	}

	writeJSON(&c.Controller, 200, "Alcance del usuario obtenido con éxito", alcance, nil)
}

// GetResolucionesByDependencia ...
// @Title GetResolucionesByDependencia
// @Description Obtiene resoluciones según el alcance del usuario y vigencia
// @Param	numero_documento	query	string	true	"Número de documento del usuario"
// @Param	roles			query	string	true	"Roles del usuario separados por coma"
// @Param	vigencia		query	int		true	"Vigencia de las resoluciones"
// @Param	id_oikos		query	int		false	"ID OIKOS para filtrar una dependencia específica"
// @Success 200 {object} map[string]interface{}
// @Failure 400 bad request
// @Failure 500 Internal server error
// @router /consulta [get]
func (c *ResolucionesPorRolController) GetResolucionesByDependencia() {
	defer helpers.ErrorController(c.Controller, "ResolucionesPorRolController")

	authContext := requireAuthenticatedContext(buildAuthenticatedContext(&c.Controller), "GetResolucionesByDependencia")
	filtro := buildFiltroConsulta(c, "vigencia")

	vigencia, errVigencia := c.GetInt("vigencia")
	idOikos, errIdOikos := c.GetInt("id_oikos")

	var dependenciaFiltro *int

	if errVigencia != nil || vigencia <= 0 {
		panic(badRequest("GetResolucionesByDependencia", validateNamedPositiveInt(filtro.Vigencia, "vigencia")))
	}

	if err := validateFiltroConsulta(filtro); err != nil {
		panic(badRequest("GetResolucionesByDependencia", err))
	}

	if errIdOikos == nil && idOikos > 0 {
		dependenciaFiltro = &idOikos
	}
	resoluciones, total, errMap := services.GetResolucionesTablaByAlcance(
		authContext.NumeroDocumento,
		authContext.Roles,
		filtro,
		dependenciaFiltro,
	)
	if errMap != nil {
		panic(errMap)
	}

	writeJSON(&c.Controller, 200, "Resoluciones obtenidas con éxito", resoluciones, map[string]interface{}{
		"Total": total,
	})
}
