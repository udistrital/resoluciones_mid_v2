package services

import (
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/astaxie/beego/logs"
	"github.com/udistrital/resoluciones_mid_v2/helpers"
	"github.com/udistrital/resoluciones_mid_v2/models"
	"github.com/xuri/excelize/v2"
)

func stripBadTZ(payload map[string]interface{}) map[string]interface{} {
	for _, k := range []string{"FechaCreacion", "FechaModificacion"} {
		if v, ok := payload[k]; ok {
			if str, ok := v.(string); ok {
				str = strings.ReplaceAll(str, "+0000 +0000", "")
				str = strings.ReplaceAll(str, "+0000", "")
				payload[k] = strings.TrimSpace(str)
			}
		}
	}
	return payload
}

func sanitizePayload(data map[string]interface{}) map[string]interface{} {
	clean := make(map[string]interface{})
	for k, v := range data {
		if v == nil {
			continue
		}
		switch val := v.(type) {
		case map[string]interface{}:
			clean[k] = sanitizePayload(val)
		case []interface{}:
			arr := []interface{}{}
			for _, item := range val {
				if sub, ok := item.(map[string]interface{}); ok {
					arr = append(arr, sanitizePayload(sub))
				} else if item != nil {
					arr = append(arr, item)
				}
			}
			clean[k] = arr
		case float64, int, bool, string:
			clean[k] = val
		default:
			clean[k] = fmt.Sprintf("%v", val)
		}
	}
	return clean
}

func coerceInt(value interface{}) int {
	switch v := value.(type) {
	case float64:
		return int(v)
	case int:
		return v
	case string:
		re := regexp.MustCompile(`\d+`)
		s := strings.Join(re.FindAllString(v, -1), "")
		if s == "" {
			return 0
		}
		n, _ := strconv.Atoi(s)
		return n
	default:
		return 0
	}
}

func normalizeHeader(h string) string {
	replacements := map[string]string{
		"á": "a", "é": "e", "í": "i", "ó": "o", "ú": "u",
		"Á": "A", "É": "E", "Í": "I", "Ó": "O", "Ú": "U",
	}
	for k, v := range replacements {
		h = strings.ReplaceAll(h, k, v)
	}
	h = strings.ToLower(strings.TrimSpace(h))
	h = strings.ReplaceAll(h, " ", "_")
	return h
}

func validarExistenciaVinculacion(crp string, vigenciaRp int) (bool, error) {
	var vincs []map[string]interface{}
	query := fmt.Sprintf("NumeroRp:%s,VigenciaRp:%d,Activo:true", crp, vigenciaRp)
	url := "vinculacion_docente?query=" + query

	if err := helpers.GetRequestNew("UrlCrudResoluciones", url, &vincs); err != nil {
		return false, err
	}

	return len(vincs) > 0, nil
}

