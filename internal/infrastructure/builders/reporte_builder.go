package builders

import (
	"fmt"
	"time"
)

// ReporteBuilder implementa el patrón Builder para construir reportes
type ReporteBuilder struct {
	reporte *Reporte
}

// Reporte representa un reporte generado por el sistema
type Reporte struct {
	ID           string                 `json:"id"`
	Tipo         string                 `json:"tipo"`
	FechaInicio  time.Time              `json:"fecha_inicio"`
	FechaFin     time.Time              `json:"fecha_fin"`
	SucursalID   uint                   `json:"sucursal_id"`
	Formato      string                 `json:"formato"`
	Datos        map[string]interface{} `json:"datos"`
	Filtros      []string               `json:"filtros"`
	Ordenamiento string                 `json:"ordenamiento"`
	GeneradoEn   time.Time              `json:"generado_en"`
	Estado       string                 `json:"estado"`
}

// NewReporteBuilder crea una nueva instancia del builder
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

// SetTipo establece el tipo de reporte
func (rb *ReporteBuilder) SetTipo(tipo string) *ReporteBuilder {
	rb.reporte.Tipo = tipo
	return rb
}

// SetFechas establece el rango de fechas del reporte
func (rb *ReporteBuilder) SetFechas(fechaInicio, fechaFin time.Time) *ReporteBuilder {
	rb.reporte.FechaInicio = fechaInicio
	rb.reporte.FechaFin = fechaFin
	return rb
}

// SetSucursal establece la sucursal para el reporte
func (rb *ReporteBuilder) SetSucursal(sucursalID uint) *ReporteBuilder {
	rb.reporte.SucursalID = sucursalID
	return rb
}

// SetFormato establece el formato del reporte
func (rb *ReporteBuilder) SetFormato(formato string) *ReporteBuilder {
	rb.reporte.Formato = formato
	return rb
}

// AddFiltro agrega un filtro al reporte
func (rb *ReporteBuilder) AddFiltro(filtro string) *ReporteBuilder {
	rb.reporte.Filtros = append(rb.reporte.Filtros, filtro)
	return rb
}

// SetOrdenamiento establece el ordenamiento del reporte
func (rb *ReporteBuilder) SetOrdenamiento(ordenamiento string) *ReporteBuilder {
	rb.reporte.Ordenamiento = ordenamiento
	return rb
}

// AddDato agrega un dato al reporte
func (rb *ReporteBuilder) AddDato(clave string, valor interface{}) *ReporteBuilder {
	rb.reporte.Datos[clave] = valor
	return rb
}

// SetEstado establece el estado del reporte
func (rb *ReporteBuilder) SetEstado(estado string) *ReporteBuilder {
	rb.reporte.Estado = estado
	return rb
}

// Build construye y retorna el reporte final
func (rb *ReporteBuilder) Build() (*Reporte, error) {
	// Validaciones
	if rb.reporte.Tipo == "" {
		return nil, fmt.Errorf("tipo de reporte es requerido")
	}

	if rb.reporte.Formato == "" {
		rb.reporte.Formato = "excel" // formato por defecto
	}

	if rb.reporte.ID == "" {
		rb.reporte.ID = fmt.Sprintf("reporte_%s_%d", rb.reporte.Tipo, time.Now().Unix())
	}

	return rb.reporte, nil
}

// BuildVentasPorSucursal construye un reporte específico de ventas por sucursal
func (rb *ReporteBuilder) BuildVentasPorSucursal(sucursalID uint, fechaInicio, fechaFin time.Time) *ReporteBuilder {
	return rb.
		SetTipo("ventas_por_sucursal").
		SetSucursal(sucursalID).
		SetFechas(fechaInicio, fechaFin).
		SetFormato("excel").
		AddFiltro("estado=completada").
		SetOrdenamiento("fecha_venta DESC")
}

// BuildStockActual construye un reporte específico de stock actual
func (rb *ReporteBuilder) BuildStockActual(sucursalID uint) *ReporteBuilder {
	return rb.
		SetTipo("stock_actual").
		SetSucursal(sucursalID).
		SetFormato("excel").
		AddFiltro("stock_actual > 0").
		SetOrdenamiento("stock_actual ASC")
}

// BuildClientesPotenciales construye un reporte específico de clientes potenciales
func (rb *ReporteBuilder) BuildClientesPotenciales(sucursalID uint, fechaInicio, fechaFin time.Time) *ReporteBuilder {
	return rb.
		SetTipo("clientes_potenciales").
		SetSucursal(sucursalID).
		SetFechas(fechaInicio, fechaFin).
		SetFormato("excel").
		AddFiltro("estado=activo").
		SetOrdenamiento("fecha_captura DESC")
}
