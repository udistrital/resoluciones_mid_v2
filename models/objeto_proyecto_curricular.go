package models

type ObjetoProyectoCurricular struct {
	Homologacion struct {
		CodigoProyecto string `json:"codigo_proyecto"`
		IdArgo         string `json:"id_argo"`
		ProyectoSnies  string `json:"proyecto_snies"`
		IDOikos        string `json:"id_oikos"`
		IDSnies        string `json:"id_snies"`
	} `json:"homologacion"`
}
