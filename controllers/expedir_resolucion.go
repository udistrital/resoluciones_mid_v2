package controllers

import (
	"encoding/json"

	"github.com/astaxie/beego"
	"github.com/udistrital/resoluciones_mid_v2/helpers"
	"github.com/udistrital/resoluciones_mid_v2/models"
)

// ExpedirResolucionController operations for ExpedirResolucion
type ExpedirResolucionController struct {
	beego.Controller
}

// URLMapping ...
func (c *ExpedirResolucionController) URLMapping() {
	c.Mapping("Expedir", c.Expedir)
	c.Mapping("ValidarDatosExpedicion", c.ValidarDatosExpedicion)
	c.Mapping("ExpedirModificacion", c.ExpedirModificacion)
	c.Mapping("Cancelar", c.Cancelar)
}

// Expedir ...
// @Title Expedir
// @Description create Expedir
// @Param	body		body 	models.ExpedicionResolucion	true		"body for Expedicion Resolucion content"
// @Success 201 {object} models.ExpedicionResolucion
// @Failure 400 bad request
// @Failure 500 Internal server error
// @router /expedir [post]
func (c *ExpedirResolucionController) Expedir() {
	defer helpers.ErrorController(c.Controller, "ExpedirResolucionController")

	if v, e := helpers.ValidarBody(c.Ctx.Input.RequestBody); !v || e != nil {
		panic(map[string]interface{}{"funcion": "Expedir", "err": helpers.ErrorBody, "status": "400"})
	}

	var m models.ExpedicionResolucion

	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &m); err == nil {
		if err := helpers.ExpedirResolucion(m); err == nil {
			c.Ctx.Output.SetStatus(201)
			c.Data["json"] = map[string]interface{}{"Success": true, "Status": "201", "Message": "Successful", "Data": m}
		} else {
			panic(err)
		}
	} else { //If 13 - Unmarshal
		panic(map[string]interface{}{"funcion": "Expedir", "err": err.Error(), "status": "400"})
	}
	c.ServeJSON()
}

// ExpedirResolucionController ...
// @Title ValidarDatosExpedicion
// @Description create ValidarDatosExpedicion
// @Param	body		body 	[]models.ExpedicionResolucion	true		"body for Validar Datos Expedición content"
// @Success 201 {object} string OK
// @Failure 400 bad request
// @Failure 500 Internal server error
// @router /validar_datos_expedicion [post]
func (c *ExpedirResolucionController) ValidarDatosExpedicion() {
	defer helpers.ErrorController(c.Controller, "ExpedirResolucionController")

	if v, e := helpers.ValidarBody(c.Ctx.Input.RequestBody); !v || e != nil {
		panic(map[string]interface{}{"funcion": "ValidarDatosExpedicion", "err": helpers.ErrorBody, "status": "400"})
	}

	var m models.ExpedicionResolucion

	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &m); err == nil {
		if err := helpers.ValidarDatosExpedicion(m); err == nil {
			c.Ctx.Output.SetStatus(201)
			c.Data["json"] = map[string]interface{}{"Success": true, "Status": "201", "Message": "Successful", "Data": "OK"}
		} else {
			panic(err)
		}
	} else {
		panic(map[string]interface{}{"funcion": "ValidarDatosExpedicion", "err": err.Error(), "status": "400"})
	}
	c.ServeJSON()
}

// ExpedirModificacion ...
// @Title ExpedirModificacion
// @Description create ExpedirModificacion
// @Param	body		body 	models.ExpedicionResolucion	true		"body for Validar Datos Expedición content"
// @Success 201 {object} models.ExpedicionResolucion
// @Failure 400 bad request
// @Failure 500 Internal server error
// @router /expedirModificacion [post]
func (c *ExpedirResolucionController) ExpedirModificacion() {
	defer helpers.ErrorController(c.Controller, "ExpedirResolucionController")

	if v, e := helpers.ValidarBody(c.Ctx.Input.RequestBody); !v || e != nil {
		panic(map[string]interface{}{"funcion": "ExpedirModificacion", "err": helpers.ErrorBody, "status": "400"})
	}

	var m models.ExpedicionResolucion
	// If 13 - Unmarshal
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &m); err == nil {
		if err := helpers.ExpedirModificacion(m); err == nil {
			c.Ctx.Output.SetStatus(201)
			c.Data["json"] = map[string]interface{}{"Success": true, "Status": "201", "Message": "Successful", "Data": m}
		} else {
			panic(err)
		}
	} else { //If 13 - Unmarshal
		panic(map[string]interface{}{"funcion": "ExpedirModificacion", "err": err.Error(), "status": "400"})
	}
	c.ServeJSON()
}

// Cancelar ...
// @Title Cancelar
// @Description create Cancelar
// @Param	body		body 	models.ExpedicionCancelacion	true		"body for Expedición a cancelar content"
// @Success 201 {object} models.ExpedicionCancelacion
// @Failure 400 bad request
// @Failure 500 Internal server error
// @router /cancelar [post]
func (c *ExpedirResolucionController) Cancelar() {
	defer helpers.ErrorController(c.Controller, "ExpedirResolucionController")

	if v, e := helpers.ValidarBody(c.Ctx.Input.RequestBody); !v || e != nil {
		panic(map[string]interface{}{"funcion": "Cancelar", "err": helpers.ErrorBody, "status": "400"})
	}

	var m models.ExpedicionCancelacion

	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &m); err == nil {
		if err := helpers.ExpedirCancelacion(m); err == nil {
			c.Ctx.Output.SetStatus(201)
			c.Data["json"] = map[string]interface{}{"Success": true, "Status": "201", "Message": "Successful", "Data": m}
		} else {
			panic(err)
		}
	} else { //If 13 - Unmarshal
		panic(map[string]interface{}{"funcion": "Cancelar", "err": err.Error(), "status": "400"})
	}
	c.ServeJSON()
}
