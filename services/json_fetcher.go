package services

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/udistrital/utils_oas/request"
)

func normalizeBaseNoProto(u string) string {
	u = strings.TrimSpace(u)
	u = strings.TrimLeft(u, "/")
	return u
}

func joinWSO2URL(protocol, base, ns, path string) string {
	protocol = strings.TrimRight(protocol, "://")
	base = strings.TrimRight(normalizeBaseNoProto(base), "/")
	ns = strings.Trim(ns, "/")
	path = strings.TrimLeft(path, "/")
	return fmt.Sprintf("%s://%s/%s/%s", protocol, base, ns, path)
}

func getJSONWithUtilOAS(url string, target interface{}) map[string]interface{} {
	if err := request.GetJson(url, target); err != nil {
		return map[string]interface{}{
			"funcion": "getJSONWithUtilOAS",
			"err":     err.Error(),
			"status":  "502",
		}
	}
	return nil
}

func getJSONWithHTTP(url string, target interface{}) map[string]interface{} {
	client := &http.Client{
		Timeout: 20 * time.Second,
	}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return map[string]interface{}{
			"funcion": "getJSONWithHTTP:newRequest",
			"err":     err.Error(),
			"status":  "500",
		}
	}

	req.Header.Set("Accept", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return map[string]interface{}{
			"funcion": "getJSONWithHTTP:do",
			"err":     err.Error(),
			"status":  "502",
		}
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return map[string]interface{}{
			"funcion": "getJSONWithHTTP:readBody",
			"err":     err.Error(),
			"status":  "502",
		}
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return map[string]interface{}{
			"funcion": "getJSONWithHTTP:statusCode",
			"err":     fmt.Sprintf("respuesta no exitosa del servicio externo: %d - %s", resp.StatusCode, strings.TrimSpace(string(body))),
			"status":  strconv.Itoa(resp.StatusCode),
		}
	}

	if err := json.Unmarshal(body, target); err != nil {
		return map[string]interface{}{
			"funcion": "getJSONWithHTTP:unmarshal",
			"err":     err.Error(),
			"status":  "502",
		}
	}

	return nil
}

func getJSON(url string, target interface{}) map[string]interface{} {
	if errMap := getJSONWithUtilOAS(url, target); errMap == nil {
		return nil
	}
	return getJSONWithHTTP(url, target)
}
