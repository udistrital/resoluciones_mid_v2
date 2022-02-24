package helpers

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	. "github.com/udistrital/golog"
	"github.com/udistrital/resoluciones_mid_v2/models"
)

func CalculoSalarios(v []models.VinculacionDocente, periodo int) (total int, outputError map[string]interface{}) {

	defer func() {
		if err := recover(); err != nil {
			outputError = map[string]interface{}{"funcion": "/CalculoSalarios", "err": err, "status": "502"}
			panic(outputError)
		}
	}()
	var totalesDisponibilidad int
	if v, err1 := CalcularSalarioPrecontratacion(v); err1 == nil {
		totalesSalario := CalcularTotalSalario(v)
		vigencia := strconv.Itoa(int(v[0].Vigencia))
		periodo := strconv.Itoa(periodo)
		// disponibilidad := strconv.Itoa(v[0].Disponibilidad)

		url := "/vinculacion_docente/get_valores_totales_x_disponibilidad/" + vigencia + "/" + periodo + "/" // + disponibilidad
		if err2 := GetRequestNew("UrlCrudResoluciones", url, &totalesDisponibilidad); err2 == nil {
			total = int(totalesSalario) + totalesDisponibilidad
		} else {
			logs.Error(err2)
			outputError = map[string]interface{}{"funcion": "/CalculoSalarios2", "err2": err2.Error(), "status": "502"}
			return total, outputError
		}
	} else {
		logs.Error(err1)
		return total, err1
	}
	return
}

func CalcularSalarioPrecontratacion(docentes_a_vincular []models.VinculacionDocente) (docentes_a_insertar []models.VinculacionDocente, outputError map[string]interface{}) {

	defer func() {
		if err := recover(); err != nil {
			outputError = map[string]interface{}{"funcion": "/CalcularSalarioPrecontratacion", "err": err, "status": "502"}
			panic(outputError)
		}
	}()
	nivelAcademico := docentes_a_vincular[0].ResolucionVinculacionDocenteId.NivelAcademico
	vigencia := strconv.Itoa(int(docentes_a_vincular[0].Vigencia))
	var a string
	var categoria string

	salarioMinimo, err1 := CargarSalarioMinimo(vigencia)
	if err1 != nil {
		logs.Error(err1)
		//outputError = map[string]interface{}{"funcion": "/CalcularSalarioPrecontratacion1", "err1": err1, "status": "502"}
		return docentes_a_insertar, err1
	}

	for x, docente := range docentes_a_vincular {
		p, err2 := EsDocentePlanta(strconv.FormatInt(int64(docente.PersonaId), 10))
		if err1 != nil {
			logs.Error(err2)
			//outputError = map[string]interface{}{"funcion": "/CalcularSalarioPrecontratacion2", "err2": err2, "status": "502"}
			return docentes_a_insertar, err2
		}
		if p && strings.ToLower(nivelAcademico) == "posgrado" {
			categoria = strings.TrimSpace(docente.Categoria) + "ud"
		} else {
			categoria = strings.TrimSpace(docente.Categoria)
		}

		var predicados string
		if strings.ToLower(nivelAcademico) == "posgrado" {
			predicados = "valor_salario_minimo(" + strconv.Itoa(salarioMinimo.Valor) + "," + vigencia + ")." + "\n"
			docente.NumeroSemanas = 1
		} else if strings.ToLower(nivelAcademico) == "pregrado" {
			a, err3 := CargarPuntoSalarial()
			if err3 != nil {
				logs.Error(err2)
				//outputError = map[string]interface{}{"funcion": "/CalcularSalarioPrecontratacion3", "err3": err3, "status": "502"}
				return docentes_a_insertar, err3
			}
			predicados = "valor_punto(" + strconv.Itoa(a.ValorPunto) + ", " + vigencia + ")." + "\n"
		}

		predicados = predicados + "categoria(" + strconv.FormatInt(int64(docente.PersonaId), 10) + "," + strings.ToLower(categoria) + ", " + vigencia + ")." + "\n"
		predicados = predicados + "vinculacion(" + strconv.FormatInt(int64(docente.PersonaId), 10) + "," + strings.ToLower(docente.ResolucionVinculacionDocenteId.Dedicacion) + ", " + vigencia + ")." + "\n"
		predicados = predicados + "horas(" + strconv.FormatInt(int64(docente.PersonaId), 10) + "," + strconv.Itoa(docente.NumeroHorasSemanales*docente.NumeroSemanas) + ", " + vigencia + ")." + "\n"
		reglasbase, err4 := CargarReglasBase("CDVE")
		beego.Info("predicados: ", predicados, "a: ", a)
		if err4 != nil {
			logs.Error(err4)
			//outputError = map[string]interface{}{"funcion": "/CalcularSalarioPrecontratacion4", "err4": err4, "status": "502"}
			return docentes_a_insertar, err4
		}
		reglasbase = reglasbase + predicados
		m := NewMachine().Consult(reglasbase)
		beego.Info("m: ", m)
		contratos := m.ProveAll("valor_contrato(" + strings.ToLower(nivelAcademico) + "," + strconv.FormatInt(int64(docente.PersonaId), 10) + "," + vigencia + ",X).")
		for _, solution := range contratos {
			a = fmt.Sprintf("%s", solution.ByName_("X"))
			//beego.Info("a: ", a)
		}
		f, err5 := strconv.ParseFloat(a, 64)
		if err5 != nil {
			logs.Error(err5)
			outputError = map[string]interface{}{"funcion": "/CalcularSalarioPrecontratacion5", "err5": err5.Error(), "status": "502"}
			return docentes_a_vincular, outputError
		}
		salario := f
		beego.Info("f: ", f, "salario: ", salario)
		docentes_a_vincular[x].ValorContrato = salario

	}
	return docentes_a_vincular, nil
}

