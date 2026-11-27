package helpers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"strconv"
	"strings"

	"github.com/astaxie/beego"
	"github.com/udistrital/resoluciones_mid_v2/models"
	"github.com/udistrital/utils_oas/request"
)

const previnculacionSinDatosMessage = "No se encontraron datos en el proceso Cargue de Soportes Previnculación para el docente"

func validarDocentesVinculables(ctx context.Context, datos models.ObjetoPrevinculaciones) map[string]interface{} {
	if len(datos.Docentes) == 0 {
		return nil
	}

	respuesta, err := consultarRequisitosVinculacion(ctx, datos)
	if err != nil {
		return err
	}

	noVinculables := construirDocentesNoVinculables(datos.Docentes, respuesta)
	if len(noVinculables) == 0 {
		return nil
	}

	return map[string]interface{}{
		"funcion":       "/validarDocentesVinculables",
		"err":           construirMensajeDocentesNoVinculables(len(datos.Docentes), noVinculables),
		"status":        fmt.Sprintf("%d", http.StatusConflict),
		"code":          "docentes_no_vinculables",
		"data":          noVinculables,
		"selectedCount": len(datos.Docentes),
	}
}

func consultarRequisitosVinculacion(ctx context.Context, datos models.ObjetoPrevinculaciones) ([]models.RequisitoVinculacion, map[string]interface{}) {
	if len(datos.Docentes) == 0 {
		return nil, nil
	}

	url := strings.TrimSpace(beego.AppConfig.String("UrlPrevinculaciones"))
	if url == "" {
		return nil, map[string]interface{}{
			"funcion": "/consultarRequisitosVinculacion",
			"err":     "configuración RESOLUCIONES_MID_V2_PREVINCULACIONES requerida para validar el proceso Cargue de Soportes Previnculación",
			"status":  fmt.Sprintf("%d", http.StatusInternalServerError),
		}
	}

	url = strings.TrimRight(url, "/") + "/previnculacion"

	periodo, err := strconv.Atoi(strings.TrimSpace(datos.Docentes[0].Periodo))
	if err != nil || periodo <= 0 {
		return nil, map[string]interface{}{
			"funcion": "/consultarRequisitosVinculacion",
			"err":     "período inválido para consultar previnculación",
			"status":  fmt.Sprintf("%d", http.StatusBadRequest),
		}
	}
	body := models.RequisitosVinculacionRequest{
		Anio:      datos.Vigencia,
		Periodo:   periodo,
		IdUsuario: construirIdUsuariosPrevinculacion(datos.Docentes),
	}

	var response interface{}
	if _, err := request.PostWithContext(ctx, url, body, &response); err != nil {
		return nil, map[string]interface{}{
			"funcion": "/consultarRequisitosVinculacion",
			"err":     fmt.Sprintf("error consultando el proceso Cargue de Soportes Previnculación: %v", err),
			"status":  fmt.Sprintf("%d", http.StatusBadGateway),
		}
	}

	return normalizarRespuestaRequisitosPrevinculacion(response)
}

func construirIdUsuariosPrevinculacion(docentes []models.CargaLectiva) []string {
	ids := make([]string, 0, len(docentes))
	seen := make(map[string]struct{}, len(docentes))

	for _, docente := range docentes {
		documento := strings.TrimSpace(docente.DocDocente)
		if documento == "" {
			continue
		}
		if _, ok := seen[documento]; ok {
			continue
		}
		seen[documento] = struct{}{}
		ids = append(ids, documento)
	}

	return ids
}

