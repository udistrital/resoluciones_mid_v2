package main

import (
	"flag"
	"net/http"
	"os"
	"strings"
	"testing"
)

var parameters struct {
	ProtocolAdmin       string
	UrlcrudResoluciones string
	UrlcrudAgora        string
	UrlcrudKronos       string
	Urlruler            string
	UrlcrudOikos        string
	UrlcrudParametros   string
	UrlmidTerceros      string
	UrlcrudWSO2         string
}

func TestMain(m *testing.M) {
	parameters.ProtocolAdmin = os.Getenv("RESOLUCIONES_MID_V2_PROTOCOL_ADMIN")
	parameters.UrlcrudResoluciones = os.Getenv("RESOLUCIONES_MID_V2_RESOLUCIONES_CRUD_URL")
	parameters.UrlcrudAgora = os.Getenv("RESOLUCIONES_MID_V2_AGORA_URL")
	parameters.UrlcrudKronos = os.Getenv("RESOLUCIONES_MID_V2_KRONOS_URL")
	parameters.Urlruler = os.Getenv("RESOLUCIONES_MID_V2_RULER_URL")
	parameters.UrlcrudOikos = os.Getenv("RESOLUCIONES_MID_V2_OIKOS_URL")
	parameters.UrlcrudParametros = os.Getenv("RESOLUCIONES_MID_V2_PARAMETROS_URL")
	parameters.UrlmidTerceros = os.Getenv("RESOLUCIONES_MID_V2_TERCEROS_MID_URL")
	parameters.UrlcrudWSO2 = os.Getenv("RESOLUCIONES_MID_V2_WSO2_URL")
	flag.Parse()
	os.Exit(m.Run())
}

func TestEndPointResolucionesCrud(t *testing.T) {
	endpoint := parameters.ProtocolAdmin + "://" + strings.Replace(parameters.UrlcrudResoluciones, "/v1/", "", 1)
	BaseTestEndpoint(t, endpoint)
}
func TestEndPointTerceros(t *testing.T) {
	endpoint := parameters.ProtocolAdmin + "://" + strings.Replace(parameters.UrlmidTerceros, "/v1/", "", 1)
	BaseTestEndpoint(t, endpoint)
}
func TestEndPointAgora(t *testing.T) {
	endpoint := parameters.ProtocolAdmin + "://" + strings.Replace(parameters.UrlcrudAgora, "/v1/", "", 1)
	BaseTestEndpoint(t, endpoint)
}
func TestEndPointKronos(t *testing.T) {
	endpoint := parameters.ProtocolAdmin + "://" + strings.Replace(parameters.UrlcrudKronos, "/v1/", "", 1)
	BaseTestEndpoint(t, endpoint)
}
func TestEndPointRuler(t *testing.T) {
	endpoint := parameters.ProtocolAdmin + "://" + strings.Replace(parameters.Urlruler, "/v1/", "", 1)
	BaseTestEndpoint(t, endpoint)
}
func TestEndPointOikos(t *testing.T) {
	endpoint := parameters.ProtocolAdmin + "://" + strings.Replace(parameters.UrlcrudOikos, "/v2/", "", 1)
	BaseTestEndpoint(t, endpoint)
}
func TestEndPointWSO2(t *testing.T) {
	endpoint := parameters.ProtocolAdmin + "://" + parameters.UrlcrudWSO2
	BaseTestEndpoint(t, endpoint)
}

func BaseTestEndpoint(t *testing.T, endpoint string) {
	t.Log(endpoint)
	if response, err := http.Get(endpoint); err == nil {
		if response.StatusCode != 200 {
			t.Error("Error TestEndpoint:", endpoint, "Estado: ", response.StatusCode)
			t.Fail()
		} else {
			t.Log("TestEndPoint", endpoint, "Finalizado Correctamente (OK)")
		}
	} else {
		t.Error("Error EndPoint:", err.Error())
		t.Fail()
	}
}
