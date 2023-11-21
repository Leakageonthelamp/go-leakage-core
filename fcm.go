package core

import (
	"context"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/messaging"
	"google.golang.org/api/option"
)

type IFMC interface {
	Send(data map[string]string, token string) error
	SendAll(data map[string]string, tokens []string) error
	SendTopic(data map[string]string, topic string) error
}

type fmcService struct {
	ctx    IContext
	Client *messaging.Client
}

func NewFMCService(ctx IContext, bgCtx context.Context, credentialJSON []byte) (IFMC, error) {
	opt := option.WithCredentialsJSON(credentialJSON)
	app, err := firebase.NewApp(bgCtx, nil, opt)
	if err != nil {
		return nil, err
	}

	client, err := app.Messaging(bgCtx)
	if err != nil {
		return nil, err
	}

	return &fmcService{
		ctx:    ctx,
		Client: client,
	}, nil
}

func (s *fmcService) Send(data map[string]string, token string) error {
	message := &messaging.Message{
		Data:  data,
		Token: token,
	}

	_, err := s.Client.Send(context.Background(), message)
	if err != nil {
		return err
	}

	return nil
}

func (s *fmcService) SendAll(data map[string]string, tokens []string) error {
	message := &messaging.MulticastMessage{
		Data:   data,
		Tokens: tokens,
	}

	_, err := s.Client.SendMulticast(context.Background(), message)
	if err != nil {
		return err
	}

	return nil
}

func (s *fmcService) SendTopic(data map[string]string, topic string) error {
	message := &messaging.Message{
		Data:  data,
		Topic: topic,
	}

	_, err := s.Client.Send(context.Background(), message)
	if err != nil {
		return err
	}

	return nil
}
