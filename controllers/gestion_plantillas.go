package controllers

import (
	"encoding/json"
	"strconv"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/udistrital/resoluciones_mid_v2/helpers"
	"github.com/udistrital/resoluciones_mid_v2/models"
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
	defer ErrorController(c)

	var m models.ContenidoResolucion

	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &m); err == nil {
		if idPlantilla, err := helpers.InsertarPlantilla(m); err == nil {
			c.Ctx.Output.SetStatus(201)
			c.Data["json"] = map[string]interface{}{"Success": true, "Status": 201, "Message": "Plantilla insertada con exito", "Data": idPlantilla}
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
// @Description get GestionPlantillas by id
// @Param	id		path 	string	true		"The key for staticblock"
// @Success 200 {object} models.ContenidoResolucion
// @Failure 400 bad request
// @Failure 500 Internal server error
// @router /:id [get]
func (c *GestionPlantillasController) GetOne() {
	defer ErrorController(c)

	idStr := c.Ctx.Input.Param(":id")
	id, err := strconv.Atoi(idStr)

	if err != nil {
		panic(map[string]interface{}{"funcion": "GetOne", "err": "Error en los parametros de ingreso", "status": "400"})
	}

	if p, err2 := helpers.CargarPlantilla(id); err2 == nil {
		c.Ctx.Output.SetStatus(200)
		c.Data["json"] = map[string]interface{}{"Success": true, "Status": 200, "Message": "Plantilla cargada con exito", "Data": p}
	} else {
		panic(err2)
	}
	c.ServeJSON()
}

// GetAll ...
// @Title GetAll
// @Description get GestionPlantillas
// @Success 200 {object} models.ContenidoResolucion
// @Failure 400 bad request
// @Failure 500 Internal server error
// @router / [get]
func (c *GestionPlantillasController) GetAll() {
	defer ErrorController(c)

	if l, err := helpers.ListarPlantillas(); err == nil {
		c.Ctx.Output.SetStatus(200)
		c.Data["json"] = map[string]interface{}{"Success": true, "Status": 200, "Message": "Plantilla cargada con exito", "Data": l}
	} else {
		panic(err)
	}
	c.ServeJSON()
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
	defer ErrorController(c)

	idStr := c.Ctx.Input.Param(":id")
	_, err := strconv.Atoi(idStr)

	if err != nil {
		panic(map[string]interface{}{"funcion": "Put", "err": "Error en los parametros de ingreso", "status": "400"})
	}

	var m models.ContenidoResolucion

	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &m); err == nil {
		if err := helpers.ActualizarPlantilla(m); err == nil {
			c.Ctx.Output.SetStatus(200)
			c.Data["json"] = map[string]interface{}{"Success": true, "Status": 200, "Message": "Plantilla actualizada con exito", "Data": m}
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
// @Description delete the GestionPlantillas
// @Param	id		path 	string	true		"The id you want to delete"
// @Success 200 {string} delete success!
// @Failure 400 bad request
// @Failure 500 Internal server error
// @router /:id [delete]
func (c *GestionPlantillasController) Delete() {
	defer ErrorController(c)

	idStr := c.Ctx.Input.Param(":id")
	id, err := strconv.Atoi(idStr)

	if err != nil {
		panic(map[string]interface{}{"funcion": "Delete", "err": "Error en los parametros de ingreso", "status": "400"})
	}

	if err2 := helpers.BorrarPlantilla(id); err == nil {
		c.Ctx.Output.SetStatus(200)
		d := map[string]interface{}{"Id": id}
		c.Data["json"] = map[string]interface{}{"Success": true, "Status": 200, "Message": "Plantilla eliminada con exito", "Data": d}
	} else {
		panic(err2)
	}
	c.ServeJSON()
}

func ErrorController(c *GestionPlantillasController) {
	if err := recover(); err != nil {
		logs.Error(err)
		localError := err.(map[string]interface{})
		c.Data["mesaage"] = (beego.AppConfig.String("appname") + "/" + "GestionPlantillasController" + "/" + (localError["funcion"]).(string))
		c.Data["data"] = (localError["err"])
		if status, ok := localError["status"]; ok {
			c.Abort(status.(string))
		} else {
			c.Abort("500")
		}
	}
}
