package models

type ContratoTitan struct {
	Id             int    `json:"Id"`
	NumeroContrato string `json:"NumeroContrato"`
	Vigencia       int    `json:"Vigencia"`
	Documento      string `json:"Documento"`
	ResolucionId   int    `json:"ResolucionId"`
	Rp             int    `json:"Rp"`
	Activo         bool   `json:"Activo"`
}
