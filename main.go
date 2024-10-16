package main

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/arravoco/hackathon_backend/config"
	consumerhandlers "github.com/arravoco/hackathon_backend/consumer_handlers"
	"github.com/arravoco/hackathon_backend/consumers"
	_ "github.com/arravoco/hackathon_backend/db"
	"github.com/arravoco/hackathon_backend/exports"
	"github.com/arravoco/hackathon_backend/instru"
	"github.com/arravoco/hackathon_backend/publishers"
	"github.com/arravoco/hackathon_backend/resources"

	//"github.com/arravoco/hackathon_backend/jobs"

	_ "net/http/pprof"

	//_ "github.com/arravoco/hackathon_backend/nsq/consumer"
	routes_v1 "github.com/arravoco/hackathon_backend/routes/v1"
	"github.com/arravoco/hackathon_backend/security"
	"github.com/labstack/echo/v4"

	//_ "github.com/newrelic/go-agent/v3/newrelic"
	"github.com/prometheus/client_golang/prometheus"
)

// @Version 1.0.0
// @Title Hackathon Backend API
// @Description API usually works as expected. But sometimes its not true.
// @ContactName David Alabi
// @ContactEmail appdev@arravo.co
// @ContactURL http://arravo.co/contact
// @TermsOfServiceUrl http://arravo.co/contact
// @LicenseName MIT
// @LicenseURL https://en.wikipedia.org/wiki/MIT_License
// @Server http://localhost:5000 Localhost
// @Server https://hackathon-backend-2cvk.onrender.com Development
func main() {
	prometheus.MustRegister(exports.MyFirstCounter)
	security.GenerateKeys()

	//txn := app.StartTransaction("")
	//txn.
	// http.HandleFunc(newrelic.WrapHandleFunc(app, "/users", usersHandler))
	/*rmqUtils.SetupDefaultQueue()
	data.SetupDefaultDataSource()
	rabbitutils.SetupDefaultRMQ()
	rabbitutils.DeclareAllQueues()
	publish.SetPublisher(&rabbitutils.RMQPublisher{})*/
	//panic("Intentionally crashed")
	e := echo.New()
	port := config.GetPort()
	routes_v1.StartAllRoutes(e)
	/*
		e.GET("/metrics", func(c echo.Context) error {
			handler := promhttp.Handler()
			handler.ServeHTTP(c.Response().Writer, c.Request())
			return nil
		})
		fmt.Println("Starting metrics")
	*/
	res := resources.GetDefaultResources()

	fn, err := instru.Setup(context.Background(), &instru.SetupOtel{
		LoggerProvider: res.LoggerProvider,
		MeterProvider:  res.MeterProvider,
		TracerProvider: res.TraceProvider,
	})

	if err != nil {
		e.Logger.Error(err)
	}
	defer func() {
		fn(context.Background())
	}()

	publishChannel, err := res.RabbitMQConn.Channel()
	if err != nil {
		res.Logger.Fatal(err.Error())
	}
	/*
		rmqPushlisher := publishers.NewRMQPublisherWithChannel(publishChannel)
		rmqPushlisher.DeclareAllExchanges()
		res.Logger.Sugar().Infoln(rmqPushlisher)
	*/

	publishers.DeclareAllExchanges(publishChannel)

	consumerChannel, err := res.RabbitMQConn.Channel()
	if err != nil {
		res.Logger.Fatal(err.Error())
	}

	time.Sleep(time.Second * 10)
	rmqConsumer := consumers.NewRMQConsumerWithChannel(consumers.CreateRMQConsumerOpts{
		Channel:       consumerChannel,
		Logger:        res.Logger,
		RMQConnection: res.RabbitMQConn,
	})

	go func() {
		err = rmqConsumer.DeclareAllQueuesParameterized(
			consumerhandlers.SendWelcomeAndEmailVerificationEmailToAdmin,
			exports.AdminsExchange,
			exports.SendAdminWelcomeEmailQueueName,
			exports.AdminSendWelcomeEmailBindingKeyName)
		if err != nil {
			fmt.Println(err)
		}
	}()

	go func() {
		err = rmqConsumer.DeclareAllQueuesParameterized(
			consumerhandlers.SendWelcomeAndEmailVerificationEmailToAdminRegisteredByAdmin,
			exports.AdminsExchange,
			exports.SendAdminRegisteredByAdminWelcomeEmailQueueName,
			exports.AdminRegisteredByAdminSendWelcomeEmailBindingKeyName)

		if err != nil {
			fmt.Println(err)
		}
	}()

	go func() {
		err = rmqConsumer.DeclareAllQueuesParameterized(
			consumerhandlers.SendWelcomeAndEmailVerificationEmailToJudgeRegisteredByAdmin,
			exports.JudgesExchange,
			exports.SendJudgeWelcomeEmailQueueName,
			exports.SendJudgeWelcomeEmailQueueBindingKeyName)

		if err != nil {
			fmt.Println(err)
		}
	}()

	go func() {
		err = rmqConsumer.DeclareAllQueuesParameterized(
			consumerhandlers.SendWelcomeAndEmailVerificationEmailToJudge,
			exports.JudgesExchange,
			exports.SendJudgeRegisteredByAdminWelcomeEmailQueueName,
			exports.JudgeRegisteredByAdminSendWelcomeEmailBindingKeyName)

		if err != nil {
			fmt.Println(err)
		}
	}()

	go func() {
		err = rmqConsumer.DeclareAllQueuesParameterized(
			consumerhandlers.SendTeamLeadWelcomeAndVerificationEmail,
			exports.ParticipantsExchange,
			exports.SendTeamLeadWelcomeEmailQueueName,
			exports.SendTeamLeadWelcomeEmailQueueBindingKeyName)
		if err != nil {
			fmt.Println(err)
		}
	}()

	go func() {
		err = rmqConsumer.DeclareAllQueuesParameterized(
			consumerhandlers.SendTeamMemberWelcomeAndVerificationEmail,
			exports.ParticipantsExchange,
			exports.SendTeamMemberWelcomeEmailQueueName,
			exports.ParticipantTeamMemberSendWelcomeEmailRoutingKeyName)
		if err != nil {
			fmt.Println(err)
		}
	}()
	go func() {
		rmqConsumer.DeclareAllQueuesParameterized(
			consumerhandlers.SendInviteEmailQueueHandler,
			exports.InvitationsExchange,
			exports.SendParticipantTeammateInvitationEmailQueueName,
			exports.ParticipantTeammateSendInvitationEmailBindingKeyName)
		if err != nil {
			fmt.Println(err)
		}
	}()
	go func() {
		rmqConsumer.DeclareAllQueuesParameterized(
			consumerhandlers.HandleUploadJudgeProfilePicConsumption,
			exports.UploadJobsExchange,
			exports.UploadJudgeProfilePicQueueName,
			exports.UploadJudgeProfilePicBindingKeyName)
		if err != nil {
			fmt.Println(err)
		}
	}()

	go func() {
		rmqConsumer.DeclareAllQueuesParameterized(
			consumerhandlers.HandleUploadJudgeProfilePicConsumption,
			exports.UploadJobsExchange,
			exports.UploadJudgeProfilePicQueueName,
			exports.UploadJudgeProfilePicBindingKeyName)
		if err != nil {
			fmt.Println(err)
		}
	}()
	//AdminRegisteredByAdminSendWelcomeEmailRoutingKeyName

	e.Logger.Fatal(e.Start(getURL(port)))
}

func getURL(port int) string {
	return strings.Join([]string{"", strconv.Itoa(port)}, ":")
}
