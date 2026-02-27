package helpers

import (
	"fmt"
	"strconv"

	"github.com/udistrital/resoluciones_mid_v2/models"
)

// ValidarLiquidacionesTitan valida que, para cada vinculacion recibida,
// exista en Titan el contrato base (inicial) asociado por NumeroContrato y Vigencia.
// Esta funcion aplica para modificaciones.
func ValidarLiquidacionesTitan(vinculaciones []*models.ContratoVinculacion) (outputError map[string]interface{}) {
	if len(vinculaciones) == 0 {
		return map[string]interface{}{
			"funcion": "/ValidarLiquidacionesTitan",
			"err":     "No se recibieron vinculaciones para validar en Titan.",
			"status":  "400",
		}
	}

	ids, errInfo := extraerIdsVinculaciones(vinculaciones)
	if errInfo != nil {
		return errInfo
	}

	return ValidarLiquidacionesTitanPorIds(ids)
}

// ValidarLiquidacionesTitanPorIds valida en Titan a partir de ids de vinculacion_docente.
// Esta funcion sirve para reutilizar la validacion en flujos como cancelacion,
// donde el request suele traer solo el Id de la vinculacion.
func ValidarLiquidacionesTitanPorIds(vinculacionIds []int) (outputError map[string]interface{}) {
	if len(vinculacionIds) == 0 {
		return map[string]interface{}{
			"funcion": "/ValidarLiquidacionesTitanPorIds",
			"err":     "No se recibieron ids de vinculaciones para validar en Titan.",
			"status":  "400",
		}
	}

	for _, vinculacionId := range vinculacionIds {
		if vinculacionId <= 0 {
			return map[string]interface{}{
				"funcion": "/ValidarLiquidacionesTitanPorIds",
				"err":     "Se recibió un id de vinculación docente inválido.",
				"status":  "400",
			}
		}

		if errInfo := validarVinculacionContraTitan(vinculacionId); errInfo != nil {
			return errInfo
		}
	}

	return nil
}

func extraerIdsVinculaciones(vinculaciones []*models.ContratoVinculacion) (ids []int, outputError map[string]interface{}) {
	ids = make([]int, 0, len(vinculaciones))

	for _, vinculacion := range vinculaciones {
		if vinculacion == nil {
			return nil, map[string]interface{}{
				"funcion": "/extraerIdsVinculaciones",
				"err":     "Se recibió una vinculación nula.",
				"status":  "400",
			}
		}

		if vinculacion.VinculacionDocente.Id <= 0 {
			return nil, map[string]interface{}{
				"funcion": "/extraerIdsVinculaciones",
				"err":     "Se recibió una vinculación docente sin id válido.",
				"status":  "400",
			}
		}

		ids = append(ids, vinculacion.VinculacionDocente.Id)
	}

	return ids, nil
}

func validarVinculacionContraTitan(vinculacionId int) (outputError map[string]interface{}) {
	var vinculacionActual models.VinculacionDocente

	url := VinculacionEndpoint + strconv.Itoa(vinculacionId)
	if err := GetRequestNew("UrlCrudResoluciones", url, &vinculacionActual); err != nil {
		return map[string]interface{}{
			"funcion": "/validarVinculacionContraTitan",
			"err":     "No fue posible consultar la vinculación docente a validar.",
			"detalle": "Vinculación docente: " + strconv.Itoa(vinculacionId),
			"status":  "502",
		}
	}

	numeroContrato, vigencia, errInfo := ObtenerContratoBaseTitan(vinculacionActual)
	if errInfo != nil {
		return errInfo
	}

	return ValidarContratoEnTitan(numeroContrato, vigencia, vinculacionActual.Id)
}

// ObtenerContratoBaseTitan resuelve el contrato base que debe existir en Titan.
// Para novedades (cancelacion/modificacion) debe validar SIEMPRE contra la vinculacion inicial.
func ObtenerContratoBaseTitan(vinculacion models.VinculacionDocente) (numeroContrato string, vigencia int, outputError map[string]interface{}) {
	contratosHistoricos := make([]models.VinculacionDocente, 0)

	if err := BuscarContratosModificar(vinculacion.Id, &contratosHistoricos); err != nil {
		return "", 0, map[string]interface{}{
			"funcion": "/ObtenerContratoBaseTitan",
			"err":     "No fue posible consultar el histórico de la vinculación para validar la liquidación en Titan.",
			"detalle": "Vinculación docente: " + strconv.Itoa(vinculacion.Id),
			"status":  "502",
		}
	}

	var vinculacionBase models.VinculacionDocente

	if len(contratosHistoricos) > 0 {
		// La última del histórico corresponde a la vinculación inicial
		vinculacionBase = contratosHistoricos[len(contratosHistoricos)-1]
	} else {
		// Si no hay histórico, la actual es la inicial
		vinculacionBase = vinculacion
	}

	if vinculacionBase.NumeroContrato == nil || *vinculacionBase.NumeroContrato == "" || vinculacionBase.Vigencia <= 0 {
		return "", 0, map[string]interface{}{
			"funcion": "/ObtenerContratoBaseTitan",
			"err":     "No se pudo identificar el contrato inicial de la vinculación para validar la liquidación en Titan.",
			"detalle": "Vinculación docente: " + strconv.Itoa(vinculacion.Id),
			"status":  "409",
		}
	}

	return *vinculacionBase.NumeroContrato, vinculacionBase.Vigencia, nil
}

func dataTitanTieneContenidoReal(data []interface{}) bool {
	for _, item := range data {
		switch v := item.(type) {
		case map[string]interface{}:
			if len(v) > 0 {
				return true
			}
		case nil:
			// no hace nada
		default:
			return true
		}
	}
	return false
}

// ValidarContratoEnTitan valida que exista al menos un contrato en Titan con el NumeroContrato y Vigencia indicados.
// Regla de negocio: si existe el contrato en Titan, se entiende que sí tiene liquidacion.
func ValidarContratoEnTitan(numeroContrato string, vigencia int, vinculacionDocenteId int) (outputError map[string]interface{}) {
	var response map[string]interface{}

	queryPath := fmt.Sprintf("contrato?query=NumeroContrato:%s,Vigencia:%d", numeroContrato, vigencia)

	if err := GetRequestLegacy("TitanCrudService", queryPath, &response); err != nil {
		return map[string]interface{}{
			"funcion": "/ValidarContratoEnTitan",
			"err":     "No fue posible validar la liquidación en Titan para la vinculación inicial.",
			"detalle": "Contrato inicial: " + numeroContrato + ", vigencia: " + strconv.Itoa(vigencia),
			"status":  "502",
		}
	}

	existeContrato := false

	if data, ok := response["Data"]; ok && data != nil {
		switch valor := data.(type) {
		case []interface{}:
			existeContrato = dataTitanTieneContenidoReal(valor)
		case map[string]interface{}:
			existeContrato = len(valor) > 0
		default:
			existeContrato = false
		}
	}

	if !existeContrato {
		return map[string]interface{}{
			"funcion": "/ValidarContratoEnTitan",
			"err":     "No se puede expedir la novedad porque la vinculación inicial no tiene liquidación registrada en Titan.",
			"detalle": "Contrato inicial: " + numeroContrato + ", vigencia: " + strconv.Itoa(vigencia),
			"status":  "409",
			"code":    "TITAN_LIQUIDACION_NO_ENCONTRADA",
			"data": map[string]interface{}{
				"contrato": numeroContrato,
				"vigencia": vigencia,
			},
		}
	}

	return nil
}
