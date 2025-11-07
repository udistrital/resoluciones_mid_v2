package helpers

import (
	"fmt"
	"math"
	"strconv"

	"github.com/astaxie/beego/logs"
	"github.com/udistrital/resoluciones_mid_v2/models"
)

func CalcularDesagregadoTitan(v models.VinculacionDocente, dedicacion, nivelAcademico string, objetoNovedad ...*models.ObjetoNovedad) (d map[string]interface{}, outputError map[string]interface{}) {
	defer func() {
		if err := recover(); err != nil {
			outputError = map[string]interface{}{"funcion": "CalcularDesagregadoTitan", "err": err, "status": "500"}
			panic(outputError)
		}
	}()
	fmt.Println("VALOR CONTRATO DES", v.ValorContrato)
	var desagregado map[string]interface{}
	datos := &models.DatosVinculacion{
		NumeroContrato: "",
		Vigencia:       v.Vigencia,
		Documento:      strconv.Itoa(int(v.PersonaId)),
		Dedicacion:     dedicacion,
		Categoria:      v.Categoria,
		NumeroSemanas:  v.NumeroSemanas,
		HorasSemanales: v.NumeroHorasSemanales,
		NivelAcademico: nivelAcademico,
		PuntoSalarial:  v.ValorPuntoSalarial,
	}

	if len(objetoNovedad) > 0 {
		datos.ObjetoNovedad = objetoNovedad[0]
	}

	if nivelAcademico == "POSGRADO" {
		datos.NumeroSemanas = 1
	}

	fmt.Println("DESAGREGADO HCS ", datos)
	if err := SendRequestNew("UrlmidTitan", "desagregado_hcs", "POST", &desagregado, &datos); err != nil {
		logs.Error(err.Error())
		panic("Consultando desagregado -> " + err.Error())
	}

	return desagregado, outputError
}

func EjecutarPreliquidacionTitan(v models.VinculacionDocente) (output map[string]interface{}) {
	defer func() {
		if err := recover(); err != nil {
			output = map[string]interface{}{
				"funcion": "EjecutarPreliquidacionTitan",
				"err":     fmt.Sprintf("%v", err),
				"status":  "500",
			}
			panic(output)
		}
	}()

	if v.NumeroContrato == nil || *v.NumeroContrato == "" {
		return map[string]interface{}{
			"status":  "error",
			"message": fmt.Sprintf("La vinculación %d no tiene número de contrato asignado", v.Id),
		}
	}

	queryContrato := fmt.Sprintf("NumeroContrato:%s,Vigencia:%d,Activo:true", *v.NumeroContrato, v.Vigencia)
	var raw interface{}
	if err := GetRequestNew("TitanCrudService", "contrato?query="+queryContrato, &raw); err == nil {
		if m, ok := raw.(map[string]interface{}); ok {
			if data, ok := m["Data"].([]interface{}); ok && len(data) > 0 {
				return map[string]interface{}{
					"status":  "omitido",
					"message": fmt.Sprintf("Contrato %s (vigencia %d) ya existe en Titan", *v.NumeroContrato, v.Vigencia),
				}
			}
		}
		if arr, ok := raw.([]interface{}); ok && len(arr) > 0 {
			return map[string]interface{}{
				"status":  "omitido",
				"message": fmt.Sprintf("Contrato %s (vigencia %d) ya existe en Titan", *v.NumeroContrato, v.Vigencia),
			}
		}
	}

	if v.NumeroRp != 0 {
		queryRp := fmt.Sprintf("Rp:%d,Vigencia:%d,Activo:true", int(v.NumeroRp), v.Vigencia)
		var rawRp interface{}
		if err := GetRequestNew("TitanCrudService", "contrato?query="+queryRp, &rawRp); err == nil {
			if m, ok := rawRp.(map[string]interface{}); ok {
				if data, ok := m["Data"].([]interface{}); ok && len(data) > 0 {
					return map[string]interface{}{
						"status":  "omitido",
						"message": fmt.Sprintf("Contrato omitido: RP %d (vigencia %d) ya existe en Titan", int(v.NumeroRp), v.Vigencia),
					}
				}
			}
			if arr, ok := rawRp.([]interface{}); ok && len(arr) > 0 {
				return map[string]interface{}{
					"status":  "omitido",
					"message": fmt.Sprintf("Contrato omitido: RP %d (vigencia %d) ya existe en Titan", int(v.NumeroRp), v.Vigencia),
				}
			}
		}
	}

	var c models.ContratoPreliquidacion
	var desagregado []models.DisponibilidadVinculacion
	var docente []models.InformacionProveedor
	var actaInicio []models.ActaInicio

	resolucion := GetResolucion(v.ResolucionVinculacionDocenteId.Id)
	resVin := GetResolucionVinculacionDocente(v.ResolucionVinculacionDocenteId.Id)

	preliquidacion := &models.ContratoPreliquidacion{
		NumeroContrato: *v.NumeroContrato,
		Vigencia:       v.Vigencia,
		Documento:      fmt.Sprintf("%.f", v.PersonaId),
		DependenciaId:  resolucion.DependenciaId,
		Rp:             int(v.NumeroRp),
		TipoNominaId:   410,
		Activo:         true,
	}

	url := "disponibilidad_vinculacion?query=Activo:true,VinculacionDocenteId.Id:" + strconv.Itoa(v.Id)
	if err := GetRequestNew("UrlcrudResoluciones", url, &desagregado); err != nil {
		panic("Error al cargar desagregado-preliquidación: " + err.Error())
	}

	desagregadoMap := map[string]float64{}
	if resVin.Dedicacion == "HCH" {
		preliquidacion.ValorContrato = math.Floor(v.ValorContrato)
		preliquidacion.TipoNominaId = 409
	} else {
		for _, d := range desagregado {
			if d.Rubro == "SueldoBasico" {
				preliquidacion.ValorContrato = d.Valor
			} else {
				desagregadoMap[d.Rubro] = d.Valor
			}
		}
		preliquidacion.Desagregado = &desagregadoMap
	}

	url2 := "informacion_proveedor?query=NumDocumento:" + preliquidacion.Documento
	if err := GetRequestLegacy("UrlcrudAgora", url2, &docente); err != nil {
		panic("Error consultando docente en Ágora: " + err.Error())
	} else if len(docente) == 0 {
		panic("No se encontró información del docente en Ágora")
	}

	url3 := "acta_inicio?query=NumeroContrato:" + *v.NumeroContrato + ",Vigencia:" + strconv.Itoa(v.Vigencia)
	if err := GetRequestLegacy("UrlcrudAgora", url3, &actaInicio); err != nil {
		panic("Error consultando acta de inicio: " + err.Error())
	} else if len(actaInicio) == 0 {
		panic("No se encontró acta de inicio asociada")
	}

	preliquidacion.FechaInicio = actaInicio[0].FechaInicio
	preliquidacion.FechaFin = actaInicio[0].FechaFin
	preliquidacion.Cdp = desagregado[0].Disponibilidad
	preliquidacion.NombreCompleto = docente[0].NomProveedor
	preliquidacion.PersonaId = docente[0].Id
	preliquidacion.NumeroSemanas = v.NumeroSemanas
	preliquidacion.ResolucionId = v.ResolucionVinculacionDocenteId.Id
	preliquidacion.Resolucion = resolucion.NumeroResolucion

	if err := SendRequestNew("UrlmidTitan", "preliquidacion", "POST", &c, &preliquidacion); err != nil {
		panic("Error enviando preliquidación a Titan: " + err.Error())
	}

	return map[string]interface{}{
		"status": "ok",
		"message": fmt.Sprintf("Preliquidación enviada correctamente a Titan para contrato %s (vigencia %d, RP %d)",
			*v.NumeroContrato, v.Vigencia, int(v.NumeroRp)),
	}
}

