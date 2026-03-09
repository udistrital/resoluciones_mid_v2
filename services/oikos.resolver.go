package services

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/astaxie/beego"
	"github.com/udistrital/resoluciones_mid_v2/helpers"
)

type DependenciaUsuario struct {
	CodigoDependencia int    `json:"codigo_dependencia"`
	IdOikos           int    `json:"id_oikos"`
	Nombre            string `json:"nombre,omitempty"`
	Rol               string `json:"rol,omitempty"`
}

type decanoFacultadXML struct {
	Decanos []struct {
		CodigoFacultad int    `xml:"codigo_facultad"`
		NombreFacultad string `xml:"facultad"`
	} `xml:"decano"`
}

type asistenteFacultadXML struct {
	Facultades []struct {
		CodigoDependencia int    `xml:"codigo_dependecia"`
		NombreDependencia string `xml:"nombre_dependencia"`
	} `xml:"facultad"`
}
type homologacionFacultadXML struct {
	IdOikos int `xml:"id_oikos"`
	IdGedep int `xml:"id_gedep"`
}

func normalizeBaseNoProto(u string) string {
	u = strings.TrimSpace(u)
	u = strings.TrimLeft(u, "/")
	return u
}

func joinWSO2URL(protocol, base, ns, path string) string {
	protocol = strings.TrimRight(protocol, "://")
	base = strings.TrimRight(normalizeBaseNoProto(base), "/")
	ns = strings.Trim(ns, "/")
	path = strings.TrimLeft(path, "/")
	return fmt.Sprintf("%s://%s/%s/%s", protocol, base, ns, path)
}

func deduplicateDependencias(items []DependenciaUsuario) []DependenciaUsuario {
	seen := make(map[string]bool)
	result := make([]DependenciaUsuario, 0)

	for _, item := range items {
		key := fmt.Sprintf("%d-%d", item.CodigoDependencia, item.IdOikos)
		if !seen[key] {
			seen[key] = true
			result = append(result, item)
		}
	}

	return result
}

func resolveDependenciasFromSGA(numeroDocumento, rol string) ([]DependenciaUsuario, map[string]interface{}) {
	protocol := beego.AppConfig.String("ProtocolAdmin")
	baseWSO2 := beego.AppConfig.String("UrlcrudWSO2")
	nsAcademica := beego.AppConfig.String("NscrudAcademica")

	rol = strings.ToUpper(strings.TrimSpace(rol))

	switch rol {

	case "DECANO":
		var dec decanoFacultadXML
		url := joinWSO2URL(protocol, baseWSO2, nsAcademica, "decano/"+numeroDocumento)

		if err := helpers.GetXml(url, &dec); err != nil {
			return nil, map[string]interface{}{
				"funcion": "resolveDependenciasFromSGA:decano",
				"err":     err.Error(),
				"status":  "502",
			}
		}

		if len(dec.Decanos) == 0 {
			return nil, map[string]interface{}{
				"funcion": "resolveDependenciasFromSGA:decano",
				"err":     "no se encontró una facultad activa para el usuario en el SGA",
				"status":  "404",
			}
		}

		dependencias := make([]DependenciaUsuario, 0)
		for _, item := range dec.Decanos {
			if item.CodigoFacultad > 0 {
				dependencias = append(dependencias, DependenciaUsuario{
					CodigoDependencia: item.CodigoFacultad,
					Nombre:            item.NombreFacultad,
					Rol:               rol,
				})
			}
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
		var asis asistenteFacultadXML
		url := joinWSO2URL(protocol, baseWSO2, nsAcademica, "asistente_facultad/"+numeroDocumento)

		if err := helpers.GetXml(url, &asis); err != nil {
			return nil, map[string]interface{}{
				"funcion": "resolveDependenciasFromSGA:asistente_decanatura",
				"err":     err.Error(),
				"status":  "502",
			}
		}

		if len(asis.Facultades) == 0 {
			return nil, map[string]interface{}{
				"funcion": "resolveDependenciasFromSGA:asistente_decanatura",
				"err":     "no se encontró una dependencia activa para el usuario en el SGA",
				"status":  "404",
			}
		}

		dependencias := make([]DependenciaUsuario, 0)
		for _, item := range asis.Facultades {
			if item.CodigoDependencia > 0 {
				dependencias = append(dependencias, DependenciaUsuario{
					CodigoDependencia: item.CodigoDependencia,
					Nombre:            item.NombreDependencia,
					Rol:               rol,
				})
			}
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

	var hom homologacionFacultadXML
	url := joinWSO2URL(protocol, baseWSO2, nsHomologacion, "facultad_oikos_gedep/"+strconv.Itoa(codigoDependencia))

	if err := helpers.GetXml(url, &hom); err != nil {
		return 0, map[string]interface{}{
			"funcion": "resolveIdOikosFromHomologacion",
			"err":     err.Error(),
			"status":  "502",
		}
	}

	if hom.IdOikos == 0 {
		return 0, map[string]interface{}{
			"funcion": "resolveIdOikosFromHomologacion",
			"err":     "no se encontró homologación Oikos para la dependencia consultada",
			"status":  "404",
		}
	}

	return hom.IdOikos, nil
}

func ResolveDependenciasByRol(numeroDocumento, rol string) ([]DependenciaUsuario, map[string]interface{}) {
	dependencias, err := resolveDependenciasFromSGA(numeroDocumento, rol)
	if err != nil {
		return nil, err
	}

	resultado := make([]DependenciaUsuario, 0)

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
