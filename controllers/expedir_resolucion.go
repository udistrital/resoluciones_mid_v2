package controllers

import (
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

	var m models.ExpedicionResolucion
	decodeJSONBody(c.Ctx.Input.RequestBody, &m, "Expedir")

	if err := helpers.ExpedirResolucion(m); err == nil {
		writeJSON(&c.Controller, 201, "Successful", m, nil)
	} else {
		panic(err)
	}
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

	var m models.ExpedicionResolucion
	decodeJSONBody(c.Ctx.Input.RequestBody, &m, "ValidarDatosExpedicion")

	if err := helpers.ValidarDatosExpedicion(m); err == nil {
		writeJSON(&c.Controller, 201, "Successful", "OK", nil)
	} else {
		panic(err)
	}
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

	var m models.ExpedicionResolucion
	decodeJSONBody(c.Ctx.Input.RequestBody, &m, "ExpedirModificacion")

	if err := helpers.ExpedirModificacion(m); err == nil {
		writeJSON(&c.Controller, 201, "Successful", m, nil)
	} else {
		panic(err)
	}
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

	var m models.ExpedicionCancelacion
	decodeJSONBody(c.Ctx.Input.RequestBody, &m, "Cancelar")

	if err := helpers.ExpedirCancelacion(m); err == nil {
		writeJSON(&c.Controller, 201, "Successful", m, nil)
	} else {
		panic(err)
	}
}
