package helpers

import (
	"bufio"
	"bytes"
	"encoding/base64"
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
		url := beego.AppConfig.String("ProtocolAdmin") + "://" + beego.AppConfig.String("UrlGestorDocumental") + "document/upload"
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
