package helpers

import (
	"fmt"
	"strconv"

	"github.com/udistrital/resoluciones_mid_v2/models"
)

// Verifica si un docente hace parte de la planta docente de la universidad
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

// Consulta la categoría del docente para una vigenca y periodo académico específicos
func BuscarCategoriaDocente(vigencia, periodo, docenteId string) (categoria models.ObjetoCategoriaDocente, outputError map[string]interface{}) {
	var respuesta models.ObjetoCategoriaDocente

	url := fmt.Sprintf("categoria_docente/%s/%s/%s", vigencia, periodo, docenteId)
	if err := GetRequestWSO2("NscrudUrano", url, &respuesta); err != nil {
		outputError = map[string]interface{}{"funcion": "/BuscarCategoriaDocente", "err": err.Error(), "status": "500"}
		return categoria, outputError
	}

	return respuesta, outputError
}

// Consulta el nombre completo y tipo de documento de un docente
func BuscarDatosPersonalesDocente(personaId float64) (p models.InformacionPersonaNatural, outputError map[string]interface{}) {
	defer func() {
		if err := recover(); err != nil {
			outputError = map[string]interface{}{"funcion": "/BuscarDatosPersonalesDocente", "err": err, "status": "500"}
			panic(outputError)
		}
	}()
	var personas []models.InformacionPersonaNatural

	// TODO: Esta consulta deberá moverse a terceros, tabla datos_identificacion
	// cuando core_amazon_crud se encuentre deprecada definitivamente y los datos de
	// expedición de cedulas se encuentren relacionados con ubicaciones_crud
	// url := "datos_identificacion?query=Numero:" + strconv.Itoa(int(personaId))
	url := "informacion_persona_natural?query=Id:" + strconv.Itoa(int(personaId))
	if err := GetRequestLegacy("UrlcrudAgora", url, &personas); err != nil {
		panic(err.Error())
	}
	persona := personas[0]
	persona.NomProveedor = fmt.Sprintf("%s %s %s %s", persona.PrimerNombre, persona.SegundoNombre, persona.PrimerApellido, persona.SegundoApellido)
	persona.CiudadExpedicionDocumento = strconv.Itoa(int(persona.IdCiudadExpedicionDocumento))

	return persona, outputError
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

	url := fmt.Sprintf("carga_lectiva/%s/%s/%d/%s/%s", vigencia, periodo, dedicacionOld, facultadOld, nivelAcademico)
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
		docentesCargaHoraria.CargasLectivas.CargaLectiva[i].CategoriaNombre = "ASOCIADO" // strings.Trim(categoria.CategoriaDocente.Categoria, " ")
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