func construirDocentesNoVinculables(docentes []models.CargaLectiva, respuesta []models.RequisitoVinculacion) []models.DocenteNoVinculable {
	docentesPorDocumento := make(map[string]models.CargaLectiva, len(docentes))
	for _, docente := range docentes {
		documento := strings.TrimSpace(docente.DocDocente)
		if documento == "" {
			continue
		}
		if _, ok := docentesPorDocumento[documento]; !ok {
			docentesPorDocumento[documento] = docente
		}
	}

	respuestaPorDocumento := make(map[string]models.RequisitoVinculacion, len(respuesta))
	for _, requisito := range respuesta {
		documento := strings.TrimSpace(requisito.IdUsuario)
		if documento == "" {
			continue
		}
		respuestaPorDocumento[documento] = requisito
	}

	documentos := make([]string, 0, len(docentesPorDocumento))
	for documento := range docentesPorDocumento {
		documentos = append(documentos, documento)
	}
	sort.Strings(documentos)

	noVinculables := make([]models.DocenteNoVinculable, 0)
	for _, documento := range documentos {
		docente := docentesPorDocumento[documento]
		requisito, ok := respuestaPorDocumento[documento]
		if !ok {
			noVinculables = append(noVinculables, models.DocenteNoVinculable{
				Documento: documento,
				Nombre:    construirNombreDocente(docente),
				Motivos:   []string{previnculacionSinDatosMessage},
			})
			continue
		}

		if strings.EqualFold(strings.TrimSpace(requisito.ResolucionVinculable), "S") {
			continue
		}

		noVinculables = append(noVinculables, models.DocenteNoVinculable{
			Documento: documento,
			Nombre:    construirNombreDocente(docente),
			Motivos:   construirMotivosNoVinculable(requisito),
		})
	}

	return noVinculables
}

func construirNombreDocente(docente models.CargaLectiva) string {
	return strings.TrimSpace(strings.Join([]string{
		strings.TrimSpace(docente.DocenteApellido),
		strings.TrimSpace(docente.DocenteNombre),
	}, " "))
}

func construirMotivosNoVinculable(requisito models.RequisitoVinculacion) []string {
	motivos := make([]string, 0)

	if estado := strings.TrimSpace(requisito.Estado); estado != "" && !strings.EqualFold(estado, "A") {
		nombreEstado := strings.TrimSpace(requisito.NombreEstado)
		if nombreEstado != "" {
			motivos = append(motivos, fmt.Sprintf("estado del Cargue de Soportes Previnculación: %s", nombreEstado))
		} else {
			motivos = append(motivos, fmt.Sprintf("estado del Cargue de Soportes Previnculación no activo: %s", estado))
		}
	}
	if !banderaAprobada(requisito.ApruebaSoporte) {
		motivos = append(motivos, "soporte no aprobado")
	}
	if !banderaAprobada(requisito.RegistroTercero) {
		motivos = append(motivos, "sin registro de tercero")
	}
	if !banderaAprobada(requisito.RegistroProveedor) {
		motivos = append(motivos, "sin registro de proveedor")
	}
	if !banderaAprobada(requisito.RegistroTipoVinculo) {
		motivos = append(motivos, "sin registro de tipo de vínculo")
	}
	if !banderaAprobada(requisito.RegistroNoCruce) {
		motivos = append(motivos, "sin validación de no cruce")
	}
	if !banderaAprobada(requisito.RegistroPreCarga) {
		motivos = append(motivos, "sin precarga")
	}
	if observacion := strings.TrimSpace(requisito.Observacion); observacion != "" {
		motivos = append(motivos, observacion)
	}
	if len(motivos) == 0 {
		motivos = append(motivos, "el proceso Cargue de Soportes Previnculación reporta docente no vinculable")
	}

	return motivos
}

func construirMensajeDocentesNoVinculables(totalSeleccionados int, docentes []models.DocenteNoVinculable) string {
	partes := make([]string, 0, len(docentes))
	for _, docente := range docentes {
		nombre := strings.TrimSpace(docente.Nombre)
		if nombre == "" {
			nombre = docente.Documento
		}
		partes = append(partes, fmt.Sprintf("%s (%s): %s", nombre, docente.Documento, strings.Join(docente.Motivos, ", ")))
	}

	return fmt.Sprintf(
		"Se seleccionaron %d docente(s). %d presentan inconsistencias en el proceso Cargue de Soportes Previnculación, por lo tanto la vinculación no continuará y no se vinculará ningún docente de la selección actual. %s",
		totalSeleccionados,
		len(docentes),
		strings.Join(partes, " | "),
	)
}

