package clientes

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type Cliente struct {
	ConsecutivoCompania         int          `gorm:"column:ConsecutivoCompania;primaryKey"`
	Consecutivo                 int          `gorm:"column:Consecutivo"`
	Codigo                      string       `gorm:"column:Codigo;primaryKey"`
	Nombre                      string       `gorm:"column:Nombre"`
	NumeroRIF                   *string      `gorm:"column:NumeroRIF"`
	NumeroNit                   *string      `gorm:"column:NumeroNit"`
	Direccion                   *string      `gorm:"column:Direccion"`
	Ciudad                      *string      `gorm:"column:Ciudad"`
	ZonaPostal                  *string      `gorm:"column:ZonaPostal"`
	Telefono                    *string      `gorm:"column:Telefono"`
	Fax                         *string      `gorm:"column:Fax"`
	Status                      *string      `gorm:"column:Status"`
	Contacto                    *string      `gorm:"column:Contacto"`
	ZonaDeCobranza              *string      `gorm:"column:ZonaDeCobranza"`
	CodigoVendedor              *string      `gorm:"column:CodigoVendedor"`
	RazonInactividad            *string      `gorm:"column:RazonInactividad"`
	Email                       *string      `gorm:"column:Email"`
	ActivarAvisoAlEscoger       string       `gorm:"column:ActivarAvisoAlEscoger"`
	TextoDelAviso               *string      `gorm:"column:TextoDelAviso"`
	CuentaContableCxc           *string      `gorm:"column:CuentaContableCxc"`
	CuentaContableIngresos      *string      `gorm:"column:CuentaContableIngresos"`
	CuentaContableAnticipo      *string      `gorm:"column:CuentaContableAnticipo"`
	InfoGalac                   *string      `gorm:"column:InfoGalac"`
	SectorDeNegocio             *string      `gorm:"column:SectorDeNegocio"`
	CodigoLote                  *string      `gorm:"column:CodigoLote"`
	NivelDePrecio               *string      `gorm:"column:NivelDePrecio"`
	Origen                      *string      `gorm:"column:Origen"`
	DiaCumpleanos               *int         `gorm:"column:DiaCumpleanos"`
	MesCumpleanos               *int         `gorm:"column:MesCumpleanos"`
	CorrespondenciaXenviar      string       `gorm:"column:CorrespondenciaXenviar"`
	EsExtranjero                string       `gorm:"column:EsExtranjero"`
	ClienteDesdeFecha           *time.Time   `gorm:"column:ClienteDesdeFecha"`
	AQueSeDedicaElCliente       *string      `gorm:"column:AQueSeDedicaElCliente"`
	NombreOperador              *string      `gorm:"column:NombreOperador"`
	FechaUltimaModificacion     *time.Time   `gorm:"column:FechaUltimaModificacion"`
	TipoDocumentoIdentificacion *string      `gorm:"column:TipoDocumentoIdentificacion"`
	TipoDeContribuyente         *string      `gorm:"column:TipoDeContribuyente"`
	CampoDefinible1             *string      `gorm:"column:CampoDefinible1"`
	FldTimeStamp                sql.NullTime `gorm:"column:fldTimeStamp"`
	ConsecutivoVendedor         int          `gorm:"column:ConsecutivoVendedor"`
}

// Logger para registrar errores en un archivo
var errorLogger *log.Logger

// Inicializa el logger para errores
func initErrorLogger() {
	// Abrir o crear el archivo de log
	logFile, err := os.OpenFile("errorCliente.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Println("Error al abrir archivo de log de clientes:", err)
		return
	}

	// Inicializar el logger con timestamp
	errorLogger = log.New(logFile, "", log.Ldate|log.Ltime)
	log.Println("Logger de errores de clientes inicializado correctamente")
}

// logError registra un error en el archivo de log
func logError(mensaje string, err error) {
	if errorLogger != nil {
		errorLogger.Printf("[ERROR] %s: %v\n", mensaje, err)
	} else {
		// Si el logger no se inicializó correctamente, usa el logger estándar
		log.Printf("[ERROR] %s: %v\n", mensaje, err)
	}
}

func ClienteRoutes(api *gin.RouterGroup, db *sql.DB) {
	// Inicializar el logger de errores
	initErrorLogger()

	api.GET("/existe-cliente", buscarClientes(db))
}

func buscarClientes(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		requestTime := time.Now()
		requestID := requestTime.Format("20060102150405")

		rif := c.Query("rif")
		exacta := c.DefaultQuery("exacta", "no")

		if rif == "" {
			logError(requestID+" - Falta el parámetro RIF", nil)
			c.JSON(http.StatusBadRequest, gin.H{
				"message":    "Falta el RIF a buscar",
				"statusCode": http.StatusBadRequest,
				"success":    false,
			})
			return
		}

		busqueda := "="
		param := rif
		if _, ok := map[byte]bool{'V': true, 'G': true, 'J': true, 'E': true}[rif[0]]; !ok || exacta != "si" {
			busqueda = "LIKE"
			param = "%" + rif + "%"
		}

		logError(requestID+" - Iniciando búsqueda de códigos de cliente con RIF: "+rif+" (modo: "+exacta+")", nil)

		query := `
            SELECT Codigo
            FROM dbo.Cliente
            WHERE NumeroRIF ` + busqueda + ` @p1
            ORDER BY FechaUltimaModificacion DESC
        `

		rows, err := db.Query(query, param)
		if err != nil {
			logError(requestID+" - Error al consultar los códigos de cliente", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"message":    "Error al consultar los clientes",
				"error":      err.Error(),
				"statusCode": http.StatusInternalServerError,
				"success":    false,
			})
			return
		}
		defer rows.Close()

		var codigos []string

		for rows.Next() {
			var codigo string
			if err := rows.Scan(&codigo); err != nil {
				logError(requestID+" - Error al leer el código del cliente", err)
				c.JSON(http.StatusInternalServerError, gin.H{
					"message":    "Error al leer los datos del cliente",
					"error":      err.Error(),
					"statusCode": http.StatusInternalServerError,
					"success":    false,
				})
				return
			}
			codigos = append(codigos, codigo)
		}

		if err = rows.Err(); err != nil {
			logError(requestID+" - Error al iterar resultados de códigos", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"message":    "Error al procesar resultados",
				"error":      err.Error(),
				"statusCode": http.StatusInternalServerError,
				"success":    false,
			})
			return
		}

		if len(codigos) == 0 {
			logError(requestID+" - No se encontraron códigos de cliente con RIF: "+rif, nil)
			c.JSON(http.StatusOK, gin.H{
				"message":    "No se encontraron códigos de cliente",
				"statusCode": http.StatusOK,
				"success":    false,
			})
			return
		}

		logError(requestID+" - Búsqueda exitosa. Códigos encontrados: "+strconv.Itoa(len(codigos)), nil)
		c.JSON(http.StatusOK, gin.H{
			"message":    "Códigos de cliente encontrados",
			"data":       codigos,
			"count":      len(codigos),
			"statusCode": http.StatusOK,
			"success":    true,
		})
	}
}
