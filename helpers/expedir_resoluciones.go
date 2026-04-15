package helpers

import (
	"fmt"
	"math"
	"strconv"
	"time"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/udistrital/resoluciones_mid_v2/models"
)

type contextoExpedicionResolucion struct {
	resolucion     models.Resolucion
	tipoResolucion models.Parametro
	facultad       models.Dependencia
	ordenadorGasto models.OrdenadorGasto
}

func resolverUsuarioExpedicion(usuarioCodificado string) (map[string]interface{}, error) {
	usuario, err := GetUsuario(usuarioCodificado)
	if err != nil {
		logs.Error(err)
		return nil, err
	}

	return usuario, nil
}

func cargarResolucionYTipo(idResolucion int) (models.Resolucion, models.Parametro, error) {
	var resolucion models.Resolucion
	var parametro models.Parametro

	if err := GetRequestNew("UrlCrudResoluciones", ResolucionEndpoint+strconv.Itoa(idResolucion), &resolucion); err != nil {
		logs.Error(err)
		return resolucion, parametro, err
	}

	url := "parametro/" + strconv.Itoa(resolucion.TipoResolucionId)
	if err := GetRequestNew("UrlcrudParametros", url, &parametro); err != nil {
		logs.Error(err)
		return resolucion, parametro, fmt.Errorf("Cargando tipo_resolucion -> %s", err.Error())
	}

	return resolucion, parametro, nil
}

func cargarOrdenadorGastoPorDependencia(dependenciaId int) (models.OrdenadorGasto, error) {
	var ordenadoresGasto []models.OrdenadorGasto
	var ordenadorGasto models.OrdenadorGasto

	url := "ordenador_gasto?query=DependenciaId:" + strconv.Itoa(dependenciaId)
	if err := GetRequestLegacy("UrlcrudCore", url, &ordenadoresGasto); err != nil {
		logs.Error(err)
		return ordenadorGasto, err
	}

	if len(ordenadoresGasto) > 0 {
		return ordenadoresGasto[0], nil
	}

	if err := GetRequestLegacy("UrlcrudCore", "ordenador_gasto/1", &ordenadorGasto); err != nil {
		logs.Error(err)
		return ordenadorGasto, err
	}

	return ordenadorGasto, nil
}

func cargarFacultadResolucion(idResolucion int) (models.ResolucionVinculacionDocente, models.Dependencia, error) {
	var resolucionVinculacion models.ResolucionVinculacionDocente
	var dependencia models.Dependencia

	if err := GetRequestNew("UrlCrudResoluciones", ResVinEndpoint+strconv.Itoa(idResolucion), &resolucionVinculacion); err != nil {
		logs.Error(err)
		return resolucionVinculacion, dependencia, err
	}

	if err := GetRequestLegacy("UrlcrudOikos", "dependencia/"+strconv.Itoa(resolucionVinculacion.FacultadId), &dependencia); err != nil {
		logs.Error(err)
		return resolucionVinculacion, dependencia, err
	}

	return resolucionVinculacion, dependencia, nil
}

func cargarContextoExpedicionResolucion(idResolucion int) (contextoExpedicionResolucion, error) {
	var contexto contextoExpedicionResolucion
	var err error

	contexto.resolucion, contexto.tipoResolucion, err = cargarResolucionYTipo(idResolucion)
	if err != nil {
		return contexto, err
	}

	contexto.ordenadorGasto, err = cargarOrdenadorGastoPorDependencia(contexto.resolucion.DependenciaId)
	if err != nil {
		return contexto, err
	}

	if _, contexto.facultad, err = cargarFacultadResolucion(idResolucion); err != nil {
		return contexto, err
	}

	return contexto, nil
}

func cargarVinculacionPorID(idVinculacion int) (models.VinculacionDocente, error) {
	var vinculacion models.VinculacionDocente
	url := VinculacionEndpoint + strconv.Itoa(idVinculacion)
	if err := GetRequestNew("UrlCrudResoluciones", url, &vinculacion); err != nil {
		logs.Error(err)
		return vinculacion, err
	}

	return vinculacion, nil
}

func cargarVinculacionPorQuery(idVinculacion int) (models.VinculacionDocente, error) {
	var vinculaciones []models.VinculacionDocente
	url := "vinculacion_docente?query=Id:" + strconv.Itoa(idVinculacion)
	if err := GetRequestNew("UrlCrudResoluciones", url, &vinculaciones); err != nil {
		logs.Error(err)
		return models.VinculacionDocente{}, err
	}
	if len(vinculaciones) == 0 {
		return models.VinculacionDocente{}, fmt.Errorf("no se encontró la vinculación %d", idVinculacion)
	}

	return vinculaciones[0], nil
}

func cargarTipoContratoAgora(idTipoContrato int) (models.TipoContrato, error) {
	var tipoContrato models.TipoContrato
	url := "tipo_contrato/" + strconv.Itoa(idTipoContrato)
	if err := GetRequestLegacy("UrlcrudAgora", url, &tipoContrato); err != nil {
		logs.Error(err)
		return tipoContrato, err
	}

	return tipoContrato, nil
}

func cargarProveedorAgora(numeroDocumento int) (models.InformacionProveedor, error) {
	var proveedores []models.InformacionProveedor
	url := "informacion_proveedor?query=NumDocumento:" + strconv.Itoa(numeroDocumento)
	if err := GetRequestLegacy("UrlcrudAgora", url, &proveedores); err != nil {
		logs.Error(err)
		return models.InformacionProveedor{}, err
	}
	if len(proveedores) == 0 {
		return models.InformacionProveedor{}, fmt.Errorf("no se encontró proveedor en agora para documento %d", numeroDocumento)
	}

	return proveedores[0], nil
}

