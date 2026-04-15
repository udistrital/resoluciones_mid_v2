package controllers

import (
	"strings"
	"time"

	"github.com/astaxie/beego"
	"github.com/udistrital/resoluciones_mid_v2/helpers"
	"github.com/udistrital/resoluciones_mid_v2/models"
	"github.com/udistrital/resoluciones_mid_v2/services"
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
	c.Mapping("ConsultarSemaforoResolucion", c.ConsultarSemaforoResolucion)
	c.Mapping("ConsultarDashboardResoluciones", c.ConsultarDashboardResoluciones)
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
	decodeJSONBody(c.Ctx.Input.RequestBody, &p, "Post")

	if r, err2 := helpers.RegistrarVinculaciones(p); err2 == nil {
		writeJSON(&c.Controller, 201, "Vinculaciones registradas con exito", r, nil)
	} else {
		panic(err2)
	}
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

	var vd models.EdicionVinculaciones
	decodeJSONBody(c.Ctx.Input.RequestBody, &vd, "ModificarVinculacion")

	if r, err2 := helpers.EditarVinculaciones(vd); err2 == nil {
		writeJSON(&c.Controller, 201, "Vinculaciones modificadas con éxito", r, nil)
	} else {
		panic(err2)
	}
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
	decodeJSONBody(c.Ctx.Input.RequestBody, &p, "ModificarVinculacion")

	if r, err2 := helpers.ModificarVinculaciones(p); err2 == nil {
		writeJSON(&c.Controller, 201, "Vinculaciones modificadas con éxito", r, nil)
	} else {
		panic(err2)
	}
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
	decodeJSONBody(c.Ctx.Input.RequestBody, &p, "Desvincular")

	if r, err2 := helpers.RegistrarCancelaciones(p); err2 == nil {
		writeJSON(&c.Controller, 201, "Cancelaciones registradas con éxito", r, nil)
	} else {
		panic(err2)
	}
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
	parsePositivePathID(&c.Controller, ":resolucion_id", "DocentesPrevinculados")

	if vinculaciones, err2 := helpers.ListarVinculaciones(resolucionId, false); err2 == nil {
		writeJSON(&c.Controller, 200, "Successful", vinculaciones, nil)
	} else {
		panic(err2)
	}
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
	parsePositivePathID(&c.Controller, ":resolucion_id", "DocentesPrevinculadosRp")

	if vinculaciones, err2 := helpers.ListarVinculaciones(resolucionId, true); err2 == nil {
		writeJSON(&c.Controller, 200, "Successful", vinculaciones, nil)
	} else {
		panic(err2)
	}
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

	vigencia := parseYearPathParam(&c.Controller, ":vigencia", "DocentesCargaHoraria")
	periodo := parseSingleDigitPositivePathParam(&c.Controller, ":periodo", "DocentesCargaHoraria")
	dedicacion := c.Ctx.Input.Param(":dedicacion")
	facultad := c.Ctx.Input.Param(":facultad")
	nivelAcademico := c.Ctx.Input.Param(":nivel_academico")
	parsePositiveIntParam(facultad, "DocentesCargaHoraria")

	if respuesta, err := helpers.ListarDocentesCargaHoraria(vigencia, periodo, dedicacion, facultad, nivelAcademico); err == nil {
		writeJSON(&c.Controller, 200, "Successful", respuesta.CargasLectivas.CargaLectiva, nil)
	} else {
		panic(err)
	}
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

	var v []models.Vinculaciones
	decodeJSONBody(c.Ctx.Input.RequestBody, &v, "InformeVinculaciones")

	if i, err2 := helpers.GenerarInformeVinculaciones(v); err2 == nil {
		writeJSON(&c.Controller, 200, "Informe generado con exito", i, nil)
	} else {
		panic(err2)
	}
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

	var v []models.Vinculaciones
	decodeJSONBody(c.Ctx.Input.RequestBody, &v, "DesvincularDocentes")

	if err2 := helpers.RetirarVinculaciones(v); err2 == nil {
		writeJSON(&c.Controller, 201, "Vinculaciones retiradas con exito", "OK", nil)
	} else {
		panic(err2)
	}
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

	var p models.ObjetoPrevinculaciones
	var total int
	decodeJSONBody(c.Ctx.Input.RequestBody, &p, "CalcularValorContratosSeleccionados")

	if v, err1 := helpers.ConstruirVinculaciones(p); err1 == nil {
		if w, err2 := helpers.CalcularSalarioPrecontratacion(v); err2 == nil {
			total = int(helpers.CalcularTotalSalarios(w))
			writeJSON(&c.Controller, 201, "Cálculo exitoso", helpers.FormatMoney(total, 2), nil)
		} else {
			panic(err2)
		}
	} else {
		panic(err1)
	}
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

	vigencia := parseYearPathParam(&c.Controller, ":vigencia", "ConsultarSemaforoDocente")
	periodo := parseSingleDigitPositivePathParam(&c.Controller, ":periodo", "ConsultarSemaforoDocente")
	docente := c.Ctx.Input.Param(":docente")
	parsePositiveIntParam(docente, "ConsultarSemaforoDocente")

	if respuesta, err := helpers.BuscarCategoriaDocente(vigencia, periodo, docente); err == nil {
		writeJSON(&c.Controller, 200, "Successful", respuesta.CategoriaDocente.Categoria, nil)
	} else {
		panic(err)
	}
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
	vigencia := parseYearPathParam(&c.Controller, ":vigencia", "ConsultarSemanasRestantes")
	contrato := c.Ctx.Input.Param(":contrato")

	fechaParsed, err := time.Parse("2006-01-02", fecha)
	vigenciaParsed := parsePositiveIntParam(vigencia, "ConsultarSemanasRestantes")
	contratoValido := strings.Contains(contrato, "DVE")

	if err != nil || !contratoValido {
		panic(map[string]interface{}{"funcion": "ConsultarSemanasRestantes", "err": helpers.ErrorParametros, "status": "400"})
	}

	if respuesta, err := helpers.CalcularNumeroSemanas(fechaParsed, contrato, vigenciaParsed); err == nil {
		writeJSON(&c.Controller, 200, "Successful", respuesta, nil)
	} else {
		panic(err.Error())
	}
}

