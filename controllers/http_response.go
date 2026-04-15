package controllers

import (
	"github.com/astaxie/beego"
	"github.com/udistrital/utils_oas/requestresponse"
)

func writeJSON(c *beego.Controller, status int, message string, data interface{}, extras map[string]interface{}) {
	response := requestresponse.APIResponseDTO(true, status, data, message)
	payload := map[string]interface{}{
		"Success": response.Success,
		"Status":  response.Status,
		"Message": response.Message,
		"Data":    response.Data,
	}

	for key, value := range extras {
		payload[key] = value
	}

	c.Ctx.Output.SetStatus(status)
	c.Data["json"] = payload
	c.ServeJSON()
}

func writeErrorJSON(c *beego.Controller, status int, message string, extras map[string]interface{}) {
	response := requestresponse.APIResponseDTO(false, status, nil, message)
	payload := map[string]interface{}{
		"Success": response.Success,
		"Status":  response.Status,
		"Message": response.Message,
		"Data":    response.Data,
	}

	for key, value := range extras {
		payload[key] = value
	}

	c.Ctx.Output.SetStatus(status)
	c.Data["json"] = payload
	c.ServeJSON()
}
