package services

import (
	"sort"
	"strconv"
	"strings"
	"sync"

	"github.com/udistrital/resoluciones_mid_v2/helpers"
	"github.com/udistrital/resoluciones_mid_v2/models"
)

type resultadoResumenResolucion struct {
	resumen *models.ResumenDashboardResolucion
	err     map[string]interface{}
}

func ConsultarDashboardResoluciones(numeroDocumento string, roles []string, vigencia int, dependenciaFiltro *int, limit int, offset int) (*models.RespuestaDashboardResoluciones, map[string]interface{}) {
	if limit <= 0 {
		limit = 10
	}
	if offset < 0 {
		offset = 0
	}

	resoluciones, errRes := GetResolucionesByAlcance(numeroDocumento, roles, vigencia, dependenciaFiltro)
	if errRes != nil {
		return nil, errRes
	}

	sort.SliceStable(resoluciones, func(i, j int) bool {
		return resoluciones[i].Id > resoluciones[j].Id
	})

	total := len(resoluciones)

	if offset > total {
		offset = total
	}

	end := offset + limit
	if end > total {
		end = total
	}

	resolucionesPaginadas := []models.Resolucion{}
	if offset < end {
		resolucionesPaginadas = resoluciones[offset:end]
	}

	resultado := make([]models.ResumenDashboardResolucion, 0, len(resolucionesPaginadas))

	resolucionesCompletas := 0
	resolucionesConPendientesTitan := 0
	resolucionesConSinRp := 0

	cacheDependencias := make(map[int]string)
	var cacheDependenciasMu sync.RWMutex

	const maxDashboardWorkers = 6

	maxWorkers := maxDashboardWorkers

	if len(resolucionesPaginadas) < maxWorkers {
		maxWorkers = len(resolucionesPaginadas)
	}
	if maxWorkers <= 0 {
		maxWorkers = 1
	}

	type trabajoResolucion struct {
		index      int
		resolucion models.Resolucion
	}

	trabajos := make(chan trabajoResolucion, len(resolucionesPaginadas))
	resultados := make(chan resultadoResumenResolucion, len(resolucionesPaginadas))

	var wg sync.WaitGroup

	for i := 0; i < maxWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for trabajo := range trabajos {
				resumen, errMap := construirResumenDashboardResolucion(
					trabajo.resolucion,
					cacheDependencias,
					&cacheDependenciasMu,
				)
				resultados <- resultadoResumenResolucion{
					resumen: resumen,
					err:     errMap,
				}
			}
		}()
	}

	for i, resolucion := range resolucionesPaginadas {
		trabajos <- trabajoResolucion{
			index:      i,
			resolucion: resolucion,
		}
	}
	close(trabajos)

	wg.Wait()
	close(resultados)

	for item := range resultados {
		if item.err != nil {
			return nil, item.err
		}
		if item.resumen == nil {
			continue
		}
		if item.resumen.Total == 0 {
			continue
		}

		switch item.resumen.EstadoGeneralCodigo {
		case "COMPLETA":
			resolucionesCompletas++
		case "CON_PENDIENTES_TITAN":
			resolucionesConPendientesTitan++
		case "CON_SIN_RP":
			resolucionesConSinRp++
		}

		resultado = append(resultado, *item.resumen)
	}

	sort.SliceStable(resultado, func(i, j int) bool {
		pi := prioridadEstadoDashboard(resultado[i].EstadoGeneralCodigo)
		pj := prioridadEstadoDashboard(resultado[j].EstadoGeneralCodigo)

		if pi != pj {
			return pi < pj
		}

		return resultado[i].ResolucionId > resultado[j].ResolucionId
	})

	totalPagina := len(resultado)

	porcentajeResolucionesCompletas := 0.0
	if totalPagina > 0 {
		porcentajeResolucionesCompletas = (float64(resolucionesCompletas) / float64(totalPagina)) * 100
	}

	respuesta := &models.RespuestaDashboardResoluciones{
		ResumenGlobal: models.ResumenGlobalDashboardResoluciones{
			TotalResoluciones:               totalPagina,
			ResolucionesCompletas:           resolucionesCompletas,
			ResolucionesConPendientesTitan:  resolucionesConPendientesTitan,
			ResolucionesConSinRp:            resolucionesConSinRp,
			PorcentajeResolucionesCompletas: porcentajeResolucionesCompletas,
		},
		Resoluciones: resultado,
		Total:        total,
		Limit:        limit,
		Offset:       offset,
	}

	return respuesta, nil
}

