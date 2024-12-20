package helpers

import (
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/astaxie/beego/logs"
	"github.com/udistrital/resoluciones_mid_v2/models"
)

// Consulta las vinculaciones asociadas a una resolución y construye un listado con la información relevante
func ListarVinculaciones(resolucionId string, rp bool) (vinculaciones []models.Vinculaciones, outputError map[string]interface{}) {
	defer func() {
		if err := recover(); err != nil {
			outputError = map[string]interface{}{"funcion": "/ListarVinculaciones", "err": err, "status": "500"}
			panic(outputError)
		}
	}()
	var previnculaciones []models.VinculacionDocente
	//var previnculacionesAux []models.VinculacionDocente
	var disponibilidad []models.DisponibilidadVinculacion
	var persona models.InformacionPersonaNatural
	var ciudad map[string]interface{}
	var err2 map[string]interface{}

	if rp {
		fmt.Println("RPS ")
		previnculaciones, outputError = PrevinculacionesRps(resolucionId)
	} else {
		previnculaciones, outputError = Previnculaciones(resolucionId)
	}

	for i := range previnculaciones {
		fmt.Println("previnculacion ", previnculaciones[i].PersonaId)
		persona, err2 = BuscarDatosPersonalesDocente(previnculaciones[i].PersonaId)
		if err2 != nil {
			panic(err2)
		}

		// TODO: Esta consulta para obtener la ciudad de expedición de un documento deberá moverse a ubicaciones_crud
		// una vez se corrigan los valores del campo correspondiente en el esquema terceros
		// teniendo en cuenta que, a la fecha (15/03/2022), core_amazon_crud está en proceso de ser deprecada
		if err3 := GetRequestLegacy("UrlcrudCore", "ciudad/"+persona.CiudadExpedicionDocumento, &ciudad); err3 != nil {
			logs.Error(err3.Error())
			panic(err3.Error())
		}

		url := "disponibilidad_vinculacion?fields=Disponibilidad&query=VinculacionDocenteId.Id:" + strconv.Itoa(previnculaciones[i].Id)
		if err4 := GetRequestNew("UrlcrudResoluciones", url, &disponibilidad); err4 != nil {
			logs.Error(err4.Error())
			panic(err4.Error())
		}

		var proycur []models.Dependencia
		url1 := "dependencia?query=Id:" + strconv.Itoa(previnculaciones[i].ProyectoCurricularId)
		if err := GetRequestLegacy("UrlcrudOikos", url1, &proycur); err != nil { // If 6
			logs.Error(proycur)
			outputError = map[string]interface{}{"funcion": "/ValidarDatosExpedicion6", "err": err.Error(), "status": "502"}
		}

		if previnculaciones[i].NumeroContrato == nil {
			previnculaciones[i].NumeroContrato = new(string)
		}
		fmt.Println("previnculaciones ", previnculaciones[i])
		fmt.Println(previnculaciones[i].ValorContrato)
		/*fmt.Println("1 ", previnculaciones[i].Id)
		fmt.Println("2 ", persona.NomProveedor)
		fmt.Println("3 ", persona.TipoDocumento.ValorParametro)
		fmt.Println("4 ", ciudad["Nombre"].(string))
		fmt.Println("5 ", previnculaciones[i].PersonaId)
		fmt.Println("6 ", previnculaciones[i].NumeroHorasSemanales)
		fmt.Println("7 ", previnculaciones[i].NumeroSemanas)
		fmt.Println("8 ", strings.Trim(previnculaciones[i].Categoria, " "))
		fmt.Println("9 ", previnculaciones[i].ResolucionVinculacionDocenteId.Dedicacion)
		fmt.Println("10 ", FormatMoney(int(previnculaciones[i].ValorContrato), 2))
		fmt.Println("11 ", *previnculaciones[i].NumeroContrato)
		fmt.Println("12 ", previnculaciones[i].Vigencia)
		fmt.Println("13 ", previnculaciones[i].ProyectoCurricularId)
		fmt.Println("14 ", disponibilidad[0].Disponibilidad)
		fmt.Println("15 ", int(previnculaciones[i].NumeroRp))*/
		vinculacion := &models.Vinculaciones{
			Id:                       previnculaciones[i].Id,
			Nombre:                   persona.NomProveedor,
			TipoDocumento:            persona.TipoDocumento.ValorParametro,
			ExpedicionDocumento:      ciudad["Nombre"].(string),
			PersonaId:                previnculaciones[i].PersonaId,
			NumeroHorasSemanales:     previnculaciones[i].NumeroHorasSemanales,
			NumeroSemanas:            previnculaciones[i].NumeroSemanas,
			Categoria:                strings.Trim(previnculaciones[i].Categoria, " "),
			Dedicacion:               previnculaciones[i].ResolucionVinculacionDocenteId.Dedicacion,
			ValorContratoFormato:     FormatMoney(int(previnculaciones[i].ValorContrato), 2),
			NumeroContrato:           *previnculaciones[i].NumeroContrato,
			Vigencia:                 previnculaciones[i].Vigencia,
			ProyectoCurricularId:     previnculaciones[i].ProyectoCurricularId,
			ProyectoCurricularNombre: proycur[0].Nombre,
			Disponibilidad:           disponibilidad[0].Disponibilidad,
			RegistroPresupuestal:     int(previnculaciones[i].NumeroRp),
		}
		vinculaciones = append(vinculaciones, *vinculacion)
	}

	if vinculaciones == nil {
		vinculaciones = []models.Vinculaciones{}
	}

	return vinculaciones, outputError
}

