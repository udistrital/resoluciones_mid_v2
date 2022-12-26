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

	// Cambiar en el futuro por terceros_mid/tipo/ordenadoresGasto y filtrar por dependenciaId
	url := "ordenador_gasto?query=DependenciaId:" + strconv.Itoa(r.DependenciaId)
	if err := GetRequestLegacy("UrlcrudCore", url, &ordenadoresGasto); err != nil {
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

	url = "contrato_general/maximo_dve"
	if err := GetRequestLegacy("UrlcrudAgora", url, &cdve); err == nil { // If 1 - consecutivo contrato_general
		numeroContratos := cdve
		fmt.Println("numeroContratos:", numeroContratos)
		// for vinculaciones
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
					sup, err := SupervisorActual(v.ResolucionVinculacionDocenteId.Id)
					if err != nil { // If 1.1.2 - supervisorActual
						fmt.Println("Error en If 1.1.2 - supervisorActual!")
						logs.Error(err)
						panic(err)
					}
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
									ai.FechaFin = CalcularFechaFin(acta.FechaInicio, v.NumeroSemanas)
									ai.FechaRegistro = time.Now()
									ai.Usuario = usuario["documento_compuesto"].(string)
									var response3 models.ActaInicio
									url = "acta_inicio"
									if err := SendRequestLegacy("UrlcrudAgora", url, "POST", &response3, &ai); err == nil { // If 1.1.7 acta_inicio
										var cd models.ContratoDisponibilidad
										cd.NumeroContrato = aux1
										fmt.Println("aux1 ", aux1)
										cd.Vigencia = aux2
										cd.Estado = true
										cd.FechaRegistro = time.Now()
										var dv models.DisponibilidadVinculacion
										var disp []models.DisponibilidadVinculacion
										url = "disponibilidad_vinculacion?query=VinculacionDocenteId.Id:" + strconv.Itoa(v.Id)
										if err := GetRequestNew("UrlCrudResoluciones", url, &disp); err == nil { // If 1.1.8 - disponibilidad_vinculacion
											dv = disp[0]
											cd.NumeroCdp = int(dv.Disponibilidad)
											var response4 models.ContratoDisponibilidad
											url = "contrato_disponibilidad"
											if err := SendRequestLegacy("UrlcrudAgora", url, "POST", &response4, &cd); err == nil { // If 1.1.9 - contrato_disponibilidad
												v.NumeroContrato = &aux1
												v.Vigencia = aux2
												v.FechaInicio = acta.FechaInicio
												url = VinculacionEndpoint + strconv.Itoa(v.Id)
												if err := SendRequestNew("UrlCrudResoluciones", url, "PUT", &response, &v); err != nil { // If 1.1.10 - vinculacion_docente
													fmt.Println("Error en If 1.1.10 - vinculacion_docente!")
													logs.Error(err)
													panic(err.Error())
												}
											} else { // If 1.1.9 -contrato_disponibilidad
												fmt.Println("Error en If 1.1.9 - contrato_disponibilidad!")
												logs.Error(err)
												panic(err.Error())
											}
										} else { // If 1.1.8 - disponibilidad_vinculacion
											fmt.Println("Error en If 1.1.8 - disponibilidad_vinculacion!")
											logs.Error(err)
											panic(err.Error())
										}
									} else { // If 1.1.7 acta_inicio
										fmt.Println("Error en If 1.1.7 - acta_inicio!")
										logs.Error(err)
										panic(err.Error())
									}
								} else { // If 1.1.6
									fmt.Println("Error en If 1.1.6 - contrato_estado!")
									logs.Error(err)
									panic(err.Error())
								}
							} else { // If 1.1.5
								fmt.Println("Error en If 1.1.5 - contrato_general (POST)!")
								logs.Error(response)
								panic(err.Error())
							}
						} else { // If 1.1.4 proveedor
							fmt.Println("Error en If 1.1.4 - proveedor vacío!")
						}
					} else { // If 1.1.3
						fmt.Println("Error en If 1.1.3 - informacion_proveedor!")
						logs.Error(proveedor)
						panic(err.Error())
					}
				} else { // If 1.1.1
					fmt.Println("Error en If 1.1.1 - tipoContrato!")
					logs.Error(tipoCon)
					panic(err.Error())
				}
			} else { // If 1.1
				fmt.Println("Error en If 1.1 - vinculacion_docente!")
				logs.Error(v)
				panic(err.Error())
			}
		} // For vinculaciones
		var r models.Resolucion
		url = ResolucionEndpoint + strconv.Itoa(m.IdResolucion)
		if err := GetRequestNew("UrlCrudResoluciones", url, &r); err == nil { // 1.2 - resolucion/
			if err := CambiarEstadoResolucion(r.Id, "REXP", m.Usuario); err == nil {
				r.FechaExpedicion = m.FechaExpedicion
				url = ResolucionEndpoint + strconv.Itoa(r.Id)
				if err := SendRequestNew("UrlCrudResoluciones", url, "PUT", &response, &r); err != nil { // If 1.2.1
					fmt.Println("Error en If 1.2.1 - resolucion/ (PUT)!")
					logs.Error(response)
					panic(err.Error())
				}
				if documento, err := AlmacenarResolucionGestorDocumental(r.Id); err != nil {
					fmt.Println("Error en If 1.2.3 - AlmacenarResolucionGestorDocumental/ (POST)!")
					logs.Error(response)
					panic(err)
				} else {
					r.NuxeoUid = documento.Enlace
					url = ResolucionEndpoint + strconv.Itoa(r.Id)
					if err := SendRequestNew("UrlCrudResoluciones", url, "PUT", &response, &r); err != nil { // If 1.2.1
						fmt.Println("Error en If 1.2.1 - resolucion/ (PUT)!")
						logs.Error(response)
						panic(err.Error())
					}
				}
			} else { // If 1.2.1
				fmt.Println("Error en If 1.2.1 - Cambiar estado/ (POST)!")
				logs.Error(response)
				panic(err.Error())
			}
		} else { // If 1.2
			fmt.Println("Error en If 1.2 - resolucion/!")
			logs.Error(r)
			panic(err.Error())
		}
	} else { // If 1
		fmt.Println("Error en If 1 - consecutivo contrato_general!")
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
	var modVin []models.ModificacionVinculacion
	var response interface{}
	var resolucion models.Resolucion
	var ordenadoresGasto []models.OrdenadorGasto
	var ordenadorGasto models.OrdenadorGasto
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
		panic(err.Error())
	}

	// Cambiar en el futuro por terceros_mid/tipo/ordenadoresGasto y filtrar por dependenciaId
	url := "ordenador_gasto?query=DependenciaId:" + strconv.Itoa(resolucion.DependenciaId)
	if err := GetRequestLegacy("UrlcrudCore", url, &ordenadoresGasto); err != nil {
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

	url = "contrato_general/maximo_dve"
	if err := GetRequestLegacy("UrlcrudAgora", url, &cdve); err == nil { // If 1 - consecutivo contrato_general
		numeroContratos := cdve
		for _, vinculacion := range vinc {
			numeroContratos = numeroContratos + 1
			v := vinculacion.VinculacionDocente
			url = VinculacionEndpoint + strconv.Itoa(v.Id)
			if err := GetRequestNew("UrlCrudResoluciones", url, &v); err == nil { // If 1.1 - vinculacion_docente
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
				contrato.ValorContrato = v.ValorContrato
				contrato.OrdenadorGasto = ordenadorGasto.Id
				sup, err := SupervisorActual(v.ResolucionVinculacionDocenteId.Id)
				if err != nil {
					fmt.Println("Error en If 1.1.2 - supervisorActual!")
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
						valoresAntes := make(map[string]float64)
						if err := CalcularTrazabilidad(strconv.Itoa(v.Id), &valoresAntes); err != nil {
							fmt.Println("Error en If 1.5 - CalcularTrazabilidad! ", err)
							panic(err.Error())
						}
						url = "modificacion_vinculacion?query=VinculacionDocenteRegistradaId.Id:" + strconv.Itoa(v.Id)
						if err := GetRequestNew("UrlCrudResoluciones", url, &modVin); err == nil { // If 1.4 - modificacion_vinculacion
							var respActaInicioAnterior []models.ActaInicio
							var actaInicioAnterior models.ActaInicio
							vinculacionModificacion := modVin[0].VinculacionDocenteRegistradaId
							vinculacionOriginal := modVin[0].VinculacionDocenteCanceladaId
							url = ResolucionEndpoint + strconv.Itoa(v.ResolucionVinculacionDocenteId.Id)
							if err := GetRequestNew("UrlCrudResoluciones", url, &resolucion); err == nil { // If 1.5 - resolucion
							} else {
								fmt.Println("Error en If 1.5 - resolucion! ", err)
								logs.Error(resolucion)
								panic(err.Error())
							}
							url = "acta_inicio?query=NumeroContrato:" + *vinculacionOriginal.NumeroContrato + ",Vigencia:" + strconv.Itoa(vinculacionOriginal.Vigencia)
							if err := GetRequestLegacy("UrlcrudAgora", url, &respActaInicioAnterior); err == nil { // If 1.6 - acta_inicio
								actaInicioAnterior = respActaInicioAnterior[0]
								semanasIniciales := vinculacionOriginal.NumeroSemanas
								semanasModificar := vinculacionModificacion.NumeroSemanas
								horasIniciales := int(valoresAntes["NumeroHorasSemanales"])
								fechaFinNuevoContrato := CalcularFechaFin(vinculacionModificacion.FechaInicio, semanasModificar)
								horasTotales := horasIniciales + vinculacionModificacion.NumeroHorasSemanales
								// Sólo si es reducción cambia la fecha fin del acta anterior y el valor del nuevo contrato
								var tipoRes models.Parametro
								url2 := ParametroEndpoint + strconv.Itoa(resolucion.TipoResolucionId)
								if err := GetRequestNew("UrlcrudParametros", url2, &tipoRes); err != nil {
									logs.Error(err)
									panic(err.Error())
								}
								if tipoRes.CodigoAbreviacion == "RRED" {
									var aini models.ActaInicio
									aini.Id = actaInicioAnterior.Id
									aini.NumeroContrato = actaInicioAnterior.NumeroContrato
									aini.Vigencia = actaInicioAnterior.Vigencia
									aini.Descripcion = actaInicioAnterior.Descripcion
									aini.FechaInicio = actaInicioAnterior.FechaInicio
									aini.FechaFin = vinculacionModificacion.FechaInicio
									aini.Usuario = usuario["documento_compuesto"].(string)
									fechaFinNuevoContrato = actaInicioAnterior.FechaFin
									beego.Info("fin nuevo ", fechaFinNuevoContrato)
									beego.Info("fin viejo", aini.FechaFin)
									url = "acta_inicio/" + strconv.Itoa(aini.Id)
									if err := SendRequestLegacy("UrlcrudAgora", url, "PUT", &response, &aini); err != nil { // If 1.7 - acta_inicio (PUT)
										fmt.Println("Error en If 1.7 - acta_inicio (PUT)! ", err)
										logs.Error(response)
										panic(err.Error())
									}
									// Calcula el valor del nuevo contrato con base en las semanas desde la fecha inicio escogida hasta la nueva fecha fin y las nuevas horas
									semanasTranscurridasDecimal := (vinculacionModificacion.FechaInicio.Sub(actaInicioAnterior.FechaInicio).Hours()) / 24 / 30 * 4 // cálculo con base en meses de 30 días y 4 semanas
									semanasTranscurridas, decimal := math.Modf(semanasTranscurridasDecimal)
									if decimal > 0 {
										semanasTranscurridas = semanasTranscurridas + 1
									}
									var semanasTranscurridasInt = int(semanasTranscurridas)
									semanasRestantes := semanasIniciales - semanasTranscurridasInt - semanasModificar
									horasTotales = horasIniciales - vinculacionModificacion.NumeroHorasSemanales
									var vinc [1]models.VinculacionDocente
									vinc[0] = models.VinculacionDocente{
										ResolucionVinculacionDocenteId: vinculacionModificacion.ResolucionVinculacionDocenteId,
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
										panic(err)
									}
									// Si es de posgrado calcula el valor que se le ha pagado hasta la fecha de inicio y se resta del total que debe quedar con la reducción
									if v.ResolucionVinculacionDocenteId.NivelAcademico == "POSGRADO" {
										diasOriginales, _ := math.Modf((actaInicioAnterior.FechaFin.Sub(actaInicioAnterior.FechaInicio).Hours()) / 24)
										diasTranscurridos, _ := math.Modf((vinculacionModificacion.FechaInicio.Sub(actaInicioAnterior.FechaInicio).Hours()) / 24)
										valorDiario := vinculacionOriginal.ValorContrato / diasOriginales
										valorPagado := valorDiario * diasTranscurridos
										salario = salario - valorPagado
									}
									contrato.ValorContrato = salario
									beego.Info(contrato.ValorContrato)
								}
								if contrato.ValorContrato > 0 {
									if err := SendRequestLegacy("UrlcrudAgora", "contrato_general", "POST", &response, &contratoGeneral); err == nil { // If 1.8 - contrato_general (POST)
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
										if err := SendRequestLegacy("UrlcrudAgora", "contrato_estado", "POST", &response, &ce); err == nil { // If 1.9 - contrato_estado (POST)
											var ai models.ActaInicio
											ai.NumeroContrato = aux1
											ai.Vigencia = aux2
											ai.Descripcion = acta.Descripcion
											ai.FechaInicio = vinculacionModificacion.FechaInicio
											ai.FechaFin = fechaFinNuevoContrato
											ai.Usuario = usuario["documento_compuesto"].(string)
											if err := SendRequestLegacy("UrlcrudAgora", "acta_inicio", "POST", &response, &ai); err == nil { // If 1.10 - acta_inicio (POST)
												var cd models.ContratoDisponibilidad
												cd.NumeroContrato = aux1
												cd.Vigencia = aux2
												cd.Estado = true
												cd.FechaRegistro = time.Now()
												var dv models.DisponibilidadVinculacion
												var disp []models.DisponibilidadVinculacion
												url = "disponibilidad_vinculacion?query=VinculacionDocenteId.Id:" + strconv.Itoa(v.Id)
												if err := GetRequestNew("UrlCrudResoluciones", url, &disp); err == nil { // If 1.11 - DisponibilidadVinculacion
													dv = disp[0]
													cd.NumeroCdp = int(dv.Disponibilidad)
													if err := SendRequestLegacy("UrlcrudAgora", "contrato_disponibilidad", "POST", &response, &cd); err == nil { // If 1.12 - contrato_disponibilidad
														vinculacionModificacion.NumeroContrato = &aux1
														vinculacionModificacion.Vigencia = aux2
														url = VinculacionEndpoint + strconv.Itoa(vinculacionModificacion.Id)
														if err := SendRequestNew("UrlCrudResoluciones", url, "PUT", &response, &vinculacionModificacion); err != nil {
															fmt.Println("Error en If 1.13 - vinculacion_docente! ", err)
															logs.Error(response)
															panic(err.Error())
														}
													} else {
														fmt.Println("Error en If 1.12 - contrato_disponibilidad! ", err)
														logs.Error(cd)
														panic(err.Error())
													}
												} else {
													fmt.Println("Error en If 1.11 - DisponibilidadVinculacion! ", err)
													logs.Error(dv)
													panic(err.Error())
												}
											} else { // If 1.10
												fmt.Println("Error en If 1.10 - acta_inicio (POST)! ", err)
												logs.Error(ai)
												panic(err.Error())
											}
										} else { // If 1.9
											fmt.Println("Error en If 1.9 - contrato_estado (POST)! ", err)
											logs.Error(ce)
											panic(err.Error())
										}
									} else { // if 1.8
										fmt.Println("Error en If 1.8 - contrato_general (POST)! ", err)
										logs.Error(contratoGeneral)
										panic(err.Error())
									}
								}
							} else {
								fmt.Println("Error en If 1.6 - acta_inicio! ", err)
								logs.Error(actaInicioAnterior)
								panic(err.Error())
							}
						} else { // If 1.4
							fmt.Println("Error en If 1.4 - modificacion_vinculacion! ", err)
							logs.Error(modVin)
							panic(err.Error())
						}
					} // If 1.3
				} else { // If 1.2
					fmt.Println("Error en If 1.2 - informacion_proveedor! ", err)
					logs.Error(v)
					panic(err.Error())
				}
			} else { // If 1.1
				fmt.Println("Error en If 1.1 - vinculacion_docente! ", err)
				logs.Error(v)
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
					fmt.Println("Error en If 2.2 - resolucion/ (PUT)!")
					logs.Error(response)
					panic(err.Error())
				}
				if documento, err := AlmacenarResolucionGestorDocumental(r.Id); err != nil {
					fmt.Println("Error en If 2.3 - AlmacenarResolucionGestorDocumental/ (POST)!")
					logs.Error(response)
					panic(err)
				} else {
					r.NuxeoUid = documento.Enlace
					url = ResolucionEndpoint + strconv.Itoa(r.Id)
					if err := SendRequestNew("UrlCrudResoluciones", url, "PUT", &response, &r); err != nil { // If 2.4
						fmt.Println("Error en If 2.4 - resolucion/ (PUT)!")
						logs.Error(response)
						panic(err.Error())
					}
				}
			} else {
				fmt.Println("Error en If 2.1 - Cambiar estado/ (POST)!")
				logs.Error(response)
				panic(err.Error())
			}
		} else { // if 2
			fmt.Println("Error en If 2 - resolucion (GET) ! ", err)
			logs.Error(r)
			panic(err.Error())
		}
	} else {
		fmt.Println("Error en If 1 - Consecutivo contrato_general! ", err)
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

	vin := m.Vinculaciones

	var response interface{}
	var usuario map[string]interface{}
	var err error

	usuario, err = GetUsuario(vin[0].ContratoCancelado.Usuario)
	if err != nil {
		logs.Error(err)
		panic(err.Error())
	}

	for _, vinculacion := range vin {
		v := vinculacion.VinculacionDocente
		contratos := new([]models.ContratoCancelar)
		if err := BuscarContratosCancelar(v.Id, contratos); err == nil { // If 1 - vinculacion_docente
			for _, contrato := range *contratos {
				contratoCancelado := &models.ContratoCancelado{
					NumeroContrato:    contrato.NumeroContrato,
					Vigencia:          contrato.Vigencia,
					FechaCancelacion:  vinculacion.ContratoCancelado.FechaCancelacion,
					MotivoCancelacion: vinculacion.ContratoCancelado.MotivoCancelacion,
					Usuario:           usuario["documento_compuesto"].(string),
					FechaRegistro:     time.Now(),
					Estado:            vinculacion.ContratoCancelado.Estado,
				}
				url := "contrato_cancelado"
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
							ce.Usuario = usuario["documento_compuesto"].(string)
							ec.Id = 7
							ce.Estado = &ec
							url = "contrato_estado"
							if err := SendRequestLegacy("UrlcrudAgora", url, "POST", &response, &ce); err == nil { // If 5 - contrato_estado
								var r models.Resolucion
								url = ResolucionEndpoint + strconv.Itoa(m.IdResolucion)
								if err := GetRequestNew("UrlCrudResoluciones", url, &r); err == nil {
									if err := CambiarEstadoResolucion(r.Id, "REXP", vinculacion.ContratoCancelado.Usuario); err == nil {
										r.FechaExpedicion = m.FechaExpedicion
										if err := SendRequestNew("UrlCrudResoluciones", url, "PUT", &response, &r); err == nil {
											if documento, err := AlmacenarResolucionGestorDocumental(r.Id); err == nil {
												r.NuxeoUid = documento.Enlace
												if err := SendRequestNew("UrlCrudResoluciones", url, "PUT", &response, &r); err != nil { // if 10
													fmt.Println("Error en if 10 - Resolucion (PUT)#2!", err)
													logs.Error(r)
													panic(err.Error())
												}
											} else { // If 9
												fmt.Println("Error en if 9 - GestorDocumental (POST)!", err)
												logs.Error(documento)
												panic(err)
											}
										} else { // If 8
											fmt.Println("Error en if 8 - resolucion (PUT)!", err)
											logs.Error(r)
											panic(err.Error())
										}
									} else { // If 7
										fmt.Println("Error en if 7 - CambiarEstadoResolucion (POST)!", err)
										logs.Error(err)
										panic(err.Error())
									}
								} else { // If 6
									fmt.Println("Error en if 6 - resolucion (get)!", err)
									logs.Error(r)
									panic(err.Error())
								}
							} else { // If 5
								fmt.Println("Error en if 5 - contrato_estado (post)!", err)
								logs.Error(response)
								panic(err.Error())
							}
						} else { // If 4
							fmt.Println("Error en if 4 - acta_inicio (put)!", err)
							logs.Error(response)
							panic(err.Error())
						}
					} else { // If 3
						fmt.Println("Error en if 3 - acta_inicio (get)!", err)
						logs.Error(ai)
						panic(err.Error())
					}
				} else { // if 2
					fmt.Println("Error en if 2 - contrato_cancelado (post)!", err)
					logs.Error(response)
					panic(err.Error())
				}
			}
		} else { // If 1
			fmt.Println("Error en if 1 - vinculacion_docente (get)!", err)
			logs.Error(v)
			panic(err.Error())
		}
	}

	return
}

// Función que recopila los contratos a cancelar de acuerdo con el histórico de modificaciones
func BuscarContratosCancelar(vinculacionId int, contratos *[]models.ContratoCancelar) error {
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

	contrato := models.ContratoCancelar{
		NumeroContrato: *modVin.VinculacionDocenteCanceladaId.NumeroContrato,
		Vigencia:       modVin.VinculacionDocenteCanceladaId.Vigencia,
	}
	*contratos = append(*contratos, contrato)

	// Segundo caso de salida
	if tipoResolucion.CodigoAbreviacion == "RVIN" || tipoResolucion.CodigoAbreviacion == "RRED" {
		return nil
	}

	// Llamada recursiva para consultar una modificación anterior hasta llegar a
	// la vinculación inicial que no tiene modificaciones
	return BuscarContratosCancelar(modVin.VinculacionDocenteCanceladaId.Id, contratos)
}