func banderaAprobada(valor string) bool {
	return strings.EqualFold(strings.TrimSpace(valor), "S")
}

func normalizarRespuestaRequisitosPrevinculacion(raw interface{}) ([]models.RequisitoVinculacion, map[string]interface{}) {
	switch typed := raw.(type) {
	case []interface{}:
		return convertirRespuestaPrevinculacion(typed)
	case map[string]interface{}:
		if data, ok := typed["Data"]; ok {
			return normalizarRespuestaRequisitosPrevinculacion(data)
		}
		if data, ok := typed["data"]; ok {
			return normalizarRespuestaRequisitosPrevinculacion(data)
		}

		detalle := extraerMensajeRespuestaPrevinculacion(typed)
		if esRespuestaSinDatosPrevinculacion(detalle) {
			return []models.RequisitoVinculacion{}, nil
		}
		if detalle == "" {
			detalle = fmt.Sprintf("respuesta inesperada de previnculación: %+v", typed)
		}

		return nil, map[string]interface{}{
			"funcion": "/consultarRequisitosVinculacion",
			"err":     fmt.Sprintf("el servicio de Cargue de Soportes Previnculación respondió un objeto no esperado: %s", detalle),
			"status":  fmt.Sprintf("%d", http.StatusBadGateway),
		}
	case nil:
		return nil, map[string]interface{}{
			"funcion": "/consultarRequisitosVinculacion",
			"err":     "el servicio de Cargue de Soportes Previnculación respondió vacío",
			"status":  fmt.Sprintf("%d", http.StatusBadGateway),
		}
	default:
		return nil, map[string]interface{}{
			"funcion": "/consultarRequisitosVinculacion",
			"err":     fmt.Sprintf("el servicio de Cargue de Soportes Previnculación respondió un formato no soportado: %T", raw),
			"status":  fmt.Sprintf("%d", http.StatusBadGateway),
		}
	}
}

func convertirRespuestaPrevinculacion(raw interface{}) ([]models.RequisitoVinculacion, map[string]interface{}) {
	bytes, err := json.Marshal(raw)
	if err != nil {
		return nil, map[string]interface{}{
			"funcion": "/consultarRequisitosVinculacion",
			"err":     fmt.Sprintf("error serializando respuesta del proceso Cargue de Soportes Previnculación: %v", err),
			"status":  fmt.Sprintf("%d", http.StatusBadGateway),
		}
	}

	var response []models.RequisitoVinculacion
	if err := json.Unmarshal(bytes, &response); err != nil {
		return nil, map[string]interface{}{
			"funcion": "/consultarRequisitosVinculacion",
			"err":     fmt.Sprintf("error interpretando respuesta del proceso Cargue de Soportes Previnculación: %v", err),
			"status":  fmt.Sprintf("%d", http.StatusBadGateway),
		}
	}

	return response, nil
}

func extraerMensajeRespuestaPrevinculacion(data map[string]interface{}) string {
	keys := []string{"Message", "message", "error", "Error", "detalle", "Detalle"}
	for _, key := range keys {
		if value, ok := data[key]; ok {
			if texto := strings.TrimSpace(fmt.Sprintf("%v", value)); texto != "" {
				return texto
			}
		}
	}
	return ""
}

func esRespuestaSinDatosPrevinculacion(mensaje string) bool {
	normalizado := strings.ToLower(strings.TrimSpace(mensaje))
	return strings.Contains(normalizado, "no se encontrarón datos asociados al reporte") ||
		strings.Contains(normalizado, "no se encontraron datos asociados al reporte")
}
