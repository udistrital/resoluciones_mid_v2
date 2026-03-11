package services

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/astaxie/beego"
	"github.com/udistrital/resoluciones_mid_v2/models"
	"github.com/udistrital/utils_oas/request"
)

var rolePriority = map[string]int{
	"ADMINISTRADOR_RESOLUCIONES": 3,
	"ASIS_FINANCIERA":            2,
	"DECANO":                     1,
	"ASISTENTE_DECANATURA":       1,
}

func normalizeBaseNoProto(u string) string {
	u = strings.TrimSpace(u)
	u = strings.TrimLeft(u, "/")
	return u
}

func normalizeRol(rol string) string {
	return strings.ToUpper(strings.TrimSpace(rol))
}

func joinWSO2URL(protocol, base, ns, path string) string {
	protocol = strings.TrimRight(protocol, "://")
	base = strings.TrimRight(normalizeBaseNoProto(base), "/")
	ns = strings.Trim(ns, "/")
	path = strings.TrimLeft(path, "/")
	return fmt.Sprintf("%s://%s/%s/%s", protocol, base, ns, path)
}

func deduplicateDependencias(items []models.DependenciaUsuario) []models.DependenciaUsuario {
	seen := make(map[string]bool)
	result := make([]models.DependenciaUsuario, 0)

	for _, item := range items {
		key := fmt.Sprintf("%d-%d", item.CodigoDependencia, item.IdOikos)
		if !seen[key] {
			seen[key] = true
			result = append(result, item)
		}
	}

	return result
}

func getHighestPriorityRol(roles []string) string {
	bestRol := ""
	bestPriority := -1

	for _, rol := range roles {
		r := normalizeRol(rol)
		priority, ok := rolePriority[r]
		if !ok {
			continue
		}

		if priority > bestPriority {
			bestPriority = priority
			bestRol = r
		}
	}

	return bestRol
}

func isGlobalRol(rol string) bool {
	switch normalizeRol(rol) {
	case "ADMINISTRADOR_RESOLUCIONES", "ASIS_FINANCIERA":
		return true
	default:
		return false
	}
}

func getJSONWithUtilOAS(url string, target interface{}) map[string]interface{} {
	if err := request.GetJson(url, target); err != nil {
		return map[string]interface{}{
			"funcion": "getJSONWithUtilOAS",
			"err":     err.Error(),
			"status":  "502",
		}
	}

	return nil
}

func getJSONWithHTTP(url string, target interface{}) map[string]interface{} {
	client := &http.Client{
		Timeout: 20 * time.Second,
	}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return map[string]interface{}{
			"funcion": "getJSONWithHTTP:newRequest",
			"err":     err.Error(),
			"status":  "500",
		}
	}

	req.Header.Set("Accept", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return map[string]interface{}{
			"funcion": "getJSONWithHTTP:do",
			"err":     err.Error(),
			"status":  "502",
		}
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return map[string]interface{}{
			"funcion": "getJSONWithHTTP:readBody",
			"err":     err.Error(),
			"status":  "502",
		}
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return map[string]interface{}{
			"funcion": "getJSONWithHTTP:statusCode",
			"err":     fmt.Sprintf("respuesta no exitosa del servicio externo: %d - %s", resp.StatusCode, strings.TrimSpace(string(body))),
			"status":  strconv.Itoa(resp.StatusCode),
		}
	}

	if err := json.Unmarshal(body, target); err != nil {
		return map[string]interface{}{
			"funcion": "getJSONWithHTTP:unmarshal",
			"err":     err.Error(),
			"status":  "502",
		}
	}

	return nil
}

func getJSON(url string, target interface{}) map[string]interface{} {
	if errMap := getJSONWithUtilOAS(url, target); errMap == nil {
		return nil
	}

	return getJSONWithHTTP(url, target)
}

