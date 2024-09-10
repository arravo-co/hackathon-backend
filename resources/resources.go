package resources

import (
	"context"
	"fmt"
	"os"

	"github.com/arravoco/hackathon_backend/db"
	"github.com/arravoco/hackathon_backend/rabbitmq"
	"github.com/labstack/echo/v4"
	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/rabbitmq/amqp091-go"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

type SetupOpts struct {
	Logger echo.Logger
}

type AppResources struct {
	RedisClient *redis.Client
	//Publisher   *publishers.RMQPublisher
	//Consumer    *consumers.RMQConsumer
	Logger       *zap.Logger
	RabbitMQConn *amqp091.Connection
	Mongo        *mongo.Database
	RelicApp     *newrelic.Application
}

var defaultRMQConn *amqp091.Connection
var defaultMongoInstance *mongo.Client
var defaultResources *AppResources

func InitializeDefaultResources() {
	var err error
	logger := GetDefaultLogger()
	rabbitMQURL := os.Getenv("RABBITMQ_URL")
	if rabbitMQURL == "" {
		logger.Fatal("Please specify rabbitMQ URL")
	}
	fmt.Println("Check if RabbitMQ connection exists")
	if defaultRMQConn == nil || defaultRMQConn.IsClosed() {
		fmt.Println("New RabbitMQ connection to be established")
		defaultRMQConn, err = rabbitmq.GetRMQConnWithURL(rabbitmq.SetupRMQConfig{
			Url: rabbitMQURL,
		})
		if err != nil {
			logger.Fatal(err.Error())
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
	defaultResources = &AppResources{
		RedisClient:  redisClient,
		Logger:       logger,
		RabbitMQConn: defaultRMQConn,
		Mongo:        dbInstance,
		RelicApp:     app,
	}
}

func GetDefaultResources() *AppResources {
	if defaultResources != nil {
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
