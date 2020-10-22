package server

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/elastic/go-elasticsearch/v5"
	"github.com/gin-gonic/gin"
	dbx "github.com/go-ozzo/ozzo-dbx"
	_ "github.com/lib/pq"
	"github.com/spf13/viper"
	"github.com/wowucco/G3/internal/product"
	producttHttp "github.com/wowucco/G3/internal/product/delivery/http"
	"github.com/wowucco/G3/internal/product/repository/psql"
	productUC "github.com/wowucco/G3/internal/product/usecase"
	"github.com/wowucco/G3/pkg/gqlgen/graph"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

type App struct {
	httpServer *http.Server

	productUC product.UseCase
	productRead product.ReadRepository

	db *dbx.DB
	es *elasticsearch.Client
}

func NewApp() *App {
	db := initDB()
	es := initElasticsearch()
	productRepo := psql.NewProductRepository(db)

	return &App{
		db: db,
		es: es,

		productUC: productUC.NewProductUseCase(productRepo),
		productRead: psql.NewProductReadRepository(db, es),
	}
}

func (app *App) Run(port string) error {
	router := gin.Default()
	router.Use(
		gin.Recovery(),
		gin.Logger(),
	)

	router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	api := router.Group("/api")

	producttHttp.RegisterHTTPEndpoints(api, app.productUC)
	graph.RegisterGraphql(api, app.productUC, app.productRead)

	app.httpServer = &http.Server{
		Addr:           fmt.Sprintf(":%v", port),
		Handler:        router,
	}

	go func() {
		if err := app.httpServer.ListenAndServe(); err != nil {
			log.Fatalf("Failed to listen and serve: %+v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, os.Interrupt)

	<-quit

	ctx, shutdown := context.WithTimeout(context.Background(), 5 * time.Second)

	defer func() {
		if err := app.db.Close(); err != nil {
			log.Fatalf("Failed to close DB connection: %+v", err)
		}
	}()

	defer shutdown()

	return app.httpServer.Shutdown(ctx)
}

func initDB()  *dbx.DB {
	// connect to the database
	db, err := dbx.MustOpen("postgres", viper.GetString("db_dns"))
	if err != nil {
		log.Fatalf("Failed to connect to DB: %+v", err)
		os.Exit(-1)
	}

	db.QueryLogFunc = logDBQuery

	return db
}

func logDBQuery(ctx context.Context, t time.Duration, sql string, rows *sql.Rows, err error) {
	fmt.Print("\n>>>\n query: ", sql, "\n<<<\n")
	fmt.Print("\n>>>\n duration: ", t, "\n<<<\n")
	if err != nil {
		fmt.Print("\n>>>\n error: ", err.Error(), "\n<<<\n")
	}
}

func initElasticsearch() *elasticsearch.Client {
	cfg := elasticsearch.Config{
		Addresses: []string{
			viper.GetString("elasticsearch.host"),
		},
	}

	es, err := elasticsearch.NewClient(cfg)

	if err != nil {
		log.Fatalf("Error creating the client: %s", err)
	}

	return es
}