func Previnculaciones(resolucionId string) (vinculacionDocente []models.VinculacionDocente, outputError map[string]interface{}) {
	defer func() {
		if err := recover(); err != nil {
			outputError = map[string]interface{}{"funcion": "/Previnculaciones", "err": err, "status": "500"}
			panic(outputError)
		}
	}()
	var previnculaciones []models.VinculacionDocente
	url := "vinculacion_docente?limit=0&sortby=ProyectoCurricularId&order=asc&query=Activo:true,ResolucionVinculacionDocenteId.Id:" + resolucionId
	if err := GetRequestNew("UrlcrudResoluciones", url, &previnculaciones); err != nil {
		logs.Error(err.Error())
		panic(err.Error())
	}
	vinculacionDocente = previnculaciones
	return vinculacionDocente, outputError
}

func PrevinculacionesRps(resolucionId string) (vinculacionDocente []models.VinculacionDocente, outputError map[string]interface{}) {
	defer func() {
		if err := recover(); err != nil {
			outputError = map[string]interface{}{"funcion": "/Previnculaciones", "err": err, "status": "500"}
			panic(outputError)
		}
	}()
	var previnculaciones []models.VinculacionDocente
	var previnculacionesAux []models.VinculacionDocente
	url := "vinculacion_docente?limit=0&sortby=ProyectoCurricularId&order=asc&query=ResolucionVinculacionDocenteId.Id:" + resolucionId
	if err := GetRequestNew("UrlcrudResoluciones", url, &previnculacionesAux); err != nil {
		logs.Error(err.Error())
		panic(err.Error())
	}
	for i := range previnculacionesAux {
		if previnculacionesAux[i].NumeroContrato != nil {
			previnculaciones = append(previnculaciones, previnculacionesAux[i])
		}
	}
	vinculacionDocente = previnculaciones
	return vinculacionDocente, outputError
}

// Desactiva las vinculaciones recibidas, si se trata de modificaciones reestablece las vinculaciones anteriores
func RetirarVinculaciones(vinculaciones []models.Vinculaciones) (outputError map[string]interface{}) {
	defer func() {
		if err := recover(); err != nil {
			outputError = map[string]interface{}{"funcion": "/RetirarVinculaciones", "err": err, "status": "500"}
			panic(outputError)
		}
	}()

	for _, vinc := range vinculaciones {
		var modificacion []models.ModificacionVinculacion
		var vinculacion models.VinculacionDocente
		var disponibilidades []models.DisponibilidadVinculacion
		var resp map[string]interface{}

		// Se consulta si hay modificaciones para elegir el procedimiento
		url := "modificacion_vinculacion?query=VinculacionDocenteRegistradaId.Id:" + strconv.Itoa(vinc.Id)
		if err := GetRequestNew("UrlcrudResoluciones", url, &modificacion); err != nil {
			panic("Consultando modificación -> " + err.Error())
		}

		if len(modificacion) == 0 {

			url := "disponibilidad_vinculacion?limit=0&query=VinculacionDocenteId.Id:" + strconv.Itoa(vinc.Id)
			if err := GetRequestNew("UrlcrudResoluciones", url, &disponibilidades); err != nil {
				panic("Consultando disponibilidades -> " + err.Error())
			}

			for _, disp := range disponibilidades {
				disp.Activo = false
				if err2 := SendRequestNew("UrlcrudResoluciones", "disponibilidad_vinculacion/"+strconv.Itoa(disp.Id), "PUT", &resp, disp); err2 != nil {
					panic("Desactivando disponibilidad -> " + err2.Error())
				}
			}

			disponibilidades[0].VinculacionDocenteId.Activo = false
			url2 := VinculacionEndpoint + strconv.Itoa(vinculacion.Id)
			if err3 := SendRequestNew("UrlcrudResoluciones", url2, "PUT", &vinculacion, disponibilidades[0].VinculacionDocenteId); err3 != nil {
				panic("Desactivando vinculacion -> " + err3.Error())
			}
		} else {
			modificacion[0].VinculacionDocenteCanceladaId.Activo = true
			modificacion[0].VinculacionDocenteRegistradaId.Activo = false
			url3 := VinculacionEndpoint + strconv.Itoa(modificacion[0].VinculacionDocenteCanceladaId.Id)
			if err4 := SendRequestNew("UrlcrudResoluciones", url3, "PUT", &vinculacion, modificacion[0].VinculacionDocenteCanceladaId); err4 != nil {
				panic("Restaurando vinculacion -> " + err4.Error())
			}
			if err5 := SendRequestNew("UrlcrudResoluciones", "modificacion_vinculacion/"+strconv.Itoa(modificacion[0].Id), "DELETE", &resp, nil); err5 != nil {
				panic("Borrando modificación -> " + err5.Error())
			}
			url3 = VinculacionEndpoint + strconv.Itoa(modificacion[0].VinculacionDocenteRegistradaId.Id)
			if err6 := SendRequestNew("UrlcrudResoluciones", url3, "PUT", &vinculacion, modificacion[0].VinculacionDocenteRegistradaId); err6 != nil {
				panic("Desactivando vinculación -> " + err6.Error())
			}

			url := "disponibilidad_vinculacion?limit=0&query=VinculacionDocenteId.Id:" + strconv.Itoa(modificacion[0].VinculacionDocenteRegistradaId.Id)
			if err := GetRequestNew("UrlcrudResoluciones", url, &disponibilidades); err != nil {
				panic("Consultando disponibilidades -> " + err.Error())
			}

			for _, disp := range disponibilidades {
				disp.Activo = false
				if err2 := SendRequestNew("UrlcrudResoluciones", "disponibilidad_vinculacion/"+strconv.Itoa(disp.Id), "PUT", &resp, disp); err2 != nil {
					panic("Desactivando disponibilidad -> " + err2.Error())
				}
			}
		}
	}

	return nil
}

