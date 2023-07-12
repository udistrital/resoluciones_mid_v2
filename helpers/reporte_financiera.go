package helpers

import (
	"fmt"
	"strconv"

	"github.com/astaxie/beego/logs"
	"github.com/udistrital/resoluciones_mid_v2/models"
)

func ReporteFinanciera(reporte models.DatosReporte) (reporteFinal []models.ReporteFinancieraFinal, outputError map[string]interface{}) {

	defer func() {
		if err := recover(); err != nil {
			outputError = map[string]interface{}{"funcion": "/ReporteFinanciera", "err": err, "status": "500"}
			panic(outputError)
		}
	}()

	var resp []models.ReporteFinanciera
	var facultad models.Dependencia
	var proyectoCurricular models.Dependencia

	url := "dependencia/" + strconv.Itoa(reporte.Facultad)
	if err2 := GetRequestLegacy("UrlcrudOikos", url, &facultad); err2 != nil {
		outputError = map[string]interface{}{"funcion": "/Obtención proyecto curricular reporte", "err": err2.Error(), "status": "500"}
		panic(outputError)
	}

	if err := SendRequestNew("UrlCrudResoluciones", "reporte_financiera", "POST", &resp, &reporte); err != nil {
		logs.Error(err)
		panic("reporte_financiera -> " + err.Error())
	}

	for i := 0; i < len(resp); i++ {
		var infoDocente models.ObjetoDocenteTg
		//var aux interface{}

		url := "dependencia/" + strconv.Itoa(resp[i].ProyectoCurricular)
		if err2 := GetRequestLegacy("UrlcrudOikos", url, &proyectoCurricular); err2 != nil {
			outputError = map[string]interface{}{"funcion": "/Obtención proyecto curricular reporte", "err": err2.Error(), "status": "500"}
			panic(outputError)
		}

		url = fmt.Sprintf("docente_tg/%d", resp[i].Cedula)
		if err2 := GetRequestWSO2("NscrudAcademica", url, &infoDocente); err2 != nil {
			panic(err2.Error())
		}
		var reporteAux models.ReporteFinancieraFinal
		reporteAux.Id = resp[i].Id
		reporteAux.Resolucion = resp[i].Resolucion
		reporteAux.Cedula = resp[i].Cedula
		reporteAux.Horas = resp[i].Horas
		reporteAux.Semanas = resp[i].Semanas
		reporteAux.Total = resp[i].Total
		reporteAux.Cdp = resp[i].Cdp
		reporteAux.SueldoBasico = resp[i].SueldoBasico
		reporteAux.PrimaNavidad = resp[i].PrimaNavidad
		reporteAux.Vacaciones = resp[i].Vacaciones
		reporteAux.PrimaVacaciones = resp[i].PrimaVacaciones
		reporteAux.Cesantias = resp[i].Cesantias
		reporteAux.InteresesCesantias = resp[i].InteresesCesantias
		reporteAux.PrimaServicios = resp[i].PrimaServicios
		reporteAux.BonificacionServicios = resp[i].BonificacionServicios
		reporteAux.Nombre = infoDocente.DocenteTg.Docente[0].Nombre
		reporteAux.ProyectoCurricular = proyectoCurricular.Nombre
		reporteAux.CodigoProyecto = resp[i].ProyectoCurricular
		reporteAux.Facultad = facultad.Nombre
		reporteFinal = append(reporteFinal, reporteAux)
	}
	return reporteFinal, nil
}
