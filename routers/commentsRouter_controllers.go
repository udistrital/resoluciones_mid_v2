package routers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context/param"
)

func init() {

    beego.GlobalControllerRouter["github.com/udistrital/resoluciones_mid_v2/controllers:ExpedirResolucionController"] = append(beego.GlobalControllerRouter["github.com/udistrital/resoluciones_mid_v2/controllers:ExpedirResolucionController"],
        beego.ControllerComments{
            Method: "Cancelar",
            Router: "/cancelar",
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/resoluciones_mid_v2/controllers:ExpedirResolucionController"] = append(beego.GlobalControllerRouter["github.com/udistrital/resoluciones_mid_v2/controllers:ExpedirResolucionController"],
        beego.ControllerComments{
            Method: "Expedir",
            Router: "/expedir",
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/resoluciones_mid_v2/controllers:ExpedirResolucionController"] = append(beego.GlobalControllerRouter["github.com/udistrital/resoluciones_mid_v2/controllers:ExpedirResolucionController"],
        beego.ControllerComments{
            Method: "ExpedirModificacion",
            Router: "/expedirModificacion",
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/resoluciones_mid_v2/controllers:ExpedirResolucionController"] = append(beego.GlobalControllerRouter["github.com/udistrital/resoluciones_mid_v2/controllers:ExpedirResolucionController"],
        beego.ControllerComments{
            Method: "ValidarDatosExpedicion",
            Router: "/validar_datos_expedicion",
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/resoluciones_mid_v2/controllers:GestionPlantillasController"] = append(beego.GlobalControllerRouter["github.com/udistrital/resoluciones_mid_v2/controllers:GestionPlantillasController"],
        beego.ControllerComments{
            Method: "Post",
            Router: "/",
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/resoluciones_mid_v2/controllers:GestionPlantillasController"] = append(beego.GlobalControllerRouter["github.com/udistrital/resoluciones_mid_v2/controllers:GestionPlantillasController"],
        beego.ControllerComments{
            Method: "GetAll",
            Router: "/",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/resoluciones_mid_v2/controllers:GestionPlantillasController"] = append(beego.GlobalControllerRouter["github.com/udistrital/resoluciones_mid_v2/controllers:GestionPlantillasController"],
        beego.ControllerComments{
            Method: "GetOne",
            Router: "/:id",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/resoluciones_mid_v2/controllers:GestionPlantillasController"] = append(beego.GlobalControllerRouter["github.com/udistrital/resoluciones_mid_v2/controllers:GestionPlantillasController"],
        beego.ControllerComments{
            Method: "Put",
            Router: "/:id",
            AllowHTTPMethods: []string{"put"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/resoluciones_mid_v2/controllers:GestionPlantillasController"] = append(beego.GlobalControllerRouter["github.com/udistrital/resoluciones_mid_v2/controllers:GestionPlantillasController"],
        beego.ControllerComments{
            Method: "Delete",
            Router: "/:id",
            AllowHTTPMethods: []string{"delete"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/resoluciones_mid_v2/controllers:GestionResolucionesController"] = append(beego.GlobalControllerRouter["github.com/udistrital/resoluciones_mid_v2/controllers:GestionResolucionesController"],
        beego.ControllerComments{
            Method: "Post",
            Router: "/",
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/resoluciones_mid_v2/controllers:GestionResolucionesController"] = append(beego.GlobalControllerRouter["github.com/udistrital/resoluciones_mid_v2/controllers:GestionResolucionesController"],
        beego.ControllerComments{
            Method: "GetAll",
            Router: "/",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/resoluciones_mid_v2/controllers:GestionResolucionesController"] = append(beego.GlobalControllerRouter["github.com/udistrital/resoluciones_mid_v2/controllers:GestionResolucionesController"],
        beego.ControllerComments{
            Method: "GetOne",
            Router: "/:id",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/resoluciones_mid_v2/controllers:GestionResolucionesController"] = append(beego.GlobalControllerRouter["github.com/udistrital/resoluciones_mid_v2/controllers:GestionResolucionesController"],
        beego.ControllerComments{
            Method: "Put",
            Router: "/:id",
            AllowHTTPMethods: []string{"put"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/resoluciones_mid_v2/controllers:GestionResolucionesController"] = append(beego.GlobalControllerRouter["github.com/udistrital/resoluciones_mid_v2/controllers:GestionResolucionesController"],
        beego.ControllerComments{
            Method: "Delete",
            Router: "/:id",
            AllowHTTPMethods: []string{"delete"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/resoluciones_mid_v2/controllers:GestionResolucionesController"] = append(beego.GlobalControllerRouter["github.com/udistrital/resoluciones_mid_v2/controllers:GestionResolucionesController"],
        beego.ControllerComments{
            Method: "ActualizarEstado",
            Router: "/actualizar_estado",
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/resoluciones_mid_v2/controllers:GestionResolucionesController"] = append(beego.GlobalControllerRouter["github.com/udistrital/resoluciones_mid_v2/controllers:GestionResolucionesController"],
        beego.ControllerComments{
            Method: "ConsultaDocente",
            Router: "/consultar_docente/:id",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/resoluciones_mid_v2/controllers:GestionResolucionesController"] = append(beego.GlobalControllerRouter["github.com/udistrital/resoluciones_mid_v2/controllers:GestionResolucionesController"],
        beego.ControllerComments{
            Method: "GenerarResolucion",
            Router: "/generar_resolucion/:id",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/resoluciones_mid_v2/controllers:GestionResolucionesController"] = append(beego.GlobalControllerRouter["github.com/udistrital/resoluciones_mid_v2/controllers:GestionResolucionesController"],
        beego.ControllerComments{
            Method: "GetResolucionesAprobadas",
            Router: "/resoluciones_aprobadas",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/resoluciones_mid_v2/controllers:GestionResolucionesController"] = append(beego.GlobalControllerRouter["github.com/udistrital/resoluciones_mid_v2/controllers:GestionResolucionesController"],
        beego.ControllerComments{
            Method: "GetResolucionesExpedidas",
            Router: "/resoluciones_expedidas",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/resoluciones_mid_v2/controllers:GestionResolucionesController"] = append(beego.GlobalControllerRouter["github.com/udistrital/resoluciones_mid_v2/controllers:GestionResolucionesController"],
        beego.ControllerComments{
            Method: "GetResolucionesInscritas",
            Router: "/resoluciones_inscritas",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/resoluciones_mid_v2/controllers:GestionVinculacionesController"] = append(beego.GlobalControllerRouter["github.com/udistrital/resoluciones_mid_v2/controllers:GestionVinculacionesController"],
        beego.ControllerComments{
            Method: "Post",
            Router: "/",
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/resoluciones_mid_v2/controllers:GestionVinculacionesController"] = append(beego.GlobalControllerRouter["github.com/udistrital/resoluciones_mid_v2/controllers:GestionVinculacionesController"],
        beego.ControllerComments{
            Method: "DocentesPrevinculados",
            Router: "/:resolucion_id",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/resoluciones_mid_v2/controllers:GestionVinculacionesController"] = append(beego.GlobalControllerRouter["github.com/udistrital/resoluciones_mid_v2/controllers:GestionVinculacionesController"],
        beego.ControllerComments{
            Method: "CalcularValorContratosSeleccionados",
            Router: "/calcular_valor_contratos_seleccionados",
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/resoluciones_mid_v2/controllers:GestionVinculacionesController"] = append(beego.GlobalControllerRouter["github.com/udistrital/resoluciones_mid_v2/controllers:GestionVinculacionesController"],
        beego.ControllerComments{
            Method: "Desvincular",
            Router: "/desvincular",
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/resoluciones_mid_v2/controllers:GestionVinculacionesController"] = append(beego.GlobalControllerRouter["github.com/udistrital/resoluciones_mid_v2/controllers:GestionVinculacionesController"],
        beego.ControllerComments{
            Method: "DesvincularDocentes",
            Router: "/desvincular_docentes",
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/resoluciones_mid_v2/controllers:GestionVinculacionesController"] = append(beego.GlobalControllerRouter["github.com/udistrital/resoluciones_mid_v2/controllers:GestionVinculacionesController"],
        beego.ControllerComments{
            Method: "DocentesCargaHoraria",
            Router: "/docentes_carga_horaria/:vigencia/:periodo/:dedicacion/:facultad/:nivel_academico",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/resoluciones_mid_v2/controllers:GestionVinculacionesController"] = append(beego.GlobalControllerRouter["github.com/udistrital/resoluciones_mid_v2/controllers:GestionVinculacionesController"],
        beego.ControllerComments{
            Method: "InformeVinculaciones",
            Router: "/informe_vinculaciones",
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/resoluciones_mid_v2/controllers:GestionVinculacionesController"] = append(beego.GlobalControllerRouter["github.com/udistrital/resoluciones_mid_v2/controllers:GestionVinculacionesController"],
        beego.ControllerComments{
            Method: "ModificarVinculacion",
            Router: "/modificar_vinculacion",
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/resoluciones_mid_v2/controllers:ServicesController"] = append(beego.GlobalControllerRouter["github.com/udistrital/resoluciones_mid_v2/controllers:ServicesController"],
        beego.ControllerComments{
            Method: "DesagregadoPlaneacion",
            Router: "/desagregado_planeacion",
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

}
