package services

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/astaxie/beego"
	"github.com/udistrital/resoluciones_mid_v2/helpers"
)

type decanoFacultadXML struct {
	Decano struct {
		CodigoFacultad int `xml:"codigo_facultad"`
	} `xml:"decano"`
}

type asistenteFacultadXML struct {
	Facultad struct {
		CodigoDependencia int `xml:"codigo_dependecia"`
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

func resolveCodigoDependencia(numeroDocumento, rol string) (int, map[string]interface{}) {
	protocol := beego.AppConfig.String("ProtocolAdmin")
	baseWSO2 := beego.AppConfig.String("UrlcrudWSO2")
	nsAcademica := beego.AppConfig.String("NscrudAcademica")

	rol = strings.ToUpper(strings.TrimSpace(rol))

	switch rol {
	case "DECANO":
		var dec decanoFacultadXML
		url := joinWSO2URL(protocol, baseWSO2, nsAcademica, "decano/"+numeroDocumento)

		if err := helpers.GetXml(url, &dec); err != nil {
			return 0, map[string]interface{}{
				"funcion": "resolveCodigoDependencia:decano",
				"err":     err.Error(),
				"status":  "502",
			}
		}

		if dec.Decano.CodigoFacultad == 0 {
			return 0, map[string]interface{}{
				"funcion": "resolveCodigoDependencia:decano",
				"err":     "no se encontró una facultad activa para el usuario en el SGA",
				"status":  "404",
			}
		}

		return dec.Decano.CodigoFacultad, nil

	case "ASISTENTE_DECANATURA":
		var asis asistenteFacultadXML
		url := joinWSO2URL(protocol, baseWSO2, nsAcademica, "asistente_facultad/"+numeroDocumento)

		if err := helpers.GetXml(url, &asis); err != nil {
			return 0, map[string]interface{}{
				"funcion": "resolveCodigoDependencia:asistente_decanatura",
				"err":     err.Error(),
				"status":  "502",
			}
		}

		if asis.Facultad.CodigoDependencia == 0 {
			return 0, map[string]interface{}{
				"funcion": "resolveCodigoDependencia:asistente_decanatura",
				"err":     "no se encontró una dependencia activa para el usuario en el SGA",
				"status":  "404",
			}
		}

		return asis.Facultad.CodigoDependencia, nil

	default:
		return 0, map[string]interface{}{
			"funcion": "resolveCodigoDependencia",
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

func ResolveOikosByRol(numeroDocumento, rol string) ([]int, map[string]interface{}) {
	codigoDependencia, err := resolveCodigoDependencia(numeroDocumento, rol)
	if err != nil {
		return nil, err
	}

	idOikos, err := resolveIdOikosFromHomologacion(codigoDependencia)
	if err != nil {
		return nil, err
	}

	return []int{idOikos}, nil
}
