package resources

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/arravoco/hackathon_backend/db"
	"github.com/arravoco/hackathon_backend/rabbitmq"
	"github.com/labstack/echo/v4"
	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/rabbitmq/amqp091-go"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.opentelemetry.io/otel/sdk/log"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/trace"
	"go.uber.org/zap"
)

type SetupOpts struct {
	Logger echo.Logger
}

type AppResources struct {
	RedisClient    *redis.Client
	TraceProvider  *trace.TracerProvider
	MeterProvider  *metric.MeterProvider
	LoggerProvider *log.LoggerProvider
	Logger         *zap.Logger
	RabbitMQConn   *amqp091.Connection
	Mongo          *mongo.Database
	RelicApp       *newrelic.Application
}

var defaultRMQConn *amqp091.Connection
var defaultMongoInstance *mongo.Client
var defaultResources *AppResources
var defaultTraceProvider *trace.TracerProvider
var defaultMeterProvider *metric.MeterProvider
var defaultLoggerProvider *log.LoggerProvider

type RabbitMQResource struct {
	Url     string
	Conn    *amqp091.Connection
	ErrChan chan *amqp091.Error
}

var defaultRabbitMQStruct *RabbitMQResource //= GetDefaultRabbitMQResource()

func GetDefaultRabbitMQResource() *RabbitMQResource {
	var errChan chan *amqp091.Error = make(chan *amqp091.Error)
	var url = os.Getenv("RABBITMQ_URL")
	conn, err := rabbitmq.GetRMQConnWithURL(rabbitmq.SetupRMQConfig{
		Url: url,
	})
	if err != nil {
		panic(err)
	}
	s := &RabbitMQResource{
		Url:     url,
		ErrChan: errChan,
		Conn:    conn,
	}
	go s.RestartDefaultRabbit()
	s.Conn.NotifyClose(errChan)
	return s
}

func InitializeDefaultResources() {
	var err error
	logger := GetDefaultLogger()
	rabbitMQURL := os.Getenv("RABBITMQ_URL")
	if rabbitMQURL == "" {
		logger.Fatal("Please specify rabbitMQ URL")
	}
	fmt.Println("Check if RabbitMQ connection exists")
	if defaultRabbitMQStruct == nil {
		defaultRabbitMQStruct = GetDefaultRabbitMQResource()
	} else {
		if defaultRabbitMQStruct.Conn.IsClosed() {
			fmt.Println("New RabbitMQ connection to be established")
			defaultRabbitMQStruct = GetDefaultRabbitMQResource()
		}
	}
	var dbInstance *mongo.Database
	if defaultMongoInstance == nil {
		//defaultMongoInstance.Connect()
	}
	db_url := os.Getenv("MONGODB_URL")
	clientOpts := options.Client().ApplyURI(db_url)
	defaultMongoInstance, err = mongo.Connect(context.Background(), clientOpts)
	if err != nil {
		panic(err)
	}
	dbInstance = defaultMongoInstance.Database("hackathons_db")

	NEW_RELIC_LICENSE_KEY := os.Getenv("NEW_RELIC_LICENSE_KEY")
	//rabbitMQURL := os.Getenv("NEW_RELIC_USER_KEY")
	app, err := newrelic.NewApplication(
		newrelic.ConfigAppName("Hackathon Backend"),
		newrelic.ConfigLicense(NEW_RELIC_LICENSE_KEY),
		newrelic.ConfigAppLogForwardingEnabled(true),
	)
	if err != nil {
		panic(err)
	}
	redisClient := db.NewRedisDefaultClient()

	//instru.Setup(context.Background(), &instru.SetupOtel{})
	/*
		defaultTraceProvider, err = instru.NewTraceProvider()
		if err != nil {
			panic(err)
		}

		defaultMeterProvider, err = instru.NewMeterProvider()
		if err != nil {
			panic(err)
		}
		defaultLoggerProvider, err = instru.NewLoggerProvider()
		if err != nil {
			panic(err)
		}
	*/
	defaultResources = &AppResources{
		RedisClient:  redisClient,
		Logger:       logger,
		RabbitMQConn: defaultRabbitMQStruct.Conn,
		Mongo:        dbInstance,
		RelicApp:     app,
		/*
			TraceProvider:  defaultTraceProvider,
			MeterProvider:  defaultMeterProvider,
			LoggerProvider: defaultLoggerProvider,
		*/
	}
}

func GetDefaultResources() *AppResources {
	if defaultResources != nil {
		fmt.Println("Resources default")
		return defaultResources
	}
	InitializeDefaultResources()
	fmt.Println("Resources fetched")
	return defaultResources
}

func GetDefaultLogger() *zap.Logger {

	logger, _ := zap.NewProduction()
	return logger
}

func (appRes *AppResources) CheckHealth() {
	//err:=appRes.Mongo.Client().Ping(context.Background(),&readpref.ReadPref{})
}

func (s *RabbitMQResource) RestartDefaultRabbit() {
	for {
		select {
		case err := <-s.ErrChan:
			fmt.Printf("RabbitMQ reason: %v\n", err.Reason)
			fmt.Printf("RabbitMQ error: %v\n", err.Error())
			fmt.Printf("can recover: %v\n", err.Recover)
			if s.Url != "" {
				for {
					conn, err := rabbitmq.GetRMQConnWithURL(rabbitmq.SetupRMQConfig{
						Url: s.Url,
					})
					if err != nil {
						fmt.Print("failed to get RMQ connection")
						time.Sleep(time.Second * 5)

						continue
					}
					var errChan chan *amqp091.Error = make(chan *amqp091.Error)
					s.Conn = conn
					s.Conn.NotifyClose(errChan)
					return
				}
			}
		}
	}
}

func (appRes *AppResources) Shutdown() {
	defaultRMQConn.Close()
	defaultMongoInstance.Disconnect(context.Background())
	defaultResources.RelicApp.Shutdown(time.Minute)
}
