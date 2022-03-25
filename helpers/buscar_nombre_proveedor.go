package helpers

import (
	"strconv"

	"github.com/udistrital/resoluciones_mid_v2/models"
)

func BuscarNombreProveedor(DocumentoIdentidad int) (nombre_prov string, outputError map[string]interface{}) {
	defer func() {
		if err := recover(); err != nil {
			outputError = map[string]interface{}{"funcion": "/BuscarNombreProveedor", "err": err, "status": "502"}
			panic(outputError)
		}
	}()
	var nom_proveedor string
	var informacion_proveedor []models.InformacionProveedor
	url := "informacion_proveedor?query=NumDocumento:" + strconv.Itoa(DocumentoIdentidad)
	if err2 := GetRequestNew("UrlcrudAgora", url, &informacion_proveedor); err2 == nil {
		if informacion_proveedor != nil {
			nom_proveedor = informacion_proveedor[0].NomProveedor
		} else {
			nom_proveedor = ""
		}
	} else {
		outputError = map[string]interface{}{"funcion": "/BuscarNombreProveedor2", "err": err2.Error(), "status": "404"}
		return nom_proveedor, outputError
	}

	return nom_proveedor, nil

}
