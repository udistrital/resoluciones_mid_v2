package controllers

import (
	"strings"

	"github.com/astaxie/beego"
	"github.com/udistrital/resoluciones_mid_v2/helpers"
	"github.com/udistrital/resoluciones_mid_v2/models"
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

func parseRolesParam(raw string) []string {
	if strings.TrimSpace(raw) == "" {
		return []string{}
	}

	parts := strings.Split(raw, ",")
	roles := make([]string, 0)

	for _, part := range parts {
		rol := strings.ToUpper(strings.TrimSpace(part))
		if rol != "" {
			roles = append(roles, rol)
		}
	}

	return roles
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return strings.TrimSpace(value)
		}
	}

	return ""
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

	numeroDocumento := strings.TrimSpace(c.GetString("numero_documento"))
	rolesRaw := c.GetString("roles")
	roles := parseRolesParam(rolesRaw)

	if numeroDocumento == "" {
		panic(map[string]interface{}{
			"funcion": "GetDependenciasByRol",
			"err":     "numero_documento es requerido",
			"status":  "400",
		})
	}

	if len(roles) == 0 {
		panic(map[string]interface{}{
			"funcion": "GetDependenciasByRol",
			"err":     "roles es requerido y debe contener al menos un rol",
			"status":  "400",
		})
	}

	alcance, errMap := services.ResolveAlcanceUsuario(numeroDocumento, roles)
	if errMap != nil {
		panic(errMap)
	}

	c.Ctx.Output.SetStatus(200)
	c.Data["json"] = map[string]interface{}{
		"Success": true,
		"Status":  200,
		"Message": "Alcance del usuario obtenido con éxito",
		"Data":    alcance,
	}
	c.ServeJSON()
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

	numeroDocumento := strings.TrimSpace(c.GetString("numero_documento"))
	rolesRaw := c.GetString("roles")
	roles := parseRolesParam(rolesRaw)
	filtro := models.Filtro{
		NumeroResolucion: firstNonEmpty(c.GetString("NumeroResolucion"), c.GetString("numero_resolucion")),
		Vigencia:         strings.TrimSpace(c.GetString("vigencia")),
		Periodo:          strings.TrimSpace(c.GetString("Periodo")),
		Semanas:          strings.TrimSpace(c.GetString("Semanas")),
		NivelAcademico:   strings.TrimSpace(c.GetString("NivelAcademico")),
		Dedicacion:       strings.TrimSpace(c.GetString("Dedicacion")),
		Estado:           strings.TrimSpace(c.GetString("Estado")),
		TipoResolucion:   strings.TrimSpace(c.GetString("TipoResolucion")),
		ExcluirTipo:      strings.TrimSpace(c.GetString("ExcluirTipo")),
	}

	vigencia, errVigencia := c.GetInt("vigencia")
	idOikos, errIdOikos := c.GetInt("id_oikos")
	limit, errLimit := c.GetInt("limit")
	offset, errOffset := c.GetInt("offset")

	var dependenciaFiltro *int

	if numeroDocumento == "" {
		panic(map[string]interface{}{
			"funcion": "GetResolucionesByDependencia",
			"err":     "numero_documento es requerido",
			"status":  "400",
		})
	}

	if len(roles) == 0 {
		panic(map[string]interface{}{
			"funcion": "GetResolucionesByDependencia",
			"err":     "roles es requerido y debe contener al menos un rol",
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

	if errLimit != nil || limit <= 0 {
		panic(map[string]interface{}{
			"funcion": "GetResolucionesByDependencia",
			"err":     "limit es requerido y debe ser válido",
			"status":  "400",
		})
	}

	if errOffset != nil || offset <= 0 {
		panic(map[string]interface{}{
			"funcion": "GetResolucionesByDependencia",
			"err":     "offset es requerido y debe ser válido",
			"status":  "400",
		})
	}

	if errIdOikos == nil && idOikos > 0 {
		dependenciaFiltro = &idOikos
	}

	filtro.Limit = c.GetString("limit")
	filtro.Offset = c.GetString("offset")
	filtro.Vigencia = strings.TrimSpace(c.GetString("vigencia"))

	resoluciones, total, errMap := services.GetResolucionesTablaByAlcance(
		numeroDocumento,
		roles,
		filtro,
		dependenciaFiltro,
	)
	if errMap != nil {
		panic(errMap)
	}

	c.Ctx.Output.SetStatus(200)
	c.Data["json"] = map[string]interface{}{
		"Success": true,
		"Status":  200,
		"Message": "Resoluciones obtenidas con éxito",
		"Data":    resoluciones,
		"Total":   total,
	}
	c.ServeJSON()
}
