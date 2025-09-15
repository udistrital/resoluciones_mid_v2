package helpers

import (
	"fmt"
	"strconv"
	"time"

	"github.com/udistrital/resoluciones_mid_v2/models"
)

func SupervisorActual(dependenciaId int) (supervisorActual models.SupervisorContrato, outputError map[string]interface{}) {
	var j []models.JefeDependencia
	var s []models.SupervisorContrato
	var fecha = time.Now().Format("2006-01-02") // -- Se debe dejar este una vez se suba
	// var fecha = "2021-01-01"
	//If Jefe_dependencia (GET)
	url := "jefe_dependencia?query=DependenciaId:" + strconv.Itoa(dependenciaId) + ",FechaFin__gte:" + fecha + ",FechaInicio__lte:" + fecha
	if err := GetRequestLegacy("UrlcrudCore", url, &j); err == nil && len(j) > 0 {
		//If Supervisor (GET)
		url = "supervisor_contrato?order=desc&sortby=Id&query=Documento:" + strconv.Itoa(j[0].TerceroId) + ",FechaFin__gte:" + fecha + ",FechaInicio__lte:" + fecha + "&CargoId.Cargo__startswith:DECANO|VICE"
		if err := GetRequestLegacy("UrlcrudAgora", url, &s); err == nil && len(s) > 0 {
			fmt.Println(s)
			return s[0], nil
		} else { //If Jefe_dependencia (GET)
			fmt.Println("No se ha encontrado supervisor activo en la fecha actual!!!", err)
			outputError = map[string]interface{}{"funcion": "/SupervisorActual3", "err": err.Error(), "status": "404"}
			return supervisorActual, outputError
		}
	} else { //If Jefe_dependencia (GET)
		fmt.Println("No se ha encontrado jefe de dependencia activo en la fecha actual!!! ", err)
		outputError = map[string]interface{}{"funcion": "/SupervisorActual2", "err": err.Error(), "status": "404"}
		return supervisorActual, outputError
	}

}

func calcularSemanasContratoDVE(FechaInicio time.Time, FechaFin time.Time) (semanas float64) {
	var a, m, d int
	var mesesContrato float64
	if FechaFin.IsZero() {
		FechaFin2 := time.Now()
		a, m, d = diff(FechaInicio, FechaFin2)
		mesesContrato = (float64(a * 12)) + float64(m) + (float64(d) / 30)

	} else {
		a, m, d = diff(FechaInicio, FechaFin)
		fmt.Println("a ", a)
		fmt.Println("m ", m)
		// dia inclusivo
		d += 1
		fmt.Println("d ", d)
		if d == 22 {
			d += 1
		}
		mesesContrato = (float64(a * 12)) + float64(m) + (float64(d) / 30)
	}
	fmt.Println(float64(int(mesesContrato)))
	if mesesContrato/float64(int(mesesContrato)) != 1 {
		return (mesesContrato * 4) + 1
	} else {
		return (mesesContrato * 4)
	}
}

// CalcularFechasContrato calcula la fecha de fin de un contrato a partir de la fecha de inicio y el numero de semanas
// Calcula las fechas reales del contrato y las fechas ajustadas para TITAN
// Ajusta la fecha inicio y fin para que inicie un lunes
func CalcularFechasContrato(fechaInicio time.Time, numeroSemanas int) models.FechasContrato {
	var resultado models.FechasContrato

	// Guardar fechas reales (calendario normal)
	resultado.FechaInicioReal = fechaInicio
	dias := numeroSemanas * 7
	resultado.FechaFinReal = fechaInicio.AddDate(0, 0, dias)
	resultado.SemanasReales = float64(dias) / 7

	// Inicializar fechas de pago
	resultado.FechaInicioPago = fechaInicio

	// Si la fecha no es lunes, ajustar al próximo lunes
	if resultado.FechaInicioPago.Weekday() != time.Monday {
		// Calcular días hasta el próximo lunes
		diasHastaLunes := int(time.Monday - resultado.FechaInicioPago.Weekday())
		if diasHastaLunes <= 0 {
			diasHastaLunes += 7
		}
		resultado.FechaInicioPago = resultado.FechaInicioPago.AddDate(0, 0, diasHastaLunes)
	}

	// Ajustar fecha inicio si cae en 31
	// if resultado.FechaInicioPago.Day() == 31 {
	// 	resultado.FechaInicioPago = resultado.FechaInicioPago.AddDate(0, 0, 1)
	// }

	// Calcular la fecha fin de pago basada en el número exacto de semanas
	diasPago := numeroSemanas * 7
	if diasPago == 0 {
		resultado.FechaFinPago = resultado.FechaInicioPago
	} else {
		resultado.FechaFinPago = resultado.FechaInicioPago.AddDate(0, 0, diasPago-1) // -1 porque el día inicial cuenta
	}

	// Si cae en 31, ajustar al 30
	if resultado.FechaFinPago.Day() == 31 {
		resultado.FechaFinPago = resultado.FechaFinPago.AddDate(0, 0, -1)
	}

	// Calcular las semanas reales (que serán iguales al número de semanas solicitado)
	diasDiferencia := resultado.FechaFinPago.Sub(resultado.FechaInicioPago).Hours() / 24
	resultado.SemanasPagoReales = (diasDiferencia + 1) / 7 // +1 para incluir el día inicial

	// Las semanas DVE serán iguales a las reales porque las fechas están alineadas
	resultado.SemanasPagoDve = calcularSemanasContratoDVE(resultado.FechaInicioPago, resultado.FechaFinPago)

	return resultado
}