func ProcesarVinculaciones(file multipart.File, fileHeader *multipart.FileHeader, vigenciaRp int) ([]models.VinculacionRpResultado, error) {
	var resultados []models.VinculacionRpResultado

	if fileHeader == nil {
		return nil, errors.New("no se recibió archivo en la solicitud")
	}

	tmpDir := os.TempDir()
	tmpPath := filepath.Join(tmpDir, fmt.Sprintf("rp_%s.xlsx", time.Now().Format("02012006_150405")))
	out, err := os.Create(tmpPath)
	if err != nil {
		return nil, fmt.Errorf("no se pudo crear archivo temporal: %v", err)
	}
	defer out.Close()

	_, _ = file.Seek(0, 0)
	_, _ = io.Copy(out, file)

	f, err := excelize.OpenFile(tmpPath)
	if err != nil {
		logs.Error("Error al abrir el archivo Excel: %v", err)
		return nil, fmt.Errorf("no se pudo leer el archivo Excel: %v", err)
	}
	defer f.Close()

	sheetName := f.GetSheetName(0)
	if sheetName == "" {
		return nil, errors.New("el archivo no contiene hojas válidas")
	}

	rows, err := f.GetRows(sheetName)
	if err != nil {
		return nil, fmt.Errorf("error leyendo filas: %v", err)
	}
	if len(rows) < 2 {
		return nil, errors.New("el archivo no contiene datos suficientes")
	}

	headers := make(map[string]int)
	for i, h := range rows[0] {
		headers[normalizeHeader(h)] = i
	}

	requeridas := []string{"cod_resolucion", "cod_facultad", "documento", "cod_proyecto", "crp"}
	for _, col := range requeridas {
		if _, ok := headers[col]; !ok {
			return nil, fmt.Errorf("falta la columna requerida: %s", col)
		}
	}

	var registros []models.VinculacionRpResultado
	for i, row := range rows[1:] {
		get := func(key string) string {
			idx := headers[key]
			if idx < len(row) {
				return strings.TrimSpace(row[idx])
			}
			return ""
		}

		reg := models.VinculacionRpResultado{
			CodResolucion: get("cod_resolucion"),
			CodFacultad:   get("cod_facultad"),
			Documento:     get("documento"),
			CodProyecto:   get("cod_proyecto"),
			CRP:           get("crp"),
			FilaExcel:     i + 2,
		}
		registros = append(registros, reg)
	}

	type conflictoInfo struct {
		CRPs  map[string]bool
		Filas []int
	}

	conflictos := make(map[string]*conflictoInfo)

	for _, r := range registros {
		resNumKey := strings.Trim(strings.ReplaceAll(r.CodResolucion, "'", ""), " ")

		keyBase := fmt.Sprintf("%s-%s-%s-%s",
			resNumKey, r.CodFacultad, r.Documento, r.CodProyecto)

		if _, ok := conflictos[keyBase]; !ok {
			conflictos[keyBase] = &conflictoInfo{
				CRPs:  make(map[string]bool),
				Filas: []int{},
			}
		}

		conflictos[keyBase].CRPs[r.CRP] = true
		conflictos[keyBase].Filas = append(conflictos[keyBase].Filas, r.FilaExcel)
	}

	llavesInvalidas := make(map[string]*conflictoInfo)
	for k, info := range conflictos {
		if len(info.CRPs) > 1 {
			llavesInvalidas[k] = info
		}
	}

	visto := make(map[string]bool)
	var registrosUnicos []models.VinculacionRpResultado
	for _, r := range registros {
		clave := fmt.Sprintf("%s-%s-%s-%s-%s",
			r.CodResolucion, r.CodFacultad, r.Documento, r.CodProyecto, r.CRP)
		if !visto[clave] {
			visto[clave] = true
			registrosUnicos = append(registrosUnicos, r)
		}
	}

	for _, res := range registrosUnicos {
		resNum := strings.Trim(strings.ReplaceAll(res.CodResolucion, "'", ""), " ")

		keyBase := fmt.Sprintf("%s-%s-%s-%s",
			resNum, res.CodFacultad, res.Documento, res.CodProyecto)

		if info, ok := llavesInvalidas[keyBase]; ok {
			crps := make([]string, 0, len(info.CRPs))
			for crp := range info.CRPs {
				crps = append(crps, crp)
			}
			sort.Strings(crps)

			filas := append([]int{}, info.Filas...)
			sort.Ints(filas)

			res.PutStatus = fmt.Sprintf(
				"CONFLICTO: llave duplicada con CRPs diferentes. CRPs=%s. Filas=%v. No se actualiza ninguna fila de esta llave.",
				strings.Join(crps, ","),
				filas,
			)
			resultados = append(resultados, res)
			continue
		}

		var resoluciones []map[string]interface{}
		queryRes := fmt.Sprintf("NumeroResolucion:%s,Vigencia:%d,DependenciaId:%s,Activo:true", resNum, vigenciaRp, res.CodFacultad)
		if err := helpers.GetRequestNew("UrlCrudResoluciones", "resolucion?query="+queryRes, &resoluciones); err != nil {
			res.PutStatus = fmt.Sprintf("Error consultando resolución: %v", err)
			resultados = append(resultados, res)
			continue
		}
		if len(resoluciones) == 0 {
			res.PutStatus = "Resolución no encontrada"
			resultados = append(resultados, res)
			continue
		}
		res.IdResolucion = fmt.Sprintf("%.0f", resoluciones[0]["Id"].(float64))

		var vinculaciones []map[string]interface{}
		queryVin := fmt.Sprintf("ResolucionVinculacionDocenteId:%s,PersonaId:%s,ProyectoCurricularId:%s",
			res.IdResolucion, res.Documento, res.CodProyecto)
		if err := helpers.GetRequestNew("UrlCrudResoluciones", "vinculacion_docente?query="+queryVin, &vinculaciones); err != nil {
			res.PutStatus = fmt.Sprintf("Error consultando vinculación: %v", err)
			resultados = append(resultados, res)
			continue
		}
		if len(vinculaciones) == 0 {
			res.PutStatus = "Vinculación no encontrada"
			resultados = append(resultados, res)
			continue
		}

		var elegida map[string]interface{}
		for _, v := range vinculaciones {
			if nc, ok := v["NumeroContrato"]; ok && nc != nil {
				elegida = v
				break
			}
		}
		if elegida == nil {
			res.PutStatus = "No hay vinculación con NumeroContrato != null (no se actualiza)"
			resultados = append(resultados, res)
			continue
		}

		res.IdVinculacion = fmt.Sprintf("%.0f", elegida["Id"].(float64))

		var vincActual map[string]interface{}
		if err := helpers.GetRequestNew("UrlCrudResoluciones", "vinculacion_docente/"+res.IdVinculacion, &vincActual); err != nil {
			res.PutStatus = fmt.Sprintf("Error GET previo: %v", err)
			resultados = append(resultados, res)
			continue
		}
		if raw, ok := vincActual["Data"]; ok {
			if dataArr, ok := raw.([]interface{}); ok && len(dataArr) > 0 {
				if m, ok := dataArr[0].(map[string]interface{}); ok {
					vincActual = m
				} else {
					res.PutStatus = "Respuesta inválida: Data[0] no es objeto"
					resultados = append(resultados, res)
					continue
				}
			} else if dataMap, ok := raw.(map[string]interface{}); ok && dataMap != nil {
				vincActual = dataMap
			} else {
				res.PutStatus = "Respuesta inválida: Data no tiene el formato esperado"
				resultados = append(resultados, res)
				continue
			}
		} else {
			if _, ok := vincActual["Id"]; !ok {
				keys := make([]string, 0, len(vincActual))
				for k := range vincActual {
					keys = append(keys, k)
				}
				sort.Strings(keys)
				res.PutStatus = fmt.Sprintf("Respuesta inválida: sin Data y sin Id. Keys=%v", keys)
				resultados = append(resultados, res)
				continue
			}
		}

		for _, key := range []string{"Message", "Status", "Success"} {
			delete(vincActual, key)
		}

		if val, ok := vincActual["ResolucionVinculacionDocenteId"].(float64); ok {
			vincActual["ResolucionVinculacionDocenteId"] = map[string]interface{}{"Id": int(val)}
		}

		vincActual["NumeroRp"] = coerceInt(res.CRP)
		vincActual["VigenciaRp"] = vigenciaRp
		vincActual["Activo"] = true

		crpNum := strconv.Itoa(coerceInt(res.CRP))
		exists, err := validarExistenciaVinculacion(crpNum, vigenciaRp)
		if err != nil {
			res.PutStatus = fmt.Sprintf("Error validando duplicado: %v", err)
			resultados = append(resultados, res)
			continue
		}
		if exists {
			res.PutStatus = fmt.Sprintf("RP duplicado (ya existe en base de datos): CRP %s - Vigencia %d", res.CRP, vigenciaRp)
			resultados = append(resultados, res)
			continue
		}

		vincActual = stripBadTZ(vincActual)
		vincActual = sanitizePayload(vincActual)

		var respPut map[string]interface{}
		err = helpers.SendRequestFull("UrlCrudResoluciones",
			"vinculacion_docente/"+res.IdVinculacion, "PUT", &respPut, vincActual)

		if err != nil {
			res.PutStatus = fmt.Sprintf("Error PUT: %v", err)
		} else {
			success := false
			message := ""
			if s, ok := respPut["Success"].(bool); ok {
				success = s
			}
			if m, ok := respPut["Message"].(string); ok {
				message = m
			}
			if success || strings.Contains(strings.ToLower(message), "update successful") {
				res.PutStatus = "OK"
			} else {
				res.PutStatus = "PUT no exitoso"
			}
		}

		resultados = append(resultados, res)
	}

	return resultados, nil
}
