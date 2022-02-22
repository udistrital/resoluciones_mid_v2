package controllers

import (
	"encoding/json"
	"strconv"

	"github.com/astaxie/beego"
	"github.com/udistrital/resoluciones_mid_v2/helpers"
	"github.com/udistrital/resoluciones_mid_v2/models"
)

// Gestion_vinculacionesController operations for Gestion_vinculaciones
type GestionVinculacionesController struct {
	beego.Controller
}

// URLMapping ...
func (c *GestionVinculacionesController) URLMapping() {
	c.Mapping("Post", c.Post)
	c.Mapping("DocentesPrevinculados", c.DocentesPrevinculados)
	c.Mapping("DocentesCargaHoraria", c.DocentesCargaHoraria)
	c.Mapping("InformeVinculaciones", c.InformeVinculaciones)
	c.Mapping("Delete", c.Delete)
}

// Post ...
// @Title Create
// @Description create Gestion_vinculaciones
// @Param	body		body 	models.Gestion_vinculaciones	true		"body for Gestion_vinculaciones content"
// @Success 201 {object} models.Gestion_vinculaciones
// @Failure 403 body is empty
// @router / [post]
func (c *GestionVinculacionesController) Post() {
	defer helpers.ErrorController(c.Controller, "GestionVinculacionesController")

}

// DocentesPrevinculados ...
// @Title DocentesPrevinculados
// @Description Docentes previnculados a una resolución
// @Param	resolucion_id		path 	string	true		"El id de la resolución"
// @Success 200 {object} []models.Vinculaciones
// @Failure 400 bad request
// @Failure 500 Internal server error
// @router /:resolucion_id [get]
func (c *GestionVinculacionesController) DocentesPrevinculados() {
	defer helpers.ErrorController(c.Controller, "GestionVinculacionesController")

	resolucionId := c.Ctx.Input.Param(":resolucion_id")
	_, err := strconv.Atoi(resolucionId)

	if err != nil {
		panic(map[string]interface{}{"funcion": "DocentesPrevinculados", "err": "Error en los parametros de ingreso", "status": "400"})
	}

	if vinculaciones, err2 := helpers.ListarVinculaciones(resolucionId); err2 == nil {
		c.Ctx.Output.SetStatus(200)
		c.Data["json"] = map[string]interface{}{"Success": true, "Status": "200", "Message": "Successful", "Data": vinculaciones}
	}
	c.ServeJSON()
}

// DocentesCargaHoraria ...
// @Title DocentesCargaHoraria
// @Description Obtiene la carga horaria de los docentes desde condor.
// @Param vigencia query string false "año a consultar"
// @Param periodo query string false "periodo a listar"
// @Param dedicacion query string false "dedicacion del docente"
// @Param facultad query string false "facultad"
// @Param nivel_academico query string false "nivel_academico"
// @Success 200 {object} []models.CargaLectiva
// @Failure 400 bad request
// @Failure 500 Internal server error
// @router /docentes_carga_horaria/:vigencia/:periodo/:dedicacion/:facultad/:nivel_academico [get]
func (c *GestionVinculacionesController) DocentesCargaHoraria() {
	defer helpers.ErrorController(c.Controller, "GestionVinculacionesController")

	vigencia := c.Ctx.Input.Param(":vigencia")
	periodo := c.Ctx.Input.Param(":periodo")
	dedicacion := c.Ctx.Input.Param(":dedicacion")
	facultad := c.Ctx.Input.Param(":facultad")
	nivelAcademico := c.Ctx.Input.Param(":nivel_academico")

	vig, err1 := strconv.Atoi(vigencia)
	per, err2 := strconv.Atoi(periodo)
	_, err4 := strconv.Atoi(facultad)

	if (err1 != nil) || (err2 != nil) || (err4 != nil) || (vig == 0) || (per == 0) || (len(vigencia) != 4) {
		panic(map[string]interface{}{"funcion": "DocentesCargaHoraria", "err": "Error en los parametros de ingreso", "status": "400"})
	}

	if respuesta, err := helpers.ListarDocentesCargaHoraria(vigencia, periodo, dedicacion, facultad, nivelAcademico); err == nil {
		c.Ctx.Output.SetStatus(200)
		c.Data["json"] = map[string]interface{}{"Success": true, "Status": "200", "Message": "Successful", "Data": respuesta.CargasLectivas.CargaLectiva}
	}
	c.ServeJSON()
}

// InformeVinculaciones ...
// @Title InformeVinculaciones
// @Description Genera un informe de las vinculaciones
// @Param	body		body 	[]models.Vinculaciones	true		"body for vinculaciones content"
// @Success 200 {string} string Base64 encoded file
// @Failure 400 bad request
// @Failure 500 Internal server error
// @router /informe_vinculaciones [post]
func (c *GestionVinculacionesController) InformeVinculaciones() {
	defer helpers.ErrorController(c.Controller, "GestionVinculacionesController")

	var v []models.Vinculaciones

	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err == nil {
		if i, err2 := helpers.GenerarInformeVinculaciones(v); err2 == nil {
			c.Ctx.Output.SetStatus(200)
			c.Data["json"] = map[string]interface{}{"Success": true, "Status": 200, "Message": "Informe generado con exito", "Data": i}
		} else {
			panic(err2)
		}
	} else {
		panic(map[string]interface{}{"funcion": "DesvincularDocentes", "err": err.Error(), "status": "400"})
	}
	c.ServeJSON()

}

// DesvincularDocentes ...
// @Title DesvincularDocentes
// @Description Elimina las vinculaciones
// @Param	body		body 	[]models.Vinculaciones	true		"body for vinculaciones content"
// @Success 201 {string} OK
// @Failure 400 bad request
// @Failure 500 Internal server error
// @router /desvincular_docentes [post]
func (c *GestionVinculacionesController) DesvincularDocentes() {
	defer helpers.ErrorController(c.Controller, "GestionVinculacionesController")

	var v []models.Vinculaciones

	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err == nil {
		if err2 := helpers.RetirarVinculaciones(v); err2 == nil {
			c.Ctx.Output.SetStatus(201)
			c.Data["json"] = map[string]interface{}{"Success": true, "Status": 200, "Message": "Vinculaciones retiradas con exito", "Data": "OK"}
		} else {
			panic(err2)
		}
	} else {
		panic(map[string]interface{}{"funcion": "DesvincularDocentes", "err": err.Error(), "status": "400"})
	}
	c.ServeJSON()
}
