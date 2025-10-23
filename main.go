package main

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
	"github.com/astaxie/beego/plugins/cors"
	_ "github.com/udistrital/resoluciones_mid_v2/routers"
	apistatus "github.com/udistrital/utils_oas/apiStatusLib"
	auditoria "github.com/udistrital/utils_oas/auditoria"
	"github.com/udistrital/utils_oas/customerrorv2"
	"github.com/udistrital/utils_oas/xray"
)

func main() {
	AllowedOrigins := []string{"*.udistrital.edu.co"}
	if beego.BConfig.RunMode == "dev" {
		AllowedOrigins = []string{"*"}
		beego.BConfig.WebConfig.DirectoryIndex = true
		beego.BConfig.WebConfig.StaticDir["/swagger"] = "swagger"
	}
	beego.InsertFilter("*", beego.BeforeRouter, cors.Allow(&cors.Options{
		AllowOrigins: AllowedOrigins,
		AllowMethods: []string{"PUT", "PATCH", "GET", "POST", "OPTIONS", "DELETE"},
		AllowHeaders: []string{"Origin", "x-requested-with",
			"content-type",
			"accept",
			"origin",
			"authorization",
			"x-csrftoken"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))
	xray.InitXRay()
	beego.ErrorController(&customerrorv2.CustomErrorController{})
	beego.InsertFilter("*", beego.BeforeExec, SecurityHeaders)
	apistatus.Init()
	auditoria.InitMiddleware()
	beego.Run()
}

func SecurityHeaders(ctx *context.Context) {
	ctx.Output.Header("Clear-Site-Data", "'cache', 'cookies', 'storage', 'executionContexts'")
	ctx.Output.Header("Cross-Origin-Embedder-Policy", "require-corp")
	ctx.Output.Header("Cross-Origin-Opener-Policy", "same-origin")
	ctx.Output.Header("Cross-Origin-Resource-Policy", "same-origin")
	ctx.Output.Header("Permissions-Policy", "geolocation=(), microphone=(), camera=()")
	ctx.Output.Header("Referrer-Policy", "no-referrer")
	ctx.Output.Header("Server", "")
	ctx.Output.Header("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
	ctx.Output.Header("X-Content-Type-Options", "nosniff")
	ctx.Output.Header("X-Frame-Options", "DENY")
	ctx.Output.Header("X-Permitted-Cross-Domain-Policies", "none")
}
