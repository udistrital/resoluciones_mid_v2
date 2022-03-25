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
	c.Mapping("ModificarVinculacion", c.ModificarVinculacion)
	c.Mapping("DocentesPrevinculados", c.DocentesPrevinculados)
	c.Mapping("DocentesCargaHoraria", c.DocentesCargaHoraria)
	c.Mapping("InformeVinculaciones", c.InformeVinculaciones)
	c.Mapping("DesvincularDocentes", c.DesvincularDocentes)
	c.Mapping("CalcularValorContratosSeleccionados", c.CalcularValorContratosSeleccionados)
}

// Post ...
// @Title Create Vinculaciones
// @Description Registra las vinculaciones calculando los valores del contrato
// @Param	body		body 	models.ObjetoPrevinculaciones	true		"body for ObjetoPrevinculaciones content"
// @Success 201 {object} []models.ObjetoPrevinculaciones
// @Failure 400 bad request
// @Failure 500 Internal server error
// @router / [post]
func (c *GestionVinculacionesController) Post() {
	defer helpers.ErrorController(c.Controller, "GestionVinculacionesController")

	var p models.ObjetoPrevinculaciones

	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &p); err == nil {
		if r, err2 := helpers.RegistrarVinculaciones(p); err == nil {
			c.Ctx.Output.SetStatus(201)
			c.Data["json"] = map[string]interface{}{"Success": true, "Status": 201, "Message": "Vinculaciones registradas con exito", "Data": r}
		} else {
			panic(err2)
		}
	} else {
		panic(map[string]interface{}{"funcion": "Post", "err": err.Error(), "status": "400"})
	}
	c.ServeJSON()
}

// ModificarVinculacion ...
// @Title ModificarVinculaciones
// @Description Registra la modificacion de una vinculacion
// @Param	body		body 	models.ObjetoModificaciones	true		"body for ObjetoModificaciones content"
// @Success 201 {object} models.VinculacionDocente "Vinculaciones modificadas con éxito"
// @Failure 400 bad request
// @Failure 500 Internal server error
// @router /modificar_vinculacion [post]
func (c *GestionVinculacionesController) ModificarVinculacion() {
	defer helpers.ErrorController(c.Controller, "GestionVinculacionesController")

	var p models.ObjetoModificaciones

	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &p); err == nil {
		if r, err2 := helpers.ModificarVinculaciones(p); err == nil {
			c.Ctx.Output.SetStatus(201)
			c.Data["json"] = map[string]interface{}{"Success": true, "Status": 201, "Message": "Vinculaciones modificadas con éxito", "Data": r}
		} else {
			panic(err2)
		}
	} else {
		panic(map[string]interface{}{"funcion": "ModificarVinculacion", "err": err.Error(), "status": "400"})
	}
	c.ServeJSON()
}

// Desvincular ...
// @Title Desvincular
// @Description Registra la cancelación de una vinculacion
// @Param	body		body 	models.ObjetoCancelaciones	true		"body for ObjetoCancelaciones content"
// @Success 201 {object} []models.VinculacionDocente "Cancelaciones registradas con éxito"
// @Failure 400 bad request
// @Failure 500 Internal server error
// @router /desvincular [post]
func (c *GestionVinculacionesController) Desvincular() {
	defer helpers.ErrorController(c.Controller, "GestionVinculacionesController")

	var p models.ObjetoCancelaciones

	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &p); err == nil {
		if r, err2 := helpers.RegistrarCancelaciones(p); err == nil {
			c.Ctx.Output.SetStatus(201)
			c.Data["json"] = map[string]interface{}{"Success": true, "Status": 201, "Message": "Cancelaciones registradas con éxito", "Data": r}
		} else {
			panic(err2)
		}
	} else {
		panic(map[string]interface{}{"funcion": "Desvincular", "err": err.Error(), "status": "400"})
	}
	c.ServeJSON()
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
			c.Data["json"] = map[string]interface{}{"Success": true, "Status": 201, "Message": "Vinculaciones retiradas con exito", "Data": "OK"}
		} else {
			panic(err2)
		}
	} else {
		panic(map[string]interface{}{"funcion": "DesvincularDocentes", "err": err.Error(), "status": "400"})
	}
	c.ServeJSON()
}

// CalcularValorContratosSeleccionados ...
// @Title CalcularValorContratosSeleccionados
// @Description Calcula el valor total de los contratos seleccionados
// @Param	body		body 	models.ObjetoPrevinculaciones	true		"body for vinculaciones content"
// @Success 201 {string} Valor total de los contratos seleccionados
// @Failure 400 bad request
// @Failure 500 Internal server error
// @router /calcular_valor_contratos_seleccionados [post]
func (c *GestionVinculacionesController) CalcularValorContratosSeleccionados() {
	defer helpers.ErrorController(c.Controller, "GestionVinculacionesController")

	var p models.ObjetoPrevinculaciones
	var total int

	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &p); err == nil {
		if v, err1 := helpers.ConstruirVinculaciones(p); err1 == nil {
			if w, err2 := helpers.CalcularSalarioPrecontratacion(v); err2 == nil {
				total = int(helpers.CalcularTotalSalarios(w))
				c.Ctx.Output.SetStatus(201)
				c.Data["json"] = map[string]interface{}{"Success": true, "Status": 201, "Message": "", "Data": helpers.FormatMoney(total, 2)}
			} else {
				panic(err2)
			}
		} else {
			panic(err1)
		}
	} else {
		panic(map[string]interface{}{"funcion": "CalcularValorContratosSeleccionados", "err": err.Error(), "status": "400"})
	}
	c.ServeJSON()
}
