package controllers

import (
	"github.com/astaxie/beego"
	"github.com/udistrital/resoluciones_mid_v2/helpers"
	"github.com/udistrital/resoluciones_mid_v2/models"
)

// GestionResolucionesController operations for Gestionresoluciones
type GestionResolucionesController struct {
	beego.Controller
}

// URLMapping ...
func (c *GestionResolucionesController) URLMapping() {
	c.Mapping("Post", c.Post)
	c.Mapping("GetOne", c.GetOne)
	c.Mapping("GetAll", c.GetAll)
	c.Mapping("Put", c.Put)
	c.Mapping("Delete", c.Delete)
	c.Mapping("ConsultaDocente", c.ConsultaDocente)
	c.Mapping("GenerarResolucion", c.GenerarResolucion)
}

// Post ...
// @Title Create
// @Description crea una nueva resolución basado en una plantilla
// @Param	body		body 	models.ContenidoResolucion	true		"body for ContenidoResolucion content"
// @Success 201 {object} int Id de la nueva resolución
// @Failure 400 bad request
// @Failure 500 Internal server error
// @router / [post]
func (c *GestionResolucionesController) Post() {
	defer helpers.ErrorController(c.Controller, "GestionResolucionesController")

	var m models.ContenidoResolucion
	decodeJSONBody(c.Ctx.Input.RequestBody, &m, "Post")

	if idResolucion, err2 := helpers.InsertarResolucion(m); err2 == nil {
		writeJSON(&c.Controller, 201, "Resolución insertada con exito", idResolucion, nil)
	} else {
		panic(err2)
	}
}

// GetOne ...
// @Title GetOne
// @Description get one ContenidoResolucion by id
// @Param	id		path 	string	true		"The key for staticblock"
// @Success 200 {object} models.ContenidoResolucion
// @Failure 400 bad request
// @Failure 500 Internal server error
// @router /:id [get]
func (c *GestionResolucionesController) GetOne() {
	defer helpers.ErrorController(c.Controller, "GestionResolucionesController")

	id := parsePositivePathID(&c.Controller, ":id", "GetOne")

	if r, err2 := helpers.CargarResolucionCompleta(id); err2 == nil {
		writeJSON(&c.Controller, 200, "Resolución cargada con exito", r, nil)
	} else {
		panic(err2)
	}
}

// GetAll ...
// @Title GetAll
// @Description Carga las resoluciones de acuerdo a los parametros recibidos
// @Param	limit			query	int		true	"Limit the size of result set. Must be an integer"
// @Param	offset			query	int		true	"Start position of result set. Must be an integer"
// @Param	NumeroResolucion	query	int	false	"Numero de resolución a buscar"
// @Param	Vigencia		query	int		false	"Año de la resolución a buscar"
// @Param	Periodo			query	int		false	"Periodo academico a buscar"
// @Param	Semanas			query	int		false	"Numero de semanas por las que se ha hecho la resolución"
// @Param	Facultad		query	string	false	"Facultad que emite la resolución"
// @Param	NivelAcademico	query	string	false	"Nivel academico de la resolución. PREGRADO o POSGRADO"
// @Param	Dedicacion		query	string	false	"Dedicación del docente"
// @Param	Estado			query	string	false	"Estado de la resolución"
// @Param	TipoResolucion	query	string	false	"Tipo de resolución a bucar"
// @Param	ExcluirTipo		query	string	false	"Tipo de resolución a excluir de la consulta"
// @Success 200 {object} []models.Resoluciones
// @Failure 400 bad request
// @Failure 500 Internal server error
// @router / [get]
func (c *GestionResolucionesController) GetAll() {
	defer helpers.ErrorController(c.Controller, "GestionResolucionesController")

	f := buildFiltroConsulta(c, "Vigencia")
	if err := validateFiltroConsulta(f); err != nil {
		panic(badRequest("GetAll", err))
	}

	if l, t, err := helpers.ListarResolucionesFiltradas(f); err == nil {
		writeJSON(&c.Controller, 200, helpers.CargaResExito, l, map[string]interface{}{"Total": t})
	} else {
		panic(err)
	}
}

