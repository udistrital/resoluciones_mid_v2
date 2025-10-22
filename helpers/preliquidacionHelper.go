package helpers

import (
	"fmt"
	"time"

	"github.com/astaxie/beego/logs"
	"github.com/udistrital/resoluciones_mid_v2/models"
)

func ProcesarPreliquidaciones(jobID string, rps []models.RpSeleccionado) {
	total := len(rps)
	completados := 0

	for _, rp := range rps {
		msg := fmt.Sprintf("Procesando vinculación %d (RP %d)", rp.VinculacionId, rp.Consecutivo)
		logs.Info(msg)
		ActualizarJob(jobID, msg, false, false)

		var vinculacion models.VinculacionDocente
		url := fmt.Sprintf("vinculacion_docente/%d", rp.VinculacionId)
		if err := GetRequestNew("UrlcrudResoluciones", url, &vinculacion); err != nil {
			errMsg := fmt.Sprintf("Error cargando vinculación %d: %v", rp.VinculacionId, err)
			ActualizarJob(jobID, errMsg, false, false)
			continue
		}

		res := EjecutarPreliquidacionTitan(vinculacion)

		completados++
		IncrementarProcesados(jobID)

		if res == nil {
			ActualizarJob(jobID,
				fmt.Sprintf("Vinculación %d (RP %d) -> Error desconocido", rp.VinculacionId, rp.Consecutivo),
				false, true)
			continue
		}

		status, _ := res["status"].(string)
		message, _ := res["message"].(string)

		switch status {
		case "ok":
			ActualizarJob(jobID, "Contrato liquidado correctamente: "+message, true, false)
		case "omitido":
			ActualizarJob(jobID, "Contrato omitido: "+message, false, true)
		default:
			ActualizarJob(jobID, "Error en contrato: "+message, false, true)
		}

		porcentaje := float64(completados) / float64(total) * 100
		progreso := fmt.Sprintf("Progreso: %d/%d (%.1f%%)", completados, total, porcentaje)
		ActualizarJob(jobID, progreso, false, false)

		time.Sleep(300 * time.Millisecond)
	}

	FinalizarJob(jobID)
	logs.Info("Job finalizado:", jobID)
}
