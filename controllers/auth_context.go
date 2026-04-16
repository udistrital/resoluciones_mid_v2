package controllers

import (
	"strings"

	"github.com/astaxie/beego"
	"github.com/udistrital/resoluciones_mid_v2/models"
)

var trustedNumeroDocumentoHeaders = []string{
	"X-Authenticated-User",
	"X-User-Document",
	"X-User",
}

var trustedRolesHeaders = []string{
	"X-Authenticated-Roles",
	"X-User-Roles",
	"X-Roles",
}

func firstHeaderValue(controller *beego.Controller, headers ...string) string {
	for _, header := range headers {
		if value := strings.TrimSpace(controller.Ctx.Input.Header(header)); value != "" {
			return value
		}
	}
	return ""
}

func buildRequestAuthContext(controller *beego.Controller) models.RequestAuthContext {
	numeroDocumento := firstHeaderValue(controller, trustedNumeroDocumentoHeaders...)
	rolesRaw := firstHeaderValue(controller, trustedRolesHeaders...)
	source := "headers"
	trusted := true

	if numeroDocumento == "" {
		numeroDocumento = strings.TrimSpace(controller.GetString("numero_documento"))
		source = "query_fallback"
		trusted = false
	}

	if strings.TrimSpace(rolesRaw) == "" {
		rolesRaw = controller.GetString("roles")
		if source == "headers" {
			source = "query_fallback"
			trusted = false
		}
	}

	return models.RequestAuthContext{
		NumeroDocumento: numeroDocumento,
		Roles:           parseRolesParam(rolesRaw),
		Source:          source,
		Trusted:         trusted,
	}
}

func requireRequestAuthContext(ctx models.RequestAuthContext, function string) models.RequestAuthContext {
	if err := validateRequiredText(ctx.NumeroDocumento, "numero_documento es requerido"); err != nil {
		panic(badRequest(function, err))
	}

	if err := validateRequiredRoles(ctx.Roles); err != nil {
		panic(badRequest(function, err))
	}

	return ctx
}