// RegistrarRps ...
// @Title RegistrarRps
// @Description registra los numeros de RP en las respectivas vinculaciones y dispara el proceso de preliquidación
// @Param	body	body 	[]models.RpSeleccionado	true		"body for vinculaciones content"
// @Success 200 {object} string Proceso iniciado
// @Failure 400 bad request
// @Failure 500 Internal server error
// @router /rp_vinculaciones [post]
func (c *GestionVinculacionesController) RegistrarRps() {
	defer helpers.ErrorController(c.Controller, "GestionVinculacionesController")

	var rps []models.RpSeleccionado
	decodeJSONBody(c.Ctx.Input.RequestBody, &rps, "RegistrarRps")

	jobID := helpers.CrearJob(len(rps))
	go helpers.ProcesarPreliquidaciones(jobID, rps)

	writeJSON(&c.Controller, 200, "Proceso de registro de RPs iniciado correctamente", nil, map[string]interface{}{
		"JobId": jobID,
		"Total": len(rps),
	})
}

// ConsultarSemaforoResolucion ...
// @Title ConsultarSemaforoResolucion
// @Description Consulta el estado del dashboard de vinculaciones de una resolución y opcionalmente filtra por documento
// @Param resolucion_id path string true "ID de la resolución"
// @Param usuario query string true "Número de documento del usuario que consulta"
// @Param roles query string true "Roles del usuario separados por coma"
// @Param numero_documento query string false "Número de documento del docente para filtrar"
// @Success 200 {object} models.RespuestaSemaforoResolucion
// @Failure 400 bad request
// @Failure 403 forbidden
// @Failure 404 not found
// @Failure 500 internal server error
// @router /consultar_semaforo_resolucion/:resolucion_id [get]
func (c *GestionVinculacionesController) ConsultarSemaforoResolucion() {
	defer helpers.ErrorController(c.Controller, "GestionVinculacionesController")

	authContext := requireAuthenticatedContext(buildAuthenticatedContext(&c.Controller), "ConsultarSemaforoResolucion")
	numeroDocumentoStr := strings.TrimSpace(c.GetString("numero_documento"))

	resolucionId := parsePositivePathID(&c.Controller, ":resolucion_id", "ConsultarSemaforoResolucion")

	numeroDocumentoFiltro := parseOptionalPositiveIntPointer(numeroDocumentoStr, "numero_documento", "ConsultarSemaforoResolucion")

	if respuesta, errMap := services.ConsultarSemaforoResolucion(resolucionId, authContext.NumeroDocumento, authContext.Roles, numeroDocumentoFiltro); errMap == nil {
		writeJSON(&c.Controller, 200, "Semáforo consultado exitosamente", respuesta, nil)
	} else {
		panic(errMap)
	}
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
		panic(badRequest("ObtenerProgreso", validateRequiredText(jobID, "Job ID requerido")))
	}

	result := helpers.ObtenerJob(jobID)
	if result["Success"].(bool) {
		c.Ctx.Output.SetStatus(200)
		c.Data["json"] = result
		c.ServeJSON()
	} else {
		writeErrorJSON(&c.Controller, 404, result["Message"].(string), nil)
	}
}

