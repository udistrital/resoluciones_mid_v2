package controllers

import (
	"encoding/json"
	"fmt"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
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
// @Success 201 {int} models.ExpedicionResolucion
// @Failure 400 bad request
// @Failure 404 aborted by server
// @router /expedir [post]
func (c *ExpedirResolucionController) Expedir() {

	defer func() {
		if err := recover(); err != nil {
			logs.Error(err)
			localError := err.(map[string]interface{})
			c.Data["mesaage"] = (beego.AppConfig.String("appname") + "/" + "ExpedirResolucionController" + "/" + (localError["funcion"]).(string))
			c.Data["data"] = (localError["err"])
			if status, ok := localError["status"]; ok {
				c.Abort(status.(string))
			} else {
				c.Abort("404")
			}
		}
	}()

	var m models.ExpedicionResolucion

	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &m); err == nil {
		if err := helpers.Expedir2(m); err == nil {
			fmt.Println("Sisi")
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
// @Success 201 {string}
// @Failure 400 bad request
// @Failure 404 aborted by server
// @router /validar_datos_expedicion [post]
func (c *ExpedirResolucionController) ValidarDatosExpedicion() {

	defer func() {
		if err := recover(); err != nil {
			logs.Error(err)
			localError := err.(map[string]interface{})
			c.Data["mesaage"] = (beego.AppConfig.String("appname") + "/" + "ExpedirResolucionController" + "/" + (localError["funcion"]).(string))
			c.Data["data"] = (localError["err"])
			if status, ok := localError["status"]; ok {
				c.Abort(status.(string))
			} else {
				c.Abort("404")
			}
		}
	}()

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
// @Param	body		body 	[]models.ExpedicionResolucion	true		"body for Validar Datos Expedición content"
// @Success 201 {int} models.ExpedicionResolucion
// @Failure 400 bad request
// @Failure 404 aborted by server
// @router /expedirModificacion [post]
func (c *ExpedirResolucionController) ExpedirModificacion() {

	defer func() {
		if err := recover(); err != nil {
			logs.Error(err)
			localError := err.(map[string]interface{})
			c.Data["mesaage"] = (beego.AppConfig.String("appname") + "/" + "ExpedirResolucionController" + "/" + (localError["funcion"]).(string))
			c.Data["data"] = (localError["err"])
			if status, ok := localError["status"]; ok {
				c.Abort(status.(string))
			} else {
				c.Abort("404")
			}
		}
	}()

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
// @Param	body		body 	[]models.ExpedicionCancelacion	true		"body for Expedición a cancelar content"
// @Success 201 {int} models.ExpedicionCancelacion
// @Failure 400 bad request
// @Failure 404 aborted by server
// @router /cancelar [post]
func (c *ExpedirResolucionController) Cancelar() {

	defer func() {
		if err := recover(); err != nil {
			logs.Error(err)
			localError := err.(map[string]interface{})
			c.Data["mesaage"] = (beego.AppConfig.String("appname") + "/" + "ExpedirResolucionController" + "/" + (localError["funcion"]).(string))
			c.Data["data"] = (localError["err"])
			if status, ok := localError["status"]; ok {
				c.Abort(status.(string))
			} else {
				c.Abort("404")
			}
		}
	}()

	var m models.ExpedicionCancelacion

	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &m); err == nil {
		if err := helpers.Cancelar(m); err == nil {
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