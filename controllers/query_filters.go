package controllers

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/udistrital/resoluciones_mid_v2/helpers"
	"github.com/udistrital/resoluciones_mid_v2/models"
)

type queryParamReader interface {
	GetString(string, ...string) string
}

func parseRolesParam(raw string) []string {
	if strings.TrimSpace(raw) == "" {
		return []string{}
	}

	parts := strings.Split(raw, ",")
	roles := make([]string, 0, len(parts))

	for _, part := range parts {
		rol := strings.ToUpper(strings.TrimSpace(part))
		if rol != "" {
			roles = append(roles, rol)
		}
	}

	return roles
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return strings.TrimSpace(value)
		}
	}

	return ""
}

func buildFiltroConsulta(reader queryParamReader, vigenciaKeys ...string) models.Filtro {
	return models.Filtro{
		Limit:            strings.TrimSpace(reader.GetString("limit")),
		Offset:           strings.TrimSpace(reader.GetString("offset")),
		NumeroResolucion: firstNonEmpty(reader.GetString("NumeroResolucion"), reader.GetString("numero_resolucion")),
		Vigencia:         firstNonEmpty(getStrings(reader, vigenciaKeys...)...),
		Periodo:          strings.TrimSpace(reader.GetString("Periodo")),
		Semanas:          strings.TrimSpace(reader.GetString("Semanas")),
		FacultadId:       firstNonEmpty(reader.GetString("Facultad"), reader.GetString("id_oikos")),
		NivelAcademico:   strings.TrimSpace(reader.GetString("NivelAcademico")),
		Dedicacion:       strings.TrimSpace(reader.GetString("Dedicacion")),
		Estado:           strings.TrimSpace(reader.GetString("Estado")),
		TipoResolucion:   strings.TrimSpace(reader.GetString("TipoResolucion")),
		ExcluirTipo:      strings.TrimSpace(reader.GetString("ExcluirTipo")),
	}
}

func getStrings(reader queryParamReader, keys ...string) []string {
	values := make([]string, 0, len(keys))
	for _, key := range keys {
		values = append(values, reader.GetString(key))
	}
	return values
}

func validateFiltroConsulta(f models.Filtro) error {
	if err := validateRequiredPositiveInt(f.Limit); err != nil {
		return err
	}
	if err := validateRequiredPositiveInt(f.Offset); err != nil {
		return err
	}
	if err := validateOptionalInt(f.NumeroResolucion); err != nil {
		return err
	}
	if err := validateOptionalInt(f.Vigencia); err != nil {
		return err
	}
	if err := validateOptionalInt(f.Periodo); err != nil {
		return err
	}
	if err := validateOptionalInt(f.Semanas); err != nil {
		return err
	}
	return nil
}

func validateRequiredText(value string, message string) error {
	if strings.TrimSpace(value) == "" {
		return &controllerValidationError{message: message}
	}
	return nil
}

func validateRequiredRoles(roles []string) error {
	if len(roles) == 0 {
		return &controllerValidationError{message: "roles es requerido y debe contener al menos un rol"}
	}
	return nil
}

func validateRequiredPositiveInt(value string) error {
	number, err := strconv.Atoi(value)
	if err != nil || number <= 0 {
		return errInvalidParams()
	}
	return nil
}

func validateNamedPositiveInt(value string, field string) error {
	number, err := strconv.Atoi(value)
	if err != nil || number <= 0 {
		return &controllerValidationError{message: fmt.Sprintf("%s es requerido y debe ser válido", field)}
	}
	return nil
}

func validateOptionalInt(value string) error {
	if strings.TrimSpace(value) == "" {
		return nil
	}
	if _, err := strconv.Atoi(value); err != nil {
		return errInvalidParams()
	}
	return nil
}

func errInvalidParams() error {
	return &controllerValidationError{message: helpers.ErrorParametros}
}

func badRequest(function string, err error) map[string]interface{} {
	return map[string]interface{}{
		"funcion": function,
		"err":     err.Error(),
		"status":  "400",
	}
}

func parseRolesRequired(raw string, function string) []string {
	roles := parseRolesParam(raw)
	if err := validateRequiredRoles(roles); err != nil {
		panic(badRequest(function, err))
	}
	return roles
}

func parseRequiredPositiveInt(value string, field string, function string) int {
	if err := validateNamedPositiveInt(value, field); err != nil {
		panic(badRequest(function, err))
	}

	number, _ := strconv.Atoi(value)
	return number
}

func parseOptionalPositiveIntPointer(value string, field string, function string) *int {
	if strings.TrimSpace(value) == "" {
		return nil
	}

	if err := validateNamedPositiveInt(value, field); err != nil {
		panic(badRequest(function, err))
	}

	number, _ := strconv.Atoi(value)
	return &number
}

func parseOptionalNonNegativeInt(value string, field string, function string, defaultValue int) int {
	if strings.TrimSpace(value) == "" {
		return defaultValue
	}

	number, err := strconv.Atoi(value)
	if err != nil || number < 0 {
		panic(badRequest(function, &controllerValidationError{message: fmt.Sprintf("El parámetro %s no es válido", field)}))
	}

	return number
}

type controllerValidationError struct {
	message string
}

func (e *controllerValidationError) Error() string {
	return e.message
}