func cargarDependenciaProyectoCurricular(idProyectoCurricular int) (models.Dependencia, error) {
	var dependencias []models.Dependencia
	url := "dependencia?query=Id:" + strconv.Itoa(idProyectoCurricular)
	if err := GetRequestLegacy("UrlcrudOikos", url, &dependencias); err != nil {
		logs.Error(err)
		return models.Dependencia{}, err
	}
	if len(dependencias) == 0 {
		return models.Dependencia{}, fmt.Errorf("dependencia incorrectamente homologada para proyecto curricular %d", idProyectoCurricular)
	}

	return dependencias[0], nil
}

func registrarEstadoContrato(numeroContrato string, vigencia int, usuario string, estadoID int) error {
	ce := models.ContratoEstado{
		NumeroContrato: numeroContrato,
		Vigencia:       vigencia,
		FechaRegistro:  time.Now(),
		Usuario:        usuario,
		Estado: &models.EstadoContrato{
			Id: estadoID,
		},
	}

	var response models.ContratoEstado
	if err := SendRequestLegacy("UrlcrudAgora", "contrato_estado", "POST", &response, &ce); err != nil {
		logs.Error(err)
		return err
	}

	return nil
}

func registrarActaInicioContrato(numeroContrato string, vigencia int, acta *models.ActaInicio, fechaFin time.Time, usuario string) error {
	ai := models.ActaInicio{
		NumeroContrato: numeroContrato,
		Vigencia:       vigencia,
		Descripcion:    acta.Descripcion,
		FechaInicio:    acta.FechaInicio,
		FechaFin:       fechaFin,
		FechaRegistro:  time.Now(),
		Usuario:        usuario,
	}

	var response models.ActaInicio
	if err := SendRequestLegacy("UrlcrudAgora", "acta_inicio", "POST", &response, &ai); err != nil {
		logs.Error(err)
		return err
	}

	return nil
}

func registrarContratoDisponibilidad(numeroContrato string, vigencia int, vinculacionID int) error {
	var disponibilidad models.DisponibilidadVinculacion
	var disponibilidades []models.DisponibilidadVinculacion

	url := "disponibilidad_vinculacion?query=VinculacionDocenteId.Id:" + strconv.Itoa(vinculacionID)
	if err := GetRequestNew("UrlCrudResoluciones", url, &disponibilidades); err != nil {
		logs.Error(err)
		return err
	}
	if len(disponibilidades) == 0 {
		return fmt.Errorf("no se encontró disponibilidad para la vinculación %d", vinculacionID)
	}

	disponibilidad = disponibilidades[0]
	contratoDisponibilidad := models.ContratoDisponibilidad{
		NumeroContrato: numeroContrato,
		Vigencia:       vigencia,
		Estado:         true,
		FechaRegistro:  time.Now(),
		NumeroCdp:      int(disponibilidad.Disponibilidad),
		VigenciaCdp:    vigencia,
	}

	var response models.ContratoDisponibilidad
	if err := SendRequestLegacy("UrlcrudAgora", "contrato_disponibilidad", "POST", &response, &contratoDisponibilidad); err != nil {
		logs.Error(err)
		return err
	}

	return nil
}

func prepararContratoGeneralBase(contrato *models.ContratoGeneral, numeroContrato string, vigencia int, valorContrato float64, ordenadorGastoID int) {
	contrato.VigenciaContrato = vigencia
	contrato.Id = numeroContrato
	contrato.FormaPago.Id = 240
	contrato.DescripcionFormaPago = "Abono a Cuenta Mensual de acuerdo a puntos y horas laboradas"
	contrato.Justificacion = "Docente de Vinculación Especial"
	contrato.UnidadEjecucion.Id = 269
	contrato.LugarEjecucion.Id = 4
	contrato.TipoControl = 181
	contrato.ClaseContratista = 33
	contrato.TipoMoneda = 137
	contrato.OrigenRecursos = 149
	contrato.OrigenPresupueso = 156
	contrato.TemaGastoInversion = 166
	contrato.TipoGasto = 146
	contrato.RegimenContratacion = 136
	contrato.Procedimiento = 132
	contrato.ModalidadSeleccion = 123
	contrato.TipoCompromiso = 35
	contrato.TipologiaContrato = 46
	contrato.FechaRegistro = time.Now()
	contrato.UnidadEjecutora = 1
	contrato.ValorContrato = valorContrato
	contrato.OrdenadorGasto = ordenadorGastoID
	contrato.Condiciones = "Sin condiciones"
}

