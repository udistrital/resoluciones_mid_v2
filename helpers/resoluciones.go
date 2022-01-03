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
	/* defer func() {
		if err := recover(); err != nil {
			outputError = map[string]interface{}{"funcion": "InsertarResolucion", "err": err, "status": "500"}
			panic(outputError)
		}
	}() */

	var plantillas []models.Parametro
	var plantilla models.ContenidoResolucion

	url := "parametro?limit=0&query=Activo:true,ParametroPadreId.Id:" + strconv.Itoa(nuevaRes.Resolucion.TipoResolucionId)
	if err := GetRequestNew("UrlcrudParametros", url, &plantillas); err != nil {
		logs.Error(err)
		panic(err.Error())
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
			TipoComponente: "Artículo",
			Activo:         true,
		}
		if err = SendRequestNew("UrlCrudResoluciones", "componente_resolucion", "POST", &art, &articulo); err == nil {
			for j, obj2 := range obj.Paragrafos {
				paragrafo := &models.ComponenteResolucion{
					Numero:                    j + 1,
					ResolucionId:              &models.Resolucion{Id: resolucionId},
					Texto:                     obj2.Texto,
					TipoComponente:            "Parágrafo",
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

func ListarResoluciones() (r []models.Resoluciones, outputError map[string]interface{}) {

	return r, outputError
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
	url := "componente_resolucion?sortby=Numero&order=asc&query=TipoComponente:Artículo,Activo:true,ResolucionId.Id:" + strconv.Itoa(ResolucionId)
	var arts []models.ComponenteResolucion
	var parag []models.ComponenteResolucion
	if err = GetRequestNew("UrlCrudResoluciones", url, &arts); err != nil {
		return articulos, err
	}
	for _, art := range arts {
		url = "componente_resolucion?sortby=Numero&order=asc&query=Activo:true,ComponenteResolucionPadre.Id:" + strconv.Itoa(art.Id)
		if err = GetRequestNew("UrlCrudResoluciones", url, &parag); err != nil {
			return articulos, err
		}
		articulos = append(articulos, models.Articulo{
			Articulo:   art,
			Paragrafos: parag,
		})
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
	url2 := "resolucion_estado?limit=0&query=Activo:true,ResolucionId.Id:" + strconv.Itoa(resolucionId)
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
