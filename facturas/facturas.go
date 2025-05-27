package facturas

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

////////////////////////////////////////////////////////
////////////////////////////////////////////////////////
////////////////////////////////////////////////////////

/*
	Logger Interno para el registro de errores y problemas.
	Solo se instancia en este modulo y generar el archivo
	errores_facturas.log
*/

// Logger para registrar errores en un archivo
var errorLogger *log.Logger

func initErrorLogger() {
	logFile, err := os.OpenFile("errores_facturas.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Println("Error al abrir archivo de log:", err)
		return
	}
	errorLogger = log.New(logFile, "", log.Ldate|log.Ltime)
	log.Println("Logger de errores inicializado correctamente")
}

func logError(mensaje string, err error) {
	if errorLogger != nil {
		errorLogger.Printf("[ERROR] %s: %v\n", mensaje, err)
	} else {
		log.Printf("[ERROR] %s: %v\n", mensaje, err)
	}
}

////////////////////////////////////////////////////////
////////////////////////////////////////////////////////
////////////////////////////////////////////////////////

func Facturas(api *gin.RouterGroup, db *sql.DB) {
	initErrorLogger()
	api.GET("/", buscarFacturas(db))
}

////////////////////////////////////////////////////////
////////////////////////////////////////////////////////
////////////////////////////////////////////////////////

/*
La funcion permite la busqueda de las facturas en
SqlServer, recibe los siguientes parametros de busqueda
mediante Querys:

mes: debe ser numerico entre 1 y  12
anio (año): año de busqueda de la factura
codigoCliente: alfanumerico
page: debe ser numerico o se setea en 1
pageZise: debe ser numerico o se setea en 1000, no se

	recomienda valores muy altos ya que tiende
	a generar fallas en las API
*/
var meses = map[string]string{
	"ENERO":      "1",
	"FEBRERO":    "2",
	"MARZO":      "3",
	"ABRIL":      "4",
	"MAYO":       "5",
	"JUNIO":      "6",
	"JULIO":      "7",
	"AGOSTO":     "8",
	"SEPTIEMBRE": "9",
	"OCTUBRE":    "10",
	"NOVIEMBRE":  "11",
	"DICIEMBRE":  "12",
}

var estadosFactura = map[string]string{
	"EMITIDA":      "0",
	"BORRADOR":     "2",
	"NOTA_CREDITO": "1",
	"0":            "0",
	"2":            "2",
	"1":            "1",
}

