package helpers

import (
	"encoding/base64"
	"encoding/json"
	"strconv"
	"time"

	"github.com/astaxie/beego/logs"
	"github.com/udistrital/resoluciones_mid_v2/models"
)

func InsertarResolucion(nuevaRes models.ContenidoResolucion) (resolucionId int, outputError map[string]interface{}) {
	defer func() {
		if err := recover(); err != nil {
			outputError = map[string]interface{}{"funcion": "InsertarResolucion", "err": err, "status": "500"}
			panic(outputError)
		}
	}()

	var plantillas []models.Parametro
	var tipoRes models.Parametro
	var plantilla models.ContenidoResolucion

	url := "parametro?limit=0&query=Activo:true,ParametroPadreId.Id:" + strconv.Itoa(nuevaRes.Resolucion.TipoResolucionId)
	if err := GetRequestNew("UrlcrudParametros", url, &plantillas); err != nil {
		logs.Error(err)
		panic(err.Error())
	}
	url2 := "parametro/" + strconv.Itoa(nuevaRes.Resolucion.TipoResolucionId)
	if err := GetRequestNew("UrlcrudParametros", url2, &tipoRes); err != nil {
		logs.Error(err)
		panic(err.Error())
	}

	if tipoRes.CodigoAbreviacion != "RVIN" {
		var anteriorResvin models.ResolucionVinculacionDocente
		url := "resolucion_vinculacion_docente/" + strconv.Itoa(nuevaRes.ResolucionAnteriorId)
		if er := GetRequestNew("UrlcrudResoluciones", url, &anteriorResvin); er != nil {
			logs.Error(er)
			panic(er.Error())
		}
		nuevaRes.Vinculacion.Dedicacion = anteriorResvin.Dedicacion
		nuevaRes.Vinculacion.NivelAcademico = anteriorResvin.NivelAcademico
		nuevaRes.Vinculacion.FacultadId = anteriorResvin.FacultadId
	}

	if existe, id := validarExistenciaPlantilla(nuevaRes, plantillas); existe {
		var p *[]models.Resolucion
		url := "resolucion?query=Activo:true,TipoResolucionId:" + strconv.Itoa(id)
		if err2 := GetRequestNew("UrlcrudResoluciones", url, &p); err2 == nil && p != nil && len(*p) > 0 {
			var err3 map[string]interface{}
			if plantilla, err3 = CargarPlantilla((*p)[0].Id); err3 != nil {
				logs.Error(err3)
				panic(err3)
			}
		} else if err2 != nil {
			panic(err2.Error())
		}
	} else {
		return 0, nil
	}

	nuevaRes.Resolucion.Titulo = plantilla.Resolucion.Titulo
	nuevaRes.Resolucion.PreambuloResolucion = plantilla.Resolucion.PreambuloResolucion
	nuevaRes.Resolucion.ConsideracionResolucion = plantilla.Resolucion.ConsideracionResolucion
	nuevaRes.Resolucion.CuadroResponsabilidades = plantilla.Resolucion.CuadroResponsabilidades
	for _, art := range plantilla.Articulos {
		articulo := models.ComponenteResolucion{
			Texto: art.Articulo.Texto,
		}
		var paragrafos []models.ComponenteResolucion
		for _, par := range art.Paragrafos {
			paragrafo := models.ComponenteResolucion{
				Texto: par.Texto,
			}
			paragrafos = append(paragrafos, paragrafo)
		}
		nuevaRes.Articulos = append(nuevaRes.Articulos, models.Articulo{
			Articulo:   articulo,
			Paragrafos: paragrafos,
		})
	}

	var err4 map[string]interface{}
	if resolucionId, err4 = InsertarResolucionCompleta(nuevaRes); err4 != nil {
		logs.Error(err4)
		panic(err4)
	}

	var decData map[string]interface{}
	if data, err6 := base64.StdEncoding.DecodeString(nuevaRes.Usuario); err6 != nil {
		panic(err6.Error())
	} else {
		if err7 := json.Unmarshal(data, &decData); err7 != nil {
			panic(err7)
		}
	}
	usuario := decData["user"].(map[string]interface{})["sub"].(string)

	if err5 := CambiarEstadoResolucion(resolucionId, "RSOL", usuario); err5 != nil {
		logs.Error(err5)
		panic(err5.Error())
	}

	if tipoRes.CodigoAbreviacion != "RVIN" {
		var modResp models.ModificacionResolucion
		modRes := models.ModificacionResolucion{
			ResolucionNuevaId:    &models.Resolucion{Id: resolucionId},
			ResolucionAnteriorId: &models.Resolucion{Id: nuevaRes.ResolucionAnteriorId},
			Activo:               true,
		}
		if err6 := SendRequestNew("UrlcrudResoluciones", "modificacion_resolucion", "POST", &modResp, &modRes); err6 != nil {
			logs.Error(err6)
			panic(err6.Error())
		}
	}

	return resolucionId, outputError
}

