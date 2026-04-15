package helpers

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/astaxie/beego/logs"
	"github.com/udistrital/resoluciones_mid_v2/models"
)

type contextoConstruccionVinculacion struct {
	tipoResolucion  models.Parametro
	valorPunto      float64
	salarioMinimoID int
}

func cargarDisponibilidadesVinculacion(idVinculacion int) ([]models.DisponibilidadVinculacion, error) {
	var disponibilidades []models.DisponibilidadVinculacion
	url := "disponibilidad_vinculacion?limit=0&query=VinculacionDocenteId.Id:" + strconv.Itoa(idVinculacion)
	if err := GetRequestNew("UrlcrudResoluciones", url, &disponibilidades); err != nil {
		return nil, err
	}
	if len(disponibilidades) == 0 {
		return nil, fmt.Errorf("no se encontraron disponibilidades para la vinculación %d", idVinculacion)
	}

	return disponibilidades, nil
}

func cargarDisponibilidadesResumenVinculacion(idVinculacion int) ([]models.DisponibilidadVinculacion, error) {
	var disponibilidades []models.DisponibilidadVinculacion
	url := "disponibilidad_vinculacion?fields=Disponibilidad&query=VinculacionDocenteId.Id:" + strconv.Itoa(idVinculacion)
	if err := GetRequestNew("UrlcrudResoluciones", url, &disponibilidades); err != nil {
		return nil, err
	}
	if len(disponibilidades) == 0 {
		return nil, fmt.Errorf("no se encontró resumen de disponibilidad para la vinculación %d", idVinculacion)
	}

	return disponibilidades, nil
}

func cargarProyectoCurricular(idProyectoCurricular int) (models.Dependencia, error) {
	var proyectos []models.Dependencia
	url := "dependencia?query=Id:" + strconv.Itoa(idProyectoCurricular)
	if err := GetRequestLegacy("UrlcrudOikos", url, &proyectos); err != nil {
		return models.Dependencia{}, err
	}
	if len(proyectos) == 0 {
		return models.Dependencia{}, fmt.Errorf("no se encontró proyecto curricular %d", idProyectoCurricular)
	}

	return proyectos[0], nil
}

func desactivarDisponibilidadesVinculacion(disponibilidades []models.DisponibilidadVinculacion) error {
	var resp map[string]interface{}
	for _, disp := range disponibilidades {
		disp.Activo = false
		if err := SendRequestNew("UrlcrudResoluciones", "disponibilidad_vinculacion/"+strconv.Itoa(disp.Id), "PUT", &resp, disp); err != nil {
			return err
		}
	}

	return nil
}

func cargarDisponibilidadesActivasVinculacion(idVinculacion int) ([]models.DisponibilidadVinculacion, error) {
	var disponibilidades []models.DisponibilidadVinculacion
	url := "disponibilidad_vinculacion?query=vinculacion_docente_id:" + strconv.Itoa(idVinculacion) + ",activo:true"
	if err := GetRequestNew("UrlCrudResoluciones", url, &disponibilidades); err != nil {
		return nil, err
	}
	if len(disponibilidades) == 0 {
		return nil, fmt.Errorf("no se encontraron disponibilidades activas para la vinculación %d", idVinculacion)
	}

	return disponibilidades, nil
}

func actualizarDisponibilidadVinculacion(disponibilidad models.DisponibilidadVinculacion) error {
	var actualizada models.DisponibilidadVinculacion
	url := "disponibilidad_vinculacion/" + strconv.Itoa(disponibilidad.Id)
	if err := SendRequestNew("UrlcrudResoluciones", url, "PUT", &actualizada, disponibilidad); err != nil {
		return err
	}

	return nil
}

func registrarDisponibilidadVinculacion(disponibilidad *models.DisponibilidadVinculacion) error {
	var registrada models.DisponibilidadVinculacion
	if err := SendRequestNew("UrlcrudResoluciones", "disponibilidad_vinculacion", "POST", &registrada, disponibilidad); err != nil {
		return err
	}

	return nil
}

func actualizarDisponibilidadesDesdeDesagregado(disponibilidades []models.DisponibilidadVinculacion, desagregado map[string]interface{}) error {
	for nombre, valor := range desagregado {
		if nombre == "NumeroContrato" || nombre == "Vigencia" {
			continue
		}

		index := -1
		for i := range disponibilidades {
			if disponibilidades[i].Rubro == nombre {
				index = i
				break
			}
		}
		if index == -1 {
			return fmt.Errorf("no se encontró rubro %s en las disponibilidades activas", nombre)
		}

		disponibilidades[index].Valor = valor.(float64)
		if err := actualizarDisponibilidadVinculacion(disponibilidades[index]); err != nil {
			return err
		}
	}

	return nil
}

