package db

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/microsoft/go-mssqldb"
)

// Configuración de conexión
type Config struct {
	Server   string
	Port     int
	User     string
	Password string
	Database string
}

// Nueva conexión a SQL Server
func NewSQLServerConnection(cfg Config) (*sql.DB, error) {
	connectionString := fmt.Sprintf("server=%s;port=%d;user id=%s;password=%s;database=%s;",
		cfg.Server,
		cfg.Port,
		cfg.User,
		cfg.Password,
		cfg.Database,
	)

	// Abrir conexión
	db, err := sql.Open("sqlserver", connectionString)
	if err != nil {
		return nil, fmt.Errorf("error al abrir conexión: %v", err)
	}

	// Verificar conexión
	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("error al conectar con la base de datos: %v", err)
	}

	log.Println("¡Conexión exitosa a SQL Server!")
	return db, nil
}

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

// Cliente representa la tabla Cliente de la base de datos
type Cliente struct {
	ConsecutivoCompania         int        `gorm:"column:ConsecutivoCompania;primaryKey"`
	Consecutivo                 int        `gorm:"column:Consecutivo"`
	Codigo                      string     `gorm:"column:Codigo;primaryKey"`
	Nombre                      string     `gorm:"column:Nombre"`
	NumeroRIF                   *string    `gorm:"column:NumeroRIF"`
	NumeroNit                   *string    `gorm:"column:NumeroNit"`
	Direccion                   *string    `gorm:"column:Direccion"`
	Ciudad                      *string    `gorm:"column:Ciudad"`
	ZonaPostal                  *string    `gorm:"column:ZonaPostal"`
	Telefono                    *string    `gorm:"column:Telefono"`
	Fax                         *string    `gorm:"column:Fax"`
	Status                      *string    `gorm:"column:Status"`
	Contacto                    *string    `gorm:"column:Contacto"`
	ZonaDeCobranza              *string    `gorm:"column:ZonaDeCobranza"`
	CodigoVendedor              *string    `gorm:"column:CodigoVendedor"`
	RazonInactividad            *string    `gorm:"column:RazonInactividad"`
	Email                       *string    `gorm:"column:Email"`
	ActivarAvisoAlEscoger       string     `gorm:"column:ActivarAvisoAlEscoger"`
	TextoDelAviso               *string    `gorm:"column:TextoDelAviso"`
	CuentaContableCxc           *string    `gorm:"column:CuentaContableCxc"`
	CuentaContableIngresos      *string    `gorm:"column:CuentaContableIngresos"`
	CuentaContableAnticipo      *string    `gorm:"column:CuentaContableAnticipo"`
	InfoGalac                   *string    `gorm:"column:InfoGalac"`
	SectorDeNegocio             *string    `gorm:"column:SectorDeNegocio"`
	CodigoLote                  *string    `gorm:"column:CodigoLote"`
	NivelDePrecio               *string    `gorm:"column:NivelDePrecio"`
	Origen                      *string    `gorm:"column:Origen"`
	DiaCumpleanos               *int       `gorm:"column:DiaCumpleanos"`
	MesCumpleanos               *int       `gorm:"column:MesCumpleanos"`
	CorrespondenciaXenviar      string     `gorm:"column:CorrespondenciaXenviar"`
	EsExtranjero                string     `gorm:"column:EsExtranjero"`
	ClienteDesdeFecha           *time.Time `gorm:"column:ClienteDesdeFecha"`
	AQueSeDedicaElCliente       *string    `gorm:"column:AQueSeDedicaElCliente"`
	NombreOperador              *string    `gorm:"column:NombreOperador"`
	FechaUltimaModificacion     *time.Time `gorm:"column:FechaUltimaModificacion"`
	TipoDocumentoIdentificacion *string    `gorm:"column:TipoDocumentoIdentificacion"`
	TipoDeContribuyente         *string    `gorm:"column:TipoDeContribuyente"`
	CampoDefinible1             *string    `gorm:"column:CampoDefinible1"`
	FldTimeStamp                time.Time  `gorm:"column:fldTimeStamp"`
	ConsecutivoVendedor         int        `gorm:"column:ConsecutivoVendedor"`
}
