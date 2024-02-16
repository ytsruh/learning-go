package data

import (
	"ytsruh.com/saas/data/model"
)

type DB struct {
	DatabaseName string
	Connection   *model.Connection
	CopySession  bool

	Users    UserServices
	Webhooks WebhookServices
}

type SessionRefresher interface {
	RefreshSession(*model.Connection, string)
}

type SessionCloser interface {
	Close()
}

type UserServices interface {
	SessionRefresher
	SessionCloser
	SignUp(email, password string) (*model.Account, error)
	AddToken(accountID, userID model.Key, name string) (*model.AccessToken, error)
	RemoveToken(accountID, userID, tokenID model.Key) error
	Auth(accountID, token string, pat bool) (*model.Account, *model.User, error)
	GetDetail(id model.Key) (*model.Account, error)
	GetByStripe(stripeId string) (*model.Account, error)
	SetSeats(id model.Key, seats int) error
	ConvertToPaid(id model.Key, stripeID, subID, plan string, yearly bool, seats int) error
	ChangePlan(id model.Key, plan string, yearly bool) error
	Cancel(id model.Key) error
}

type AdminServices interface {
	SessionRefresher
	SessionCloser
	LogRequests(reqs []model.APIRequest) error
}

type WebhookServices interface {
	SessionRefresher
	SessionCloser
	Add(accountID model.Key, events, url string) error
	List(accountID model.Key) ([]model.Webhook, error)
	Delete(accountID model.Key, event, url string) error
	AllSubscriptions(event string) ([]model.Webhook, error)
}