func InsertarResolucionCompleta(v models.ContenidoResolucion) (resolucionId int, outputError map[string]interface{}) {
	defer func() {
		if err := recover(); err != nil {
			outputError = map[string]interface{}{"funcion": "InsertarResolucionCompleta", "err": err, "status": "500"}
			panic(outputError)
		}
	}()

	var resp models.Resolucion
	v.Resolucion.Activo = true
	v.Resolucion.DependenciaId = v.Vinculacion.FacultadId
	v.Resolucion.Vigencia, _, _ = time.Now().Date()

	if err := SendRequestNew("UrlCrudResoluciones", "resolucion", "POST", &resp, &v.Resolucion); err != nil {
		logs.Error(err)
		panic("resolucion -> " + err.Error())
	}
	v.Vinculacion.Id = resp.Id
	v.Vinculacion.Activo = true
	resolucionId = resp.Id

	var resVin models.ResolucionVinculacionDocente
	if err2 := SendRequestNew("UrlCrudResoluciones", "resolucion_vinculacion_docente", "POST", &resVin, &v.Vinculacion); err2 != nil {
		logs.Error(err2)
		panic("resolucion_vinculacion -> " + err2.Error())
	}
	if err3 := InsertarArticulos(v.Articulos, resolucionId); err3 != nil {
		logs.Error(err3)
		panic("Insertar articulos -> " + err3.Error())
	}

	return resolucionId, outputError
}

func InsertarArticulos(articulos []models.Articulo, resolucionId int) (err error) {
	var art models.ComponenteResolucion
	var par models.ComponenteResolucion
	for i, obj := range articulos {
		articulo := &models.ComponenteResolucion{
			Numero:         i + 1,
			ResolucionId:   &models.Resolucion{Id: resolucionId},
			Texto:          obj.Articulo.Texto,
			TipoComponente: "Articulo",
			Activo:         true,
		}
		if err = SendRequestNew("UrlCrudResoluciones", "componente_resolucion", "POST", &art, &articulo); err == nil {
			for j, obj2 := range obj.Paragrafos {
				paragrafo := &models.ComponenteResolucion{
					Numero:                    j + 1,
					ResolucionId:              &models.Resolucion{Id: resolucionId},
					Texto:                     obj2.Texto,
					TipoComponente:            "Paragrafo",
					Activo:                    true,
					ComponenteResolucionPadre: &models.ComponenteResolucion{Id: art.Id},
				}
				if err = SendRequestNew("UrlCrudResoluciones", "componente_resolucion", "POST", &par, &paragrafo); err != nil {
					return err
				}
			}
		} else {
			return err
		}
	}
	return err
}