func construirResumenDashboardResolucion(
	resolucion models.Resolucion,
	cacheDependencias map[int]string,
	cacheDependenciasMu *sync.RWMutex,
) (*models.ResumenDashboardResolucion, map[string]interface{}) {
	vinculaciones, errVin := helpers.Previnculaciones(strconv.Itoa(resolucion.Id))
	if errVin != nil {
		return nil, errVin
	}

	contratosTitan, errTitan := helpers.ObtenerContratosTitanPorResolucion(resolucion.Id)
	if errTitan != nil {
		return nil, errTitan
	}

	dependenciaNombre := obtenerNombreDependencia(resolucion.DependenciaId, cacheDependencias, cacheDependenciasMu)

	mapaTitan := make(map[string]int)
	for _, contrato := range contratosTitan {
		contratoNumero := strings.TrimSpace(contrato.NumeroContrato)
		if contratoNumero != "" && contrato.Activo {
			mapaTitan[contratoNumero]++
		}
	}

	total := 0
	totalConRp := 0
	completas := 0
	pendientesTitan := 0
	sinRp := 0

	for _, vinculacion := range vinculaciones {
		total++

		numeroContrato := ""
		if vinculacion.NumeroContrato != nil {
			numeroContrato = strings.TrimSpace(*vinculacion.NumeroContrato)
		}

		tieneRpResoluciones := vinculacion.NumeroRp > 0 &&
			vinculacion.VigenciaRp > 0 &&
			numeroContrato != ""

		if !tieneRpResoluciones {
			sinRp++
			continue
		}

		totalConRp++

		if mapaTitan[numeroContrato] > 0 {
			completas++
		} else {
			pendientesTitan++
		}
	}

	porcentajeCompletas := 0.0
	porcentajePendientesTitan := 0.0
	porcentajeSinRp := 0.0

	if total > 0 {
		porcentajeCompletas = (float64(completas) / float64(total)) * 100
		porcentajePendientesTitan = (float64(pendientesTitan) / float64(total)) * 100
		porcentajeSinRp = (float64(sinRp) / float64(total)) * 100
	}

	estadoGeneralCodigo := "SIN_VINCULACIONES"
	estadoGeneralNombre := "La resolución no tiene vinculaciones"

	if total > 0 {
		switch {
		case sinRp > 0:
			estadoGeneralCodigo = "CON_SIN_RP"
			estadoGeneralNombre = "Tiene vinculaciones sin RP"
		case pendientesTitan > 0:
			estadoGeneralCodigo = "CON_PENDIENTES_TITAN"
			estadoGeneralNombre = "Tiene vinculaciones pendientes en Titan"
		case completas == total:
			estadoGeneralCodigo = "COMPLETA"
			estadoGeneralNombre = "Todas las vinculaciones están completas"
		default:
			estadoGeneralCodigo = "EN_GESTION"
			estadoGeneralNombre = "Resolución en gestión"
		}
	}

	resumen := &models.ResumenDashboardResolucion{
		ResolucionId:              resolucion.Id,
		NumeroResolucion:          resolucion.NumeroResolucion,
		Vigencia:                  resolucion.Vigencia,
		DependenciaId:             resolucion.DependenciaId,
		DependenciaNombre:         dependenciaNombre,
		Total:                     total,
		TotalConRp:                totalConRp,
		Completas:                 completas,
		PendientesTitan:           pendientesTitan,
		SinRp:                     sinRp,
		PorcentajeCompletas:       porcentajeCompletas,
		PorcentajePendientesTitan: porcentajePendientesTitan,
		PorcentajeSinRp:           porcentajeSinRp,
		EstadoGeneralCodigo:       estadoGeneralCodigo,
		EstadoGeneralNombre:       estadoGeneralNombre,
	}

	return resumen, nil
}

func obtenerNombreDependencia(
	dependenciaId int,
	cacheDependencias map[int]string,
	cacheDependenciasMu *sync.RWMutex,
) string {
	if dependenciaId <= 0 {
		return ""
	}

	cacheDependenciasMu.RLock()
	if nombre, ok := cacheDependencias[dependenciaId]; ok {
		cacheDependenciasMu.RUnlock()
		return nombre
	}
	cacheDependenciasMu.RUnlock()

	var dependencia models.Dependencia
	nombre := ""

	if err := helpers.GetRequestLegacy("UrlcrudOikos", "dependencia/"+strconv.Itoa(dependenciaId), &dependencia); err == nil {
		nombre = dependencia.Nombre
	}

	cacheDependenciasMu.Lock()
	cacheDependencias[dependenciaId] = nombre
	cacheDependenciasMu.Unlock()

	return nombre
}

func prioridadEstadoDashboard(estado string) int {
	switch estado {
	case "CON_SIN_RP":
		return 1
	case "CON_PENDIENTES_TITAN":
		return 2
	case "COMPLETA":
		return 3
	default:
		return 4
	}
}
