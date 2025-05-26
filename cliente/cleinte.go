package clientes

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

/////////////////////////////////////////////////////
/////////////////////////////////////////////////////
/////////////////////////////////////////////////////

/*
	Logger Interno para el registro de errores y problemas.
	Solo se instancia en este modulo y generar el archivo
	errores_cliente.log
*/

// Logger para registrar errores en un archivo
var errorLogger *log.Logger

// Inicializa el logger para errores
func initErrorLogger() {
	// Abrir o crear el archivo de log
	logFile, err := os.OpenFile("errores_cliente.log.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
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

/////////////////////////////////////////////////////
/////////////////////////////////////////////////////
/////////////////////////////////////////////////////

func ClienteRoutes(api *gin.RouterGroup, db *sql.DB) {
	// Inicializar el logger de errores
	initErrorLogger()

	api.GET("/existe-cliente", buscarClientes(db))
}

/////////////////////////////////////////////////////
/////////////////////////////////////////////////////
/////////////////////////////////////////////////////

// Función para construir la consulta de búsqueda de clientes
func buildClientQuery(rif string, exacta string) (string, string) {
	busqueda := "="
	param := rif
	// La inicial del documento debe ser uno de los siguiente valores
	if _, ok := map[byte]bool{'V': true, 'G': true, 'J': true, 'E': true}[rif[0]]; !ok || exacta != "si" {
		// Busqueda aproximada o exacta
		busqueda = "LIKE"
		param = "%" + rif + "%"
	}

	return busqueda, param
}

///////////////////////////////////////////////////////////////

/*
Función para obtener el codigo del cliente
Ejecuta un query en la BD
Retorna un Slice de Clientes
*/
func obtenerCodigosCliente(db *sql.DB, requestID, query string, param string) ([]string, error) {
	rows, err := db.Query(query, param)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var codigos []string
	for rows.Next() {
		var codigo string
		if err := rows.Scan(&codigo); err != nil {
			return nil, err
		}
		codigos = append(codigos, codigo)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return codigos, nil
}

///////////////////////////////////////////////////////////////

/*
Función para obtener detalles completos de clientes
Ejecuta un query en la BD
Retorna un Slice de Clientes
*/
func obtenerDetallesClientes(db *sql.DB, requestID string, codigos []string) ([]Cliente, error) {
	// Construir la consulta IN para múltiples códigos
	codigosPlaceholders := make([]string, len(codigos))
	for i := range codigos {
		codigosPlaceholders[i] = fmt.Sprintf("@p%d", i+1)
	}

	query := `
        SELECT 
            ConsecutivoCompania, Consecutivo, Codigo, Nombre, NumeroRIF, NumeroNit, 
            Direccion, Ciudad, ZonaPostal, Telefono, Fax, Status, Contacto, 
            ZonaDeCobranza, CodigoVendedor, RazonInactividad, Email, 
            ActivarAvisoAlEscoger, TextoDelAviso, CuentaContableCxc, 
            CuentaContableIngresos, CuentaContableAnticipo, InfoGalac, 
            SectorDeNegocio, CodigoLote, NivelDePrecio, Origen, DiaCumpleanos, 
            MesCumpleanos, CorrespondenciaXenviar, EsExtranjero, ClienteDesdeFecha, 
            AQueSeDedicaElCliente, NombreOperador, FechaUltimaModificacion, 
            TipoDocumentoIdentificacion, TipoDeContribuyente, CampoDefinible1, 
            ConsecutivoVendedor
        FROM dbo.Cliente
        WHERE Codigo IN (` + strings.Join(codigosPlaceholders, ",") + `)
    `

	// Convertir codigos a interface{}
	params := make([]interface{}, len(codigos))
	for i, v := range codigos {
		params[i] = v
	}

	rows, err := db.Query(query, params...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var clientes []Cliente
	for rows.Next() {
		var cliente Cliente
		err := rows.Scan(
			&cliente.ConsecutivoCompania, &cliente.Consecutivo, &cliente.Codigo,
			&cliente.Nombre, &cliente.NumeroRIF, &cliente.NumeroNit,
			&cliente.Direccion, &cliente.Ciudad, &cliente.ZonaPostal,
			&cliente.Telefono, &cliente.Fax, &cliente.Status, &cliente.Contacto,
			&cliente.ZonaDeCobranza, &cliente.CodigoVendedor,
			&cliente.RazonInactividad, &cliente.Email,
			&cliente.ActivarAvisoAlEscoger, &cliente.TextoDelAviso,
			&cliente.CuentaContableCxc, &cliente.CuentaContableIngresos,
			&cliente.CuentaContableAnticipo, &cliente.InfoGalac,
			&cliente.SectorDeNegocio, &cliente.CodigoLote,
			&cliente.NivelDePrecio, &cliente.Origen, &cliente.DiaCumpleanos,
			&cliente.MesCumpleanos, &cliente.CorrespondenciaXenviar,
			&cliente.EsExtranjero, &cliente.ClienteDesdeFecha,
			&cliente.AQueSeDedicaElCliente, &cliente.NombreOperador,
			&cliente.FechaUltimaModificacion,
			&cliente.TipoDocumentoIdentificacion,
			&cliente.TipoDeContribuyente, &cliente.CampoDefinible1,
			&cliente.ConsecutivoVendedor,
		)
		if err != nil {
			return nil, err
		}
		clientes = append(clientes, cliente)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return clientes, nil
}

///////////////////////////////////////////////////////////////

/*
	Permite la obtencion del codigo del cliente.
	Se debe enviar la CI/RIF en el query.
	Tambien se puede solicitar una busqueda exacta.

	NOTA: Dado las inconsistencia en Galac infomado por
	el departamento administrativo, este endpoint retorna los
	resultados en un slice/array/lista de coincidencias, dado
	que nos infomarn de documentos de identidad dupplicados.
*/

// Función principal de búsqueda de clientes
func buscarClientes(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		requestTime := time.Now()
		requestID := requestTime.Format("20060102150405")

		// -----  Verificacion de parametros de busqueda ----- //
		// Obtencion de Querys
		rif := c.Query("rif")
		exacta := c.DefaultQuery("exacta", "no")
		cliente := c.DefaultQuery("cliente", "no")
		codigo := c.Query("codigo") // Nota: cambié "Codigo" a minúsculas para consistencia

		// Verificar que al menos uno de los parámetros (rif o codigo) esté presente
		if rif == "" && codigo == "" {
			logError(requestID+" - Falta el parámetro RIF o Código", nil)
			c.JSON(http.StatusBadRequest, gin.H{
				"message":    "Falta el RIF o Código a buscar",
				"statusCode": http.StatusBadRequest,
				"success":    false,
			})
			return
		}

		var query string
		var param string
		var err error
		var codigos []string

		// -----  Setting del Query ----- //
		if codigo != "" {
			// Búsqueda por código de cliente (siempre exacta)
			logError(requestID+" - Iniciando búsqueda de cliente por código: "+codigo, nil)
			query = `
				SELECT Codigo
				FROM dbo.Cliente
				WHERE Codigo = @p1
				ORDER BY FechaUltimaModificacion DESC
			`
			param = codigo
		} else {
			// Búsqueda por RIF (como antes)
			logError(requestID+" - Iniciando búsqueda de códigos de cliente con RIF: "+rif+" (modo: "+exacta+")", nil)
			query, param = buildClientQuery(rif, exacta)
			query = `
				SELECT Codigo
				FROM dbo.Cliente
				WHERE NumeroRIF ` + query + ` @p1
				ORDER BY FechaUltimaModificacion DESC
			`
		}

		// ----- Busqueda de códigos ---- //
		codigos, err = obtenerCodigosCliente(db, requestID, query, param)
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

		// Resto del código permanece igual...
		// ----- Manejo de resultado de búsqueda ---- //
		if len(codigos) == 0 {
			logError(requestID+" - No se encontraron códigos de cliente", nil)
			c.JSON(http.StatusOK, gin.H{
				"message":    "No se encontraron códigos de cliente",
				"statusCode": http.StatusOK,
				"success":    false,
			})
			return
		}

		// ----- Procesamiento de respuesta ----- //
		// Si se solicitan detalles completos de cliente
		if cliente == "si" {
			// Obtener detalles completos de clientes
			clientes, err := obtenerDetallesClientes(db, requestID, codigos)
			if err != nil {
				logError(requestID+" - Error al obtener detalles de clientes", err)
				c.JSON(http.StatusInternalServerError, gin.H{
					"message":    "Error al obtener detalles de clientes",
					"error":      err.Error(),
					"statusCode": http.StatusInternalServerError,
					"success":    false,
				})
				return
			}

			logError(requestID+" - Búsqueda exitosa. Clientes encontrados: "+strconv.Itoa(len(clientes)), nil)
			c.JSON(http.StatusOK, gin.H{
				"message":    "Detalles de clientes encontrados",
				"data":       clientes,
				"count":      len(clientes),
				"statusCode": http.StatusOK,
				"success":    true,
			})
			return
		}

		// Respuesta por defecto con solo códigos
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

/////////////////////////////////////////////////////
/////////////////////////////////////////////////////
/////////////////////////////////////////////////////

/*
	El struct representa el modelo extraido de la BD.
	Se utilizo gorm para realizar la extraccion del modelo
	pero despues resulto dificil su uso, debido a eso se
	utilizo otra libreria para la ejecucion del SQL
	sin embargo se deja todo el codigo como esta,
*/

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
