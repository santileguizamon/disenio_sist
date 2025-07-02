package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"sistema-gestion-informacion/internal/infrastructure/builders"
	"sistema-gestion-informacion/internal/infrastructure/events"
)

// ProcesamientoHandler maneja las peticiones de procesamiento de datos
type ProcesamientoHandler struct {
	eventBus *events.EventBus
}

// NewProcesamientoHandler crea una nueva instancia del handler
func NewProcesamientoHandler(eventBus *events.EventBus) *ProcesamientoHandler {
	return &ProcesamientoHandler{
		eventBus: eventBus,
	}
}

// ProcesarDatos godoc
// @Summary Procesar datos
// @Description Recibe datos crudos, los procesa y depura usando el servicio y el EventBus
// @Tags procesamiento
// @Accept json
// @Produce json
// @Param datos body map[string]interface{} true "Datos crudos a procesar"
// @Success 200 {object} map[string]string
// @Router /api/procesar [post]
func (h *ProcesamientoHandler) ProcesarDatos(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	var datosCrudos map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&datosCrudos); err != nil {
		http.Error(w, "JSON inválido", http.StatusBadRequest)
		return
	}

	// Simular procesamiento real usando el servicio (aquí se puede integrar ProcesadorDatosService)
	// Por simplicidad, guardamos los datos depurados en memoria
	procesados := procesarYDepurar(datosCrudos)
	storeDatosProcesados(procesados)

	// --- Generar reporte usando el builder ---
	reporteBuilder := builders.NewReporteBuilder().SetTipo("reporte_procesamiento").SetFormato("excel")
	reporteBuilder.AddDato("datos_procesados", procesados)
	reporte, _ := reporteBuilder.Build()
	storeUltimoReporte(reporte)
	// --- Fin generación de reporte ---

	h.eventBus.Publish(events.CreateEvent(
		events.EventDatosProcesados,
		map[string]interface{}{
			"registros_final": len(procesados),
			"timestamp":       time.Now(),
		},
		"procesamiento_handler",
	))

	response := map[string]string{
		"status":  "Proceso completado exitosamente",
		"message": "Datos procesados y depurados y reporte generado",
		"time":    time.Now().Format(time.RFC3339),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// GET /api/datos-procesados
func (h *ProcesamientoHandler) GetDatosProcesados(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(getDatosProcesados())
}

// GET /api/reporte
func (h *ProcesamientoHandler) GetUltimoReporte(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(getUltimoReporte())
}

// --- Simulación de almacenamiento en memoria ---
var (
	datosProcesadosMem []map[string]interface{}
	ultimoReporteMem   *builders.Reporte
)

func storeDatosProcesados(datos []map[string]interface{}) {
	datosProcesadosMem = datos
}

func getDatosProcesados() []map[string]interface{} {
	return datosProcesadosMem
}

func storeUltimoReporte(reporte *builders.Reporte) {
	ultimoReporteMem = reporte
}

func getUltimoReporte() *builders.Reporte {
	return ultimoReporteMem
}

// procesarYDepurar simula el procesamiento y depuración de datos crudos
func procesarYDepurar(datos map[string]interface{}) []map[string]interface{} {
	// Aquí puedes integrar la lógica real de ProcesadorDatosService
	// Por ahora, simplemente envolvemos los datos en un slice
	return []map[string]interface{}{datos}
}