func ReliquidarContratoCancelado(cancelacion models.VinculacionDocente, cancelado models.VinculacionDocente) (outputError map[string]interface{}) {
	defer func() {
		if err := recover(); err != nil {
			outputError = map[string]interface{}{"funcion": "ReliquidarContratoCancelado", "err": err, "status": "500"}
			panic(outputError)
		}
	}()
	var c models.ContratoPreliquidacion
	var desagregado, err map[string]interface{}
	valores := make(map[string]float64)
	var sueldoBasico float64

	contratoReliquidar := &models.ContratoCancelacion{
		NumeroContrato: *cancelado.NumeroContrato,
		Vigencia:       cancelado.Vigencia,
		Semanas:        cancelado.NumeroSemanas - cancelacion.NumeroSemanas,
		FechaAnulacion: cancelacion.FechaInicio,
		Documento:      strconv.Itoa(int(cancelacion.PersonaId)),
	}

	if cancelado.ResolucionVinculacionDocenteId.Dedicacion != "HCH" {
		cancelado.NumeroSemanas -= cancelacion.NumeroSemanas
		dedicacion := cancelado.ResolucionVinculacionDocenteId.Dedicacion
		nivel := cancelado.ResolucionVinculacionDocenteId.NivelAcademico
		if nivel == "POSGRADO" {
			cancelado.NumeroHorasSemanales -= cancelacion.NumeroHorasSemanales
		} else {
			cancelado.NumeroHorasSemanales = cancelacion.NumeroHorasSemanales
		}
		if desagregado, err = CalcularDesagregadoTitan(cancelado, dedicacion, nivel); err != nil {
			panic(err)
		}

		for concepto, valor := range desagregado {
			if concepto != "NumeroContrato" && concepto != "Vigencia" && concepto != "SueldoBasico" {
				valores[concepto] = valor.(float64)
			}
			if concepto == "SueldoBasico" {
				sueldoBasico = valor.(float64)
			}
		}
		contratoReliquidar.ValorContrato = sueldoBasico
		contratoReliquidar.Desagregado = &valores
	}

	contratoReliquidar.NivelAcademico = cancelado.ResolucionVinculacionDocenteId.NivelAcademico
	if err2 := SendRequestNew("UrlmidTitan", "novedadVE/aplicar_anulacion", "POST", &c, &contratoReliquidar); err2 != nil {
		panic("Reliquidando -> " + err2.Error())
	}

	return
}

func ReducirContratosTitan(reduccion *models.Reduccion, modificacion *models.VinculacionDocente, valorReduccion float64) (outputError map[string]interface{}) {
	defer func() {
		if err := recover(); err != nil {
			outputError = map[string]interface{}{"funcion": "ReducirContratosTitan", "err": err, "status": "500"}
			panic(outputError)
		}
	}()
	var c models.ContratoPreliquidacion

	if reduccion.ContratoNuevo != nil && reduccion.ContratoNuevo.ValorContratoReduccion == 0 {
		reduccion.ContratoNuevo.ValorContratoReduccion = valorReduccion
	}
	if err2 := SendRequestNew("UrlmidTitan", "novedadVE/aplicar_reduccion", "POST", &c, &reduccion); err2 != nil {
		panic("Reliquidando -> " + err2.Error())
	}

	return nil
}
