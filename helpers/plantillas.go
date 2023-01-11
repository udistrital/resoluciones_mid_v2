package helpers

import (
	"strconv"

	"github.com/astaxie/beego/logs"
	"github.com/udistrital/resoluciones_mid_v2/models"
)

func InsertarPlantilla(plantilla models.ContenidoResolucion) (PlantillaId int, outputError map[string]interface{}) {
	defer func() {
		if err := recover(); err != nil {
			outputError = map[string]interface{}{"funcion": "InsertarPlantilla", "err": err, "status": "500"}
			panic(outputError)
		}
	}()

	var tipos []models.Parametro
	url := "parametro?limit=0&query=ParametroPadreId.Id:" + strconv.Itoa(plantilla.Resolucion.TipoResolucionId)
	if err := GetRequestNew("UrlcrudParametros", url, &tipos); err != nil {
		logs.Error(err)
		panic(err.Error())
	}

	if existe, id := validarExistenciaPlantilla(plantilla, tipos); existe {
		var p *[]models.Resolucion
		url := "resolucion?query=Activo:true,TipoResolucionId:" + strconv.Itoa(id)
		if err4 := GetRequestNew("UrlCrudResoluciones", url, &p); err4 == nil && p != nil && len(*p) > 0 {
			return 0, nil
		} else if err4 != nil {
			panic(err4.Error())
		} else {
			plantilla.Resolucion.TipoResolucionId = id
		}
	} else {
		var err2 error
		if plantilla.Resolucion.TipoResolucionId, err2 = insertarTipoPlantilla(plantilla); err2 != nil {
			logs.Error(err2)
			panic("insertarTipoPlantilla -> " + err2.Error())
		}
	}

	var err3 map[string]interface{}
	if PlantillaId, err3 = InsertarResolucionCompleta(plantilla); err3 != nil {
		logs.Error(err3)
		panic(err3)
	}

	return PlantillaId, outputError
}

func ActualizarPlantilla(plantilla models.ContenidoResolucion) (outputError map[string]interface{}) {
	defer func() {
		if err := recover(); err != nil {
			outputError = map[string]interface{}{"funcion": "ActualizarPlantilla", "err": err, "status": "500"}
			panic(outputError)
		}
	}()

	var tipos []models.Parametro
	url := "parametro?limit=0&query=ParametroPadreId.Id:" + strconv.Itoa(plantilla.Resolucion.TipoResolucionId)
	if err := GetRequestNew("UrlcrudParametros", url, &tipos); err != nil {
		logs.Error(err)
		panic(err.Error())
	}

	if existe, id := validarExistenciaPlantilla(plantilla, tipos); existe {
		plantilla.Resolucion.TipoResolucionId = id
		if err2 := ActualizarResolucionCompleta(plantilla); err2 != nil {
			panic(err2)
		}
	}

	return nil
}

func CargarPlantilla(PlantillaId int) (plantilla models.ContenidoResolucion, outputError map[string]interface{}) {
	defer func() {
		if err := recover(); err != nil {
			outputError = map[string]interface{}{"funcion": "CargarPlantilla", "err": err, "status": "500"}
			panic(outputError)
		}
	}()

	var err map[string]interface{}
	if plantilla, err = CargarResolucionCompleta(PlantillaId); err != nil {
		logs.Error(err)
		panic(err)
	}

	var tipoRes models.Parametro
	if err2 := GetRequestNew("UrlcrudParametros", ParametroEndpoint+strconv.Itoa(plantilla.Resolucion.TipoResolucionId), &tipoRes); err2 != nil {
		logs.Error(err2)
		panic(err2.Error())
	}
	plantilla.Resolucion.TipoResolucionId = tipoRes.ParametroPadreId.Id

	return plantilla, outputError
}