func registrarDisponibilidadesDesdeDesagregado(idVinculacion int, numeroDisponibilidad int, desagregado map[string]interface{}) error {
	for nombre, valor := range desagregado {
		if nombre == "NumeroContrato" || nombre == "Vigencia" {
			continue
		}

		dispVinculacion := &models.DisponibilidadVinculacion{
			Disponibilidad:       numeroDisponibilidad,
			Rubro:                nombre,
			NombreRubro:          "",
			VinculacionDocenteId: &models.VinculacionDocente{Id: idVinculacion},
			Activo:               true,
			Valor:                valor.(float64),
		}
		if err := registrarDisponibilidadVinculacion(dispVinculacion); err != nil {
			return err
		}
	}

	return nil
}

func registrarDisponibilidadSueldoBasico(idVinculacion int, numeroDisponibilidad int, valor float64) error {
	dispVinculacion := &models.DisponibilidadVinculacion{
		Disponibilidad:       numeroDisponibilidad,
		Rubro:                "SueldoBasico",
		NombreRubro:          "",
		VinculacionDocenteId: &models.VinculacionDocente{Id: idVinculacion},
		Activo:               true,
		Valor:                valor,
	}

	return registrarDisponibilidadVinculacion(dispVinculacion)
}

func cargarContextoConstruccionVinculacion(data models.ObjetoPrevinculaciones) (contextoConstruccionVinculacion, error) {
	var contexto contextoConstruccionVinculacion
	vigencia := strconv.Itoa(data.Vigencia)

	if data.ResolucionData.NivelAcademico == "PREGRADO" {
		var resolucion models.Resolucion
		if err := GetRequestNew("UrlCrudResoluciones", "resolucion/"+strconv.Itoa(data.ResolucionData.Id), &resolucion); err != nil {
			return contexto, err
		}
		if err := GetRequestNew("UrlCrudParametros", ParametroEndpoint+strconv.Itoa(resolucion.TipoResolucionId), &contexto.tipoResolucion); err != nil {
			return contexto, err
		}
		if contexto.tipoResolucion.CodigoAbreviacion == "RVIN" {
			_, valorPunto, err := CargarParametroPeriodo(vigencia, "PSAL")
			if err != nil {
				return contexto, fmt.Errorf("%v", err)
			}
			contexto.valorPunto = valorPunto
		}
	}

	if data.ResolucionData.NivelAcademico == "POSGRADO" {
		salarioMinimoID, _, err := CargarParametroPeriodo(vigencia, "SMMLV")
		if err != nil {
			return contexto, fmt.Errorf("%v", err)
		}
		contexto.salarioMinimoID = salarioMinimoID
	}

	return contexto, nil
}

func construirVinculacionDocenteBase(docente models.CargaLectiva, data models.ObjetoPrevinculaciones, contexto contextoConstruccionVinculacion) (models.VinculacionDocente, error) {
	docDocente, e1 := strconv.Atoi(docente.DocDocente)
	horas, e2 := strconv.Atoi(docente.HorasLectivas)
	proyecto, e3 := strconv.Atoi(docente.IDProyecto)
	dedicacionID, e4 := strconv.Atoi(docente.IDTipoVinculacion)
	if e1 != nil || e2 != nil || e3 != nil || e4 != nil {
		return models.VinculacionDocente{}, fmt.Errorf("error de conversión de datos")
	}

	return models.VinculacionDocente{
		Vigencia:                       data.Vigencia,
		PersonaId:                      float64(docDocente),
		NumeroContrato:                 nil,
		NumeroHorasSemanales:           horas,
		NumeroSemanas:                  data.NumeroSemanas,
		ResolucionVinculacionDocenteId: data.ResolucionData,
		ProyectoCurricularId:           proyecto,
		Categoria:                      docente.CategoriaNombre,
		DependenciaAcademica:           docente.DependenciaAcademica,
		DedicacionId:                   dedicacionID,
		Activo:                         true,
		ValorPuntoSalarial:             contexto.valorPunto,
		SalarioMinimoId:                contexto.salarioMinimoID,
	}, nil
}

