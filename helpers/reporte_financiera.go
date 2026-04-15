package helpers

import (
	"fmt"
	"strconv"

	"github.com/astaxie/beego/logs"
	"github.com/udistrital/resoluciones_mid_v2/models"
)

func cargarDependenciaReporte(id int, function string) (models.Dependencia, map[string]interface{}) {
	var dependencia models.Dependencia
	url := "dependencia/" + strconv.Itoa(id)
	if err := GetRequestLegacy("UrlcrudOikos", url, &dependencia); err != nil {
		return models.Dependencia{}, map[string]interface{}{
			"funcion": function,
			"err":     err.Error(),
			"status":  "500",
		}
	}
	return dependencia, nil
}

func cargarDocenteReporte(documento int) (models.ObjetoDocenteTg, error) {
	var infoDocente models.ObjetoDocenteTg
	url := fmt.Sprintf("docente/%d", documento)
	if err := GetRequestWSO2("NscrudAcademica", url, &infoDocente); err != nil {
		return models.ObjetoDocenteTg{}, err
	}
	return infoDocente, nil
}

func construirRegistroReporteFinanciera(
	item models.ReporteFinanciera2,
	facultad models.Dependencia,
	proyectoCurricular models.Dependencia,
	infoDocente models.ObjetoDocenteTg,
	codigoFacultad int,
) models.ReporteFinancieraFinal2 {
	return models.ReporteFinancieraFinal2{
		Id:                    item.Id,
		Resolucion:            item.Resolucion,
		DocumentoDocente:      item.DocumentoDocente,
		Horas:                 item.Horas,
		Semanas:               item.Semanas,
		Total:                 item.Total,
		Cdp:                   item.Cdp,
		Rp:                    item.Rp,
		Vigencia:              item.Vigencia,
		Periodo:               item.Periodo,
		NivelAcademico:        item.NivelAcademico,
		TipoVinculacion:       item.TipoVinculacion,
		Sueldobasico:          item.Sueldobasico,
		Primanavidad:          item.Primanavidad,
		Vacaciones:            item.Vacaciones,
		Primavacaciones:       item.Primavacaciones,
		Cesantias:             item.Cesantias,
		TipoResolucion:        item.TipoResolucion,
		Interesescesantias:    item.Interesescesantias,
		Primaservicios:        item.Primaservicios,
		Bonificacionservicios: item.Bonificacionservicios,
		Nombre:                infoDocente.DocenteTg.Docente[0].Nombre,
		ProyectoCurricular:    proyectoCurricular.Nombre,
		CodigoProyecto:        item.Proyectocurricular,
		Facultad:              facultad.Nombre,
		CodigoFacultad:        codigoFacultad,
	}
}

func ReporteFinanciera(reporte models.DatosReporte) (reporteFinal []models.ReporteFinancieraFinal2, outputError map[string]interface{}) {

	defer func() {
		if err := recover(); err != nil {
			outputError = map[string]interface{}{"funcion": "/ReporteFinanciera", "err": err, "status": "500"}
			panic(outputError)
		}
	}()

	var resp []models.ReporteFinanciera2
	facultad, errMap := cargarDependenciaReporte(reporte.Facultad, "/Obtención facultad reporte")
	if errMap != nil {
		outputError = errMap
		panic(outputError)
	}

	if err := SendRequestNew("UrlCrudResoluciones", "reporte_financiera/all", "POST", &resp, &reporte); err != nil {
		logs.Error(err)
		panic("reporte_financiera -> " + err.Error())
	}

	for i := 0; i < len(resp); i++ {
		proyectoCurricular, errMap := cargarDependenciaReporte(resp[i].Proyectocurricular, "/Obtención proyecto curricular reporte")
		if errMap != nil {
			outputError = errMap
			panic(outputError)
		}

		infoDocente, err2 := cargarDocenteReporte(resp[i].DocumentoDocente)
		if err2 != nil {
			panic(err2.Error())
		}

		reporteFinal = append(reporteFinal, construirRegistroReporteFinanciera(
			resp[i],
			facultad,
			proyectoCurricular,
			infoDocente,
			reporte.Facultad,
		))
	}
	return reporteFinal, nil
}
