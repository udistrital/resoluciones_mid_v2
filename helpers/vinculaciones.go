package helpers

import (
	"strconv"
	"strings"

	"github.com/udistrital/resoluciones_mid_v2/models"
)

func ListarVinculaciones(resolucionId string) (vinculaciones []models.Vinculaciones, outputError map[string]interface{}) {
	defer func() {
		if err := recover(); err != nil {
			outputError = map[string]interface{}{"funcion": "/ListarVinculaciones", "err": err, "status": "500"}
			panic(outputError)
		}
	}()
	var previnculaciones []models.VinculacionDocente

	url := "vinculacion_docente?limit=0&sortby=ProyectoCurricularId&order=asc&query=Activo:true,ResolucionVinculacionDocenteId.Id:" + resolucionId
	if err := GetRequestNew("UrlcrudResoluciones", url, &previnculaciones); err != nil {
		panic(err.Error())
	}

	for i := range previnculaciones {
		/*
		 * Buscar de agora:
		 *		Nombre completo
		 *		Tipo documento
		 *		Lugar expedicion documento
		 */

		vinculacion := &models.Vinculaciones{
			Id:                   previnculaciones[i].Id,
			Nombre:               "",
			PersonaId:            previnculaciones[i].PersonaId,
			NumeroHorasSemanales: previnculaciones[i].NumeroHorasSemanales,
			NumeroSemanas:        previnculaciones[i].NumeroSemanas,
			Categoria:            strings.Trim(previnculaciones[i].Categoria, " "),
			Dedicacion:           previnculaciones[i].ResolucionVinculacionDocenteId.Dedicacion,
			ValorContratoFormato: FormatMoney(int(previnculaciones[i].ValorContrato), 2),
			NumeroContrato:       previnculaciones[i].NumeroContrato,
			Vigencia:             previnculaciones[i].Vigencia,
			ProyectoCurricularId: previnculaciones[i].ProyectoCurricularId,
		}
		vinculaciones = append(vinculaciones, *vinculacion)
	}

	if vinculaciones == nil {
		vinculaciones = []models.Vinculaciones{}
	}

	return vinculaciones, outputError
}

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
