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

	mapaTitan := construirMapaContratosTitan(contratosTitan)

	resultado := make([]models.EstadoSemaforoVinculacion, 0, len(vinculaciones))

	total := 0
	totalConRp := 0
	completas := 0
	pendientesTitan := 0
	sinRp := 0

	for _, vinculacion := range vinculaciones {
		estado := clasificarEstadoSemaforoVinculacion(vinculacion, mapaTitan)
		numeroDocumento := estado.NumeroDocumento

		if numeroDocumentoFiltro != nil && *numeroDocumentoFiltro > 0 && numeroDocumento != *numeroDocumentoFiltro {
			continue
		}

		total++
		if estado.TieneRpResoluciones {
			totalConRp++
		}

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
