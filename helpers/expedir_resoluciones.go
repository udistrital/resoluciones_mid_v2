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

func Expedir2(m models.ExpedicionResolucion) (outputError map[string]interface{}) {

	defer func() {
		if err := recover(); err != nil {
			outputError = map[string]interface{}{"funcion": "/Expedir", "err": err, "status": "502"}
			panic(outputError)
		}
	}()

	var cdve int
	var tipoCon models.TipoContrato
	var proveedor []models.InformacionProveedor
	var response interface{}

	vigencia, _, _ := time.Now().Date()
	v := m.Vinculaciones

	var respuestaPeticion map[string]interface{}
	url := "contrato_general/maximo_dve"
	if err := GetRequestNew("UrlcrudAgora", url, &cdve); err == nil { // If 1 - consecutivo contrato_general
		numeroContratos := cdve
		fmt.Println("numeroContratos:", numeroContratos)
		// for vinculaciones
		for _, vinculacion := range *v { // For vinculaciones
			numeroContratos = numeroContratos + 1
			v := vinculacion.VinculacionDocente
			idVinculacionDocente := strconv.Itoa(v.Id)
			url = "vinculacion_docente/" + idVinculacionDocente
			if err := GetRequestNew("UrlCrudResoluciones", url, &respuestaPeticion); err == nil { // If 1.1 - vinculacion_docente
				contrato := vinculacion.ContratoGeneral
				url = "tipo_contrato/" + strconv.Itoa(contrato.TipoContrato.Id)
				if err := GetRequestNew("UrlcrudAgora", url, &tipoCon); err == nil { // If 1.1.1 - tipoContrato
					var sup models.SupervisorContrato
					acta := vinculacion.ActaInicio
					aux1 := 181
					contrato.VigenciaContrato = vigencia
					contrato.Id = "DVE" + strconv.Itoa(numeroContratos)
					contrato.FormaPago.Id = 240
					contrato.DescripcionFormaPago = "Abono a Cuenta Mensual de acuerdo a puntos y horas laboradas"
					contrato.Justificacion = "Docente de Vinculacion Especial"
					contrato.UnidadEjecucion.Id = 269
					contrato.LugarEjecucion.Id = 4
					contrato.TipoControl = aux1
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
					sup, err := SupervisorActual(v.ResolucionVinculacionDocenteId.Id)
					if err != nil { // If 1.1.2 - supervisorActual
						fmt.Println("Error en If 1.1.2 - supervisorActual!")
						logs.Error(err)
						outputError = map[string]interface{}{"funcion": "/Expedir ", "err": err, "status": "502"}
						return outputError
					}
					contrato.Supervisor = &sup
					contrato.Condiciones = "Sin condiciones"
					url = "informacion_proveedor?query=NumDocumento:" + strconv.Itoa(contrato.Contratista)
					if err := GetRequestNew("UrlcrudAgora", url, &proveedor); err == nil { // If 1.1.3 - informacion_proveedor
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
									"Id":           6,
									"TipoContrato": "Contrato de Prestación de Servicios Profesionales o Apoyo a la Gestión",
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
							fmt.Println(contratoGeneral)
							url = "contrato_general"
							if err := SendRequestNew("UrlcrudAgora", url, "POST", &response, contratoGeneral); err == nil { // If 1.1.5 contrato_general
								aux1 := contrato.Id
								aux2 := contrato.VigenciaContrato
								var ce models.ContratoEstado
								var ec models.EstadoContrato
								ce.NumeroContrato = aux1
								ce.Vigencia = aux2
								ce.FechaRegistro = time.Now()
								ec.Id = 4
								ce.Estado = &ec
								var response2 models.ContratoEstado
								url = "contrato_estado"
								if err := SendRequestNew("UrlcrudAgora", url, "POST", &response2, &ce); err == nil { // If 1.1.6 contrato_estado
									a := vinculacion.VinculacionDocente
									var ai models.ActaInicio
									ai.NumeroContrato = aux1
									ai.Vigencia = aux2
									ai.Descripcion = acta.Descripcion
									ai.FechaInicio = acta.FechaInicio
									ai.FechaFin = acta.FechaFin
									ai.FechaFin = CalcularFechaFin(acta.FechaInicio, a.NumeroSemanas)
									ai.FechaRegistro = time.Now()
									var response3 models.ActaInicio
									url = "acta_inicio"
									if err := SendRequestNew("UrlcrudAgora", url, "POST", &response3, &ai); err == nil { // If 1.1.7 acta_inicio
										var cd models.ContratoDisponibilidad
										cd.NumeroContrato = aux1
										fmt.Println("aux1 ", aux1)
										cd.Vigencia = aux2
										cd.Estado = true
										cd.FechaRegistro = time.Now()
										var dv models.DisponibilidadVinculacion
										url = "disponibilidad_vinculacion?query=VinculacionDocenteId.Id:" + strconv.Itoa(v.Id)
										if err := GetRequestNew("UrlCrudResoluciones", url, &dv); err == nil { // If 1.1.8 - disponibilidad_vinculacion
											cd.NumeroCdp = int(dv.Disponibilidad)
											var response4 models.ContratoDisponibilidad
											url = "contrato_disponibilidad"
											if err := SendRequestNew("UrlcrudAgora", url, "POST", &response4, &cd); err == nil { // If 1.1.9 - contrato_disponibilidad
												a.PuntoSalarialId = vinculacion.VinculacionDocente.PuntoSalarialId
												a.SalarioMinimoId = vinculacion.VinculacionDocente.SalarioMinimoId
												v := a
												v.NumeroContrato = aux1
												v.Vigencia = aux2
												v.FechaInicio = acta.FechaInicio
												url = "vinculacion_docente/" + strconv.Itoa(v.Id)
												if err := SendRequestNew("UrlCrudResoluciones", url, "PUT", &response, &v); err == nil { // If 1.1.10 - vinculacion_docente
													fmt.Println(response)
													fmt.Println("Vinculación docente actualizada")
												} else { // If 1.1.10 - vinculacion_docente
													fmt.Println("Error en If 1.1.10 - vinculacion_docente!")
													logs.Error(err)
													outputError = map[string]interface{}{"funcion": "/Expedir ", "err": err.Error(), "status": "502"}
													return outputError
												}
											} else { // If 1.1.9 -contrato_disponibilidad
												fmt.Println("Error en If 1.1.9 - contrato_disponibilidad!")
												logs.Error(err)
												outputError = map[string]interface{}{"funcion": "/Expedir ", "err": err.Error(), "status": "502"}
												return outputError
											}
										} else { // If 1.1.8 - disponibilidad_vinculacion
											fmt.Println("Error en If 1.1.8 - disponibilidad_vinculacion!")
											logs.Error(err)
											outputError = map[string]interface{}{"funcion": "/Expedir ", "err": err.Error(), "status": "502"}
											return outputError
										}
									} else { // If 1.1.7 acta_inicio
										fmt.Println("Error en If 1.1.7 - acta_inicio!")
										logs.Error(err)
										outputError = map[string]interface{}{"funcion": "/Expedir ", "err": err.Error(), "status": "502"}
										return outputError
									}
								} else { // If 1.1.6
									fmt.Println("Error en If 1.1.6 - contrato_estado!")
									logs.Error(err)
									outputError = map[string]interface{}{"funcion": "/Expedir ", "err": err.Error(), "status": "502"}
									return outputError
								}
							} else { // If 1.1.5
								fmt.Println("Error en If 1.1.5 - contrato_general (POST)!")
								logs.Error(err)
								outputError = map[string]interface{}{"funcion": "/Expedir ", "err": err.Error(), "status": "502"}
								return outputError
							}
						} else { // If 1.1.4 proveedor
							fmt.Println("Error en If 1.1.4 - proveedor vacío!")
						}
					} else { // If 1.1.3
						fmt.Println("Error en If 1.1.3 - informacion_proveedor!")
						logs.Error(err)
						outputError = map[string]interface{}{"funcion": "/Expedir ", "err": err.Error(), "status": "502"}
						return outputError
					}
				} else { // If 1.1.1
					fmt.Println("Error en If 1.1.1 - tipoContrato!")
					logs.Error(err)
					outputError = map[string]interface{}{"funcion": "/Expedir ", "err": err.Error(), "status": "502"}
					return outputError
				}
			} else { // If 1.1
				fmt.Println("Error en If 1.1 - vinculacion_docente!")
				logs.Error(err)
				outputError = map[string]interface{}{"funcion": "/Expedir ", "err": err.Error(), "status": "502"}
				return outputError
			}
		} // For vinculaciones
		var r models.Resolucion
		var rest models.ResolucionEstado
		r.Id = m.IdResolucion
		idResolucionDVE := strconv.Itoa(m.IdResolucion)
		url = "resolucion/" + idResolucionDVE
		if err := GetRequestNew("UrlCrudResoluciones", url, &r); err == nil { // 1.2 - resolucion/
			r.FechaExpedicion = m.FechaExpedicion
			url = "resolucion/" + strconv.Itoa(r.Id)
			if err := SendRequestNew("UrlCrudResoluciones", url, "PUT", &response, &r); err == nil { // If 1.2.1
				var e models.ResolucionEstado
				var er models.EstadoResolucion
				e.ResolucionId = &r
				er.Id = 2
				url = "resolucion_estado?query=Activo:true,ResolucionId.Id:" + strconv.Itoa(r.Id)
				if err := GetRequestNew("UrlCrudResoluciones", url, &rest); err == nil { // If 1.2.2
					e.EstadoResolucionId = rest.EstadoResolucionId
				} else { // If 1.2.2
					fmt.Println("Error en If 1.2.2 - resolucion_estado/ (GET)!")
					logs.Error(err)
					outputError = map[string]interface{}{"funcion": "/Expedir ", "err": err.Error(), "status": "502"}
					return outputError
				}
				url = "resolucion_estado"
				if err := SendRequestNew("UrlCrudResoluciones", url, "POST", &response, &e); err == nil { // If 1.2.3
					fmt.Println("Expedición exitosa")
				} else { // If 1.2.3
					fmt.Println("Error en If 1.2.3 - resolucion_estado/ (POST)!")
					logs.Error(err)
					outputError = map[string]interface{}{"funcion": "/Expedir ", "err": err.Error(), "status": "502"}
					return outputError
				}
			} else { // If 1.2.1
				fmt.Println("Error en If 1.2.1 - resolucion/ (PUT)!")
				logs.Error(err)
				outputError = map[string]interface{}{"funcion": "/Expedir ", "err": err.Error(), "status": "502"}
				return outputError
			}
		} else { // If 1.2
			fmt.Println("Error en If 1.2 - resolucion/!")
			logs.Error(err)
			outputError = map[string]interface{}{"funcion": "/Expedir ", "err": err.Error(), "status": "502"}
			return outputError
		}
	} else { // If 1
		fmt.Println("Error en If 1 - consecutivo contrato_general!")
		logs.Error(err)
		outputError = map[string]interface{}{"funcion": "/Expedir ", "err": err.Error(), "status": "502"}
		return outputError
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

	v := m.Vinculaciones
	beego.Info(v)

	for _, vinculacion := range *v {
		v := vinculacion.VinculacionDocente
		idVinculacionDocente := strconv.Itoa(v.Id)
		url := "vinculacion_docente/" + idVinculacionDocente
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

		// var dispoap []models.DisponibilidadApropiacion
		// url = "disponibilidad_apropiacion?query=Id:" + strconv.Itoa(v.Disponibilidad)
		// if err := GetRequestLegacy("UrlcrudKronos", url, &dispoap); err != nil { // If 4 - disponibilidad_apropiacion
		// 	beego.Error("Error en If 4 - Disponibilidad no válida asociada al docente identificado con " + strconv.Itoa(contrato.Contratista) + " en Ágora")
		// 	logs.Error(dispoap)
		// 	outputError = map[string]interface{}{"funcion": "/ValidarDatosExpedicion3", "err": "No existe el docente con este numero de documento", "status": "502"}
		// 	return outputError
		// }
		// if dispoap == nil {
		// 	beego.Error("Error en If 5 - Disponibilidad no válida asociada al docente identificado con " + strconv.Itoa(contrato.Contratista) + " en Ágora")
		// 	logs.Error(dispoap)
		// 	outputError = map[string]interface{}{"funcion": "/ValidarDatosExpedicion3", "err": "No existe el docente con este numero de documento", "status": "502"}
		// 	return outputError
		// }

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
		beego.Info(proycur)
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
	var modVin []models.ModificacionVinculacion
	var response interface{}
	var resolucion models.Resolucion
	var respuesta_peticion map[string]interface{}
	vigencia, _, _ := time.Now().Date()
	v := m.Vinculaciones

	url := "contrato_general/maximo_dve"
	if err := GetRequestLegacy("UrlcrudAgora", url, &cdve); err == nil { // If 1 - consecutivo contrato_general
		numeroContratos := cdve
		for _, vinculacion := range *v {
			numeroContratos = numeroContratos + 1
			v := vinculacion.VinculacionDocente
			idvinculaciondocente := strconv.Itoa(v.Id)
			url = "vinculacion_docente/" + idvinculaciondocente
			if err := GetRequestNew("UrlCrudResoluciones", url, &v); err == nil { // If 1.1 - vinculacion_docente
				contrato := vinculacion.ContratoGeneral
				var sup models.SupervisorContrato
				acta := vinculacion.ActaInicio
				aux1 := 181
				contrato.VigenciaContrato = vigencia
				contrato.Id = "DVE" + strconv.Itoa(numeroContratos)
				contrato.FormaPago.Id = 240
				contrato.DescripcionFormaPago = "Abono a Cuenta Mensual de acuerdo a puntos y horas laboradas"
				contrato.Justificacion = "Docente de Vinculacion Especial"
				contrato.UnidadEjecucion.Id = 269
				contrato.LugarEjecucion.Id = 4
				contrato.TipoControl = aux1
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

				sup, err := SupervisorActual(v.ResolucionVinculacionDocenteId.Id)
				if err != nil {
					fmt.Println("Error en SupervisorActual")
					return err
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
								"Id":           6,
								"TipoContrato": "Contrato de Prestación de Servicios Profesionales o Apoyo a la Gestión",
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
						url = "modificacion_vinculacion?query=VinculacionDocenteRegistradaId.Id:" + strconv.Itoa(v.Id)
						if err := GetRequestNew("UrlCrudResoluciones", url, &modVin); err == nil { // If 1.4 - modificacion_vinculacion
							var actaInicioAnterior []models.ActaInicio
							vinculacionModificacion := modVin[0].VinculacionDocenteRegistradaId
							vinculacionOriginal := modVin[0].VinculacionDocenteCanceladaId
							url = "resolucion/" + strconv.Itoa(v.ResolucionVinculacionDocenteId.Id)
							if err := GetRequestNew("UrlCrudResoluciones", url, &resolucion); err == nil { // If 1.5 - resolucion
							} else {
								fmt.Println("Error en If 1.5 - resolucion! ", err)
								logs.Error(resolucion)
								outputError = map[string]interface{}{"funcion": "/ExpedirModificacion", "err": err.Error(), "status": "502"}
								return outputError
							}
							url = "acta_inicio?query=NumeroContrato:" + modVin[0].VinculacionDocenteCanceladaId.NumeroContrato + ",Vigencia:" + strconv.Itoa(int(modVin[0].VinculacionDocenteCanceladaId.Vigencia))
							if err := GetRequestLegacy("UrlcrudAgora", url, &actaInicioAnterior); err == nil { // If 1.6 - acta_inicio
								semanasIniciales := vinculacionOriginal.NumeroSemanas
								semanasModificar := vinculacionModificacion.NumeroSemanas
								horasIniciales := vinculacionOriginal.NumeroHorasSemanales
								fechaFinNuevoContrato := CalcularFechaFin(acta.FechaInicio, semanasModificar)
								horasTotales := horasIniciales + vinculacionModificacion.NumeroHorasSemanales
								// Sólo si es reducción cambia la fecha fin del acta anterior y el valor del nuevo contrato
								if resolucion.TipoResolucionId == 4 {
									var aini models.ActaInicio
									aini.Id = actaInicioAnterior[0].Id
									aini.NumeroContrato = actaInicioAnterior[0].NumeroContrato
									aini.Vigencia = actaInicioAnterior[0].Vigencia
									aini.Descripcion = actaInicioAnterior[0].Descripcion
									aini.FechaInicio = actaInicioAnterior[0].FechaInicio
									aini.FechaFin = acta.FechaInicio
									fechaFinNuevoContrato = actaInicioAnterior[0].FechaFin
									beego.Info("fin nuevo ", fechaFinNuevoContrato)
									beego.Info("fin viejo", aini.FechaFin)
									url = "acta_inicio/" + strconv.Itoa(aini.Id)
									if err := SendRequestLegacy("UrlcrudAgora", url, "PUT", &response, &aini); err == nil { // If 1.7 - acta_inicio (PUT)
										fmt.Println("Acta anterior cancelada en la fecha indicada")
									} else {
										fmt.Println("Error en If 1.7 - acta_inicio (PUT)! ", err)
										logs.Error(response)
										outputError = map[string]interface{}{"funcion": "/ExpedirModificacion", "err": err.Error(), "status": "502"}
										return outputError
									}
									// Calcula el valor del nuevo contrato con base en las semanas desde la fecha inicio escogida hasta la nueva fecha fin y las nuevas horas
									semanasTranscurridasDecimal := (acta.FechaInicio.Sub(actaInicioAnterior[0].FechaInicio).Hours()) / 24 / 30 * 4 // cálculo con base en meses de 30 días y 4 semanas
									semanasTranscurridas, decimal := math.Modf(semanasTranscurridasDecimal)
									if decimal > 0 {
										semanasTranscurridas = semanasTranscurridas + 1
									}
									var semanasTranscurridasInt = int(semanasTranscurridas)
									semanasRestantes := semanasIniciales - semanasTranscurridasInt - semanasModificar
									horasTotales = horasIniciales - vinculacionModificacion.NumeroHorasSemanales
									var vinc [1]models.VinculacionDocente
									vinc[0] = models.VinculacionDocente{
										ResolucionVinculacionDocenteId: &models.ResolucionVinculacionDocente{Id: m.IdResolucion},
										PersonaId:                      v.PersonaId,
										NumeroHorasSemanales:           horasTotales,
										NumeroSemanas:                  semanasModificar,
										DedicacionId:                   v.DedicacionId,
										ProyectoCurricularId:           v.ProyectoCurricularId,
										Categoria:                      v.Categoria,
										Vigencia:                       v.Vigencia,
									}
									salario, err := CalcularValorContratoReduccion(vinc, semanasRestantes, horasIniciales, v.ResolucionVinculacionDocenteId.NivelAcademico, resolucion.Periodo)
									if err != nil {
										fmt.Println("Error en cálculo del contrato reducción!", err)
										outputError = map[string]interface{}{"funcion": "/ExpedirModificacion17", "err": err, "status": "502"}
										return outputError
									}
									// Si es de posgrado calcula el valor que se le ha pagado hasta la fecha de inicio y se resta del total que debe quedar con la reducción
									if v.ResolucionVinculacionDocenteId.NivelAcademico == "POSGRADO" {
										diasOriginales, _ := math.Modf((actaInicioAnterior[0].FechaFin.Sub(actaInicioAnterior[0].FechaInicio).Hours()) / 24)
										diasTranscurridos, _ := math.Modf((acta.FechaInicio.Sub(actaInicioAnterior[0].FechaInicio).Hours()) / 24)
										valorDiario := vinculacionOriginal.ValorContrato / diasOriginales
										valorPagado := valorDiario * diasTranscurridos
										salario = salario - valorPagado
									}
									contrato.ValorContrato = salario
									beego.Info(contrato.ValorContrato)
								}
								if contrato.ValorContrato > 0 {
									url = "contrato_general"
									if err := SendRequestLegacy("UrlcrudAgora", url, "POST", &response, &contratoGeneral); err == nil { // If 1.8 - contrato_general (POST)
										aux1 := contrato.Id
										aux2 := contrato.VigenciaContrato
										var ce models.ContratoEstado
										var ec models.EstadoContrato
										ce.NumeroContrato = aux1
										ce.Vigencia = aux2
										ce.FechaRegistro = time.Now()
										ec.Id = 4
										ce.Estado = &ec
										url = "contrato_estado"
										if err := SendRequestLegacy("UrlcrudAgora", url, "POST", &response, &ce); err == nil { // If 1.9 - contrato_estado (POST)
											var ai models.ActaInicio
											ai.NumeroContrato = aux1
											ai.Vigencia = aux2
											ai.Descripcion = acta.Descripcion
											ai.FechaInicio = acta.FechaInicio
											ai.FechaFin = fechaFinNuevoContrato
											beego.Info("inicio ", ai.FechaInicio, " fin ", ai.FechaFin)
											url = "acta_inicio"
											if err := SendRequestLegacy("UrlcrudAgora", url, "POST", &response, &ai); err == nil { // If 1.10 - acta_inicio (POST)
												var cd models.ContratoDisponibilidad
												cd.NumeroContrato = aux1
												cd.Vigencia = aux2
												cd.Estado = true
												cd.FechaRegistro = time.Now()
												var dv models.DisponibilidadVinculacion
												url = "disponibilidad_vinculacion?query=VinculacionDocenteId.Id:" + strconv.Itoa(v.Id)
												if err := GetRequestNew("UrlCrudResoluciones", url, &dv); err == nil { // If 1.11 - DisponibilidadVinculacion
													cd.NumeroCdp = int(dv.Disponibilidad)
													url = "contrato_disponibilidad"
													if err := SendRequestNew("UrlcrudAgora", url, "POST", &response, &cd); err == nil { // If 1.12 - contrato_disponibilidad
														vinculacionModificacion.PuntoSalarialId = vinculacion.VinculacionDocente.PuntoSalarialId
														vinculacionModificacion.SalarioMinimoId = vinculacion.VinculacionDocente.SalarioMinimoId
														vinculacionModificacion.NumeroContrato = aux1
														vinculacionModificacion.Vigencia = aux2
														url = "vinculacion_docente/" + strconv.Itoa(vinculacionModificacion.Id)
														if err := SendRequestNew("UrlCrudResoluciones", url, "PUT", &response, &vinculacionModificacion); err == nil {
															fmt.Println(response)
															fmt.Println("Vinculación docente actualizada")
														} else {
															fmt.Println("Error en If 1.13 - vinculacion_docente! ", err)
															logs.Error(response)
															outputError = map[string]interface{}{"funcion": "/ExpedirModificacion", "err": err.Error(), "status": "502"}
															return outputError
														}
													} else {
														fmt.Println("Error en If 1.12 - contrato_disponibilidad! ", err)
														logs.Error(cd)
														outputError = map[string]interface{}{"funcion": "/ExpedirModificacion", "err": err.Error(), "status": "502"}
														return outputError
													}
												} else {
													fmt.Println("Error en If 1.11 - DisponibilidadVinculacion! ", err)
													logs.Error(dv)
													outputError = map[string]interface{}{"funcion": "/ExpedirModificacion", "err": err.Error(), "status": "502"}
													return outputError
												}
											} else { // If 1.10
												fmt.Println("Error en If 1.10 - acta_inicio (POST)! ", err)
												logs.Error(ai)
												outputError = map[string]interface{}{"funcion": "/ExpedirModificacion", "err": err.Error(), "status": "502"}
												return outputError
											}
										} else { // If 1.9
											fmt.Println("Error en If 1.9 - contrato_estado (POST)! ", err)
											logs.Error(ce)
											outputError = map[string]interface{}{"funcion": "/ExpedirModificacion", "err": err.Error(), "status": "502"}
											return outputError
										}
									} else { // if 1.8
										fmt.Println("Error en If 1.8 - contrato_general (POST)! ", err)
										logs.Error(contratoGeneral)
										outputError = map[string]interface{}{"funcion": "/ExpedirModificacion", "err": err.Error(), "status": "502"}
										return outputError
									}
								}
							} else {
								fmt.Println("Error en If 1.6 - acta_inicio! ", err)
								logs.Error(actaInicioAnterior)
								outputError = map[string]interface{}{"funcion": "/ExpedirModificacion", "err": err.Error(), "status": "502"}
								return outputError
							}
						} else { // If 1.4
							fmt.Println("Error en If 1.4 - modificacion_vinculacion! ", err)
							logs.Error(modVin)
							outputError = map[string]interface{}{"funcion": "/ExpedirModificacion", "err": err.Error(), "status": "502"}
							return outputError
						}
					} // If 1.3
				} else { // If 1.2
					fmt.Println("Error en If 1.2 - informacion_proveedor! ", err)
					logs.Error(v)
					outputError = map[string]interface{}{"funcion": "/ExpedirModificacion", "err": err.Error(), "status": "502"}
					return outputError
				}
			} else { // If 1.1
				fmt.Println("Error en If 1.1 - vinculacion_docente! ", err)
				logs.Error(v)
				outputError = map[string]interface{}{"funcion": "/ExpedirModificacion", "err": err.Error(), "status": "502"}
				return outputError
			}
		}
		var r models.Resolucion
		r.Id = m.IdResolucion
		idResolucionDVE := strconv.Itoa(m.IdResolucion)
		url = "resolucion/" + idResolucionDVE
		if err := GetRequestNew("UrlCrudResoluciones", url, &r); err == nil { // If 2 - resolucion (GET)
			r.FechaExpedicion = m.FechaExpedicion
			url = "resolucion/" + strconv.Itoa(r.Id)
			if err := SendRequestNew("UrlCrudResoluciones", url, "PUT", &respuesta_peticion, &r); err == nil { // If 2.1 - resolucion (PUT)
				var e models.ResolucionEstado
				var er models.EstadoResolucion
				e.ResolucionId = &r
				er.Id = 2
				e.EstadoResolucionId = er.Id
				url = "resolucion_estado"
				if err := SendRequestNew("UrlCrudResoluciones", url, "POST", &response, &e); err == nil { // If 2.2 - resolucion_estado (POST)
					fmt.Println("Expedición exitosa")
				} else {
					fmt.Println("Error en If 2.2 - resolucion_estado (POST) ! ", err)
					logs.Error(response)
					outputError = map[string]interface{}{"funcion": "/ExpedirModificacion", "err": err.Error(), "status": "502"}
					return outputError
				}
			} else { // If 2.1
				fmt.Println("Error en If 2.1 - resolucion (PUT) ! ", err)
				logs.Error(respuesta_peticion)
				outputError = map[string]interface{}{"funcion": "/ExpedirModificacion", "err": err.Error(), "status": "502"}
				return outputError
			}
		} else { // if 2
			fmt.Println("Error en If 2 - resolucion (GET) ! ", err)
			logs.Error(r)
			outputError = map[string]interface{}{"funcion": "/ExpedirModificacion", "err": err.Error(), "status": "502"}
			return outputError
		}
	} else {
		fmt.Println("Error en If 1 - Consecutivo contrato_general! ", err)
		logs.Error(cdve)
		outputError = map[string]interface{}{"funcion": "/ExpedirModificacion", "err": err.Error(), "status": "502"}
		return outputError
	}
	return
}

