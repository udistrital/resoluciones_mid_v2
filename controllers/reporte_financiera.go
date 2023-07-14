package controllers

import (
	"encoding/json"

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

	if v, e := helpers.ValidarBody(c.Ctx.Input.RequestBody); !v || e != nil {
		panic(map[string]interface{}{"funcion": "ReporteFinanciera", "err": helpers.ErrorBody, "status": "400"})
	}

	var p models.DatosReporte

	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &p); err == nil {
		if r, err2 := helpers.ReporteFinanciera(p); err == nil {
			c.Ctx.Output.SetStatus(201)
			c.Data["json"] = map[string]interface{}{"Success": true, "Status": 201, "Message": "Reporte Generado con Exito", "Data": r}
		} else {
			panic(err2)
		}
	} else {
		panic(map[string]interface{}{"funcion": "ReporteFinanciera", "err": err.Error(), "status": "400"})
	}
	c.ServeJSON()
}