func construirPayloadContratoGeneral(contrato *models.ContratoGeneral, contratistaID int, supervisor models.SupervisorContrato, tipoContrato models.TipoContrato) map[string]interface{} {
	return map[string]interface{}{
		"Id":               contrato.Id,
		"VigenciaContrato": contrato.VigenciaContrato,
		"ObjetoContrato":   contrato.ObjetoContrato,
		"PlazoEjecucion":   contrato.PlazoEjecucion,
		"FormaPago": map[string]interface{}{
			"Id":                240,
			"Descripcion":       "TRANSACCIÓN",
			"CodigoContraloria": "'",
			"EstadoRegistro":    true,
			"FechaRegistro":     "2016-10-25T00:00:00Z",
		},
		"OrdenadorGasto":         contrato.OrdenadorGasto,
		"SedeSolicitante":        contrato.SedeSolicitante,
		"DependenciaSolicitante": contrato.DependenciaSolicitante,
		"Contratista":            contratistaID,
		"UnidadEjecucion": map[string]interface{}{
			"Id":                269,
			"Descripcion":       "Semana(s)",
			"CodigoContraloria": "'",
			"EstadoRegistro":    true,
			"FechaRegistro":     "2018-03-20T00:00:00Z",
		},
		"ValorContrato":        int(contrato.ValorContrato),
		"Justificacion":        contrato.Justificacion,
		"DescripcionFormaPago": contrato.DescripcionFormaPago,
		"Condiciones":          contrato.Condiciones,
		"UnidadEjecutora":      contrato.UnidadEjecutora,
		"FechaRegistro":        contrato.FechaRegistro.Format(time.RFC3339),
		"TipologiaContrato":    contrato.TipologiaContrato,
		"TipoCompromiso":       contrato.TipoCompromiso,
		"ModalidadSeleccion":   contrato.ModalidadSeleccion,
		"Procedimiento":        contrato.Procedimiento,
		"RegimenContratacion":  contrato.RegimenContratacion,
		"TipoGasto":            contrato.TipoGasto,
		"TemaGastoInversion":   contrato.TemaGastoInversion,
		"OrigenPresupueso":     contrato.OrigenPresupueso,
		"OrigenRecursos":       contrato.OrigenRecursos,
		"TipoMoneda":           contrato.TipoMoneda,
		"TipoControl":          contrato.TipoControl,
		"Observaciones":        contrato.Observaciones,
		"Supervisor": map[string]interface{}{
			"Id":                    supervisor.Id,
			"Nombre":                supervisor.Nombre,
			"Documento":             supervisor.Documento,
			"Cargo":                 supervisor.Cargo,
			"SedeSupervisor":        supervisor.SedeSupervisor,
			"DependenciaSupervisor": supervisor.DependenciaSupervisor,
			"Tipo":                  supervisor.Tipo,
			"Estado":                supervisor.Estado,
			"DigitoVerificacion":    supervisor.DigitoVerificacion,
			"FechaInicio":           supervisor.FechaInicio,
			"FechaFin":              supervisor.FechaFin,
			"CargoId": map[string]interface{}{
				"Id": supervisor.CargoId.Id,
			},
		},
		"ClaseContratista": contrato.ClaseContratista,
		"TipoContrato": map[string]interface{}{
			"Id":           tipoContrato.Id,
			"TipoContrato": tipoContrato.TipoContrato,
			"Estado":       true,
		},
		"LugarEjecucion": map[string]interface{}{
			"Id":          4,
			"Direccion":   "CALLE 40 A No 13-09",
			"Sede":        "00IP",
			"Dependencia": "DEP39",
			"Ciudad":      96,
		},
	}
}

func registrarContratoGeneral(contrato *models.ContratoGeneral, contratistaID int, supervisor models.SupervisorContrato, tipoContrato models.TipoContrato) error {
	var response interface{}
	payload := construirPayloadContratoGeneral(contrato, contratistaID, supervisor, tipoContrato)
	if err := SendRequestLegacy("UrlcrudAgora", "contrato_general", "POST", &response, payload); err != nil {
		logs.Error(payload)
		logs.Error(response)
		return err
	}

	return nil
}

func actualizarResolucionExpedida(resolucion *models.Resolucion, usuario string, fechaExpedicion time.Time, normalizarRango bool) error {
	var response interface{}

	if err := CambiarEstadoResolucion(resolucion.Id, "REXP", usuario); err != nil {
		logs.Error(response)
		return err
	}

	resolucion.FechaExpedicion = fechaExpedicion
	if normalizarRango {
		resolucion.FechaInicio, resolucion.FechaFin = NormalizarFechasResolucion(resolucion.FechaInicio, resolucion.FechaFin)
	} else if fecha := NormalizarFechaTimezone(resolucion.FechaInicio); fecha != nil {
		resolucion.FechaInicio = fecha
	}

	url := ResolucionEndpoint + strconv.Itoa(resolucion.Id)
	if err := SendRequestNew("UrlCrudResoluciones", url, "PUT", &response, resolucion); err != nil {
		logs.Error(response)
		return err
	}

	documento, outputError := AlmacenarResolucionGestorDocumental(resolucion.Id)
	if outputError != nil {
		logs.Error(response)
		return fmt.Errorf("%v", outputError)
	}

	resolucion.NuxeoUid = documento.Enlace
	if err := SendRequestNew("UrlCrudResoluciones", url, "PUT", &response, resolucion); err != nil {
		logs.Error(response)
		return err
	}

	return nil
}

func notificarDocentesAsync(datosCorreo []models.EmailData, codigoAbreviacion string) {
	go func() {
		if err := NotificarDocentes(datosCorreo, codigoAbreviacion); err != nil {
			logs.Error(err)
		}
	}()
}

func persistirContratoExpedido(numeroContrato string, vigencia int, usuario string, acta *models.ActaInicio, fechaFin time.Time, vinculacionID int) error {
	if err := registrarEstadoContrato(numeroContrato, vigencia, usuario, 4); err != nil {
		return err
	}

	if err := registrarActaInicioContrato(numeroContrato, vigencia, acta, fechaFin, usuario); err != nil {
		return err
	}

	if err := registrarContratoDisponibilidad(numeroContrato, vigencia, vinculacionID); err != nil {
		return err
	}

	return nil
}

