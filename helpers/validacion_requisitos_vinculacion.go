package helpers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"strings"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/udistrital/resoluciones_mid_v2/models"
	"github.com/udistrital/utils_oas/request"
	"github.com/udistrital/utils_oas/ssm"
)

const (
	odinLoginPath                 = "auth/login"
	odinRequisitosVinculacionPath = "gen/apis?api=api_requisitos_vinculacion&proc=vinculacion_docente"
	odinSinDatosMessage           = "No se encontraron datos en el proceso Cargue de Soportes Previnculación para el docente"
	odinUserKey                   = "OdinServicioOATIUser"
	odinPasswordKey               = "OdinServicioOATIPassword"
	odinVersionKey                = "OdinServicioOATIVersion"
)

type odinConfig struct {
	BaseURL  string
	Username string
	Password string
	Version  string
}

func validarDocentesVinculables(ctx context.Context, datos models.ObjetoPrevinculaciones) map[string]interface{} {
	if len(datos.Docentes) == 0 {
		return nil
	}

	config, errMap := resolverConfiguracionOdin(ctx)
	if errMap != nil {
		return errMap
	}

	token, err := autenticarOdin(ctx, config.BaseURL, config.Username, config.Password, config.Version)
	if err != nil {
		return err
	}

	respuesta, err := consultarRequisitosVinculacion(ctx, config.BaseURL, token, datos)
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

func resolverConfiguracionOdin(ctx context.Context) (odinConfig, map[string]interface{}) {
	config := odinConfig{
		BaseURL: beego.AppConfig.String("UrlOdinServicioOATI"),
	}

	config.Username = beego.AppConfig.String(odinUserKey)
	password := beego.AppConfig.String(odinPasswordKey)
	config.Password = password
	config.Version = beego.AppConfig.String(odinVersionKey)

	logs.Info("ODIN auth config inicial: base_url=%q username_set=%t password_set=%t version_set=%t parameter_store=%q",
		config.BaseURL, config.Username != "", config.Password != "", config.Version != "", beego.AppConfig.String("parameterStore"))

	if config.BaseURL == "" {
		logs.Error("ODIN auth config incompleta: base_url vacía")
		return odinConfig{}, map[string]interface{}{
			"funcion": "/resolverConfiguracionOdin",
			"err":     "configuración incompleta para validar el proceso Cargue de Soportes Previnculación",
			"status":  fmt.Sprintf("%d", http.StatusInternalServerError),
		}
	}

	if config.Username != "" && config.Password != "" && config.Version != "" {
		logs.Info("ODIN auth config resuelta desde app.conf/env")
		return config, nil
	}

	parameterStore := beego.AppConfig.String("parameterStore")
	if parameterStore == "" {
		logs.Error("ODIN auth config incompleta: parameterStore vacío y faltan credenciales")
		return odinConfig{}, map[string]interface{}{
			"funcion": "/resolverConfiguracionOdin",
			"err":     "configuración incompleta para validar el proceso Cargue de Soportes Previnculación",
			"status":  fmt.Sprintf("%d", http.StatusInternalServerError),
		}
	}

	logs.Info("ODIN auth config: intentando resolver credenciales desde Parameter Store")

	usernameValue, err := resolverParametroOdin(ctx, parameterStore, odinUserKey)
	if err != nil {
		return odinConfig{}, err
	}

	passwordValue, err := resolverParametroOdin(ctx, parameterStore, odinPasswordKey)
	if err != nil {
		return odinConfig{}, err
	}

	versionValue, err := resolverParametroOdin(ctx, parameterStore, odinVersionKey)
	if err != nil {
		return odinConfig{}, err
	}

	config.Username = usernameValue
	config.Password = passwordValue
	config.Version = versionValue

	logs.Info("ODIN auth config resuelta desde Parameter Store: base_url=%q username_set=%t password_set=%t version_set=%t",
		config.BaseURL, config.Username != "", config.Password != "", config.Version != "")

	return config, nil
}

func resolverParametroOdin(ctx context.Context, parameterStore, parameterName string) (string, map[string]interface{}) {
	path := fmt.Sprintf("/%s/%s/%s", parameterStore, beego.AppConfig.String("appname"), parameterName)
	logs.Info("ODIN auth parameter lookup: parameter=%s path=%s", parameterName, path)
	value, err := ssm.GetParameterFromParameterStore(ctx, path)
	if err != nil {
		logs.Error("ODIN auth parameter lookup fallo: parameter=%s path=%s err=%v", parameterName, path, err)
		return "", map[string]interface{}{
			"funcion": "/resolverCredencialesOdin",
			"err":     fmt.Sprintf("error consultando %s en Parameter Store: %v", parameterName, err),
			"status":  fmt.Sprintf("%d", http.StatusBadGateway),
		}
	}

	logs.Info("ODIN auth parameter lookup ok: parameter=%s path=%s value_set=%t", parameterName, path, strings.TrimSpace(value) != "")

	return strings.TrimSpace(value), nil
}

func autenticarOdin(ctx context.Context, baseURL, username, password, version string) (string, map[string]interface{}) {
	loginURL := unirURL(baseURL, odinLoginPath)
	logs.Info("ODIN auth login request: url=%s username=%s version=%s", loginURL, username, version)
	body := models.OdinAuthRequest{
		Username: username,
		Password: password,
		Version:  version,
	}

	var response models.OdinAuthResponse
	statusCode, err := request.PostWithContext(ctx, loginURL, body, &response)
	if err != nil {
		logs.Error("ODIN auth login fallo: url=%s status=%d err=%v", loginURL, statusCode, err)
		return "", map[string]interface{}{
			"funcion": "/autenticarOdin",
			"err":     fmt.Sprintf("error autenticando el proceso Cargue de Soportes Previnculación: %v", err),
			"status":  fmt.Sprintf("%d", http.StatusBadGateway),
		}
	}
	logs.Info("ODIN auth login ok: url=%s status=%d token_set=%t access_token_set=%t", loginURL, statusCode, strings.TrimSpace(response.Token) != "", strings.TrimSpace(response.AccessToken) != "")

	token := strings.TrimSpace(response.Token)
	if token == "" {
		token = strings.TrimSpace(response.AccessToken)
	}
	if token == "" {
		logs.Error("ODIN auth login sin token: url=%s", loginURL)
		return "", map[string]interface{}{
			"funcion": "/autenticarOdin",
			"err":     "el servicio de Cargue de Soportes Previnculación no retornó token de autenticación",
			"status":  fmt.Sprintf("%d", http.StatusBadGateway),
		}
	}

	return token, nil
}

func consultarRequisitosVinculacion(ctx context.Context, baseURL, token string, datos models.ObjetoPrevinculaciones) ([]models.OdinRequisitoVinculacion, map[string]interface{}) {
	if len(datos.Docentes) == 0 {
		return nil, nil
	}

	body := models.OdinRequisitosVinculacionRequest{
		Parametros: models.OdinRequisitosVinculacionParametros{
			Anio:      fmt.Sprintf("%d", datos.Vigencia),
			Periodo:   strings.TrimSpace(datos.Docentes[0].Periodo),
			IdUsuario: construirIdUsuariosOdin(datos.Docentes),
		},
	}

	url := unirURL(baseURL, odinRequisitosVinculacionPath)
	var response interface{}
	headerAnterior := request.GetHeader()
	request.SetHeader("Bearer " + token)
	defer request.SetHeader(headerAnterior)

	if err := request.SendJson(url, http.MethodPost, &response, body); err != nil {
		return nil, map[string]interface{}{
			"funcion": "/consultarRequisitosVinculacion",
			"err":     fmt.Sprintf("error consultando el proceso Cargue de Soportes Previnculación: %v", err),
			"status":  fmt.Sprintf("%d", http.StatusBadGateway),
		}
	}

	respuestaNormalizada, errMap := normalizarRespuestaRequisitosODIN(response)
	if errMap != nil {
		return nil, errMap
	}

	return respuestaNormalizada, nil
}

func construirIdUsuariosOdin(docentes []models.CargaLectiva) string {
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
		ids = append(ids, fmt.Sprintf("'%s'", documento))
	}

	return strings.Join(ids, ",")
}