func ListarResoluciones(limit, offset int) (listaRes []models.Resoluciones, total int, outputError map[string]interface{}) {
	defer func() {
		if err := recover(); err != nil {
			outputError = map[string]interface{}{"funcion": "ListarResoluciones", "err": err, "status": "500"}
			panic(outputError)
		}
	}()

	var res []models.Resolucion
	var resv []models.ResolucionVinculacionDocente
	var rest []models.ResolucionEstado
	var dep models.Dependencia
	var estado models.Parametro
	var tipo models.Parametro
	var err error

	url0 := "resolucion_vinculacion_docente?limit=0&fields=Id"
	if err = GetRequestNew("UrlcrudResoluciones", url0, &res); err != nil {
		logs.Error(err)
		panic(err.Error())
	}
	total = len(res)

	url := "resolucion?query=Activo:true&order=desc&sortby=Id&limit=" + strconv.Itoa(limit) + "&offset=" + strconv.Itoa(10*(offset-1))
	if err = GetRequestNew("UrlcrudResoluciones", url, &res); err != nil {
		logs.Error(err)
		panic(err.Error())
	}

	url2 := "resolucion_vinculacion_docente?query=Activo:true&order=desc&sortby=Id&limit=" + strconv.Itoa(limit) + "&offset=" + strconv.Itoa(10*(offset-1))
	if err = GetRequestNew("UrlcrudResoluciones", url2, &resv); err != nil {
		logs.Error(err)
		panic(err.Error())
	}

	for i := range res {
		url3 := "resolucion_estado?order=desc&sortby=Id&query=Activo:true,ResolucionId.Id:" + strconv.Itoa(res[i].Id)
		if err = GetRequestNew("UrlcrudResoluciones", url3, &rest); err != nil {
			panic(err.Error())
		}

		if len(rest) > 0 {
			if err = GetRequestNew("UrlcrudParametros", "parametro/"+strconv.Itoa(rest[0].EstadoResolucionId), &estado); err != nil {
				panic(err.Error())
			}
		} else {
			continue
		}

		if err = GetRequestLegacy("UrlcrudOikos", "dependencia/"+strconv.Itoa(resv[i].FacultadId), &dep); err != nil {
			panic(err.Error())
		}

		if err = GetRequestNew("UrlcrudParametros", "parametro/"+strconv.Itoa(res[i].TipoResolucionId), &tipo); err != nil {
			panic(err.Error())
		}

		resolucion := &models.Resoluciones{
			Id:               res[i].Id,
			NumeroResolucion: res[i].NumeroResolucion,
			Vigencia:         res[i].Vigencia,
			Periodo:          res[i].Periodo,
			Semanas:          res[i].NumeroSemanas,
			NivelAcademico:   resv[i].NivelAcademico,
			Dedicacion:       resv[i].Dedicacion,
			Facultad:         dep.Nombre,
			Estado:           estado.Nombre,
			TipoResolucion:   tipo.Nombre,
		}

		listaRes = append(listaRes, *resolucion)
	}

	if listaRes == nil {
		listaRes = []models.Resoluciones{}
	}

	return listaRes, total, outputError
}