// Construye un arreglo de estructuras de tipo models.VinculacionDocente con los datos de docentes y la resolución
func ConstruirVinculaciones(d models.ObjetoPrevinculaciones) (v []models.VinculacionDocente, outputError map[string]interface{}) {
	defer func() {
		if err := recover(); err != nil {
			outputError = map[string]interface{}{"funcion": "/ConstruirVinculaciones", "err": err, "status": "500"}
			panic(outputError)
		}
	}()

	var vinculaciones []models.VinculacionDocente
	var resolucion models.Resolucion
	var tipoResolucion models.Parametro
	var valorPunto float64
	vigencia := strconv.Itoa(d.Vigencia)
	for i := range d.Docentes {
		docDocente, e1 := strconv.Atoi(d.Docentes[i].DocDocente)
		horas, e2 := strconv.Atoi(d.Docentes[i].HorasLectivas)
		proyecto, e3 := strconv.Atoi(d.Docentes[i].IDProyecto)
		dedicacionId, e4 := strconv.Atoi(d.Docentes[i].IDTipoVinculacion)
		if e1 != nil || e2 != nil || e3 != nil || e4 != nil {
			panic("Error de conversión de datos")
		}
		vinculacion := &models.VinculacionDocente{
			Vigencia:                       d.Vigencia,
			PersonaId:                      float64(docDocente),
			NumeroContrato:                 nil,
			NumeroHorasSemanales:           horas,
			NumeroSemanas:                  d.NumeroSemanas,
			ResolucionVinculacionDocenteId: d.ResolucionData,
			ProyectoCurricularId:           proyecto,
			Categoria:                      d.Docentes[i].CategoriaNombre,
			DependenciaAcademica:           d.Docentes[i].DependenciaAcademica,
			DedicacionId:                   dedicacionId,
			Activo:                         true,
		}

		if d.ResolucionData.NivelAcademico == "PREGRADO" {
			//1.)Validar el tipo de resolucion
			// 1.1 ) Se obtiene la resolucion
			url := "resolucion/" + strconv.Itoa(vinculacion.ResolucionVinculacionDocenteId.Id)
			err := GetRequestNew("UrlCrudResoluciones", url, &resolucion)
			if err != nil {
				outputError = map[string]interface{}{"funcion": "/ConstruirVinculaciones-res", "err": err, "status": "500"}
				panic(outputError)
			}

			// 1.2 ) Se obtiene el tipo de resolucion
			url2 := ParametroEndpoint + strconv.Itoa(resolucion.TipoResolucionId)
			err2 := GetRequestNew("UrlCrudParametros", url2, &tipoResolucion)
			if err2 != nil {
				outputError = map[string]interface{}{"funcion": "/ConstruirVinculaciones-tipores", "err": err2, "status": "500"}
				panic(outputError)
			}

			//Se valida si es una resolucion de vinculacion
			if tipoResolucion.CodigoAbreviacion == "RVIN" {
				// Se aplica el valor del punto de parametros
				fmt.Println("entra")
				_, vP, err3 := CargarParametroPeriodo(vigencia, "PSAL")
				if err3 != nil {
					logs.Error(err3)
					panic(err3)
				}
				valorPunto = vP
			}
			vinculacion.ValorPuntoSalarial = valorPunto
		}
		if d.ResolucionData.NivelAcademico == "POSGRADO" {
			salarioMinimoId, _, err2 := CargarParametroPeriodo(vigencia, "SMMLV")
			if err2 != nil {
				logs.Error(err2)
				panic(err2)
			}
			vinculacion.SalarioMinimoId = salarioMinimoId
		}

		vinculaciones = append(vinculaciones, *vinculacion)
	}

	return vinculaciones, outputError
}

