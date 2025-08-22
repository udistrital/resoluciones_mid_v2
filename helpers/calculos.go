package helpers

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/astaxie/beego/logs"
	. "github.com/udistrital/golog"
	"github.com/udistrital/resoluciones_mid_v2/models"
)

// Calcula el valor del contrato para cada docente utilizando el conjunto de reglas CDVE
func CalcularSalarioPrecontratacion(docentesVincular []models.VinculacionDocente) (docentesInsertar []models.VinculacionDocente, outputError map[string]interface{}) {
	defer func() {
		if err := recover(); err != nil {
			outputError = map[string]interface{}{"funcion": "/CalcularSalarioPrecontratacion", "err": err, "status": "500"}
			panic(outputError)
		}
	}()

	nivelAcademico := docentesVincular[0].ResolucionVinculacionDocenteId.NivelAcademico
	vigencia := strconv.Itoa(int(docentesVincular[0].Vigencia))
	var a string
	var categoria string

	_, salarioMinimo, err1 := CargarParametroPeriodo(vigencia, "SMMLV")
	if err1 != nil {
		logs.Error(err1)
		panic(err1)
	}

	for x, docente := range docentesVincular {
		p, err2 := EsDocentePlanta(strconv.Itoa(int(docente.PersonaId)))
		if err2 != nil {
			logs.Error(err2)
			panic(err2)
		}
		if p && strings.ToLower(nivelAcademico) == "posgrado" {
			categoria = strings.TrimSpace(docente.Categoria) + "ud"
		} else {
			categoria = strings.TrimSpace(docente.Categoria)
		}

		var predicados string
		if strings.ToLower(nivelAcademico) == "posgrado" {
			predicados = "valor_salario_minimo(" + strconv.Itoa(int(salarioMinimo)) + "," + vigencia + ")." + "\n"
			docente.NumeroSemanas = 1
		} else if strings.ToLower(nivelAcademico) == "pregrado" {
			predicados = "valor_punto(" + strconv.Itoa(int(docente.ValorPuntoSalarial)) + "," + vigencia + ")." + "\n"
		}

		predicados = predicados + "categoria(" + strconv.FormatInt(int64(docente.PersonaId), 10) + "," + strings.ToLower(categoria) + "," + vigencia + ")." + "\n"
		predicados = predicados + "vinculacion(" + strconv.FormatInt(int64(docente.PersonaId), 10) + "," + strings.ToLower(docente.ResolucionVinculacionDocenteId.Dedicacion) + "," + vigencia + ")." + "\n"
		predicados = predicados + "horas(" + strconv.FormatInt(int64(docente.PersonaId), 10) + "," + strconv.Itoa(docente.NumeroHorasSemanales*docente.NumeroSemanas) + "," + vigencia + ")." + "\n"
		reglasbase, err4 := CargarReglasBase("CDVE")
		if err4 != nil {
			logs.Error(err4)
			panic(err4)
		}

		reglasbase = reglasbase + predicados
		m := NewMachine().Consult(reglasbase)
		contratos := m.ProveAll("valor_contrato(" + strings.ToLower(nivelAcademico) + "," + strconv.FormatInt(int64(docente.PersonaId), 10) + "," + vigencia + ",X).")

		for _, solution := range contratos {
			a = fmt.Sprintf("%s", solution.ByName_("X"))
		}
		salario, err5 := strconv.ParseFloat(a, 64)
		if err5 != nil {
			logs.Error(err5)
			panic(err5.Error())
		}
		docentesVincular[x].ValorContrato = salario

	}
	return docentesVincular, nil
}

// Calcula el valor de la modificación del contrato de un docente
func CalcularValorContratoReduccion(v [1]models.VinculacionDocente, semanasRestantes int, horasOriginales int, nivelAcademico string) (salarioTotal float64, outputError map[string]interface{}) {
	defer func() {
		if err := recover(); err != nil {
			outputError = map[string]interface{}{"funcion": "/CalcularValorContratoReduccion", "err": err, "status": "500"}
			panic(outputError)
		}
	}()
	var d []models.VinculacionDocente
	var salarioSemanasReducidas float64
	var salarioSemanasRestantes float64

	d = append(d, v[0])

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
			panic(err)
		}
		salarioSemanasRestantes = docentes[0].ValorContrato
	}

	salarioTotal = salarioSemanasReducidas + salarioSemanasRestantes
	return salarioTotal, outputError
}

