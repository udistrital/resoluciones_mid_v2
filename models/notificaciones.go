package models

type TemplatedEmail struct {
	Source              string         `json:"Source"`
	Template            string         `json:"Template"`
	Destinations        []Destinations `json:"Destinations"`
	DefaultTemplateData TemplateData   `json:"DefaultTemplateData"`
}

type Destination struct {
	BccAddresses []string `json:"BccAddresses"`
	CcAddresses  []string `json:"CcAddresses"`
	ToAddresses  []string `json:"ToAddresses"`
}
type Destinations struct {
	Destination             Destination   `json:"Destination"`
	ReplacementTemplateData TemplateData  `json:"ReplacementTemplateData"`
	Attachments             []Attachments `json:"Attachments"`
}

type TemplateData struct {
	NumeroResolucion string `json:"numero_resolucion"`
	Facultad         string `json:"facultad"`
	NumeroContrato   string `json:"numero_contrato"`
}

type Content struct {
	Data string `json:"Data"`
}

type Body struct {
	Html Content `json:"Html"`
	Text Content `json:"Text"`
}

type Attachments struct {
	ContentType string `json:"ContentType"`
	FileName    string `json:"FileName"`
	Base64File  string `json:"Base64File"`
}
type EmailResponse struct {
	Result struct {
		MessageId string `json:"MessageId"`
	} `json:"Result"`
}

type EmailData struct {
	Documento        string
	ContratoId       string
	Facultad         string
	NumeroResolucion string
	Correo           string
}
