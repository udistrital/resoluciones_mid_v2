package services

import "github.com/udistrital/resoluciones_mid_v2/models"

func TestHookNormalizeRol(rol string) string {
	return normalizeRol(rol)
}

func TestHookJoinWSO2URL(protocol, base, ns, path string) string {
	return joinWSO2URL(protocol, base, ns, path)
}

func TestHookGetHighestPriorityRol(roles []string) string {
	return getHighestPriorityRol(roles)
}

func TestHookIsGlobalRol(rol string) bool {
	return isGlobalRol(rol)
}

func TestHookDeduplicateDependencias(items []models.DependenciaUsuario) []models.DependenciaUsuario {
	return deduplicateDependencias(items)
}

func TestHookConstruirMapaContratosTitan(contratos []models.ContratoTitan) map[string]bool {
	return construirMapaContratosTitan(contratos)
}

func TestHookResumirEstadoRp(vinculacion models.VinculacionDocente) resumenEstadoRp {
	return resumirEstadoRp(vinculacion)
}

func TestHookClasificarEstadoSemaforoVinculacion(vinculacion models.VinculacionDocente, titanPorContrato map[string]bool) models.EstadoSemaforoVinculacion {
	return clasificarEstadoSemaforoVinculacion(vinculacion, titanPorContrato)
}