// Calcula el desagregado general unitario para los parámetros recibidos
func CalcularComponentesSalario(d []models.ObjetoDesagregado) (d2 []map[string]interface{}, outputError map[string]interface{}) {
	defer func() {
		if err := recover(); err != nil {
			outputError = map[string]interface{}{"funcion": "CalcularComponentesSalario", "err": err, "status": "500"}
			panic(outputError)
		}
	}()

	vigencia := strconv.Itoa(d[0].Vigencia)

	predicadosPrestaciones, err := obtenerReglasDesagregado()
	if err != nil {
		logs.Error(err)
		panic(err)
	}

	_, puntoSalarial, err := CargarParametroPeriodo(vigencia, "PSAL")
	if err != nil {
		logs.Error(err)
		panic(err)
	}

	_, salarioMinimo, err2 := CargarParametroPeriodo(vigencia, "SMMLV")
	if err2 != nil {
		logs.Error(err2)
		panic(err2)
	}

	reglas1, err3 := CargarReglasBase("CDVE")
	if err3 != nil {
		logs.Error(err3)
		panic(err3)
	}

	reglas2, err4 := CargarReglasBase("HCS")
	if err4 != nil {
		logs.Error(err4)
		panic(err4)
	}

	reglas3, err5 := CargarReglasBase("SeguridadSocial")
	if err5 != nil {
		logs.Error(err5)
		panic(err5)
	}

	predicadosBase := "valor_punto(" + fmt.Sprintf("%.f", puntoSalarial) + ", " + vigencia + ").\n"
	predicadosBase = predicadosBase + predicadosPrestaciones
	predicadosBase = predicadosBase + "valor_salario_minimo(" + fmt.Sprintf("%.f", salarioMinimo) + ", " + vigencia + ").\n"
	predicadosBase = predicadosBase + "sueldo_basico(N,D,C,V,S):-(N==pregrado->valor_punto(X,V);valor_salario_minimo(X,V)),factor(N,D,C,Y,V),(D==tco->T is Y*(X/160);D==mto->T is Y*(X/80);T is X*Y),S is T.\n"
	predicadosBase = predicadosBase + "porcentaje_concepto_v2(P, CO, T, UT, R) :- (UT = semanas -> (T @>= 24.0 -> porcentaje_devengo_v2_mayor(P, CO, V) ; porcentaje_devengo_v2_menor(P, CO, V)) ; UT = meses -> (T @>= 6.0 -> porcentaje_devengo_v2_mayor(P, CO, V) ; porcentaje_devengo_v2_menor(P, CO, V))), R is V.\n"
	predicadosBase = predicadosBase + "subrubro_desagregado(N,D,C,V,T,CP,R):-sueldo_basico(N,D,C,V,S),porcentaje_concepto_v2(V, CP, T, semanas, X), CO is S * X, R is (CO rnd 0).\n"
	predicadosBase = predicadosBase + "subrubro_desagregado2(N,D,C,V,CP,R):-sueldo_basico(N,D,C,V,S),concepto_aporte(CP,X,planta,2388),T is S * X, R is (T rnd 0).\n"
	predicadosBase = predicadosBase + "subrubro_salud(N,D,C,V,R):-sueldo_basico(N,D,C,V,S),salud(V,X),T is S * (X/100), R is (T rnd 0).\n"

	reglas := reglas1 + reglas2 + reglas3 + predicadosBase

	m := NewMachine().Consult(reglas)

	resultados := make([]map[string]interface{}, len(d))

	for i, obj := range d {

		resultado := map[string]interface{}{}
		resultados[i] = make(map[string]interface{})

		resultado["Vigencia"] = (d[i].Vigencia)
		resultado["Categoria"] = (d[i].Categoria)
		resultado["Dedicacion"] = (d[i].Dedicacion)
		resultado["NivelAcademico"] = (d[i].NivelAcademico)
		resultado["EsDePlanta"] = (d[i].EsDePlanta)

		if strings.ToLower(obj.NivelAcademico) == "posgrado" && obj.EsDePlanta {
			obj.Categoria = obj.Categoria + "UD"
		}

		salarios := m.ProveAll("sueldo_basico(" +
			strings.ToLower(obj.NivelAcademico) + "," +
			strings.ToLower(obj.Dedicacion) + "," +
			strings.ToLower(obj.Categoria) + "," +
			vigencia + ",S).")
		for _, salario := range salarios {
			resultado["salarioBasico"], _ = strconv.ParseFloat(fmt.Sprintf("%s", salario.ByName_("S")), 64)
		}

		prestaciones := m.ProveAll("subrubro_desagregado(" +
			strings.ToLower(obj.NivelAcademico) + "," +
			strings.ToLower(obj.Dedicacion) + "," +
			strings.ToLower(obj.Categoria) + "," +
			vigencia + "," +
			strconv.Itoa(int(obj.Semanas)) +
			",CP,R).")

		for _, res := range prestaciones {
			nombre := fmt.Sprintf("%s", res.ByName_("CP"))
			parts := strings.Split(nombre, "_")
			for i := 1; i < len(parts); i++ {
				if len(parts[i]) > 0 {
					parts[i] = strings.Title(parts[i])
				}
			}
			n := strings.Join(parts, "")
			valor, _ := strconv.ParseFloat(fmt.Sprintf("%s", res.ByName_("R")), 64)
			resultado[n] = valor
		}

		aportes := m.ProveAll("subrubro_desagregado2(" +
			strings.ToLower(obj.NivelAcademico) + "," +
			strings.ToLower(obj.Dedicacion) + "," +
			strings.ToLower(obj.Categoria) + "," +
			vigencia + ",CP,R).")

		for _, res := range aportes {
			nombre := fmt.Sprintf("%s", res.ByName_("CP"))
			valor, _ := strconv.ParseFloat(fmt.Sprintf("%s", res.ByName_("R")), 64)
			if valor != 0 {
				resultado[nombre] = valor
			}
		}

		salud := m.ProveAll("subrubro_salud(" +
			strings.ToLower(obj.NivelAcademico) + "," +
			strings.ToLower(obj.Dedicacion) + "," +
			strings.ToLower(obj.Categoria) + "," +
			vigencia + ",R).")

		for _, res := range salud {
			valor, _ := strconv.ParseFloat(fmt.Sprintf("%s", res.ByName_("R")), 64)
			if valor != 0 {
				resultado["salud"] = valor
			}
		}
		resultados[i] = resultado
	}
	return resultados, outputError
}

