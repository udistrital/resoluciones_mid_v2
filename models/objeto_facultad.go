package models

type ObjetoFacultad struct {
	Homologacion struct {
		IdOikos string `json:"id_oikos"`
		IdGeDep string `json:"id_gedep"`
	} `json:"homologacion"`
}
