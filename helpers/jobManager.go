package helpers

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

type Job struct {
	ID         string
	Total      int
	Procesados int
	Exitosos   int
	Omitidos   int
	Errores    int
	Estado     string
	Mensajes   []string
	Inicio     time.Time
	Fin        *time.Time
}

var jobs = make(map[string]*Job)
var mu sync.Mutex

func CrearJob(total int) string {
	mu.Lock()
	defer mu.Unlock()

	jobID := RandString(10)
	jobs[jobID] = &Job{
		ID:         jobID,
		Total:      total,
		Procesados: 0,
		Exitosos:   0,
		Omitidos:   0,
		Errores:    0,
		Estado:     "En progreso",
		Mensajes:   []string{},
		Inicio:     time.Now(),
	}
	return jobID
}

func ActualizarJob(jobID string, mensaje string, exito bool, omitido bool) {
	mu.Lock()
	defer mu.Unlock()

	job, ok := jobs[jobID]
	if !ok {
		return
	}

	job.Mensajes = append(job.Mensajes, mensaje)

	if exito {
		job.Exitosos++
	} else if omitido {
		job.Omitidos++
	} else {
		job.Errores++
	}

	if job.Procesados == job.Total {
		job.Estado = "Completado"
		now := time.Now()
		job.Fin = &now
	}
}

func IncrementarProcesados(jobID string) {
	mu.Lock()
	defer mu.Unlock()

	if job, ok := jobs[jobID]; ok {
		job.Procesados++
	}
}

func ObtenerJob(jobID string) map[string]interface{} {
	mu.Lock()
	defer mu.Unlock()

	job, ok := jobs[jobID]
	if !ok {
		return map[string]interface{}{
			"Success": false,
			"Message": fmt.Sprintf("No se encontró el job con ID %s", jobID),
		}
	}

	porcentaje := 0.0
	if job.Total > 0 {
		porcentaje = float64(job.Procesados) / float64(job.Total) * 100
	}

	return map[string]interface{}{
		"Success":          true,
		"JobId":            job.ID,
		"Estado":           job.Estado,
		"Inicio":           job.Inicio,
		"Fin":              job.Fin,
		"Total":            job.Total,
		"Procesados":       job.Procesados,
		"Exitosos":         job.Exitosos,
		"Omitidos":         job.Omitidos,
		"Errores":          job.Errores,
		"Mensajes":         job.Mensajes,
		"PorcentajeAvance": porcentaje,
	}
}

func RandString(n int) string {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	rand.Seed(time.Now().UnixNano())
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func FinalizarJob(jobID string) {
	mu.Lock()
	defer mu.Unlock()

	job, ok := jobs[jobID]
	if !ok {
		return
	}

	job.Estado = "Completado"
	now := time.Now()
	job.Fin = &now
	job.Mensajes = append(job.Mensajes, "Proceso finalizado correctamente.")
}