// Carga el parámetro para el periodo/vigencia indicado y extrae el valor correspondiente
func CargarParametroPeriodo(vigencia, codigo string) (id int, parametro float64, outputError map[string]interface{}) {
	var s []models.ParametroPeriodo
	var valor map[string]interface{}
	var url string
	if codigo == "PSAL" {
		url = "parametro_periodo?order=desc&sortby=Id&query=ParametroId.CodigoAbreviacion:" + codigo
	} else {
		url = "parametro_periodo?query=ParametroId.CodigoAbreviacion:" + codigo + ",PeriodoId.Year:" + vigencia
	}
	if err := GetRequestNew("UrlcrudParametros", url, &s); err != nil {
		outputError = map[string]interface{}{"funcion": "/CargarParametroPeriodo", "err": err.Error(), "status": "500"}
		return id, parametro, outputError
	} else if len(s) == 0 {
		outputError = map[string]interface{}{"funcion": "/CargarParametroPeriodo", "err": "No se pudo cargar el parámetro " + codigo, "status": "500"}
		return id, parametro, outputError
	}
	if err2 := json.Unmarshal([]byte(s[0].Valor), &valor); err2 != nil {
		outputError = map[string]interface{}{"funcion": "/CargarParametroPeriodo-parse", "err": err2.Error(), "status": "500"}
		return id, parametro, outputError
	}

	return s[0].Id, valor["Valor"].(float64), outputError
}

// Calcula la sumatoria del valor de los contratos de una resolución
func CalcularTotalSalarios(v []models.VinculacionDocente) (total float64) {
	var sumatoria float64
	for _, docente := range v {
		sumatoria = sumatoria + docente.ValorContrato
	}

	return sumatoria
}