func ListarResolucionesExpedidas(limit, offset int) (listaRes []models.Resoluciones, total int, outputError map[string]interface{}) {
	defer func() {
		if err := recover(); err != nil {
			outputError = map[string]interface{}{"funcion": "ListarResolucionesExpedidas", "err": err, "status": "500"}
			panic(outputError)
		}
	}()
	var estado []models.Parametro
	var resv models.ResolucionVinculacionDocente
	var rest []models.ResolucionEstado
	var dep models.Dependencia
	var tipo models.Parametro
	var err error

	url := "parametro?query=CodigoAbreviacion:REXP"
	if err = GetRequestNew("UrlcrudParametros", url, &estado); err != nil {
		panic(err.Error())
	}

	if len(estado) > 0 {
		url1 := "resolucion_estado?fields=Id&limit=0&query=Activo:true,EstadoResolucionId:" + strconv.Itoa(estado[0].Id)
		if err = GetRequestNew("UrlCrudResoluciones", url1, &rest); err != nil {
			panic(err.Error())
		}
		total = len(rest)

		url2 := "resolucion_estado?order=desc&sortby=Id&query=Activo:true,EstadoResolucionId:" + strconv.Itoa(estado[0].Id) + "&limit=" + strconv.Itoa(limit) + "&offset=" + strconv.Itoa(10*(offset-1))
		if err = GetRequestNew("UrlCrudResoluciones", url2, &rest); err != nil {
			panic(err.Error())
		}
		for i := range rest {
			url3 := "resolucion_vinculacion_docente/" + strconv.Itoa(rest[i].ResolucionId.Id)
			if err = GetRequestNew("UrlCrudResoluciones", url3, &resv); err != nil {
				panic(err.Error())
			}

			if err = GetRequestLegacy("UrlcrudOikos", "dependencia/"+strconv.Itoa(resv.FacultadId), &dep); err != nil {
				panic(err.Error())
			}

			if err = GetRequestNew("UrlcrudParametros", "parametro/"+strconv.Itoa(rest[i].ResolucionId.TipoResolucionId), &tipo); err != nil {
				panic(err.Error())
			}

			resolucion := &models.Resoluciones{
				Id:               rest[i].ResolucionId.Id,
				NumeroResolucion: rest[i].ResolucionId.NumeroResolucion,
				Vigencia:         rest[i].ResolucionId.Vigencia,
				Periodo:          rest[i].ResolucionId.Periodo,
				Semanas:          rest[i].ResolucionId.NumeroSemanas,
				NivelAcademico:   resv.NivelAcademico,
				Dedicacion:       resv.Dedicacion,
				Facultad:         dep.Nombre,
				Estado:           estado[0].Nombre,
				TipoResolucion:   tipo.Nombre,
			}

			listaRes = append(listaRes, *resolucion)
		}
	}

	if listaRes == nil {
		listaRes = []models.Resoluciones{}
	}

	return listaRes, total, outputError
}

func CargarResolucionCompleta(ResolucionId int) (resolucion models.ContenidoResolucion, outputError map[string]interface{}) {
	defer func() {
		if err := recover(); err != nil {
			outputError = map[string]interface{}{"funcion": "CargarResolucionCompleta", "err": err, "status": "500"}
			panic(outputError)
		}
	}()

	var err error

	if err = GetRequestNew("UrlCrudResoluciones", "resolucion/"+strconv.Itoa(ResolucionId), &resolucion.Resolucion); err != nil {
		panic(err.Error())
	}
	if err = GetRequestNew("UrlCrudResoluciones", "resolucion_vinculacion_docente/"+strconv.Itoa(ResolucionId), &resolucion.Vinculacion); err != nil {
		panic(err.Error())
	}
	if resolucion.Articulos, err = CargarArticulos(ResolucionId); err != nil {
		panic(err.Error())
	}

	return resolucion, outputError
}

func CargarArticulos(ResolucionId int) (articulos []models.Articulo, err error) {
	url := "componente_resolucion?limit=0&sortby=Numero&order=asc&query=TipoComponente:Articulo,Activo:true,ResolucionId.Id:" + strconv.Itoa(ResolucionId)
	var arts []models.ComponenteResolucion
	if err = GetRequestNew("UrlCrudResoluciones", url, &arts); err != nil {
		return articulos, err
	}
	for _, art := range arts {
		var parag []models.ComponenteResolucion
		url = "componente_resolucion?limit=0&sortby=Numero&order=asc&query=Activo:true,ComponenteResolucionPadre.Id:" + strconv.Itoa(art.Id)
		if err = GetRequestNew("UrlCrudResoluciones", url, &parag); err != nil {
			return articulos, err
		}
		articulo := &models.Articulo{
			Articulo:   art,
			Paragrafos: parag,
		}
		articulos = append(articulos, *articulo)
	}
	return articulos, err
}

