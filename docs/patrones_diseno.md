# Patrones de Diseño Implementados

## Descripción General

Este documento describe los patrones de diseño implementados en el Sistema de Procesamiento de Datos, siguiendo los principios SOLID y las mejores prácticas de arquitectura de software. El sistema utiliza una arquitectura basada en eventos y procesamiento en memoria.

## 1. Patrón Singleton

### Propósito
Garantizar que una clase tenga una única instancia y proporcionar un punto de acceso global a ella.

### Implementación

#### EventBus Singleton
```go
// internal/infrastructure/events/event_bus.go
var (
    eventBusInstance *EventBus
    eventBusOnce     sync.Once
)

func GetEventBusInstance() *EventBus {
    eventBusOnce.Do(func() {
        eventBusInstance = NewEventBus()
    })
    return eventBusInstance
}

func NewEventBus() *EventBus {
    return &EventBus{
        subscribers: make(map[string][]EventHandler),
    }
}
```

### Ventajas
- **Control de acceso**: Garantiza una única instancia del bus de eventos
- **Lazy initialization**: La instancia se crea solo cuando se necesita
- **Thread-safe**: Implementación segura para concurrencia usando `sync.Once`
- **Gestión de recursos**: Centraliza la gestión del sistema de eventos

### Uso en el Sistema
- **EventBus**: Una única instancia global del bus de eventos compartida por toda la aplicación
- **Punto de acceso centralizado**: Todos los componentes usan la misma instancia para publicar y suscribirse a eventos

## 2. Patrón Observer (Event-Driven Architecture)

### Propósito
Definir una dependencia uno-a-muchos entre objetos, de modo que cuando un objeto cambia de estado, todos sus dependientes son notificados y actualizados automáticamente.

### Implementación

#### EventBus (Subject)
```go
// internal/infrastructure/events/event_bus.go
type EventBus struct {
    subscribers map[string][]EventHandler
    mutex       sync.RWMutex
}

type EventHandler interface {
    Handle(event Event) error
    GetEventType() string
}

func (eb *EventBus) Subscribe(eventType string, handler EventHandler) {
    eb.mutex.Lock()
    defer eb.mutex.Unlock()
    eb.subscribers[eventType] = append(eb.subscribers[eventType], handler)
}

func (eb *EventBus) Publish(event Event) error {
    eb.mutex.RLock()
    defer eb.mutex.RUnlock()
    
    handlers, exists := eb.subscribers[event.Type]
    if !exists {
        return nil
    }
    
    for _, handler := range handlers {
        if err := handler.Handle(event); err != nil {
            return err
        }
    }
    return nil
}
```

#### Handlers (Observers)
```go
// internal/infrastructure/events/event_bus.go
type DatosProcesadosHandler struct{}

func (h *DatosProcesadosHandler) Handle(event Event) error {
    fmt.Printf("🔔 Datos procesados: %d registros\n", event.Data["registros_procesados"])
    return nil
}

func (h *DatosProcesadosHandler) GetEventType() string {
    return "datos.procesados"
}

type EventLogger struct{}

func (h *EventLogger) Handle(event Event) error {
    fmt.Printf("📝 Evento registrado: %s - %v\n", event.Type, event.Data)
    return nil
}

func (h *EventLogger) GetEventType() string {
    return "*" // Maneja todos los eventos
}

type EventMetrics struct{}

func (h *EventMetrics) Handle(event Event) error {
    fmt.Printf("📊 Métrica actualizada: %s\n", event.Type)
    return nil
}

func (h *EventMetrics) GetEventType() string {
    return "*" // Maneja todos los eventos
}
```

### Ventajas
- **Desacoplamiento**: Los componentes no dependen directamente entre sí
- **Extensibilidad**: Fácil agregar nuevos manejadores de eventos
- **Escalabilidad**: Procesamiento asíncrono de eventos
- **Mantenibilidad**: Cambios en un componente no afectan otros

### Uso en el Sistema
- **Procesamiento de datos**: Notificación cuando se completan procesos
- **Logging**: Registro automático de eventos importantes
- **Métricas**: Recopilación automática de estadísticas
- **Generación de reportes**: Activación automática tras procesamiento

## 3. Patrón Builder

### Propósito
Separar la construcción de un objeto complejo de su representación, permitiendo crear diferentes representaciones del mismo proceso de construcción.

### Implementación

