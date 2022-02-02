package helpers

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/udistrital/resoluciones_mid_v2/models"
)

func EsDocentePlanta(docenteId string) (planta bool, outputError map[string]interface{}) {
	var docente models.ObjetoDocentePlanta
	var esPlanta bool
	url := "consultar_datos_docente/" + docenteId
	if err := GetRequestWSO2("NscrudAcademica", url, &docente); err != nil {
		outputError = map[string]interface{}{"funcion": "/EsDocentePlanta", "err": err.Error(), "status": "500"}
		return false, outputError
	}
	esPlanta = docente.DocenteCollection.Docente[0].Planta == "true"

	return esPlanta, nil
}

func BuscarCategoriaDocente(vigencia, periodo, docenteId string) (categoria models.ObjetoCategoriaDocente, outputError map[string]interface{}) {
	var respuesta models.ObjetoCategoriaDocente

	url := fmt.Sprintf("categoria_docente/%s/%s/%s", vigencia, periodo, docenteId)
	if err := GetRequestWSO2("NscrudUrano", url, &respuesta); err != nil {
		outputError = map[string]interface{}{"funcion": "/BuscarCategoriaDocente", "err": err.Error(), "status": "500"}
		return categoria, outputError
	}

	return respuesta, outputError
}

func ListarDocentesHorasLectivas(vigencia, periodo, dedicacion, facultad, nivelAcademico string) (cargaLectiva models.ObjetoCargaLectiva, outputError map[string]interface{}) {
	defer func() {
		if err := recover(); err != nil {
			outputError = map[string]interface{}{"funcion": "/ListarDocentesHorasLectivas", "err": err, "status": "500"}
			panic(outputError)
		}
	}()
	dedicacionOld := HomologarDedicacionNombre(dedicacion, "old")
	facultadOld, err := HomologarFacultad("new", facultad)
	if err != nil {
		panic(err)
	}
	var docentesCarga models.ObjetoCargaLectiva

	url := fmt.Sprintf("carga_lectiva/%s/%s/%d/%s/%s", vigencia, "1", dedicacionOld, facultadOld, nivelAcademico)
	if err2 := GetRequestWSO2("NscrudAcademica", url, &docentesCarga); err2 != nil {
		panic(err2.Error())
	}

	return docentesCarga, outputError
}

func ListarDocentesCargaHoraria(vigencia, periodo, dedicacion, facultad, nivelAcademico string) (cargaHoraria models.ObjetoCargaLectiva, outputError map[string]interface{}) {
	defer func() {
		if err := recover(); err != nil {
			outputError = map[string]interface{}{"funcion": "/ListarDocentesCargaHoraria", "err": err, "status": "500"}
			panic(outputError)
		}
	}()

	docentesCargaHoraria, err := ListarDocentesHorasLectivas(vigencia, periodo, dedicacion, facultad, nivelAcademico)
	if err != nil {
		panic(err)
	}

	for i := range docentesCargaHoraria.CargasLectivas.CargaLectiva {
		categoria, err2 := BuscarCategoriaDocente(vigencia, periodo, docentesCargaHoraria.CargasLectivas.CargaLectiva[i].DocDocente)
		if err2 != nil {
			panic(err2)
		}
		docentesCargaHoraria.CargasLectivas.CargaLectiva[i].CategoriaNombre = strings.Trim(categoria.CategoriaDocente.Categoria, " ")
		docentesCargaHoraria.CargasLectivas.CargaLectiva[i].IDCategoria = categoria.CategoriaDocente.IDCategoria
		docentesCargaHoraria.CargasLectivas.CargaLectiva[i].NombreTipoVinculacion = dedicacion
		docentesCargaHoraria.CargasLectivas.CargaLectiva[i].IDTipoVinculacion = fmt.Sprintf("%d", HomologarDedicacionNombre(dedicacion, "new"))
		if dedicacion == "TCO" {
			docentesCargaHoraria.CargasLectivas.CargaLectiva[i].HorasLectivas = "20"
		}
		if dedicacion == "MTO" {
			docentesCargaHoraria.CargasLectivas.CargaLectiva[i].HorasLectivas = "40"
		}
		docentesCargaHoraria.CargasLectivas.CargaLectiva[i].IDFacultad = facultad

		proyectoCurricularOld := docentesCargaHoraria.CargasLectivas.CargaLectiva[i].IDProyecto
		docentesCargaHoraria.CargasLectivas.CargaLectiva[i].DependenciaAcademica, _ = strconv.Atoi(proyectoCurricularOld)
		proyectoCurricularNew, err3 := HomologarProyectoCurricular(proyectoCurricularOld)
		if err3 != nil {
			panic(err3)
		}
		docentesCargaHoraria.CargasLectivas.CargaLectiva[i].IDProyecto = proyectoCurricularNew

	}

	return docentesCargaHoraria, outputError
}
