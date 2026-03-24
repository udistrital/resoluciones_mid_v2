package services

import (
	"strconv"

	"github.com/udistrital/resoluciones_mid_v2/helpers"
	"github.com/udistrital/resoluciones_mid_v2/models"
)

func ConsultarSemaforoResolucion(resolucionId int, numeroDocUsuario string, roles []string, numeroDocumentoFiltro *int) (*models.RespuestaSemaforoResolucion, map[string]interface{}) {
	permitido, errPermiso := UsuarioPuedeConsultarResolucion(resolucionId, numeroDocUsuario, roles)
	if errPermiso != nil {
		return nil, errPermiso
	}
	if !permitido {
		return nil, map[string]interface{}{
			"funcion": "ConsultarSemaforoResolucion",
			"err":     "No autorizado para consultar esta resolución",
			"status":  "403",
		}
	}

	vinculaciones, errVin := helpers.Previnculaciones(strconv.Itoa(resolucionId))
	if errVin != nil {
		return nil, errVin
	}

	contratosTitan, errTitan := helpers.ObtenerContratosTitanPorResolucion(resolucionId)
	if errTitan != nil {
		return nil, errTitan
	}

	mapaTitan := make(map[string]bool)
	for _, contrato := range contratosTitan {
		if contrato.NumeroContrato != "" && contrato.Activo {
			mapaTitan[contrato.NumeroContrato] = true
		}
	}

	resultado := make([]models.EstadoSemaforoVinculacion, 0, len(vinculaciones))

	total := 0
	totalConRp := 0
	completas := 0
	pendientesTitan := 0
	sinRp := 0

	for _, vinculacion := range vinculaciones {
		numeroDocumento := int(vinculacion.PersonaId)

		if numeroDocumentoFiltro != nil && *numeroDocumentoFiltro > 0 && numeroDocumento != *numeroDocumentoFiltro {
			continue
		}

		estado := models.EstadoSemaforoVinculacion{
			VinculacionId:       vinculacion.Id,
			NumeroDocumento:     numeroDocumento,
			Vigencia:            vinculacion.Vigencia,
			NumeroRp:            int(vinculacion.NumeroRp),
			VigenciaRp:          int(vinculacion.VigenciaRp),
			TieneRpResoluciones: false,
			TieneRpTitan:        false,
			EstadoCodigo:        "SIN_RP",
			EstadoNombre:        "Sin RP cargado en resoluciones",
			Prioridad:           3,
		}

		if vinculacion.NumeroContrato != nil {
			estado.NumeroContrato = *vinculacion.NumeroContrato
		}

		tieneRpResoluciones := vinculacion.NumeroRp > 0 &&
			vinculacion.VigenciaRp > 0 &&
			estado.NumeroContrato != ""

		if tieneRpResoluciones {
			totalConRp++
			estado.TieneRpResoluciones = true
			estado.EstadoCodigo = "PENDIENTE_TITAN"
			estado.EstadoNombre = "Cargado en resoluciones, pendiente en Titan"
			estado.Prioridad = 2

			if mapaTitan[estado.NumeroContrato] {
				estado.TieneRpTitan = true
				estado.EstadoCodigo = "COMPLETO"
				estado.EstadoNombre = "Cargado en resoluciones y Titan"
				estado.Prioridad = 1
			}
		}

		total++

		switch estado.EstadoCodigo {
		case "COMPLETO":
			completas++
		case "PENDIENTE_TITAN":
			pendientesTitan++
		case "SIN_RP":
			sinRp++
		}

		resultado = append(resultado, estado)
	}

	if numeroDocumentoFiltro != nil && *numeroDocumentoFiltro > 0 && len(resultado) == 0 {
		return nil, map[string]interface{}{
			"funcion": "ConsultarSemaforoResolucion",
			"err":     "No se encontraron vinculaciones para el número de documento indicado en la resolución consultada",
			"status":  "404",
		}
	}

	porcentajeCompletas := 0.0
	porcentajePendientesTitan := 0.0
	porcentajeSinRp := 0.0

	if total > 0 {
		porcentajeCompletas = (float64(completas) / float64(total)) * 100
		porcentajePendientesTitan = (float64(pendientesTitan) / float64(total)) * 100
		porcentajeSinRp = (float64(sinRp) / float64(total)) * 100
	}

	respuesta := &models.RespuestaSemaforoResolucion{
		Resumen: models.ResumenSemaforoResolucion{
			ResolucionId:              resolucionId,
			Total:                     total,
			TotalConRp:                totalConRp,
			Completas:                 completas,
			PendientesTitan:           pendientesTitan,
			SinRp:                     sinRp,
			PorcentajeCompletas:       porcentajeCompletas,
			PorcentajePendientesTitan: porcentajePendientesTitan,
			PorcentajeSinRp:           porcentajeSinRp,
		},
		Detalle: resultado,
	}

	return respuesta, nil
}
