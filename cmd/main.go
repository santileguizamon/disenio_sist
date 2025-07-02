package main

import (
	"log"
	"net/http"
	"os"
	"time"

	_ "sistema-gestion-informacion/docs" // Documentación generada por swag

	"github.com/joho/godotenv"
	httpSwagger "github.com/swaggo/http-swagger"

	"sistema-gestion-informacion/internal/infrastructure/events"
	"sistema-gestion-informacion/internal/interfaces/handlers"
)

// @title Sistema de Gestión de Información API
// @version 1.0
// @description API RESTful con arquitectura dirigida por eventos para gestión de clientes potenciales
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /
// @schemes http
func main() {
	// Cargar variables de entorno
	if err := godotenv.Load(); err != nil {
		log.Println("No se encontró archivo .env, usando variables de entorno del sistema")
	}

	// Inicializar base de datos (Singleton)
	// dbInstance := database.GetInstance()
	// dsn := getDSN()
	// db, err := dbInstance.Connect(dsn)
	// if err != nil {
	// 	log.Fatalf("❌ Error conectando a la base de datos: %v", err)
	// }

	// Auto migrar modelos
	// log.Println("🔄 Migrando esquemas de base de datos...")
	// entities := []interface{}{
	// 	&entities.Producto{},
	// 	&entities.Venta{},
	// 	&entities.DetalleVenta{},
	// 	&entities.Sucursal{},
	// }

	// for _, entity := range entities {
	// 	if err := db.AutoMigrate(entity); err != nil {
	// 		log.Fatalf("❌ Error migrando entidad: %v", err)
	// 	}
	// }
	// log.Println("✅ Migración completada")

	// Inicializar bus de eventos (Singleton)
	eventBus := events.GetEventBusInstance()

	// Registrar manejadores de eventos
	registerEventHandlers(eventBus)

	// Crear handlers
	// clienteHandler := handlers.NewClienteHandler(db, eventBus)
	procesamientoHandler := handlers.NewProcesamientoHandler(eventBus)

	// Configurar rutas con HTTP nativo
	mux := http.NewServeMux()

	// Ruta de procesamiento (POST)
	mux.HandleFunc("/api/procesar", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			procesamientoHandler.ProcesarDatos(w, r)
		} else {
			http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		}
	})

	// Ruta para consultar datos depurados (GET)
	mux.HandleFunc("/api/datos-procesados", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			procesamientoHandler.GetDatosProcesados(w, r)
		} else {
			http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		}
	})

	// Ruta para consultar el último reporte generado (GET)
	mux.HandleFunc("/api/reporte", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			procesamientoHandler.GetUltimoReporte(w, r)
		} else {
			http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		}
	})

	// Ruta de salud
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"status":"ok","service":"Sistema de Gestión de Información"}`))
		} else {
			http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		}
	})

	// Ruta raíz con información del sistema
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet && r.URL.Path == "/" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{
				"name": "Sistema de Gestión de Información",
				"version": "1.0.0",
				"description": "API RESTful con arquitectura dirigida por eventos",
				"endpoints": {
					"clientes": "/api/clientes",
					"procesar": "/api/procesar",
					"health": "/health",
					"swagger": "/swagger/"
				}
			}`))
		} else {
			http.NotFound(w, r)
		}
	})

	// Swagger UI
	mux.HandleFunc("/swagger/", httpSwagger.Handler(
		httpSwagger.URL("http://localhost:8080/swagger/doc.json"),
	))

	// Configurar servidor HTTP
	port := getEnv("PORT", "8080")
	server := &http.Server{
		Addr:         ":" + port,
		Handler:      mux,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	log.Printf("🚀 Servidor iniciando en http://localhost:%s", port)
	log.Printf("📚 Documentación Swagger disponible en http://localhost:%s/swagger/", port)

	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("❌ Error iniciando servidor: %v", err)
	}
}

// getDSN obtiene la cadena de conexión a la base de datos
func getDSN() string {
	host := getEnv("DB_HOST", "localhost")
	port := getEnv("DB_PORT", "3306")
	user := getEnv("DB_USER", "root")
	password := getEnv("DB_PASSWORD", "")
	database := getEnv("DB_NAME", "sistema_gestion_informacion")

	return user + ":" + password + "@tcp(" + host + ":" + port + ")/" + database + "?charset=utf8mb4&parseTime=True&loc=Local"
}

// getEnv obtiene una variable de entorno con valor por defecto
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// registerEventHandlers registra los manejadores de eventos
func registerEventHandlers(eventBus *events.EventBus) {
	// Registrar logger de eventos
	eventLogger := &events.EventLogger{}
	eventBus.Subscribe("*", eventLogger)

	// Registrar métricas de eventos
	eventMetrics := events.NewEventMetrics()
	eventBus.Subscribe("*", eventMetrics)

	// Registrar manejadores específicos para el pipeline de procesamiento
	eventBus.Subscribe(events.EventDatosRecolectados, &DatosRecolectadosHandler{})
	eventBus.Subscribe(events.EventDatosProcesados, &DatosProcesadosHandler{})
	eventBus.Subscribe(events.EventDatosPersistidos, &DatosPersistidosHandler{})
	eventBus.Subscribe(events.EventReporteGenerado, &ReporteGeneradoHandler{})

	log.Println("✅ Manejadores de eventos registrados")
}

// Handlers específicos para eventos
type DatosRecolectadosHandler struct{}

func (h *DatosRecolectadosHandler) Handle(event events.Event) error {
	log.Printf("📥 Datos recolectados: %v", event.Data)
	return nil
}

func (h *DatosRecolectadosHandler) GetEventType() string {
	return events.EventDatosRecolectados
}

type DatosProcesadosHandler struct{}

func (h *DatosProcesadosHandler) Handle(event events.Event) error {
	log.Printf("⚙️ Datos procesados: %v", event.Data)
	return nil
}

func (h *DatosProcesadosHandler) GetEventType() string {
	return events.EventDatosProcesados
}

type DatosPersistidosHandler struct{}

func (h *DatosPersistidosHandler) Handle(event events.Event) error {
	log.Printf("💾 Datos persistidos: %v", event.Data)
	return nil
}

func (h *DatosPersistidosHandler) GetEventType() string {
	return events.EventDatosPersistidos
}

type ReporteGeneradoHandler struct{}

func (h *ReporteGeneradoHandler) Handle(event events.Event) error {
	log.Printf("📊 Reporte generado: %v", event.Data)
	return nil
}

func (h *ReporteGeneradoHandler) GetEventType() string {
	return events.EventReporteGenerado
}
