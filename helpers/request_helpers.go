package helpers

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/astaxie/beego"
	"github.com/udistrital/resoluciones_mid_v2/models"
	"github.com/udistrital/utils_oas/formatdata"
	utilsrequest "github.com/udistrital/utils_oas/request"
	"github.com/udistrital/utils_oas/xray"
)

// Envia una petición con datos al endpoint indicado y extrae la respuesta del campo Data para retornarla
func SendRequestNew(endpoint string, route string, trequest string, target interface{}, datajson interface{}) error {
	url := beego.AppConfig.String("ProtocolAdmin") + "://" + beego.AppConfig.String(endpoint) + route

	var response map[string]interface{}
	var err error
	err = SendJson(url, trequest, &response, &datajson)
	err = ExtractData(response, target, err)
	return err
}

func SendRequestFull(endpoint string, route string, trequest string, target interface{}, datajson interface{}) error {
	url := beego.AppConfig.String("ProtocolAdmin") + "://" + beego.AppConfig.String(endpoint) + route

	var b bytes.Buffer
	if datajson != nil {
		if err := json.NewEncoder(&b).Encode(datajson); err != nil {
			return fmt.Errorf("error codificando JSON: %v", err)
		}
	}

	client := &http.Client{Timeout: 60 * time.Second}
	req, err := http.NewRequest(trequest, url, &b)
	if err != nil {
		return fmt.Errorf("error creando request: %v", err)
	}
	req.Header.Set("Accept", AppJson)
	req.Header.Add("Content-Type", AppJson)

	seg := xray.BeginSegmentSec(req)
	resp, err := client.Do(req)
	xray.UpdateSegment(resp, err, seg)
	if err != nil {
		return fmt.Errorf("error ejecutando request: %v", err)
	}
	defer resp.Body.Close()

	bodyBytes := new(bytes.Buffer)
	if _, err := bodyBytes.ReadFrom(resp.Body); err != nil {
		return fmt.Errorf("error leyendo respuesta: %v", err)
	}

	if err := json.Unmarshal(bodyBytes.Bytes(), target); err != nil {
		return fmt.Errorf("respuesta no estándar o error de parseo: %v", err)
	}

	return nil
}

// Envia una petición con datos a endponts que responden con el body sin encapsular
func SendRequestLegacy(endpoint string, route string, trequest string, target interface{}, datajson interface{}) error {
	url := beego.AppConfig.String("ProtocolAdmin") + "://" + beego.AppConfig.String(endpoint) + route
	if err := SendJson(url, trequest, &target, &datajson); err != nil {
		return err
	}
	return nil
}

// Envia una petición al endpoint indicado y extrae la respuesta del campo Data para retornarla
func GetRequestNew(endpoint string, route string, target interface{}) error {
	url := beego.AppConfig.String("ProtocolAdmin") + "://" + beego.AppConfig.String(endpoint) + route
	var response map[string]interface{}
	var err error
	err = GetJson(url, &response)
	err = ExtractData(response, &target, err)
	return err
}

// Envia una petición a endponts que responden con el body sin encapsular
func GetRequestLegacy(endpoint string, route string, target interface{}) error {
	url := beego.AppConfig.String("ProtocolAdmin") + "://" + beego.AppConfig.String(endpoint) + route
	if err := GetJson(url, target); err != nil {
		return err
	}
	return nil
}

func GetRequestWSO2(service string, route string, target interface{}) error {
	url := beego.AppConfig.String("ProtocolAdmin") + "://" +
		beego.AppConfig.String("UrlcrudWSO2") +
		beego.AppConfig.String(service) + "/" + route
	if response, err := GetJsonWSO2Test(url, &target); response == 200 && err == nil {
		return nil
	} else {
		return err
	}
}

// Esta función extrae la información cuando se recibe encapsulada en una estructura
// y da manejo a las respuestas que contienen arreglos de objetos vacíos
func ExtractData(respuesta map[string]interface{}, v interface{}, err2 error) error {
	var err error

	if err2 != nil {
		return err2
	}
	if respuesta["Success"] == false {
		err = errors.New(fmt.Sprint(respuesta["Data"], respuesta["Message"]))
		panic(err)
	}
	datatype := fmt.Sprintf("%v", respuesta["Data"])
	switch datatype {
	case "map[]", "[map[]]":
	default:
		err = formatdata.FillStruct(respuesta["Data"], &v)
		respuesta = nil
	}
	return err
}

func SendJson(url string, trequest string, target interface{}, datajson interface{}) error {
	b := new(bytes.Buffer)
	if datajson != nil {
		if err := json.NewEncoder(b).Encode(datajson); err != nil {
			beego.Error(err)
		}
	}
	client := &http.Client{}
	req, err := http.NewRequest(trequest, url, b)
	req.Header.Set("Accept", AppJson)
	req.Header.Add("Content-Type", AppJson)
	seg := xray.BeginSegmentSec(req)
	r, err := client.Do(req)
	xray.UpdateSegment(r, err, seg)
	if err != nil {
		beego.Error("error", err)
		return err
	}
	defer func() {
		if err := r.Body.Close(); err != nil {
			beego.Error(err)
		}
	}()

	return json.NewDecoder(r.Body).Decode(target)
}

func GetJsonTest(url string, target interface{}) (status int, err error) {
	return utilsrequest.GetJsonTest2(url, target)
}

func GetJson(url string, target interface{}) error {
	return utilsrequest.GetJson(url, target)
}

func GetXml(url string, target interface{}) error {
	req, _ := http.NewRequest("GET", url, nil)
	seg := xray.BeginSegmentSec(req)
	client := &http.Client{}
	r, err := client.Do(req)
	xray.UpdateSegment(r, err, seg)
	if err != nil {
		return err
	}
	defer func() {
		if err := r.Body.Close(); err != nil {
			beego.Error(err)
		}
	}()

	return xml.NewDecoder(r.Body).Decode(target)
}

func GetJsonWSO2(urlp string, target interface{}) error {
	return utilsrequest.GetJsonWSO2(urlp, target)
}

func GetJsonWSO2Test(urlp string, target interface{}) (status int, err error) {
	return utilsrequest.GetJsonWSO2Test(urlp, target)
}

func GetTipoResolucion(id int) (tipoResolucion models.Parametro) {
	var tipoResolucionAux models.Parametro
	var resolucionAux models.Resolucion
	err2 := GetRequestNew("UrlCrudResoluciones", "resolucion/"+strconv.Itoa(id), &resolucionAux)
	if err2 != nil {
		panic(err2.Error())
	}
	err3 := GetRequestNew("UrlcrudParametros", "parametro/"+strconv.Itoa(resolucionAux.TipoResolucionId), &tipoResolucionAux)
	if err3 != nil {
		panic(err3.Error())
	}

	return tipoResolucionAux
}

func GetResolucion(id int) (resolucion models.Resolucion) {
	urlAux := "resolucion/" + strconv.Itoa(id)
	if err := GetRequestNew("UrlcrudResoluciones", urlAux, &resolucion); err != nil {
		panic("Consultando resolución -> " + err.Error())
	}
	return resolucion
}