func EditarVinculaciones(vd models.EdicionVinculaciones) (v []models.VinculacionDocente, outputError map[string]interface{}) {
	defer func() {
		if err := recover(); err != nil {
			outputError = map[string]interface{}{"funcion": "/EditarVinculaciones", "err": err, "status": "500"}
			panic(outputError)
		}
	}()
	var vinculacionesModificadas []models.VinculacionDocente
	var err map[string]interface{}
	var disVinc []models.DisponibilidadVinculacion

	for i := range vd.Vinculaciones {
		vd.Vinculaciones[i].NumeroSemanas = vd.Semanas
	}
	if vd.Vinculaciones, err = CalcularSalarioPrecontratacion(vd.Vinculaciones); err != nil {
		panic(err)
	}
	for i := range vd.Vinculaciones {
		var vRegistrada models.VinculacionDocente
		url := VinculacionEndpoint + strconv.Itoa(vd.Vinculaciones[i].Id)
		if err2 := SendRequestNew("UrlcrudResoluciones", url, "PUT", &vRegistrada, vd.Vinculaciones[i]); err2 != nil {
			logs.Error(err2.Error())
			panic("Registrando vinculacion -> " + err2.Error())
		}
		vinculacionesModificadas = append(vinculacionesModificadas, vRegistrada)
	}
	for j := range vinculacionesModificadas {
		if err := GetRequestNew("UrlCrudResoluciones", "disponibilidad_vinculacion?query=vinculacion_docente_id:"+strconv.Itoa(vinculacionesModificadas[j].Id)+",activo:true", &disVinc); err != nil {
			panic(err.Error())
		}

		if vd.Dedicacion != "HCH" {
			desagregado, err := CalcularDesagregadoTitan(vinculacionesModificadas[j], vd.Dedicacion, vd.NivelAcademico)
			if err != nil {
				panic(err)
			}

			for nombre, valor := range desagregado {
				if nombre != "NumeroContrato" && nombre != "Vigencia" {
					var index int
					for k := range disVinc {
						if disVinc[k].Rubro == nombre {
							index = k
						}
					}
					disVinc[index].Valor = valor.(float64)
					var dvModificada models.DisponibilidadVinculacion
					url := "disponibilidad_vinculacion/" + strconv.Itoa(disVinc[index].Id)
					if err3 := SendRequestNew("UrlcrudResoluciones", url, "PUT", &dvModificada, disVinc[index]); err3 != nil {
						logs.Error(err3.Error())
						panic("Modificando disponibilidad vinculación -> " + err3.Error())
					}
				}
			}
		} else {
			var index int
			disVinc[index].Valor = vinculacionesModificadas[j].ValorContrato
			var dvModificada models.DisponibilidadVinculacion
			url := "disponibilidad_vinculacion/" + strconv.Itoa(disVinc[index].Id)
			if err3 := SendRequestNew("UrlcrudResoluciones", url, "PUT", &dvModificada, disVinc[index]); err3 != nil {
				logs.Error(err3.Error())
				panic("Modificando disponibilidad vinculación -> " + err3.Error())
			}
		}
	}

	return vinculacionesModificadas, outputError
}

// Registra en el CRUD a traves de POST las vinculaciones de los docentes y la disponibilidad correspondiente con los rubros elegidos
func RegistrarVinculaciones(d models.ObjetoPrevinculaciones) (v []models.VinculacionDocente, outputError map[string]interface{}) {
	defer func() {
		if err := recover(); err != nil {
			outputError = map[string]interface{}{"funcion": "/RegistrarVinculaciones", "err": err, "status": "500"}
			panic(outputError)
		}
	}()
	var vinculaciones []models.VinculacionDocente
	var err map[string]interface{}

	if vinculaciones, err = ConstruirVinculaciones(d); err != nil {
		panic(err)
	}

	if vinculaciones, err = CalcularSalarioPrecontratacion(vinculaciones); err != nil {
		panic(err)
	}
	var vinculacionesRegistradas []models.VinculacionDocente
	for i := range vinculaciones {
		var vRegistrada models.VinculacionDocente
		if err2 := SendRequestNew("UrlcrudResoluciones", "vinculacion_docente", "POST", &vRegistrada, &vinculaciones[i]); err2 != nil {
			logs.Error(err2.Error())
			panic("Registrando vinculacion -> " + err2.Error())
		}
		vinculacionesRegistradas = append(vinculacionesRegistradas, vRegistrada)
	}

	var dvRegistrada models.DisponibilidadVinculacion

	for j := range vinculacionesRegistradas {

		if d.ResolucionData.Dedicacion != "HCH" {

			// var desagregado models.DesagregadoContrato
			desagregado, err := CalcularDesagregadoTitan(vinculacionesRegistradas[j], d.ResolucionData.Dedicacion, d.ResolucionData.NivelAcademico)
			if err != nil {
				panic(err)
			}

			for _, disponibilidad := range d.Disponibilidad {
				// for _, rubro := range disponibilidad.Afectacion {
				// TODO La idea es cruzar los rubros (Afectacion) seleccionados en la Disponibilidad con los valores calculados para cada uno
				// una vez salga kronos a producción, de manera que el valor calculado con Titan se corresponda con el rubro de Kronos
				for nombre, valor := range desagregado {
					if nombre != "NumeroContrato" && nombre != "Vigencia" {
						dispVinculacion := &models.DisponibilidadVinculacion{
							Disponibilidad:       int(disponibilidad.Consecutivo),
							Rubro:                nombre,
							NombreRubro:          "", // rubro.Padre,
							VinculacionDocenteId: &models.VinculacionDocente{Id: vinculacionesRegistradas[j].Id},
							Activo:               true,
							Valor:                valor.(float64),
						}

						if err3 := SendRequestNew("UrlcrudResoluciones", "disponibilidad_vinculacion", "POST", &dvRegistrada, &dispVinculacion); err3 != nil {
							logs.Error(err3.Error())
							panic("Registrando disponibilidad -> " + err3.Error())
						}
					}
				}
			}
		} else {
			for _, disponibilidad := range d.Disponibilidad {
				dispVinculacion := &models.DisponibilidadVinculacion{
					Disponibilidad:       int(disponibilidad.Consecutivo),
					Rubro:                "SueldoBasico", // nombre, // rubro.Padre,
					NombreRubro:          "",
					VinculacionDocenteId: &models.VinculacionDocente{Id: vinculacionesRegistradas[j].Id},
					Activo:               true,
					Valor:                vinculacionesRegistradas[j].ValorContrato,
				}
				if err3 := SendRequestNew("UrlcrudResoluciones", "disponibilidad_vinculacion", "POST", &dvRegistrada, &dispVinculacion); err3 != nil {
					logs.Error(err3.Error())
					panic("Registrando disponibilidad -> " + err3.Error())
				}
			}
		}
	}

	return vinculacionesRegistradas, outputError
}

