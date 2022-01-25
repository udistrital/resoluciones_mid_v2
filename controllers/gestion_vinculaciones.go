package controllers

import (
	"strconv"

	"github.com/astaxie/beego"
	"github.com/udistrital/resoluciones_mid_v2/helpers"
)

// Gestion_vinculacionesController operations for Gestion_vinculaciones
type GestionVinculacionesController struct {
	beego.Controller
}

// URLMapping ...
func (c *GestionVinculacionesController) URLMapping() {
	c.Mapping("Post", c.Post)
	c.Mapping("GetOne", c.DocentesPrevinculados)
	c.Mapping("GetAll", c.DocentesCargaHoraria)
	c.Mapping("Put", c.Put)
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
// @Success 200 {object} []models.VinculacionDocente
// @Failure 400 bad request
// @Failure 500 Internal server error
// @router /docentes_previnculados/:resolucion_id [get]
func (c *GestionVinculacionesController) DocentesPrevinculados() {
	defer helpers.ErrorController(c.Controller, "GestionVinculacionesController")

}

// DocentesCargaHoraria ...
// @Title DocentesCargaHoraria
// @Description Obtiene la carga horaria de los docentes desde condor.
// @Param vigencia query string false "año a consultar"
// @Param periodo query string false "periodo a listar"
// @Param dedicacion query string false "dedicacion del docente"
// @Param facultad query string false "facultad"
// @Param nivel_academico query string false "nivel_academico"
// @Success 200 {object} []models.DocentesCargaHoraria.CargasLectivas.CargaLectiva
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

// Put ...
// @Title Put
// @Description update the Gestion_vinculaciones
// @Param	id		path 	string	true		"The id you want to update"
// @Param	body		body 	models.Gestion_vinculaciones	true		"body for Gestion_vinculaciones content"
// @Success 200 {object} models.Gestion_vinculaciones
// @Failure 403 :id is not int
// @router /:id [put]
func (c *GestionVinculacionesController) Put() {
	defer helpers.ErrorController(c.Controller, "GestionVinculacionesController")

}

// Delete ...
// @Title Delete
// @Description delete the Gestion_vinculaciones
// @Param	id		path 	string	true		"The id you want to delete"
// @Success 200 {string} delete success!
// @Failure 403 id is empty
// @router /:id [delete]
func (c *GestionVinculacionesController) Delete() {
	defer helpers.ErrorController(c.Controller, "GestionVinculacionesController")

}
