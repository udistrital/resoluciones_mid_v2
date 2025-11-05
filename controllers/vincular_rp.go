package controllers

import (
	"net/http"
	"strconv"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/udistrital/resoluciones_mid_v2/services"
)

type VincularRpController struct {
	beego.Controller
}

func (c *VincularRpController) URLMapping() {
	c.Mapping("Post", c.Post)
}

// @Title Procesar archivo RP
// @Description Procesa un archivo Excel con CRP y actualiza vinculaciones docentes
// @Param	file	formData	file	true	"Archivo Excel (.xlsx)"
// @Param	vigenciaRp	formData	int	true	"Año de vigencia del RP"
// @Success 200 {object} map[string]interface{}
// @Failure 400 archivo inválido o error al procesar
// @router / [post]
func (c *VincularRpController) Post() {
	logs.Info("Inicio del endpoint /v1/vinculacion_rp/ [POST]")

	file, header, err := c.GetFile("file")
	if err != nil {
		logs.Error("Error al obtener archivo del request:", err)
		c.Data["json"] = map[string]interface{}{
			"Success": false,
			"Message": "No se pudo leer el archivo",
			"Error":   err.Error(),
		}
		c.ServeJSON()
		return
	}
	defer file.Close()

	vigenciaStr := c.GetString("vigenciaRp")
	if vigenciaStr == "" {
		c.Data["json"] = map[string]interface{}{
			"Success": false,
			"Message": "Debe indicar la vigencia del RP (vigenciaRp)",
		}
		c.ServeJSON()
		return
	}

	vigenciaRp, err := strconv.Atoi(vigenciaStr)
	if err != nil {
		c.Data["json"] = map[string]interface{}{
			"Success": false,
			"Message": "El valor de 'vigenciaRp' debe ser un número entero",
			"Error":   err.Error(),
		}
		c.ServeJSON()
		return
	}

	resultados, err := services.ProcesarVinculaciones(file, header, vigenciaRp)
	if err != nil {
		c.Data["json"] = map[string]interface{}{
			"Success": false,
			"Message": "Error al procesar el archivo",
			"Error":   err.Error(),
		}
		c.ServeJSON()
		return
	}

	response := map[string]interface{}{
		"Success": true,
		"Message": "Procesadas " + strconv.Itoa(len(resultados)) + " filas",
		"Data":    resultados,
	}

	c.Ctx.Output.SetStatus(http.StatusOK)
	c.Data["json"] = response
	c.ServeJSON()
}
