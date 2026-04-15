package controllers

import (
	"strings"

	"github.com/astaxie/beego"
	"github.com/udistrital/resoluciones_mid_v2/models"
)

type authContextReader interface {
	GetString(string, ...string) string
}

func buildAuthenticatedContextFromRequest(reader authContextReader) models.AuthenticatedContext {
	numeroDocumento := strings.TrimSpace(reader.GetString("numero_documento"))
	roles := parseRolesParam(reader.GetString("roles"))

	return models.AuthenticatedContext{
		NumeroDocumento: numeroDocumento,
		Roles:           roles,
		Source:          "query",
	}
}

func requireAuthenticatedContext(ctx models.AuthenticatedContext, function string) models.AuthenticatedContext {
	if err := validateRequiredText(ctx.NumeroDocumento, "numero_documento es requerido"); err != nil {
		panic(badRequest(function, err))
	}

	if err := validateRequiredRoles(ctx.Roles); err != nil {
		panic(badRequest(function, err))
	}

	return ctx
}

func buildAuthenticatedContext(controller *beego.Controller) models.AuthenticatedContext {
	return buildAuthenticatedContextFromRequest(controller)
}

