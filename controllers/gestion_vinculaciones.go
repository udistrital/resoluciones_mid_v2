package controllers

import (
	"encoding/json"
	"strconv"
	"strings"
	"time"

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
	c.Mapping("EditarVinculaciones", c.EditarVinculaciones)
	c.Mapping("ModificarVinculacion", c.ModificarVinculacion)
	c.Mapping("DocentesPrevinculados", c.DocentesPrevinculados)
	c.Mapping("DocentesPrevinculadosRp", c.DocentesPrevinculadosRp)
	c.Mapping("DocentesCargaHoraria", c.DocentesCargaHoraria)
	c.Mapping("InformeVinculaciones", c.InformeVinculaciones)
	c.Mapping("DesvincularDocentes", c.DesvincularDocentes)
	c.Mapping("CalcularValorContratosSeleccionados", c.CalcularValorContratosSeleccionados)
	c.Mapping("ConsultarSemaforoDocente", c.ConsultarSemaforoDocente)
	c.Mapping("ConsultarSemanasRestantes", c.ConsultarSemanasRestantes)
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

	if v, e := helpers.ValidarBody(c.Ctx.Input.RequestBody); !v || e != nil {
		panic(map[string]interface{}{"funcion": "Post", "err": helpers.ErrorBody, "status": "400"})
	}

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

// EditarVinculacion ...
// @Title EditarVinculacion
// @Description Modifica las vinculaciones docente
// @Param	body		body 	models.EdicionVinculaciones	true		"body for ObjetoModificaciones content"
// @Success 201 {object} []models.EdicionVinculaciones "Vinculaciones modificadas con éxito"
// @Failure 400 bad request
// @Failure 500 Internal server error
// @router /editar_vinculaciones [post]
func (c *GestionVinculacionesController) EditarVinculaciones() {
	defer helpers.ErrorController(c.Controller, "GestionVinculacionesController")

	if v, e := helpers.ValidarBody(c.Ctx.Input.RequestBody); !v || e != nil {
		panic(map[string]interface{}{"funcion": "ModificarVinculacion", "err": helpers.ErrorBody, "status": "400"})
	}

	var vd models.EdicionVinculaciones

	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &vd); err == nil {
		if r, err2 := helpers.EditarVinculaciones(vd); err == nil {
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

	if v, e := helpers.ValidarBody(c.Ctx.Input.RequestBody); !v || e != nil {
		panic(map[string]interface{}{"funcion": "ModificarVinculacion", "err": helpers.ErrorBody, "status": "400"})
	}

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

	if v, e := helpers.ValidarBody(c.Ctx.Input.RequestBody); !v || e != nil {
		panic(map[string]interface{}{"funcion": "Desvincular", "err": helpers.ErrorBody, "status": "400"})
	}

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
	id, err := strconv.Atoi(resolucionId)

	if err != nil || id <= 0 {
		panic(map[string]interface{}{"funcion": "DocentesPrevinculados", "err": helpers.ErrorParametros, "status": "400"})
	}

	if vinculaciones, err2 := helpers.ListarVinculaciones(resolucionId, false); err2 == nil {
		c.Ctx.Output.SetStatus(200)
		c.Data["json"] = map[string]interface{}{"Success": true, "Status": "200", "Message": "Successful", "Data": vinculaciones}
	} else {
		panic(err2)
	}
	c.ServeJSON()
}

// DocentesPrevinculadosRp ...
// @Title DocentesPrevinculadosRp
// @Description Docentes previnculados a una resolución con rps
// @Param	resolucion_id		path 	string	true		"El id de la resolución"
// @Success 200 {object} []models.Vinculaciones
// @Failure 400 bad request
// @Failure 500 Internal server error
// @router /rp/:resolucion_id [get]
func (c *GestionVinculacionesController) DocentesPrevinculadosRp() {
	defer helpers.ErrorController(c.Controller, "GestionVinculacionesController")

	resolucionId := c.Ctx.Input.Param(":resolucion_id")
	id, err := strconv.Atoi(resolucionId)

	if err != nil || id <= 0 {
		panic(map[string]interface{}{"funcion": "DocentesPrevinculadosRp", "err": helpers.ErrorParametros, "status": "400"})
	}

	if vinculaciones, err2 := helpers.ListarVinculaciones(resolucionId, true); err2 == nil {
		c.Ctx.Output.SetStatus(200)
		c.Data["json"] = map[string]interface{}{"Success": true, "Status": "200", "Message": "Successful", "Data": vinculaciones}
	} else {
		panic(err2)
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
// @Success 200 {object} []models.CargaLectiva Carga horaria de los docentes organizada
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
		panic(map[string]interface{}{"funcion": "DocentesCargaHoraria", "err": helpers.ErrorParametros, "status": "400"})
	}

	if respuesta, err := helpers.ListarDocentesCargaHoraria(vigencia, periodo, dedicacion, facultad, nivelAcademico); err == nil {
		c.Ctx.Output.SetStatus(200)
		c.Data["json"] = map[string]interface{}{"Success": true, "Status": "200", "Message": "Successful", "Data": respuesta.CargasLectivas.CargaLectiva}
	} else {
		panic(err)
	}
	c.ServeJSON()
}

// InformeVinculaciones ...
// @Title InformeVinculaciones
// @Description Genera un informe de las vinculaciones
// @Param	body		body 	[]models.Vinculaciones	true		"body for vinculaciones content"
// @Success 200 {object} string Base64 encoded file
// @Failure 400 bad request
// @Failure 500 Internal server error
// @router /informe_vinculaciones [post]
func (c *GestionVinculacionesController) InformeVinculaciones() {
	defer helpers.ErrorController(c.Controller, "GestionVinculacionesController")

	if v, e := helpers.ValidarBody(c.Ctx.Input.RequestBody); !v || e != nil {
		panic(map[string]interface{}{"funcion": "InformeVinculaciones", "err": helpers.ErrorBody, "status": "400"})
	}

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
// @Success 201 {object} string OK
// @Failure 400 bad request
// @Failure 500 Internal server error
// @router /desvincular_docentes [post]
func (c *GestionVinculacionesController) DesvincularDocentes() {
	defer helpers.ErrorController(c.Controller, "GestionVinculacionesController")

	if v, e := helpers.ValidarBody(c.Ctx.Input.RequestBody); !v || e != nil {
		panic(map[string]interface{}{"funcion": "DesvincularDocentes", "err": helpers.ErrorBody, "status": "400"})
	}

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
// @Success 201 {object} string Valor total de los contratos seleccionados
// @Failure 400 bad request
// @Failure 500 Internal server error
// @router /calcular_valor_contratos_seleccionados [post]
func (c *GestionVinculacionesController) CalcularValorContratosSeleccionados() {
	defer helpers.ErrorController(c.Controller, "GestionVinculacionesController")

	if v, e := helpers.ValidarBody(c.Ctx.Input.RequestBody); !v || e != nil {
		panic(map[string]interface{}{"funcion": "CalcularValorContratosSeleccionados", "err": helpers.ErrorBody, "status": "400"})
	}

	var p models.ObjetoPrevinculaciones
	var total int

	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &p); err == nil {
		if v, err1 := helpers.ConstruirVinculaciones(p); err1 == nil {
			if w, err2 := helpers.CalcularSalarioPrecontratacion(v); err2 == nil {
				total = int(helpers.CalcularTotalSalarios(w))
				c.Ctx.Output.SetStatus(201)
				c.Data["json"] = map[string]interface{}{"Success": true, "Status": 201, "Message": "Cálculo exitoso", "Data": helpers.FormatMoney(total, 2)}
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

// ConsultarSemaforoDocente ...
// @Title ConsultarSemaforoDocente
// @Description Consulta el estado del semaforo del docente en condor
// @Param vigencia query string false "año a consultar"
// @Param periodo query string false "periodo a listar"
// @Param docente query string false "documento del docente a consultar"
// @Success 200 {object} string categoria del docente
// @Failure 400 bad request
// @Failure 500 Internal server error
// @router /consultar_semaforo_docente/:vigencia/:periodo/:docente [get]
func (c *GestionVinculacionesController) ConsultarSemaforoDocente() {
	defer helpers.ErrorController(c.Controller, "GestionVinculacionesController")

	vigencia := c.Ctx.Input.Param(":vigencia")
	periodo := c.Ctx.Input.Param(":periodo")
	docente := c.Ctx.Input.Param(":docente")

	vig, err1 := strconv.Atoi(vigencia)
	per, err2 := strconv.Atoi(periodo)
	doc, err3 := strconv.Atoi(docente)

	if (err1 != nil) || (err2 != nil) || (err3 != nil) || (vig == 0) || (per == 0) || (doc == 0) || (len(vigencia) != 4) || (len(periodo) != 1) {
		panic(map[string]interface{}{"funcion": "ConsultarSemaforoDocente", "err": helpers.ErrorParametros, "status": "400"})
	}

	if respuesta, err := helpers.BuscarCategoriaDocente(vigencia, periodo, docente); err == nil {
		c.Ctx.Output.SetStatus(200)
		c.Data["json"] = map[string]interface{}{"Success": true, "Status": "200", "Message": "Successful", "Data": respuesta.CategoriaDocente.Categoria}
	} else {
		panic(err)
	}
	c.ServeJSON()
}

// ConsultarSemanasRestantes ...
// @Title ConsultarSemanasRestantes
// @Description Consulta el numero de semanas restantes de un contrato específico
// @Param fecha query string true "Documento del docente a consultar"
// @Param vigencia query string true "Año de la vinculación"
// @Param contrato query string true "Número de contrato de la vinculación"
// @Success 200 {object} int Numero de semanas resultantes
// @Failure 400 bad request
// @Failure 500 Internal server error
// @router /consultar_semanas_restantes/:fecha/:vigencia/:contrato [get]
func (c *GestionVinculacionesController) ConsultarSemanasRestantes() {
	defer helpers.ErrorController(c.Controller, "GestionVinculacionesController")

	fecha := c.Ctx.Input.Param(":fecha")
	vigencia := c.Ctx.Input.Param(":vigencia")
	contrato := c.Ctx.Input.Param(":contrato")

	fechaParsed, err := time.Parse("2006-01-02", fecha)
	vigenciaParsed, err2 := strconv.Atoi(vigencia)
	contratoValido := strings.Contains(contrato, "DVE")

	if (err != nil) || (err2 != nil) || (vigenciaParsed == 0) || (len(vigencia) != 4) || !contratoValido {
		panic(map[string]interface{}{"funcion": "ConsultarSemanasRestantes", "err": helpers.ErrorParametros, "status": "400"})
	}

	if respuesta, err := helpers.CalcularNumeroSemanas(fechaParsed, contrato, vigenciaParsed); err == nil {
		c.Ctx.Output.SetStatus(200)
		c.Data["json"] = map[string]interface{}{"Success": true, "Status": "200", "Message": "Successful", "Data": respuesta}
	} else {
		panic(err.Error())
	}

	c.ServeJSON()
}

// RegistrarRps ...
// @Title RegistrarRps
// @Description registra los numeros de RP en las respectivas vinculaciones y dispara el proceso de preliquidación
// @Param	body	body 	[]models.RpSeleccionado	true		"body for vinculaciones content"
// @Success 202 {object} string Proceso iniciado
// @Failure 400 bad request
// @Failure 500 Internal server error
// @router /rp_vinculaciones [post]
func (c *GestionVinculacionesController) RegistrarRps() {
	defer helpers.ErrorController(c.Controller, "GestionVinculacionesController")

	var rps []models.RpSeleccionado
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &rps); err != nil {
		panic(map[string]interface{}{"funcion": "RegistrarRps", "err": err.Error(), "status": "400"})
	}

	jobID := helpers.CrearJob(len(rps))
	go helpers.ProcesarPreliquidaciones(jobID, rps)

	c.Ctx.Output.SetStatus(202)
	c.Data["json"] = map[string]interface{}{
		"Success": true,
		"Status":  202,
		"Message": "Proceso de registro de RPs iniciado correctamente",
		"JobId":   jobID,
		"Total":   len(rps),
	}
	c.ServeJSON()
}

// ObtenerProgreso ...
// @Title ObtenerProgreso
// @Description Consulta el estado y progreso de un job en ejecución
// @Param	jobId	path	string	true	"ID del job"
// @Success 200 {object} map[string]interface{}
// @Failure 404 Job no encontrado
// @router /progreso/:jobId [get]
func (c *GestionVinculacionesController) ObtenerProgreso() {
	defer helpers.ErrorController(c.Controller, "GestionVinculacionesController")

	jobID := c.Ctx.Input.Param(":jobId")
	if jobID == "" {
		panic(map[string]interface{}{"funcion": "ObtenerProgreso", "err": "Job ID requerido", "status": "400"})
	}

	result := helpers.ObtenerJob(jobID)
	if result["Success"].(bool) {
		c.Data["json"] = result
		c.ServeJSON()
	} else {
		c.Ctx.Output.SetStatus(404)
		c.Data["json"] = result
		c.ServeJSON()
	}
}
