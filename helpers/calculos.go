package helpers

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/astaxie/beego/logs"
	. "github.com/udistrital/golog"
	"github.com/udistrital/resoluciones_mid_v2/models"
)

func CalcularComponentesSalario(d []models.ObjetoDesagregado) (d2 []map[string]interface{}, outputError map[string]interface{}) {
	defer func() {
		if err := recover(); err != nil {
			outputError = map[string]interface{}{"funcion": "CalcularComponentesSalario", "err": err, "status": "500"}
			panic(outputError)
		}
	}()

	vigencia := strconv.Itoa(d[0].Vigencia)

	puntoSalarial, err := CargarParametroPeriodo(vigencia, "PSAL")
	if err != nil {
		logs.Error(err)
		panic(err)
	}

	salarioMinimo, err2 := CargarParametroPeriodo(vigencia, "SMMLV")
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
	predicadosBase = predicadosBase + "valor_salario_minimo(" + fmt.Sprintf("%.f", salarioMinimo) + ", " + vigencia + ").\n"
	predicadosBase = predicadosBase + "sueldo_basico(N,D,C,V,S):-(N==pregrado->valor_punto(X,V);valor_salario_minimo(X,V)),factor(N,D,C,Y,V),(D==tco->T is Y*(X/160);D==mto->T is Y*(X/80);T is X*Y),S is T.\n"
	predicadosBase = predicadosBase + "subrubro_desagregado(N,D,C,V,CP,R):-sueldo_basico(N,D,C,V,S),porcentaje_devengo_v2(V,CP,X), T is S * X, R is (T rnd 0).\n"
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
			vigencia + ",CP,R).")

		for _, res := range prestaciones {
			nombre := fmt.Sprintf("%s", res.ByName_("CP"))
			valor, _ := strconv.ParseFloat(fmt.Sprintf("%s", res.ByName_("R")), 64)
			resultado[nombre] = valor
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

func CargarParametroPeriodo(vigencia, codigo string) (parametro float64, outputError map[string]interface{}) {
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
		return parametro, outputError
	}
	if err2 := json.Unmarshal([]byte(s[0].Valor), &valor); err2 != nil {
		outputError = map[string]interface{}{"funcion": "/CargarParametroPeriodo-parse", "err": err2.Error(), "status": "500"}
		return parametro, outputError
	}

	return valor["Valor"].(float64), outputError
}

func CalcularTotalSalarios(v []models.VinculacionDocente) (total float64) {
	var sumatoria float64
	for _, docente := range v {
		sumatoria = sumatoria + docente.ValorContrato
	}

	return sumatoria
}

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
