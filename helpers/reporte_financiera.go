package helpers

import (
	"fmt"
	"strconv"

	"github.com/astaxie/beego/logs"
	"github.com/udistrital/resoluciones_mid_v2/models"
)

func ReporteFinanciera(reporte models.DatosReporte) (reporteFinal []models.ReporteFinancieraFinal2, outputError map[string]interface{}) {

	defer func() {
		if err := recover(); err != nil {
			outputError = map[string]interface{}{"funcion": "/ReporteFinanciera", "err": err, "status": "500"}
			panic(outputError)
		}
	}()

	var resp []models.ReporteFinanciera2
	var facultad models.Dependencia
	var proyectoCurricular models.Dependencia

	url := "dependencia/" + strconv.Itoa(reporte.Facultad)
	if err2 := GetRequestLegacy("UrlcrudOikos", url, &facultad); err2 != nil {
		outputError = map[string]interface{}{"funcion": "/Obtención proyecto curricular reporte", "err": err2.Error(), "status": "500"}
		panic(outputError)
	}

	if err := SendRequestNew("UrlCrudResoluciones", "reporte_financiera/all", "POST", &resp, &reporte); err != nil {
		logs.Error(err)
		panic("reporte_financiera -> " + err.Error())
	}

	for i := 0; i < len(resp); i++ {
		var infoDocente models.ObjetoDocenteTg
		//var aux interface{}

		url := "dependencia/" + strconv.Itoa(resp[i].Proyectocurricular)
		if err2 := GetRequestLegacy("UrlcrudOikos", url, &proyectoCurricular); err2 != nil {
			outputError = map[string]interface{}{"funcion": "/Obtención proyecto curricular reporte", "err": err2.Error(), "status": "500"}
			panic(outputError)
		}

		url = fmt.Sprintf("docente/%d", resp[i].DocumentoDocente)
		if err2 := GetRequestWSO2("NscrudAcademica", url, &infoDocente); err2 != nil {
			panic(err2.Error())
		}
		var reporteAux models.ReporteFinancieraFinal2
		reporteAux.Id = resp[i].Id
		reporteAux.Resolucion = resp[i].Resolucion
		reporteAux.DocumentoDocente = resp[i].DocumentoDocente
		reporteAux.Horas = resp[i].Horas
		reporteAux.Semanas = resp[i].Semanas
		reporteAux.Total = resp[i].Total
		reporteAux.Cdp = resp[i].Cdp
		reporteAux.Rp = resp[i].Rp
		reporteAux.Vigencia = resp[i].Vigencia
		reporteAux.Periodo = resp[i].Periodo
		reporteAux.NivelAcademico = resp[i].NivelAcademico
		reporteAux.TipoVinculacion = resp[i].TipoVinculacion
		reporteAux.Sueldobasico = resp[i].Sueldobasico
		reporteAux.Primanavidad = resp[i].Primanavidad
		reporteAux.Vacaciones = resp[i].Vacaciones
		reporteAux.Primavacaciones = resp[i].Primavacaciones
		reporteAux.Cesantias = resp[i].Cesantias
		reporteAux.TipoResolucion = resp[i].TipoResolucion
		reporteAux.Interesescesantias = resp[i].Interesescesantias
		reporteAux.Primaservicios = resp[i].Primaservicios
		reporteAux.Bonificacionservicios = resp[i].Bonificacionservicios
		reporteAux.Nombre = infoDocente.DocenteTg.Docente[0].Nombre
		reporteAux.ProyectoCurricular = proyectoCurricular.Nombre
		reporteAux.CodigoProyecto = resp[i].Proyectocurricular
		reporteAux.Facultad = facultad.Nombre
		reporteAux.CodigoFacultad = reporte.Facultad
		reporteFinal = append(reporteFinal, reporteAux)
	}
	return reporteFinal, nil
}
