package main

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	clientes "github.com/desarrolladoresnet/api_galac_bd/cliente"
	"github.com/desarrolladoresnet/api_galac_bd/facturas"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/microsoft/go-mssqldb"
)

func main() {
	// Configuración de conexión
	connString := fmt.Sprintf("server=%s;user id=%s;password=%s;database=%s;TrustServerCertificate=true;",
		"WIN-754KG2T0CQI\\GALACSQLX17",
		"userdb",
		"123456",
		"SAWDB",
	)

	// Establecer conexión
	db, err := sql.Open("sqlserver", connString)
	if err != nil {
		log.Fatal("Error al abrir conexión:", err.Error())
	}
	defer db.Close()

	// Verificar conexión
	err = db.Ping()
	if err != nil {
		log.Fatal("Error al conectar:", err.Error())
	}
	fmt.Println("¡Conexión establecida correctamente!")

	// Inicializar Gin
	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"}, // ✅ solo tu app de Vite
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization", "X-API-Key"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Ruta básica
	router.GET("/ping", func(c *gin.Context) {
		c.String(200, "Hello World")
	})

	api_facturas := router.Group("facturas")
	facturas.Facturas(api_facturas, db)

	api_clientes := router.Group("clientes")
	clientes.ClienteRoutes(api_clientes, db)
	// Iniciar servidor
	router.Run(":5000")
}
