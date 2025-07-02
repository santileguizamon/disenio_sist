package entities

import (
	"time"
)

// Sucursal representa una sucursal en el sistema
type Sucursal struct {
	Nombre        string    `json:"nombre"`
	Direccion     string    `json:"direccion"`
	Telefono      string    `json:"telefono"`
	Email         string    `json:"email"`
	Ciudad        string    `json:"ciudad"`
	Estado        string    `json:"estado"`
	FechaApertura time.Time `json:"fecha_apertura"`

	// Configuración de integración
	APIEndpoint          string    `json:"api_endpoint"`
	APIKey               string    `json:"api_key"`
	APISecret            string    `json:"api_secret"`
	TipoSistema          string    `json:"tipo_sistema"`  // 'api', 'csv', 'excel', 'database'
	Configuracion        string    `json:"configuracion"` // JSON con configuración específica
	UltimaSincronizacion time.Time `json:"ultima_sincronizacion"`
}

// ObtenerParametros retorna los parámetros de conexión para la sucursal
func (s *Sucursal) ObtenerParametros() map[string]interface{} {
	return map[string]interface{}{
		"endpoint": s.APIEndpoint,
		"api_key":  s.APIKey,
		"secret":   s.APISecret,
		"tipo":     s.TipoSistema,
		"config":   s.Configuracion,
	}
}

// EsActiva verifica si la sucursal está activa
func (s *Sucursal) EsActiva() bool {
	return s.Estado == "activa"
}

// ActualizarSincronizacion actualiza la fecha de última sincronización
func (s *Sucursal) ActualizarSincronizacion() {
	s.UltimaSincronizacion = time.Now()
}

// EsValida verifica si la sucursal tiene los datos mínimos requeridos
func (s *Sucursal) EsValida() bool {
	return s.Nombre != "" && s.TipoSistema != ""
}

// ConfiguracionSistema representa la configuración específica de cada sistema
type ConfiguracionSistema struct {
	TipoSistema             string            `json:"tipo_sistema"`
	Parametros              map[string]string `json:"parametros"`
	MapeoCampos             map[string]string `json:"mapeo_campos"`
	Filtros                 []string          `json:"filtros"`
	IntervaloSincronizacion int               `json:"intervalo_sincronizacion"` // en minutos
}
