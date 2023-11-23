package core

import (
	"context"
	"net/http"
	"time"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/messaging"
	"google.golang.org/api/option"
)

var FCMServiceConnectError = Error{
	Status:  http.StatusBadGateway,
	Code:    "FCM_CONNECT_ERROR",
	Message: "failure to connect to fcm service",
}

var FCMServiceSendError = Error{
	Status:  http.StatusBadGateway,
	Code:    "FCM_SEND_ERROR",
	Message: "failure to send notification",
}

type IFMCMessage struct {
	Title    string `json:"title"`
	Body     string `json:"body"`
	ImageURL string `json:"image_url"`

	Data map[string]string `json:"data"`
}

type IFMCPayload struct {
	Token   string       `json:"token"`
	Message *IFMCMessage `json:"payload"`
}

type IFMC interface {
	SendSimpleMessage(tokens []string, payload *IFMCMessage) IError
	SendSimpleMessages(payload []IFMCPayload) IError
	SendTopic(topic string, payload map[string]string) IError
}

type fmcService struct {
	ctx    IContext
	client *messaging.Client
}

func NewFMC(ctx IContext) IFMC {
	return &fmcService{
		ctx: ctx,
	}
}

func (s *fmcService) connect() IError {
	if s.client != nil {
		return nil
	}

	credential := s.ctx.ENV().Config().FirebaseCredential

	opt := option.WithCredentialsJSON([]byte(credential))
	ctx, cancel := s.getContext()
	defer cancel()

	app, err := firebase.NewApp(ctx, nil, opt)
	if err != nil {
		return s.ctx.NewError(err, FCMServiceConnectError)
	}

	client, err := app.Messaging(ctx)
	if err != nil {
		return s.ctx.NewError(err, FCMServiceConnectError)
	}

	s.client = client
	return nil
}

func (s *fmcService) SendSimpleMessage(tokens []string, message *IFMCMessage) IError {
	ierr := s.connect()
	if ierr != nil {
		return s.ctx.NewError(ierr, ierr)
	}

	_message := &messaging.MulticastMessage{
		Notification: &messaging.Notification{
			Title:    message.Title,
			Body:     message.Body,
			ImageURL: message.ImageURL,
		},
		Data:   message.Data,
		Tokens: tokens,
	}

	ctx, cancel := s.getContext()
	defer cancel()

	_, err := s.client.SendMulticast(ctx, _message)
	if err != nil {
		return s.ctx.NewError(err, FCMServiceSendError)
	}

	return nil
}

func (s *fmcService) SendSimpleMessages(payload []IFMCPayload) IError {
	ierr := s.connect()
	if ierr != nil {
		return s.ctx.NewError(ierr, ierr)
	}

	messages := make([]*messaging.Message, 0)
	for _, p := range payload {
		messages = append(messages, &messaging.Message{
			Notification: &messaging.Notification{
				Title:    p.Message.Title,
				Body:     p.Message.Body,
				ImageURL: p.Message.ImageURL,
			},
			Data:  p.Message.Data,
			Token: p.Token,
		})

	}

	ctx, cancel := s.getContext()
	defer cancel()

	_, err := s.client.SendAll(ctx, messages)
	if err != nil {
		return s.ctx.NewError(err, FCMServiceSendError)
	}

	return nil
}

func (s *fmcService) SendTopic(topic string, payload map[string]string) IError {
	ierr := s.connect()
	if ierr != nil {
		return s.ctx.NewError(ierr, ierr)
	}

	message := &messaging.Message{
		Data:  payload,
		Topic: topic,
	}

	ctx, cancel := s.getContext()
	defer cancel()

	_, err := s.client.Send(ctx, message)
	if err != nil {
		return s.ctx.NewError(err, FCMServiceSendError)
	}

	return nil
}

func (s *fmcService) getContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), 10*time.Second)
}