func Cancelar(m models.ExpedicionCancelacion) (outputError map[string]interface{}) {

	defer func() {
		if err := recover(); err != nil {
			outputError = map[string]interface{}{"funcion": "/Cancelar", "err": err, "status": "502"}
			panic(outputError)
		}
	}()

	v := m.Vinculaciones
	var contratoCancelado models.ContratoCancelado
	var response interface{}

	for _, vinculacion := range *v {
		v := vinculacion.VinculacionDocente
		idVinculacionDocente := strconv.Itoa(v.Id)
		url := "vinculacion_docente/" + idVinculacionDocente
		if err := GetRequestNew("UrlCrudResoluciones", url, &v); err == nil { // If 1 - vinculacion_docente
			contratoCancelado.NumeroContrato = v.NumeroContrato
			contratoCancelado.Vigencia = int(v.Vigencia)
			contratoCancelado.FechaCancelacion = vinculacion.ContratoCancelado.FechaCancelacion
			contratoCancelado.MotivoCancelacion = vinculacion.ContratoCancelado.MotivoCancelacion
			contratoCancelado.Usuario = vinculacion.ContratoCancelado.Usuario
			contratoCancelado.FechaRegistro = time.Now()
			contratoCancelado.Estado = vinculacion.ContratoCancelado.Estado
			url = "contrato_cancelado"
			if err := SendRequestLegacy("UrlcrudAgora", url, "POST", &response, &contratoCancelado); err == nil { // If 2 -contrato_cancelado (post)
				var ai []models.ActaInicio
				url = "acta_inicio?query=NumeroContrato:" + contratoCancelado.NumeroContrato + ",Vigencia:" + strconv.Itoa(contratoCancelado.Vigencia)
				if err := GetRequestLegacy("UrlcrudAgora", url, &ai); err == nil { // If 3 - acta_inicio (get)
					ai[0].FechaFin = CalcularFechaFin(ai[0].FechaInicio, v.NumeroSemanas)
					url = "acta_inicio/" + strconv.Itoa(ai[0].Id)
					if err := SendRequestLegacy("UrlcrudAgora", url, "PUT", &response, &ai[0]); err == nil { // If 4 - acta_inicio
						var ce models.ContratoEstado
						var ec models.EstadoContrato
						ce.NumeroContrato = contratoCancelado.NumeroContrato
						ce.Vigencia = contratoCancelado.Vigencia
						ce.FechaRegistro = time.Now()
						ec.Id = 7
						ce.Estado = &ec
						url = "contrato_estado"
						if err := SendRequestLegacy("UrlcrudAgora", url, "POST", &response, &ce); err == nil { // If 5 - contrato_estado

						} else { // If 5
							fmt.Println("Error en if 4 - contrato_estado (post)!", err)
							logs.Error(contratoCancelado)
							outputError = map[string]interface{}{"funcion": "/Cancelar", "err": err.Error(), "status": "502"}
							return outputError
						}
					} else { // If 4
						fmt.Println("Error en if 4 - acta_inicio (put)!", err)
						logs.Error(contratoCancelado)
						outputError = map[string]interface{}{"funcion": "/Cancelar", "err": err.Error(), "status": "502"}
						return outputError
					}
				} else { // If 3
					fmt.Println("Error en if 3 - acta_inicio (get)!", err)
					logs.Error(contratoCancelado)
					outputError = map[string]interface{}{"funcion": "/Cancelar", "err": err.Error(), "status": "502"}
					return outputError
				}
			} else { // if 2
				fmt.Println("Error en if 2 - contrato_cancelado (post)!", err)
				logs.Error(contratoCancelado)
				outputError = map[string]interface{}{"funcion": "/Cancelar", "err": err.Error(), "status": "502"}
				return outputError
			}
		} else { // If 1
			fmt.Println("Error en if 1 - vinculacion_docente (get)!", err)
			logs.Error(contratoCancelado)
			outputError = map[string]interface{}{"funcion": "/Cancelar", "err": err.Error(), "status": "502"}
			return outputError
		}
	}

	return
}