func ExpedirResolucion(m models.ExpedicionResolucion) (outputError map[string]interface{}) {
	defer func() {
		if err := recover(); err != nil {
			outputError = map[string]interface{}{"funcion": "/ExpedirResolucion", "err": err, "status": "502"}
			panic(outputError)
		}
	}()

	var cdve int
	var response interface{}
	var datosCorreo []models.EmailData
	var contexto contextoExpedicionResolucion

	vigencia, _, _ := time.Now().Date()
	vin := m.Vinculaciones

	usuario, err := resolverUsuarioExpedicion(m.Usuario)
	if err != nil {
		panic(err.Error())
	}

	contexto, err = cargarContextoExpedicionResolucion(m.IdResolucion)
	if err != nil {
		panic(err.Error())
	}
	r := contexto.resolucion
	parametro := contexto.tipoResolucion
	ordenadorGasto := contexto.ordenadorGasto
	dep := contexto.facultad

	url := ""

	url = "contrato_general/maximo_dve"
	if err := GetRequestLegacy("UrlcrudAgora", url, &cdve); err == nil { // If 1 - consecutivo contrato_general
		numeroContratos := cdve
		for _, vinculacion := range vin { // For vinculaciones
			numeroContratos = numeroContratos + 1
			if v, err := cargarVinculacionPorID(vinculacion.VinculacionDocente.Id); err == nil { // If 1.1 - vinculacion_docente
				contrato := vinculacion.ContratoGeneral
				if tipoCon, err := cargarTipoContratoAgora(contrato.TipoContrato.Id); err == nil { // If 1.1.1 - tipoContrato
					var sup models.SupervisorContrato
					acta := vinculacion.ActaInicio
					prepararContratoGeneralBase(contrato, "DVE"+strconv.Itoa(numeroContratos), vigencia, v.ValorContrato, ordenadorGasto.Id)
					sup, err := SupervisorActual(r.DependenciaId)
					if err != nil { // If 1.1.2 - supervisorActual
						logs.Error(err)
						panic(err)
					}
					datoMensaje := models.EmailData{
						Documento:        strconv.FormatFloat(v.PersonaId, 'f', 0, 64),
						ContratoId:       contrato.Id,
						Facultad:         dep.Nombre,
						NumeroResolucion: r.NumeroResolucion,
					}
					datosCorreo = append(datosCorreo, datoMensaje)
					contrato.Supervisor = &sup
					if proveedor, err := cargarProveedorAgora(contrato.Contratista); err == nil { // If 1.1.3 - informacion_proveedor
						if proveedor.Id > 0 { // If 1.1.4 - proveedor
							temp := proveedor.Id
							if err := registrarContratoGeneral(contrato, temp, sup, tipoCon); err == nil { // If 1.1.5 contrato_general
								aux1 := contrato.Id
								aux2 := contrato.VigenciaContrato
								fechasContrato := CalcularFechasContrato(acta.FechaInicio, v.NumeroSemanas)
								if err := persistirContratoExpedido(aux1, aux2, usuario["documento_compuesto"].(string), acta, fechasContrato.FechaFinPago, v.Id); err == nil { // If 1.1.6 - 1.1.9
									v.NumeroContrato = &aux1
									v.Vigencia = aux2
									v.FechaInicio = acta.FechaInicio
									if fecha := NormalizarFechaTimezone(&v.FechaInicio); fecha != nil {
										v.FechaInicio = *fecha
									}
									url = VinculacionEndpoint + strconv.Itoa(v.Id)
									if err := SendRequestNew("UrlCrudResoluciones", url, "PUT", &response, &v); err != nil { // If 1.1.10 - vinculacion_docente
										logs.Error(err)
										panic(err.Error())
									}
								} else { // If 1.1.6
									logs.Error(err)
									panic(err.Error())
								}
							} else { // If 1.1.5
								logs.Error(response)
								panic(err.Error())
							}
						} else { // If 1.1.4 proveedor
						}
					} else { // If 1.1.3
						panic(err.Error())
					}
				} else { // If 1.1.1
					logs.Error(tipoCon)
					panic(err.Error())
				}
			} else { // If 1.1
				logs.Error(v)
				panic(err.Error())
			}
		} // For vinculaciones
		if err := actualizarResolucionExpedida(&r, m.Usuario, m.FechaExpedicion, true); err != nil {
			panic(err)
		}
		notificarDocentesAsync(datosCorreo, parametro.CodigoAbreviacion)
	} else { // If 1
		logs.Error(cdve)
		panic(err.Error())
	}
	return nil
}

func ValidarDatosExpedicion(m models.ExpedicionResolucion) (outputError map[string]interface{}) {
	defer func() {
		if err := recover(); err != nil {
			outputError = map[string]interface{}{"funcion": "/ValidarDatosExpedicion", "err": err, "status": "502"}
			panic(outputError)
		}
	}()

	vinc := m.Vinculaciones

	for _, vinculacion := range vinc {
		v, err := cargarVinculacionPorID(vinculacion.VinculacionDocente.Id)
		if err != nil { // If 1- vinculacion_docente
			beego.Error("Error en If 1 - Previnculación no valida", err)
			logs.Error(v)
			outputError = map[string]interface{}{"funcion": "/ValidarDatosExpedicion", "err": err.Error(), "status": "502"}
			return outputError
		}

		contrato := vinculacion.ContratoGeneral
		if _, err := cargarProveedorAgora(contrato.Contratista); err != nil { // if 2 - informacion_proveedor
			beego.Error("Error en If 2 - Docente no válido en Ágora, se encuentra identificado con el documento número ", strconv.Itoa(contrato.Contratista), err)
			outputError = map[string]interface{}{"funcion": "/ValidarDatosExpedicion", "err": err.Error(), "status": "502"}
			return outputError
		}

		if _, err := cargarDependenciaProyectoCurricular(v.ProyectoCurricularId); err != nil { // If 6
			beego.Error("Error en If 6 - Dependencia incorrectamente homologada asociada al docente identificado con "+strconv.Itoa(contrato.Contratista)+" en Ágora", err)
			outputError = map[string]interface{}{"funcion": "/ValidarDatosExpedicion6", "err": err.Error(), "status": "502"}
			return outputError
		}
	}
	return
}