#### ReporteBuilder
```go
// internal/infrastructure/builders/reporte_builder.go
type ReporteBuilder struct {
    reporte *Reporte
}

type Reporte struct {
    ID          uint
    Tipo        string
    GeneradoEn  time.Time
    Estado      string
    FechaInicio time.Time
    FechaFin    time.Time
    Datos       map[string]interface{}
    Filtros     []string
    Resumen     string
}

func NewReporteBuilder() *ReporteBuilder {
    return &ReporteBuilder{
        reporte: &Reporte{
            GeneradoEn: time.Now(),
            Estado:     "pendiente",
            Datos:      make(map[string]interface{}),
            Filtros:    make([]string, 0),
        },
    }
}

func (rb *ReporteBuilder) SetTipo(tipo string) *ReporteBuilder {
    rb.reporte.Tipo = tipo
    return rb
}

func (rb *ReporteBuilder) SetFechas(fechaInicio, fechaFin time.Time) *ReporteBuilder {
    rb.reporte.FechaInicio = fechaInicio
    rb.reporte.FechaFin = fechaFin
    return rb
}

func (rb *ReporteBuilder) SetDatos(datos map[string]interface{}) *ReporteBuilder {
    rb.reporte.Datos = datos
    return rb
}

func (rb *ReporteBuilder) SetFiltros(filtros []string) *ReporteBuilder {
    rb.reporte.Filtros = filtros
    return rb
}

func (rb *ReporteBuilder) SetResumen(resumen string) *ReporteBuilder {
    rb.reporte.Resumen = resumen
    return rb
}

func (rb *ReporteBuilder) Build() (*Reporte, error) {
    // Validaciones
    if rb.reporte.Tipo == "" {
        return nil, fmt.Errorf("tipo de reporte es requerido")
    }
    
    if rb.reporte.FechaInicio.IsZero() || rb.reporte.FechaFin.IsZero() {
        return nil, fmt.Errorf("fechas de inicio y fin son requeridas")
    }
    
    if rb.reporte.FechaInicio.After(rb.reporte.FechaFin) {
        return nil, fmt.Errorf("fecha de inicio no puede ser posterior a fecha de fin")
    }
    
    // Marcar como completado
    rb.reporte.Estado = "completado"
    
    return rb.reporte, nil
}
```

### Uso del Builder
```go
// Ejemplo de uso en el servicio
reporte, err := NewReporteBuilder().
    SetTipo("ventas").
    SetFechas(fechaInicio, fechaFin).
    SetDatos(map[string]interface{}{
        "total_ventas": 502.50,
        "cantidad_transacciones": 1,
        "producto_mas_vendido": "Producto A",
    }).
    SetFiltros([]string{"fecha_inicio", "fecha_fin"}).
    SetResumen("Reporte de ventas del período").
    Build()
```

### Ventajas
- **Flexibilidad**: Permite construir objetos complejos paso a paso
- **Reutilización**: Métodos de construcción reutilizables
- **Validación**: Validación en cada paso de la construcción
- **Legibilidad**: Código más legible y expresivo

### Uso en el Sistema
- **Reportes**: Construcción de reportes con diferentes configuraciones
- **Validación**: Validación de datos durante la construcción
- **Configuración**: Construcción de configuraciones complejas

## Aplicación de Principios SOLID

### 1. Single Responsibility Principle (SRP)
- **EventBus**: Responsable solo de la gestión de eventos
- **ReporteBuilder**: Responsable solo de la construcción de reportes
- **Handlers**: Cada handler tiene una responsabilidad específica

### 2. Open/Closed Principle (OCP)
- **EventBus**: Abierto para extensión (nuevos handlers), cerrado para modificación
- **ReporteBuilder**: Abierto para extensión (nuevos métodos), cerrado para modificación

### 3. Liskov Substitution Principle (LSP)
- **EventHandler**: Cualquier implementación puede sustituir a otra
- **ReporteBuilder**: Diferentes builders pueden sustituirse entre sí

### 4. Interface Segregation Principle (ISP)
- **EventHandler**: Interfaz específica y cohesiva
- **ReporteBuilder**: Métodos específicos para construcción

### 5. Dependency Inversion Principle (DIP)
- **EventBus**: Depende de abstracciones (EventHandler)
- **Servicios**: Dependen de abstracciones, no de implementaciones concretas

## Beneficios de los Patrones Implementados

### 1. Arquitectura Limpia
- **Separación de responsabilidades**: Cada componente tiene una función específica
- **Independencia de frameworks**: Lógica de negocio independiente
- **Testabilidad**: Fácil testing de cada componente
- **Independencia de base de datos**: Procesamiento en memoria

### 2. Escalabilidad
- **Horizontal**: Fácil agregar nuevos handlers de eventos
- **Vertical**: Fácil optimizar componentes específicos
- **Concurrencia**: Patrones thread-safe para operaciones concurrentes

### 3. Mantenibilidad
- **Código limpio**: Estructura clara y organizada
- **Documentación**: Patrones bien documentados
- **Estándares**: Seguimiento de mejores prácticas
- **Refactoring**: Fácil modificación sin afectar otros componentes

## Flujo de Eventos en el Sistema

1. **Recepción de datos**: `POST /api/procesar` recibe datos crudos
2. **Procesamiento**: El servicio procesa y depura los datos
3. **Publicación de evento**: Se publica el evento `datos.procesados`
4. **Notificación a handlers**: Todos los handlers suscritos reciben el evento
5. **Generación de reporte**: Se construye un reporte usando el Builder
6. **Publicación de evento**: Se publica el evento `reporte.generado`

## Conclusión

La implementación de estos patrones de diseño proporciona al sistema:

1. **Robustez**: Manejo robusto de errores y casos edge
2. **Flexibilidad**: Fácil adaptación a nuevos requerimientos
3. **Escalabilidad**: Capacidad de crecer sin comprometer rendimiento
4. **Mantenibilidad**: Código fácil de mantener y evolucionar
5. **Testabilidad**: Componentes fáciles de probar
6. **Reutilización**: Componentes reutilizables en diferentes contextos

Estos patrones forman la base de una arquitectura sólida que cumple con los requisitos del sistema y permite el crecimiento futuro. 