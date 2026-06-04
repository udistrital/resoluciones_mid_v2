package helpers

import "github.com/udistrital/resoluciones_mid_v2/models"

func TestHookPlantillaNotificacionResolucion(tipoResolucion string) string {
	return plantillaNotificacionResolucion(tipoResolucion)
}

func TestHookConstruirResumenVinculacion(previnculacion models.VinculacionDocente, persona models.InformacionPersonaNatural, ciudadExpedicion string, proyectoCurricular models.Dependencia, disponibilidad models.DisponibilidadVinculacion) models.Vinculaciones {
	return construirResumenVinculacion(previnculacion, persona, ciudadExpedicion, proyectoCurricular, disponibilidad)
}

func TestHookConstruirIdUsuariosOdin(docentes []models.CargaLectiva) string {
	return construirIdUsuariosOdin(docentes)
}

func TestHookConstruirMotivosNoVinculable(requisito models.OdinRequisitoVinculacion) []string {
	return construirMotivosNoVinculable(requisito)
}

func TestHookConstruirDocentesNoVinculables(docentes []models.CargaLectiva, respuesta []models.OdinRequisitoVinculacion) []models.DocenteNoVinculable {
	return construirDocentesNoVinculables(docentes, respuesta)
}

func TestHookNormalizarRespuestaRequisitosODIN(raw interface{}) ([]models.OdinRequisitoVinculacion, map[string]interface{}) {
	return normalizarRespuestaRequisitosODIN(raw)
}
