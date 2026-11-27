package helpers

import "github.com/udistrital/resoluciones_mid_v2/models"

func TestHookPlantillaNotificacionResolucion(tipoResolucion string) string {
	return plantillaNotificacionResolucion(tipoResolucion)
}

func TestHookConstruirResumenVinculacion(previnculacion models.VinculacionDocente, persona models.InformacionPersonaNatural, ciudadExpedicion string, proyectoCurricular models.Dependencia, disponibilidad models.DisponibilidadVinculacion) models.Vinculaciones {
	return construirResumenVinculacion(previnculacion, persona, ciudadExpedicion, proyectoCurricular, disponibilidad)
}

func TestHookConstruirIdUsuariosPrevinculacion(docentes []models.CargaLectiva) []string {
	return construirIdUsuariosPrevinculacion(docentes)
}

func TestHookConstruirMotivosNoVinculable(requisito models.RequisitoVinculacion) []string {
	return construirMotivosNoVinculable(requisito)
}

func TestHookConstruirDocentesNoVinculables(docentes []models.CargaLectiva, respuesta []models.RequisitoVinculacion) []models.DocenteNoVinculable {
	return construirDocentesNoVinculables(docentes, respuesta)
}

func TestHookNormalizarRespuestaRequisitosPrevinculacion(raw interface{}) ([]models.RequisitoVinculacion, map[string]interface{}) {
	return normalizarRespuestaRequisitosPrevinculacion(raw)
}
