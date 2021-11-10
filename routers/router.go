// @APIVersion 1.0.0
// @Title Resoluciones MID API Versión 2
// @Description API MID para el sistema de Resoluciones en su nueva versión
// @Contact computo@udistrital.edu.co
// @TermsOfServiceUrl http://beego.me/
// @License Apache 2.0
// @LicenseUrl http://www.apache.org/licenses/LICENSE-2.0.html
package routers

import (
	"github.com/astaxie/beego"
	"github.com/udistrital/resoluciones_mid_v2/controllers"
)

func init() {
	ns := beego.NewNamespace("/v1",
		beego.NSNamespace("/gestion_resoluciones",
			beego.NSInclude(
				&controllers.GestionResolucionesController{},
			),
		),
	)
	beego.AddNamespace(ns)
}