func resolveDependenciasFromSGA(numeroDocumento, rol string) ([]models.DependenciaUsuario, map[string]interface{}) {
	protocol := beego.AppConfig.String("ProtocolAdmin")
	baseWSO2 := beego.AppConfig.String("UrlcrudWSO2")
	nsAcademica := beego.AppConfig.String("NscrudAcademica")

	rol = normalizeRol(rol)

	switch rol {
	case "DECANO":
		var dec models.DecanoFacultadResponse
		url := joinWSO2URL(protocol, baseWSO2, nsAcademica, "decano/"+numeroDocumento)

		if errMap := getJSON(url, &dec); errMap != nil {
			errMap["funcion"] = "resolveDependenciasFromSGA:decano"
			return nil, errMap
		}

		if len(dec.Facultad.Decano) == 0 {
			return nil, map[string]interface{}{
				"funcion": "resolveDependenciasFromSGA:decano",
				"err":     "no se encontró una facultad activa para el usuario en el SGA",
				"status":  "404",
			}
		}

		dependencias := make([]models.DependenciaUsuario, 0)
		for _, item := range dec.Facultad.Decano {
			codigo, err := strconv.Atoi(strings.TrimSpace(item.CodigoFacultad))
			if err != nil || codigo <= 0 {
				continue
			}

			dependencias = append(dependencias, models.DependenciaUsuario{
				CodigoDependencia: codigo,
				Nombre:            item.NombreFacultad,
				Rol:               rol,
			})
		}

		if len(dependencias) == 0 {
			return nil, map[string]interface{}{
				"funcion": "resolveDependenciasFromSGA:decano",
				"err":     "no se encontró una facultad activa para el usuario en el SGA",
				"status":  "404",
			}
		}

		return deduplicateDependencias(dependencias), nil

	case "ASISTENTE_DECANATURA":
		var asis models.AsistenteFacultadResponse
		url := joinWSO2URL(protocol, baseWSO2, nsAcademica, "asistente_facultad/"+numeroDocumento)

		if errMap := getJSON(url, &asis); errMap != nil {
			errMap["funcion"] = "resolveDependenciasFromSGA:asistente_decanatura"
			return nil, errMap
		}

		if len(asis.Asistente.Facultad) == 0 {
			return nil, map[string]interface{}{
				"funcion": "resolveDependenciasFromSGA:asistente_decanatura",
				"err":     "no se encontró una dependencia activa para el usuario en el SGA",
				"status":  "404",
			}
		}

		dependencias := make([]models.DependenciaUsuario, 0)
		for _, item := range asis.Asistente.Facultad {
			codigo, err := strconv.Atoi(strings.TrimSpace(item.CodigoDependencia))
			if err != nil || codigo <= 0 {
				continue
			}

			dependencias = append(dependencias, models.DependenciaUsuario{
				CodigoDependencia: codigo,
				Nombre:            item.NombreDependencia,
				Rol:               rol,
			})
		}

		if len(dependencias) == 0 {
			return nil, map[string]interface{}{
				"funcion": "resolveDependenciasFromSGA:asistente_decanatura",
				"err":     "no se encontró una dependencia activa para el usuario en el SGA",
				"status":  "404",
			}
		}

		return deduplicateDependencias(dependencias), nil

	default:
		return nil, map[string]interface{}{
			"funcion": "resolveDependenciasFromSGA",
			"err":     "rol no soportado",
			"status":  "400",
		}
	}
}

func resolveIdOikosFromHomologacion(codigoDependencia int) (int, map[string]interface{}) {
	protocol := beego.AppConfig.String("ProtocolAdmin")
	baseWSO2 := beego.AppConfig.String("UrlcrudWSO2")
	nsHomologacion := beego.AppConfig.String("NscrudHomologacion")

	var hom models.HomologacionFacultadResponse
	url := joinWSO2URL(protocol, baseWSO2, nsHomologacion, "facultad_oikos_gedep/"+strconv.Itoa(codigoDependencia))

	if errMap := getJSON(url, &hom); errMap != nil {
		errMap["funcion"] = "resolveIdOikosFromHomologacion"
		return 0, errMap
	}

	if strings.TrimSpace(hom.Homologacion.IdOikos) == "" {
		return 0, map[string]interface{}{
			"funcion": "resolveIdOikosFromHomologacion",
			"err":     "no se encontró homologación Oikos para la dependencia consultada",
			"status":  "404",
		}
	}

	idOikos, err := strconv.Atoi(strings.TrimSpace(hom.Homologacion.IdOikos))
	if err != nil || idOikos <= 0 {
		return 0, map[string]interface{}{
			"funcion": "resolveIdOikosFromHomologacion",
			"err":     "el id_oikos recibido no es válido",
			"status":  "502",
		}
	}

	return idOikos, nil
}

func ResolveDependenciasByRol(numeroDocumento, rol string) ([]models.DependenciaUsuario, map[string]interface{}) {
	dependencias, err := resolveDependenciasFromSGA(numeroDocumento, rol)
	if err != nil {
		return nil, err
	}

	resultado := make([]models.DependenciaUsuario, 0)

	for _, dep := range dependencias {
		idOikos, errMap := resolveIdOikosFromHomologacion(dep.CodigoDependencia)
		if errMap != nil {
			return nil, errMap
		}

		dep.IdOikos = idOikos
		resultado = append(resultado, dep)
	}

	return deduplicateDependencias(resultado), nil
}

func ResolveAlcanceUsuario(numeroDocumento string, roles []string) (models.AlcanceUsuario, map[string]interface{}) {
	rolPrincipal := getHighestPriorityRol(roles)

	if rolPrincipal == "" {
		return models.AlcanceUsuario{}, map[string]interface{}{
			"funcion": "ResolveAlcanceUsuario",
			"err":     "el usuario no tiene roles soportados",
			"status":  "400",
		}
	}

	if isGlobalRol(rolPrincipal) {
		return models.AlcanceUsuario{
			RolPrincipal: rolPrincipal,
			EsGlobal:     true,
			Dependencias: []models.DependenciaUsuario{},
		}, nil
	}

	dependencias, err := ResolveDependenciasByRol(numeroDocumento, rolPrincipal)
	if err != nil {
		return models.AlcanceUsuario{}, err
	}

	return models.AlcanceUsuario{
		RolPrincipal: rolPrincipal,
		EsGlobal:     false,
		Dependencias: dependencias,
	}, nil
}
