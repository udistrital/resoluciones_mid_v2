package controllers

import (
	"encoding/json"

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
}

// Post ...
// @Title Create
// @Description crea una nueva resolución basado en una plantilla
// @Param	body		body 	models.ContenidoResolucion	true		"body for ContenidoResolucion content"
// @Success 201 {int} Id de la nueva resolución
// @Failure 400 bad request
// @Failure 500 Internal server error
// @router / [post]
func (c *GestionResolucionesController) Post() {
	defer helpers.ErrorController(c.Controller, "GestionResolucionesController")

	var m models.ContenidoResolucion

	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &m); err == nil {
		if idResolucion, err2 := helpers.InsertarResolucion(m); err2 == nil {
			c.Ctx.Output.SetStatus(201)
			c.Data["json"] = map[string]interface{}{"Success": true, "Status": 201, "Message": "Plantilla insertada con exito", "Data": idResolucion}
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
// @Success 200 {object} []models.ContenidoResolucion
// @Failure 400 bad request
// @Failure 500 Internal server error
// @router / [get]
func (c *GestionResolucionesController) GetAll() {
	defer helpers.ErrorController(c.Controller, "GestionResolucionesController")
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
}

// Delete ...
// @Title Delete
// @Description delete the Gestionresoluciones
// @Param	id		path 	string	true		"The id you want to delete"
// @Success 200 {int} id
// @Failure 400 bad request
// @Failure 500 Internal server error
// @router /:id [delete]
func (c *GestionResolucionesController) Delete() {
	defer helpers.ErrorController(c.Controller, "GestionResolucionesController")
}
