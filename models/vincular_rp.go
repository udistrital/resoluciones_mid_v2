package models

type VinculacionRpResultado struct {
	CodResolucion string `json:"cod_resolucion"`
	CodFacultad   string `json:"cod_facultad"`
	Documento     string `json:"documento"`
	CodProyecto   string `json:"cod_proyecto"`
	CRP           string `json:"crp"`
	IdResolucion  string `json:"id_resolucion"`
	IdVinculacion string `json:"id_vinculacion"`
	PutStatus     string `json:"put_status"`
	FilaExcel     int    `json:"fila_excel"`
}
