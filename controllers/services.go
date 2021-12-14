package controllers

import (
	"encoding/json"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
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
// @Success 201 {object} models.ObjetoDesagregado
// @Failure 400 bad request
// @Failure 500 Internal server error
// @router /desagregado_planeacion [post]
func (c *ServicesController) DesagregadoPlaneacion() {
	defer ErrorController2(c)

	var d []models.ObjetoDesagregado

	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &d); err == nil {
		if dd, err2 := helpers.CalcularComponentesSalario(d); err2 == nil {
			c.Ctx.Output.SetStatus(201)
			c.Data["json"] = map[string]interface{}{"Success": true, "Status": 201, "Message": "Cálculos realizados con éxito", "Data": dd}
		} else {
			panic(err2)
		}
	} else {
		panic(map[string]interface{}{"funcion": "DesagregadoPlaneacion", "err": err.Error(), "status": "400"})
	}
	c.ServeJSON()
}

func ErrorController2(c *ServicesController) {
	if err := recover(); err != nil {
		logs.Error(err)
		localError := err.(map[string]interface{})
		c.Data["mesaage"] = (beego.AppConfig.String("appname") + "/" + "ServicesController" + "/" + (localError["funcion"]).(string))
		c.Data["data"] = (localError["err"])
		if status, ok := localError["status"]; ok {
			c.Abort(status.(string))
		} else {
			c.Abort("500")
		}
	}
}