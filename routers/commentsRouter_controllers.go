package routers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context/param"
)

func init() {

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
            Method: "ConsultaDocente",
            Router: "/consultar_docente/:id",
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

    beego.GlobalControllerRouter["github.com/udistrital/resoluciones_mid_v2/controllers:ServicesController"] = append(beego.GlobalControllerRouter["github.com/udistrital/resoluciones_mid_v2/controllers:ServicesController"],
        beego.ControllerComments{
            Method: "DesagregadoPlaneacion",
            Router: "/desagregado_planeacion",
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

}
