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
	"github.com/wowucco/G3/internal/checkout"
	checkoutHttp "github.com/wowucco/G3/internal/checkout/delivery/http"
	"github.com/wowucco/G3/internal/checkout/repository"
	"github.com/wowucco/G3/internal/checkout/strategy"
	"github.com/wowucco/G3/internal/checkout/usecase"
	"github.com/wowucco/G3/internal/delivery"
	_deliveryRepo "github.com/wowucco/G3/internal/delivery/repository"
	"github.com/wowucco/G3/internal/menu"
	_menuRepo "github.com/wowucco/G3/internal/menu/repository/psql"
	"github.com/wowucco/G3/internal/product"
	productHttp "github.com/wowucco/G3/internal/product/delivery/http"
	_productRepo "github.com/wowucco/G3/internal/product/repository/psql"
	productUC "github.com/wowucco/G3/internal/product/usecase"
	"github.com/wowucco/G3/pkg/gqlgen/graph"
	"github.com/wowucco/G3/pkg/http/middleware"
	"github.com/wowucco/G3/pkg/notification"
	"github.com/wowucco/G3/pkg/payments/liqpay"
	"github.com/wowucco/G3/pkg/payments/privatPay"
	"github.com/wowucco/G3/pkg/sms"
	smsMock "github.com/wowucco/G3/pkg/sms/mock"
	smsClub "github.com/wowucco/G3/pkg/sms/smsclub"
	telegram2 "github.com/wowucco/G3/pkg/telegram"
	"github.com/wowucco/go-novaposhta"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

type App struct {
	httpServer *http.Server

	productUC    product.UseCase
	productRead  product.ReadRepository
	menuRead     menu.ReadRepository
	deliveryRead delivery.DeliveryReadRepository

	orderManage checkout.IOrderUseCase

	db *dbx.DB
	es *elasticsearch.Client

	sms          sms.Client
	smsChan      chan sms.Message
	telegram     telegram2.Client
	telegramChan chan telegram2.Message
}

func NewApp() *App {
	db := initDB()
	es := initElasticsearch()
	np := initNovaposhtaClient()

	productRepo := _productRepo.NewProductRepository(db)
	productRead := _productRepo.NewProductReadRepository(db, es)
	deliveryRead := _deliveryRepo.NewDeliveryReadRepository(db, es, np)

	smsChan := make(chan sms.Message, 1)
	telegramChan := make(chan telegram2.Message, 1)

	return &App{
		db: db,
		es: es,

		productUC:    productUC.NewProductUseCase(productRepo),
		productRead:  productRead,
		menuRead:     _menuRepo.NewMenuReadRepository(db),
		deliveryRead: deliveryRead,

		orderManage: usecase.NewOrderUseCase(
			repository.NewOrderRepository(db),
			productRead,
			deliveryRead,
			repository.NewPaymentRepository(db),
			initNotificationService(smsChan, telegramChan),
			initPaymentContext(db),
		),

		smsChan:      smsChan,
		telegramChan: telegramChan,
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

	platformAuth := middleware.TokenAuthMiddleware(viper.GetString("auth.api_id"), viper.GetString("auth.api_code"))

	productHttp.RegisterHTTPEndpoints(api, app.productUC, platformAuth)
	checkoutHttp.RegisterHTTPEndpoints(api, app.orderManage, platformAuth)
	graph.RegisterGraphql(api, app.productUC, app.productRead, app.menuRead, app.deliveryRead)

	app.httpServer = &http.Server{
		Addr:    fmt.Sprintf(":%v", port),
		Handler: router,
	}

	initSmsListening(app.smsChan)
	initTelegramListening(app.telegramChan)

	go func() {
		if err := app.httpServer.ListenAndServe(); err != nil {
			log.Fatalf("Failed to listen and serve: %+v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, os.Interrupt)

	<-quit

	ctx, shutdown := context.WithTimeout(context.Background(), 5*time.Second)

	defer func() {
		if err := app.db.Close(); err != nil {
			log.Fatalf("Failed to close DB connection: %+v", err)
		}
	}()

	defer shutdown()

	return app.httpServer.Shutdown(ctx)
}

func initDB() *dbx.DB {
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

func initNovaposhtaClient() *novaposhta.Client {
	cfg := novaposhta.Config{
		ApiKey: viper.GetString("novaposhta.api_key"),
	}

	np, err := novaposhta.NewClient(cfg)

	if err != nil {
		log.Fatalf("Error creating the novaposhta client: %s", err)
	}

	return np
}

func initSmsListening(ch <-chan sms.Message) {

	var c sms.Client

	switch viper.GetString("sms.provider") {
	case "smsclub":
		c = smsClub.NewClient(smsClub.Config{
			Token: viper.GetString("sms.smsclub_token"),
			From:  viper.GetString("sms.smsclub_alfaname"),
		})
	default:
		c = smsMock.NewClient()
	}

	go func(cl sms.Client, ch <-chan sms.Message) {
		for {
			msg := <-ch
			_, _ = cl.Send(msg)
		}
	}(c, ch)
}

func initTelegramListening(ch <-chan telegram2.Message) {

	var cl telegram2.Client

	switch viper.GetString("telegram.mode") {
	case "real":
		c, e := telegram2.NewTelegramClient(telegram2.Config{
			ApiUrl: viper.GetString("telegram.api_url"),
			BotId:  viper.GetString("telegram.bot_id"),
		})

		if e != nil {
			panic(e)
		}

		cl = c
	default:
		c, _ := telegram2.NewMock()
		cl = c
	}

	go func(cl telegram2.Client, ch <-chan telegram2.Message) {
		for {
			msg := <-ch
			_, _ = cl.Send(msg)
		}
	}(cl, ch)
}
func initNotificationService(smsChan chan sms.Message, telegramChan chan telegram2.Message) *notification.Service {

	telegramChats := map[string]string{
		notification.TelegramOrderChat:  viper.GetString("telegram.order_chat_id"),
		notification.TelegramRecallChat: viper.GetString("telegram.recall_chat_id"),
	}

	return notification.NewNotificationService(
		smsChan,
		telegramChan,
		telegramChats,
		viper.GetString("payments.card_number"),
		viper.GetString("domain.bo_order_link_mask"),
		viper.GetString("domain.web_product_link_mask"),
	)
}

func initPaymentContext(db *dbx.DB) *strategy.PaymentContext {

	r := repository.NewPaymentRepository(db)

	lp := liqpay.NewClient(
		viper.GetString("payments.liqpay.public"),
		viper.GetString("payments.liqpay.private"),
		viper.GetString("payments.liqpay.callback_url"),
		viper.GetString("payments.liqpay.return_url"),
	)

	ppp := privatPay.NewClient(privatPay.Config{
		StoreId:     viper.GetString("payments.privat_pay.store_id"),
		Password:    viper.GetString("payments.privat_pay.passport"),
		Min:         viper.GetInt("payments.privat_pay.min_parts"),
		Max:         viper.GetInt("payments.privat_pay.max_parts"),
		ResponseUrl: viper.GetString("payments.privat_pay.callback_url"),
		RedirectUrl: viper.GetString("payments.privat_pay.return_url"),
	})

	p2p := strategy.NewP2PStrategy(r, lp)
	pp := strategy.NewPartsPayStrategy(r, ppp)
	d := strategy.NewDefaultStrategy(r)

	return strategy.NewPaymentContext(p2p, pp, d)
}