// Registra la modificación de una vinculación asociada a la vinculación original
func ModificarVinculaciones(obj models.ObjetoModificaciones) (v models.VinculacionDocente, outputError map[string]interface{}) {
	defer func() {
		if err := recover(); err != nil {
			outputError = map[string]interface{}{"funcion": "/ModificarVinculaciones", "err": err, "status": "500"}
			panic(outputError)
		}
	}()
	var vinculacion models.VinculacionDocente
	var vin []models.VinculacionDocente
	var desagregado map[string]interface{}
	var err map[string]interface{}

	// Recuperación de la vinculación original
	url := VinculacionEndpoint + strconv.Itoa(obj.CambiosVinculacion.VinculacionOriginal.Id)
	if err := GetRequestNew("UrlcrudResoluciones", url, &vinculacion); err != nil {
		panic("Cargando vinculacion original -> " + err.Error())
	}

	// Si solo se modificaron las horas, las semanas son las que falten para terminar
	if obj.CambiosVinculacion.NumeroSemanas == 0 {
		var err2 error
		obj.CambiosVinculacion.NumeroSemanas, err2 = CalcularNumeroSemanas(obj.CambiosVinculacion.FechaInicio, *vinculacion.NumeroContrato, vinculacion.Vigencia)
		if err2 != nil {
			panic("Error en acta de inicio " + err2.Error())
		}
	} else if obj.CambiosVinculacion.NumeroHorasSemanales == 0 {
		// Si solo se modificaron las semanas, las horas son las mismas de la vinc original
		// Aplica solo para cancelaciones de pregrado
		valores := make(map[string]float64)
		if err := CalcularTrazabilidad(strconv.Itoa(vinculacion.Id), &valores); err != nil {
			logs.Error("Error en trazabilidad -> " + err.Error())
			panic("Error en trazabilidad -> " + err.Error())
		}
		var tipoResolucion models.Parametro
		var resolucion models.Resolucion
		obj.CambiosVinculacion.NumeroHorasSemanales = int(valores["NumeroHorasSemanales"])
		err2 := GetRequestNew("UrlCrudResoluciones", "resolucion/"+strconv.Itoa(vinculacion.ResolucionVinculacionDocenteId.Id), &resolucion)
		if err2 != nil {
			panic(err2.Error())
		}
		err3 := GetRequestNew("UrlcrudParametros", "parametro/"+strconv.Itoa(resolucion.TipoResolucionId), &tipoResolucion)
		if err3 != nil {
			panic(err3.Error())
		}
		if tipoResolucion.CodigoAbreviacion == "RVIN" || tipoResolucion.CodigoAbreviacion == "RADD" {
			obj.CambiosVinculacion.NumeroHorasSemanales += obj.CambiosVinculacion.VinculacionOriginal.NumeroHorasSemanales
		} else {
			obj.CambiosVinculacion.NumeroHorasSemanales -= obj.CambiosVinculacion.VinculacionOriginal.NumeroHorasSemanales
		}
	}

	// Creación de la nueva vinculación
	nuevaVinculacion := models.VinculacionDocente{
		Vigencia:                       obj.CambiosVinculacion.VinculacionOriginal.Vigencia,
		PersonaId:                      obj.CambiosVinculacion.VinculacionOriginal.PersonaId,
		NumeroHorasSemanales:           obj.CambiosVinculacion.NumeroHorasSemanales,
		NumeroHorasTrabajadas:          obj.CambiosVinculacion.NumeroHorasTrabajadas,
		NumeroSemanas:                  obj.CambiosVinculacion.NumeroSemanas,
		ResolucionVinculacionDocenteId: obj.ResolucionNuevaId,
		DedicacionId:                   vinculacion.DedicacionId,
		ProyectoCurricularId:           vinculacion.ProyectoCurricularId,
		Categoria:                      vinculacion.Categoria,
		DependenciaAcademica:           vinculacion.DependenciaAcademica,
		ValorPuntoSalarial:             vinculacion.ValorPuntoSalarial,
		SalarioMinimoId:                vinculacion.SalarioMinimoId,
		FechaInicio:                    obj.CambiosVinculacion.FechaInicio,
		Activo:                         true,
	}

	vin = append(vin, nuevaVinculacion)
	//fmt.Println("VIN ", vin)
	// calculo del valor del contrato para la nueva vinculación
	if vin, err = CalcularSalarioPrecontratacion(vin); err != nil {
		panic(err)
	}

	nuevaVinculacion = vin[0]

	// Si el documento es RP se almacenan los datos relevantes
	if obj.CambiosVinculacion.DocPresupuestal != nil && obj.CambiosVinculacion.DocPresupuestal.Tipo == "rp" {
		nuevaVinculacion.NumeroRp = obj.CambiosVinculacion.DocPresupuestal.Consecutivo
		nuevaVinculacion.VigenciaRp = float64(obj.CambiosVinculacion.DocPresupuestal.Vigencia)
	} else {
		nuevaVinculacion.NumeroRp = float64(obj.CambiosVinculacion.VinculacionOriginal.RegistroPresupuestal)
		nuevaVinculacion.VigenciaRp = float64(obj.CambiosVinculacion.VinculacionOriginal.Vigencia)
	}

	tipoResolucion := GetTipoResolucion(obj.ResolucionNuevaId.Id)

	if tipoResolucion.CodigoAbreviacion == "RADD" {
		nuevaVinculacion.NumeroRp = 0
		nuevaVinculacion.VigenciaRp = 0
	}

	// Se desactiva la vinculación original, asi no estará disponible para ser modificada
	var vinc *models.VinculacionDocente
	vinculacion.Activo = false
	if err2 := SendRequestNew("UrlcrudResoluciones", url, "PUT", &vinc, &vinculacion); err2 != nil {
		panic("Desactivando vinculacion -> " + err2.Error())
	}
	vinc = nil

	// Se registra la nueva vinculación
	if err3 := SendRequestNew("UrlcrudResoluciones", "vinculacion_docente", "POST", &vinc, &nuevaVinculacion); err3 != nil {
		panic("Registrando nueva vinculacion -> " + err3.Error())
	}

	// Se crea y se registra la modificación de la vinculación
	var modvinc models.ModificacionVinculacion
	modificacionVinculacion := models.ModificacionVinculacion{
		ModificacionResolucionId:       &models.ModificacionResolucion{Id: obj.ModificacionResolucionId},
		VinculacionDocenteCanceladaId:  &models.VinculacionDocente{Id: vinculacion.Id},
		VinculacionDocenteRegistradaId: &models.VinculacionDocente{Id: vinc.Id},
		Horas:                          float64(obj.CambiosVinculacion.NumeroHorasSemanales),
		Activo:                         true,
	}

	if err4 := SendRequestNew("UrlcrudResoluciones", "modificacion_vinculacion", "POST", &modvinc, &modificacionVinculacion); err4 != nil {
		panic("Registrando modificacion -> " + err4.Error())
	}

	if obj.ResolucionNuevaId.Dedicacion != "HCH" {
		desagregado, err = CalcularDesagregadoTitan(*vinc, obj.ResolucionNuevaId.Dedicacion, obj.ResolucionNuevaId.NivelAcademico, tipoResolucion.CodigoAbreviacion)
		if err != nil {
			panic(err)
		}

		var dvRegistrada models.DisponibilidadVinculacion
		// Se registran los rubros de la disponibilidad segun el caso
		if obj.CambiosVinculacion.DocPresupuestal == nil || obj.CambiosVinculacion.DocPresupuestal.Tipo == "rp" {
			// Si no se cambia la disponibilidad se usa la misma de la vinculación original
			var dv []models.DisponibilidadVinculacion

			url := "disponibilidad_vinculacion?limit=0&query=VinculacionDocenteId.Id:" + strconv.Itoa(vinculacion.Id)
			if err5 := GetRequestNew("UrlcrudResoluciones", url, &dv); err5 != nil {
				panic("Cargando disponibilidad_vinculacion -> " + err5.Error())
			}
			for i := range dv {
				nuevaDv := &models.DisponibilidadVinculacion{
					Disponibilidad:       dv[i].Disponibilidad,
					Rubro:                dv[i].Rubro,
					NombreRubro:          dv[i].NombreRubro,
					VinculacionDocenteId: &models.VinculacionDocente{Id: vinc.Id},
					Activo:               true,
					Valor:                0,
				}
				// Se coloca la validacion de prima de servicios para la reduccion del contrato de un TCO a
				// Menos de 24 semanas
				if (dv[i].Rubro == "PrimaServicios") && ((obj.CambiosVinculacion.VinculacionOriginal.NumeroSemanas-obj.CambiosVinculacion.NumeroSemanas < 24) && obj.CambiosVinculacion.VinculacionOriginal.Dedicacion == "TCO") {
					nuevaDv.Valor = dv[i].Valor
				} else {
					nuevaDv.Valor = desagregado[dv[i].Rubro].(float64)
				}
				if err6 := SendRequestNew("UrlcrudResoluciones", "disponibilidad_vinculacion", "POST", &dvRegistrada, &nuevaDv); err6 != nil {
					panic("Registrando disponibilidad -> " + err6.Error())
				}
			}
		} else {
			disponibilidad := obj.CambiosVinculacion.DocPresupuestal
			// for _, rubro := range disponibilidad.Afectacion {
			// TODO La idea es cruzar los rubros (Afectacion) seleccionados en la Disponibilidad con los valores calculados para cada uno
			// una vez salga kronos a producción, de manera que el valor calculado con Titan se corresponda con el rubro de Kronos
			for nombre, valor := range desagregado {
				if nombre != "NumeroContrato" && nombre != "Vigencia" {
					nuevaDv := &models.DisponibilidadVinculacion{
						Disponibilidad:       int(disponibilidad.Consecutivo),
						Rubro:                nombre,
						NombreRubro:          "", // rubro.Padre,
						VinculacionDocenteId: &models.VinculacionDocente{Id: vinc.Id},
						Activo:               true,
						Valor:                valor.(float64),
					}
					if err6 := SendRequestNew("UrlcrudResoluciones", "disponibilidad_vinculacion", "POST", &dvRegistrada, &nuevaDv); err6 != nil {
						panic("Registrando disponibilidad -> " + err6.Error())
					}
				}
			}
		}
	} else {
		var dvRegistrada models.DisponibilidadVinculacion
		var numeroDisponibilidad int
		// Se registran los rubros de la disponibilidad segun el caso
		if obj.CambiosVinculacion.DocPresupuestal == nil || obj.CambiosVinculacion.DocPresupuestal.Tipo == "rp" {
			// Si no se cambia la disponibilidad se usa la misma de la vinculación original
			var dv []models.DisponibilidadVinculacion

			url := "disponibilidad_vinculacion?limit=0&query=VinculacionDocenteId.Id:" + strconv.Itoa(vinculacion.Id)
			if err5 := GetRequestNew("UrlcrudResoluciones", url, &dv); err5 != nil {
				panic("Cargando disponibilidad_vinculacion -> " + err5.Error())
			}
			numeroDisponibilidad = dv[0].Disponibilidad
		} else {
			numeroDisponibilidad = int(obj.CambiosVinculacion.DocPresupuestal.Consecutivo)
		}
		nuevaDv := &models.DisponibilidadVinculacion{
			Disponibilidad:       numeroDisponibilidad,
			Rubro:                "SueldoBasico",
			NombreRubro:          "",
			VinculacionDocenteId: &models.VinculacionDocente{Id: vinc.Id},
			Activo:               true,
			Valor:                nuevaVinculacion.ValorContrato,
		}

		if err3 := SendRequestNew("UrlcrudResoluciones", "disponibilidad_vinculacion", "POST", &dvRegistrada, &nuevaDv); err3 != nil {
			logs.Error(err3.Error())
			panic("Registrando disponibilidad -> " + err3.Error())
		}
	}

	return *vinc, outputError
}

