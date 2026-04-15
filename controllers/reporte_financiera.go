package controllers

import (
	"github.com/astaxie/beego"
	"github.com/udistrital/resoluciones_mid_v2/helpers"
	"github.com/udistrital/resoluciones_mid_v2/models"
)

// Reporte_financieraController operations for Reporte_financiera
type ReporteFinancieraController struct {
	beego.Controller
}

// URLMapping ...
func (c *ReporteFinancieraController) URLMapping() {
	c.Mapping("ModificarVinculacion", c.ReporteFinanciera)
}

// ReporteFinanciera ...
// @Title ReporteFinanciera
// @Description Genera el reporte para financiera
// @Param	body		body 	models.DatosReporte	true		"body for DatosReporte content"
// @Success 201 {object} models.ReporteFinancieraFinal "Reporte Financiera generado con exito"
// @Failure 400 bad request
// @Failure 500 Internal server error
// @router / [post]
func (c *ReporteFinancieraController) ReporteFinanciera() {
	defer helpers.ErrorController(c.Controller, "ReporteFinancieraController")

	var p models.DatosReporte
	decodeJSONBody(c.Ctx.Input.RequestBody, &p, "ReporteFinanciera")

	if r, err2 := helpers.ReporteFinanciera(p); err2 == nil {
		writeJSON(&c.Controller, 201, "Reporte Generado con Exito", r, nil)
	} else {
		panic(err2)
	}
}