func NotificarDocentes(datosCorreo []models.EmailData, tipoResolucion string) (outputError map[string]interface{}) {

	var emailRes models.EmailResponse

	if updatedDatosCorreos, err := ObtenerCorreoDocentes(datosCorreo); err != nil {
		fmt.Println("No se ha podido obtener los correos de los docentes", err)
		outputError = map[string]interface{}{"funcion": "/NotificarDocentes", "err": err, "status": "400"}
		return outputError
	} else {
		var destinationsArray = []models.Destinations{}
		for _, datoCorreo := range updatedDatosCorreos {
			destinations := models.Destinations{
				Destination: models.Destination{
					ToAddresses:  []string{datoCorreo.Correo},
					BccAddresses: nil,
					CcAddresses:  nil,
				},
				ReplacementTemplateData: models.TemplateData{
					Facultad:         datoCorreo.Facultad,
					NumeroContrato:   datoCorreo.ContratoId,
					NumeroResolucion: datoCorreo.NumeroResolucion,
				},
				Attachments: []models.Attachments{},
			}
			destinationsArray = append(destinationsArray, destinations)
		}
		emailBody := models.TemplatedEmail{
			Source:       "notificacion_resoluciones@udistrital.edu.co",
			Template:     "",
			Destinations: destinationsArray,
			DefaultTemplateData: models.TemplateData{
				Facultad:         "",
				NumeroContrato:   "",
				NumeroResolucion: "",
			},
		}
		if tipoResolucion == "RVIN" {
			emailBody.Template = "RESOLUCIONES_VINCULACION_PLANTILLA"
		} else if tipoResolucion == "RCAN" {
			emailBody.Template = "RESOLUCIONES_CANCELACION_PLANTILLA"
		} else if tipoResolucion == "RRED" {
			emailBody.Template = "RESOLUCIONES_REDUCCION_PLANTILLA"
		} else if tipoResolucion == "RADD" {
			emailBody.Template = "RESOLUCIONES_ADICION_PLANTILLA"
		}
		url := "email/enviar_templated_email"
		if err := SendRequestNew("UrlMidNotificaciones", url, "POST", &emailRes, emailBody); err != nil {
			fmt.Println("No se ha podido enviar el correo a los docentes ", err)
			outputError = map[string]interface{}{"funcion": "/NotificarDocentes", "err": err.Error(), "status": "400"}
		}
	}
	fmt.Println("outputError", outputError)
	return outputError
}

func ObtenerCorreoDocentes(datosCorreo []models.EmailData) (updatedDatosCorreo []models.EmailData, outputError map[string]interface{}) {
	var response struct {
		Email string `json:"email"`
	}
	type Documento struct {
		Numero string `json:"numero"`
	}
	var docSrtct Documento
	for i, datos := range datosCorreo {
		docSrtct = Documento{Numero: datos.Documento}
		url := "token/documentoToken"
		if err := SendRequestLegacy("UrlMidAutenticacion", url, "POST", &response, &docSrtct); err == nil {
			if response.Email != "" {
				datosCorreo[i].Correo = response.Email
			}
		} else {
			fmt.Println("No se ha encontrado información del usuario", err)
			outputError = map[string]interface{}{"funcion": "/ObtenerCorreoDocentes", "err": err.Error(), "status": "404"}
		}
	}
	updatedDatosCorreo = datosCorreo
	return updatedDatosCorreo, outputError
}