// Registra la cancelación de las vinculaciones seleccionadas como modificaciones
func RegistrarCancelaciones(p models.ObjetoCancelaciones) (v []models.VinculacionDocente, outputError map[string]interface{}) {
	defer func() {
		if err := recover(); err != nil {
			outputError = map[string]interface{}{"funcion": "/RegistrarCancelaciones", "err": err, "status": "500"}
			panic(outputError)
		}
	}()
	var cancelacionesRegistradas []models.VinculacionDocente
	var actaInicio []models.ActaInicio
	var vinculacionDocente models.VinculacionDocente
	var resolucion models.Resolucion
	var parametro models.Parametro
	for i := range p.CambiosVinculacion {
		fmt.Println("FEHCA INICIO ", p.CambiosVinculacion[i].FechaInicio)
		vin := p.CambiosVinculacion[i].VinculacionOriginal
		url := "acta_inicio?query=NumeroContrato:" + vin.NumeroContrato + ",Vigencia:" + strconv.Itoa(vin.Vigencia)
		if err := GetRequestLegacy("UrlcrudAgora", url, &actaInicio); err != nil {
			panic("Acta de inicio -> " + err.Error())
		}
		//url3 := "resolucion_estado?order=desc&sortby=Id&query=Activo:true,ResolucionId.Id:" + strconv.Itoa(res[i].Id)
		url1 := "vinculacion_docente/" + strconv.Itoa(vin.Id)
		if err := GetRequestNew("UrlcrudResoluciones", url1, &vinculacionDocente); err != nil {
			panic(err.Error())
		}

		/*if err := GetRequestLegacy("UrlCrudResoluciones", url1, &vinculacionDocente); err != nil {
			panic("Vinculación docente -> " + err.Error())
		}*/
		url2 := "resolucion/" + strconv.Itoa(vinculacionDocente.ResolucionVinculacionDocenteId.Id)
		if err := GetRequestNew("UrlcrudResoluciones", url2, &resolucion); err != nil {
			panic(err.Error())
		}
		url3 := "parametro/" + strconv.Itoa(resolucion.TipoResolucionId)
		if err := GetRequestNew("UrlcrudParametros", url3, &parametro); err != nil {
			logs.Error(err)
			panic("Cargando tipo_resolucion -> " + err.Error())
		}
		semanasFinales := vin.NumeroSemanas - p.CambiosVinculacion[i].NumeroSemanas
		if parametro.CodigoAbreviacion == "RADD" || parametro.CodigoAbreviacion == "RRED" {
			semanasFinales -= 1
		}
		p.CambiosVinculacion[i].FechaInicio = CalcularFechaFin(actaInicio[0].FechaInicio, semanasFinales)
		cancelacion := models.ObjetoModificaciones{
			CambiosVinculacion:       &p.CambiosVinculacion[i],
			ResolucionNuevaId:        p.ResolucionNuevaId,
			ModificacionResolucionId: p.ModificacionResolucionId,
		}
		if cancelacionRegistrada, err := ModificarVinculaciones(cancelacion); err != nil {
			panic(err)
		} else {
			cancelacionesRegistradas = append(cancelacionesRegistradas, cancelacionRegistrada)
		}
	}

	return cancelacionesRegistradas, outputError
}