func ListarPlantillas() (lista []models.Plantilla, outputError map[string]interface{}) {
	defer func() {
		if err := recover(); err != nil {
			outputError = map[string]interface{}{"funcion": "ListarPlantillas", "err": err, "status": "500"}
			panic(outputError)
		}
	}()
	var res []models.Resolucion
	var resv models.ResolucionVinculacionDocente
	var dep models.Dependencia
	var tipos []models.Parametro

	url := "parametro?limit=0&query=CodigoAbreviacion:RTP,Activo:true"
	if err := GetRequestNew("UrlcrudParametros", url, &tipos); err != nil {
		logs.Error(err)
		panic(err.Error())
	}

	for _, tipoPlantilla := range tipos {
		url := "resolucion?query=Activo:true,TipoResolucionId:" + strconv.Itoa(tipoPlantilla.Id)
		if err2 := GetRequestNew("UrlCrudResoluciones", url, &res); err2 != nil {
			logs.Error(err2)
			panic(err2.Error())
		}
		if len(res) != 0 {
			if err3 := GetRequestNew("UrlCrudResoluciones", ResVinEndpoint+strconv.Itoa(res[0].Id), &resv); err3 != nil {
				logs.Error(err3)
				panic(err3.Error())
			}
			if err4 := GetRequestLegacy("UrlcrudOikos", "dependencia/"+strconv.Itoa(resv.FacultadId), &dep); err4 != nil {
				logs.Error(err4)
				panic(err4.Error())
			}
			plantilla := models.Plantilla{
				Id:             res[0].Id,
				Dedicacion:     resv.Dedicacion,
				NivelAcademico: resv.NivelAcademico,
				Facultad:       dep.Nombre,
				TipoResolucion: tipoPlantilla.ParametroPadreId.Nombre,
			}
			lista = append(lista, plantilla)
		}
	}
	if lista == nil {
		lista = []models.Plantilla{}
	}
	return lista, outputError
}

func BorrarPlantilla(PlantillaId int) (outputError map[string]interface{}) {
	defer func() {
		if err := recover(); err != nil {
			outputError = map[string]interface{}{"funcion": "BorrarPlantilla", "err": err, "status": "500"}
			panic(outputError)
		}
	}()

	var err error
	var respuesta map[string]interface{}
	var respuesta2 map[string]interface{}
	if plantilla, err2 := CargarResolucionCompleta(PlantillaId); err2 == nil {
		plantilla.Resolucion.Activo = false
		plantilla.Vinculacion.Activo = false

		url := "resolucion/" + strconv.Itoa(plantilla.Resolucion.Id)
		if err = SendRequestNew("UrlCrudResoluciones", url, "PUT", &respuesta, &plantilla.Resolucion); err != nil {
			logs.Error(err)
			panic(err.Error())
		}
		url = ResVinEndpoint + strconv.Itoa(plantilla.Vinculacion.Id)
		if err = SendRequestNew("UrlCrudResoluciones", url, "PUT", &respuesta2, &plantilla.Vinculacion); err != nil {
			logs.Error(err)
			panic(err.Error())
		}

		if err = EliminarArticulos(plantilla.Articulos); err != nil {
			logs.Error(err)
			panic(err.Error())
		}
	} else {
		panic(err2)
	}

	return nil
}

func insertarTipoPlantilla(plantilla models.ContenidoResolucion) (TipoId int, err error) {
	var tipo []models.TipoParametro
	if err = GetRequestNew("UrlcrudParametros", "tipo_parametro?query=CodigoAbreviacion:TR", &tipo); err != nil {
		return TipoId, err
	}
	nombre := plantilla.Vinculacion.NivelAcademico + " " + plantilla.Vinculacion.Dedicacion + " " + strconv.Itoa(plantilla.Vinculacion.FacultadId)

	var resp models.Parametro
	tipoPlantilla := models.Parametro{
		Nombre:            nombre,
		Descripcion:       nombre,
		CodigoAbreviacion: "RTP",
		Activo:            true,
		TipoParametroId:   &models.TipoParametro{Id: tipo[0].Id},
		ParametroPadreId:  &models.Parametro{Id: plantilla.Resolucion.TipoResolucionId},
	}

	if err = SendRequestNew("UrlcrudParametros", "parametro", "POST", &resp, &tipoPlantilla); err == nil {
		TipoId = resp.Id
	}
	return TipoId, err
}

func validarExistenciaPlantilla(plantilla models.ContenidoResolucion, tipos []models.Parametro) (existe bool, id int) {
	existe = false
	nombre := plantilla.Vinculacion.NivelAcademico + " " + plantilla.Vinculacion.Dedicacion + " " + strconv.Itoa(plantilla.Vinculacion.FacultadId)
	for _, tipo := range tipos {
		existe = existe || nombre == tipo.Nombre
		if existe {
			id = tipo.Id
			break
		}
	}
	return existe, id
}
