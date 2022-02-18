package helpers

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"path/filepath"
	"strconv"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/phpdave11/gofpdf"
	"github.com/udistrital/resoluciones_mid_v2/models"
)

func GenerarResolucion(resolucionId int) (encodedPdf string, outputError map[string]interface{}) {
	/* defer func() {
		if err := recover(); err != nil {
			outputError = map[string]interface{}{"funcion": "GenerarResolucion", "err": err, "status": "500"}
			panic(outputError)
		}
	}() */
	var pdf *gofpdf.Fpdf
	var err3 map[string]interface{}

	if contenidoResolucion, err := CargarResolucionCompleta(resolucionId); err != nil {
		panic(err)
	} else {
		if vinculaciones, err2 := ListarVinculaciones(strconv.Itoa(resolucionId)); err2 != nil {
			panic(err2)
		} else {
			if pdf, err3 = ConstruirDocumentoResolucion(contenidoResolucion, vinculaciones); err3 != nil {
				panic(err3)
			}
			if pdf.Err() {
				logs.Error(pdf.Error())
				panic(pdf.Error())
			}
			if pdf.Ok() {
				encodedPdf = encodePDF(pdf)
			}
		}
	}
	return
}

func GenerarInformeVinculaciones(vinculaciones []models.Vinculaciones) (encodedPdf string, outputError map[string]interface{}) {
	defer func() {
		if err := recover(); err != nil {
			outputError = map[string]interface{}{"funcion": "GenerarResolucion", "err": err, "status": "500"}
			panic(outputError)
		}
	}()
	var err map[string]interface{}

	fontPath := filepath.Join(beego.AppConfig.String("StaticPath"), "fonts")
	fontSize := 12.0
	lineHeight := 5.0

	pdf := gofpdf.New("L", "mm", "A4", fontPath)
	pdf.AddUTF8Font("Calibri", "", "calibri.ttf")
	pdf.SetMargins(20, 20, 20)
	pdf.AddPage()
	pdf.SetFont("Calibri", "", fontSize)

	pdf, err = ConstruirTablaVinculaciones(pdf, vinculaciones, lineHeight, fontSize, "RVIN")
	if err != nil {
		panic(err)
	}
	if pdf.Err() {
		logs.Error(pdf.Error())
		panic(pdf.Error())
	}
	if pdf.Ok() {
		encodedPdf = encodePDF(pdf)
	}
	return
}

