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
		writeErrorJSON(&c.Controller, http.StatusBadRequest, "No se pudo leer el archivo", map[string]interface{}{"Error": err.Error()})
		return
	}
	defer file.Close()

	vigenciaStr := c.GetString("vigenciaRp")
	if vigenciaStr == "" {
		writeErrorJSON(&c.Controller, http.StatusBadRequest, "Debe indicar la vigencia del RP (vigenciaRp)", nil)
		return
	}

	vigenciaRp, err := strconv.Atoi(vigenciaStr)
	if err != nil {
		writeErrorJSON(&c.Controller, http.StatusBadRequest, "El valor de 'vigenciaRp' debe ser un número entero", map[string]interface{}{"Error": err.Error()})
		return
	}

	resultados, err := services.ProcesarVinculaciones(file, header, vigenciaRp)
	if err != nil {
		writeErrorJSON(&c.Controller, http.StatusBadRequest, "Error al procesar el archivo", map[string]interface{}{"Error": err.Error()})
		return
	}

	writeJSON(&c.Controller, http.StatusOK, "Procesadas "+strconv.Itoa(len(resultados))+" filas", resultados, nil)
}
