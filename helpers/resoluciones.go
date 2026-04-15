package helpers

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"
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
	url2 := ParametroEndpoint + strconv.Itoa(nuevaRes.Resolucion.TipoResolucionId)
	if err := GetRequestNew("UrlcrudParametros", url2, &tipoRes); err != nil {
		logs.Error(err)
		panic(err.Error())
	}

	if tipoRes.CodigoAbreviacion != "RVIN" {
		var anteriorResvin models.ResolucionVinculacionDocente
		var anteriorRes models.Resolucion
		url := ResVinEndpoint + strconv.Itoa(nuevaRes.ResolucionAnteriorId)
		if er := GetRequestNew("UrlcrudResoluciones", url, &anteriorResvin); er != nil {
			logs.Error(er)
			panic(er.Error())
		}
		url = ResolucionEndpoint + strconv.Itoa(nuevaRes.ResolucionAnteriorId)
		if err := GetRequestNew("UrlcrudResoluciones", url, &anteriorRes); err != nil {
			logs.Error(err)
			panic(err.Error())
		}
		nuevaRes.Vinculacion.Dedicacion = anteriorResvin.Dedicacion
		nuevaRes.Vinculacion.NivelAcademico = anteriorResvin.NivelAcademico
		nuevaRes.Vinculacion.FacultadId = anteriorResvin.FacultadId
		nuevaRes.Resolucion.DependenciaId = anteriorRes.DependenciaId
		nuevaRes.Resolucion.DependenciaFirmaId = anteriorRes.DependenciaFirmaId
		nuevaRes.Resolucion.Periodo = anteriorRes.Periodo
		nuevaRes.Resolucion.NumeroSemanas = anteriorRes.NumeroSemanas
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
		} else if p == nil || len(*p) == 0 {
			return 0, nil
		}
	} else {
		return 0, nil
	}

	nuevaRes.Resolucion.Titulo = plantilla.Resolucion.Titulo
	nuevaRes.Resolucion.PreambuloResolucion = plantilla.Resolucion.PreambuloResolucion
	nuevaRes.Resolucion.ConsideracionResolucion = plantilla.Resolucion.ConsideracionResolucion
	nuevaRes.Resolucion.CuadroResponsabilidades = plantilla.Resolucion.CuadroResponsabilidades
	nuevaRes.Resolucion.FechaInicio = plantilla.Resolucion.FechaInicio
	nuevaRes.Resolucion.FechaFin = plantilla.Resolucion.FechaFin
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

	if err5 := CambiarEstadoResolucion(resolucionId, "RSOL", nuevaRes.Usuario); err5 != nil {
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
	v.Resolucion.Vigencia, _, _ = time.Now().Date()
	if v.Resolucion.VigenciaCarga == 0 {
		v.Resolucion.VigenciaCarga = v.Resolucion.Vigencia
	}
	if v.Resolucion.PeriodoCarga == 0 {
		v.Resolucion.PeriodoCarga = v.Resolucion.Periodo
	}

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

// Genera un listado de resoluciones teniendo en cuenta la query
func ListarResoluciones(query string) (listaRes []models.Resoluciones, outputError map[string]interface{}) {
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

	if err = GetRequestNew("UrlcrudResoluciones", "resolucion"+query, &res); err != nil {
		logs.Error(err)
		panic(err.Error())
	}

	if err = GetRequestNew("UrlcrudResoluciones", "resolucion_vinculacion_docente"+query, &resv); err != nil {
		logs.Error(err)
		panic(err.Error())
	}

	for i := range res {
		url3 := "resolucion_estado?order=desc&sortby=Id&query=Activo:true,ResolucionId.Id:" + strconv.Itoa(res[i].Id)
		if err = GetRequestNew("UrlcrudResoluciones", url3, &rest); err != nil {
			panic(err.Error())
		}

		if len(rest) > 0 {
			if err = GetRequestNew("UrlcrudParametros", ParametroEndpoint+strconv.Itoa(rest[0].EstadoResolucionId), &estado); err != nil {
				panic(err.Error())
			}
		} else {
			continue
		}

		if err = GetRequestLegacy("UrlcrudOikos", "dependencia/"+strconv.Itoa(resv[i].FacultadId), &dep); err != nil {
			panic(err.Error())
		}

		if err = GetRequestNew("UrlcrudParametros", ParametroEndpoint+strconv.Itoa(res[i].TipoResolucionId), &tipo); err != nil {
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
	return listaRes, outputError
}

func resolverTipoResolucionFiltro(f *models.Filtro) error {
	var parametros []models.Parametro

	if len(f.TipoResolucion) > 0 {
		params := url.Values{}
		params.Add("query", "Nombre:"+f.TipoResolucion)
		ruta := "parametro?" + params.Encode()
		if err := GetRequestNew("UrlcrudParametros", ruta, &parametros); err != nil {
			logs.Error(err)
			return err
		}
		if len(parametros) > 0 {
			f.TipoResolucion = strconv.Itoa(parametros[0].Id)
		}
	}

	if len(f.ExcluirTipo) > 0 || len(f.TipoResolucion) == 0 {
		ruta := "parametro?query=TipoParametroId.CodigoAbreviacion:TR"
		if err := GetRequestNew("UrlcrudParametros", ruta, &parametros); err != nil {
			logs.Error(err)
			return err
		}

		ids := make([]string, 0, len(parametros))
		for _, param := range parametros {
			if param.ParametroPadreId != nil {
				continue
			}
			if len(f.ExcluirTipo) > 0 && param.CodigoAbreviacion == f.ExcluirTipo {
				continue
			}
			ids = append(ids, strconv.Itoa(param.Id))
		}

		if len(f.ExcluirTipo) > 0 || len(f.TipoResolucion) == 0 {
			f.TipoResolucion = strings.Join(ids, "|")
		}
	}

	return nil
}

func resolverFiltroFacultad(facultad string) (string, error) {
	if len(facultad) == 0 {
		return "", nil
	}

	if _, err := strconv.Atoi(facultad); err == nil {
		return "DependenciaId:" + facultad, nil
	}

	var dependencias []models.Dependencia
	rutaDependencias := "dependencia?limit=0&query=Nombre__icontains:" + url.QueryEscape(facultad)
	if err := GetRequestLegacy("UrlcrudOikos", rutaDependencias, &dependencias); err != nil {
		logs.Error(err)
		return "", err
	}

	if len(dependencias) == 0 {
		return "", nil
	}

	ids := make([]string, 0, len(dependencias))
	for _, dependencia := range dependencias {
		ids = append(ids, strconv.Itoa(dependencia.Id))
	}

	return "DependenciaId.in:" + strings.Join(ids, "|"), nil
}

func appendResolucionFilters(query string, f models.Filtro, filtroFacultad string, nested bool) string {
	prefix := ""
	if nested {
		prefix = "ResolucionId."
	}

	if len(f.NumeroResolucion) > 0 {
		query += prefix + "NumeroResolucion:" + f.NumeroResolucion + ","
	}
	if len(f.Vigencia) > 0 {
		query += prefix + "Vigencia:" + f.Vigencia + ","
	}
	if len(f.Periodo) > 0 {
		query += prefix + "Periodo:" + f.Periodo + ","
	}
	if len(f.Semanas) > 0 {
		query += prefix + "NumeroSemanas:" + f.Semanas + ","
	}
	if filtroFacultad != "" {
		query += prefix + filtroFacultad + ","
	}
	if len(f.TipoResolucion) > 0 {
		query += prefix + "TipoResolucionId.in:" + f.TipoResolucion + ","
	}

	return strings.TrimSuffix(query, ",")
}

func consultarIdsResolucionPorEstado(f models.Filtro, queryBase string, filtroFacultad string) ([]string, error) {
	var parametros []models.Parametro
	var estados []models.ResolucionEstado

	param := url.Values{}
	param.Add("query", "Nombre.in:"+f.Estado)
	ruta := "parametro?limit=0&" + param.Encode()
	if err := GetRequestNew("UrlcrudParametros", ruta, &parametros); err != nil {
		return nil, err
	}

	idsEstado := make([]string, 0, len(parametros))
	for _, parametro := range parametros {
		idsEstado = append(idsEstado, strconv.Itoa(parametro.Id))
	}

	query := queryBase + "EstadoResolucionId.in:" + strings.Join(idsEstado, "|") + ","
	query += "ResolucionId.Activo:true,"
	query = appendResolucionFilters(query, f, filtroFacultad, true)

	ruta = "resolucion_estado?" + query + "&limit=0&fields=ResolucionId"
	if err := GetRequestNew("UrlcrudResoluciones", ruta, &estados); err != nil {
		logs.Error(err)
		return nil, err
	}

	return extraerIdsResolucionDesdeEstado(estados), nil
}

func consultarIdsResolucion(f models.Filtro, queryBase string, filtroFacultad string) ([]string, error) {
	var resoluciones []models.Resolucion

	query := appendResolucionFilters(queryBase, f, filtroFacultad, false)
	ruta := "resolucion?" + query + "&limit=0&fields=Id"
	if err := GetRequestNew("UrlcrudResoluciones", ruta, &resoluciones); err != nil {
		logs.Error(err)
		return nil, err
	}

	return extraerIdsResolucion(resoluciones), nil
}

func consultarIdsResolucionVinculacion(f models.Filtro, queryBase string, resolucionIds []string) ([]string, error) {
	var vinculaciones []models.ResolucionVinculacionDocente

	query := queryBase
	if len(f.Dedicacion) > 0 {
		query += "Dedicacion:" + f.Dedicacion + ","
	}
	if len(f.NivelAcademico) > 0 {
		query += "NivelAcademico:" + f.NivelAcademico + ","
	}
	query += "Id.in:" + strings.Join(resolucionIds, "|")
	query = strings.TrimSuffix(query, ",")

	ruta := "resolucion_vinculacion_docente?" + query + "&limit=0&fields=Id"
	if err := GetRequestNew("UrlCrudResoluciones", ruta, &vinculaciones); err != nil {
		logs.Error(err)
		return nil, err
	}

	return extraerIdsResolucionVinculacion(vinculaciones), nil
}

func extraerIdsResolucionDesdeEstado(estados []models.ResolucionEstado) []string {
	ids := make([]string, 0, len(estados))
	for i := range estados {
		ids = append(ids, strconv.Itoa(estados[i].ResolucionId.Id))
	}
	return ids
}

func extraerIdsResolucion(resoluciones []models.Resolucion) []string {
	ids := make([]string, 0, len(resoluciones))
	for i := range resoluciones {
		ids = append(ids, strconv.Itoa(resoluciones[i].Id))
	}
	return ids
}

func extraerIdsResolucionVinculacion(vinculaciones []models.ResolucionVinculacionDocente) []string {
	ids := make([]string, 0, len(vinculaciones))
	for i := range vinculaciones {
		ids = append(ids, strconv.Itoa(vinculaciones[i].Id))
	}
	return ids
}

func calcularVentanaConsulta(limit int, offset int) int {
	return (limit * (offset - 1)) + 100
}

func limitarIdsConsulta(ids []string, maximo int) []string {
	if len(ids) < maximo {
		maximo = len(ids)
	}
	return ids[:maximo]
}

func construirQueryFinalResoluciones(limit int, offset int, queryBase string, ids []string) string {
	listadoCompleto := fmt.Sprintf("Id.in:%s", strings.Join(ids, "|"))
	return fmt.Sprintf("?limit=%d&offset=%d&%s%s", limit, limit*(offset-1), queryBase, listadoCompleto)
}

func cargarResolucionesPaginadas(limit int, offset int, queryBase string, ids []string) ([]models.Resoluciones, map[string]interface{}) {
	queryFinal := construirQueryFinalResoluciones(limit, offset, queryBase, ids)
	return ListarResoluciones(queryFinal)
}

func obtenerIdsResolucionFiltrados(f *models.Filtro, queryBase string) ([]string, error) {
	if err := resolverTipoResolucionFiltro(f); err != nil {
		return nil, err
	}

	filtroFacultadResolucion, err := resolverFiltroFacultad(f.FacultadId)
	if err != nil {
		return nil, err
	}
	if len(f.FacultadId) > 0 && filtroFacultadResolucion == "" {
		return []string{}, nil
	}

	if len(f.Estado) > 0 {
		return consultarIdsResolucionPorEstado(*f, queryBase, filtroFacultadResolucion)
	}

	return consultarIdsResolucion(*f, queryBase, filtroFacultadResolucion)
}

func obtenerIdsResolucionVinculacionFiltrados(f models.Filtro, queryBase string, ids []string) ([]string, error) {
	if len(ids) == 0 {
		return []string{}, nil
	}

	return consultarIdsResolucionVinculacion(f, queryBase, ids)
}

func listarResolucionesVaciasSiNoHayIds(ids []string, outputError map[string]interface{}) ([]models.Resoluciones, int, bool) {
	if len(ids) == 0 {
		return []models.Resoluciones{}, 0, true
	}

	return nil, 0, false
}

// Procesa los filtros de busqueda por medio de los parametros de consulta y genera la query correspondiente
func ListarResolucionesFiltradas(f models.Filtro) (listaRes []models.Resoluciones, total int, outputError map[string]interface{}) {
	defer func() {
		if err := recover(); err != nil {
			outputError = map[string]interface{}{"funcion": "ListarResolucionesFiltradas", "err": err, "status": "500"}
			panic(outputError)
		}
	}()

	var listado []string
	var err2 map[string]interface{}

	limit, _ := strconv.Atoi(f.Limit)
	offset, _ := strconv.Atoi(f.Offset)
	queryBase := "order=desc&sortby=Id&query=Activo:true,"
	offset2 := calcularVentanaConsulta(limit, offset)

	var err error
	if listado, err = obtenerIdsResolucionFiltrados(&f, queryBase); err != nil {
		panic(err.Error())
	}

	if vacio, totalVacio, ok := listarResolucionesVaciasSiNoHayIds(listado, outputError); ok {
		return vacio, totalVacio, outputError
	}

	if listado, err = obtenerIdsResolucionVinculacionFiltrados(f, queryBase, listado); err != nil {
		panic(err.Error())
	}

	if vacio, totalVacio, ok := listarResolucionesVaciasSiNoHayIds(listado, outputError); ok {
		return vacio, totalVacio, outputError
	}

	total = len(listado)

	if listaRes, err2 = cargarResolucionesPaginadas(limit, offset, queryBase, limitarIdsConsulta(listado, offset2)); err2 != nil {
		panic(err2)
	}

	if listaRes == nil {
		listaRes = []models.Resoluciones{}
	}

	return listaRes, total, outputError
}

func GetResolucionVinculacionDocente(id int) (rv models.ResolucionVinculacionDocente) {
	url := ResVinEndpoint + strconv.Itoa(id)
	if err := GetRequestNew("UrlCrudResoluciones", url, &rv); err != nil {
		panic("Error consultando resolucion_vinculacion_docente -> " + err.Error())
	}
	return
}

// Recupera toda la información de una resolución a partir de su Id
func CargarResolucionCompleta(ResolucionId int) (resolucion models.ContenidoResolucion, outputError map[string]interface{}) {
	defer func() {
		if err := recover(); err != nil {
			outputError = map[string]interface{}{"funcion": "CargarResolucionCompleta", "err": err, "status": "500"}
			panic(outputError)
		}
	}()

	var err error
	var modificacionResolucion []models.ModificacionResolucion
	var tipoResolucion models.Parametro

	if err = GetRequestNew("UrlCrudResoluciones", ResolucionEndpoint+strconv.Itoa(ResolucionId), &resolucion.Resolucion); err != nil {
		panic(err.Error())
	}
	if err = GetRequestNew("UrlCrudResoluciones", ResVinEndpoint+strconv.Itoa(ResolucionId), &resolucion.Vinculacion); err != nil {
		panic(err.Error())
	}
	if err := GetRequestNew("UrlcrudParametros", ParametroEndpoint+strconv.Itoa(resolucion.Resolucion.TipoResolucionId), &tipoResolucion); err != nil {
		panic(map[string]interface{}{"funcion": "/ConstruirDocumentoResolucion-param", "err": err.Error(), "status": "500"})
	}
	if tipoResolucion.CodigoAbreviacion != "RVIN" && tipoResolucion.CodigoAbreviacion != "RTP" {
		if err = GetRequestNew("UrlCrudResoluciones", "modificacion_resolucion?query=resolucionNuevaId:"+strconv.Itoa(ResolucionId), &modificacionResolucion); err != nil {
			panic(err.Error())
		}
		resolucion.ResolucionAnteriorId = modificacionResolucion[0].ResolucionAnteriorId.Id
	}
	if resolucion.Articulos, err = CargarArticulos(ResolucionId); err != nil {
		panic(err.Error())
	}

	return resolucion, outputError
}

// Recupera toda la información de los articulos y paragrafos de una resolución a partir de su Id
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

	url := ResolucionEndpoint + strconv.Itoa(r.Resolucion.Id)
	if err = SendRequestNew("UrlCrudResoluciones", url, "PUT", &respuesta, &r.Resolucion); err != nil {
		logs.Error(err)
		panic(err.Error())
	}
	url = ResVinEndpoint + strconv.Itoa(r.Vinculacion.Id)
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

	if nombreUsuario, err := GetUsuario(usuario); err == nil {
		var respNuevoEstado models.ResolucionEstado
		nuevoEstado := models.ResolucionEstado{}
		nuevoEstado.Activo = true
		nuevoEstado.EstadoResolucionId = objEstado[0].Id
		nuevoEstado.ResolucionId = &models.Resolucion{Id: resolucionId}

		var ok bool
		if nuevoEstado.Usuario, ok = nombreUsuario["sub"].(string); !ok {
			nuevoEstado.Usuario = ""
		}

		err = SendRequestNew("UrlcrudResoluciones", "resolucion_estado", "POST", &respNuevoEstado, &nuevoEstado)
		if err != nil {
			return err
		}
	} else {
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
	if vinculaciones, err4 := ListarVinculaciones(strconv.Itoa(ResolucionId), false); err4 != nil {
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

	url1 := "vinculacion_docente?limit=0&order=desc&sortby=Id&query=PersonaId:" + strconv.Itoa(DocenteId)
	if err = GetRequestNew("UrlcrudResoluciones", url1, &vinculaciones); err != nil {
		panic(err.Error())
	}
	if len(estado) > 0 {
		for i := range vinculaciones {
			url2 := "resolucion_estado?fields=fill&query=Activo:true,EstadoResolucionId:" + strconv.Itoa(estado[0].Id) + ",ResolucionId.Id:" + strconv.Itoa(vinculaciones[i].ResolucionVinculacionDocenteId.Id)
			if err = GetRequestNew("UrlCrudResoluciones", url2, &rest); err != nil {
				panic(err.Error())
			}

			if err = GetRequestLegacy("UrlcrudOikos", "dependencia/"+strconv.Itoa(vinculaciones[i].ResolucionVinculacionDocenteId.FacultadId), &dep); err != nil {
				panic(err.Error())
			}

			if len(rest) > 0 {
				if err = GetRequestNew("UrlcrudParametros", ParametroEndpoint+strconv.Itoa(rest[0].ResolucionId.TipoResolucionId), &tipo); err != nil {
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