func ExpedirModificacion(m models.ExpedicionResolucion) (outputError map[string]interface{}) {
	defer func() {
		if err := recover(); err != nil {
			if errMap, ok := err.(map[string]interface{}); ok {
				outputError = errMap
			} else {
				outputError = map[string]interface{}{
					"funcion": "/ExpedirCancelacion",
					"err":     err,
					"status":  "502",
				}
			}
			panic(outputError)
		}
	}()

	var temp int
	var cdve int
	var response interface{}
	var tipoRes models.Parametro
	var reduccion *models.Reduccion
	var resolucion models.Resolucion
	var ordenadorGasto models.OrdenadorGasto
	var datosCorreo []models.EmailData
	var parametro models.Parametro
	var dep models.Dependencia
	var contexto contextoExpedicionResolucion
	vigencia, _, _ := time.Now().Date()
	vinc := m.Vinculaciones

	var usuario map[string]interface{}
	var err error

	usuario, err = resolverUsuarioExpedicion(m.Usuario)
	if err != nil {
		panic(err.Error())
	}

	contexto, err = cargarContextoExpedicionResolucion(m.IdResolucion)
	if err != nil {
		panic(err.Error())
	}
	resolucion = contexto.resolucion
	parametro = contexto.tipoResolucion
	tipoRes = contexto.tipoResolucion

	url := ""

	if tipoRes.CodigoAbreviacion != "RVIN" {
		if err := ValidarLiquidacionesTitan(vinc); err != nil {
			panic(err)
		}
	}

	ordenadorGasto = contexto.ordenadorGasto
	dep = contexto.facultad
	url = "contrato_general/maximo_dve"
	if err := GetRequestLegacy("UrlcrudAgora", url, &cdve); err == nil { // If 1 - consecutivo contrato_general
		numeroContratos := cdve
		for _, vinculacion := range vinc {
			numeroContratos = numeroContratos + 1
			contrato := vinculacion.ContratoGeneral
			if modificacion, err := cargarVinculacionPorQuery(vinculacion.VinculacionDocente.Id); err == nil { // If 1.1 - vinculacion_docente
				if tipoCon, err := cargarTipoContratoAgora(contrato.TipoContrato.Id); err == nil {
					contrato := vinculacion.ContratoGeneral
					var sup models.SupervisorContrato
					acta := vinculacion.ActaInicio
					prepararContratoGeneralBase(contrato, "DVE"+strconv.Itoa(numeroContratos), vigencia, math.Floor(modificacion.ValorContrato), ordenadorGasto.Id)
					sup, err := SupervisorActual(resolucion.DependenciaId)
					if err != nil {
						logs.Error(err)
						panic(err)
					}
					contrato.Supervisor = &sup
					if proveedor, err := cargarProveedorAgora(contrato.Contratista); err == nil { // If 1.2
						if proveedor.Id > 0 { // If 1.3 - proveedor != nil
							temp = proveedor.Id
							if tipoRes.CodigoAbreviacion == "RRED" {
								//horasFinales := 0
								horasReducir := modificacion.NumeroHorasSemanales
								dedicacion := modificacion.ResolucionVinculacionDocenteId.Dedicacion
								nivel := modificacion.ResolucionVinculacionDocenteId.NivelAcademico
								reduccion = &models.Reduccion{
									Vigencia:          modificacion.Vigencia,
									Documento:         fmt.Sprintf("%.f", modificacion.PersonaId),
									FechaReduccion:    modificacion.FechaInicio,
									Semanas:           modificacion.NumeroSemanas,
									SemanasAnteriores: 0,
									NivelAcademico:    nivel,
								}

								contratosAnteriores := new([]models.VinculacionDocente)
								if err := BuscarContratosModificar(modificacion.Id, contratosAnteriores); err != nil {
									panic(err.Error())
								}
								var respActaInicioAnterior []models.ActaInicio
								var actaInicioAnterior models.ActaInicio
								var resolucion models.Resolucion
								var parametro models.Parametro
								var ultimaVinculacon models.VinculacionDocente
								var aux = 0
								//horasAnterior := 0
								horasNuevo := 0
								for _, subcontratoAux := range *contratosAnteriores {
									if aux == 0 {
										ultimaVinculacon = subcontratoAux
									}
									url := "resolucion/" + strconv.Itoa(subcontratoAux.ResolucionVinculacionDocenteId.Id)
									if err := GetRequestNew("UrlCrudResoluciones", url, &resolucion); err != nil {
										logs.Error(err.Error())
										panic(err.Error())
									}
									url3 := "parametro/" + strconv.Itoa(resolucion.TipoResolucionId)
									if err := GetRequestNew("UrlcrudParametros", url3, &parametro); err != nil {
										logs.Error(err)
										panic("Cargando tipo_resolucion -> " + err.Error())
									}
									if parametro.CodigoAbreviacion == "RVIN" || parametro.CodigoAbreviacion == "RADD" {
										horasNuevo += subcontratoAux.NumeroHorasSemanales
									} else if parametro.CodigoAbreviacion == "RRED" {
										horasNuevo -= subcontratoAux.NumeroHorasSemanales
									}
									aux += 1
									reduccion.SemanasAnteriores = ultimaVinculacon.NumeroSemanas
									url = fmt.Sprintf("acta_inicio?query=NumeroContrato:%s,Vigencia:%d", *ultimaVinculacon.NumeroContrato, ultimaVinculacon.Vigencia)
									if err := GetRequestLegacy("UrlcrudAgora", url, &respActaInicioAnterior); err != nil {
										panic(err.Error())
									} else if len(respActaInicioAnterior) == 0 {
										panic("Acta de inicio no encontrada")
									}
									horasNuevoAux := horasNuevo - horasReducir
									actaInicioAnterior = respActaInicioAnterior[0]
									if actaInicioAnterior.FechaInicio.Before(modificacion.FechaInicio) || actaInicioAnterior.FechaInicio.Equal(modificacion.FechaInicio) {
										valores := make(map[string]float64)
										contratoReducir := &models.ContratoReducir{
											NumeroContratoOriginal: *subcontratoAux.NumeroContrato,
										}
										if modificacion.ResolucionVinculacionDocenteId.Dedicacion != "HCH" {
											// calcular el desagregado del resto de cada contrato
											//var subcontratoAux = ultimaVinculacon
											var desagregado, err map[string]interface{}
											if nivel == "POSGRADO" {
												/*horasXSemana := subcontrato.NumeroHorasSemanales / subcontrato.NumeroSemanas
												horasAntesReduccion := horasXSemana * (subcontrato.NumeroSemanas - modificacion.NumeroSemanas)
												subcontrato.NumeroHorasSemanales = horasAntesReduccion*/
												subcontratoAux.NumeroHorasSemanales = modificacion.NumeroHorasTrabajadas

											} else {
												//ultimaVinculacon.NumeroSemanas = ultimaVinculacon.NumeroSemanas - modificacion.NumeroSemanas
												subcontratoAux.NumeroSemanas = subcontratoAux.NumeroSemanas - modificacion.NumeroSemanas
												//subcontratoAux.NumeroHorasSemanales = horasAnterior
											}
											if desagregado, err = CalcularDesagregadoTitan(subcontratoAux, dedicacion, nivel); err != nil {
												panic(err)
											}
											for concepto, valor := range desagregado {
												if concepto != "NumeroContrato" && concepto != "Vigencia" {
													if concepto == "SueldoBasico" {
														contratoReducir.ValorContratoReducido = valor.(float64)
													} else {
														valores[concepto] = valor.(float64)
													}
												}
											}
											contratoReducir.DesagregadoOriginal = &valores
										} else {
											var vin []models.VinculacionDocente
											//var desagregado, err map[string]interface{}
											// subcontrato.NumeroSemanas = subcontrato.NumeroSemanas - modificacion.NumeroSemanas
											nuevaVinculacion := models.VinculacionDocente{
												Vigencia:                       ultimaVinculacon.Vigencia,
												PersonaId:                      ultimaVinculacon.PersonaId,
												NumeroHorasSemanales:           subcontratoAux.NumeroHorasSemanales,
												NumeroSemanas:                  subcontratoAux.NumeroSemanas - modificacion.NumeroSemanas,
												ResolucionVinculacionDocenteId: ultimaVinculacon.ResolucionVinculacionDocenteId,
												Categoria:                      ultimaVinculacon.Categoria,
												Activo:                         true,
												ValorPuntoSalarial:             ultimaVinculacon.ValorPuntoSalarial,
											}
											if nivel == "POSGRADO" {
												horasXSemana := subcontratoAux.NumeroHorasSemanales / subcontratoAux.NumeroSemanas
												// horasRestantesTotales := subcontrato.NumeroHorasSemanales - modificacion.NumeroHorasSemanales
												horasAntesReduccion := horasXSemana * (subcontratoAux.NumeroSemanas - modificacion.NumeroSemanas)
												nuevaVinculacion.NumeroHorasSemanales = horasAntesReduccion
											} else {
												nuevaVinculacion.NumeroSemanas = subcontratoAux.NumeroSemanas - modificacion.NumeroSemanas
											}
											vin = append(vin, nuevaVinculacion)
											if w, err2 := CalcularSalarioPrecontratacion(vin); err2 == nil {
												vin = nil
												contratoReducir.ValorContratoReducido = w[0].ValorContrato
											} else {
												panic(err2)
											}
										}
										reduccion.ContratosOriginales = append(reduccion.ContratosOriginales, *contratoReducir)
										// actualizacion acta_inicio
										// fechaFinOriginal := actaInicioAnterior.FechaFin
										actaInicioAnterior.FechaFin = modificacion.FechaInicio
										actaInicioAnterior.Usuario = usuario["documento_compuesto"].(string)
										url = "acta_inicio/" + strconv.Itoa(actaInicioAnterior.Id)
										if err := SendRequestLegacy("UrlcrudAgora", url, "PUT", &response, &actaInicioAnterior); err != nil {
											panic(err.Error())
										}
										if horasNuevoAux > 0 && horasNuevoAux != 0 {
											// Calcula el valor del nuevo contrato con base en las semanas desde la fecha inicio escogida hasta la nueva fecha fin y las nuevas horas
											semanasTranscurridas := math.Ceil(modificacion.FechaInicio.Sub(actaInicioAnterior.FechaInicio).Hours() / (24 * 7)) // cálculo con base en semanas
											semanasRestantes := ultimaVinculacon.NumeroSemanas - modificacion.NumeroSemanas - int(semanasTranscurridas)
											var vinc [1]models.VinculacionDocente
											vinc[0] = models.VinculacionDocente{
												ResolucionVinculacionDocenteId: modificacion.ResolucionVinculacionDocenteId,
												PersonaId:                      modificacion.PersonaId,
												NumeroHorasSemanales:           horasNuevoAux,
												NumeroSemanas:                  modificacion.NumeroSemanas,
												Vigencia:                       modificacion.Vigencia,
												Categoria:                      modificacion.Categoria,
												ValorPuntoSalarial:             modificacion.ValorPuntoSalarial,
											}
											if nivel == "POSGRADO" {
												vinc[0].NumeroHorasSemanales = horasNuevo - modificacion.NumeroHorasSemanales - modificacion.NumeroHorasTrabajadas
											}
											salario, err := CalcularValorContratoReduccion(vinc, semanasRestantes, ultimaVinculacon.NumeroHorasSemanales, nivel)
											if err != nil {
												panic(err)
											}

											contrato.ValorContrato = math.Floor(salario)
											beego.Info(contrato.ValorContrato)
											// el subcontrato actual es reducido parcialmente y los siguientes no deben ser afectados
											var desagregadoNuevo, err2 map[string]interface{}
											if desagregadoNuevo, err2 = CalcularDesagregadoTitan(vinc[0], dedicacion, nivel); err2 != nil {
												panic(err)
											}
											reduccion.ContratoNuevo = &models.ContratoReducido{}
											valoresNuevo := make(map[string]float64)
											for concepto, valor := range desagregadoNuevo {
												if concepto != "NumeroContrato" && concepto != "Vigencia" {
													if concepto == "SueldoBasico" {
														reduccion.ContratoNuevo.ValorContratoReduccion = valor.(float64)
													} else {
														valoresNuevo[concepto] = valor.(float64)
													}
												}
											}
											reduccion.ContratoNuevo.DesagregadoReduccion = &valoresNuevo
											break
										} else if horasNuevoAux == 0 {
											break
										}
									}
								}
							}
							if contrato.ValorContrato > 0 {
								if err := registrarContratoGeneral(contrato, temp, sup, tipoCon); err == nil { // If 1.8 - contrato_general (POST)
									numContrato := contrato.Id
									vigencia := contrato.VigenciaContrato
									actaModificacion := &models.ActaInicio{
										Descripcion: acta.Descripcion,
										FechaInicio: modificacion.FechaInicio,
									}
									fechasContrato := CalcularFechasContrato(modificacion.FechaInicio, modificacion.NumeroSemanas)
									if err := persistirContratoExpedido(numContrato, vigencia, usuario["documento_compuesto"].(string), actaModificacion, fechasContrato.FechaFinPago, modificacion.Id); err == nil { // If 1.9 - 1.12
										modificacion.NumeroContrato = &numContrato
										modificacion.Vigencia = vigencia
										url = VinculacionEndpoint + strconv.Itoa(modificacion.Id)
										if fecha := NormalizarFechaTimezone(&modificacion.FechaInicio); fecha != nil {
											modificacion.FechaInicio = *fecha
										}
										if err := SendRequestNew("UrlCrudResoluciones", url, "PUT", &response, &modificacion); err != nil {
											logs.Error(response)
											panic(err.Error())
										}
										if tipoRes.CodigoAbreviacion == "RRED" {
											if reduccion.ContratoNuevo != nil {
												reduccion.ContratoNuevo.NumeroContratoReduccion = numContrato
												reduccion.ContratoNuevo.NumeroResolucion = resolucion.NumeroResolucion
												reduccion.ContratoNuevo.IdResolucion = resolucion.Id
											}
										}
									} else { // If 1.9
										logs.Error(err)
										panic(err.Error())
									}
								} else { // if 1.8
									panic(err.Error())
								}
							} else {
								reduccion.ContratoNuevo = nil
							}
						}
					} else { // If 1.2
						panic(err.Error())
					}
					if tipoRes.CodigoAbreviacion == "RRED" {
						if err := ReducirContratosTitan(reduccion, &modificacion, contrato.ValorContrato); err != nil {
							panic(err)
						}
					}
					datoMensaje := models.EmailData{
						Documento:        strconv.FormatFloat(modificacion.PersonaId, 'f', 0, 64),
						ContratoId:       *modificacion.NumeroContrato,
						Facultad:         dep.Nombre,
						NumeroResolucion: resolucion.NumeroResolucion,
					}
					datosCorreo = append(datosCorreo, datoMensaje)
				}
			} else { // If 1.1
				panic(err.Error())
			}
		}
		if err := actualizarResolucionExpedida(&resolucion, m.Usuario, m.FechaExpedicion, false); err != nil {
			panic(err)
		}
		notificarDocentesAsync(datosCorreo, parametro.CodigoAbreviacion)
	} else {
		logs.Error(cdve)
		panic(err.Error())
	}
	return
}

