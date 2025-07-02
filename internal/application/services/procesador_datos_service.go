package services

import (
	"context"
	"fmt"
	"log"
	"time"

	"sistema-gestion-informacion/internal/infrastructure/events"
)

// ProcesadorDatosService implementa la lógica de procesamiento de datos
type ProcesadorDatosService struct {
	eventBus *events.EventBus
}

// NewProcesadorDatosService crea una nueva instancia del servicio
func NewProcesadorDatosService(eventBus *events.EventBus) *ProcesadorDatosService {
	return &ProcesadorDatosService{
		eventBus: eventBus,
	}
}

// DatosCrudos representa los datos brutos recibidos de las fuentes
type DatosCrudos struct {
	Origen     string                   `json:"origen"`
	Tipo       string                   `json:"tipo"` // 'cliente', 'venta', 'producto', 'stock'
	Datos      []map[string]interface{} `json:"datos"`
	Timestamp  time.Time                `json:"timestamp"`
	SucursalID uint                     `json:"sucursal_id"`
}

// RegistroProcesado representa un registro después del procesamiento
type RegistroProcesado struct {
	Entidad     string                 `json:"entidad"`
	Datos       map[string]interface{} `json:"datos"`
	Validado    bool                   `json:"validado"`
	Errores     []string               `json:"errores"`
	Enriquecido bool                   `json:"enriquecido"`
}

// ProcesarLote procesa un lote de datos brutos
func (pds *ProcesadorDatosService) ProcesarLote(ctx context.Context, datosCrudos *DatosCrudos) error {
	log.Printf("Iniciando procesamiento de lote desde %s", datosCrudos.Origen)

	// Publicar evento de inicio de procesamiento
	pds.eventBus.Publish(events.CreateEvent(
		events.EventDatosRecolectados,
		map[string]interface{}{
			"origen":      datosCrudos.Origen,
			"tipo":        datosCrudos.Tipo,
			"cantidad":    len(datosCrudos.Datos),
			"sucursal_id": datosCrudos.SucursalID,
		},
		"procesador_datos",
	))

	// Normalizar datos
	datosNormalizados, err := pds.normalizarDatos(datosCrudos)
	if err != nil {
		pds.publicarError("error_normalizacion", err)
		return fmt.Errorf("error normalizando datos: %v", err)
	}

	// Validar datos
	datosValidados, err := pds.validarDatos(datosNormalizados)
	if err != nil {
		pds.publicarError("error_validacion", err)
		return fmt.Errorf("error validando datos: %v", err)
	}

	// Enriquecer datos
	datosEnriquecidos, err := pds.enriquecerDatos(datosValidados)
	if err != nil {
		pds.publicarError("error_enriquecimiento", err)
		return fmt.Errorf("error enriqueciendo datos: %v", err)
	}

	// Eliminar duplicados
	datosFinales, err := pds.eliminarDuplicados(datosEnriquecidos)
	if err != nil {
		pds.publicarError("error_deduplicacion", err)
		return fmt.Errorf("error eliminando duplicados: %v", err)
	}

	// Persistir datos
	err = pds.persistirDatos(ctx, datosFinales)
	if err != nil {
		pds.publicarError("error_persistencia", err)
		return fmt.Errorf("error persistiendo datos: %v", err)
	}

	// Publicar evento de procesamiento completado
	pds.eventBus.Publish(events.CreateEvent(
		events.EventDatosProcesados,
		map[string]interface{}{
			"origen":            datosCrudos.Origen,
			"tipo":              datosCrudos.Tipo,
			"registros_inicial": len(datosCrudos.Datos),
			"registros_final":   len(datosFinales),
			"sucursal_id":       datosCrudos.SucursalID,
		},
		"procesador_datos",
	))

	log.Printf("Procesamiento completado: %d registros procesados", len(datosFinales))
	return nil
}

