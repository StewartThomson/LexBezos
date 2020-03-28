package main

import (
	"database/sql"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/dghubble/oauth1"
	_ "github.com/go-sql-driver/mysql"
	"go.uber.org/zap"
	"os"
)

const EnvDBHost = "LB_DB_HOST"
const EnvDBUser = "LB_DB_USER"
const EnvDBPass = "LB_DB_PASS"

const EnvConsumerKey = "LB_CONSUMER_KEY"
const EnvConsumerSecret = "LB_CONSUMER_SECRET"
const EnvAccessToken = "LB_ACCESS_TOKEN"
const EnvAccessSecret = "LB_ACCESS_SECRET"

type input struct {
	Debug bool `json:"debug"`
}

var logger *zap.SugaredLogger

func handler(request input) error {
	var err error
	var l *zap.Logger
	if request.Debug == true {
		l, err = zap.NewDevelopment()
	} else {
		l, err = zap.NewProduction()
	}
	if err != nil {
		return err
	}
	logger = l.Sugar()
	defer logger.Sync()

	logger.Debug("Starting function")

	db, err := sql.Open("mysql", os.Getenv(EnvDBUser)+":"+os.Getenv(EnvDBPass)+"@tcp("+os.Getenv(EnvDBHost)+")/")
	if err != nil {
		return err
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		return err
	}

	id, content, err := GetTweet(db)
	if err != nil {
		return err
	}

	if id == 0 {
		logger.Info("No new tweets to send. Ending.")
		return nil
	}

	config := oauth1.NewConfig(os.Getenv(EnvConsumerKey), os.Getenv(EnvConsumerSecret))
	token := oauth1.NewToken(os.Getenv(EnvAccessToken), os.Getenv(EnvAccessSecret))

	httpClient := config.Client(oauth1.NoContext, token)

	twitID, err := SendTweet(httpClient, content)
	if err != nil {
		return err
	}

	err = UpdateTwitInfo(db, id, twitID)
	if err != nil {
		return err
	}
	return nil
}

func main() {
	lambda.Start(handler)
}
