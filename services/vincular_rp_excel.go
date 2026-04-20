package services

import (
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/astaxie/beego/logs"
	"github.com/udistrital/resoluciones_mid_v2/models"
	"github.com/xuri/excelize/v2"
)

func cargarRegistrosRpDesdeArchivo(file multipart.File, fileHeader *multipart.FileHeader) ([]models.VinculacionRpResultado, error) {
	if fileHeader == nil {
		return nil, errors.New("no se recibió archivo en la solicitud")
	}

	tmpDir := os.TempDir()
	tmpPath := filepath.Join(tmpDir, fmt.Sprintf("rp_%s.xlsx", time.Now().Format("02012006_150405")))
	out, err := os.Create(tmpPath)
	if err != nil {
		return nil, fmt.Errorf("no se pudo crear archivo temporal: %v", err)
	}
	defer out.Close()

	_, _ = file.Seek(0, 0)
	_, _ = io.Copy(out, file)

	f, err := excelize.OpenFile(tmpPath)
	if err != nil {
		logs.Error("Error al abrir el archivo Excel: %v", err)
		return nil, fmt.Errorf("no se pudo leer el archivo Excel: %v", err)
	}
	defer f.Close()

	sheetName := f.GetSheetName(0)
	if sheetName == "" {
		return nil, errors.New("el archivo no contiene hojas válidas")
	}

	rows, err := f.GetRows(sheetName)
	if err != nil {
		return nil, fmt.Errorf("error leyendo filas: %v", err)
	}
	if len(rows) < 2 {
		return nil, errors.New("el archivo no contiene datos suficientes")
	}

	headers, err := extraerHeadersRp(rows[0])
	if err != nil {
		return nil, err
	}

	return construirRegistrosRp(rows[1:], headers), nil
}

func construirEstadoConflictoRp(info *conflictoInfo) string {
	crps := make([]string, 0, len(info.CRPs))
	for crp := range info.CRPs {
		crps = append(crps, crp)
	}
	sort.Strings(crps)

	filas := append([]int{}, info.Filas...)
	sort.Ints(filas)

	return fmt.Sprintf(
		"CONFLICTO: llave duplicada con CRPs diferentes. CRPs=%s. Filas=%v. No se actualiza ninguna fila de esta llave.",
		strings.Join(crps, ","),
		filas,
	)
}
