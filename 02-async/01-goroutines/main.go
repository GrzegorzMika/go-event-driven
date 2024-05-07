package main

import (
	"fmt"
	"time"
)

type User struct {
	Email string
}

type UserRepository interface {
	CreateUserAccount(u User) error
}

type NotificationsClient interface {
	SendNotification(u User) error
}

type NewsletterClient interface {
	AddToNewsletter(u User) error
}

type Handler struct {
	repository          UserRepository
	newsletterClient    NewsletterClient
	notificationsClient NotificationsClient
}

func NewHandler(
	repository UserRepository,
	newsletterClient NewsletterClient,
	notificationsClient NotificationsClient,
) Handler {
	return Handler{
		repository:          repository,
		newsletterClient:    newsletterClient,
		notificationsClient: notificationsClient,
	}
}

func retryWrapper(user User, fn func(user User) error) {
	for {
		err := fn(user)
		if err == nil {
			return
		}
		fmt.Printf("failed to retry: %s\n", err)
		time.Sleep(time.Second)
	}
}

func (h Handler) SignUp(u User) error {
	if err := h.repository.CreateUserAccount(u); err != nil {
		return err
	}

	go retryWrapper(u, h.newsletterClient.AddToNewsletter)
	go retryWrapper(u, h.notificationsClient.SendNotification)
	return nil
}
