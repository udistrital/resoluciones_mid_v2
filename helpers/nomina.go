package helpers

import (
	"fmt"
	"strconv"

	"github.com/astaxie/beego/logs"
	"github.com/udistrital/resoluciones_mid_v2/models"
)

// Envia a Titan la información necesaria para calcular el valor de un contrato desagregado por rubros
func CalcularDesagregadoTitan(v models.VinculacionDocente, dedicacion, nivelAcademico string) (d map[string]interface{}, outputError map[string]interface{}) {
	defer func() {
		if err := recover(); err != nil {
			outputError = map[string]interface{}{"funcion": "CalcularDesagregadoTitan", "err": err, "status": "500"}
			panic(outputError)
		}
	}()
	fmt.Println("VALOR CONTRATO DES", v.ValorContrato)
	// var desagregado models.DesagregadoContrato
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

// Envía a Titan la información para la preliquidación de nómina para los docentes recien contratados con RP actualizado
func EjecutarPreliquidacionTitan(v models.VinculacionDocente) (outputError map[string]interface{}) {
	defer func() {
		if err := recover(); err != nil {
			outputError = map[string]interface{}{"funcion": "EjecutarPreliquidacionTitan", "err": err, "status": "500"}
			panic(outputError)
		}
	}()

	var c models.ContratoPreliquidacion
	var desagregado []models.DisponibilidadVinculacion
	var docente []models.InformacionProveedor
	var actaInicio []models.ActaInicio

	preliquidacion := &models.ContratoPreliquidacion{
		NumeroContrato: *v.NumeroContrato,
		Vigencia:       v.Vigencia,
		Documento:      fmt.Sprintf("%.f", v.PersonaId),
		DependenciaId:  v.ResolucionVinculacionDocenteId.FacultadId,
		Rp:             int(v.NumeroRp),
		TipoNominaId:   410,
		Activo:         true,
	}

	url := "disponibilidad_vinculacion?query=Activo:true,VinculacionDocenteId.Id:" + strconv.Itoa(v.Id)
	if err := GetRequestNew("UrlcrudResoluciones", url, &desagregado); err != nil {
		panic("Cargando desagregado-preliq -> " + err.Error())
	}

	desagregadoMap := map[string]float64{}
	if v.ResolucionVinculacionDocenteId.Dedicacion == "HCH" {
		preliquidacion.ValorContrato = v.ValorContrato
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
	if err2 := GetRequestLegacy("UrlcrudAgora", url2, &docente); err2 != nil {
		panic("Info docente -> " + err2.Error())
	} else if len(docente) == 0 {
		panic("Info docente -> No se encontró información del docente en Agora!!")
	}

	url3 := "acta_inicio?query=NumeroContrato:" + *v.NumeroContrato + ",Vigencia:" + strconv.Itoa(v.Vigencia)
	if err3 := GetRequestLegacy("UrlcrudAgora", url3, &actaInicio); err3 != nil {
		panic("Acta inicio -> " + err3.Error())
	} else if len(actaInicio) == 0 {
		panic("Acta inicio -> No se pudo encontrar el acta de inicio")
	}

	preliquidacion.FechaInicio = actaInicio[0].FechaInicio
	preliquidacion.FechaFin = actaInicio[0].FechaFin
	preliquidacion.Cdp = desagregado[0].Disponibilidad
	preliquidacion.NombreCompleto = docente[0].NomProveedor
	preliquidacion.PersonaId = docente[0].Id
	preliquidacion.NumeroSemanas = v.NumeroSemanas
	preliquidacion.Resolucion = v.ResolucionVinculacionDocenteId.Id

	if err2 := SendRequestNew("UrlmidTitan", "preliquidacion", "POST", &c, &preliquidacion); err2 != nil {
		panic("Preliquidando -> " + err2.Error())
	}

	return
}

// Al cancelar una vinculación se debe ajustar la liquidación del contrato en Titan
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

	contratoReliquidar := &models.ContratoCancelacion{
		NumeroContrato: *cancelado.NumeroContrato,
		Vigencia:       cancelado.Vigencia,
		//ValorContrato:  cancelado.ValorContrato - cancelacion.ValorContrato,
		FechaAnulacion: cancelacion.FechaInicio,
		Documento:      strconv.Itoa(int(cancelacion.PersonaId)),
	}

	// calcular el desagregado de la cancelación individual
	if cancelado.ResolucionVinculacionDocenteId.Dedicacion != "HCH" {
		cancelado.NumeroSemanas = cancelado.NumeroSemanas - cancelacion.NumeroSemanas
		dedicacion := cancelado.ResolucionVinculacionDocenteId.Dedicacion
		nivel := cancelado.ResolucionVinculacionDocenteId.NivelAcademico
		if desagregado, err = CalcularDesagregadoTitan(cancelado, dedicacion, nivel); err != nil {
			panic(err)
		}

		for concepto, valor := range desagregado {
			if concepto != "NumeroContrato" && concepto != "Vigencia" && concepto != "SueldoBasico" {
				valores[concepto] = valor.(float64)
			}
		}
		contratoReliquidar.Desagregado = &valores
	}

	fmt.Println("APLICAR ANULACIÓN ", contratoReliquidar)
	fmt.Println("APLICAR ANULACIÓN ", contratoReliquidar.Desagregado)
	if err2 := SendRequestNew("UrlmidTitan", "novedadVE/aplicar_anulacion", "POST", &c, &contratoReliquidar); err2 != nil {
		panic("Reliquidando -> " + err2.Error())
	}

	return
}

// Envía a Titan la información de lso contratos afectados en una reducción
func ReducirContratosTitan(reduccion *models.Reduccion, modificacion *models.VinculacionDocente, valorReduccion float64) (outputError map[string]interface{}) {
	defer func() {
		if err := recover(); err != nil {
			outputError = map[string]interface{}{"funcion": "ReducirContratosTitan", "err": err, "status": "500"}
			panic(outputError)
		}
	}()
	var c models.ContratoPreliquidacion

	JsonDebug(reduccion)
	if reduccion.ContratoNuevo != nil && reduccion.ContratoNuevo.ValorContratoReduccion == 0 {
		reduccion.ContratoNuevo.ValorContratoReduccion = valorReduccion
	}
	fmt.Println("APLICAR REDUCCIÓN ", reduccion)
	if reduccion.ContratoNuevo != nil {
		fmt.Println("APLICAR REDUCCIÓN ", reduccion.ContratoNuevo.DesagregadoReduccion)
	}
	if err2 := SendRequestNew("UrlmidTitan", "novedadVE/aplicar_reduccion", "POST", &c, &reduccion); err2 != nil {
		panic("Reliquidando -> " + err2.Error())
	}

	return nil
}
