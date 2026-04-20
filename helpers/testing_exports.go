package helpers

import "github.com/udistrital/resoluciones_mid_v2/models"

func TestHookPlantillaNotificacionResolucion(tipoResolucion string) string {
	return plantillaNotificacionResolucion(tipoResolucion)
}

func TestHookConstruirResumenVinculacion(previnculacion models.VinculacionDocente, persona models.InformacionPersonaNatural, ciudadExpedicion string, proyectoCurricular models.Dependencia, disponibilidad models.DisponibilidadVinculacion) models.Vinculaciones {
	return construirResumenVinculacion(previnculacion, persona, ciudadExpedicion, proyectoCurricular, disponibilidad)
}