func ConstruirDocumentoResolucion(datos models.ContenidoResolucion, vinculaciones []models.Vinculaciones) (doc *gofpdf.Fpdf, outputError map[string]interface{}) {
	fontPath := filepath.Join(beego.AppConfig.String("StaticPath"), "fonts")
	imgPath := filepath.Join(beego.AppConfig.String("StaticPath"), "img")
	fontSize := 11.0
	lineHeight := 4.0
	fecha := datos.Resolucion.FechaExpedicion
	fechaParsed := fmt.Sprintf("(%s %02d de %d)", TranslateMonth(fecha.Month().String()), fecha.Day(), fecha.Year())

	var tipoResolucion models.Parametro
	if err := GetRequestNew("UrlcrudParametros", "parametro/"+strconv.Itoa(datos.Resolucion.TipoResolucionId), &tipoResolucion); err != nil {
		outputError = map[string]interface{}{"funcion": "/ConstruirDocumentoResolucion-param", "err": err.Error(), "status": "500"}
		return doc, outputError
	}

	/*
		 * Consultar ordenadores del gasto por medio de Terceros MID y filtrar usando dependenciaId
		var ordenadoresGasto []map[string]interface{}
		if errr := GetRequestLegacy("UrlmidTerceros", "tipo/ordenadoresGasto", &ordenadoresGasto); errr != nil {
			outputError = map[string]interface{}{"funcion": "/ConstruirDocumentoResolucion-ter", "err": errr.Error(), "status": "500"}
			return doc, outputError
		}
	*/

	// Consultar supervisor de contrato de Argo con la dependencia homologada
	var ordenadorGasto models.SupervisorContrato
	var ordenadoresGasto []models.SupervisorContrato
	if dependenciaId, errr := HomologarFacultad("new", strconv.Itoa(datos.Resolucion.DependenciaId)); errr != nil {
		return doc, errr
	} else {
		url := "supervisor_contrato?limit=1&sortby=Id&order=desc&query=DependenciaSupervisor:DEP" + dependenciaId
		if nErr := GetRequestLegacy("UrlcrudAgora", url, &ordenadoresGasto); nErr != nil {

		} else {
			ordenadorGasto = ordenadoresGasto[0]
		}
	}

	pdf := gofpdf.New("P", "mm", "A4", fontPath)
	pdf.AddUTF8Font("Calibri", "", "calibri.ttf")
	pdf.AddUTF8Font("Calibri-Bold", "B", "calibrib.ttf")
	pdf.AddUTF8Font("MinionPro-BoldCn", "B", "MinionPro-BoldCn.ttf")
	pdf.AddUTF8Font("MinionPro-MediumCn", "", "MinionPro-MediumCn.ttf")
	pdf.AddUTF8Font("MinionProBoldItalic", "BI", "MinionProBoldItalic.ttf")

	pdf.SetTopMargin(85)

	pdf.SetHeaderFuncMode(func() {

		pdf.SetLeftMargin(10)
		pdf.SetRightMargin(10)

		pdf.ImageOptions(filepath.Join(imgPath, "escudo.png"), 82, 8, 45, 45, false, gofpdf.ImageOptions{ImageType: "PNG", ReadDpi: true}, 0, "")
		pdf.SetY(55)
		pdf.SetFont("MinionPro-BoldCn", "B", fontSize)
		pdf.WriteAligned(0, lineHeight+1, "RESOLUCIÓN Nº "+datos.Resolucion.NumeroResolucion, "C")
		pdf.Ln(lineHeight)
		pdf.WriteAligned(0, lineHeight+1, fechaParsed, "C")
		pdf.Ln(lineHeight * 2)

		pdf.SetFont("MinionProBoldItalic", "BI", fontSize)
		pdf.WriteAligned(0, lineHeight+1, datos.Resolucion.Titulo, "C")
		pdf.Ln(lineHeight * 2)
	}, true)

	pdf.SetFooterFunc(func() {
		pdf.SetY(-15)
		pdf.SetFont("Calibri", "", 8)
		pdf.WriteAligned(0, lineHeight-1, fmt.Sprintf("Página %d de {nb}", pdf.PageNo()), "R")
	})

	pdf.SetAcceptPageBreakFunc(func() bool {
		y := pdf.GetY()
		_, h := pdf.GetPageSize()
		_, _, _, b := pdf.GetMargins()
		p := pdf.PageNo()
		fmt.Println(p, y, h-b-lineHeight, h, b, (lineHeight * 2))
		return y >= h-b-(lineHeight*2)
	})

	pdf.AliasNbPages("")
	pdf.AddPage()

	pdf.SetAutoPageBreak(false, 25)

	pdf.SetLeftMargin(20)
	pdf.SetRightMargin(20)

	pdf.Ln(lineHeight)

	pdf.SetFont("Calibri", "", fontSize)
	pdf.WriteAligned(0, lineHeight, datos.Resolucion.PreambuloResolucion, "L")
	pdf.Ln(lineHeight * 2)

	pdf.SetFont("Calibri-Bold", "B", fontSize)
	pdf.WriteAligned(0, lineHeight, "CONSIDERANDO", "C")
	pdf.Ln(lineHeight * 2)

	pdf.SetFont("Calibri", "", fontSize)
	pdf.Write(lineHeight, datos.Resolucion.ConsideracionResolucion)
	pdf.Ln(lineHeight * 2)

	pdf.SetFont("Calibri-Bold", "B", fontSize)
	pdf.WriteAligned(0, lineHeight, "RESUELVE", "C")
	pdf.Ln(lineHeight * 2)

	for _, articulo := range datos.Articulos {

		pdf.SetLeftMargin(20)
		pdf.SetRightMargin(20)

		pdf.SetFont("Calibri-Bold", "B", fontSize)
		pdf.Write(lineHeight, fmt.Sprintf("ARTÍCULO %dº. ", articulo.Articulo.Numero))

		pdf.SetFont("Calibri", "", fontSize)
		pdf.Write(lineHeight, articulo.Articulo.Texto)
		pdf.Ln(lineHeight)

		if articulo.Articulo.Numero == 1 {
			pdf.SetLeftMargin(10)
			pdf.SetRightMargin(10)

			pdf, outputError = ConstruirTablaVinculaciones(pdf, vinculaciones, lineHeight, fontSize, tipoResolucion.CodigoAbreviacion)
			if outputError != nil {
				return pdf, outputError
			}

			pdf.SetLeftMargin(20)
			pdf.SetRightMargin(20)
		}

		for _, paragrafo := range articulo.Paragrafos {

			pdf.SetFont("Calibri-Bold", "B", fontSize)
			pdf.Write(lineHeight, "PARÁGRAFO. ")

			pdf.SetFont("Calibri", "", fontSize)
			pdf.Write(lineHeight, paragrafo.Texto)
			pdf.Ln(lineHeight)
		}
	}

	pdf.Ln(lineHeight)
	pdf.SetFont("Calibri-Bold", "B", fontSize)
	pdf.WriteAligned(0, lineHeight, "COMUNÍQUESE Y CÚMPLASE", "C")
	pdf.Ln(lineHeight * 2)

	pdf.SetFont("Calibri", "", fontSize)
	pdf.Write(lineHeight, fmt.Sprintf("Dado en Bogotá D.C., a los %d dias del mes de %s de %d", fecha.Day(), TranslateMonth(fecha.Month().String()), fecha.Year()))
	_, h := pdf.GetPageSize()
	_, _, _, b := pdf.GetMargins()
	if pdf.GetY() > h-b-(lineHeight*10) {
		pdf.AddPage()
	}
	pdf.Ln(lineHeight * 10)

	pdf.SetFont("MinionPro-MediumCn", "", fontSize)
	pdf.WriteAligned(0, lineHeight, ordenadorGasto.Nombre, "C")
	pdf.Ln(lineHeight)
	pdf.WriteAligned(0, lineHeight, ordenadorGasto.Cargo, "C")
	pdf.Ln(lineHeight * 2)

	var cuadroResponsabilidades []map[string]interface{}
	if len(datos.Resolucion.CuadroResponsabilidades) > 0 {
		if err := json.Unmarshal([]byte(datos.Resolucion.CuadroResponsabilidades), &cuadroResponsabilidades); err != nil {
			outputError = map[string]interface{}{"funcion": "ConstruirDocumentoResolucion", "err": err.Error(), "status": "500"}
			return nil, outputError
		}
	} else {
		cuadroResponsabilidades = make([]map[string]interface{}, 4)
	}

	pdf = ConstruirCuadroResp(pdf, cuadroResponsabilidades, true)

	return pdf, nil
}

