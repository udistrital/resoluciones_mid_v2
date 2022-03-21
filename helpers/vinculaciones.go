package helpers

import (
	"strconv"
	"strings"

	"github.com/astaxie/beego/logs"
	"github.com/udistrital/resoluciones_mid_v2/models"
)

// Consulta las vinculaciones asociadas a una resolución y construye un listado con la información relevante
func ListarVinculaciones(resolucionId string) (vinculaciones []models.Vinculaciones, outputError map[string]interface{}) {
	defer func() {
		if err := recover(); err != nil {
			outputError = map[string]interface{}{"funcion": "/ListarVinculaciones", "err": err, "status": "500"}
			panic(outputError)
		}
	}()
	var previnculaciones []models.VinculacionDocente
	var disponibilidad []models.DisponibilidadVinculacion
	var persona models.InformacionPersonaNatural
	var ciudad []map[string]interface{}
	var err2 map[string]interface{}

	url := "vinculacion_docente?limit=0&sortby=ProyectoCurricularId&order=asc&query=Activo:true,ResolucionVinculacionDocenteId.Id:" + resolucionId
	if err := GetRequestNew("UrlcrudResoluciones", url, &previnculaciones); err != nil {
		logs.Error(err.Error())
		panic(err.Error())
	}

	for i := range previnculaciones {
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

		vinculacion := &models.Vinculaciones{
			Id:                   previnculaciones[i].Id,
			Nombre:               persona.NomProveedor,
			TipoDocumento:        persona.TipoDocumento.ValorParametro,
			ExpedicionDocumento:  ciudad[0]["Nombre"].(string),
			PersonaId:            previnculaciones[i].PersonaId,
			NumeroHorasSemanales: previnculaciones[i].NumeroHorasSemanales,
			NumeroSemanas:        previnculaciones[i].NumeroSemanas,
			Categoria:            strings.Trim(previnculaciones[i].Categoria, " "),
			Dedicacion:           previnculaciones[i].ResolucionVinculacionDocenteId.Dedicacion,
			ValorContratoFormato: FormatMoney(int(previnculaciones[i].ValorContrato), 2),
			NumeroContrato:       previnculaciones[i].NumeroContrato,
			Vigencia:             previnculaciones[i].Vigencia,
			ProyectoCurricularId: previnculaciones[i].ProyectoCurricularId,
			Disponibilidad:       disponibilidad[0].Disponibilidad,
			RegistroPresupuestal: int(previnculaciones[i].NumeroRp),
		}
		vinculaciones = append(vinculaciones, *vinculacion)
	}

	if vinculaciones == nil {
		vinculaciones = []models.Vinculaciones{}
	}

	return vinculaciones, outputError
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
			url2 := "vinculacion_docente/" + strconv.Itoa(vinculacion.Id)
			if err3 := SendRequestNew("UrlcrudResoluciones", url2, "PUT", &vinculacion, disponibilidades[0].VinculacionDocenteId); err3 != nil {
				panic("Desactivando vinculacion -> " + err3.Error())
			}
		} else {
			modificacion[0].VinculacionDocenteCanceladaId.Activo = true
			modificacion[0].VinculacionDocenteRegistradaId.Activo = false
			url3 := "vinculacion_docente/" + strconv.Itoa(modificacion[0].VinculacionDocenteCanceladaId.Id)
			if err4 := SendRequestNew("UrlcrudResoluciones", url3, "PUT", &vinculacion, modificacion[0].VinculacionDocenteCanceladaId); err4 != nil {
				panic("Restaurando vinculacion -> " + err4.Error())
			}
			if err5 := SendRequestNew("UrlcrudResoluciones", "modificacion_vinculacion/"+strconv.Itoa(modificacion[0].Id), "DELETE", &resp, nil); err5 != nil {
				panic("Borrando modificación -> " + err5.Error())
			}
			url3 = "vinculacion_docente/" + strconv.Itoa(modificacion[0].VinculacionDocenteRegistradaId.Id)
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
			NumeroHorasSemanales:           horas,
			NumeroSemanas:                  d.NumeroSemanas,
			ResolucionVinculacionDocenteId: d.ResolucionData,
			ProyectoCurricularId:           proyecto,
			Categoria:                      d.Docentes[i].CategoriaNombre,
			DependenciaAcademica:           d.Docentes[i].DependenciaAcademica,
			DedicacionId:                   dedicacionId,
			Activo:                         true,
		}
		vinculaciones = append(vinculaciones, *vinculacion)
	}

	return vinculaciones, outputError
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
		for _, disponibilidad := range d.Disponibilidad {
			for _, rubro := range disponibilidad.Afectacion {
				dispVinculacion := &models.DisponibilidadVinculacion{
					Disponibilidad:       int(disponibilidad.Consecutivo),
					Rubro:                rubro.Padre,
					VinculacionDocenteId: &vinculacionesRegistradas[j],
					Activo:               true,
					Valor:                0,
				}

				if err3 := SendRequestNew("UrlcrudResoluciones", "disponibilidad_vinculacion", "POST", &dvRegistrada, &dispVinculacion); err3 != nil {
					logs.Error(err3.Error())
					panic("Registrando disponibilidad -> " + err3.Error())
				}
			}
		}
	}

	return v, outputError
}
