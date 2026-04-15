package controllers

import "github.com/udistrital/resoluciones_mid_v2/models"

type testHookQueryParamReader struct {
	values map[string]string
}

func (r testHookQueryParamReader) GetString(key string, defaults ...string) string {
	if value, ok := r.values[key]; ok {
		return value
	}
	if len(defaults) > 0 {
		return defaults[0]
	}
	return ""
}

func TestHookBuildFiltroConsulta(values map[string]string, vigenciaKeys ...string) models.Filtro {
	return buildFiltroConsulta(testHookQueryParamReader{values: values}, vigenciaKeys...)
}

func TestHookValidateFiltroConsulta(filtro models.Filtro) error {
	return validateFiltroConsulta(filtro)
}

func TestHookValidateRequiredText(value string, message string) error {
	return validateRequiredText(value, message)
}

func TestHookValidateRequiredRoles(roles []string) error {
	return validateRequiredRoles(roles)
}

func TestHookValidateNamedPositiveInt(value string, field string) error {
	return validateNamedPositiveInt(value, field)
}
