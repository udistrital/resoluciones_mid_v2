package services

import (
	"fmt"
	"strings"

	"github.com/udistrital/resoluciones_mid_v2/models"
)

var rolePriority = map[string]int{
	"ADMINISTRADOR_RESOLUCIONES": 3,
	"ASIS_FINANCIERA":            2,
	"DECANO":                     1,
	"ASISTENTE_DECANATURA":       1,
}

func normalizeRol(rol string) string {
	return strings.ToUpper(strings.TrimSpace(rol))
}

func deduplicateDependencias(items []models.DependenciaUsuario) []models.DependenciaUsuario {
	seen := make(map[string]bool)
	result := make([]models.DependenciaUsuario, 0)

	for _, item := range items {
		key := fmt.Sprintf("%d-%d", item.CodigoDependencia, item.IdOikos)
		if !seen[key] {
			seen[key] = true
			result = append(result, item)
		}
	}

	return result
}

func getHighestPriorityRol(roles []string) string {
	bestRol := ""
	bestPriority := -1

	for _, rol := range roles {
		r := normalizeRol(rol)
		priority, ok := rolePriority[r]
		if !ok {
			continue
		}

		if priority > bestPriority {
			bestPriority = priority
			bestRol = r
		}
	}

	return bestRol
}

func isGlobalRol(rol string) bool {
	switch normalizeRol(rol) {
	case "ADMINISTRADOR_RESOLUCIONES", "ASIS_FINANCIERA":
		return true
	default:
		return false
	}
}

func DependenciaPermitida(idOikos int, dependencias []models.DependenciaUsuario) bool {
	for _, dep := range dependencias {
		if dep.IdOikos == idOikos {
			return true
		}
	}
	return false
}
