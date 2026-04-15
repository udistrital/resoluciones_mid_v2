package controllers

import (
	"github.com/astaxie/beego"
	"github.com/udistrital/resoluciones_mid_v2/helpers"
	"github.com/udistrital/resoluciones_mid_v2/models"
)

// ServicesController operations for Services
type ServicesController struct {
	beego.Controller
}

// URLMapping ...
func (c *ServicesController) URLMapping() {
	c.Mapping("DesagregadoPlaneacion", c.DesagregadoPlaneacion)

}

// Post ...
// @Title Create
// @Description Genera el detalle desagregado de salario y sus prestaciones segun los parámetros indicados
// @Param	body		body 	models.ObjetoDesagregado	true		"body for DesagregadoPlaneacion content"
// @Success 201 {object} []models.ObjetoDesagregado
// @Failure 400 bad request
// @Failure 500 Internal server error
// @router /desagregado_planeacion [post]
func (c *ServicesController) DesagregadoPlaneacion() {
	defer helpers.ErrorController(c.Controller, "ServicesController")

	var d []models.ObjetoDesagregado

	decodeJSONBody(c.Ctx.Input.RequestBody, &d, "DesagregadoPlaneacion")

	if dd, err2 := helpers.CalcularComponentesSalario(d); err2 == nil {
		writeJSON(&c.Controller, 201, "Cálculos realizados con éxito", dd, nil)
	} else {
		panic(err2)
	}
}
