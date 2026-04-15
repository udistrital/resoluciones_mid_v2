package services

import (
	"strings"

	"github.com/udistrital/resoluciones_mid_v2/models"
)

type resumenEstadoRp struct {
	NumeroContrato     string
	TieneRpResoluciones bool
}

func construirMapaContratosTitan(contratos []models.ContratoTitan) map[string]bool {
	resultado := make(map[string]bool)

	for _, contrato := range contratos {
		numeroContrato := strings.TrimSpace(contrato.NumeroContrato)
		if numeroContrato != "" && contrato.Activo {
			resultado[numeroContrato] = true
		}
	}

	return resultado
}

func construirMapaConteoContratosTitan(contratos []models.ContratoTitan) map[string]int {
	resultado := make(map[string]int)

	for _, contrato := range contratos {
		numeroContrato := strings.TrimSpace(contrato.NumeroContrato)
		if numeroContrato != "" && contrato.Activo {
			resultado[numeroContrato]++
		}
	}

	return resultado
}

func resumirEstadoRp(vinculacion models.VinculacionDocente) resumenEstadoRp {
	numeroContrato := ""
	if vinculacion.NumeroContrato != nil {
		numeroContrato = strings.TrimSpace(*vinculacion.NumeroContrato)
	}

	return resumenEstadoRp{
		NumeroContrato: numeroContrato,
		TieneRpResoluciones: vinculacion.NumeroRp > 0 &&
			vinculacion.VigenciaRp > 0 &&
			numeroContrato != "",
	}
}

func clasificarEstadoSemaforoVinculacion(vinculacion models.VinculacionDocente, titanPorContrato map[string]bool) models.EstadoSemaforoVinculacion {
	estado := models.EstadoSemaforoVinculacion{
		VinculacionId:       vinculacion.Id,
		NumeroDocumento:     int(vinculacion.PersonaId),
		Vigencia:            vinculacion.Vigencia,
		NumeroRp:            int(vinculacion.NumeroRp),
		VigenciaRp:          int(vinculacion.VigenciaRp),
		TieneRpResoluciones: false,
		TieneRpTitan:        false,
		EstadoCodigo:        "SIN_RP",
		EstadoNombre:        "Sin RP cargado en resoluciones",
		Prioridad:           3,
	}

	resumenRp := resumirEstadoRp(vinculacion)
	estado.NumeroContrato = resumenRp.NumeroContrato

	if !resumenRp.TieneRpResoluciones {
		return estado
	}

	estado.TieneRpResoluciones = true
	estado.EstadoCodigo = "PENDIENTE_TITAN"
	estado.EstadoNombre = "Cargado en resoluciones, pendiente en Titan"
	estado.Prioridad = 2

	if titanPorContrato[estado.NumeroContrato] {
		estado.TieneRpTitan = true
		estado.EstadoCodigo = "COMPLETO"
		estado.EstadoNombre = "Cargado en resoluciones y Titan"
		estado.Prioridad = 1
	}

	return estado
}
