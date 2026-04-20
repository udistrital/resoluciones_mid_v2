package helpers

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"math/big"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/udistrital/utils_oas/formatdata"
)

const (
	ErrorParametros     string = "Error en los parametros de ingreso"
	ErrorBody           string = "Cuerpo de la peticion invalido"
	CargaResExito       string = "Resoluciones cargadas con exito"
	CampoMeses          string = "%d semanas"
	PasaA               string = "Pasa a %d"
	ResolucionEndpoint  string = "resolucion/"
	ParametroEndpoint   string = "parametro/"
	VinculacionEndpoint string = "vinculacion_docente/"
	ResVinEndpoint      string = "resolucion_vinculacion_docente/"
	AppJson             string = "application/json"
	Calibri             string = "Calibri"
	CalibriBold         string = "Calibri-Bold"
	MinionProBoldCn     string = "MinionPro-BoldCn"
	MinionProMediumCn   string = "MinionPro-MediumCn"
	MinionProBoldItalic string = "MinionProBoldItalic"
)

func JsonDebug(i interface{}) {
	formatdata.JsonPrint(i)
}

func iguales(a interface{}, b interface{}) bool {
	return reflect.DeepEqual(a, b)
}

func diff(a, b time.Time) (year, month, day int) {
	if a.Location() != b.Location() {
		b = b.In(a.Location())
	}
	if a.After(b) {
		a, b = b, a
	}
	oneDay := time.Hour * 5
	a = a.Add(oneDay)
	b = b.Add(oneDay)
	y1, M1, d1 := a.Date()
	y2, M2, d2 := b.Date()

	year = int(y2 - y1)
	month = int(M2 - M1)
	day = int(d2 - d1)

	// Normalize negative values

	if day < 0 {
		// days in month:
		t := time.Date(y1, M1, 32, 0, 0, 0, 0, time.UTC)
		day += 32 - t.Day()
		month--
	}
	if month < 0 {
		month += 12
		year--
	}

	return
}

func FormatMoney(value interface{}, Precision int) string {
	formattedNumber := FormatNumber(value, Precision, ",", ".")
	return FormatMoneyString(formattedNumber, Precision)
}

func FormatMoneyString(formattedNumber string, Precision int) string {
	var format string

	zero := "0"
	if Precision > 0 {
		zero += "." + strings.Repeat("0", Precision)
	}

	format = "%s%v"
	result := strings.Replace(format, "%s", "$", -1)
	result = strings.Replace(result, "%v", formattedNumber, -1)

	return result
}

func FormatNumber(value interface{}, precision int, thousand string, decimal string) string {
	v := reflect.ValueOf(value)
	var x string
	switch v.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		x = fmt.Sprintf("%d", v.Int())
		if precision > 0 {
			x += "." + strings.Repeat("0", precision)
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		x = fmt.Sprintf("%d", v.Uint())
		if precision > 0 {
			x += "." + strings.Repeat("0", precision)
		}
	case reflect.Float32, reflect.Float64:
		x = fmt.Sprintf(fmt.Sprintf("%%.%df", precision), v.Float())
	case reflect.Ptr:
		switch v.Type().String() {
		case "*big.Rat":
			x = value.(*big.Rat).FloatString(precision)

		default:
			panic("Unsupported type - " + v.Type().String())
		}
	default:
		panic("Unsupported type - " + v.Kind().String())
	}

	return formatNumberString(x, precision, thousand, decimal)
}

func formatNumberString(x string, precision int, thousand string, decimal string) string {
	lastIndex := strings.Index(x, ".") - 1
	if lastIndex < 0 {
		lastIndex = len(x) - 1
	}

	var buffer []byte
	var strBuffer bytes.Buffer

	j := 0
	for i := lastIndex; i >= 0; i-- {
		j++
		buffer = append(buffer, x[i])

		if j == 3 && i > 0 && !(i == 1 && x[0] == '-') {
			buffer = append(buffer, ',')
			j = 0
		}
	}

	for i := len(buffer) - 1; i >= 0; i-- {
		strBuffer.WriteByte(buffer[i])
	}
	result := strBuffer.String()

	if thousand != "," {
		result = strings.Replace(result, ",", thousand, -1)
	}

	extra := x[lastIndex+1:]
	if decimal != "." {
		extra = strings.Replace(extra, ".", decimal, 1)
	}

	return result + extra
}

// Valida que el body recibido en la petición tenga contenido válido
func ValidarBody(body []byte) (valid bool, err error) {
	var test interface{}
	if err = json.Unmarshal(body, &test); err != nil {
		return false, err
	} else {
		content := fmt.Sprintf("%v", test)
		switch content {
		case "map[]", "[map[]]": // body vacio
			return false, nil
		}
	}
	return true, nil
}

// Quita el formato de moneda a un string y lo convierte en valor flotante
func DeformatNumber(formatted string) (number float64) {
	formatted = strings.ReplaceAll(formatted, ",", "")
	formatted = strings.Trim(formatted, "$")
	number, _ = strconv.ParseFloat(formatted, 64)
	return
}

// Obtiene los datos del usuario autenticado
func GetUsuario(usuario string) (nombreUsuario map[string]interface{}, err error) {
	if len(usuario) > 0 {
		var decData map[string]interface{}
		if data, err6 := base64.StdEncoding.DecodeString(usuario); err6 != nil {
			return nombreUsuario, err6
		} else {
			if err7 := json.Unmarshal(data, &decData); err7 != nil {
				return nombreUsuario, err7
			}
		}
		nombreUsuario = decData["user"].(map[string]interface{})
	}
	return nombreUsuario, err
}
