package utils

import (
	"context"
	"os"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/messaging"
	"google.golang.org/api/option"
)

var client *messaging.Client

func InitNotifier() error {
	if client != nil {
		return nil
	}
	ctx := context.Background()
	opt := option.WithCredentialsFile(os.Getenv("FIREBASE_PATH"))
	app, err := firebase.NewApp(ctx, nil, opt)

	if err != nil {
		return err
	}

	client, err = app.Messaging(ctx)
	if err != nil {
		return err
	}

	return nil
}

func SendMessage(token string, plantName string) {
	if client == nil {
		InitNotifier()
	}

	ctx := context.Background()

	message := &messaging.Message{
		Token: token,
		Data: map[string]string{
			"title": "Hello from Plantie!",
			"body":  "Time to water your plant " + plantName,
		},
	}
	client.Send(ctx, message)
}
