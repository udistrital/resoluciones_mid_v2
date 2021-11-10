package controllers

import (
	"github.com/astaxie/beego"
)

// GestionResolucionesController operations for GestionResoluciones
type GestionResolucionesController struct {
	beego.Controller
}

// URLMapping ...
func (c *GestionResolucionesController) URLMapping() {
	c.Mapping("Post", c.Post)
	c.Mapping("GetAll", c.GetAll)
	c.Mapping("Put", c.Put)
}

// Post ...
// @Title Create
// @Description create GestionResoluciones
// @Param	body		body 	models.GestionResoluciones	true		"body for GestionResoluciones content"
// @Success 201 {object} models.GestionResoluciones
// @Failure 403 body is empty
// @router / [post]
func (c *GestionResolucionesController) Post() {

}

// GetAll ...
// @Title GetAll
// @Description get GestionResoluciones
// @Param	query	query	string	false	"Filter. e.g. col1:v1,col2:v2 ..."
// @Param	fields	query	string	false	"Fields returned. e.g. col1,col2 ..."
// @Param	sortby	query	string	false	"Sorted-by fields. e.g. col1,col2 ..."
// @Param	order	query	string	false	"Order corresponding to each sortby field, if single value, apply to all sortby fields. e.g. desc,asc ..."
// @Param	limit	query	string	false	"Limit the size of result set. Must be an integer"
// @Param	offset	query	string	false	"Start position of result set. Must be an integer"
// @Success 200 {object} models.GestionResoluciones
// @Failure 403
// @router / [get]
func (c *GestionResolucionesController) GetAll() {

}

// Put ...
// @Title Put
// @Description update the GestionResoluciones
// @Param	id		path 	string	true		"The id you want to update"
// @Param	body		body 	models.GestionResoluciones	true		"body for GestionResoluciones content"
// @Success 200 {object} models.GestionResoluciones
// @Failure 403 :id is not int
// @router /:id [put]
func (c *GestionResolucionesController) Put() {

}