// ConsultarDashboardResoluciones ...
// @Title ConsultarDashboardResoluciones
// @Description Consulta el dashboard general de resoluciones con su avance de vinculaciones
// @Param numero_documento query string true "Número de documento del usuario que consulta"
// @Param roles query string true "Roles del usuario separados por coma"
// @Param vigencia query string true "Vigencia a consultar"
// @Param id_oikos query string false "Dependencia Oikos para filtrar"
// @Param limit query int false "Cantidad de registros por página"
// @Param offset query int false "Posición inicial para paginación"
// @Success 200 {object} models.RespuestaDashboardResoluciones
// @Failure 400 bad request
// @Failure 403 forbidden
// @Failure 500 internal server error
// @router /dashboard_resoluciones [get]
func (c *GestionVinculacionesController) ConsultarDashboardResoluciones() {
	defer helpers.ErrorController(c.Controller, "GestionVinculacionesController")

	authContext := requireAuthenticatedContext(buildAuthenticatedContext(&c.Controller), "ConsultarDashboardResoluciones")
	vigenciaStr := strings.TrimSpace(c.GetString("vigencia"))
	idOikosStr := strings.TrimSpace(c.GetString("id_oikos"))
	limitStr := strings.TrimSpace(c.GetString("limit"))
	offsetStr := strings.TrimSpace(c.GetString("offset"))

	if err := validateRequiredText(vigenciaStr, helpers.ErrorParametros); err != nil {
		panic(badRequest("ConsultarDashboardResoluciones", err))
	}

	vigencia := parseRequiredPositiveInt(vigenciaStr, "vigencia", "ConsultarDashboardResoluciones")

	limit := 10
	if limitStr != "" {
		limit = parseRequiredPositiveInt(limitStr, "limit", "ConsultarDashboardResoluciones")
	}

	offset := parseOptionalNonNegativeInt(offsetStr, "offset", "ConsultarDashboardResoluciones", 0)

	dependenciaFiltro := parseOptionalPositiveIntPointer(idOikosStr, "id_oikos", "ConsultarDashboardResoluciones")

	if respuesta, errMap := services.ConsultarDashboardResoluciones(authContext.NumeroDocumento, authContext.Roles, vigencia, dependenciaFiltro, limit, offset); errMap == nil {
		writeJSON(&c.Controller, 200, "Dashboard de resoluciones consultado exitosamente", respuesta, nil)
	} else {
		panic(errMap)
	}
}
