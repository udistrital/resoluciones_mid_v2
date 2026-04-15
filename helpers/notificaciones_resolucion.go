package helpers

import "github.com/udistrital/resoluciones_mid_v2/models"

func NotificarDocentes(datosCorreo []models.EmailData, tipoResolucion string) (outputError map[string]interface{}) {
	var emailRes models.EmailResponse

	if updatedDatosCorreos, err := ObtenerCorreoDocentes(datosCorreo); err != nil {
		outputError = map[string]interface{}{"funcion": "/NotificarDocentes", "err": err, "status": "400"}
		return outputError
	} else {
		emailBody := models.TemplatedEmail{
			Source:       "notificacion_resoluciones@udistrital.edu.co",
			Template:     plantillaNotificacionResolucion(tipoResolucion),
			Destinations: construirDestinosCorreo(updatedDatosCorreos),
			DefaultTemplateData: models.TemplateData{
				Facultad:         "",
				NumeroContrato:   "",
				NumeroResolucion: "",
			},
		}
		url := "email/enviar_templated_email"
		if err := SendRequestNew("UrlMidNotificaciones", url, "POST", &emailRes, emailBody); err != nil {
			outputError = map[string]interface{}{"funcion": "/NotificarDocentes", "err": err.Error(), "status": "400"}
		}
	}
	return outputError
}

func construirDestinosCorreo(datosCorreo []models.EmailData) []models.Destinations {
	destinations := make([]models.Destinations, 0, len(datosCorreo))
	for _, datoCorreo := range datosCorreo {
		destinations = append(destinations, models.Destinations{
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
		})
	}

	return destinations
}

func plantillaNotificacionResolucion(tipoResolucion string) string {
	switch tipoResolucion {
	case "RVIN":
		return "RESOLUCIONES_VINCULACION_PLANTILLA"
	case "RCAN":
		return "RESOLUCIONES_CANCELACION_PLANTILLA"
	case "RRED":
		return "RESOLUCIONES_REDUCCION_PLANTILLA"
	case "RADD":
		return "RESOLUCIONES_ADICION_PLANTILLA"
	default:
		return ""
	}
}

func obtenerCorreoDocentePorDocumento(numeroDocumento string) (string, error) {
	var response struct {
		Email string `json:"email"`
	}
	type Documento struct {
		Numero string `json:"numero"`
	}

	body := Documento{Numero: numeroDocumento}
	if err := SendRequestLegacy("UrlMidAutenticacion", "token/documentoToken", "POST", &response, &body); err != nil {
		return "", err
	}

	return response.Email, nil
}

func ObtenerCorreoDocentes(datosCorreo []models.EmailData) (updatedDatosCorreo []models.EmailData, outputError map[string]interface{}) {
	for i, datos := range datosCorreo {
		if correo, err := obtenerCorreoDocentePorDocumento(datos.Documento); err == nil {
			if correo != "" {
				datosCorreo[i].Correo = correo
			}
		} else {
			outputError = map[string]interface{}{"funcion": "/ObtenerCorreoDocentes", "err": err.Error(), "status": "404"}
		}
	}
	updatedDatosCorreo = datosCorreo
	return updatedDatosCorreo, outputError
}
