package email

import (
	"testing"

	"github.com/arravoco/hackathon_backend/testdbsetup"
)

func TestSendEmail(t *testing.T) {
	testdbsetup.SetupDefaultTestEnv()
	err := SendEmailHtml(&SendEmailHtmlData{
		Email:   "temitope.alabi@arravo.co",
		Subject: "Test Subject",
		Message: "Test Message",
	})
	if err != nil {
		t.Fatal(err.Error())
	}
}