func buscarFacturas(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		mesStr := c.Query("mes")
		anioStr := c.Query("anio")
		codigoCliente := c.Query("codigoCliente") // <-- NUEVO
		pageStr := c.DefaultQuery("page", "1")
		pageSizeStr := c.DefaultQuery("pageSize", "1000")
		odooQuery := c.Query("odoo")
		mesNombre := strings.ToUpper(c.Query("mesNombre"))         // Ejemplo: "ABRIL"
		estadoFactura := strings.ToUpper(c.Query("estadoFactura")) // 0 = Emitida, 2 = Borrador, 1 = Nota de Credito en status Factura
		numeroControl := strings.ToLower(c.Query("numeroControl")) // Convertir a minúsculas para comparación

		requestTime := time.Now()
		requestID := requestTime.Format("20060102150405")
		logError(requestID+" - Iniciando consulta de facturas", nil)

		// Verificar si solo se quieren números de control
		soloNumerosControl := numeroControl == "si" || numeroControl == "true"

		// ------- Seteo de los parametros de Busqueda ----- //

		// Convertir page y pageSize a int
		page, err := strconv.Atoi(pageStr)
		if err != nil || page < 1 {
			page = 1
		}

		pageSize, err := strconv.Atoi(pageSizeStr)
		if err != nil || pageSize < 1 {
			pageSize = 1000
		}

		offset := (page - 1) * pageSize

		baseQuery := `
            FROM dbo.factura
            WHERE 1=1
        `

		var filterQuery string
		params := []interface{}{}

		var mes, anio int
		var hayFiltroFecha bool

		// ------ Busqueda por Estado de Factura ------ //
		if estadoFactura != "" {
			// Verificar si el estado es válido (número o nombre)
			codigoEstado, ok := estadosFactura[estadoFactura]
			if !ok {
				mensaje := "Estado de factura inválido. Use: 0 (EMITIDA), 2 (BORRADOR), 1 (NOTA_CREDITO) o sus nombres"
				logError(requestID+" - "+mensaje, nil)
				estadoFactura = ""
			}

			// Agregar filtro para buscar por estado
			filterQuery += " AND StatusFactura = @estadoFactura"
			params = append(params, sql.Named("estadoFactura", codigoEstado))
		}

		// ------ Busqueda por Nombre de Mes en Observaciones ------ //
		if mesNombre != "" {
			// Verificar si el nombre del mes es válido
			if _, ok := meses[mesNombre]; !ok {
				mensaje := "Nombre de mes inválido. Use: ENERO, FEBRERO, ..., DICIEMBRE"
				logError(requestID+" - "+mensaje, nil)
				mesNombre = "" // Si el mes no es valido, se deja vacio el campo
			}

			// Agregar filtro para buscar el mes en observaciones
			filterQuery += " AND Observaciones LIKE @mesObs"
			params = append(params, sql.Named("mesObs", "%"+mesNombre+"%"))

			// Si quieres buscar exactamente el patrón "MES-Suscripcion:"
			// filterQuery += " AND Observaciones LIKE @mesObs"
			// params = append(params, sql.Named("mesObs", "%"+mesNombre+"-Suscripcion:%"))
		}

		// ------ Busqueda por Año y Mes ------ //
		if mesStr != "" {
			mes, err = strconv.Atoi(mesStr)
			if err != nil || mes < 1 || mes > 12 {
				mensaje := "Mes inválido, debe ser un número entre 1 y 12"
				logError(requestID+" - "+mensaje, err)
				c.JSON(http.StatusBadRequest, gin.H{"error": mensaje})
				return
			}
			hayFiltroFecha = true
		}

		if anioStr != "" {
			anio, err = strconv.Atoi(anioStr)
			if err != nil || anio < 1900 || anio > 2100 {
				mensaje := "Año inválido, debe ser un número entre 1900 y 2100"
				logError(requestID+" - "+mensaje, err)
				c.JSON(http.StatusBadRequest, gin.H{"error": mensaje})
				return
			}
			hayFiltroFecha = true
		}

		if hayFiltroFecha {
			if mesStr != "" && anioStr != "" {
				filterQuery += " AND MONTH(Fecha) = @mes AND YEAR(Fecha) = @anio"
				params = append(params, sql.Named("mes", mes), sql.Named("anio", anio))
			} else if mesStr != "" {
				filterQuery += " AND MONTH(Fecha) = @mes"
				params = append(params, sql.Named("mes", mes))
			} else if anioStr != "" {
				filterQuery += " AND YEAR(Fecha) = @anio"
				params = append(params, sql.Named("anio", anio))
			}
		}

		// ----- Busca por campo observacion ------ //
		// Jembi lo quiere porque aqui colocan el Codigo SUB de Odoo
		if odooQuery != "" {
			filterQuery += " AND Observaciones LIKE @odoo"
			params = append(params, sql.Named("odoo", "%"+odooQuery+"%"))
		}

		// NUEVO: filtro por CodigoCliente si viene en la query
		if codigoCliente != "" {
			filterQuery += " AND CodigoCliente = @codigoCliente"
			params = append(params, sql.Named("codigoCliente", codigoCliente))
		}
		// ------- Seteo de los parametros de Busqueda FIN ----- //

		// ------- Ejecucion de las Busquedas ----- //

		// Obtener el TOTAL DE REGISTROS //
		countQuery := "SELECT COUNT(*) " + baseQuery + filterQuery
		var total int
		err = db.QueryRow(countQuery, params...).Scan(&total)
		if err != nil {
			mensaje := "Error al obtener cantidad total de facturas"
			logError(requestID+" - "+mensaje, err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": mensaje + ": " + err.Error()})
			return
		}

		logError(requestID+" - Total de facturas encontradas: "+strconv.Itoa(total), nil)

		// CONSULTA PRINCIPAL - Modificada según si solo queremos números de control
		var query string
		if soloNumerosControl {
			// Solo seleccionar NumeroControl
			query = `
                SELECT DISTINCT NumeroControl
            ` + baseQuery + filterQuery + `
                AND NumeroControl IS NOT NULL 
                AND NumeroControl != ''
                ORDER BY NumeroControl
                OFFSET @offset ROWS FETCH NEXT @pageSize ROWS ONLY
            `
		} else {
			// Consulta completa original
			query = `
                SELECT
                    ConsecutivoCompania, Numero, Fecha, CodigoCliente, CodigoVendedor, Observaciones, TotalMontoExento,
                    TotalBaseImponible, TotalRenglones, TotalIVA, TotalFactura, PorcentajeDescuento, CodigoNota1,
                    CodigoNota2, Moneda, NivelDePrecio, ReservarMercancia, FechaDeRetiro, CodigoAlmacen, StatusFactura,
                    TipoDeDocumento, InsertadaManualmente, FacturaHistorica, Cancelada, UsarDireccionFiscal,
                    NoDirDespachoAimprimir, CambioABolivares, MontoDelAbono, FechaDeVencimiento, CondicionesDePago,
                    FormaDeLaInicial, PorcentajeDeLaInicial, NumeroDeCuotas, MontoDeLasCuotas, MontoUltimaCuota,
                    Talonario, FormaDePago, NumDiasDeVencimiento1aCuota, EditarMontoCuota, NumeroControl,
                    TipoDeTransaccion, NumeroFacturaAfectada, NumeroPlanillaExportacion, TipoDeVenta, UsaMaquinaFiscal,
                    CodigoMaquinaRegistradora, NumeroDesde, NumeroHasta, NumeroControlHasta, MontoIvaRetenido,
                    FechaAplicacionRetIVA, NumeroComprobanteRetIVA, FechaComprobanteRetIVA, SeRetuvoIVA,
                    FacturaConPreciosSinIva, VueltoDelCobroDirecto, ConsecutivoCaja, GeneraCobroDirecto,
                    FechaDeFacturaAfectada, FechaDeEntrega, PorcentajeDescuento1, PorcentajeDescuento2,
                    MontoDescuento1, MontoDescuento2, CodigoLote, Devolucion, PorcentajeAlicuota1, PorcentajeAlicuota2,
                    PorcentajeAlicuota3, MontoIVAAlicuota1, MontoIVAAlicuota2, MontoIVAAlicuota3, MontoGravableAlicuota1,
                    MontoGravableAlicuota2, MontoGravableAlicuota3, RealizoCierreZ, NumeroComprobanteFiscal,
                    SerialMaquinaFiscal, AplicarPromocion, RealizoCierreX, HoraModificacion, FormaDeCobro,
                    OtraFormaDeCobro, NoCotizacionDeOrigen, NoContrato, ConsecutivoVehiculo, ConsecutivoAlmacen,
                    NumeroResumenDiario, NoControlDespachoDeOrigen, ImprimeFiscal, EsDiferida, EsOriginalmenteDiferida,
                    SeContabilizoIvaDiferido, AplicaDecretoIvaEspecial, EsGeneradaPorPuntoDeVenta, CambioMonedaCXC,
                    CambioMostrarTotalEnDivisas, CodigoMonedaDeCobro, GeneradaPorNotaEntrega, EmitidaEnFacturaNumero,
                    CodigoMoneda, NombreOperador, FechaUltimaModificacion, NumeroParaResumen, NroDiasMantenerCambioAMonedaLocal,
                    FechaLimiteCambioAMonedaLocal, GeneradoPor, BaseImponibleIGTF, IGTFML, IGTFME, AlicuotaIGTF,
                    MotivoDeAnulacion, ProveedorImprentaDigital, ConsecutivoVendedor, ImprentaDigitalGUID
             ` + baseQuery + filterQuery + `
                ORDER BY Fecha DESC
                OFFSET @offset ROWS FETCH NEXT @pageSize ROWS ONLY
            `
		}

		params = append(params, sql.Named("offset", offset), sql.Named("pageSize", pageSize))

		rows, err := db.Query(query, params...)
		if err != nil {
			mensaje := "Error al ejecutar consulta de facturas"
			logError(requestID+" - "+mensaje, err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": mensaje + ": " + err.Error()})
			return
		}
		defer rows.Close()
		// ------- Ejecucion de las Busquedas FIN ----- //

		// ------- Formateo de los resultados ----- //
		if soloNumerosControl {
			// Solo números de control
			var numerosControl []string

			for rows.Next() {
				var numeroCtrl string
				err := rows.Scan(&numeroCtrl)
				if err != nil {
					mensaje := "Error al leer números de control"
					logError(requestID+" - "+mensaje, err)
					c.JSON(http.StatusInternalServerError, gin.H{"error": mensaje + ": " + err.Error()})
					return
				}
				numerosControl = append(numerosControl, numeroCtrl)
			}

			// Respuesta simplificada para números de control
			totalPages := (total + pageSize - 1) / pageSize
			respuesta := gin.H{
				"count":       len(numerosControl),
				"total":       total,
				"page":        page,
				"pageSize":    pageSize,
				"totalPages":  totalPages,
				"isFirstPage": page == 1,
				"isLastPage":  page == totalPages,
				"hasNextPage": page < totalPages,
				"hasPrevPage": page > 1,
				"data":        numerosControl, // Array de strings en lugar de objetos
			}

			// Agregar filtros aplicados
			if estadoFactura != "" || mesNombre != "" || hayFiltroFecha || odooQuery != "" || codigoCliente != "" {
				respuesta["filtros"] = map[string]interface{}{}
				if estadoFactura != "" {
					respuesta["filtros"].(map[string]interface{})["estadoFactura"] = estadoFactura
				}
				if mesNombre != "" {
					respuesta["filtros"].(map[string]interface{})["mesNombre"] = mesNombre
				}
				if hayFiltroFecha {
					respuesta["filtros"].(map[string]interface{})["mes"] = mes
					respuesta["filtros"].(map[string]interface{})["anio"] = anio
				}
				if odooQuery != "" {
					respuesta["filtros"].(map[string]interface{})["odoo"] = odooQuery
				}
				if codigoCliente != "" {
					respuesta["filtros"].(map[string]interface{})["codigoCliente"] = codigoCliente
				}
			}

			logError(requestID+" - Consulta completada (solo números de control). Página: "+strconv.Itoa(page)+", Números devueltos: "+strconv.Itoa(len(numerosControl))+" / "+strconv.Itoa(total), nil)
			c.JSON(http.StatusOK, respuesta)
			return
		}

		// Lógica original para facturas completas
		var facturas []Factura

		for rows.Next() {
			var factura Factura
			err := rows.Scan(
				// todos los campos igual que antes
				&factura.ConsecutivoCompania, &factura.Numero, &factura.Fecha, &factura.CodigoCliente, &factura.CodigoVendedor,
				&factura.Observaciones, &factura.TotalMontoExento, &factura.TotalBaseImponible, &factura.TotalRenglones,
				&factura.TotalIVA, &factura.TotalFactura, &factura.PorcentajeDescuento, &factura.CodigoNota1, &factura.CodigoNota2,
				&factura.Moneda, &factura.NivelDePrecio, &factura.ReservarMercancia, &factura.FechaDeRetiro, &factura.CodigoAlmacen,
				&factura.StatusFactura, &factura.TipoDeDocumento, &factura.InsertadaManualmente, &factura.FacturaHistorica,
				&factura.Cancelada, &factura.UsarDireccionFiscal, &factura.NoDirDespachoAimprimir, &factura.CambioABolivares,
				&factura.MontoDelAbono, &factura.FechaDeVencimiento, &factura.CondicionesDePago, &factura.FormaDeLaInicial,
				&factura.PorcentajeDeLaInicial, &factura.NumeroDeCuotas, &factura.MontoDeLasCuotas, &factura.MontoUltimaCuota,
				&factura.Talonario, &factura.FormaDePago, &factura.NumDiasDeVencimiento1aCuota, &factura.EditarMontoCuota,
				&factura.NumeroControl, &factura.TipoDeTransaccion, &factura.NumeroFacturaAfectada, &factura.NumeroPlanillaExportacion,
				&factura.TipoDeVenta, &factura.UsaMaquinaFiscal, &factura.CodigoMaquinaRegistradora, &factura.NumeroDesde,
				&factura.NumeroHasta, &factura.NumeroControlHasta, &factura.MontoIvaRetenido, &factura.FechaAplicacionRetIVA,
				&factura.NumeroComprobanteRetIVA, &factura.FechaComprobanteRetIVA, &factura.SeRetuvoIVA,
				&factura.FacturaConPreciosSinIva, &factura.VueltoDelCobroDirecto, &factura.ConsecutivoCaja,
				&factura.GeneraCobroDirecto, &factura.FechaDeFacturaAfectada, &factura.FechaDeEntrega,
				&factura.PorcentajeDescuento1, &factura.PorcentajeDescuento2, &factura.MontoDescuento1,
				&factura.MontoDescuento2, &factura.CodigoLote, &factura.Devolucion, &factura.PorcentajeAlicuota1,
				&factura.PorcentajeAlicuota2, &factura.PorcentajeAlicuota3, &factura.MontoIVAAlicuota1, &factura.MontoIVAAlicuota2,
				&factura.MontoIVAAlicuota3, &factura.MontoGravableAlicuota1, &factura.MontoGravableAlicuota2,
				&factura.MontoGravableAlicuota3, &factura.RealizoCierreZ, &factura.NumeroComprobanteFiscal,
				&factura.SerialMaquinaFiscal, &factura.AplicarPromocion, &factura.RealizoCierreX, &factura.HoraModificacion,
				&factura.FormaDeCobro, &factura.OtraFormaDeCobro, &factura.NoCotizacionDeOrigen, &factura.NoContrato,
				&factura.ConsecutivoVehiculo, &factura.ConsecutivoAlmacen, &factura.NumeroResumenDiario,
				&factura.NoControlDespachoDeOrigen, &factura.ImprimeFiscal, &factura.EsDiferida, &factura.EsOriginalmenteDiferida,
				&factura.SeContabilizoIvaDiferido, &factura.AplicaDecretoIvaEspecial, &factura.EsGeneradaPorPuntoDeVenta,
				&factura.CambioMonedaCXC, &factura.CambioMostrarTotalEnDivisas, &factura.CodigoMonedaDeCobro,
				&factura.GeneradaPorNotaEntrega, &factura.EmitidaEnFacturaNumero, &factura.CodigoMoneda, &factura.NombreOperador,
				&factura.FechaUltimaModificacion, &factura.NumeroParaResumen, &factura.NroDiasMantenerCambioAMonedaLocal,
				&factura.FechaLimiteCambioAMonedaLocal, &factura.GeneradoPor, &factura.BaseImponibleIGTF, &factura.IGTFML,
				&factura.IGTFME, &factura.AlicuotaIGTF, &factura.MotivoDeAnulacion, &factura.ProveedorImprentaDigital,
				&factura.ConsecutivoVendedor, &factura.ImprentaDigitalGUID,
			)
			if err != nil {
				mensaje := "Error al leer datos de facturas"
				logError(requestID+" - "+mensaje, err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": mensaje + ": " + err.Error()})
				return
			}
			facturas = append(facturas, factura)
		}
		// ------- Formateo de los resultados FIN ----- //

		// ------- Seteo de la respuesta ----- //
		totalPages := (total + pageSize - 1) / pageSize

		respuesta := gin.H{
			"count":       len(facturas),
			"total":       total,
			"page":        page,
			"pageSize":    pageSize,
			"totalPages":  totalPages,
			"isFirstPage": page == 1,
			"isLastPage":  page == totalPages,
			"hasNextPage": page < totalPages,
			"hasPrevPage": page > 1,
			"data":        facturas,
		}
		if estadoFactura != "" {
			if respuesta["filtros"] == nil {
				respuesta["filtros"] = map[string]interface{}{}
			}
			respuesta["filtros"].(map[string]interface{})["estadoFactura"] = estadoFactura
		}

		if mesNombre != "" {
			if respuesta["filtros"] == nil {
				respuesta["filtros"] = map[string]interface{}{}
			}
			respuesta["filtros"].(map[string]interface{})["mesNombre"] = mesNombre
		}

		if hayFiltroFecha {
			respuesta["filtros"] = map[string]interface{}{
				"mes":  mes,
				"anio": anio,
			}
		}

		if odooQuery != "" {
			if respuesta["filtros"] == nil {
				respuesta["filtros"] = map[string]interface{}{}
			}
			respuesta["filtros"].(map[string]interface{})["odoo"] = odooQuery
		}

		if codigoCliente != "" {
			if respuesta["filtros"] == nil {
				respuesta["filtros"] = map[string]interface{}{}
			}
			respuesta["filtros"].(map[string]interface{})["codigoCliente"] = codigoCliente
		}

		logError(requestID+" - Consulta completada. Página: "+strconv.Itoa(page)+", Registros devueltos: "+strconv.Itoa(len(facturas))+" / "+strconv.Itoa(total), nil)

		c.JSON(http.StatusOK, respuesta)
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

type Factura struct {
	ConsecutivoCompania               int        `gorm:"column:ConsecutivoCompania;primaryKey"`
	Numero                            string     `gorm:"column:Numero;primaryKey"`
	Fecha                             time.Time  `gorm:"column:Fecha"`
	CodigoCliente                     *string    `gorm:"column:CodigoCliente"`
	CodigoVendedor                    *string    `gorm:"column:CodigoVendedor"`
	Observaciones                     *string    `gorm:"column:Observaciones"`
	TotalMontoExento                  *float64   `gorm:"column:TotalMontoExento"`
	TotalBaseImponible                *float64   `gorm:"column:TotalBaseImponible"`
	TotalRenglones                    *float64   `gorm:"column:TotalRenglones"`
	TotalIVA                          *float64   `gorm:"column:TotalIVA"`
	TotalFactura                      *float64   `gorm:"column:TotalFactura"`
	PorcentajeDescuento               *float64   `gorm:"column:PorcentajeDescuento"`
	CodigoNota1                       *string    `gorm:"column:CodigoNota1"`
	CodigoNota2                       *string    `gorm:"column:CodigoNota2"`
	Moneda                            *string    `gorm:"column:Moneda"`
	NivelDePrecio                     *string    `gorm:"column:NivelDePrecio"`
	ReservarMercancia                 string     `gorm:"column:ReservarMercancia"`
	FechaDeRetiro                     *time.Time `gorm:"column:FechaDeRetiro"`
	CodigoAlmacen                     *string    `gorm:"column:CodigoAlmacen"`
	StatusFactura                     *string    `gorm:"column:StatusFactura"`
	TipoDeDocumento                   string     `gorm:"column:TipoDeDocumento;primaryKey"`
	InsertadaManualmente              string     `gorm:"column:InsertadaManualmente"`
	FacturaHistorica                  string     `gorm:"column:FacturaHistorica"`
	Cancelada                         string     `gorm:"column:Cancelada"`
	UsarDireccionFiscal               string     `gorm:"column:UsarDireccionFiscal"`
	NoDirDespachoAimprimir            *int       `gorm:"column:NoDirDespachoAimprimir"`
	CambioABolivares                  *float64   `gorm:"column:CambioABolivares"`
	MontoDelAbono                     *float64   `gorm:"column:MontoDelAbono"`
	FechaDeVencimiento                *time.Time `gorm:"column:FechaDeVencimiento"`
	CondicionesDePago                 *string    `gorm:"column:CondicionesDePago"`
	FormaDeLaInicial                  *string    `gorm:"column:FormaDeLaInicial"`
	PorcentajeDeLaInicial             *float64   `gorm:"column:PorcentajeDeLaInicial"`
	NumeroDeCuotas                    *int       `gorm:"column:NumeroDeCuotas"`
	MontoDeLasCuotas                  *float64   `gorm:"column:MontoDeLasCuotas"`
	MontoUltimaCuota                  *float64   `gorm:"column:MontoUltimaCuota"`
	Talonario                         *string    `gorm:"column:Talonario"`
	FormaDePago                       *string    `gorm:"column:FormaDePago"`
	NumDiasDeVencimiento1aCuota       *int       `gorm:"column:NumDiasDeVencimiento1aCuota"`
	EditarMontoCuota                  *string    `gorm:"column:EditarMontoCuota"`
	NumeroControl                     *string    `gorm:"column:NumeroControl"`
	TipoDeTransaccion                 *string    `gorm:"column:TipoDeTransaccion"`
	NumeroFacturaAfectada             *string    `gorm:"column:NumeroFacturaAfectada"`
	NumeroPlanillaExportacion         *string    `gorm:"column:NumeroPlanillaExportacion"`
	TipoDeVenta                       *string    `gorm:"column:TipoDeVenta"`
	UsaMaquinaFiscal                  *string    `gorm:"column:UsaMaquinaFiscal"`
	CodigoMaquinaRegistradora         *string    `gorm:"column:CodigoMaquinaRegistradora"`
	NumeroDesde                       *string    `gorm:"column:NumeroDesde"`
	NumeroHasta                       *string    `gorm:"column:NumeroHasta"`
	NumeroControlHasta                *string    `gorm:"column:NumeroControlHasta"`
	MontoIvaRetenido                  *float64   `gorm:"column:MontoIvaRetenido"`
	FechaAplicacionRetIVA             *time.Time `gorm:"column:FechaAplicacionRetIVA"`
	NumeroComprobanteRetIVA           *int       `gorm:"column:NumeroComprobanteRetIVA"`
	FechaComprobanteRetIVA            *time.Time `gorm:"column:FechaComprobanteRetIVA"`
	SeRetuvoIVA                       *string    `gorm:"column:SeRetuvoIVA"`
	FacturaConPreciosSinIva           string     `gorm:"column:FacturaConPreciosSinIva"`
	VueltoDelCobroDirecto             *float64   `gorm:"column:VueltoDelCobroDirecto"`
	ConsecutivoCaja                   *int       `gorm:"column:ConsecutivoCaja"`
	GeneraCobroDirecto                string     `gorm:"column:GeneraCobroDirecto"`
	FechaDeFacturaAfectada            time.Time  `gorm:"column:FechaDeFacturaAfectada"`
	FechaDeEntrega                    *time.Time `gorm:"column:FechaDeEntrega"`
	PorcentajeDescuento1              *float64   `gorm:"column:PorcentajeDescuento1"`
	PorcentajeDescuento2              *float64   `gorm:"column:PorcentajeDescuento2"`
	MontoDescuento1                   *float64   `gorm:"column:MontoDescuento1"`
	MontoDescuento2                   *float64   `gorm:"column:MontoDescuento2"`
	CodigoLote                        *string    `gorm:"column:CodigoLote"`
	Devolucion                        string     `gorm:"column:Devolucion"`
	PorcentajeAlicuota1               *float64   `gorm:"column:PorcentajeAlicuota1"`
	PorcentajeAlicuota2               *float64   `gorm:"column:PorcentajeAlicuota2"`
	PorcentajeAlicuota3               *float64   `gorm:"column:PorcentajeAlicuota3"`
	MontoIVAAlicuota1                 *float64   `gorm:"column:MontoIVAAlicuota1"`
	MontoIVAAlicuota2                 *float64   `gorm:"column:MontoIVAAlicuota2"`
	MontoIVAAlicuota3                 *float64   `gorm:"column:MontoIVAAlicuota3"`
	MontoGravableAlicuota1            *float64   `gorm:"column:MontoGravableAlicuota1"`
	MontoGravableAlicuota2            *float64   `gorm:"column:MontoGravableAlicuota2"`
	MontoGravableAlicuota3            *float64   `gorm:"column:MontoGravableAlicuota3"`
	RealizoCierreZ                    string     `gorm:"column:RealizoCierreZ"`
	NumeroComprobanteFiscal           *string    `gorm:"column:NumeroComprobanteFiscal"`
	SerialMaquinaFiscal               *string    `gorm:"column:SerialMaquinaFiscal"`
	AplicarPromocion                  string     `gorm:"column:AplicarPromocion"`
	RealizoCierreX                    string     `gorm:"column:RealizoCierreX"`
	HoraModificacion                  *string    `gorm:"column:HoraModificacion"`
	FormaDeCobro                      string     `gorm:"column:FormaDeCobro"`
	OtraFormaDeCobro                  *string    `gorm:"column:OtraFormaDeCobro"`
	NoCotizacionDeOrigen              *string    `gorm:"column:NoCotizacionDeOrigen"`
	NoContrato                        *string    `gorm:"column:NoContrato"`
	ConsecutivoVehiculo               *int       `gorm:"column:ConsecutivoVehiculo"`
	ConsecutivoAlmacen                int        `gorm:"column:ConsecutivoAlmacen"`
	NumeroResumenDiario               *string    `gorm:"column:NumeroResumenDiario"`
	NoControlDespachoDeOrigen         *string    `gorm:"column:NoControlDespachoDeOrigen"`
	ImprimeFiscal                     string     `gorm:"column:ImprimeFiscal"`
	EsDiferida                        string     `gorm:"column:EsDiferida"`
	EsOriginalmenteDiferida           string     `gorm:"column:EsOriginalmenteDiferida"`
	SeContabilizoIvaDiferido          string     `gorm:"column:SeContabilizoIvaDiferido"`
	AplicaDecretoIvaEspecial          string     `gorm:"column:AplicaDecretoIvaEspecial"`
	EsGeneradaPorPuntoDeVenta         string     `gorm:"column:EsGeneradaPorPuntoDeVenta"`
	CambioMonedaCXC                   float64    `gorm:"column:CambioMonedaCXC"`
	CambioMostrarTotalEnDivisas       float64    `gorm:"column:CambioMostrarTotalEnDivisas"`
	CodigoMonedaDeCobro               *string    `gorm:"column:CodigoMonedaDeCobro"`
	GeneradaPorNotaEntrega            *string    `gorm:"column:GeneradaPorNotaEntrega"`
	EmitidaEnFacturaNumero            *string    `gorm:"column:EmitidaEnFacturaNumero"`
	CodigoMoneda                      string     `gorm:"column:CodigoMoneda"`
	NombreOperador                    *string    `gorm:"column:NombreOperador"`
	FechaUltimaModificacion           *time.Time `gorm:"column:FechaUltimaModificacion"`
	NumeroParaResumen                 *int       `gorm:"column:NumeroParaResumen"`
	NroDiasMantenerCambioAMonedaLocal *int       `gorm:"column:NroDiasMantenerCambioAMonedaLocal"`
	FechaLimiteCambioAMonedaLocal     *time.Time `gorm:"column:FechaLimiteCambioAMonedaLocal"`
	FldTimeStamp                      time.Time  `gorm:"column:fldTimeStamp"`
	GeneradoPor                       *string    `gorm:"column:GeneradoPor"`
	BaseImponibleIGTF                 *float64   `gorm:"column:BaseImponibleIGTF"`
	IGTFML                            *float64   `gorm:"column:IGTFML"`
	IGTFME                            *float64   `gorm:"column:IGTFME"`
	AlicuotaIGTF                      *float64   `gorm:"column:AlicuotaIGTF"`
	MotivoDeAnulacion                 *string    `gorm:"column:MotivoDeAnulacion"`
	ProveedorImprentaDigital          string     `gorm:"column:ProveedorImprentaDigital"`
	ConsecutivoVendedor               int        `gorm:"column:ConsecutivoVendedor"`
	ImprentaDigitalGUID               *string    `gorm:"column:ImprentaDigitalGUID"`
}
