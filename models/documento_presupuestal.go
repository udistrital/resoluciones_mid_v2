package models

type DocumentoPresupuestal struct {
	Tipo          string
	AfectacionIds []string
	Afectacion    []MovimientoRubro
	FechaRegistro string
	Estado        string
	ValorActual   float64
	ValorInicial  float64
	Vigencia      int
	Consecutivo   float64
}

type MovimientoRubro struct {
	IDPsql        int
	Tipo          string
	Padre         string
	FechaRegistro string
	Estado        string
	ValorActual   float64
	ValorInicial  float64
}
