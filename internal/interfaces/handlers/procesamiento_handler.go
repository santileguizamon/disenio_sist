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

// Estructuras para documentación Swagger
type DatosProcesamientoRequest struct {
	Datos []map[string]interface{} `json:"datos" example:"[{\"producto_id\":1,\"nombre\":\"Producto A\",\"precio\":100.50,\"sucursal_id\":1,\"cantidad\":5}]"`
}

type ProcesamientoResponse struct {
	Status  string `json:"status" example:"Proceso completado exitosamente"`
	Message string `json:"message" example:"Datos procesados y depurados y reporte generado"`
	Time    string `json:"time" example:"2024-01-15T10:30:00Z"`
}

type DatosProcesadosResponse struct {
	Productos  []map[string]interface{} `json:"productos"`
	Ventas     []map[string]interface{} `json:"ventas"`
	Sucursales []map[string]interface{} `json:"sucursales"`
}

type ReporteResponse struct {
	ID          uint                   `json:"ID" example:"1"`
	Tipo        string                 `json:"Tipo" example:"ventas"`
	GeneradoEn  time.Time              `json:"GeneradoEn" example:"2024-01-15T12:00:00Z"`
	Estado      string                 `json:"Estado" example:"completado"`
	FechaInicio time.Time              `json:"FechaInicio" example:"2024-01-01T00:00:00Z"`
	FechaFin    time.Time              `json:"FechaFin" example:"2024-01-15T23:59:59Z"`
	Datos       map[string]interface{} `json:"Datos"`
	Filtros     []string               `json:"Filtros" example:"[\"fecha_inicio\",\"fecha_fin\"]"`
	Resumen     string                 `json:"Resumen" example:"Reporte de ventas del período"`
}

type ErrorResponse struct {
	Error   string `json:"error" example:"Bad Request"`
	Message string `json:"message" example:"JSON inválido"`
	Time    string `json:"time" example:"2024-01-15T10:30:00Z"`
}

// ProcesarDatos godoc
// @Summary Procesar datos crudos
// @Description Recibe datos crudos, los procesa y depura, y dispara eventos de notificación
// @Tags procesamiento
// @Accept json
// @Produce json
// @Param request body DatosProcesamientoRequest true "Datos crudos a procesar"
// @Success 200 {object} ProcesamientoResponse
// @Failure 400 {object} ErrorResponse
// @Failure 405 {object} ErrorResponse
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

	response := ProcesamientoResponse{
		Status:  "Proceso completado exitosamente",
		Message: "Datos procesados y depurados y reporte generado",
		Time:    time.Now().Format(time.RFC3339),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// GetDatosProcesados godoc
// @Summary Consultar datos procesados
// @Description Obtiene los datos procesados y depurados almacenados en memoria
// @Tags procesamiento
// @Produce json
// @Success 200 {object} DatosProcesadosResponse
// @Failure 405 {object} ErrorResponse
// @Router /api/datos-procesados [get]
func (h *ProcesamientoHandler) GetDatosProcesados(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	// Simular datos de respuesta estructurados
	response := DatosProcesadosResponse{
		Productos: []map[string]interface{}{
			{
				"ID":           1,
				"SKU":          "PROD-001",
				"Nombre":       "Producto A",
				"Descripcion":  "Descripción del producto A",
				"Categoria":    "general",
				"Fabricante":   "Fabricante XYZ",
				"Precio":       100.50,
				"PrecioOferta": 90.45,
				"StockMinimo":  10,
				"StockActual":  45,
			},
		},
		Ventas: []map[string]interface{}{
			{
				"ID":         1,
				"SucursalID": 1,
				"FechaVenta": "2024-01-15T10:30:00Z",
				"MetodoPago": "efectivo",
				"Total":      502.50,
				"DetallesVenta": []map[string]interface{}{
					{
						"ProductoID":     1,
						"Cantidad":       5,
						"PrecioUnitario": 100.50,
						"Descuento":      0.00,
					},
				},
			},
		},
		Sucursales: []map[string]interface{}{
			{
				"ID":          1,
				"Nombre":      "Sucursal Centro",
				"Direccion":   "Av. Principal 123",
				"Telefono":    "+1234567890",
				"Email":       "centro@empresa.com",
				"Ciudad":      "Ciudad Principal",
				"TipoSistema": "api",
				"ApiEndpoint": "https://api.sucursal.com",
				"ApiKey":      "key123",
				"ApiSecret":   "secret123",
			},
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// GetUltimoReporte godoc
// @Summary Obtener último reporte
// @Description Obtiene el último reporte generado usando el patrón Builder
// @Tags procesamiento
// @Produce json
// @Success 200 {object} ReporteResponse
// @Failure 405 {object} ErrorResponse
// @Router /api/reporte [get]
func (h *ProcesamientoHandler) GetUltimoReporte(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	// Simular reporte de respuesta
	response := ReporteResponse{
		ID:          1,
		Tipo:        "ventas",
		GeneradoEn:  time.Now(),
		Estado:      "completado",
		FechaInicio: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
		FechaFin:    time.Date(2024, 1, 15, 23, 59, 59, 0, time.UTC),
		Datos: map[string]interface{}{
			"total_ventas":           502.50,
			"cantidad_transacciones": 1,
			"producto_mas_vendido":   "Producto A",
			"sucursal_principal":     "Sucursal Centro",
		},
		Filtros: []string{"fecha_inicio", "fecha_fin"},
		Resumen: "Reporte de ventas del período 2024-01-01 al 2024-01-15",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
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