// Unifica los valores de la vinculación atraves de las diferentes modificaciones que ha tenido
func CalcularTrazabilidad(vinculacionId string, valoresAntes *map[string]float64) error {
	var modificaciones []models.ModificacionVinculacion
	var modVin models.ModificacionVinculacion
	var resolucion models.Resolucion
	var tipoResolucion models.Parametro

	url := "modificacion_vinculacion?query=VinculacionDocenteRegistradaId.Id:" + vinculacionId
	if err := GetRequestNew("UrlCrudResoluciones", url, &modificaciones); err != nil {
		logs.Error(err.Error())
		return err
	}
	fmt.Println("Modificaciones ", modificaciones)
	// Caso de salida
	if len(modificaciones) == 0 {
		return nil
	}

	modVin = modificaciones[0]
	vinculacionAnteriorId := strconv.Itoa(modVin.VinculacionDocenteCanceladaId.Id)

	var desagregadoAntes []models.DisponibilidadVinculacion
	url2 := "disponibilidad_vinculacion?query=Activo:true,VinculacionDocenteId.Id:" + vinculacionAnteriorId
	if err2 := GetRequestNew("UrlCrudResoluciones", url2, &desagregadoAntes); err2 != nil {
		logs.Error(err2.Error())
		return err2
	}

	urlaux := "resolucion/" + strconv.Itoa(modVin.VinculacionDocenteCanceladaId.ResolucionVinculacionDocenteId.Id)
	if erraux := GetRequestNew("UrlCrudResoluciones", urlaux, &resolucion); erraux != nil {
		logs.Error(erraux.Error())
		return erraux
	}

	url3 := ParametroEndpoint + strconv.Itoa(resolucion.TipoResolucionId)
	if err3 := GetRequestNew("UrlcrudParametros", url3, &tipoResolucion); err3 != nil {
		logs.Error(err3.Error())
		return err3
	}

	for _, disp := range desagregadoAntes {
		if tipoResolucion.CodigoAbreviacion == "RVIN" || tipoResolucion.CodigoAbreviacion == "RADD" {
			(*valoresAntes)[disp.Rubro] = disp.Valor + (*valoresAntes)[disp.Rubro]
		} else {
			(*valoresAntes)[disp.Rubro] = (*valoresAntes)[disp.Rubro] - disp.Valor
		}
	}
	fmt.Println("Tipo de resolución ", tipoResolucion)
	switch tipoResolucion.CodigoAbreviacion {
	case "RCAN":
		(*valoresAntes)["NumeroSemanas"] = float64(modVin.VinculacionDocenteCanceladaId.NumeroSemanas) - (*valoresAntes)["NumeroSemanas"]
		(*valoresAntes)["ValorContrato"] = float64(modVin.VinculacionDocenteCanceladaId.ValorContrato) - (*valoresAntes)["ValorContrato"]
		(*valoresAntes)["NumeroHorasSemanales"] = float64(modVin.VinculacionDocenteCanceladaId.NumeroHorasSemanales)
		break
	case "RRED":
		(*valoresAntes)["NumeroHorasSemanales"] = (*valoresAntes)["NumeroHorasSemanales"] - float64(modVin.VinculacionDocenteCanceladaId.NumeroHorasSemanales)
		(*valoresAntes)["ValorContrato"] = (*valoresAntes)["ValorContrato"] - float64(modVin.VinculacionDocenteCanceladaId.ValorContrato)
		break
	default:
		(*valoresAntes)["NumeroHorasSemanales"] = float64(modVin.VinculacionDocenteCanceladaId.NumeroHorasSemanales) + (*valoresAntes)["NumeroHorasSemanales"]
		(*valoresAntes)["ValorContrato"] = float64(modVin.VinculacionDocenteCanceladaId.ValorContrato) + (*valoresAntes)["ValorContrato"]
		(*valoresAntes)["NumeroSemanas"] = float64(modVin.VinculacionDocenteCanceladaId.NumeroSemanas)
		break
	}

	// Llamada recursiva para consultar una modificación anterior hasta llegar a
	// la vinculación inicial que no tiene modificaciones
	return CalcularTrazabilidad(vinculacionAnteriorId, valoresAntes)

}

