package helpers

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/phpdave11/gofpdf"
	"github.com/udistrital/resoluciones_mid_v2/models"
)

var meses = map[string][]string{
	"es": {"Enero", "Febrero", "Marzo", "Abril", "Mayo", "Junio", "Julio", "Agosto", "Septiembre", "Octubre", "Noviembre", "Diciembre"},
}

// Función que orquesta el proceso de generación de la resolución en formato pdf
func GenerarResolucion(resolucionId int) (encodedPdf string, outputError map[string]interface{}) {
	defer func() {
		if err := recover(); err != nil {
			outputError = map[string]interface{}{"funcion": "GenerarResolucion", "err": err, "status": "500"}
			panic(outputError)
		}
	}()
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

// Esta función genera un documento en formato pdf con las tablas de las vinculaciones proporcionadas
func GenerarInformeVinculaciones(vinculaciones []models.Vinculaciones) (encodedPdf string, outputError map[string]interface{}) {
	defer func() {
		if err := recover(); err != nil {
			outputError = map[string]interface{}{"funcion": "GenerarInformeVinculaciones", "err": err, "status": "500"}
			panic(outputError)
		}
	}()
	var err map[string]interface{}
	var v []models.VinculacionDocente

	url := "vinculacion_docente?query=Id:" + strconv.Itoa(vinculaciones[0].Id)
	if err := GetRequestNew("UrlcrudResoluciones", url, &v); err != nil {
		logs.Error(err.Error())
		panic(err.Error())
	}

	fontPath := filepath.Join(beego.AppConfig.String("StaticPath"), "fonts")
	fontSize := 12.0
	lineHeight := 4.0

	pdf := gofpdf.New("L", "mm", "A4", fontPath)
	pdf.AddUTF8Font(Calibri, "", "calibri.ttf")
	pdf.SetMargins(20, 20, 20)
	pdf.AddPage()
	pdf.SetFont(Calibri, "", fontSize)

	pdf, err = ConstruirTablaVinculaciones(pdf, vinculaciones, lineHeight, fontSize, "RVIN", v[0].ResolucionVinculacionDocenteId.NivelAcademico)
	if err != nil {
		logs.Error(err)
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

func obtenerNombreMes(m time.Month, idioma string) string {
	// Verifica si el idioma está soportado
	nombresMeses, ok := meses[idioma]
	if !ok {
		return ""
	}

	// Resta 1 al valor del mes porque en Go los meses empiezan desde 1
	mes := int(m) - 1

	// Verifica si el mes está dentro del rango válido
	if mes < 0 || mes >= len(nombresMeses) {
		return ""
	}

	return nombresMeses[mes]
}

// Esta función genera un documento en formato pdf con la información de la resolución registrada en la base de datos
func ConstruirDocumentoResolucion(datos models.ContenidoResolucion, vinculaciones []models.Vinculaciones) (doc *gofpdf.Fpdf, outputError map[string]interface{}) {
	defer func() {
		if err := recover(); err != nil {
			outputError = map[string]interface{}{"funcion": "ConstruirDocumentoResolucion", "err": err, "status": "500"}
			panic(outputError)
		}
	}()

	var fecha time.Time
	fontPath := filepath.Join(beego.AppConfig.String("StaticPath"), "fonts")
	imgPath := filepath.Join(beego.AppConfig.String("StaticPath"), "img")
	fontSize := 11.0
	lineHeight := 4.0
	fecha = datos.Resolucion.FechaExpedicion
	if datos.Resolucion.FechaExpedicion.IsZero() {
		fecha = time.Now()
	}
	fechaParsed := fmt.Sprintf("(%s %02d de %d)", TranslateMonth(fecha.Month().String()), fecha.Day(), fecha.Year())

	var tipoResolucion models.Parametro
	if err := GetRequestNew("UrlcrudParametros", ParametroEndpoint+strconv.Itoa(datos.Resolucion.TipoResolucionId), &tipoResolucion); err != nil {
		panic(map[string]interface{}{"funcion": "/ConstruirDocumentoResolucion-param", "err": err.Error(), "status": "500"})
	}

	/*
		 * Consultar ordenadores del gasto por medio de Terceros MID y filtrar usando dependenciaId
		var ordenadoresGasto []map[string]interface{}
		if errr := GetRequestLegacy("UrlmidTerceros", "tipo/ordenadoresGasto", &ordenadoresGasto); errr != nil {
			outputError = map[string]interface{}{"funcion": "/ConstruirDocumentoResolucion-ter", "err": errr.Error(), "status": "500"}
			panic(outputError)
		}
	*/

	// Consultar ordenador del gasto por core_amazon_crud
	var ordenadorGasto models.OrdenadorGasto
	var ordenadoresGasto []models.OrdenadorGasto
	url := "ordenador_gasto?query=DependenciaId:" + strconv.Itoa(datos.Resolucion.DependenciaFirmaId)
	if err := GetRequestLegacy("UrlcrudCore", url, &ordenadoresGasto); err != nil {
		logs.Error(err)
		panic(err.Error())
	} else {
		if len(ordenadoresGasto) > 0 {
			ordenadorGasto = ordenadoresGasto[0]
		} else {
			if err := GetRequestLegacy("UrlcrudCore", "ordenador_gasto/1", &ordenadorGasto); err != nil {
				logs.Error(err)
				panic(err.Error())
			}
		}
		var jefeDependencia []models.JefeDependencia
		var fechaActual = time.Now().Format("2006-01-02") // -- Se debe dejar este una vez se suba
		// var fechaActual = "2021-01-01"
		url2 := "jefe_dependencia?query=DependenciaId:" + strconv.Itoa(datos.Resolucion.DependenciaFirmaId) + ",FechaFin__gte:" + fechaActual + ",FechaInicio__lte:" + fechaActual
		if err := GetRequestLegacy("UrlcrudCore", url2, &jefeDependencia); err != nil {
			logs.Error(err)
			panic(err.Error())
		}
		if len(jefeDependencia) > 0 {
			if ordenador, err2 := BuscarDatosPersonalesDocente(float64(jefeDependencia[0].TerceroId)); err2 == nil {
				ordenadorGasto.NombreOrdenador = ordenador.NomProveedor
			} else {
				logs.Error(err2)
				panic(err2)
			}
		} else {
			panic("No se encontró jefe para la dependencia en el periodo actual")
		}
	}

	pdf := gofpdf.New("P", "mm", "A4", fontPath)
	pdf.AddUTF8Font(Calibri, "", "calibri.ttf")
	pdf.AddUTF8Font(CalibriBold, "B", "calibrib.ttf")
	pdf.AddUTF8Font(MinionProBoldCn, "B", "MinionPro-BoldCn.ttf")
	pdf.AddUTF8Font(MinionProMediumCn, "", "MinionPro-MediumCn.ttf")
	pdf.AddUTF8Font(MinionProBoldItalic, "BI", "MinionProBoldItalic.ttf")

	pdf.SetTopMargin(85)

	pdf.SetHeaderFuncMode(func() {

		pdf.SetLeftMargin(10)
		pdf.SetRightMargin(10)

		pdf.ImageOptions(filepath.Join(imgPath, "escudo.png"), 82, 8, 45, 45, false, gofpdf.ImageOptions{ImageType: "PNG", ReadDpi: true}, 0, "")
		pdf.SetY(55)
		pdf.SetFont(MinionProBoldCn, "B", fontSize)
		pdf.WriteAligned(0, lineHeight+1, "RESOLUCIÓN Nº "+datos.Resolucion.NumeroResolucion, "C")
		pdf.Ln(lineHeight)
		pdf.WriteAligned(0, lineHeight+1, fechaParsed, "C")
		pdf.Ln(lineHeight * 2)

		pdf.SetFont(MinionProBoldItalic, "BI", fontSize)
		pdf.WriteAligned(0, lineHeight+1, datos.Resolucion.Titulo, "C")
		pdf.Ln(lineHeight * 2)
	}, true)

	pdf.SetFooterFunc(func() {
		pdf.SetY(-15)
		pdf.SetFont(Calibri, "", 8)
		pdf.WriteAligned(0, lineHeight-1, fmt.Sprintf("Página %d de {nb}", pdf.PageNo()), "R")
	})

	pdf.SetAcceptPageBreakFunc(func() bool {
		y := pdf.GetY()
		_, h := pdf.GetPageSize()
		_, _, _, b := pdf.GetMargins()
		// p := pdf.PageNo()
		// fmt.Println(p, y, h-b-lineHeight, h, b, (lineHeight * 2))
		return y >= h-b-(lineHeight*2)
	})

	pdf.AliasNbPages("")
	pdf.AddPage()

	pdf.SetAutoPageBreak(false, 25)

	pdf.SetLeftMargin(20)
	pdf.SetRightMargin(20)

	pdf.Ln(lineHeight)

	pdf.SetFont(Calibri, "", fontSize)
	pdf.MultiCell(0, lineHeight, datos.Resolucion.PreambuloResolucion, "", "J", false)
	pdf.Ln(lineHeight * 2)

	pdf.SetFont(CalibriBold, "B", fontSize)
	pdf.WriteAligned(0, lineHeight, "CONSIDERANDO", "C")
	pdf.Ln(lineHeight * 2)

	pdf.SetFont(Calibri, "", fontSize)
	pdf.MultiCell(0, lineHeight, datos.Resolucion.ConsideracionResolucion, "", "J", false)
	pdf.Ln(lineHeight * 2)

	pdf.SetFont(CalibriBold, "B", fontSize)
	pdf.WriteAligned(0, lineHeight, "RESUELVE", "C")
	pdf.Ln(lineHeight * 2)
	for _, articulo := range datos.Articulos {
		pdf.SetLeftMargin(20)
		pdf.SetRightMargin(20)

		x := pdf.GetX()
		pdf.SetFont(CalibriBold, "B", fontSize)
		tArticulo := fmt.Sprintf("ARTÍCULO %dº.   ", articulo.Articulo.Numero)
		if articulo.Articulo.Numero == 1 && tipoResolucion.CodigoAbreviacion == "RVIN" {
			var inicio *time.Time = datos.Resolucion.FechaInicio
			var diaInicio = inicio.Day()
			var mesInicio = obtenerNombreMes(inicio.Month(), "es")
			var añoInicio = inicio.Year()
			var stringInicio string = strconv.Itoa(diaInicio) + " de " + mesInicio + " del " + strconv.Itoa(añoInicio)
			var fin *time.Time = datos.Resolucion.FechaFin
			var diaFin = fin.Day()
			var mesFin = obtenerNombreMes(fin.Month(), "es")
			var añoFin = fin.Year()
			var stringFin string = strconv.Itoa(diaFin) + " de " + mesFin + " del " + strconv.Itoa(añoFin)
			var replaceI string = "$i"
			var replaceF string = "$f"

			articulo.Articulo.Texto = strings.Replace(articulo.Articulo.Texto, replaceI, stringInicio, 1)
			articulo.Articulo.Texto = strings.Replace(articulo.Articulo.Texto, replaceF, stringFin, 1)
		}
		tArticuloLen := pdf.GetStringWidth(tArticulo)
		pdf.Write(lineHeight, tArticulo)

		pdf.SetX(x)
		pdf.SetFont(Calibri, "", fontSize)
		pdf.MultiCell(0, lineHeight, strings.Repeat(" ", int(tArticuloLen))+articulo.Articulo.Texto, "", "J", false)
		pdf.Ln(lineHeight)

		if articulo.Articulo.Numero == 1 {
			pdf.SetLeftMargin(10)
			pdf.SetRightMargin(10)

			if datos.Vinculacion.Dedicacion != "HCH" {
				pdf, outputError = ConstruirVinculacionesDesagregado(pdf, vinculaciones, lineHeight, fontSize, tipoResolucion.CodigoAbreviacion, datos.Vinculacion.NivelAcademico)
				if outputError != nil {
					logs.Error(outputError)
					panic(outputError)
				}
			} else {
				pdf, outputError = ConstruirTablaVinculaciones(pdf, vinculaciones, lineHeight, fontSize, tipoResolucion.CodigoAbreviacion, datos.Vinculacion.NivelAcademico)
				if outputError != nil {
					logs.Error(outputError)
					panic(outputError)
				}
			}

			pdf.SetLeftMargin(20)
			pdf.SetRightMargin(20)
		}

		for _, paragrafo := range articulo.Paragrafos {

			x := pdf.GetX()
			pdf.SetFont(CalibriBold, "B", fontSize)
			tParagrafo := "PARÁGRAFO.  "
			tParagrafoLen := pdf.GetStringWidth(tParagrafo)
			pdf.Write(lineHeight, tParagrafo)

			pdf.SetX(x)
			pdf.SetFont(Calibri, "", fontSize)
			pdf.MultiCell(0, lineHeight, strings.Repeat(" ", int(tParagrafoLen))+paragrafo.Texto, "", "J", false)
			pdf.Ln(lineHeight)
		}
	}

	pdf.Ln(lineHeight)
	pdf.SetFont(CalibriBold, "B", fontSize)
	pdf.WriteAligned(0, lineHeight, "COMUNÍQUESE Y CÚMPLASE", "C")
	pdf.Ln(lineHeight * 2)

	pdf.SetFont(Calibri, "", fontSize)
	pdf.Write(lineHeight, fmt.Sprintf("Dado en Bogotá D.C., a los %d dias del mes de %s de %d", fecha.Day(), TranslateMonth(fecha.Month().String()), fecha.Year()))
	_, h := pdf.GetPageSize()
	_, _, _, b := pdf.GetMargins()
	if pdf.GetY() > h-b-(lineHeight*10) {
		pdf.AddPage()
	}
	pdf.Ln(lineHeight * 10)

	pdf.SetFont(MinionProBoldCn, "B", fontSize)
	pdf.WriteAligned(0, lineHeight, strings.ToUpper(ordenadorGasto.NombreOrdenador), "C")
	pdf.Ln(lineHeight)
	pdf.SetFont(MinionProMediumCn, "", fontSize)
	pdf.WriteAligned(0, lineHeight, strings.ToUpper(ordenadorGasto.Cargo), "C")
	pdf.Ln(lineHeight * 2)

	var cuadroResponsabilidades []map[string]interface{}
	if len(datos.Resolucion.CuadroResponsabilidades) > 0 {
		if err := json.Unmarshal([]byte(datos.Resolucion.CuadroResponsabilidades), &cuadroResponsabilidades); err != nil {
			logs.Error(err.Error())
			panic(err.Error())
		}
	} else {
		cuadroResponsabilidades = make([]map[string]interface{}, 4)
	}

	var err map[string]interface{}
	if pdf, err = ConstruirCuadroResp(pdf, cuadroResponsabilidades, true); err != nil {
		logs.Error(err)
		panic(err)
	}

	return pdf, outputError
}

// Genera las tablas de las vinculaciones por proyecto curricular de acuerdo al tipo de resolución
func ConstruirTablaVinculaciones(pdf *gofpdf.Fpdf, vinculaciones []models.Vinculaciones, lineHeight, fontSize float64, tipoRes, nivel string) (doc *gofpdf.Fpdf, outputError map[string]interface{}) {
	defer func() {
		if err := recover(); err != nil {
			outputError = map[string]interface{}{"funcion": "ConstruirTablaVinculaciones", "err": err, "status": "500"}
			panic(outputError)
		}
	}()

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
				panic(outputError)
			}
			pdf.Ln(lineHeight * 2)
			pdf.SetFont(Calibri, "", fontSize)
			pdf.Write(lineHeight, proyectoCurricular.Nombre)
			pdf.SetFont(Calibri, "", fontSize-3)
			pdf.Ln(lineHeight * 2)

			// Encabezados
			pdf.CellFormat(w+4, lineHeight*2, "Nombre", "1", 0, "C", false, 0, "")
			pdf.CellFormat(w+2, lineHeight*2, "Tipo Documento", "1", 0, "C", false, 0, "")
			pdf.CellFormat(w, lineHeight*2, "Cédula", "1", 0, "C", false, 0, "")
			pdf.CellFormat(w, lineHeight*2, "Expedida", "1", 0, "C", false, 0, "")
			pdf.CellFormat(w, lineHeight*2, "Categoría", "1", 0, "C", false, 0, "")
			pdf.CellFormat(w-1, lineHeight*2, "Dedicación", "1", 0, "C", false, 0, "")
			x, y := pdf.GetXY()
			if nivel == "PREGRADO" {
				pdf.MultiCell(w-2, lineHeight, "Horas semanales", "1", "C", false)
			} else {
				pdf.MultiCell(w-2, lineHeight, "Horas semestrales", "1", "C", false)
			}
			if pdf.GetY()-y > lineHeight {
				pdf.SetXY(x+w-2, y)
			}
			x, y = pdf.GetXY()
			pdf.MultiCell(w-1, lineHeight, "Periodo de vinculación", "1", "C", false)
			if pdf.GetY()-y > lineHeight {
				pdf.SetXY(x+w-1, y)
			}
			pdf.CellFormat(w+1, lineHeight*2, "Valor total", "1", 0, "C", false, 0, "")
			if tipoRes == "RADD" {
				pdf.CellFormat(w+2, lineHeight*2, "Valor a adicionar", "1", 0, "C", false, 0, "")
			}
			if tipoRes == "RCAN" || tipoRes == "RRED" {
				pdf.CellFormat(w+2, lineHeight*2, "Valor a reversar", "1", 0, "C", false, 0, "")
			}
			if tipoRes == "RVIN" || tipoRes == "RADD" {
				pdf.CellFormat(7, lineHeight*2, "CDP", "1", 0, "C", false, 0, "")
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
		// Nombre
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
		docHeight := 2 * (cellHeight / lineHeight)
		x, y = pdf.GetXY()
		pdf.MultiCell(w+2, docHeight, vinc.TipoDocumento, "1", "C", false)
		if pdf.GetY()-y > docHeight {
			pdf.SetXY(x+w+2, y)
		}
		pdf.CellFormat(w, cellHeight, fmt.Sprintf("%.f", vinc.PersonaId), "1", 0, "C", false, 0, "")
		pdf.CellFormat(w, cellHeight, vinc.ExpedicionDocumento, "1", 0, "C", false, 0, "")
		pdf.CellFormat(w, cellHeight, vinc.Categoria, "1", 0, "C", false, 0, "")
		pdf.CellFormat(w-1, cellHeight, vinc.Dedicacion, "1", 0, "C", false, 0, "")

		valoresAntes := make(map[string]float64)
		if err := CalcularTrazabilidad(strconv.Itoa(vinc.Id), &valoresAntes); err != nil {
			logs.Error("Error en trazabilidad -> " + err.Error())
			panic("Error en trazabilidad -> " + err.Error())
		}
		cellHeight2 := cellHeight - minHeight + lineHeight

		switch tipoRes {
		case "RVIN":
			pdf.CellFormat(w-2, cellHeight, strconv.Itoa(vinc.NumeroHorasSemanales), "1", 0, "C", false, 0, "")
			pdf.CellFormat(w-1, cellHeight, fmt.Sprintf(CampoMeses, vinc.NumeroSemanas), "1", 0, "C", false, 0, "")
			pdf.CellFormat(w+1, cellHeight, vinc.ValorContratoFormato, "1", 0, "C", false, 0, "")
			break
		case "RCAN":
			if nivel == "PREGRADO" {
				pdf.CellFormat(w-2, cellHeight, strconv.Itoa(vinc.NumeroHorasSemanales), "1", 0, "C", false, 0, "")
				x, y = pdf.GetXY()
				pdf.MultiCell(w-1, lineHeight, fmt.Sprintf(CampoMeses, vinc.NumeroSemanas), "1", "C", false)
				pdf.SetX(x)
				pdf.MultiCell(w-1, lineHeight, "Pasa a", "TLR", "C", false)
				pdf.SetX(x)
				pdf.MultiCell(w-1, lineHeight, fmt.Sprintf(CampoMeses, int((valoresAntes["NumeroSemanas"]-float64(vinc.NumeroSemanas)))), "BLR", "C", false)
				if pdf.GetY()-y > lineHeight {
					pdf.SetXY(x+w-1, y)
				}
			} else {
				x, y = pdf.GetXY()
				pdf.MultiCell(w-2, cellHeight2, strconv.Itoa(vinc.NumeroHorasSemanales), "1", "C", false)
				pdf.SetX(x)
				pdf.MultiCell(w-2, lineHeight, "Pasa a", "TLR", "C", false)
				pdf.SetX(x)
				pdf.MultiCell(w-2, lineHeight, fmt.Sprintf("%.f", valoresAntes["NumeroHorasSemanales"]-float64(vinc.NumeroHorasSemanales)), "BLR", "C", false)
				if pdf.GetY()-y > lineHeight {
					pdf.SetXY(x+w-2, y)
				}
				pdf.CellFormat(w-1, cellHeight, fmt.Sprintf(CampoMeses, vinc.NumeroSemanas), "1", 0, "C", false, 0, "")
			}

			x, y = pdf.GetXY()
			pdf.MultiCell(w+1, cellHeight2, FormatMoney(valoresAntes["SueldoBasico"], 2), "1", "C", false)
			pdf.SetX(x)
			pdf.MultiCell(w+1, lineHeight, "Pasa a", "TLR", "C", false)
			pdf.SetX(x)
			valorContrato := DeformatNumber(vinc.ValorContratoFormato)
			pdf.MultiCell(w+1, lineHeight, FormatMoney(valoresAntes["SueldoBasico"]-valorContrato, 2), "BLR", "C", false)
			if pdf.GetY()-y > lineHeight {
				pdf.SetXY(x+w+1, y)
			}
			pdf.CellFormat(w+2, cellHeight, vinc.ValorContratoFormato, "1", 0, "C", false, 0, "")
			break
		default:
			// Horas semanales|semestrales
			horasAntes := int(valoresAntes["NumeroHorasSemanales"])
			horasDespues := vinc.NumeroHorasSemanales
			horasTotales := 0
			x, y = pdf.GetXY()
			pdf.MultiCell(w-2, lineHeight, strconv.Itoa(horasAntes), "1", "C", false)
			pdf.SetX(x)
			pdf.MultiCell(w-2, lineHeight, strconv.Itoa(horasDespues), "1", "C", false)
			pdf.SetX(x)
			if tipoRes == "RADD" {
				horasTotales = horasAntes + horasDespues
			} else {
				horasTotales = horasAntes - horasDespues
			}
			pdf.MultiCell(w-2, lineHeight, fmt.Sprintf("Total %d", horasTotales), "1", "C", false)
			if pdf.GetY()-y > lineHeight {
				pdf.SetXY(x+w-2, y)
			}

			// Periodo de vinculacion
			x, y = pdf.GetXY()
			pdf.MultiCell(w-1, lineHeight, fmt.Sprintf(CampoMeses, int(valoresAntes["NumeroSemanas"])), "1", "C", false)
			pdf.SetX(x)
			pdf.MultiCell(w-1, lineHeight, fmt.Sprintf(CampoMeses, vinc.NumeroSemanas), "1", "C", false)
			pdf.SetX(x)
			pdf.MultiCell(w-1, lineHeight, "", "1", "C", false)
			if pdf.GetY()-y > lineHeight {
				pdf.SetXY(x+w-1, y)
			}

			// Valor contrato
			x, y = pdf.GetXY()
			pdf.MultiCell(w+1, lineHeight, FormatMoney(valoresAntes["SueldoBasico"], 2), "1", "C", false)
			pdf.SetX(x)
			pdf.MultiCell(w+1, lineHeight, "", "1", "C", false)
			pdf.SetX(x)
			valorContratoCambio := DeformatNumber(vinc.ValorContratoFormato)
			valorFinalContrato := 0.0
			if tipoRes == "RADD" {
				valorFinalContrato = valoresAntes["SueldoBasico"] + valorContratoCambio
			} else {
				valorFinalContrato = valoresAntes["SueldoBasico"] - valorContratoCambio
			}
			pdf.MultiCell(w+1, lineHeight, FormatMoney(valorFinalContrato, 2), "1", "C", false)
			if pdf.GetY()-y > lineHeight {
				pdf.SetXY(x+w+1, y)
			}

			//Valor Adicionar|Reversar
			x, y = pdf.GetXY()
			pdf.MultiCell(w+2, lineHeight, "", "1", "C", false)
			pdf.SetX(x)
			pdf.MultiCell(w+2, lineHeight, vinc.ValorContratoFormato, "1", "C", false)
			pdf.SetX(x)
			pdf.MultiCell(w+2, lineHeight, "", "1", "C", false)
			if pdf.GetY()-y > lineHeight {
				pdf.SetXY(x+w+2, y)
			}
			break
		}

		if tipoRes == "RVIN" || tipoRes == "RADD" {
			pdf.CellFormat(7, cellHeight, strconv.Itoa(vinc.Disponibilidad), "1", 0, "C", false, 0, "")
		} else {
			pdf.CellFormat(7, cellHeight, strconv.Itoa(vinc.RegistroPresupuestal), "1", 0, "C", false, 0, "")
		}
		pdf.Ln(-1)
	}
	pdf.Ln(lineHeight)
	return pdf, outputError
}

// Construye las tablas de vinculaciones con los rubros presupuestales
func ConstruirVinculacionesDesagregado(pdf *gofpdf.Fpdf, vinculaciones []models.Vinculaciones, lineHeight, fontSize float64, tipoRes, nivel string) (doc *gofpdf.Fpdf, outputError map[string]interface{}) {
	defer func() {
		if err := recover(); err != nil {
			outputError = map[string]interface{}{"funcion": "ConstruirVinculacionesDesagregado", "err": err, "status": "500"}
			panic(outputError)
		}
	}()
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
				panic(outputError)
			}
			pdf.Ln(lineHeight * 2)
			pdf.SetFont(Calibri, "", fontSize)
			pdf.Write(lineHeight, proyectoCurricular.Nombre)
			pdf.SetFont(Calibri, "", fontSize-3.5)
			pdf.Ln(lineHeight * 2)

		}
		// Encabezados
		pdf.CellFormat(w+4, lineHeight*2, "Nombre", "1", 0, "C", false, 0, "")
		pdf.CellFormat(w+2, lineHeight*2, "Tipo Documento", "1", 0, "C", false, 0, "")
		pdf.CellFormat(w-2, lineHeight*2, "Cédula", "1", 0, "C", false, 0, "")
		pdf.CellFormat(w-1, lineHeight*2, "Expedida", "1", 0, "C", false, 0, "")
		pdf.CellFormat(w+1, lineHeight*2, "Categoría", "1", 0, "C", false, 0, "")
		if tipoRes == "RVIN" {
			pdf.CellFormat(w+1, lineHeight*2, "Dedicación", "1", 0, "C", false, 0, "")
		} else {
			pdf.CellFormat(w-2, lineHeight*2, "Dedicación", "1", 0, "C", false, 0, "")
		}

		horas := ""
		if nivel == "PREGRADO" {
			horas = "Horas semanales"
		} else {
			horas = "Horas semestrales"
		}

		if tipoRes == "RVIN" {
			x, y := pdf.GetXY()
			pdf.MultiCell(w-3, lineHeight, horas, "1", "C", false)
			if pdf.GetY()-y > lineHeight {
				pdf.SetXY(x+w-3, y)
			}
			x, y = pdf.GetXY()
			pdf.MultiCell(w, lineHeight, "Periodo de vinculación", "1", "C", false)
			if pdf.GetY()-y > lineHeight {
				pdf.SetXY(x+w, y)
			}
			pdf.CellFormat(7, lineHeight*2, "CDP", "1", 0, "C", false, 0, "")
			pdf.CellFormat(w+1, lineHeight*2, "Valor total", "1", 0, "C", false, 0, "")
		}

		if tipoRes == "RADD" || tipoRes == "RRED" {
			valor := ""
			if tipoRes == "RADD" {
				pdf.CellFormat(w-2, lineHeight*2, "CDP", "1", 0, "C", false, 0, "")
				valor = "Valor total a adicionar"
			} else {
				pdf.CellFormat(w-2, lineHeight*2, "CRP", "1", 0, "C", false, 0, "")
				valor = "Valor total a reducir"
			}
			x, y := pdf.GetXY()
			pdf.MultiCell((w*2)-2, lineHeight*2, valor, "1", "C", false)
			if pdf.GetY()-y > lineHeight {
				pdf.SetXY(x+(w*2)-2, y)
			}
			x, y = pdf.GetXY()
			pdf.MultiCell(w-3, lineHeight, horas, "1", "C", false)
			if pdf.GetY()-y > lineHeight {
				pdf.SetXY(x+w-3, y)
			}
			x, y = pdf.GetXY()
			pdf.MultiCell(w-2, lineHeight, "Periodo de vinculación", "1", "C", false)
			if pdf.GetY()-y > lineHeight {
				pdf.SetXY(x+w-2, y+lineHeight)
			}
		}

		if tipoRes == "RCAN" {
			pdf.CellFormat(w-2, lineHeight*2, "CRP", "1", 0, "C", false, 0, "")
			x, y := pdf.GetXY()
			pdf.MultiCell(w-3, lineHeight, horas, "1", "C", false)
			if pdf.GetY()-y > lineHeight {
				pdf.SetXY(x+w-3, y)
			}
			pdf.CellFormat(w+1, lineHeight*2, "Valor a reversar", "1", 0, "C", false, 0, "")
			x, y = pdf.GetXY()
			pdf.MultiCell(w, lineHeight, "Periodo de vinculación", "1", "C", false)
			if pdf.GetY()-y > lineHeight {
				pdf.SetXY(x+w, y+lineHeight)
			}
		}
		pdf.Ln(-1)

		splitText := pdf.SplitLines([]byte(vinc.Nombre), w+4)
		cellHeight := lineHeight * 4 // float64(len(splitText)) * lineHeight
		for float64(len(splitText)) < 4 {
			splitText = append(splitText, []byte(""))
		}
		if cellHeight < minHeight {
			cellHeight = minHeight
		}
		if cellHeight > maxHeight {
			maxHeight = cellHeight
		}
		_, h := pdf.GetPageSize()
		_, _, _, b := pdf.GetMargins()
		if pdf.GetY() > h-b-(cellHeight) {
			pdf.AddPage()
		}
		// Nombre
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
		// Tipo documento
		cellHeight = lineHeight * 2
		x, y = pdf.GetXY()
		pdf.MultiCell(w+2, lineHeight, vinc.TipoDocumento, "1", "C", false)
		if pdf.GetY()-y > lineHeight {
			pdf.SetXY(x+w+2, y)
		}
		pdf.CellFormat(w-2, cellHeight, fmt.Sprintf("%.f", vinc.PersonaId), "1", 0, "C", false, 0, "")
		pdf.CellFormat(w-1, cellHeight, vinc.ExpedicionDocumento, "1", 0, "C", false, 0, "")
		pdf.CellFormat(w+1, cellHeight, vinc.Categoria, "1", 0, "C", false, 0, "")
		if tipoRes == "RVIN" {
			pdf.CellFormat(w+1, cellHeight, vinc.Dedicacion, "1", 0, "C", false, 0, "")
		} else {
			pdf.CellFormat(w-2, cellHeight, vinc.Dedicacion, "1", 0, "C", false, 0, "")
		}

		valorHoras := strconv.Itoa(vinc.NumeroHorasSemanales)

		if tipoRes == "RVIN" {
			var desagregado []models.DisponibilidadVinculacion
			url := "disponibilidad_vinculacion?query=Activo:true,VinculacionDocenteId.Id:" + strconv.Itoa(vinc.Id)
			if err := GetRequestNew("UrlCrudResoluciones", url, &desagregado); err != nil {
				logs.Error(err.Error())
				panic(err.Error())
			}
			valores := map[string]float64{}
			for _, disp := range desagregado {
				valores[disp.Rubro] = disp.Valor
			}

			pdf.CellFormat(w-3, cellHeight, valorHoras, "1", 0, "C", false, 0, "")
			pdf.CellFormat(w, cellHeight, fmt.Sprintf(CampoMeses, vinc.NumeroSemanas), "1", 0, "C", false, 0, "")

			cellHeight = lineHeight * 4
			pdf.CellFormat(7, cellHeight, strconv.Itoa(vinc.Disponibilidad), "1", 0, "C", false, 0, "")
			pdf.CellFormat(w+1, cellHeight, vinc.ValorContratoFormato, "1", 0, "C", false, 0, "")
			pdf.Ln(-1)

			x, y = pdf.GetXY()
			pdf.SetXY(x+w+4, y-(2*lineHeight))
			pdf.CellFormat(w+2, lineHeight, "Sueldo", "1", 0, "C", false, 0, "")
			pdf.CellFormat(w-2, lineHeight, "Prima navidad", "1", 0, "C", false, 0, "")
			pdf.CellFormat(w-1, lineHeight, "Vacaciones", "1", 0, "C", false, 0, "")
			pdf.CellFormat(w+1, lineHeight, "Prima Vacaciones", "1", 0, "C", false, 0, "")
			pdf.CellFormat(w+1, lineHeight, "Interes Cesantias", "1", 0, "C", false, 0, "")
			pdf.CellFormat(w-3, lineHeight, "Cesantias", "1", 0, "C", false, 0, "")
			pdf.CellFormat(w, lineHeight, "Prima servicios", "1", 0, "C", false, 0, "")
			pdf.Ln(-1)
			pdf.SetXY(x+w+4, y-lineHeight)
			pdf.CellFormat(w+2, lineHeight, FormatMoney(valores["SueldoBasico"], 2), "1", 0, "C", false, 0, "")
			pdf.CellFormat(w-2, lineHeight, FormatMoney(valores["PrimaNavidad"], 2), "1", 0, "C", false, 0, "")
			pdf.CellFormat(w-1, lineHeight, FormatMoney(valores["Vacaciones"], 2), "1", 0, "C", false, 0, "")
			pdf.CellFormat(w+1, lineHeight, FormatMoney(valores["PrimaVacaciones"], 2), "1", 0, "C", false, 0, "")
			pdf.CellFormat(w+1, lineHeight, FormatMoney(valores["InteresesCesantias"], 2), "1", 0, "C", false, 0, "")
			pdf.CellFormat(w-3, lineHeight, FormatMoney(valores["Cesantias"], 2), "1", 0, "C", false, 0, "")
			pdf.CellFormat(w, lineHeight, FormatMoney(valores["PrimaServicios"], 2), "1", 0, "C", false, 0, "")
		}

		if tipoRes == "RADD" || tipoRes == "RRED" {
			cellHeight = lineHeight * 2
			if tipoRes == "RADD" {
				pdf.CellFormat(w-2, cellHeight, strconv.Itoa(vinc.Disponibilidad), "1", 0, "C", false, 0, "")
			} else {
				pdf.CellFormat(w-2, cellHeight, strconv.Itoa(vinc.RegistroPresupuestal), "1", 0, "C", false, 0, "")
			}
			pdf.CellFormat((w*2)-2, cellHeight, vinc.ValorContratoFormato, "1", 0, "C", false, 0, "")
			pdf.CellFormat(w-3, cellHeight*2, valorHoras, "1", 0, "C", false, 0, "")
			pdf.CellFormat(w-2, cellHeight*2, fmt.Sprintf(CampoMeses, vinc.NumeroSemanas), "1", 0, "C", false, 0, "")
			pdf.Ln(-1)
		}

		if tipoRes == "RCAN" {
			cellHeight = lineHeight * 2
			pdf.CellFormat(w-2, cellHeight, strconv.Itoa(vinc.RegistroPresupuestal), "1", 0, "C", false, 0, "")
			pdf.CellFormat(w-3, cellHeight, valorHoras, "1", 0, "C", false, 0, "")
			pdf.CellFormat(w+1, cellHeight, vinc.ValorContratoFormato, "1", 0, "C", false, 0, "")
			pdf.CellFormat(w, cellHeight*2, fmt.Sprintf(CampoMeses, vinc.NumeroSemanas), "1", 0, "C", false, 0, "")
			pdf.Ln(-1)
		}

		if tipoRes != "RVIN" {

			valoresAntes := make(map[string]float64)
			if err := CalcularTrazabilidad(strconv.Itoa(vinc.Id), &valoresAntes); err != nil {
				logs.Error("Error en trazabilidad -> " + err.Error())
				panic("Error en trazabilidad -> " + err.Error())
			}

			var desagregadoDespues []models.DisponibilidadVinculacion
			url3 := "disponibilidad_vinculacion?query=Activo:true,VinculacionDocenteId.Id:" + strconv.Itoa(vinc.Id)
			if err := GetRequestNew("UrlCrudResoluciones", url3, &desagregadoDespues); err != nil {
				logs.Error(err.Error())
				panic(err.Error())
			}
			valoresDespues := map[string]float64{}
			for _, disp := range desagregadoDespues {
				valoresDespues[disp.Rubro] = disp.Valor
			}

			x, y = pdf.GetXY()
			pdf.SetXY(x+w+4, y-cellHeight)
			pdf.CellFormat(w+2, cellHeight, "Sueldo", "1", 0, "C", false, 0, "")

			x, y = pdf.GetXY()
			pdf.MultiCell(w-2, lineHeight, "Prima navidad", "1", "C", false)
			if pdf.GetY()-y > lineHeight {
				pdf.SetXY(x+w-2, y)
			}
			pdf.CellFormat(w-1, cellHeight, "Vacaciones", "1", 0, "C", false, 0, "")

			x, y = pdf.GetXY()
			pdf.MultiCell(w+1, lineHeight, "Prima vacaciones", "1", "C", false)
			if pdf.GetY()-y > lineHeight {
				pdf.SetXY(x+w+1, y)
			}

			x, y = pdf.GetXY()
			pdf.MultiCell(w-2, lineHeight, "Intereses cesantias", "1", "C", false)
			if pdf.GetY()-y > lineHeight {
				pdf.SetXY(x+w-2, y)
			}
			pdf.CellFormat(w-2, cellHeight, "Cesantias", "1", 0, "C", false, 0, "")

			x, y = pdf.GetXY()
			pdf.MultiCell(w-3, lineHeight, "Prima servicios", "1", "C", false)
			if pdf.GetY()-y > lineHeight {
				pdf.SetXY(x+w-3, y)
			}

			pdf.CellFormat(w+1, cellHeight, "Totales", "1", 0, "C", false, 0, "")
			pdf.Ln(-1)

			pdf.CellFormat(w+4, lineHeight, "Valores originales", "1", 0, "C", false, 0, "")
			pdf.CellFormat(w+2, lineHeight, FormatMoney(valoresAntes["SueldoBasico"], 2), "1", 0, "C", false, 0, "")
			pdf.CellFormat(w-2, lineHeight, FormatMoney(valoresAntes["PrimaNavidad"], 2), "1", 0, "C", false, 0, "")
			pdf.CellFormat(w-1, lineHeight, FormatMoney(valoresAntes["Vacaciones"], 2), "1", 0, "C", false, 0, "")
			pdf.CellFormat(w+1, lineHeight, FormatMoney(valoresAntes["PrimaVacaciones"], 2), "1", 0, "C", false, 0, "")
			pdf.CellFormat(w-2, lineHeight, FormatMoney(valoresAntes["InteresesCesantias"], 2), "1", 0, "C", false, 0, "")
			pdf.CellFormat(w-2, lineHeight, FormatMoney(valoresAntes["Cesantias"], 2), "1", 0, "C", false, 0, "")
			pdf.CellFormat(w-3, lineHeight, FormatMoney(valoresAntes["PrimaServicios"], 2), "1", 0, "C", false, 0, "")
			pdf.CellFormat(w+1, lineHeight, FormatMoney(valoresAntes["ValorContrato"], 2), "1", 0, "C", false, 0, "")

			if tipoRes == "RADD" || tipoRes == "RRED" {
				pdf.CellFormat(w-3, lineHeight, strconv.Itoa(int(valoresAntes["NumeroHorasSemanales"])), "1", 0, "C", false, 0, "")
				pdf.CellFormat(w-2, lineHeight*3, fmt.Sprintf(CampoMeses, int(valoresAntes["NumeroSemanas"])), "1", 0, "C", false, 0, "")
				x, y = pdf.GetXY()
				pdf.SetXY(x, y-lineHeight*2)
			}

			filaValores := "Valores a "
			switch tipoRes {
			case "RADD":
				filaValores += "adicionar"
				break
			case "RRED":
				filaValores += "reducir"
				break
			case "RCAN":
				filaValores += "reversar"
				pdf.CellFormat(w, lineHeight, fmt.Sprintf(CampoMeses, int(valoresAntes["NumeroSemanas"])), "1", 0, "C", false, 0, "")
				break
			}
			pdf.Ln(-1)

			pdf.CellFormat(w+4, lineHeight, filaValores, "1", 0, "C", false, 0, "")
			pdf.CellFormat(w+2, lineHeight, FormatMoney(valoresDespues["SueldoBasico"], 2), "1", 0, "C", false, 0, "")
			pdf.CellFormat(w-2, lineHeight, FormatMoney(valoresDespues["PrimaNavidad"], 2), "1", 0, "C", false, 0, "")
			pdf.CellFormat(w-1, lineHeight, FormatMoney(valoresDespues["Vacaciones"], 2), "1", 0, "C", false, 0, "")
			pdf.CellFormat(w+1, lineHeight, FormatMoney(valoresDespues["PrimaVacaciones"], 2), "1", 0, "C", false, 0, "")
			pdf.CellFormat(w-2, lineHeight, FormatMoney(valoresDespues["InteresesCesantias"], 2), "1", 0, "C", false, 0, "")
			pdf.CellFormat(w-2, lineHeight, FormatMoney(valoresDespues["Cesantias"], 2), "1", 0, "C", false, 0, "")
			pdf.CellFormat(w-3, lineHeight, FormatMoney(valoresDespues["PrimaServicios"], 2), "1", 0, "C", false, 0, "")
			pdf.CellFormat(w+1, lineHeight, vinc.ValorContratoFormato, "1", 0, "C", false, 0, "")
			switch tipoRes {
			case "RADD":
				pdf.CellFormat(w-3, lineHeight*2, fmt.Sprintf(PasaA, int(valoresAntes["NumeroHorasSemanales"]+float64(vinc.NumeroHorasSemanales))), "1", 0, "C", false, 0, "")
			case "RRED":
				pdf.CellFormat(w-3, lineHeight*2, fmt.Sprintf(PasaA, int(valoresAntes["NumeroHorasSemanales"]-float64(vinc.NumeroHorasSemanales))), "1", 0, "C", false, 0, "")
				break
			case "RCAN":
				pdf.CellFormat(w, lineHeight*2, fmt.Sprintf(PasaA, int(valoresAntes["NumeroSemanas"]-float64(vinc.NumeroSemanas))), "1", 0, "C", false, 0, "")
				break
			}
			x, y = pdf.GetXY()
			pdf.SetXY(x, y-lineHeight)
			pdf.Ln(-1)

			valorContrato := DeformatNumber(vinc.ValorContratoFormato)
			pdf.CellFormat(w+4, lineHeight, "Pasa a", "1", 0, "C", false, 0, "")
			if tipoRes == "RADD" {
				pdf.CellFormat(w+2, lineHeight, FormatMoney(valoresAntes["SueldoBasico"]+valoresDespues["SueldoBasico"], 2), "1", 0, "C", false, 0, "")
				pdf.CellFormat(w-2, lineHeight, FormatMoney(valoresAntes["PrimaNavidad"]+valoresDespues["PrimaNavidad"], 2), "1", 0, "C", false, 0, "")
				pdf.CellFormat(w-1, lineHeight, FormatMoney(valoresAntes["Vacaciones"]+valoresDespues["Vacaciones"], 2), "1", 0, "C", false, 0, "")
				pdf.CellFormat(w+1, lineHeight, FormatMoney(valoresAntes["PrimaVacaciones"]+valoresDespues["PrimaVacaciones"], 2), "1", 0, "C", false, 0, "")
				pdf.CellFormat(w-2, lineHeight, FormatMoney(valoresAntes["InteresesCesantias"]+valoresDespues["InteresesCesantias"], 2), "1", 0, "C", false, 0, "")
				pdf.CellFormat(w-2, lineHeight, FormatMoney(valoresAntes["Cesantias"]+valoresDespues["Cesantias"], 2), "1", 0, "C", false, 0, "")
				pdf.CellFormat(w-3, lineHeight, FormatMoney(valoresAntes["PrimaServicios"]+valoresDespues["PrimaServicios"], 2), "1", 0, "C", false, 0, "")
				pdf.CellFormat(w+1, lineHeight, FormatMoney(valoresAntes["ValorContrato"]+valorContrato, 2), "1", 0, "C", false, 0, "")
			} else {
				pdf.CellFormat(w+2, lineHeight, FormatMoney(valoresAntes["SueldoBasico"]-valoresDespues["SueldoBasico"], 2), "1", 0, "C", false, 0, "")
				pdf.CellFormat(w-2, lineHeight, FormatMoney(valoresAntes["PrimaNavidad"]-valoresDespues["PrimaNavidad"], 2), "1", 0, "C", false, 0, "")
				pdf.CellFormat(w-1, lineHeight, FormatMoney(valoresAntes["Vacaciones"]-valoresDespues["Vacaciones"], 2), "1", 0, "C", false, 0, "")
				pdf.CellFormat(w+1, lineHeight, FormatMoney(valoresAntes["PrimaVacaciones"]-valoresDespues["PrimaVacaciones"], 2), "1", 0, "C", false, 0, "")
				pdf.CellFormat(w-2, lineHeight, FormatMoney(valoresAntes["InteresesCesantias"]-valoresDespues["InteresesCesantias"], 2), "1", 0, "C", false, 0, "")
				pdf.CellFormat(w-2, lineHeight, FormatMoney(valoresAntes["Cesantias"]-valoresDespues["Cesantias"], 2), "1", 0, "C", false, 0, "")
				pdf.CellFormat(w-3, lineHeight, FormatMoney(valoresAntes["PrimaServicios"]-valoresDespues["PrimaServicios"], 2), "1", 0, "C", false, 0, "")
				pdf.CellFormat(w+1, lineHeight, FormatMoney(valoresAntes["ValorContrato"]-valorContrato, 2), "1", 0, "C", false, 0, "")
			}
			pdf.Ln(-1)

		}

		pdf.Ln(-1)
		pdf.Ln(lineHeight * 2)
	}

	return pdf, outputError
}

// Genera la tabla del cuadro de responsabilidades que va l final de cada resolución
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
	//pdf.OutputFileAndClose("resolucion.pdf") // para guardar el archivo localmente
	pdf.Output(writer)
	writer.Flush()
	encodedFile := base64.StdEncoding.EncodeToString(buffer.Bytes())
	return encodedFile
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

// Realiza el proceso de almacenar la resolución a traves del gestór documental
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
