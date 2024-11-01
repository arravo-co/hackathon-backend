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
		Message: ` 
		<!DOCTYPE html>
		<html lang="en">
		<head>
		<meta charset="UTF-8">
		<meta name="viewport" content="width=device-width, initial-scale=1.0">
		<title>Arravo Hackathon Registration</title>
		</head>
		<body
		style="
		font-family: Google Sans, roboto, Noto Sans Myanmar UI, arial, sans-serif;
			margin: 0;
			padding: 0;
			background-color: #f8f9fa;">
		<div style="
		max-width: 600px;
		margin: 0 auto;
		background-color: #fff;
		border-radius: 8px;
		text-align: center;
		position: relative;">
		<h1 style="
    font-size: 18px;">Welcome aboard Arravo's 2024 Hackathon!</h1>
    <p style="
    text-align: center;
    color: #333;
    font-size: 1.3em;
    line-height: 1.5em;
    margin-top: 0;
    margin-bottom: 0;
    font-weight: 300;">You're one step away</p>
    <p style="
    text-align: justify;
    color: #333;
    font-size: 1.3em;
    line-height: 1.5em;
    margin-top: 0px;
    margin-bottom: 10px;
    font-weight: 300;
    padding: 20px">We're super pumped to have you leading your team through this exhilarating journey. <br> <br> Before we dive into the action, we just need to confirm your email address. Give that button a click below to get verified and set for success:</p>
    <h1 style="
    font-size: 18px;">Verify your email address:</h1>
    <div>
      <a href=https://hackathon-backend-2cvk.onrender.com/api/v1/auth/verification/email/completion?token&#61;yWttag2gfLBoKBoJfRPuQRMCMl1ki3De4FCPrJug5dB2TfX-p2szQR9Wf0Z5J33W-FUziZGxdZM3IxoQO5qw7dJnDzhGSqUJsXZ68glaJ9LEk0suwp17-q7PxNEe0hNB9Rub5ILZiVs8oyS_ZMzT3SLUgnzh71oMd7ZEEY4Rout2tRxj6TQpYutFd2DbhN_-6Dqslct8siV_A4I94srpV5ORQxZsCfnDk01ilgH2dVGUtIY5YFPE5z5Movo0A-ljeirqQF2KEAbC6HtfbhqFq0paiUun1dRYWRHxHBrTrc_5QbdJutMzq35QFsDTTRNxaRs9LfEcIuFhndcaF3CpcdCeGW5GAYdY7VjQpGRjZylve-oB2hBPYxuNUMEkMlg9cedVb8uqZa1XHOJypY7_js6VLbioRUEKVSHXyQTmDxK3zxxHmkaOM5jMsVwU95YVRDq-caW3cxSZCkZDBOs2WmRFUWDAqeffPngGGso_3cxz7WLuKNo2VOVLavZ2WzVzDkjfIbg7s2Gj6EnQg8f4sJO6bwN8QMJxFim7nN3PseNNeL4erdwO5f4O9BxTjAh33npuf5EZGAJDkiW6UB3b8ZsbsnRCQfk7eZtFVqr-Bi_UbsIBo5Tn2M2KMHQtjE_JMHRVpoA00-B7RATa6fNG9y26SObR0IhRSLY1LD1KJ3E style="
      display: inline-block;
      padding: 15px 45px;
      background-color: #553C9A;
      color: #fff;
      text-decoration: none;
      font-weight: 600;
      font-size: 13px;
      border-radius: 25px;
      box-shadow: 0 8px 16px rgba(0, 0, 0, 0.3);
      transition: box-shadow 0.3s ease;"
       onmouseover="this.style.backgroundColor='#805AD5'; this.style.boxShadow='0 12px 24px rgba(0, 0, 0, 0.4)'; this.style.cursor='pointer';"
       onmouseout="this.style.backgroundColor='#553C9A'; this.style.boxShadow='0 8px 16px rgba(0, 0, 0, 0.3)'; this.style.cursor='auto';"
      >Verify</a>
      <p style="
      font-size: 1.2em;
      line-height: 1.5em;
      margin-top: 0;
      margin-bottom: 10px;
      font-weight: 300;
      padding: 20px;">Alternatively, you can use the following token to verify your email: <br> Verification Token: <span style="font-weight: bold;">737752</span></p>
    </div>
</body>
</html>
		
		`,
	})
	if err != nil {
		t.Fatal(err.Error())
	}
}
