package controllers

import (
	"strconv"

	"github.com/astaxie/beego"
	"github.com/udistrital/resoluciones_mid_v2/helpers"
	"github.com/udistrital/resoluciones_mid_v2/models"
	"github.com/udistrital/resoluciones_mid_v2/services"
)

// GestionPlantillasController operations for GestionPlantillas
type GestionPlantillasController struct {
	beego.Controller
}

// URLMapping ...
func (c *GestionPlantillasController) URLMapping() {
	c.Mapping("Post", c.Post)
	c.Mapping("GetOne", c.GetOne)
	c.Mapping("GetAll", c.GetAll)
	c.Mapping("Put", c.Put)
	c.Mapping("Delete", c.Delete)
	c.Mapping("CalculoFechaFin", c.CalculoFechaFin)
}

// Post ...
// @Title Create
// @Description create GestionPlantillas
// @Param	body		body 	models.ContenidoResolucion	true		"body for GestionPlantillas content"
// @Success 201 {object} models.ContenidoResolucion
// @Failure 400 bad request
// @Failure 500 Internal server error
// @router / [post]
func (c *GestionPlantillasController) Post() {
	defer helpers.ErrorController(c.Controller, "GestionPlantillasController")

	var m models.ContenidoResolucion
	decodeJSONBody(c.Ctx.Input.RequestBody, &m, "Post")

	if idPlantilla, err := helpers.InsertarPlantilla(m); err == nil {
		writeJSON(&c.Controller, 201, "Plantilla insertada con exito", idPlantilla, nil)
	} else {
		panic(err)
	}
}

// GetOne ...
// @Title GetOne
// @Description get GestionPlantillas by id
// @Param	id		path 	string	true		"Id de la plantilla a consultar"
// @Success 200 {object} models.ContenidoResolucion
// @Failure 400 bad request
// @Failure 500 Internal server error
// @router /:id [get]
func (c *GestionPlantillasController) GetOne() {
	defer helpers.ErrorController(c.Controller, "GestionPlantillasController")

	id := parsePositivePathID(&c.Controller, ":id", "GetOne")

	if p, err2 := helpers.CargarPlantilla(id); err2 == nil {
		writeJSON(&c.Controller, 200, "Plantilla cargada con exito", p, nil)
	} else {
		panic(err2)
	}
}

// GetAll ...
// @Title GetAll
// @Description get GestionPlantillas
// @Param	numero_documento	query	string	true	"Número de documento del usuario"
// @Param	roles			query	string	true	"Roles del usuario separados por coma"
// @Param	Facultad		query	string	false	"Id de facultad para filtrar plantillas"
// @Success 200 {object} []models.ContenidoResolucion
// @Failure 400 bad request
// @Failure 500 Internal server error
// @router / [get]
func (c *GestionPlantillasController) GetAll() {
	defer helpers.ErrorController(c.Controller, "GestionPlantillasController")

	authContext := requireAuthenticatedContext(buildAuthenticatedContext(&c.Controller), "GetAll")

	facultadIdStr := c.GetString("Facultad")
	var dependenciaFiltro *int

	if facultadIdStr != "" {
		if err := validateNamedPositiveInt(facultadIdStr, "Facultad"); err != nil {
			panic(badRequest("GetAll", err))
		}
		facultadId, _ := strconv.Atoi(facultadIdStr)
		dependenciaFiltro = &facultadId
	}

	if l, err := services.GetPlantillasByAlcance(authContext.NumeroDocumento, authContext.Roles, dependenciaFiltro); err == nil {
		writeJSON(&c.Controller, 200, "Plantillas consultadas con exito", l, nil)
	} else {
		panic(err)
	}
}

// Put ...
// @Title Put
// @Description update the GestionPlantillas
// @Param	id		path 	string	true		"The id you want to update"
// @Param	body		body 	models.ContenidoResolucion	true		"body for GestionPlantillas content"
// @Success 200 {object} models.ContenidoResolucion
// @Failure 400 bad request
// @Failure 500 Internal server error
// @router /:id [put]
func (c *GestionPlantillasController) Put() {
	defer helpers.ErrorController(c.Controller, "GestionPlantillasController")

	parsePositivePathID(&c.Controller, ":id", "Put")

	var m models.ContenidoResolucion
	decodeJSONBody(c.Ctx.Input.RequestBody, &m, "Put")

	if err := helpers.ActualizarPlantilla(m); err == nil {
		writeJSON(&c.Controller, 200, "Plantilla actualizada con exito", m, nil)
	} else {
		panic(err)
	}
}

// Delete ...
// @Title Delete
// @Description post the fehcaFin
// @Param	fecha_inicio		path 	string	true		"Fecha de inicio"
// @Param	numerosemanas		path 	string	true		"Numero de semanas"
// @Success 200 {object} int Id de la resolucion anulada
// @Failure 400 bad request
// @Failure 500 Internal server error
// @router /:id [delete]
func (c *GestionPlantillasController) Delete() {
	defer helpers.ErrorController(c.Controller, "GestionPlantillasController")

	id := parsePositivePathID(&c.Controller, ":id", "Delete")

	if err2 := helpers.BorrarPlantilla(id); err2 == nil {
		d := map[string]interface{}{"Id": id}
		writeJSON(&c.Controller, 200, "Plantilla eliminada con exito", d, nil)
	} else {
		panic(err2)
	}
}

// CalculoFechaFin ...
// @Title CalculoFechaFin
// @Description calcula Fecha Fin
// @Param	body		body 	models.FechaFin	true		"body for FechaFin content"
// @Success 201 {object} models.FechaFinCalculada
// @Failure 400 bad request
// @Failure 500 Internal server error
// @router /calculo_fecha_fin [post]
func (c *GestionPlantillasController) CalculoFechaFin() {
	defer helpers.ErrorController(c.Controller, "GestionPlantillasController")

	var m models.FechaFin
	decodeJSONBody(c.Ctx.Input.RequestBody, &m, "Post")

	fechaFin := helpers.CalcularFechasContrato(m.FechaInicio, m.NumeroSemanas)
	writeJSON(&c.Controller, 201, "Fechas Calculadas con Exito", fechaFin, nil)
}