func cargarTipoResolucionDesdeVinculacion(vinculacion models.VinculacionDocente) (models.Parametro, error) {
	var resolucion models.Resolucion
	var tipoResolucion models.Parametro

	if err := GetRequestNew("UrlCrudResoluciones", "resolucion/"+strconv.Itoa(vinculacion.ResolucionVinculacionDocenteId.Id), &resolucion); err != nil {
		return tipoResolucion, err
	}
	if err := GetRequestNew("UrlcrudParametros", "parametro/"+strconv.Itoa(resolucion.TipoResolucionId), &tipoResolucion); err != nil {
		return tipoResolucion, err
	}

	return tipoResolucion, nil
}

func resolverNumeroDisponibilidad(vinculacionOriginalID int, docPresupuestal *models.DocumentoPresupuestal) (int, error) {
	if docPresupuestal != nil && docPresupuestal.Tipo != "rp" {
		return int(docPresupuestal.Consecutivo), nil
	}

	disponibilidades, err := cargarDisponibilidadesVinculacion(vinculacionOriginalID)
	if err != nil {
		return 0, err
	}

	return disponibilidades[0].Disponibilidad, nil
}

func aplicarDocumentoPresupuestalVinculacion(vinculacion *models.VinculacionDocente, original *models.Vinculaciones, docPresupuestal *models.DocumentoPresupuestal, tipoResolucion models.Parametro) {
	if docPresupuestal != nil && docPresupuestal.Tipo == "rp" {
		vinculacion.NumeroRp = docPresupuestal.Consecutivo
		vinculacion.VigenciaRp = float64(docPresupuestal.Vigencia)
	} else if original != nil {
		vinculacion.NumeroRp = float64(original.RegistroPresupuestal)
		vinculacion.VigenciaRp = float64(original.Vigencia)
	}

	if tipoResolucion.CodigoAbreviacion == "RADD" {
		vinculacion.NumeroRp = 0
		vinculacion.VigenciaRp = 0
	}
}

func resolverNumeroHorasModificacion(cambios *models.CambioVinculacion, vinculacionOriginal models.VinculacionDocente) error {
	if cambios.NumeroHorasSemanales != 0 {
		return nil
	}

	valores := make(map[string]float64)
	if err := CalcularTrazabilidad(strconv.Itoa(vinculacionOriginal.Id), &valores); err != nil {
		logs.Error("Error en trazabilidad -> " + err.Error())
		return fmt.Errorf("error en trazabilidad -> %s", err.Error())
	}

	cambios.NumeroHorasSemanales = int(valores["NumeroHorasSemanales"])
	tipoResolucion, err := cargarTipoResolucionDesdeVinculacion(vinculacionOriginal)
	if err != nil {
		return err
	}

	if tipoResolucion.CodigoAbreviacion == "RVIN" || tipoResolucion.CodigoAbreviacion == "RADD" {
		cambios.NumeroHorasSemanales += cambios.VinculacionOriginal.NumeroHorasSemanales
	} else {
		cambios.NumeroHorasSemanales -= cambios.VinculacionOriginal.NumeroHorasSemanales
	}

	return nil
}

func construirNuevaVinculacionModificacion(cambios *models.CambioVinculacion, resolucionNueva *models.ResolucionVinculacionDocente, vinculacionOriginal models.VinculacionDocente) models.VinculacionDocente {
	return models.VinculacionDocente{
		Vigencia:                       cambios.VinculacionOriginal.Vigencia,
		PersonaId:                      cambios.VinculacionOriginal.PersonaId,
		NumeroHorasSemanales:           cambios.NumeroHorasSemanales,
		NumeroHorasTrabajadas:          cambios.NumeroHorasTrabajadas,
		NumeroSemanas:                  cambios.NumeroSemanas,
		ResolucionVinculacionDocenteId: resolucionNueva,
		DedicacionId:                   vinculacionOriginal.DedicacionId,
		ProyectoCurricularId:           vinculacionOriginal.ProyectoCurricularId,
		Categoria:                      vinculacionOriginal.Categoria,
		DependenciaAcademica:           vinculacionOriginal.DependenciaAcademica,
		ValorPuntoSalarial:             vinculacionOriginal.ValorPuntoSalarial,
		SalarioMinimoId:                vinculacionOriginal.SalarioMinimoId,
		FechaInicio:                    cambios.FechaInicio,
		Activo:                         true,
	}
}