// normalizarDatos convierte los datos a un formato estándar
func (pds *ProcesadorDatosService) normalizarDatos(datosCrudos *DatosCrudos) ([]map[string]interface{}, error) {
	log.Printf("Normalizando %d registros", len(datosCrudos.Datos))

	var datosNormalizados []map[string]interface{}

	for _, dato := range datosCrudos.Datos {
		normalizado := make(map[string]interface{})

		// Normalizar fechas
		if fecha, ok := dato["fecha"]; ok {
			if fechaStr, ok := fecha.(string); ok {
				if fechaTime, err := time.Parse("2006-01-02", fechaStr); err == nil {
					normalizado["fecha"] = fechaTime
				}
			}
		}

		// Normalizar valores numéricos
		if precio, ok := dato["precio"]; ok {
			if precioFloat, ok := precio.(float64); ok {
				normalizado["precio"] = precioFloat
			}
		}

		// Normalizar strings (trim, lowercase)
		for key, value := range dato {
			if str, ok := value.(string); ok {
				normalizado[key] = pds.normalizarString(str)
			} else {
				normalizado[key] = value
			}
		}

		datosNormalizados = append(datosNormalizados, normalizado)
	}

	return datosNormalizados, nil
}

// validarDatos valida la integridad de los datos
func (pds *ProcesadorDatosService) validarDatos(datos []map[string]interface{}) ([]map[string]interface{}, error) {
	log.Printf("Validando %d registros", len(datos))

	var datosValidados []map[string]interface{}

	for _, dato := range datos {
		if pds.esRegistroValido(dato) {
			datosValidados = append(datosValidados, dato)
		} else {
			log.Printf("Registro inválido descartado: %v", dato)
		}
	}

	return datosValidados, nil
}

// enriquecerDatos enriquece los datos con información adicional
func (pds *ProcesadorDatosService) enriquecerDatos(datos []map[string]interface{}) ([]map[string]interface{}, error) {
	log.Printf("Enriqueciendo %d registros", len(datos))

	for i, dato := range datos {
		// Enriquecer productos con información de APIs externas
		if sku, ok := dato["sku"].(string); ok {
			if infoProducto, err := pds.obtenerInfoProducto(sku); err == nil {
				for key, value := range infoProducto {
					dato[key] = value
				}
			}
		}

		// Agregar metadatos
		dato["enriquecido_en"] = time.Now()
		dato["fuente_enriquecimiento"] = "api_externa"

		datos[i] = dato
	}

	return datos, nil
}

// eliminarDuplicados elimina registros duplicados
func (pds *ProcesadorDatosService) eliminarDuplicados(datos []map[string]interface{}) ([]map[string]interface{}, error) {
	log.Printf("Eliminando duplicados de %d registros", len(datos))

	seen := make(map[string]bool)
	var datosUnicos []map[string]interface{}

	for _, dato := range datos {
		key := pds.generarClaveUnica(dato)
		if !seen[key] {
			seen[key] = true
			datosUnicos = append(datosUnicos, dato)
		}
	}

	return datosUnicos, nil
}

// persistirDatos persiste los datos en la base de datos
func (pds *ProcesadorDatosService) persistirDatos(ctx context.Context, datos []map[string]interface{}) error {
	log.Printf("Persistiendo %d registros", len(datos))

	for _, dato := range datos {
		if err := pds.persistirRegistro(ctx, dato); err != nil {
			log.Printf("Error persistiendo registro: %v", err)
			continue
		}
	}

	// Publicar evento de persistencia completada
	pds.eventBus.Publish(events.CreateEvent(
		events.EventDatosPersistidos,
		map[string]interface{}{
			"registros_persistidos": len(datos),
		},
		"procesador_datos",
	))

	return nil
}

// Métodos auxiliares
func (pds *ProcesadorDatosService) normalizarString(str string) string {
	// Implementar normalización de strings
	return str
}

func (pds *ProcesadorDatosService) esRegistroValido(dato map[string]interface{}) bool {
	// Implementar validación de registros
	return true
}

func (pds *ProcesadorDatosService) obtenerInfoProducto(sku string) (map[string]interface{}, error) {
	// Implementar consulta a API externa
	return make(map[string]interface{}), nil
}

func (pds *ProcesadorDatosService) generarClaveUnica(dato map[string]interface{}) string {
	// Implementar generación de clave única
	return fmt.Sprintf("%v", dato)
}

func (pds *ProcesadorDatosService) persistirRegistro(ctx context.Context, dato map[string]interface{}) error {
	// Implementar persistencia específica según el tipo de dato
	return nil
}

func (pds *ProcesadorDatosService) publicarError(tipo string, err error) {
	pds.eventBus.Publish(events.CreateEvent(
		events.EventErrorProcesamiento,
		map[string]interface{}{
			"tipo_error": tipo,
			"mensaje":    err.Error(),
		},
		"procesador_datos",
	))
}
