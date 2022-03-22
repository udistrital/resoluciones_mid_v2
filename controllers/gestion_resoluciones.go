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
	c.Mapping("GetResolucionesExpedidas", c.GetResolucionesExpedidas)
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

	if err != nil {
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
// @Description get all ContenidoResolucion
// @Param	query	query	string	false	"Filter. e.g. col1:v1,col2:v2 ..."
// @Param	fields	query	string	false	"Fields returned. e.g. col1,col2 ..."
// @Param	sortby	query	string	false	"Sorted-by fields. e.g. col1,col2 ..."
// @Param	order	query	string	false	"Order corresponding to each sortby field, if single value, apply to all sortby fields. e.g. desc,asc ..."
// @Param	limit	query	string	false	"Limit the size of result set. Must be an integer"
// @Param	offset	query	string	false	"Start position of result set. Must be an integer"
// @Success 200 {object} []models.Resoluciones
// @Failure 400 bad request
// @Failure 500 Internal server error
// @router / [get]
func (c *GestionResolucionesController) GetAll() {
	defer helpers.ErrorController(c.Controller, "GestionResolucionesController")

	limit, err1 := c.GetInt("limit")
	offset, err2 := c.GetInt("offset")

	if err1 != nil || err2 != nil {
		panic(map[string]interface{}{"funcion": "GetAll", "err": helpers.ErrorParametros, "status": "400"})
	}

	if l, t, err := helpers.ListarResoluciones(limit, offset); err == nil {
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

	if err != nil {
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

// GetResolucionesAprobadas ...
// @Title GetResolucionesAprobadas
// @Description get Resoluciones aprobadas
// @Success 200 {object} []models.Resoluciones
// @Failure 400 bad request
// @Failure 500 Internal server error
// @router /resoluciones_aprobadas [get]
func (c *GestionResolucionesController) GetResolucionesAprobadas() {
	defer helpers.ErrorController(c.Controller, "GestionResolucionesController")

	query := c.GetString("query")
	facultad := c.GetString("facultad")
	tipoRes := c.GetString("tipoRes")
	nivelA := c.GetString("nivelA")
	dedicacion := c.GetString("dedicacion")
	estadoRes := c.GetString("estadoRes")

	limit, err1 := c.GetInt("limit")
	offset, err2 := c.GetInt("offset")

	if err1 != nil || err2 != nil {
		panic(map[string]interface{}{"funcion": "GetResolucionesAprobadas", "err": helpers.ErrorParametros, "status": "400"})
	}
	if l, t, err := helpers.ListarResolucionesAprobadas(query, facultad, tipoRes, nivelA, dedicacion, estadoRes, limit, offset); err == nil {
		c.Ctx.Output.SetStatus(200)
		c.Data["json"] = map[string]interface{}{"Success": true, "Status": 200, "Message": helpers.CargaResExito, "Total": t, "Data": l}
	} else {
		panic(err)
	}
	c.ServeJSON()
}

// GetResolucionesExpedidas ...
// @Title GetResolucionesExpedidas
// @Description get Resoluciones expedidas
// @Success 200 {object} []models.Resoluciones
// @Failure 400 bad request
// @Failure 500 Internal server error
// @router /resoluciones_expedidas [get]
func (c *GestionResolucionesController) GetResolucionesExpedidas() {
	defer helpers.ErrorController(c.Controller, "GestionResolucionesController")

	limit, err1 := c.GetInt("limit")
	offset, err2 := c.GetInt("offset")

	if err1 != nil || err2 != nil {
		panic(map[string]interface{}{"funcion": "GetResolucionesExpedidas", "err": helpers.ErrorParametros, "status": "400"})
	}

	if l, t, err := helpers.ListarResolucionesExpedidas(limit, offset); err == nil {
		c.Ctx.Output.SetStatus(200)
		c.Data["json"] = map[string]interface{}{"Success": true, "Status": 200, "Message": helpers.CargaResExito, "Total": t, "Data": l}
	} else {
		panic(err)
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

	if err != nil {
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