func cargarCiudadExpedicionDocumento(idCiudad string) (string, error) {
	var ciudad map[string]interface{}
	if err := GetRequestLegacy("UrlcrudCore", "ciudad/"+idCiudad, &ciudad); err != nil {
		return "", err
	}

	nombre, ok := ciudad["Nombre"].(string)
	if !ok || nombre == "" {
		return "", fmt.Errorf("ciudad de expedición inválida para id %s", idCiudad)
	}

	return nombre, nil
}

func construirResumenVinculacion(previnculacion models.VinculacionDocente, persona models.InformacionPersonaNatural, ciudadExpedicion string, proyectoCurricular models.Dependencia, disponibilidad models.DisponibilidadVinculacion) models.Vinculaciones {
	if previnculacion.NumeroContrato == nil {
		previnculacion.NumeroContrato = new(string)
	}

	return models.Vinculaciones{
		Id:                       previnculacion.Id,
		Nombre:                   persona.NomProveedor,
		TipoDocumento:            persona.TipoDocumento.ValorParametro,
		ExpedicionDocumento:      ciudadExpedicion,
		PersonaId:                previnculacion.PersonaId,
		NumeroHorasSemanales:     previnculacion.NumeroHorasSemanales,
		NumeroSemanas:            previnculacion.NumeroSemanas,
		Categoria:                strings.Trim(previnculacion.Categoria, " "),
		Dedicacion:               previnculacion.ResolucionVinculacionDocenteId.Dedicacion,
		ValorContratoFormato:     FormatMoney(int(previnculacion.ValorContrato), 2),
		NumeroContrato:           *previnculacion.NumeroContrato,
		Vigencia:                 previnculacion.Vigencia,
		ProyectoCurricularId:     previnculacion.ProyectoCurricularId,
		ProyectoCurricularNombre: proyectoCurricular.Nombre,
		Disponibilidad:           disponibilidad.Disponibilidad,
		RegistroPresupuestal:     int(previnculacion.NumeroRp),
	}
}

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

	if rp {
		previnculaciones, outputError = PrevinculacionesRps(resolucionId)
	} else {
		previnculaciones, outputError = Previnculaciones(resolucionId)
	}

	for i := range previnculaciones {
		persona, err2 := BuscarDatosPersonalesDocente(previnculaciones[i].PersonaId)
		if err2 != nil {
			panic(err2)
		}

		// TODO: Esta consulta para obtener la ciudad de expedición de un documento deberá moverse a ubicaciones_crud
		// una vez se corrigan los valores del campo correspondiente en el esquema terceros
		// teniendo en cuenta que, a la fecha (15/03/2022), core_amazon_crud está en proceso de ser deprecada
		ciudadExpedicion, err3 := cargarCiudadExpedicionDocumento(persona.CiudadExpedicionDocumento)
		if err3 != nil {
			logs.Error(err3.Error())
			panic(err3.Error())
		}

		disponibilidad, err4 := cargarDisponibilidadesResumenVinculacion(previnculaciones[i].Id)
		if err4 != nil {
			logs.Error(err4.Error())
			panic(err4.Error())
		}

		proycur, err := cargarProyectoCurricular(previnculaciones[i].ProyectoCurricularId)
		if err != nil { // If 6
			outputError = map[string]interface{}{"funcion": "/ValidarDatosExpedicion6", "err": err.Error(), "status": "502"}
		}

		vinculacion := construirResumenVinculacion(previnculaciones[i], persona, ciudadExpedicion, proycur, disponibilidad[0])
		vinculaciones = append(vinculaciones, vinculacion)
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
		var resp map[string]interface{}

		// Se consulta si hay modificaciones para elegir el procedimiento
		url := "modificacion_vinculacion?query=VinculacionDocenteRegistradaId.Id:" + strconv.Itoa(vinc.Id)
		if err := GetRequestNew("UrlcrudResoluciones", url, &modificacion); err != nil {
			panic("Consultando modificación -> " + err.Error())
		}

		if len(modificacion) == 0 {
			disponibilidades, err := cargarDisponibilidadesVinculacion(vinc.Id)
			if err != nil {
				panic("Consultando disponibilidades -> " + err.Error())
			}

			if err := desactivarDisponibilidadesVinculacion(disponibilidades); err != nil {
				panic("Desactivando disponibilidad -> " + err.Error())
			}

			disponibilidades[0].VinculacionDocenteId.Activo = false
			url2 := VinculacionEndpoint + strconv.Itoa(disponibilidades[0].VinculacionDocenteId.Id)
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

			disponibilidades, err := cargarDisponibilidadesVinculacion(modificacion[0].VinculacionDocenteRegistradaId.Id)
			if err != nil {
				panic("Consultando disponibilidades -> " + err.Error())
			}

			if err := desactivarDisponibilidadesVinculacion(disponibilidades); err != nil {
				panic("Desactivando disponibilidad -> " + err.Error())
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
	contexto, err := cargarContextoConstruccionVinculacion(d)
	if err != nil {
		outputError = map[string]interface{}{"funcion": "/ConstruirVinculaciones-contexto", "err": err, "status": "500"}
		panic(outputError)
	}

	for i := range d.Docentes {
		vinculacion, err := construirVinculacionDocenteBase(d.Docentes[i], d, contexto)
		if err != nil {
			outputError = map[string]interface{}{"funcion": "/ConstruirVinculaciones-base", "err": err, "status": "500"}
			panic(outputError)
		}

		vinculaciones = append(vinculaciones, vinculacion)
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
		disVinc, err := cargarDisponibilidadesActivasVinculacion(vinculacionesModificadas[j].Id)
		if err != nil {
			panic(err.Error())
		}

		if vd.Dedicacion != "HCH" {
			desagregado, err := CalcularDesagregadoTitan(vinculacionesModificadas[j], vd.Dedicacion, vd.NivelAcademico)
			if err != nil {
				panic(err)
			}

			if err := actualizarDisponibilidadesDesdeDesagregado(disVinc, desagregado); err != nil {
				logs.Error(err.Error())
				panic("Modificando disponibilidad vinculación -> " + err.Error())
			}
		} else {
			disVinc[0].Valor = vinculacionesModificadas[j].ValorContrato
			if err := actualizarDisponibilidadVinculacion(disVinc[0]); err != nil {
				logs.Error(err.Error())
				panic("Modificando disponibilidad vinculación -> " + err.Error())
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
				if err := registrarDisponibilidadesDesdeDesagregado(vinculacionesRegistradas[j].Id, int(disponibilidad.Consecutivo), desagregado); err != nil {
					logs.Error(err.Error())
					panic("Registrando disponibilidad -> " + err.Error())
				}
			}
		} else {
			for _, disponibilidad := range d.Disponibilidad {
				if err := registrarDisponibilidadSueldoBasico(vinculacionesRegistradas[j].Id, int(disponibilidad.Consecutivo), vinculacionesRegistradas[j].ValorContrato); err != nil {
					logs.Error(err.Error())
					panic("Registrando disponibilidad -> " + err.Error())
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
		if err := resolverNumeroHorasModificacion(obj.CambiosVinculacion, vinculacion); err != nil {
			panic(err.Error())
		}
	}

	nuevaVinculacion := construirNuevaVinculacionModificacion(obj.CambiosVinculacion, obj.ResolucionNuevaId, vinculacion)

	if fecha := NormalizarFechaTimezone(&nuevaVinculacion.FechaInicio); fecha != nil {
		nuevaVinculacion.FechaInicio = *fecha
	}

	vin = append(vin, nuevaVinculacion)
	// calculo del valor del contrato para la nueva vinculación
	if vin, err = CalcularSalarioPrecontratacion(vin); err != nil {
		panic(err)
	}

	nuevaVinculacion = vin[0]

	tipoResolucion := GetTipoResolucion(obj.ResolucionNuevaId.Id)
	aplicarDocumentoPresupuestalVinculacion(&nuevaVinculacion, obj.CambiosVinculacion.VinculacionOriginal, obj.CambiosVinculacion.DocPresupuestal, tipoResolucion)

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
		var objetoNovedad models.ObjetoNovedad
		objetoNovedad.TipoResolucion = tipoResolucion.CodigoAbreviacion
		/* se necesita obtener la diferencia del desagregado a restar a la vinculacion original pero la aplicacion de porcentajes depende de cuantas semanas quedara
		durando el contrato despues de la cancelacion*/
		if obj.CambiosVinculacion.VinculacionOriginal.NumeroSemanas-obj.CambiosVinculacion.NumeroSemanas < 0 {
			objetoNovedad.SemanasNuevas = 0
		} else {
			objetoNovedad.SemanasNuevas = obj.CambiosVinculacion.VinculacionOriginal.NumeroSemanas - obj.CambiosVinculacion.NumeroSemanas
		}
		objetoNovedad.VinculacionOriginal = obj.CambiosVinculacion.VinculacionOriginal.NumeroContrato
		objetoNovedad.VigenciaVinculacionOriginal = obj.CambiosVinculacion.VinculacionOriginal.Vigencia
		desagregado, err = CalcularDesagregadoTitan(*vinc, obj.ResolucionNuevaId.Dedicacion, obj.ResolucionNuevaId.NivelAcademico, &objetoNovedad)
		if err != nil {
			panic(err)
		}

		// Se registran los rubros de la disponibilidad segun el caso
		if obj.CambiosVinculacion.DocPresupuestal == nil || obj.CambiosVinculacion.DocPresupuestal.Tipo == "rp" {
			// Si no se cambia la disponibilidad se usa la misma de la vinculación original
			dv, err5 := cargarDisponibilidadesVinculacion(vinculacion.Id)
			if err5 != nil {
				panic("Cargando disponibilidad_vinculacion -> " + err5.Error())
			}
			for i := range dv {
				if err6 := registrarDisponibilidadVinculacion(&models.DisponibilidadVinculacion{
					Disponibilidad:       dv[i].Disponibilidad,
					Rubro:                dv[i].Rubro,
					NombreRubro:          dv[i].NombreRubro,
					VinculacionDocenteId: &models.VinculacionDocente{Id: vinc.Id},
					Activo:               true,
					Valor:                desagregado[dv[i].Rubro].(float64),
				}); err6 != nil {
					panic("Registrando disponibilidad -> " + err6.Error())
				}
			}
		} else {
			disponibilidad := obj.CambiosVinculacion.DocPresupuestal
			if err6 := registrarDisponibilidadesDesdeDesagregado(vinc.Id, int(disponibilidad.Consecutivo), desagregado); err6 != nil {
				panic("Registrando disponibilidad -> " + err6.Error())
			}
		}
	} else {
		numeroDisponibilidad, err := resolverNumeroDisponibilidad(vinculacion.Id, obj.CambiosVinculacion.DocPresupuestal)
		if err != nil {
			panic("Cargando disponibilidad_vinculacion -> " + err.Error())
		}
		if err3 := registrarDisponibilidadSueldoBasico(vinc.Id, numeroDisponibilidad, nuevaVinculacion.ValorContrato); err3 != nil {
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
	for i := range p.CambiosVinculacion {
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
		parametro, err := cargarTipoResolucionDesdeVinculacion(vinculacionDocente)
		if err != nil {
			logs.Error(err)
			panic("Cargando tipo_resolucion -> " + err.Error())
		}
		semanasFinales := vin.NumeroSemanas - p.CambiosVinculacion[i].NumeroSemanas
		if parametro.CodigoAbreviacion == "RADD" || parametro.CodigoAbreviacion == "RRED" {
			semanasFinales -= 1
		}
		fechasContrato := CalcularFechasContrato(actaInicio[0].FechaInicio, semanasFinales)
		p.CambiosVinculacion[i].FechaInicio = fechasContrato.FechaFinPago
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
	if diferencia < 1 {
		numeroSemanas = 0
	} else {
		numeroSemanas = int(diferencia.Hours()/24/7) + 1
	}
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

func ObtenerVinculacionPorResolucionYDocente(resolucionId int, numeroDocumento int) (vinculacion *models.VinculacionDocente, outputError map[string]interface{}) {
	defer func() {
		if err := recover(); err != nil {
			outputError = map[string]interface{}{
				"funcion": "/ObtenerVinculacionPorResolucionYDocente",
				"err":     err,
				"status":  "500",
			}
		}
	}()

	var vinculaciones []models.VinculacionDocente
	url := "vinculacion_docente?query=ResolucionVinculacionDocenteId.Id:" +
		strconv.Itoa(resolucionId) + ",PersonaId:" + strconv.Itoa(numeroDocumento)

	if err := GetRequestNew("UrlcrudResoluciones", url, &vinculaciones); err != nil {
		return nil, map[string]interface{}{
			"funcion": "/ObtenerVinculacionPorResolucionYDocente",
			"err":     err.Error(),
			"status":  "502",
		}
	}

	if len(vinculaciones) == 0 {
		return nil, map[string]interface{}{
			"funcion": "/ObtenerVinculacionPorResolucionYDocente",
			"err":     "No se encontró la vinculación para la resolución y docente indicados",
			"status":  "404",
		}
	}

	return &vinculaciones[0], nil
}
