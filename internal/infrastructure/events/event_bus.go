package events

import (
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"
)

// EventBus implementa el patrón Observer para manejo de eventos
type EventBus struct {
	subscribers map[string][]EventHandler
	mutex       sync.RWMutex
}

// Event representa un evento en el sistema
type Event struct {
	ID        string                 `json:"id"`
	Type      string                 `json:"type"`
	Data      map[string]interface{} `json:"data"`
	Timestamp time.Time              `json:"timestamp"`
	Source    string                 `json:"source"`
	Priority  int                    `json:"priority"` // 1: baja, 2: media, 3: alta
}

// EventHandler define la interfaz para los manejadores de eventos
type EventHandler interface {
	Handle(event Event) error
	GetEventType() string
}

// NewEventBus crea una nueva instancia del bus de eventos
func NewEventBus() *EventBus {
	return &EventBus{
		subscribers: make(map[string][]EventHandler),
	}
}

// Subscribe registra un manejador de eventos para un tipo específico
func (eb *EventBus) Subscribe(eventType string, handler EventHandler) {
	eb.mutex.Lock()
	defer eb.mutex.Unlock()

	eb.subscribers[eventType] = append(eb.subscribers[eventType], handler)
	log.Printf("Manejador registrado para evento: %s", eventType)
}

// Unsubscribe remueve un manejador de eventos
func (eb *EventBus) Unsubscribe(eventType string, handler EventHandler) {
	eb.mutex.Lock()
	defer eb.mutex.Unlock()

	if handlers, exists := eb.subscribers[eventType]; exists {
		for i, h := range handlers {
			if h == handler {
				eb.subscribers[eventType] = append(handlers[:i], handlers[i+1:]...)
				break
			}
		}
	}
}

// Publish publica un evento a todos los suscriptores
func (eb *EventBus) Publish(event Event) error {
	eb.mutex.RLock()
	handlers, exists := eb.subscribers[event.Type]
	eb.mutex.RUnlock()

	if !exists {
		log.Printf("No hay manejadores registrados para el evento: %s", event.Type)
		return nil
	}

	log.Printf("Publicando evento: %s a %d manejadores", event.Type, len(handlers))

	var wg sync.WaitGroup
	errors := make(chan error, len(handlers))

	for _, handler := range handlers {
		wg.Add(1)
		go func(h EventHandler) {
			defer wg.Done()
			if err := h.Handle(event); err != nil {
				errors <- fmt.Errorf("error en manejador %T: %v", h, err)
			}
		}(handler)
	}

	wg.Wait()
	close(errors)

	// Recolectar errores
	var errs []error
	for err := range errors {
		errs = append(errs, err)
	}

	if len(errs) > 0 {
		return fmt.Errorf("errores en el procesamiento del evento: %v", errs)
	}

	return nil
}

// CreateEvent crea un nuevo evento
func CreateEvent(eventType string, data map[string]interface{}, source string) Event {
	return Event{
		ID:        fmt.Sprintf("event_%d", time.Now().UnixNano()),
		Type:      eventType,
		Data:      data,
		Timestamp: time.Now(),
		Source:    source,
		Priority:  2, // prioridad media por defecto
	}
}

// EventTypes define los tipos de eventos del sistema
const (
	EventDatosRecolectados        = "datos_recolectados"
	EventDatosProcesados          = "datos_procesados"
	EventDatosPersistidos         = "datos_persistidos"
	EventReporteGenerado          = "reporte_generado"
	EventErrorProcesamiento       = "error_procesamiento"
	EventSincronizacionCompletada = "sincronizacion_completada"
	EventClientePotencialCreado   = "cliente_potencial_creado"
	EventVentaRegistrada          = "venta_registrada"
	EventStockActualizado         = "stock_actualizado"
)

// EventBusSingleton implementa el patrón Singleton para el bus de eventos
var (
	eventBusInstance *EventBus
	eventBusOnce     sync.Once
)

// GetEventBusInstance retorna la instancia única del bus de eventos
func GetEventBusInstance() *EventBus {
	eventBusOnce.Do(func() {
		eventBusInstance = NewEventBus()
	})
	return eventBusInstance
}

// EventLogger implementa un manejador de eventos para logging
type EventLogger struct{}

func (el *EventLogger) Handle(event Event) error {
	eventJSON, _ := json.MarshalIndent(event, "", "  ")
	log.Printf("Evento recibido: %s\n%s", event.Type, string(eventJSON))
	return nil
}

func (el *EventLogger) GetEventType() string {
	return "*" // maneja todos los tipos de eventos
}

// EventMetrics implementa un manejador de eventos para métricas
type EventMetrics struct {
	eventCounts map[string]int
	mutex       sync.RWMutex
}

func NewEventMetrics() *EventMetrics {
	return &EventMetrics{
		eventCounts: make(map[string]int),
	}
}

func (em *EventMetrics) Handle(event Event) error {
	em.mutex.Lock()
	defer em.mutex.Unlock()

	em.eventCounts[event.Type]++
	log.Printf("Métrica actualizada - Evento %s: %d total", event.Type, em.eventCounts[event.Type])
	return nil
}

func (em *EventMetrics) GetEventType() string {
	return "*" // maneja todos los tipos de eventos
}

func (em *EventMetrics) GetEventCount(eventType string) int {
	em.mutex.RLock()
	defer em.mutex.RUnlock()
	return em.eventCounts[eventType]
}
