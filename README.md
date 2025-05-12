# API Facturas Galac 🖨️

# Introducción

La API de facturas Galac tiene por objetivo el permitir la lectura de las facturas/borradores presentes en el sistema dado que, no poseemos acceso a Galac mediante API para poder automatizar los procesos para una rápida facturación como lo requiere el SENIAT.

Para lograr este objetivo la API también permite la lectura del modelo Clientes.

La API no posee (ni debe poseer) métodos de escritura, dado que contraviene las ultimas normas del SENIAT (mayo 2025), además de que no sabemos como puede resultar inconveniente en los procesos de Galac.

## Tecnologias

- Go
- Gin
- database/sql (Std library)
- [go-mssqldb](http://github.com/microsoft/go-mssqldb)

**Nota:** No se usa docker dado que la API se aloja en un servidor Windows.

# Uso

La API posee 3 endpoints, sin embargo el primero solo sirve para verificar el estado de la API.

Todos los endpoints usa método GET.

## Clientes 👤

**Endpoint:** {URL}/clientes/existe-cliente

**Querys:**

- rif: Formato **V123456** ⇒ Obligatorio
- cliente: **si** ⇒ No obligatorio, si se setea en si, se trae structs de la entidad Cliente.

**Nota:** Se retorna siempre un Array/Lista/Slice de resultados dado que según informa administración, hay cliente duplicados.

**Ejemplo:**  [http://localhost:5000/clientes/existe-cliente?rif=V072533034](http://localhost:5000/clientes/existe-cliente?rif=V072533034)
Respuesta:

Como no se setea el valor **cliente,** retorna el código del cliente únicamente.

```json
{
  "count": 1,
  "data": [
    "\"A91XXF0EJ"
  ],
  "message": "Códigos de cliente encontrados",
  "statusCode": 200,
  "success": true
}
```

**Ejemplo:**  [http://localhost:5000/clientes/existe-cliente?rif=V072533034&cliente=si](http://localhost:5000/clientes/existe-cliente?rif=V072533034&cliente=si)
Respuesta:
Al setear el valor **cliente** en “si”, retorna el struct cliente.

```json
{
  "count": 1,
  "data": [
    {
      "ConsecutivoCompania": 2,
      "Consecutivo": 53904,
      "Codigo": "\"A91XXF0EJ",
      "Nombre": "MARIA DEL VALLE CASTELLANOS PEREZ",
      "NumeroRIF": "V072533034",
      "NumeroNit": "",
      "Direccion": "SECTOR INDIANA SUR AV INTERCOMUNAL CASA #22-A",
      "Ciudad": "VALENCIA",
      "ZonaPostal": "",
      "Telefono": "584144065239",
      "Fax": "",
      "Status": "0",
      "Contacto": "MARIA DEL VALLE CASTELLANOS PEREZ",
      "ZonaDeCobranza": "MARACAY",
      "CodigoVendedor": "00001",
      "RazonInactividad": "",
      "Email": "mariavallitacastellano@gmail.com",
      "ActivarAvisoAlEscoger": "N",
      "TextoDelAviso": "",
      "CuentaContableCxc": "",
      "CuentaContableIngresos": "",
      "CuentaContableAnticipo": "",
      "InfoGalac": " ",
      "SectorDeNegocio": "No Asignado",
      "CodigoLote": "0000000154",
      "NivelDePrecio": "0",
      "Origen": "0",
      "DiaCumpleanos": 0,
      "MesCumpleanos": 0,
      "CorrespondenciaXenviar": "N",
      "EsExtranjero": "N",
      "ClienteDesdeFecha": "2025-04-01T00:00:00Z",
      "AQueSeDedicaElCliente": "",
      "NombreOperador": "HRIVAS",
      "FechaUltimaModificacion": "2025-04-21T00:00:00Z",
      "TipoDocumentoIdentificacion": " ",
      "TipoDeContribuyente": "0",
      "CampoDefinible1": "",
      "FldTimeStamp": {
        "Time": "0001-01-01T00:00:00Z",
        "Valid": false
      },
      "ConsecutivoVendedor": 1
    }
  ],
  "message": "Detalles de clientes encontrados",
  "statusCode": 200,
  "success": true
}
```

## Facturas 📝

**Endpoint:** {URL}/facturas

**Querys:**

- mes: numérico
- año: numérico
- codigoCliente: alfanumérico
- page: numérico, por defecto 1
- pageSize: numérico por defecto 1000

**Nota:** Si no se envían valores por query la API se traerá los últimos mil valores en paginas de mil, siendo los datos presentes la primera pagina.

**Ejemplo:**   [http://localhost:5000/facturas](http://localhost:5000/facturas/?page=1)
Respuesta:

```json
{
  "count": 1000,
  "data": [
    {
      "ConsecutivoCompania": 2,
      "Numero": "00476373",
      "Fecha": "2025-04-29T00:00:00Z",
      "CodigoCliente": "NSUB75885",
      "CodigoVendedor": "00001",
      "Observaciones": "ABRIL-Suscripcion: SUB75885",
      "TotalMontoExento": 0,
      "TotalBaseImponible": 1750.25,
      "TotalRenglones": 1750.25,
      "TotalIVA": 280.04,
      "TotalFactura": 2030.29,
      "PorcentajeDescuento": 0,
      "CodigoNota1": "",
      "CodigoNota2": "",
      "Moneda": "Bolívar",
      "NivelDePrecio": "0",
      "ReservarMercancia": "N",
      "FechaDeRetiro": "2025-03-28T00:00:00Z",
      "CodigoAlmacen": "UNICO",
      "StatusFactura": "0",
      "TipoDeDocumento": "0",
      "InsertadaManualmente": "N",
      "FacturaHistorica": "N",
      "Cancelada": "N",
      "UsarDireccionFiscal": "N",
      "NoDirDespachoAimprimir": 0,
      "CambioABolivares": 1,
      "MontoDelAbono": 0,
      "FechaDeVencimiento": "2025-04-29T00:00:00Z",
      "CondicionesDePago": "Contado",
      "FormaDeLaInicial": "0",
      "PorcentajeDeLaInicial": 0,
      "NumeroDeCuotas": 1,
      "MontoDeLasCuotas": 1750.25,
      "MontoUltimaCuota": 1750.25,
      "Talonario": "0",
      "FormaDePago": "0",
      "NumDiasDeVencimiento1aCuota": 30,
      "EditarMontoCuota": "N",
      "NumeroControl": "00-00397233",
      "TipoDeTransaccion": "0",
      "NumeroFacturaAfectada": "",
      "NumeroPlanillaExportacion": "",
      "TipoDeVenta": "0",
      "UsaMaquinaFiscal": "N",
      "CodigoMaquinaRegistradora": "",
      "NumeroDesde": "",
      "NumeroHasta": "",
      "NumeroControlHasta": "",
      "MontoIvaRetenido": 0,
      "FechaAplicacionRetIVA": "2025-03-28T00:00:00Z",
      "NumeroComprobanteRetIVA": 0,
      "FechaComprobanteRetIVA": "2025-03-28T00:00:00Z",
      "SeRetuvoIVA": "N",
      "FacturaConPreciosSinIva": "N",
      "VueltoDelCobroDirecto": null,
      "ConsecutivoCaja": 0,
      "GeneraCobroDirecto": "N",
      "FechaDeFacturaAfectada": "2025-04-29T00:00:00Z",
      "FechaDeEntrega": "2025-03-28T00:00:00Z",
      "PorcentajeDescuento1": 0,
      "PorcentajeDescuento2": 0,
      "MontoDescuento1": 0,
      "MontoDescuento2": 0,
      "CodigoLote": "0000000181",
      "Devolucion": "N",
      "PorcentajeAlicuota1": 16,
      "PorcentajeAlicuota2": 8,
      "PorcentajeAlicuota3": 31,
      "MontoIVAAlicuota1": 280.04,
      "MontoIVAAlicuota2": 0,
      "MontoIVAAlicuota3": 0,
      "MontoGravableAlicuota1": 1750.25,
      "MontoGravableAlicuota2": 0,
      "MontoGravableAlicuota3": 0,
      "RealizoCierreZ": "N",
      "NumeroComprobanteFiscal": "0",
      "SerialMaquinaFiscal": "",
      "AplicarPromocion": "N",
      "RealizoCierreX": "N",
      "HoraModificacion": "00:00",
      "FormaDeCobro": "0",
      "OtraFormaDeCobro": "",
      "NoCotizacionDeOrigen": "",
      "NoContrato": "",
      "ConsecutivoVehiculo": 0,
      "ConsecutivoAlmacen": 1,
      "NumeroResumenDiario": "",
      "NoControlDespachoDeOrigen": "",
      "ImprimeFiscal": "N",
      "EsDiferida": "N",
      "EsOriginalmenteDiferida": "N",
      "SeContabilizoIvaDiferido": "N",
      "AplicaDecretoIvaEspecial": "N",
      "EsGeneradaPorPuntoDeVenta": "N",
      "CambioMonedaCXC": 1,
      "CambioMostrarTotalEnDivisas": 1,
      "CodigoMonedaDeCobro": "VED",
      "GeneradaPorNotaEntrega": "0",
      "EmitidaEnFacturaNumero": "",
      "CodigoMoneda": "VED",
      "NombreOperador": "HRIVAS",
      "FechaUltimaModificacion": "2025-04-29T00:00:00Z",
      "NumeroParaResumen": 0,
      "NroDiasMantenerCambioAMonedaLocal": 0,
      "FechaLimiteCambioAMonedaLocal": "2025-03-28T00:00:00Z",
      "FldTimeStamp": "0001-01-01T00:00:00Z",
      "GeneradoPor": "0",
      "BaseImponibleIGTF": 0,
      "IGTFML": 0,
      "IGTFME": 0,
      "AlicuotaIGTF": 3,
      "MotivoDeAnulacion": "",
      "ProveedorImprentaDigital": "1",
      "ConsecutivoVendedor": 1,
      "ImprentaDigitalGUID": ""
      },
    ...
  ],
  "hasNextPage": true,
  "hasPrevPage": false,
  "isFirstPage": true,
  "isLastPage": false,
  "page": 1,
  "pageSize": 1000,
  "total": 458591,
  "totalPages": 459
}
```