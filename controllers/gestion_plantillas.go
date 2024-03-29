package controllers

import (
	"encoding/json"
	"strconv"

	"github.com/astaxie/beego"
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

	if v, e := helpers.ValidarBody(c.Ctx.Input.RequestBody); !v || e != nil {
		panic(map[string]interface{}{"funcion": "Post", "err": helpers.ErrorBody, "status": "400"})
	}

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
// @Param	id		path 	string	true		"Id de la plantilla a consultar"
// @Success 200 {object} models.ContenidoResolucion
// @Failure 400 bad request
// @Failure 500 Internal server error
// @router /:id [get]
func (c *GestionPlantillasController) GetOne() {
	defer helpers.ErrorController(c.Controller, "GestionPlantillasController")

	idStr := c.Ctx.Input.Param(":id")
	id, err := strconv.Atoi(idStr)

	if err != nil || id <= 0 {
		panic(map[string]interface{}{"funcion": "GetOne", "err": helpers.ErrorParametros, "status": "400"})
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
// @Success 200 {object} []models.ContenidoResolucion
// @Failure 400 bad request
// @Failure 500 Internal server error
// @router / [get]
func (c *GestionPlantillasController) GetAll() {
	defer helpers.ErrorController(c.Controller, "GestionPlantillasController")

	if l, err := helpers.ListarPlantillas(); err == nil {
		c.Ctx.Output.SetStatus(200)
		c.Data["json"] = map[string]interface{}{"Success": true, "Status": 200, "Message": "Plantillas consultadas con exito", "Data": l}
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
	defer helpers.ErrorController(c.Controller, "GestionPlantillasController")

	idStr := c.Ctx.Input.Param(":id")
	_, err := strconv.Atoi(idStr)

	if err != nil {
		panic(map[string]interface{}{"funcion": "Put", "err": helpers.ErrorParametros, "status": "400"})
	}

	if v, e := helpers.ValidarBody(c.Ctx.Input.RequestBody); !v || e != nil {
		panic(map[string]interface{}{"funcion": "Put", "err": helpers.ErrorBody, "status": "400"})
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
// @Description post the fehcaFin
// @Param	fecha_inicio		path 	string	true		"Fecha de inicio"
// @Param	numerosemanas		path 	string	true		"Numero de semanas"
// @Success 200 {object} int Id de la resolucion anulada
// @Failure 400 bad request
// @Failure 500 Internal server error
// @router /:id [delete]
func (c *GestionPlantillasController) Delete() {
	defer helpers.ErrorController(c.Controller, "GestionPlantillasController")

	idStr := c.Ctx.Input.Param(":id")
	id, err := strconv.Atoi(idStr)

	if err != nil {
		panic(map[string]interface{}{"funcion": "Delete", "err": helpers.ErrorParametros, "status": "400"})
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

	if v, e := helpers.ValidarBody(c.Ctx.Input.RequestBody); !v || e != nil {
		panic(map[string]interface{}{"funcion": "Post", "err": helpers.ErrorBody, "status": "400"})
	}

	var m models.FechaFin

	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &m); err == nil {
		if fechaFin := helpers.CalcularFechaFin(m.FechaInicio, m.NumeroSemanas); err == nil {
			c.Ctx.Output.SetStatus(201)
			c.Data["json"] = map[string]interface{}{"Success": true, "Status": 201, "Message": "Plantilla insertada con exito", "Data": fechaFin}
		} else {
			panic(err)
		}
	} else {
		panic(map[string]interface{}{"funcion": "Post", "err": err.Error(), "status": "400"})
	}
	c.ServeJSON()
}
