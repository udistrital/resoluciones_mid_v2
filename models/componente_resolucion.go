package models

type ComponenteResolucion struct {
	Id                        int
	Numero                    int
	ResolucionId              *Resolucion
	Texto                     string
	TipoComponente            string
	ComponenteResolucionPadre *ComponenteResolucion
	Activo                    bool
	FechaCreacion             string
	FechaModificacion         string
}