// Calcula el numero de semanas entre la fecha recibida y la fecha fin de la vinculación dada
func CalcularNumeroSemanas(fechaInicio time.Time, NumeroContrato string, Vigencia int) (numeroSemanas int, err error) {
	var actaInicio []models.ActaInicio

	url2 := "acta_inicio?query=NumeroContrato:" + NumeroContrato + ",Vigencia:" + strconv.Itoa(Vigencia)
	if err = GetRequestLegacy("UrlcrudAgora", url2, &actaInicio); err != nil {
		return numeroSemanas, err
	}
	diferencia := actaInicio[0].FechaFin.Sub(fechaInicio)
	dif := diferencia.Hours() / (24 * 7)
	numeroSemanas = int(math.Ceil(dif))
	return
}

// Registra numero y vigencia de RP en las vinculaciones con el id correspondiente
func RegistrarVinculacionesRp(registros []models.RpSeleccionado) (outputError map[string]interface{}) {
	defer func() {
		if err := recover(); err != nil {
			outputError = map[string]interface{}{"funcion": "/RegistrarVinculacionesRp", "err": err, "status": "500"}
			panic(outputError)
		}
	}()

	var v *models.VinculacionDocente
	var v1 []models.VinculacionDocente
	for _, rp := range registros {
		// Recuperación de la vinculación original
		v = nil
		url := "vinculacion_docente?query=Id:" + strconv.Itoa(rp.VinculacionId)
		if err := GetRequestNew("UrlcrudResoluciones", url, &v1); err != nil {
			panic("Cargando vinculacion original -> " + err.Error())
		} else if len(v1) == 0 {
			panic("No se encontró la vinculacion original")
		}
		v = &v1[0]
		v.NumeroRp = float64(rp.Consecutivo)
		v.VigenciaRp = float64(rp.Vigencia)

		// Ejecutar preliquidacion
		if err2 := EjecutarPreliquidacionTitan(*v); err2 != nil {
			panic(err2)
		}
	}

	return nil
}
