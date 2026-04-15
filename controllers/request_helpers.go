package controllers

import (
	"encoding/json"
	"strconv"
	"strings"

	"github.com/astaxie/beego"
	"github.com/udistrital/resoluciones_mid_v2/helpers"
)

func parsePositivePathID(c *beego.Controller, param string, function string) int {
	id, err := strconv.Atoi(c.Ctx.Input.Param(param))
	if err != nil || id <= 0 {
		panic(map[string]interface{}{"funcion": function, "err": helpers.ErrorParametros, "status": "400"})
	}
	return id
}

func decodeJSONBody(body []byte, target interface{}, function string) {
	if valid, err := helpers.ValidarBody(body); !valid || err != nil {
		panic(map[string]interface{}{"funcion": function, "err": helpers.ErrorBody, "status": "400"})
	}

	if err := json.Unmarshal(body, target); err != nil {
		panic(map[string]interface{}{"funcion": function, "err": err.Error(), "status": "400"})
	}
}

func parseYearPathParam(c *beego.Controller, param string, function string) string {
	value := strings.TrimSpace(c.Ctx.Input.Param(param))
	if len(value) != 4 {
		panic(map[string]interface{}{"funcion": function, "err": helpers.ErrorParametros, "status": "400"})
	}
	parsePositiveIntParam(value, function)
	return value
}

func parsePositiveIntParam(value string, function string) int {
	number, err := strconv.Atoi(value)
	if err != nil || number <= 0 {
		panic(map[string]interface{}{"funcion": function, "err": helpers.ErrorParametros, "status": "400"})
	}
	return number
}

func parseSingleDigitPositivePathParam(c *beego.Controller, param string, function string) string {
	value := strings.TrimSpace(c.Ctx.Input.Param(param))
	if len(value) != 1 {
		panic(map[string]interface{}{"funcion": function, "err": helpers.ErrorParametros, "status": "400"})
	}
	parsePositiveIntParam(value, function)
	return value
}