func ConstruirTablaVinculaciones(pdf *gofpdf.Fpdf, vinculaciones []models.Vinculaciones, lineHeight, fontSize float64, tipoRes string) (doc *gofpdf.Fpdf, outputError map[string]interface{}) {
	var proyectoCurricular models.Dependencia
	w := 18.0
	minHeight := 3.0 * lineHeight
	if tipoRes == "RVIN" {
		w = 20.0
		minHeight = 2.0 * lineHeight
	}
	for _, vinc := range vinculaciones {
		maxHeight := lineHeight
		if proyectoCurricular.Id != vinc.ProyectoCurricularId {
			url := "dependencia/" + strconv.Itoa(int(vinc.ProyectoCurricularId))
			if err2 := GetRequestLegacy("UrlcrudOikos", url, &proyectoCurricular); err2 != nil {
				outputError = map[string]interface{}{"funcion": "/ConstruirTablaVinculaciones-dep", "err": err2.Error(), "status": "500"}
				return doc, outputError
			}
			pdf.Ln(lineHeight * 2)
			pdf.SetFont("Calibri", "", fontSize)
			pdf.Write(lineHeight, proyectoCurricular.Nombre)
			pdf.SetFont("Calibri", "", fontSize-3)
			pdf.Ln(lineHeight * 2)

			pdf.CellFormat(w+4, lineHeight*2, "Nombre", "1", 0, "C", false, 0, "")
			pdf.CellFormat(w+2, lineHeight*2, "Tipo Documento", "1", 0, "C", false, 0, "")
			pdf.CellFormat(w, lineHeight*2, "Cédula", "1", 0, "C", false, 0, "")
			pdf.CellFormat(w, lineHeight*2, "Expedida", "1", 0, "C", false, 0, "")
			pdf.CellFormat(w, lineHeight*2, "Categoría", "1", 0, "C", false, 0, "")
			pdf.CellFormat(w-1, lineHeight*2, "Dedicación", "1", 0, "C", false, 0, "")
			x, y := pdf.GetXY()
			pdf.MultiCell(w-3, lineHeight, "Horas semanales", "1", "C", false)
			if pdf.GetY()-y > lineHeight {
				pdf.SetXY(x+w-3, y)
			}
			x, y = pdf.GetXY()
			pdf.MultiCell(w, lineHeight, "Periodo de vinculación", "1", "C", false)
			if pdf.GetY()-y > lineHeight {
				pdf.SetXY(x+w, y)
			}
			pdf.CellFormat(w+1, lineHeight*2, "Valor total", "1", 0, "C", false, 0, "")
			if tipoRes == "RCAN" {
				pdf.CellFormat(w+1, lineHeight*2, "Valor a reversar", "1", 0, "C", false, 0, "")
			}
			if tipoRes == "RVIN" || tipoRes == "RADD" {
				pdf.CellFormat(7, lineHeight*2, "CPD", "1", 0, "C", false, 0, "")
			} else {
				pdf.CellFormat(7, lineHeight*2, "CRP", "1", 0, "C", false, 0, "")
			}
			pdf.Ln(-1)
		}

		splitText := pdf.SplitLines([]byte(vinc.Nombre), w+4)
		cellHeight := float64(len(splitText)) * lineHeight
		if cellHeight < minHeight {
			cellHeight = minHeight
			splitText = append(splitText, []byte(""))
		}

		if cellHeight > maxHeight {
			maxHeight = cellHeight
		}
		_, h := pdf.GetPageSize()
		_, _, _, b := pdf.GetMargins()
		if pdf.GetY() > h-b-(cellHeight) {
			pdf.AddPage()
		}
		x, y := pdf.GetXY()
		for i := range splitText {
			border := "LR"
			switch i {
			case 0:
				border = border + "T"
				break
			case len(splitText) - 1:
				border = border + "B"
			}
			pdf.MultiCell(w+4, lineHeight, string(splitText[i]), border, "C", false)

		}
		if pdf.GetY()-y > lineHeight {
			pdf.SetXY(x+w+4, y)
		}
		pdf.CellFormat(w+2, cellHeight, "", "1", 0, "C", false, 0, "")
		pdf.CellFormat(w, cellHeight, fmt.Sprintf("%.f", vinc.PersonaId), "1", 0, "C", false, 0, "")
		pdf.CellFormat(w, cellHeight, "", "1", 0, "C", false, 0, "")
		pdf.CellFormat(w, cellHeight, vinc.Categoria, "1", 0, "C", false, 0, "")
		pdf.CellFormat(w-1, cellHeight, vinc.Dedicacion, "1", 0, "C", false, 0, "")
		pdf.CellFormat(w-3, cellHeight, strconv.Itoa(vinc.NumeroHorasSemanales), "1", 0, "C", false, 0, "")
		switch tipoRes {
		case "RVIN":
			pdf.CellFormat(w, cellHeight, "", "1", 0, "C", false, 0, "")
			pdf.CellFormat(w+1, cellHeight, vinc.ValorContratoFormato, "1", 0, "C", false, 0, "")
			break
		case "RCAN":
			x, y = pdf.GetXY()
			pdf.MultiCell(w, lineHeight, "x meses", "1", "C", false)
			pdf.SetX(x)
			pdf.MultiCell(w, lineHeight, "Pasa a", "TLR", "C", false)
			pdf.SetX(x)
			pdf.MultiCell(w, lineHeight, "x meses", "BLR", "C", false)
			if pdf.GetY()-y > lineHeight {
				pdf.SetXY(x+w, y)
			}
			x, y = pdf.GetXY()
			pdf.MultiCell(w+1, lineHeight, vinc.ValorContratoFormato, "1", "C", false)
			pdf.SetX(x)
			pdf.MultiCell(w+1, lineHeight, "Pasa a", "TLR", "C", false)
			pdf.SetX(x)
			pdf.MultiCell(w+1, lineHeight, "$nuevo valor", "BLR", "C", false)
			if pdf.GetY()-y > lineHeight {
				pdf.SetXY(x+w+1, y)
			}
			pdf.CellFormat(w+1, cellHeight, "contrato-nuevo", "1", 0, "C", false, 0, "")
			break
		default:
			x, y = pdf.GetXY()
			pdf.MultiCell(w, lineHeight, "x meses", "1", "C", false)
			pdf.SetX(x)
			pdf.MultiCell(w, lineHeight, "horas +/-", "1", "C", false)
			pdf.SetX(x)
			pdf.MultiCell(w, lineHeight, "horas tt", "1", "C", false)
			if pdf.GetY()-y > lineHeight {
				pdf.SetXY(x+w, y)
			}
			x, y = pdf.GetXY()
			pdf.MultiCell(w+1, lineHeight, vinc.ValorContratoFormato, "1", "C", false)
			pdf.SetX(x)
			pdf.MultiCell(w+1, lineHeight, "meses new", "1", "C", false)
			pdf.SetX(x)
			pdf.MultiCell(w+1, lineHeight, "", "1", "C", false)
			if pdf.GetY()-y > lineHeight {
				pdf.SetXY(x+w+1, y)
			}
			break
		}

		pdf.CellFormat(7, cellHeight, strconv.Itoa(vinc.Disponibilidad), "1", 0, "C", false, 0, "")
		pdf.Ln(-1)
	}
	pdf.Ln(lineHeight)
	return pdf, outputError
}

func ConstruirCuadroResp(pdf *gofpdf.Fpdf, data []map[string]interface{}, resp bool) *gofpdf.Fpdf {

	headers := []string{"Funcion", "Nombre", "Cargo", "Firma"}

	pdf.SetFont("Calibri", "", 6)
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

	pdf.SetFont("Calibri", "", 6)
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

	return pdf
}

func encodePDF(pdf *gofpdf.Fpdf) string {
	var buffer bytes.Buffer
	writer := bufio.NewWriter(&buffer)
	//pdf.OutputFileAndClose("resolucion.pdf") // para guardar el archivo localmente
	pdf.Output(writer)
	writer.Flush()
	encodedFile := base64.StdEncoding.EncodeToString(buffer.Bytes())
	return encodedFile
}

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