package controllers

import (
	"encoding/json"
	"strconv"

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

	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &m); err == nil {
		if idResolucion, err2 := helpers.InsertarResolucion(m); err2 == nil {
			c.Ctx.Output.SetStatus(201)
			c.Data["json"] = map[string]interface{}{"Success": true, "Status": 201, "Message": "Resolución insertada con exito", "Data": idResolucion}
		} else {
			panic(err)
		}
	} else {
		panic(map[string]interface{}{"funcion": "Post", "err": err.Error(), "status": "400"})
	}
	c.ServeJSON()
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

	idStr := c.Ctx.Input.Param(":id")
	id, err := strconv.Atoi(idStr)

	if err != nil || id <= 0 {
		panic(map[string]interface{}{"funcion": "GetOne", "err": helpers.ErrorParametros, "status": "400"})
	}

	if r, err2 := helpers.CargarResolucionCompleta(id); err2 == nil {
		c.Ctx.Output.SetStatus(200)
		c.Data["json"] = map[string]interface{}{"Success": true, "Status": 200, "Message": "Resolución cargada con exito", "Data": r}
	} else {
		panic(err2)
	}
	c.ServeJSON()
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

	var f models.Filtro
	var limit, offset int
	var err1, err2, err3, err4, err5, err6, err7 error

	f.Limit = c.GetString("limit")
	f.Offset = c.GetString("offset")
	f.NumeroResolucion = c.GetString("NumeroResolucion")
	f.Vigencia = c.GetString("Vigencia")
	f.Periodo = c.GetString("Periodo")
	f.Semanas = c.GetString("Semanas")
	f.FacultadId = c.GetString("Facultad")
	f.NivelAcademico = c.GetString("NivelAcademico")
	f.Dedicacion = c.GetString("Dedicacion")
	f.Estado = c.GetString("Estado")
	f.TipoResolucion = c.GetString("TipoResolucion")
	f.ExcluirTipo = c.GetString("ExcluirTipo")

	if len(f.Limit) > 0 {
		limit, err1 = strconv.Atoi(f.Limit)
	}
	if len(f.Offset) > 0 {
		offset, err2 = strconv.Atoi(f.Offset)
	}
	if len(f.NumeroResolucion) > 0 {
		_, err3 = strconv.Atoi(f.NumeroResolucion)
	}
	if len(f.Vigencia) > 0 {
		_, err4 = strconv.Atoi(f.Vigencia)
	}
	if len(f.Periodo) > 0 {
		_, err5 = strconv.Atoi(f.Periodo)
	}
	if len(f.Semanas) > 0 {
		_, err6 = strconv.Atoi(f.Semanas)
	}
	if len(f.FacultadId) > 0 {
		_, err7 = strconv.Atoi(f.FacultadId)
	}

	if err1 != nil || err2 != nil || err3 != nil || err4 != nil || err5 != nil || err6 != nil || err7 != nil || limit <= 0 || offset <= 0 {
		panic(map[string]interface{}{"funcion": "GetAll", "err": helpers.ErrorParametros, "status": "400"})
	}

	if l, t, err := helpers.ListarResolucionesFiltradas(f); err == nil {
		c.Ctx.Output.SetStatus(200)
		c.Data["json"] = map[string]interface{}{"Success": true, "Status": 200, "Message": helpers.CargaResExito, "Total": t, "Data": l}
	} else {
		panic(err)
	}
	c.ServeJSON()
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

	idStr := c.Ctx.Input.Param(":id")
	_, err := strconv.Atoi(idStr)

	if err != nil {
		panic(map[string]interface{}{"funcion": "Put", "err": helpers.ErrorParametros, "status": "400"})
	}

	var r models.ContenidoResolucion

	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &r); err == nil {
		if err := helpers.ActualizarResolucionCompleta(r); err == nil {
			c.Ctx.Output.SetStatus(200)
			c.Data["json"] = map[string]interface{}{"Success": true, "Status": 200, "Message": "Resolución actualizada con exito", "Data": r}
		} else {
			panic(err)
		}
	} else {
		panic(map[string]interface{}{"funcion": "Put", "err": err.Error(), "status": "400"})
	}
	c.ServeJSON()
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

	idStr := c.Ctx.Input.Param(":id")
	id, err := strconv.Atoi(idStr)

	if err != nil {
		panic(map[string]interface{}{"funcion": "Delete", "err": helpers.ErrorParametros, "status": "400"})
	}

	if err2 := helpers.AnularResolucion(id); err == nil {
		c.Ctx.Output.SetStatus(200)
		d := map[string]interface{}{"Id": id}
		c.Data["json"] = map[string]interface{}{"Success": true, "Status": 200, "Message": "Resolución anulada con exito", "Data": d}
	} else {
		panic(err2)
	}
	c.ServeJSON()
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

	idStr := c.Ctx.Input.Param(":id")
	id, err := strconv.Atoi(idStr)

	if err != nil || id <= 0 {
		panic(map[string]interface{}{"funcion": "ConsultaDocente", "err": helpers.ErrorParametros, "status": "400"})
	}

	if r, err2 := helpers.ConsultaDocente(id); err2 == nil {
		c.Ctx.Output.SetStatus(200)
		c.Data["json"] = map[string]interface{}{"Success": true, "Status": 200, "Message": helpers.CargaResExito, "Data": r}
	} else {
		panic(err2)
	}
	c.ServeJSON()
}

// GenerarResolucion ...
// @Title GenerarResolucion
// @Description Genera el documento PDF de la resolución
// @Param	id		path 	string	true		"id de la resolución"
// @Success 200 {string} string Base64 encoded file
// @Failure 400 bad request
// @Failure 500 Internal server error
// @router /generar_resolucion/:id [get]
func (c *GestionResolucionesController) GenerarResolucion() {
	defer helpers.ErrorController(c.Controller, "GestionResolucionesController")

	idStr := c.Ctx.Input.Param(":id")
	id, err := strconv.Atoi(idStr)

	if err != nil || id <= 0 {
		panic(map[string]interface{}{"funcion": "GenerarResolucion", "err": helpers.ErrorParametros, "status": "400"})
	}

	if r, err2 := helpers.GenerarResolucion(id); err2 == nil {
		c.Ctx.Output.SetStatus(200)
		c.Data["json"] = map[string]interface{}{"Success": true, "Status": 200, "Message": helpers.CargaResExito, "Data": r}
	} else {
		panic(err2)
	}
	c.ServeJSON()
}

// ActualizarEstado ...
// @Title ActualizarEstado
// @Description Modifica el estado de una resolución
// @Param	body		body 	models.NuevoEstadoResolucion	true		"body for NuevoEstadoResolucion content"
// @Success 201 {string} string		Nuevo estado de la resolución
// @Failure 400 bad request
// @Failure 500 Internal server error
// @router /actualizar_estado [post]
func (c *GestionResolucionesController) ActualizarEstado() {
	defer helpers.ErrorController(c.Controller, "GestionResolucionesController")

	var m models.NuevoEstadoResolucion

	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &m); err == nil {
		if err2 := helpers.CambiarEstadoResolucion(m.ResolucionId, m.Estado, m.Usuario); err2 == nil {
			c.Ctx.Output.SetStatus(201)
			c.Data["json"] = map[string]interface{}{"Success": true, "Status": 201, "Message": "Resolución insertada con exito", "Data": "OK"}
		} else {
			panic(err)
		}
	} else {
		panic(map[string]interface{}{"funcion": "ActualizarEstado", "err": err.Error(), "status": "400"})
	}
	c.ServeJSON()
}