func ExpedirCancelacion(m models.ExpedicionCancelacion) (outputError map[string]interface{}) {
	defer func() {
		if err := recover(); err != nil {
			if errMap, ok := err.(map[string]interface{}); ok {
				outputError = errMap
			} else {
				outputError = map[string]interface{}{
					"funcion": "/ExpedirCancelacion",
					"err":     err,
					"status":  "502",
				}
			}
			panic(outputError)
		}
	}()

	cancelaciones := m.Vinculaciones

	vinculacionIds := make([]int, 0, len(cancelaciones))
	for _, can := range cancelaciones {
		if can == nil {
			panic(map[string]interface{}{
				"funcion": "/ExpedirCancelacion",
				"err":     "Se recibió una cancelación nula",
				"status":  "400",
			})
		}
		vinculacionIds = append(vinculacionIds, can.VinculacionDocente.Id)
	}

	if err := ValidarLiquidacionesTitanPorIds(vinculacionIds); err != nil {
		panic(err)
	}

	var response interface{}
	var usuario map[string]interface{}
	var datosCorreo []models.EmailData
	var parametro models.Parametro
	var resolucion models.Resolucion
	var dep models.Dependencia
	var contexto contextoExpedicionResolucion
	var err error

	usuario, err = resolverUsuarioExpedicion(cancelaciones[0].ContratoCancelado.Usuario)
	if err != nil {
		panic(err.Error())
	}

	contexto, err = cargarContextoExpedicionResolucion(m.IdResolucion)
	if err != nil {
		panic(err.Error())
	}
	resolucion = contexto.resolucion
	parametro = contexto.tipoResolucion
	dep = contexto.facultad

	for _, can := range cancelaciones {
		v, err := cargarVinculacionPorQuery(can.VinculacionDocente.Id)
		if err != nil {
			panic("Vinculacion (cancelacion) -> " + err.Error())
		}
		contratos := new([]models.VinculacionDocente)
		if err := BuscarContratosModificar(v.Id, contratos); err == nil { // If 1 - vinculacion_docente
			for _, contrato := range *contratos {
				url := "acta_inicio?query=NumeroContrato:" + *contrato.NumeroContrato + ",Vigencia:" + strconv.Itoa(contrato.Vigencia)
				var ai []models.ActaInicio
				if err := GetRequestLegacy("UrlcrudAgora", url, &ai); err != nil {
					panic("Acta de inicio -> " + err.Error())
				} else if len(ai) == 0 {
					panic("Acta de inicio no encontrada")
				}
				actaInicio := ai[0]
				if actaInicio.FechaFin.After(v.FechaInicio) {
					contratoCancelado := &models.ContratoCancelado{
						NumeroContrato:    *contrato.NumeroContrato,
						Vigencia:          contrato.Vigencia,
						FechaCancelacion:  v.FechaInicio,
						MotivoCancelacion: can.ContratoCancelado.MotivoCancelacion,
						Usuario:           usuario["documento_compuesto"].(string),
						FechaRegistro:     time.Now(),
						Estado:            can.ContratoCancelado.Estado,
					}
					url = "contrato_cancelado"
					if err := SendRequestLegacy("UrlcrudAgora", url, "POST", &response, &contratoCancelado); err == nil { // If 2 -contrato_cancelado (post)
						actaInicio.FechaFin = v.FechaInicio
						url = "acta_inicio/" + strconv.Itoa(actaInicio.Id)
						if err := SendRequestLegacy("UrlcrudAgora", url, "PUT", &response, &actaInicio); err == nil { // If 4 - acta_inicio
							if err := registrarEstadoContrato(contratoCancelado.NumeroContrato, contratoCancelado.Vigencia, usuario["documento_compuesto"].(string), 7); err == nil { // If 5 - contrato_estado
								datoMensaje := models.EmailData{
									Documento:        strconv.FormatFloat(v.PersonaId, 'f', 0, 64),
									ContratoId:       "",
									Facultad:         dep.Nombre,
									NumeroResolucion: resolucion.NumeroResolucion,
								}
								datosCorreo = append(datosCorreo, datoMensaje)
								if err := actualizarResolucionExpedida(&resolucion, can.ContratoCancelado.Usuario, m.FechaExpedicion, false); err == nil {
									if err := ReliquidarContratoCancelado(v, contrato); err != nil {
										panic(err)
									}
								} else { // If 7
									logs.Error(err)
									panic(err.Error())
								}
							} else { // If 5
								logs.Error(response)
								panic(err.Error())
							}
						} else { // If 4
							logs.Error(response)
							panic(err.Error())
						}
					} else { // if 2
						logs.Error(response)
						panic(err.Error())
					}
				}
			}
		} else { // If 1
			logs.Error(v)
			panic(err.Error())
		}
	}
	notificarDocentesAsync(datosCorreo, parametro.CodigoAbreviacion)

	return
}