// Carga el conjunto de reglas del API Ruler del dominio indicado
func CargarReglasBase(dominio string) (reglas string, outputError map[string]interface{}) {
	var reglasbase string = ``
	var v []models.Predicado

	url := "predicado?query=Dominio.Nombre:" + dominio + "&limit=-1"
	if err := GetRequestLegacy("Urlruler", url, &v); err != nil {
		outputError = map[string]interface{}{"funcion": "/CargarReglasBase", "err": err.Error(), "status": "500"}
		return reglasbase, outputError
	}
	reglasbase = reglasbase + FormatoReglas(v)

	return reglasbase, nil
}

// Compila un conjunto de reglas en forma de cadena de texto para su uso con el motor de reglas
func FormatoReglas(v []models.Predicado) (reglas string) {
	var arregloReglas = make([]string, len(v))
	reglas = ""

	for i := 0; i < len(v); i++ {
		arregloReglas[i] = v[i].Nombre
	}

	for i := 0; i < len(arregloReglas); i++ {
		reglas = reglas + arregloReglas[i] + "\n"
	}
	return
}

// Genera las reglas de porcentajes de desagregado
func obtenerReglasDesagregado() (predicadosString string, outputError map[string]interface{}) {
	// se debe obtener desde parametros los valores de porcentaje de prestaciones y cargar los predicados dinamicos
	// obtener el periodo vigente para app de resoluciones
	// contemplar agregar el aplicacion_id para crear periodos exclusivos para resoluciones
	var predicados []models.Predicado
	var periodo []models.Periodo
	anoActual := time.Now().Year()
	url1 := "/periodo?limit=-1&query=year:" + strconv.Itoa(anoActual) + ",codigo_abreviacion:PA,activo:true"
	if err1 := GetRequestNew("UrlcrudParametros", url1, &periodo); err1 == nil {
		var parametro []models.Parametro
		url2 := "/parametro?limit=-1&query=codigo_abreviacion:PDVE,activo:true"
		if err2 := GetRequestNew("UrlcrudParametros", url2, &parametro); err2 == nil {
			var parametroPeriodo []models.ParametroPeriodo
			url3 := "/parametro_periodo?limit=-1&query=parametro_id:" + strconv.Itoa(parametro[0].Id) + ",periodo_id:" + strconv.Itoa(periodo[0].Id) + ",activo:true"
			if err3 := GetRequestNew("UrlcrudParametros", url3, &parametroPeriodo); err3 == nil {
				for _, pp := range parametroPeriodo {
					var valores map[string]map[string]float64
					json.Unmarshal([]byte(pp.Valor), &valores)
					for concepto, porcentajes := range valores {
						if mayor, ok := porcentajes["porcentaje_mayor"]; ok {
							predicados = append(predicados, models.Predicado{Nombre: "porcentaje_devengo_v2_mayor(" + strconv.Itoa(anoActual) + "," + strings.ToLower(concepto) + "," + fmt.Sprintf("%.5f", mayor) + ")."})
						}
						if menor, ok := porcentajes["porcentaje_menor"]; ok {
							predicados = append(predicados, models.Predicado{Nombre: "porcentaje_devengo_v2_menor(" + strconv.Itoa(anoActual) + "," + strings.ToLower(concepto) + "," + fmt.Sprintf("%.5f", menor) + ")."})
						}
					}
				}
			} else {
				outputError = map[string]interface{}{"funcion": "/CargarParametroPeriodo-reglasDesagregado", "err": err3.Error(), "status": "500"}
				return predicadosString, outputError
			}
		} else {
			outputError = map[string]interface{}{"funcion": "/CargarParametroPeriodo-reglasDesagregado", "err": err2.Error(), "status": "500"}
			return predicadosString, outputError
		}
	} else {
		outputError = map[string]interface{}{"funcion": "/CargarParametroPeriodo-reglasDesagregado", "err": err1.Error(), "status": "500"}
		return predicadosString, outputError
	}
	predicadosString = FormatoReglas(predicados)
	return predicadosString, outputError
}
