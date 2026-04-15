package services

import (
	"errors"
	"fmt"
	"mime/multipart"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/udistrital/resoluciones_mid_v2/helpers"
	"github.com/udistrital/resoluciones_mid_v2/models"
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

type conflictoInfo struct {
	CRPs  map[string]bool
	Filas []int
}

func extraerHeadersRp(headerRow []string) (map[string]int, error) {
	headers := make(map[string]int)
	for i, h := range headerRow {
		headers[normalizeHeader(h)] = i
	}

	requeridas := []string{"cod_resolucion", "cod_facultad", "documento", "cod_proyecto", "crp"}
	for _, col := range requeridas {
		if _, ok := headers[col]; !ok {
			return nil, fmt.Errorf("falta la columna requerida: %s", col)
		}
	}

	return headers, nil
}

func construirRegistrosRp(rows [][]string, headers map[string]int) []models.VinculacionRpResultado {
	registros := make([]models.VinculacionRpResultado, 0, len(rows))
	for i, row := range rows {
		get := func(key string) string {
			idx := headers[key]
			if idx < len(row) {
				return strings.TrimSpace(row[idx])
			}
			return ""
		}

		registros = append(registros, models.VinculacionRpResultado{
			CodResolucion: get("cod_resolucion"),
			CodFacultad:   get("cod_facultad"),
			Documento:     get("documento"),
			CodProyecto:   get("cod_proyecto"),
			CRP:           get("crp"),
			FilaExcel:     i + 2,
		})
	}

	return registros
}

func detectarConflictosRp(registros []models.VinculacionRpResultado) map[string]*conflictoInfo {
	conflictos := make(map[string]*conflictoInfo)

	for _, r := range registros {
		resNumKey := strings.Trim(strings.ReplaceAll(r.CodResolucion, "'", ""), " ")
		keyBase := fmt.Sprintf("%s-%s-%s-%s", resNumKey, r.CodFacultad, r.Documento, r.CodProyecto)

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

	return llavesInvalidas
}

func deduplicarRegistrosRp(registros []models.VinculacionRpResultado) []models.VinculacionRpResultado {
	visto := make(map[string]bool)
	resultado := make([]models.VinculacionRpResultado, 0, len(registros))

	for _, r := range registros {
		clave := fmt.Sprintf("%s-%s-%s-%s-%s", r.CodResolucion, r.CodFacultad, r.Documento, r.CodProyecto, r.CRP)
		if !visto[clave] {
			visto[clave] = true
			resultado = append(resultado, r)
		}
	}

	return resultado
}

func cargarPayloadVinculacionRp(idVinculacion string) (map[string]interface{}, error) {
	var vincActual map[string]interface{}
	if err := helpers.GetRequestNew("UrlCrudResoluciones", "vinculacion_docente/"+idVinculacion, &vincActual); err != nil {
		return nil, fmt.Errorf("Error GET previo: %v", err)
	}
	if raw, ok := vincActual["Data"]; ok {
		if dataArr, ok := raw.([]interface{}); ok && len(dataArr) > 0 {
			if m, ok := dataArr[0].(map[string]interface{}); ok {
				vincActual = m
			} else {
				return nil, errors.New("Respuesta inválida: Data[0] no es objeto")
			}
		} else if dataMap, ok := raw.(map[string]interface{}); ok && dataMap != nil {
			vincActual = dataMap
		} else {
			return nil, errors.New("Respuesta inválida: Data no tiene el formato esperado")
		}
	} else if _, ok := vincActual["Id"]; !ok {
		keys := make([]string, 0, len(vincActual))
		for k := range vincActual {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		return nil, fmt.Errorf("Respuesta inválida: sin Data y sin Id. Keys=%v", keys)
	}

	for _, key := range []string{"Message", "Status", "Success"} {
		delete(vincActual, key)
	}

	if val, ok := vincActual["ResolucionVinculacionDocenteId"].(float64); ok {
		vincActual["ResolucionVinculacionDocenteId"] = map[string]interface{}{"Id": int(val)}
	}

	return vincActual, nil
}

func resolverEstadoPutRp(respPut map[string]interface{}) string {
	success := false
	message := ""
	if s, ok := respPut["Success"].(bool); ok {
		success = s
	}
	if m, ok := respPut["Message"].(string); ok {
		message = m
	}
	if success || strings.Contains(strings.ToLower(message), "update successful") {
		return "OK"
	}
	return "PUT no exitoso"
}

func resolverResolucionRp(res *models.VinculacionRpResultado, vigenciaRp int) error {
	resNum := strings.Trim(strings.ReplaceAll(res.CodResolucion, "'", ""), " ")

	var resoluciones []map[string]interface{}
	queryRes := fmt.Sprintf("NumeroResolucion:%s,Vigencia:%d,DependenciaId:%s,Activo:true", resNum, vigenciaRp, res.CodFacultad)
	if err := helpers.GetRequestNew("UrlCrudResoluciones", "resolucion?query="+queryRes, &resoluciones); err != nil {
		return fmt.Errorf("Error consultando resolución: %v", err)
	}
	if len(resoluciones) == 0 {
		return errors.New("Resolución no encontrada")
	}

	res.IdResolucion = fmt.Sprintf("%.0f", resoluciones[0]["Id"].(float64))
	return nil
}

func resolverVinculacionObjetivoRp(res *models.VinculacionRpResultado) error {
	var vinculaciones []map[string]interface{}
	queryVin := fmt.Sprintf("ResolucionVinculacionDocenteId:%s,PersonaId:%s,ProyectoCurricularId:%s",
		res.IdResolucion, res.Documento, res.CodProyecto)
	if err := helpers.GetRequestNew("UrlCrudResoluciones", "vinculacion_docente?query="+queryVin, &vinculaciones); err != nil {
		return fmt.Errorf("Error consultando vinculación: %v", err)
	}
	if len(vinculaciones) == 0 {
		return errors.New("Vinculación no encontrada")
	}

	for _, v := range vinculaciones {
		if nc, ok := v["NumeroContrato"]; ok && nc != nil {
			res.IdVinculacion = fmt.Sprintf("%.0f", v["Id"].(float64))
			return nil
		}
	}

	return errors.New("No hay vinculación con NumeroContrato != null (no se actualiza)")
}

func aplicarRpAVinculacion(res models.VinculacionRpResultado, vigenciaRp int) string {
	vincActual, err := cargarPayloadVinculacionRp(res.IdVinculacion)
	if err != nil {
		return err.Error()
	}

	vincActual["NumeroRp"] = coerceInt(res.CRP)
	vincActual["VigenciaRp"] = vigenciaRp

	crpNum := strconv.Itoa(coerceInt(res.CRP))
	exists, err := validarExistenciaVinculacion(crpNum, vigenciaRp)
	if err != nil {
		return fmt.Sprintf("Error validando duplicado: %v", err)
	}
	if exists {
		return fmt.Sprintf("RP duplicado (ya existe en base de datos): CRP %s - Vigencia %d", res.CRP, vigenciaRp)
	}

	vincActual = stripBadTZ(vincActual)
	vincActual = sanitizePayload(vincActual)

	var respPut map[string]interface{}
	err = helpers.SendRequestFull("UrlCrudResoluciones",
		"vinculacion_docente/"+res.IdVinculacion, "PUT", &respPut, vincActual)
	if err != nil {
		return fmt.Sprintf("Error PUT: %v", err)
	}

	return resolverEstadoPutRp(respPut)
}

func ProcesarVinculaciones(file multipart.File, fileHeader *multipart.FileHeader, vigenciaRp int) ([]models.VinculacionRpResultado, error) {
	var resultados []models.VinculacionRpResultado

	registros, err := cargarRegistrosRpDesdeArchivo(file, fileHeader)
	if err != nil {
		return nil, err
	}

	llavesInvalidas := detectarConflictosRp(registros)
	registrosUnicos := deduplicarRegistrosRp(registros)

	for _, res := range registrosUnicos {
		resNum := strings.Trim(strings.ReplaceAll(res.CodResolucion, "'", ""), " ")

		keyBase := fmt.Sprintf("%s-%s-%s-%s",
			resNum, res.CodFacultad, res.Documento, res.CodProyecto)

		if info, ok := llavesInvalidas[keyBase]; ok {
			res.PutStatus = construirEstadoConflictoRp(info)
			resultados = append(resultados, res)
			continue
		}

		if err := resolverResolucionRp(&res, vigenciaRp); err != nil {
			res.PutStatus = err.Error()
			resultados = append(resultados, res)
			continue
		}

		if err := resolverVinculacionObjetivoRp(&res); err != nil {
			res.PutStatus = err.Error()
			resultados = append(resultados, res)
			continue
		}

		res.PutStatus = aplicarRpAVinculacion(res, vigenciaRp)

		resultados = append(resultados, res)
	}

	return resultados, nil
}