func ActualizarResolucionCompleta(r models.ContenidoResolucion) (outputError map[string]interface{}) {
	defer func() {
		if err := recover(); err != nil {
			outputError = map[string]interface{}{"funcion": "ActualizarResolucionCompleta", "err": err, "status": "500"}
			panic(outputError)
		}
	}()

	var err error
	var respuesta map[string]interface{}
	var respuesta2 map[string]interface{}

	url := "resolucion/" + strconv.Itoa(r.Resolucion.Id)
	if err = SendRequestNew("UrlCrudResoluciones", url, "PUT", &respuesta, &r.Resolucion); err != nil {
		logs.Error(err)
		panic(err.Error())
	}
	url = "resolucion_vinculacion_docente/" + strconv.Itoa(r.Vinculacion.Id)
	if err = SendRequestNew("UrlCrudResoluciones", url, "PUT", &respuesta2, &r.Vinculacion); err != nil {
		logs.Error(err)
		panic(err.Error())
	}
	if err2 := ModificarArticulos(r.Articulos, r.Resolucion.Id); err2 != nil {
		logs.Error(err)
		panic(err.Error())
	}

	return outputError
}

func ModificarArticulos(artNuevos []models.Articulo, resolucionId int) (err error) {
	if artAnteriores, err2 := CargarArticulos(resolucionId); err2 == nil {
		if !iguales(artAnteriores, artNuevos) {
			if err = EliminarArticulos(artAnteriores); err != nil {
				return err
			}
			if err = InsertarArticulos(artNuevos, resolucionId); err != nil {
				return err
			}
		}
	} else {
		return err2
	}
	return nil
}

func EliminarArticulos(articulos []models.Articulo) (err error) {
	var resp map[string]interface{}
	for _, articulo := range articulos {
		if articulo.Paragrafos != nil {
			for _, paragrafoOld := range articulo.Paragrafos {
				err := SendRequestNew("UrlCrudResoluciones", "componente_resolucion/"+strconv.Itoa(paragrafoOld.Id), "DELETE", resp, nil)
				if err != nil {
					return err
				}
			}
		}
		err2 := SendRequestNew("UrlCrudResoluciones", "componente_resolucion/"+strconv.Itoa(articulo.Articulo.Id), "DELETE", resp, nil)
		if err2 != nil {
			return err2
		}
	}
	return nil
}

func CambiarEstadoResolucion(resolucionId int, estado, usuario string) (err error) {
	var objEstado []models.Parametro
	url := "parametro?query=CodigoAbreviacion:" + estado
	if err = GetRequestNew("UrlcrudParametros", url, &objEstado); err != nil {
		return err
	}

	var estados []models.ResolucionEstado
	url2 := "resolucion_estado?order=desc&sortby=Id&limit=0&query=Activo:true,ResolucionId.Id:" + strconv.Itoa(resolucionId)
	if err = GetRequestNew("UrlcrudResoluciones", url2, &estados); err != nil {
		return err
	}

	if len(estados) > 0 {
		estadoAnterior := estados[0]
		estadoAnterior.Activo = false
		url3 := "resolucion_estado/" + strconv.Itoa(estadoAnterior.Id)
		if err = SendRequestNew("UrlcrudResoluciones", url3, "PUT", &estadoAnterior, &estadoAnterior); err != nil {
			return err
		}
	}

	var respNuevoEstado models.ResolucionEstado
	nuevoEstado := models.ResolucionEstado{}
	nuevoEstado.Activo = true
	nuevoEstado.EstadoResolucionId = objEstado[0].Id
	nuevoEstado.ResolucionId = &models.Resolucion{Id: resolucionId}
	nuevoEstado.Usuario = usuario
	err = SendRequestNew("UrlcrudResoluciones", "resolucion_estado", "POST", &respNuevoEstado, &nuevoEstado)
	if err != nil {
		return err
	}

	return nil
}

