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

func ExpedirResolucion(m models.ExpedicionResolucion) (outputError map[string]interface{}) {
	defer func() {
		if err := recover(); err != nil {
			outputError = map[string]interface{}{"funcion": "/ExpedirResolucion", "err": err, "status": "502"}
			panic(outputError)
		}
	}()

	var cdve int
	var tipoCon models.TipoContrato
	var proveedor []models.InformacionProveedor
	var response interface{}
	var usuario map[string]interface{}
	var ordenadoresGasto []models.OrdenadorGasto
	var ordenadorGasto models.OrdenadorGasto
	var r models.Resolucion
	var err error
	var datosCorreo []models.EmailData
	var resv models.ResolucionVinculacionDocente
	var dep models.Dependencia
	var parametro models.Parametro

	vigencia, _, _ := time.Now().Date()
	vin := m.Vinculaciones

	usuario, err = GetUsuario(m.Usuario)
	if err != nil {
		logs.Error(err)
		panic(err.Error())
	}

	if err := GetRequestNew("UrlCrudResoluciones", ResolucionEndpoint+strconv.Itoa(m.IdResolucion), &r); err != nil {
		logs.Error(err)
		panic(err.Error())
	}

	url := "parametro/" + strconv.Itoa(r.TipoResolucionId)
	if err := GetRequestNew("UrlcrudParametros", url, &parametro); err != nil {
		logs.Error(err)
		panic("Cargando tipo_resolucion -> " + err.Error())
	}

	// Cambiar en el futuro por terceros_mid/tipo/ordenadoresGasto y filtrar por dependenciaId
	url2 := "ordenador_gasto?query=DependenciaId:" + strconv.Itoa(r.DependenciaId)
	if err := GetRequestLegacy("UrlcrudCore", url2, &ordenadoresGasto); err != nil {
		logs.Error(err)
		panic(err.Error())
	} else {
		if len(ordenadoresGasto) > 0 {
			ordenadorGasto = ordenadoresGasto[0]
		} else {
			if err := GetRequestLegacy("UrlcrudCore", "ordenador_gasto/1", &ordenadorGasto); err != nil {
				logs.Error(err)
				panic(err.Error())
			}
		}
	}

	if err := GetRequestNew("UrlCrudResoluciones", ResVinEndpoint+strconv.Itoa(m.IdResolucion), &resv); err != nil {
		logs.Error(err)
		panic(err.Error())
	}

	if err := GetRequestLegacy("UrlcrudOikos", "dependencia/"+strconv.Itoa(resv.FacultadId), &dep); err != nil {
		logs.Error(err)
		panic(err.Error())
	}

	url = "contrato_general/maximo_dve"
	if err := GetRequestLegacy("UrlcrudAgora", url, &cdve); err == nil { // If 1 - consecutivo contrato_general
		numeroContratos := cdve
		for _, vinculacion := range vin { // For vinculaciones
			numeroContratos = numeroContratos + 1
			var v models.VinculacionDocente
			url = VinculacionEndpoint + strconv.Itoa(vinculacion.VinculacionDocente.Id)
			if err := GetRequestNew("UrlCrudResoluciones", url, &v); err == nil { // If 1.1 - vinculacion_docente
				contrato := vinculacion.ContratoGeneral
				url = "tipo_contrato/" + strconv.Itoa(contrato.TipoContrato.Id)
				if err := GetRequestLegacy("UrlcrudAgora", url, &tipoCon); err == nil { // If 1.1.1 - tipoContrato
					var sup models.SupervisorContrato
					acta := vinculacion.ActaInicio
					contrato.VigenciaContrato = vigencia
					contrato.Id = "DVE" + strconv.Itoa(numeroContratos)
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
					contrato.ValorContrato = v.ValorContrato
					contrato.OrdenadorGasto = ordenadorGasto.Id
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
					contrato.Condiciones = "Sin condiciones"
					url = "informacion_proveedor?query=NumDocumento:" + strconv.Itoa(contrato.Contratista)
					if err := GetRequestLegacy("UrlcrudAgora", url, &proveedor); err == nil { // If 1.1.3 - informacion_proveedor
						if proveedor != nil { // If 1.1.4 - proveedor
							temp := proveedor[0].Id
							contratoGeneral := make(map[string]interface{})
							contratoGeneral = map[string]interface{}{
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
								"Contratista":            temp,
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
									"Id":                    sup.Id,
									"Nombre":                sup.Nombre,
									"Documento":             sup.Documento,
									"Cargo":                 sup.Cargo,
									"SedeSupervisor":        sup.SedeSupervisor,
									"DependenciaSupervisor": sup.DependenciaSupervisor,
									"Tipo":                  sup.Tipo,
									"Estado":                sup.Estado,
									"DigitoVerificacion":    sup.DigitoVerificacion,
									"FechaInicio":           sup.FechaInicio,
									"FechaFin":              sup.FechaFin,
									"CargoId": map[string]interface{}{
										"Id": sup.CargoId.Id,
									},
								},
								"ClaseContratista": contrato.ClaseContratista,
								"TipoContrato": map[string]interface{}{
									"Id":           tipoCon.Id,
									"TipoContrato": tipoCon.TipoContrato,
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
							url = "contrato_general"
							if err := SendRequestLegacy("UrlcrudAgora", url, "POST", &response, contratoGeneral); err == nil { // If 1.1.5 contrato_general
								aux1 := contrato.Id
								aux2 := contrato.VigenciaContrato
								var ce models.ContratoEstado
								var ec models.EstadoContrato
								ce.NumeroContrato = aux1
								ce.Vigencia = aux2
								ce.FechaRegistro = time.Now()
								ce.Usuario = usuario["documento_compuesto"].(string)
								ec.Id = 4
								ce.Estado = &ec
								var response2 models.ContratoEstado
								url = "contrato_estado"
								if err := SendRequestLegacy("UrlcrudAgora", url, "POST", &response2, &ce); err == nil { // If 1.1.6 contrato_estado
									var ai models.ActaInicio
									ai.NumeroContrato = aux1
									ai.Vigencia = aux2
									ai.Descripcion = acta.Descripcion
									ai.FechaInicio = acta.FechaInicio
									ai.FechaFin = acta.FechaFin
									fechasContrato := CalcularFechasContrato(acta.FechaInicio, v.NumeroSemanas)
									ai.FechaFin = fechasContrato.FechaFinPago
									ai.FechaRegistro = time.Now()
									ai.Usuario = usuario["documento_compuesto"].(string)
									var response3 models.ActaInicio
									url = "acta_inicio"
									if err := SendRequestLegacy("UrlcrudAgora", url, "POST", &response3, &ai); err == nil { // If 1.1.7 acta_inicio
										var cd models.ContratoDisponibilidad
										cd.NumeroContrato = aux1
										cd.Vigencia = aux2
										cd.Estado = true
										cd.FechaRegistro = time.Now()
										var dv models.DisponibilidadVinculacion
										var disp []models.DisponibilidadVinculacion
										url = "disponibilidad_vinculacion?query=VinculacionDocenteId.Id:" + strconv.Itoa(v.Id)
										if err := GetRequestNew("UrlCrudResoluciones", url, &disp); err == nil { // If 1.1.8 - disponibilidad_vinculacion
											dv = disp[0]
											cd.NumeroCdp = int(dv.Disponibilidad)
											cd.VigenciaCdp = aux2
											var response4 models.ContratoDisponibilidad
											url = "contrato_disponibilidad"
											if err := SendRequestLegacy("UrlcrudAgora", url, "POST", &response4, &cd); err == nil { // If 1.1.9 - contrato_disponibilidad
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
											} else { // If 1.1.9 -contrato_disponibilidad
												logs.Error(err)
												panic(err.Error())
											}
										} else { // If 1.1.8 - disponibilidad_vinculacion
											logs.Error(err)
											panic(err.Error())
										}
									} else { // If 1.1.7 acta_inicio
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
						logs.Error(proveedor)
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
		var r models.Resolucion
		url = ResolucionEndpoint + strconv.Itoa(m.IdResolucion)
		if err := GetRequestNew("UrlCrudResoluciones", url, &r); err == nil { // 1.2 - resolucion/
			if err := CambiarEstadoResolucion(r.Id, "REXP", m.Usuario); err == nil {
				r.FechaExpedicion = m.FechaExpedicion
				r.FechaInicio, r.FechaFin = NormalizarFechasResolucion(r.FechaInicio, r.FechaFin)
				url = ResolucionEndpoint + strconv.Itoa(r.Id)
				if err := SendRequestNew("UrlCrudResoluciones", url, "PUT", &response, &r); err != nil { // If 1.2.1
					logs.Error(response)
					panic(err.Error())
				}
				if documento, err := AlmacenarResolucionGestorDocumental(r.Id); err != nil {
					logs.Error(response)
					panic(err)
				} else {
					r.NuxeoUid = documento.Enlace
					url = ResolucionEndpoint + strconv.Itoa(r.Id)
					if err := SendRequestNew("UrlCrudResoluciones", url, "PUT", &response, &r); err != nil { // If 1.2.1
						logs.Error(response)
						panic(err.Error())
					}
				}
				go func() {
					if err := NotificarDocentes(datosCorreo, parametro.CodigoAbreviacion); err != nil {
						logs.Error(err)
					}
				}()
			} else { // If 1.2.1
				logs.Error(response)
				panic(err.Error())
			}
		} else { // If 1.2
			logs.Error(r)
			panic(err.Error())
		}
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
		v := vinculacion.VinculacionDocente
		url := VinculacionEndpoint + strconv.Itoa(v.Id)
		if err := GetRequestNew("UrlCrudResoluciones", url, &v); err != nil { // If 1- vinculacion_docente
			beego.Error("Error en If 1 - Previnculación no valida", err)
			logs.Error(v)
			outputError = map[string]interface{}{"funcion": "/ValidarDatosExpedicion", "err": err.Error(), "status": "502"}
			return outputError
		}

		contrato := vinculacion.ContratoGeneral
		var proveedor []models.InformacionProveedor
		url = "informacion_proveedor?query=NumDocumento:" + strconv.Itoa(contrato.Contratista)
		if err := GetRequestLegacy("UrlcrudAgora", url, &proveedor); err == nil { // if 2 - informacion_proveedor
		} else {
			beego.Error("Error en If 2 - Docente no válido en Ágora, se encuentra identificado con el documento número ", strconv.Itoa(contrato.Contratista), err)
			logs.Error(proveedor)
			outputError = map[string]interface{}{"funcion": "/ValidarDatosExpedicion", "err": err.Error(), "status": "502"}
			return outputError
		}
		if proveedor == nil { // If 3 - proveedor
			beego.Error("Error en If 3 - No existe el docente con número de documento " + strconv.Itoa(contrato.Contratista) + " en Ágora")
			logs.Error(proveedor)
			outputError = map[string]interface{}{"funcion": "/ValidarDatosExpedicion3", "err": "No existe el docente con este numero de documento", "status": "502"}
			return outputError
		}

		var proycur []models.Dependencia
		url = "dependencia?query=Id:" + strconv.Itoa(v.ProyectoCurricularId)
		if err := GetRequestLegacy("UrlcrudOikos", url, &proycur); err != nil { // If 6
			beego.Error("Error en If 6 - Dependencia incorrectamente homologada asociada al docente identificado con "+strconv.Itoa(contrato.Contratista)+" en Ágora", err)
			logs.Error(proycur)
			outputError = map[string]interface{}{"funcion": "/ValidarDatosExpedicion6", "err": err.Error(), "status": "502"}
			return outputError
		}
		if proycur == nil {
			beego.Error("Error en If 7 - Dependencia incorrectamente homologada asociada al docente identificado con " + strconv.Itoa(contrato.Contratista) + " en Ágora")
			logs.Error(proycur)
			outputError = map[string]interface{}{"funcion": "/ValidarDatosExpedicion6", "err": "Dependencia incorrectamente homologada asociada al docente", "status": "502"}
			return outputError
		}
	}
	return
}

func ExpedirModificacion(m models.ExpedicionResolucion) (outputError map[string]interface{}) {
	defer func() {
		if err := recover(); err != nil {
			outputError = map[string]interface{}{"funcion": "/ExpedirModificacion", "err": err, "status": "502"}
			panic(outputError)
		}
	}()

	var temp int
	var cdve int
	var proveedor []models.InformacionProveedor
	var response interface{}
	var tipoRes models.Parametro
	var reduccion *models.Reduccion
	var resolucion models.Resolucion
	var ordenadoresGasto []models.OrdenadorGasto
	var ordenadorGasto models.OrdenadorGasto
	var tipoCon models.TipoContrato
	var datosCorreo []models.EmailData
	var parametro models.Parametro
	var resv models.ResolucionVinculacionDocente
	var dep models.Dependencia
	vigencia, _, _ := time.Now().Date()
	vinc := m.Vinculaciones

	var usuario map[string]interface{}
	var err error

	usuario, err = GetUsuario(m.Usuario)
	if err != nil {
		logs.Error(err)
		panic(err.Error())
	}

	if err := GetRequestNew("UrlCrudResoluciones", ResolucionEndpoint+strconv.Itoa(m.IdResolucion), &resolucion); err != nil {
		logs.Error(err)
		panic("cargando resolucion -> " + err.Error())
	}

	url := "parametro/" + strconv.Itoa(resolucion.TipoResolucionId)
	if err := GetRequestNew("UrlcrudParametros", url, &parametro); err != nil {
		logs.Error(err)
		panic("Cargando tipo_resolucion -> " + err.Error())
	}

	if err := GetRequestNew("UrlcrudParametros", ParametroEndpoint+strconv.Itoa(resolucion.TipoResolucionId), &tipoRes); err != nil {
		logs.Error(err)
		panic("Cargando tipo_resolucion -> " + err.Error())
	}

	// Cambiar en el futuro por terceros_mid/tipo/ordenadoresGasto y filtrar por dependenciaId
	url2 := "ordenador_gasto?query=DependenciaId:" + strconv.Itoa(resolucion.DependenciaId)
	if err := GetRequestLegacy("UrlcrudCore", url2, &ordenadoresGasto); err != nil {
		logs.Error(err)
		panic(err.Error())
	} else {
		if len(ordenadoresGasto) > 0 {
			ordenadorGasto = ordenadoresGasto[0]
		} else {
			if err := GetRequestLegacy("UrlcrudCore", "ordenador_gasto/1", &ordenadorGasto); err != nil {
				logs.Error(err)
				panic(err.Error())
			}
		}
	}
	if err := GetRequestNew("UrlCrudResoluciones", ResVinEndpoint+strconv.Itoa(m.IdResolucion), &resv); err != nil {
		logs.Error(err)
		panic(err.Error())
	}

	if err := GetRequestLegacy("UrlcrudOikos", "dependencia/"+strconv.Itoa(resv.FacultadId), &dep); err != nil {
		logs.Error(err)
		panic(err.Error())
	}
	url = "contrato_general/maximo_dve"
	if err := GetRequestLegacy("UrlcrudAgora", url, &cdve); err == nil { // If 1 - consecutivo contrato_general
		numeroContratos := cdve
		for _, vinculacion := range vinc {
			numeroContratos = numeroContratos + 1
			var vi []models.VinculacionDocente
			url = "vinculacion_docente?query=Id:" + strconv.Itoa(vinculacion.VinculacionDocente.Id)
			contrato := vinculacion.ContratoGeneral
			if err := GetRequestNew("UrlCrudResoluciones", url, &vi); err == nil { // If 1.1 - vinculacion_docente
				url = "tipo_contrato/" + strconv.Itoa(contrato.TipoContrato.Id)
				if err := GetRequestLegacy("UrlcrudAgora", url, &tipoCon); err == nil {
					modificacion := vi[0]
					contrato := vinculacion.ContratoGeneral
					var sup models.SupervisorContrato
					acta := vinculacion.ActaInicio
					contrato.VigenciaContrato = vigencia
					contrato.Id = "DVE" + strconv.Itoa(numeroContratos)
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
					contrato.ValorContrato = math.Floor(modificacion.ValorContrato)
					contrato.OrdenadorGasto = ordenadorGasto.Id
					sup, err := SupervisorActual(resolucion.DependenciaId)
					if err != nil {
						logs.Error(err)
						panic(err)
					}
					contrato.Supervisor = &sup
					contrato.Condiciones = "Sin condiciones"
					url = "informacion_proveedor?query=NumDocumento:" + strconv.Itoa(contrato.Contratista)
					if err := GetRequestLegacy("UrlcrudAgora", url, &proveedor); err == nil { // If 1.2
						if proveedor != nil { // If 1.3 - proveedor != nil
							temp = proveedor[0].Id
							contratoGeneral := make(map[string]interface{})
							contratoGeneral = map[string]interface{}{
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
								"Contratista":            temp,
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
									"Id":                    sup.Id,
									"Nombre":                sup.Nombre,
									"Documento":             sup.Documento,
									"Cargo":                 sup.Cargo,
									"SedeSupervisor":        sup.SedeSupervisor,
									"DependenciaSupervisor": sup.DependenciaSupervisor,
									"Tipo":                  sup.Tipo,
									"Estado":                sup.Estado,
									"DigitoVerificacion":    sup.DigitoVerificacion,
									"FechaInicio":           sup.FechaInicio,
									"FechaFin":              sup.FechaFin,
									"CargoId": map[string]interface{}{
										"Id": sup.CargoId.Id,
									},
								},
								"ClaseContratista": contrato.ClaseContratista,
								"TipoContrato": map[string]interface{}{
									"Id":           tipoCon.Id,
									"TipoContrato": tipoCon.TipoContrato,
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
								contratoGeneral["ValorContrato"] = int(contrato.ValorContrato)
								if err := SendRequestLegacy("UrlcrudAgora", "contrato_general", "POST", &response, &contratoGeneral); err == nil { // If 1.8 - contrato_general (POST)
									numContrato := contrato.Id
									vigencia := contrato.VigenciaContrato
									var ce models.ContratoEstado
									var ec models.EstadoContrato
									ce.NumeroContrato = numContrato
									ce.Vigencia = vigencia
									ce.FechaRegistro = time.Now()
									ce.Usuario = usuario["documento_compuesto"].(string)
									ec.Id = 4
									ce.Estado = &ec
									if err := SendRequestLegacy("UrlcrudAgora", "contrato_estado", "POST", &response, &ce); err == nil { // If 1.9 - contrato_estado (POST)
										var ai models.ActaInicio
										ai.NumeroContrato = numContrato
										ai.Vigencia = vigencia
										ai.Descripcion = acta.Descripcion
										ai.FechaInicio = modificacion.FechaInicio
										fechasContrato := CalcularFechasContrato(modificacion.FechaInicio, modificacion.NumeroSemanas)
										ai.FechaFin = fechasContrato.FechaFinPago
										ai.Usuario = usuario["documento_compuesto"].(string)
										ai.FechaRegistro = time.Now()
										if err := SendRequestLegacy("UrlcrudAgora", "acta_inicio", "POST", &response, &ai); err == nil { // If 1.10 - acta_inicio (POST)
											var cd models.ContratoDisponibilidad
											cd.NumeroContrato = numContrato
											cd.Vigencia = vigencia
											cd.Estado = true
											cd.FechaRegistro = time.Now()
											var dv models.DisponibilidadVinculacion
											var disp []models.DisponibilidadVinculacion
											url = "disponibilidad_vinculacion?query=VinculacionDocenteId.Id:" + strconv.Itoa(modificacion.Id)
											if err := GetRequestNew("UrlCrudResoluciones", url, &disp); err == nil { // If 1.11 - DisponibilidadVinculacion
												dv = disp[0]
												cd.NumeroCdp = int(dv.Disponibilidad)
												cd.VigenciaCdp = vigencia
												if err := SendRequestLegacy("UrlcrudAgora", "contrato_disponibilidad", "POST", &response, &cd); err == nil { // If 1.12 - contrato_disponibilidad
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
												} else {
													logs.Error(cd)
													panic(err.Error())
												}
											} else {
												logs.Error(dv)
												panic(err.Error())
											}
										} else { // If 1.10
											logs.Error(ai)
											panic(err.Error())
										}
									} else { // If 1.9
										logs.Error(ce)
										panic(err.Error())
									}
								} else { // if 1.8
									logs.Error(contratoGeneral)
									panic(err.Error())
								}
							} else {
								reduccion.ContratoNuevo = nil
							}
						}
					} else { // If 1.2
						logs.Error(vi)
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
				logs.Error(vi)
				panic(err.Error())
			}
		}
		var r models.Resolucion
		url = ResolucionEndpoint + strconv.Itoa(m.IdResolucion)
		if err := GetRequestNew("UrlCrudResoluciones", url, &r); err == nil { // If 2 - resolucion (GET)
			if err := CambiarEstadoResolucion(r.Id, "REXP", m.Usuario); err == nil { // If 2.1 - Cambiar estado
				r.FechaExpedicion = m.FechaExpedicion
				url = ResolucionEndpoint + strconv.Itoa(r.Id)
				if err := SendRequestNew("UrlCrudResoluciones", url, "PUT", &response, &r); err != nil { // If 2.2
					logs.Error(response)
					panic(err.Error())
				}
				if documento, err := AlmacenarResolucionGestorDocumental(r.Id); err != nil {
					logs.Error(response)
					panic(err)
				} else {
					r.NuxeoUid = documento.Enlace
					url = ResolucionEndpoint + strconv.Itoa(r.Id)
					if err := SendRequestNew("UrlCrudResoluciones", url, "PUT", &response, &r); err != nil { // If 2.4
						logs.Error(response)
						panic(err.Error())
					}
				}
				go func() {
					if err := NotificarDocentes(datosCorreo, parametro.CodigoAbreviacion); err != nil {
						logs.Error(err)
					}
				}()
			} else {
				logs.Error(response)
				panic(err.Error())
			}
		} else { // if 2
			logs.Error(r)
			panic(err.Error())
		}
	} else {
		logs.Error(cdve)
		panic(err.Error())
	}
	return
}

func ExpedirCancelacion(m models.ExpedicionCancelacion) (outputError map[string]interface{}) {
	defer func() {
		if err := recover(); err != nil {
			outputError = map[string]interface{}{"funcion": "/ExpedirCancelacion", "err": err, "status": "502"}
			panic(outputError)
		}
	}()

	cancelaciones := m.Vinculaciones

	var response interface{}
	var usuario map[string]interface{}
	var datosCorreo []models.EmailData
	var parametro models.Parametro
	var resv models.ResolucionVinculacionDocente
	var err error

	usuario, err = GetUsuario(cancelaciones[0].ContratoCancelado.Usuario)
	if err != nil {
		logs.Error(err)
		panic(err.Error())
	}

	for _, can := range cancelaciones {
		var v models.VinculacionDocente
		var v2 []models.VinculacionDocente
		url := "vinculacion_docente?query=Id:" + strconv.Itoa(can.VinculacionDocente.Id)
		if err := GetRequestNew("UrlCrudResoluciones", url, &v2); err != nil {
			panic("Vinculacion (cancelacion) -> " + err.Error())
		}
		v = v2[0]
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
							ce := models.ContratoEstado{
								NumeroContrato: contratoCancelado.NumeroContrato,
								Vigencia:       contratoCancelado.Vigencia,
								FechaRegistro:  time.Now(),
								Usuario:        usuario["documento_compuesto"].(string),
								Estado: &models.EstadoContrato{
									Id: 7,
								},
							}
							url = "contrato_estado"
							if err := SendRequestLegacy("UrlcrudAgora", url, "POST", &response, &ce); err == nil { // If 5 - contrato_estado
								var r models.Resolucion
								var dep models.Dependencia
								urlRes := ResolucionEndpoint + strconv.Itoa(m.IdResolucion)
								if err := GetRequestNew("UrlCrudResoluciones", urlRes, &r); err == nil {
									if err := GetRequestNew("UrlCrudResoluciones", ResVinEndpoint+strconv.Itoa(r.Id), &resv); err != nil {
										logs.Error(err)
										panic(err.Error())
									}
									if err := GetRequestLegacy("UrlcrudOikos", "dependencia/"+strconv.Itoa(resv.FacultadId), &dep); err != nil {
										logs.Error(err)
										panic(err.Error())
									}
									url := "parametro/" + strconv.Itoa(r.TipoResolucionId)
									if err := GetRequestNew("UrlcrudParametros", url, &parametro); err != nil {
										logs.Error(err)
										panic("Cargando tipo_resolucion -> " + err.Error())
									}

									datoMensaje := models.EmailData{
										Documento:        strconv.FormatFloat(v.PersonaId, 'f', 0, 64),
										ContratoId:       "",
										Facultad:         dep.Nombre,
										NumeroResolucion: r.NumeroResolucion,
									}
									datosCorreo = append(datosCorreo, datoMensaje)
									if err := CambiarEstadoResolucion(r.Id, "REXP", can.ContratoCancelado.Usuario); err == nil {
										r.FechaExpedicion = m.FechaExpedicion
										if fecha := NormalizarFechaTimezone(r.FechaInicio); fecha != nil {
											r.FechaInicio = fecha
										}
										if err := SendRequestNew("UrlCrudResoluciones", urlRes, "PUT", &response, &r); err == nil {
											if documento, err := AlmacenarResolucionGestorDocumental(r.Id); err == nil {
												r.NuxeoUid = documento.Enlace
												if err := SendRequestNew("UrlCrudResoluciones", urlRes, "PUT", &response, &r); err == nil { // if 10
													if err := ReliquidarContratoCancelado(v, contrato); err != nil {
														panic(err)
													}
												} else {
													logs.Error(r)
													panic(err.Error())
												}
											} else { // If 9
												logs.Error(documento)
												panic(err)
											}
										} else { // If 8
											logs.Error(r)
											panic(err.Error())
										}
									} else { // If 7
										logs.Error(err)
										panic(err.Error())
									}
								} else { // If 6
									logs.Error(r)
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
	go func() {
		if err := NotificarDocentes(datosCorreo, parametro.CodigoAbreviacion); err != nil {
			logs.Error(err)
		}
	}()

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
