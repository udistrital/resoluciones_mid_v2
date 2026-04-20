package services

import (
	"strconv"
	"strings"

	"github.com/astaxie/beego"
	"github.com/udistrital/resoluciones_mid_v2/models"
)

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
