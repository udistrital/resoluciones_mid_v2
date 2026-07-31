package helpers

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"fmt"
	"strconv"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/phpdave11/gofpdf"
	"github.com/udistrital/resoluciones_mid_v2/models"
)

// Genera la tabla del cuadro de responsabilidades que va al final de cada resolución
func ConstruirCuadroResp(pdf *gofpdf.Fpdf, data []map[string]interface{}, resp bool) (p *gofpdf.Fpdf, outputError map[string]interface{}) {
	defer func() {
		if err := recover(); err != nil {
			outputError = map[string]interface{}{"funcion": "ConstruirCuadroResp", "err": err, "status": "500"}
			panic(outputError)
		}
	}()

	headers := []string{"Funcion", "Nombre", "Cargo", "Firma"}

	pdf.SetFont(Calibri, "", 6)
	for i, str := range headers {
		w := 42.0
		if i == 0 {
			w = w / 2
		}
		if i == 1 {
			w = w * 1.5
		}
		pdf.CellFormat(w, 4, str, "1", 0, "C", false, 0, "")
	}
	pdf.Ln(-1)

	pdf.SetFont(Calibri, "", 6)
	for _, fila := range data {
		for i, str := range headers {
			w := 42.0
			if i == 0 {
				w = w / 2
			}
			if i == 1 {
				w = w * 1.5
			}
			if _, ok := fila[str]; ok {
				pdf.CellFormat(w, 4, fila[str].(string), "1", 0, "C", false, 0, "")
			} else {
				pdf.CellFormat(w, 4, "", "1", 0, "C", false, 0, "")
			}
		}
		pdf.Ln(-1)
	}

	return pdf, outputError
}

// Codifica el documento pdf en formato Base64
func encodePDF(pdf *gofpdf.Fpdf) string {
	var buffer bytes.Buffer
	writer := bufio.NewWriter(&buffer)
	pdf.Output(writer)
	writer.Flush()
	return base64.StdEncoding.EncodeToString(buffer.Bytes())
}

// Para un mes en inglés retorna el nombre del mes en español
func TranslateMonth(engMonth string) (espMonth string) {
	meses := map[string]string{
		"January":   "Enero",
		"February":  "Febrero",
		"March":     "Marzo",
		"April":     "Abril",
		"May":       "Mayo",
		"June":      "Junio",
		"July":      "Julio",
		"August":    "Agosto",
		"September": "Septiembre",
		"October":   "Octubre",
		"November":  "Noviembre",
		"December":  "Diciembre",
	}
	espMonth, _ = meses[engMonth]
	return
}

// Realiza el proceso de almacenar la resolución a traves del gestor documental
func AlmacenarResolucionGestorDocumental(resolucionId int) (documento models.Documento, outputError map[string]interface{}) {
	var doc models.DocumentoContainer
	if documentoGenerado, err := GenerarResolucion(resolucionId); err == nil {
		data := make([]map[string]interface{}, 0)
		data = append(data, map[string]interface{}{
			"IdTipoDocumento": 22,
			"file":            documentoGenerado,
			"nombre":          "ResolucionDVE" + strconv.Itoa(resolucionId),
			"descripcion":     "Resolución de vinculación especial",
			"metadatos":       map[string]interface{}{},
		})
		url := joinConfiguredURL(beego.AppConfig.String("UrlGestorDocumental"), "document/upload")
		if err := SendJson(url, "POST", &doc, data); err != nil {
			logs.Error(err.Error())
			outputError = map[string]interface{}{"funcion": "/AlmacenarResolucionGestorDocumental ", "err": err.Error(), "status": "500"}
		}
	} else {
		logs.Error(err)
		outputError = map[string]interface{}{"funcion": "/AlmacenarResolucionGestorDocumental ", "err": err, "status": "500"}
	}
	if doc.Status != "200" {
		outputError = map[string]interface{}{"funcion": "/AlmacenarResolucionGestorDocumental ", "err": doc.Error, "status": doc.Status}
	}
	return doc.Res, outputError
}

// RecuperarDocumentoResolucion genera y almacena el PDF de una resolución que
// ya fue expedida, y persiste el enlace devuelto por el gestor documental sin
// modificar su estado ni volver a ejecutar el flujo de expedición.
func RecuperarDocumentoResolucion(resolucionId int) (uid string, outputError map[string]interface{}) {
	var resolucion models.Resolucion
	url := ResolucionEndpoint + strconv.Itoa(resolucionId)

	if err := GetRequestNew("UrlCrudResoluciones", url, &resolucion); err != nil {
		return "", map[string]interface{}{
			"funcion": "/RecuperarDocumentoResolucion",
			"err":     err.Error(),
			"status":  "500",
		}
	}

	// Permite reintentar el endpoint de forma segura cuando el UID ya fue guardado.
	if resolucion.NuxeoUid != "" {
		return resolucion.NuxeoUid, nil
	}

	documento, errMap := AlmacenarResolucionGestorDocumental(resolucionId)
	if errMap != nil {
		return "", errMap
	}
	if documento.Enlace == "" {
		return "", map[string]interface{}{
			"funcion": "/RecuperarDocumentoResolucion",
			"err":     "el gestor documental no devolvió el UID del documento",
			"status":  "502",
		}
	}

	resolucion.NuxeoUid = documento.Enlace
	var response interface{}
	if err := SendRequestNew("UrlCrudResoluciones", url, "PUT", &response, &resolucion); err != nil {
		return documento.Enlace, map[string]interface{}{
			"funcion": "/RecuperarDocumentoResolucion",
			"err":     fmt.Sprintf("el documento fue creado con UID %s, pero no fue posible guardarlo en la resolución: %v", documento.Enlace, err),
			"status":  "500",
			"uid":     documento.Enlace,
		}
	}

	return documento.Enlace, nil
}