// Función que recopila los contratos a cancelar de acuerdo con el histórico de modificaciones
func BuscarContratosModificar(vinculacionId int, contratos *[]models.VinculacionDocente) error {
	var modificaciones []models.ModificacionVinculacion
	var modVin models.ModificacionVinculacion
	var tipoResolucion models.Parametro

	url := "modificacion_vinculacion?query=VinculacionDocenteRegistradaId.Id:" + strconv.Itoa(vinculacionId)
	if err := GetRequestNew("UrlCrudResoluciones", url, &modificaciones); err != nil {
		logs.Error(err.Error())
		return err
	}

	// Caso de salida
	if len(modificaciones) == 0 {
		return nil
	}

	modVin = modificaciones[0]

	url2 := ParametroEndpoint + strconv.Itoa(modVin.ModificacionResolucionId.ResolucionAnteriorId.TipoResolucionId)
	if err2 := GetRequestNew("UrlcrudParametros", url2, &tipoResolucion); err2 != nil {
		logs.Error(err2.Error())
		return err2
	}

	*contratos = append(*contratos, *modVin.VinculacionDocenteCanceladaId)

	// Segundo caso de salida
	if tipoResolucion.CodigoAbreviacion == "RVIN" {
		return nil
	}

	// Llamada recursiva para consultar una modificación anterior hasta llegar a
	// la vinculación inicial que no tiene modificaciones
	return BuscarContratosModificar(modVin.VinculacionDocenteCanceladaId.Id, contratos)
}
