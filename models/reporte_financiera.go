package models

type ReporteFinanciera struct {
	Id                    int
	Resolucion            string
	Cedula                int
	Horas                 int
	Semanas               int
	Total                 float64
	Cdp                   int
	SueldoBasico          float64
	PrimaNavidad          float64
	Vacaciones            float64
	PrimaVacaciones       float64
	Cesantias             float64
	InteresesCesantias    float64
	PrimaServicios        float64
	BonificacionServicios float64
	ProyectoCurricular    int
}

type ReporteFinancieraFinal struct {
	Id                    int
	Resolucion            string
	Cedula                int
	Horas                 int
	Semanas               int
	Total                 float64
	Cdp                   int
	SueldoBasico          float64
	PrimaNavidad          float64
	Vacaciones            float64
	PrimaVacaciones       float64
	Cesantias             float64
	InteresesCesantias    float64
	PrimaServicios        float64
	BonificacionServicios float64
	Nombre                string
	ProyectoCurricular    string
	CodigoProyecto        int
	Facultad              string
}

type ObjetoDocenteTg struct {
	DocenteTg struct {
		Docente []Docente `json:"docente"`
	} `json:"docentes"`
}

type Docente struct {
	Id     string `json:"id"`
	Nombre string `json:"NOMBRE"`
}