func construirDocentesNoVinculables(docentes []models.CargaLectiva, respuesta []models.OdinRequisitoVinculacion) []models.DocenteNoVinculable {
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

	respuestaPorDocumento := make(map[string]models.OdinRequisitoVinculacion, len(respuesta))
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
				Motivos:   []string{odinSinDatosMessage},
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

func construirMotivosNoVinculable(requisito models.OdinRequisitoVinculacion) []string {
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

func unirURL(baseURL, path string) string {
	base := strings.TrimRight(strings.TrimSpace(baseURL), "/")
	segmento := strings.TrimLeft(strings.TrimSpace(path), "/")
	return base + "/" + segmento
}

func normalizarRespuestaRequisitosODIN(raw interface{}) ([]models.OdinRequisitoVinculacion, map[string]interface{}) {
	switch typed := raw.(type) {
	case []interface{}:
		return convertirRespuestaODIN(typed)
	case map[string]interface{}:
		if data, ok := typed["Data"]; ok {
			return normalizarRespuestaRequisitosODIN(data)
		}
		if data, ok := typed["data"]; ok {
			return normalizarRespuestaRequisitosODIN(data)
		}

		detalle := extraerMensajeRespuestaODIN(typed)
		if esRespuestaSinDatosODIN(detalle) {
			return []models.OdinRequisitoVinculacion{}, nil
		}
		if detalle == "" {
			detalle = fmt.Sprintf("respuesta inesperada de ODIN: %+v", typed)
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

func convertirRespuestaODIN(raw interface{}) ([]models.OdinRequisitoVinculacion, map[string]interface{}) {
	bytes, err := json.Marshal(raw)
	if err != nil {
		return nil, map[string]interface{}{
			"funcion": "/consultarRequisitosVinculacion",
			"err":     fmt.Sprintf("error serializando respuesta del proceso Cargue de Soportes Previnculación: %v", err),
			"status":  fmt.Sprintf("%d", http.StatusBadGateway),
		}
	}

	var response []models.OdinRequisitoVinculacion
	if err := json.Unmarshal(bytes, &response); err != nil {
		return nil, map[string]interface{}{
			"funcion": "/consultarRequisitosVinculacion",
			"err":     fmt.Sprintf("error interpretando respuesta del proceso Cargue de Soportes Previnculación: %v", err),
			"status":  fmt.Sprintf("%d", http.StatusBadGateway),
		}
	}

	return response, nil
}

func extraerMensajeRespuestaODIN(data map[string]interface{}) string {
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

func esRespuestaSinDatosODIN(mensaje string) bool {
	normalizado := strings.ToLower(strings.TrimSpace(mensaje))
	return strings.Contains(normalizado, "no se encontrarón datos asociados al reporte") ||
		strings.Contains(normalizado, "no se encontraron datos asociados al reporte")
}
