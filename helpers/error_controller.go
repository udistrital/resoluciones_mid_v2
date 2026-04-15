package helpers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/udistrital/utils_oas/requestresponse"
	"github.com/udistrital/utils_oas/xray"
)

// Manejo único de errores para controladores sin repetir código
func ErrorController(c beego.Controller, controller string) {
	if err := recover(); err != nil {
		logs.Error(err)

		statusStr := "500"
		message := "Error interno del servidor"
		errorPath := beego.AppConfig.String("appname") + "/" + controller
		var errorDetail interface{} = err
		var errorMap map[string]interface{}

		if localError, ok := err.(map[string]interface{}); ok {
			errorMap = localError

			if v, ok := localError["status"].(string); ok && v != "" {
				statusStr = v
			}

			if v, ok := localError["err"]; ok {
				message = fmt.Sprintf("%v", v)
				errorDetail = v
			}

			if v, ok := localError["funcion"].(string); ok && v != "" {
				errorPath = beego.AppConfig.String("appname") + "/" + controller + "/" + v
			}
		} else {
			message = fmt.Sprintf("%v", err)
		}

		statusCode, convErr := strconv.Atoi(statusStr)
		if convErr != nil {
			statusCode = http.StatusInternalServerError
		}

		response := requestresponse.APIResponseDTO(false, statusCode, nil, message)
		payload := map[string]interface{}{
			"Success": response.Success,
			"Status":  statusStr,
			"Message": response.Message,
			"Error": map[string]interface{}{
				"path":   errorPath,
				"detail": errorDetail,
			},
		}

		if errorMap != nil {
			if code, ok := errorMap["code"]; ok {
				payload["Code"] = code
			}
			if data, ok := errorMap["data"]; ok {
				payload["Data"] = data
			}
			if detalle, ok := errorMap["detalle"]; ok {
				payload["Detalle"] = detalle
			}
		}

		xray.EndSegmentErr(statusCode, message)
		c.Ctx.Output.SetStatus(statusCode)
		c.Data["json"] = payload
		c.ServeJSON()
	}
}