func AnularResolucion(ResolucionId int) (outputError map[string]interface{}) {
	defer func() {
		if err := recover(); err != nil {
			outputError = map[string]interface{}{"funcion": "AnularResolucion", "err": err, "status": "500"}
			panic(outputError)
		}
	}()

	if resolucion, err := CargarResolucionCompleta(ResolucionId); err != nil {
		panic(err)
	} else {
		resolucion.Resolucion.Activo = false
		resolucion.Vinculacion.Activo = false
		for _, art := range resolucion.Articulos {
			art.Articulo.Activo = false
			for _, par := range art.Paragrafos {
				par.Activo = false
			}
		}
		if err2 := ActualizarResolucionCompleta(resolucion); err2 != nil {
			panic(err2)
		}
		if err3 := CambiarEstadoResolucion(ResolucionId, "RANU", ""); err3 != nil {
			logs.Error(err3)
			panic(err3.Error())
		}
	}
	if vinculaciones, err4 := ListarVinculaciones(strconv.Itoa(ResolucionId)); err4 != nil {
		panic(err4)
	} else {
		if len(vinculaciones) > 0 {
			if err5 := RetirarVinculaciones(vinculaciones); err5 != nil {
				panic(err5)
			}
		}
	}

	return nil
}

func ConsultaDocente(DocenteId int) (listaRes []models.Resoluciones, outputError map[string]interface{}) {
	defer func() {
		if err := recover(); err != nil {
			outputError = map[string]interface{}{"funcion": "ConsultaDocente", "err": err, "status": "500"}
			panic(outputError)
		}
	}()

	var vinculaciones []models.VinculacionDocente
	var rest []models.ResolucionEstado
	var estado []models.Parametro
	var dep models.Dependencia
	var tipo models.Parametro
	var err error

	url := "parametro?query=CodigoAbreviacion:REXP"
	if err = GetRequestNew("UrlcrudParametros", url, &estado); err != nil {
		panic(err.Error())
	}

	url1 := "vinculacion_docente?limit=0&query=Activo:true,PersonaId:" + strconv.Itoa(DocenteId)
	if err = GetRequestNew("UrlcrudResoluciones", url1, &vinculaciones); err != nil {
		panic(err.Error())
	}
	if len(estado) > 0 {
		for i := range vinculaciones {
			url2 := "resolucion_estado?query=Activo:true,EstadoResolucionId:" + strconv.Itoa(estado[0].Id) + ",ResolucionId.Id:" + strconv.Itoa(vinculaciones[i].ResolucionVinculacionDocenteId.Id)
			if err = GetRequestNew("UrlCrudResoluciones", url2, &rest); err != nil {
				panic(err.Error())
			}

			if err = GetRequestLegacy("UrlcrudOikos", "dependencia/"+strconv.Itoa(vinculaciones[i].ResolucionVinculacionDocenteId.FacultadId), &dep); err != nil {
				panic(err.Error())
			}

			if len(rest) > 0 {
				if err = GetRequestNew("UrlcrudParametros", "parametro/"+strconv.Itoa(rest[0].ResolucionId.TipoResolucionId), &tipo); err != nil {
					panic(err.Error())
				}

				resolucion := &models.Resoluciones{
					Id:               rest[0].ResolucionId.Id,
					NumeroResolucion: rest[0].ResolucionId.NumeroResolucion,
					Vigencia:         rest[0].ResolucionId.Vigencia,
					Periodo:          rest[0].ResolucionId.Periodo,
					Semanas:          rest[0].ResolucionId.NumeroSemanas,
					NivelAcademico:   vinculaciones[i].ResolucionVinculacionDocenteId.NivelAcademico,
					Dedicacion:       vinculaciones[i].ResolucionVinculacionDocenteId.Dedicacion,
					Facultad:         dep.Nombre,
					Estado:           estado[0].Nombre,
					TipoResolucion:   tipo.Nombre,
				}

				listaRes = append(listaRes, *resolucion)
			}
		}
	}

	if listaRes == nil {
		listaRes = []models.Resoluciones{}
	}

	return listaRes, outputError
}