// Put ...
// @Title Put
// @Description update the ContenidoResolucion
// @Param	id		path 	string	true		"The id you want to update"
// @Param	body		body 	models.ContenidoResolucion	true		"body for Gestionresoluciones content"
// @Success 200 {object} models.ContenidoResolucion
// @Failure 400 bad request
// @Failure 500 Internal server error
// @router /:id [put]
func (c *GestionResolucionesController) Put() {
	defer helpers.ErrorController(c.Controller, "GestionResolucionesController")

	parsePositivePathID(&c.Controller, ":id", "Put")

	var r models.ContenidoResolucion
	decodeJSONBody(c.Ctx.Input.RequestBody, &r, "Put")

	if err := helpers.ActualizarResolucionCompleta(r); err == nil {
		writeJSON(&c.Controller, 200, "Resolución actualizada con exito", r, nil)
	} else {
		panic(err)
	}
}

// Delete ...
// @Title Delete
// @Description delete the Gestionresoluciones
// @Param	id		path 	string	true		"The id you want to delete"
// @Success 200 {object} int Id de la resolucion anulada
// @Failure 400 bad request
// @Failure 500 Internal server error
// @router /:id [delete]
func (c *GestionResolucionesController) Delete() {
	defer helpers.ErrorController(c.Controller, "GestionResolucionesController")

	id := parsePositivePathID(&c.Controller, ":id", "Delete")

	if err2 := helpers.AnularResolucion(id); err2 == nil {
		d := map[string]interface{}{"Id": id}
		writeJSON(&c.Controller, 200, "Resolución anulada con exito", d, nil)
	} else {
		panic(err2)
	}
}

// ConsultaDocente ...
// @Title ConsultaDocente
// @Description get Resoluciones by id del docente
// @Param	id		path 	string	true		"id del docente"
// @Success 200 {object} []models.Resoluciones
// @Failure 400 bad request
// @Failure 500 Internal server error
// @router /consultar_docente/:id [get]
func (c *GestionResolucionesController) ConsultaDocente() {
	defer helpers.ErrorController(c.Controller, "GestionResolucionesController")

	id := parsePositivePathID(&c.Controller, ":id", "ConsultaDocente")

	if r, err2 := helpers.ConsultaDocente(id); err2 == nil {
		writeJSON(&c.Controller, 200, helpers.CargaResExito, r, nil)
	} else {
		panic(err2)
	}
}

// GenerarResolucion ...
// @Title GenerarResolucion
// @Description Genera el documento PDF de la resolución
// @Param	id		path 	string	true		"id de la resolución"
// @Success 200 {object} string Base64 encoded file
// @Failure 400 bad request
// @Failure 500 Internal server error
// @router /generar_resolucion/:id [get]
func (c *GestionResolucionesController) GenerarResolucion() {
	defer helpers.ErrorController(c.Controller, "GestionResolucionesController")

	id := parsePositivePathID(&c.Controller, ":id", "GenerarResolucion")

	if r, err2 := helpers.GenerarResolucion(id); err2 == nil {
		writeJSON(&c.Controller, 200, helpers.CargaResExito, r, nil)
	} else {
		panic(err2)
	}
}

// ActualizarEstado ...
// @Title ActualizarEstado
// @Description Modifica el estado de una resolución
// @Param	body		body 	models.NuevoEstadoResolucion	true		"body for NuevoEstadoResolucion content"
// @Success 201 {object} string		Nuevo estado de la resolución
// @Failure 400 bad request
// @Failure 500 Internal server error
// @router /actualizar_estado [post]
func (c *GestionResolucionesController) ActualizarEstado() {
	defer helpers.ErrorController(c.Controller, "GestionResolucionesController")

	var m models.NuevoEstadoResolucion
	decodeJSONBody(c.Ctx.Input.RequestBody, &m, "ActualizarEstado")

	if err2 := helpers.CambiarEstadoResolucion(m.ResolucionId, m.Estado, m.Usuario); err2 == nil {
		writeJSON(&c.Controller, 201, "Resolución insertada con exito", "OK", nil)
	} else {
		panic(err2)
	}
}
