package helpers

import (
	"errors"
	"net/http"
	"time"

	"github.com/udistrital/utils_oas/request"
)

func SafeRequest(method, url string, data interface{}, result interface{}) error {
	var err error
	for attempt := 1; attempt <= 3; attempt++ {
		switch method {
		case http.MethodGet:
			err = request.GetJson(url, &result)
		case http.MethodPut:
			err = request.SendJson(url, method, nil, data)
		default:
			return errors.New("método HTTP no soportado")
		}

		if err == nil {
			return nil
		}

		time.Sleep(3 * time.Second)
	}

	return errors.New("no se pudo conectar después de 3 intentos")
}

func SanitizePayload(data map[string]interface{}) map[string]interface{} {
	clean := make(map[string]interface{})
	for k, v := range data {
		if v != nil {
			clean[k] = v
		}
	}
	return clean
}
