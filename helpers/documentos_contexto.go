package helpers

import (
	"strconv"
	"time"

	"github.com/udistrital/resoluciones_mid_v2/models"
)

func cargarTipoResolucionDocumento(tipoResolucionID int) (models.Parametro, map[string]interface{}) {
	var tipoResolucion models.Parametro
	if err := GetRequestNew("UrlcrudParametros", ParametroEndpoint+strconv.Itoa(tipoResolucionID), &tipoResolucion); err != nil {
		return models.Parametro{}, map[string]interface{}{"funcion": "/ConstruirDocumentoResolucion-param", "err": err.Error(), "status": "500"}
	}
	return tipoResolucion, nil
}

func cargarOrdenadorResolucion(dependenciaFirmaID int) (models.OrdenadorGasto, map[string]interface{}) {
	var ordenadorGasto models.OrdenadorGasto
	var ordenadoresGasto []models.OrdenadorGasto
	url := "ordenador_gasto?query=DependenciaId:" + strconv.Itoa(dependenciaFirmaID)
	if err := GetRequestLegacy("UrlcrudCore", url, &ordenadoresGasto); err != nil {
		return models.OrdenadorGasto{}, map[string]interface{}{"funcion": "/ConstruirDocumentoResolucion-ordenador", "err": err.Error(), "status": "500"}
	}

	if len(ordenadoresGasto) > 0 {
		ordenadorGasto = ordenadoresGasto[0]
	} else {
		if err := GetRequestLegacy("UrlcrudCore", "ordenador_gasto/1", &ordenadorGasto); err != nil {
			return models.OrdenadorGasto{}, map[string]interface{}{"funcion": "/ConstruirDocumentoResolucion-ordenador", "err": err.Error(), "status": "500"}
		}
	}

	var jefeDependencia []models.JefeDependencia
	fechaActual := time.Now().Format("2006-01-02")
	urlJefe := "jefe_dependencia?query=DependenciaId:" + strconv.Itoa(dependenciaFirmaID) + ",FechaFin__gte:" + fechaActual + ",FechaInicio__lte:" + fechaActual
	if err := GetRequestLegacy("UrlcrudCore", urlJefe, &jefeDependencia); err != nil {
		return models.OrdenadorGasto{}, map[string]interface{}{"funcion": "/ConstruirDocumentoResolucion-jefe", "err": err.Error(), "status": "500"}
	}
	if len(jefeDependencia) == 0 {
		return models.OrdenadorGasto{}, map[string]interface{}{"funcion": "/ConstruirDocumentoResolucion-jefe", "err": "No se encontró jefe para la dependencia en el periodo actual", "status": "500"}
	}

	ordenador, err := BuscarDatosPersonalesDocente(float64(jefeDependencia[0].TerceroId))
	if err != nil {
		return models.OrdenadorGasto{}, map[string]interface{}{"funcion": "/ConstruirDocumentoResolucion-jefe", "err": err, "status": "500"}
	}
	ordenadorGasto.NombreOrdenador = ordenador.NomProveedor

	return ordenadorGasto, nil
}

func cargarFacultadDocumento(facultadID int) (models.Dependencia, map[string]interface{}) {
	var facultad models.Dependencia
	url := "dependencia/" + strconv.Itoa(facultadID)
	if err := GetRequestLegacy("UrlcrudOikos", url, &facultad); err != nil {
		return models.Dependencia{}, map[string]interface{}{"funcion": "/ConstruirTablaVinculaciones-dep", "err": err.Error(), "status": "500"}
	}
	return facultad, nil
}
