package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/arravoco/hackathon_backend/config"
	"github.com/arravoco/hackathon_backend/db"
	_ "github.com/arravoco/hackathon_backend/db"
	"github.com/arravoco/hackathon_backend/exports"
	"github.com/arravoco/hackathon_backend/jobs"
	"github.com/arravoco/hackathon_backend/rmqUtils"

	//"github.com/arravoco/hackathon_backend/jobs"

	_ "net/http/pprof"

	//_ "github.com/arravoco/hackathon_backend/nsq/consumer"
	routes_v1 "github.com/arravoco/hackathon_backend/routes/v1"
	"github.com/arravoco/hackathon_backend/security"
	"github.com/labstack/echo/v4"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
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

	e := echo.New()
	port := config.GetPort()
	routes_v1.StartAllRoutes(e)
	e.GET("/metrics", func(c echo.Context) error {
		handler := promhttp.Handler()
		handler.ServeHTTP(c.Response().Writer, c.Request())
		return nil
	})
	fmt.Println("Starting metrics")
	e.Logger.Info(port)

	db.SetupRedis()
	/*rmqUtils.SetupDefaultQueue()
	data.SetupDefaultDataSource()
	rabbitutils.SetupDefaultRMQ()
	rabbitutils.DeclareAllQueues()
	publish.SetPublisher(&rabbitutils.RMQPublisher{})
	go rabbitutils.ListenToAllQueues()*/
	//panic("Intentionally crashed")
	e.Logger.Fatal(e.Start(getURL(port)))
}

func getURL(port int) string {
	return strings.Join([]string{"", strconv.Itoa(port)}, ":")
}

func startAllJobs() {

	judgeCreatedByAdminWelcomeEmailTaskConsumer, err := jobs.StartConsumingJudgeCreatedByAdminWelcomeEmailQueue()
	if err != nil {
		fmt.Println(err.Error())
	}

	adminCreatedByAdminWelcomeEmailTaskConsumer, err := jobs.StartConsumingAdminCreatedByAdminWelcomeEmailQueue()
	if err != nil {
		fmt.Println(err.Error())
	}

	adminWelcomeEmailTaskConsumer, err := jobs.StartAdminWelcomeEmailQueue()
	if err != nil {
		fmt.Println(err.Error())
	}

	invitelistTaskConsumer, err := jobs.StartConsumingInviteTaskQueue()
	if err != nil {
		fmt.Println(err.Error())
	}

	for {
		select {
		case err := <-rmqUtils.ErrCh:
			fmt.Println(err.Error())
		case <-adminCreatedByAdminWelcomeEmailTaskConsumer.Ch:
			fmt.Println("'adminCreatedByAdminWelcomeEmailTaskConsumer' task completed successfully")

		case <-adminWelcomeEmailTaskConsumer.Ch:
			fmt.Println("'adminWelcomeEmailTaskConsumer' Task completed")
		case <-judgeCreatedByAdminWelcomeEmailTaskConsumer.Ch:
			fmt.Println("mail to judge created by admin list task completed successfully")
		case <-invitelistTaskConsumer.Ch:
			fmt.Println("Invite list task completed successfully")
		}
	}
}
