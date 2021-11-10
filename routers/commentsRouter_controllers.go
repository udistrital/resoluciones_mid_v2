package routers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context/param"
)

func init() {

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
            Method: "Put",
            Router: "/:id",
            AllowHTTPMethods: []string{"put"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

}