func CargarSalarioMinimo(vigencia string) (p models.SalarioMinimo, outputError map[string]interface{}) {

	var v []models.SalarioMinimo
	url := "/salario_minimo?limit=1&query=Vigencia:" + vigencia
	if err := GetRequestLegacy("UrlcrudCore", url, &v); err != nil {
		outputError = map[string]interface{}{"funcion": "/CargarSalarioMinimo", "err": err.Error(), "status": "404"}
		return v[0], outputError
	}
	return v[0], nil
}

func EsDocentePlanta(idPersona string) (docentePlanta bool, outputError map[string]interface{}) {

	var temp map[string]interface{}
	var esDePlanta bool
	url := "consultar_datos_docente/" + idPersona
	if err := GetRequestLegacy("UrlcrudWSO2", url, &temp); err != nil {
		outputError = map[string]interface{}{"funcion": "/EsDocentePlanta1", "err": err.Error(), "status": "404"}
		return false, outputError
	}
	jsonDocentes, err1 := json.Marshal(temp)
	if err1 != nil {
		outputError = map[string]interface{}{"funcion": "/EsDocentePlanta2", "err": "Error en codificación de datos", "status": "404"}
		return false, outputError
	}
	var tempDocentes models.ObjetoDocentePlanta
	err2 := json.Unmarshal(jsonDocentes, &tempDocentes)
	if err2 != nil {
		outputError = map[string]interface{}{"funcion": "/EsDocentePlanta3", "err": "Error en decodificación de datos", "status": "404"}
		return false, outputError
	}
	if tempDocentes.DocenteCollection.Docente[0].Planta == "true" {
		esDePlanta = true
	} else {
		esDePlanta = false
	}

	return esDePlanta, nil
}

func CargarPuntoSalarial() (p models.PuntoSalarial, outputError map[string]interface{}) {
	var v []models.PuntoSalarial
	url := "/punto_salarial?sortby=Vigencia&order=desc&limit=1"
	if err := GetRequestLegacy("UrlcrudCore", url, &v); err != nil {
		outputError = map[string]interface{}{"funcion": "/CargarPuntoSalarial", "err": err.Error(), "status": "404"}
		return v[0], outputError
	}
	return v[0], nil
}

func CalcularTotalSalario(v []models.VinculacionDocente) (total float64) {
	var sumatoria float64
	for _, docente := range v {
		sumatoria = sumatoria + docente.ValorContrato
	}

	return sumatoria
}

func CalcularValorContratoReduccion(v [1]models.VinculacionDocente, semanasRestantes int, horasOriginales int, nivelAcademico string, periodo int) (salarioTotal float64, outputError map[string]interface{}) {
	var d []models.VinculacionDocente
	var salarioSemanasReducidas float64
	var salarioSemanasRestantes float64

	jsonEjemplo, err1 := json.Marshal(v)
	if err1 != nil {
		outputError = map[string]interface{}{"funcion": "/CalcularValorContratoReduccion1", "err": err1.Error(), "status": "404"}
		return salarioTotal, outputError
	}
	err2 := json.Unmarshal(jsonEjemplo, &d)
	if err2 != nil {
		outputError = map[string]interface{}{"funcion": "/CalcularValorContratoReduccion2", "err": err2.Error(), "status": "404"}
		return salarioTotal, outputError
	}

	docentes, err := CalcularSalarioPrecontratacion(d)
	if err != nil {
		return salarioTotal, err
	}
	salarioSemanasReducidas = docentes[0].ValorContrato
	//Para posgrados no se deben tener en cuenta las semanas restantes
	if semanasRestantes > 0 && nivelAcademico == "PREGRADO" {
		d[0].NumeroSemanas = semanasRestantes
		d[0].NumeroHorasSemanales = horasOriginales
		docentes, err := CalcularSalarioPrecontratacion(d)
		if err != nil {
			return salarioTotal, err
		}
		salarioSemanasRestantes = docentes[0].ValorContrato
	}
	beego.Info("reducidas ", salarioSemanasReducidas, "restantes ", salarioSemanasRestantes)
	salarioTotal = salarioSemanasReducidas + salarioSemanasRestantes
	return salarioTotal, nil
}
