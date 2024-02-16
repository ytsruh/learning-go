package controllers

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"ytsruh.com/saas/data"
	"ytsruh.com/saas/data/model"
	"ytsruh.com/saas/engine"
	"ytsruh.com/saas/queue"

	"github.com/gorilla/mux"
	"github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/card"
	"github.com/stripe/stripe-go/customer"
	"github.com/stripe/stripe-go/invoice"
	"github.com/stripe/stripe-go/sub"
)

func init() {
	stripe.Key = os.Getenv("STRIPE_KEY")
}

// Billing handles everything related to the billing requests
type Billing struct{}

// BillingOverview represents if an account is a paid customer or not
type BillingOverview struct {
	StripeID       string            `json:"stripeId"`
	Plan           string            `json:"plan"`
	IsYearly       bool              `json:"isYearly"`
	IsNew          bool              `json:"isNew"`
	Cards          []BillingCardData `json:"cards"`
	CostForNewUser int               `json:"costForNewUser"`
	CurrentPlan    *data.BillingPlan `json:"currentPlan"`
	Seats          int               `json:"seats"`
	Logins         []model.User      `json:"logins"`
}

// BillingCardData represents a Stripe credit card
type BillingCardData struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	Number     string `json:"number"`
	Month      string `json:"month"`
	Year       string `json:"year"`
	CVC        string `json:"cvc"`
	Brand      string `json:"brand"`
	Expiration string `json:"expiration"`
}

func newBilling() *engine.Route {
	var b interface{} = User{}
	return &engine.Route{
		Logger:      true,
		MinimumRole: model.RoleAdmin,
		Handler:     b.(http.Handler),
	}
}

func (b Billing) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var head string
	head, r.URL.Path = engine.ShiftPath(r.URL.Path)
	if r.Method == http.MethodGet {
		if head == "overview" {
			b.overview(w, r)
		} else if head == "invoices" {
			head, r.URL.Path = engine.ShiftPath(r.URL.Path)
			if head == "" {
				b.invoices(w, r)
			} else if head == "next" {
				b.getNextInvoice(w, r)
				return
			}
		}
	} else if r.Method == http.MethodPost {
		if head == "start" {
			b.start(w, r)
		} else if head == "changeplan" {
			b.changePlan(w, r)
		} else if head == "webhooks" {
			b.stripe(w, r)
		}
	}
}

func (b Billing) overview(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	keys := ctx.Value(engine.ContextAuth).(engine.Auth)
	db := ctx.Value(engine.ContextDatabase).(*data.DB)

	// this struct will be returned should we be a paid customer or not
	ov := BillingOverview{}

	// Get the current account
	account, err := db.Users.GetDetail(keys.AccountID)
	if err != nil {
		engine.Respond(w, r, http.StatusNotFound, err)
		return
	}

	ov.Logins = account.Users

	// get all logins for user roles
	for _, l := range account.Users {
		if l.Role < model.RoleFree {
			ov.Seats++
		}
	}

	if len(account.StripeID) == 0 {
		ov.IsNew = true

		// if they are on trial, we set the current plan to
		// that so the UI can based permissions on that plan.
		if account.TrialInfo.IsTrial {
			if p, ok := data.GetPlan(account.TrialInfo.Plan); ok {
				ov.CurrentPlan = &p
			}
		}

		engine.Respond(w, r, http.StatusOK, ov)
		return
	}

	// getting stripe customer
	cus, err := customer.Get(account.StripeID, nil)
	if err != nil {
		engine.Respond(w, r, http.StatusBadRequest, err)
		return
	}

	ov.StripeID = cus.ID
	ov.Plan = account.Plan
	ov.IsYearly = account.IsYearly

	if p, ok := data.GetPlan(account.Plan); ok {
		ov.CurrentPlan = &p
	}

	cards := card.List(&stripe.CardListParams{Customer: stripe.String(account.StripeID)})
	for cards.Next() {
		c := cards.Card()
		if !c.Deleted {
			ov.Cards = append(ov.Cards, BillingCardData{
				ID:         c.ID,
				Name:       c.Name,
				Number:     c.Last4,
				Month:      fmt.Sprintf("%d", c.ExpMonth),
				Year:       fmt.Sprintf("%d", c.ExpYear),
				Expiration: fmt.Sprintf("%d / %d", c.ExpMonth, c.ExpYear),
				Brand:      string(c.Brand),
			})
		}
	}

	engine.Respond(w, r, http.StatusOK, ov)
}

func (b Billing) changeQuantity(stripeID, subID string, qty int) error {
	p := &stripe.SubscriptionParams{Customer: stripe.String(stripeID), Quantity: stripe.Int64(int64(qty))}
	_, err := sub.Update(subID, p)
	return err
}

func (b Billing) userRoleChanged(db data.DB, accountID model.Key, oldRole, newRole model.Roles) (paid bool, err error) {
	acct, err := db.Users.GetDetail(accountID)
	if err != nil {
		return false, err
	}

	// if this is a paid account
	if acct.IsPaid() {
		// if they were a free user
		if oldRole == model.RoleFree {
			// and are now a paid user, we need to +1 qty and prepare the invoice
			if newRole == model.RoleAdmin || newRole == model.RoleUser {
				paid = true

				// we increase the seats number for this account
				acct.Seats++

				// try to change their subscription (+1 qty)
				if err = b.changeQuantity(acct.StripeID, acct.SubscriptionID, acct.Seats); err != nil {
					return
				}

				// ensure that the charges will be immediate and not on next billing date
				if err := queue.Enqueue(queue.TaskCreateInvoice, acct.StripeID); err != nil {
					return paid, err
				}

				if err = db.Users.SetSeats(acct.ID, acct.Seats); err != nil {
					return false, err
				}
			}
		} else {
			// they were a paid user, now they are set as free
			if newRole == model.RoleFree {
				acct.Seats--

				if err = b.changeQuantity(acct.StripeID, acct.SubscriptionID, acct.Seats); err != nil {
					return
				}

				if err = db.Users.SetSeats(acct.ID, acct.Seats); err != nil {
					return false, err
				}
			}
		}
	}
	return false, nil
}

// BillingNewCustomer represents data sent to api for creating a new customer
type BillingNewCustomer struct {
	Plan     string          `json:"plan"`
	Card     BillingCardData `json:"card"`
	Coupon   string          `json:"coupon"`
	Zip      string          `json:"zip"`
	IsYearly bool            `json:"yearly"`
}

func (b Billing) start(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	keys := ctx.Value(engine.ContextAuth).(engine.Auth)
	db := ctx.Value(engine.ContextDatabase).(*data.DB)

	var data BillingNewCustomer
	if err := engine.ParseBody(r.Body, &data); err != nil {
		engine.Respond(w, r, http.StatusBadRequest, err)
		return
	}

	p := &stripe.CustomerParams{Email: stripe.String(keys.Email)}
	p.SetSource(&stripe.CardParams{
		Name:       stripe.String(data.Card.Name),
		Number:     stripe.String(data.Card.Number),
		ExpMonth:   stripe.String(data.Card.Month),
		ExpYear:    stripe.String(data.Card.Year),
		CVC:        stripe.String(data.Card.CVC),
		AddressZip: stripe.String(data.Zip),
	})

	c, err := customer.New(p)
	if err != nil {
		engine.Respond(w, r, http.StatusInternalServerError, err)
		return
	}

	acct, err := db.Users.GetDetail(keys.AccountID)
	if err != nil {
		engine.Respond(w, r, http.StatusInternalServerError, err)
		return
	}

	seats := 0
	for _, u := range acct.Users {
		if u.Role < model.RoleFree {
			seats++
		}
	}

	plan := data.Plan
	if data.IsYearly {
		plan += "_yearly"
	}

	// Coupon:   "PRELAUNCH11",
	subp := &stripe.SubscriptionParams{
		Customer: stripe.String(c.ID),
		Plan:     stripe.String(plan),
		Quantity: stripe.Int64(int64(seats)),
	}

	if len(data.Coupon) > 0 {
		subp.Coupon = stripe.String(data.Coupon)
	}

	s, err := sub.New(subp)
	if err != nil {
		engine.Respond(w, r, http.StatusInternalServerError, err)
		return
	}

	acct.TrialInfo.IsTrial = false
	if err := db.Users.ConvertToPaid(acct.ID, c.ID, s.ID, data.Plan, data.IsYearly, seats); err != nil {
		engine.Respond(w, r, http.StatusInternalServerError, err)
		return
	}

	ov := BillingOverview{}
	ov.StripeID = c.ID
	ov.Plan = data.Plan
	ov.IsYearly = data.IsYearly
	ov.Seats = seats

	acct.StripeID = c.ID
	acct.SubscribedOn = time.Now()
	acct.SubscriptionID = s.ID
	acct.Plan = data.Plan
	acct.IsYearly = data.IsYearly
	acct.Seats = seats

	//TODO: Trigger a new customer event

	engine.Respond(w, r, http.StatusOK, ov)
}

func (b Billing) changePlan(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	keys := ctx.Value(engine.ContextAuth).(engine.Auth)
	db := ctx.Value(engine.ContextDatabase).(*data.DB)

	var data = new(struct {
		Plan     string `json:"plan"`
		IsYearly bool   `json:"isYearly"`
	})
	if err := engine.ParseBody(r.Body, &data); err != nil {
		engine.Respond(w, r, http.StatusBadRequest, err)
		return
	}

	account, err := db.Users.GetDetail(keys.AccountID)
	if err != nil {
		engine.Respond(w, r, http.StatusInternalServerError, err)
		return
	}

	plan := data.Plan

	newLevel, currentLevel := 0, 0
	if len(plan) == 0 || plan == "free" {
		newLevel = 0
	} else if strings.HasPrefix(plan, "starter") {
		newLevel = 1
	} else if strings.HasPrefix(plan, "pro") {
		newLevel = 2
	} else {
		newLevel = 3
	}

	if strings.HasPrefix(account.Plan, "starter") {
		currentLevel = 1
	} else if strings.HasPrefix(account.Plan, "pro") {
		currentLevel = 2
	} else {
		currentLevel = 3
	}

	// did they cancelled
	if newLevel == 0 {
		// we need to cancel their subscriptions
		if _, err := sub.Cancel(account.SubscriptionID, nil); err != nil {
			engine.Respond(w, r, http.StatusInternalServerError, err)
			return
		}

		if err := db.Users.Cancel(account.ID); err != nil {
			engine.Respond(w, r, http.StatusInternalServerError, err)
			return
		}
	} else {
		if data.IsYearly {
			plan += "_yearly"
		}

		seats := 0
		for _, u := range account.Users {
			if u.Role < model.RoleFree {
				seats++
			}
		}

		subParams := &stripe.SubscriptionParams{Customer: stripe.String(account.StripeID),
			Plan:     stripe.String(plan),
			Quantity: stripe.Int64(int64(seats)),
		}
		// if we upgrade we need to change billing cycle date
		upgraded := false
		if newLevel > currentLevel {
			upgraded = true
		} else if account.IsYearly == false && data.IsYearly {
			upgraded = true
		}

		if upgraded {
			// queue an invoice create for this upgrade
			queue.Enqueue(queue.TaskCreateInvoice, account.StripeID)
		}

		if _, err := sub.Update(account.SubscriptionID, subParams); err != nil {
			engine.Respond(w, r, http.StatusInternalServerError, err)
			return
		}

		if err := db.Users.ChangePlan(account.ID, plan, data.IsYearly); err != nil {
			engine.Respond(w, r, http.StatusInternalServerError, err)
			return
		}
		engine.Respond(w, r, http.StatusOK, true)
	}
}

func (b Billing) updateCard(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	keys := ctx.Value(engine.ContextAuth).(engine.Auth)
	db := ctx.Value(engine.ContextDatabase).(*data.DB)

	account, err := db.Users.GetDetail(keys.AccountID)
	if err != nil {
		engine.Respond(w, r, http.StatusBadRequest, err)
		return
	}

	var data BillingCardData
	if err := engine.ParseBody(r.Body, &data); err != nil {
		engine.Respond(w, r, http.StatusInternalServerError, err)
		return
	}

	p := &stripe.CardParams{
		Customer: stripe.String(account.StripeID),
		ExpMonth: stripe.String(data.Month),
		ExpYear:  stripe.String(data.Month),
		CVC:      stripe.String(data.CVC),
	}
	if c, err := card.Update(data.ID, p); err != nil {
		engine.Respond(w, r, http.StatusInternalServerError, err)
	} else {
		card := BillingCardData{
			ID:         c.ID,
			Name:       c.Name,
			Number:     c.Last4,
			Month:      fmt.Sprintf("%d", c.ExpMonth),
			Year:       fmt.Sprintf("%d", c.ExpYear),
			Expiration: fmt.Sprintf("%d / %d", c.ExpMonth, c.ExpYear),
			Brand:      string(c.Brand),
		}
		engine.Respond(w, r, http.StatusOK, card)
	}
}

func (b Billing) addCard(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	keys := ctx.Value(engine.ContextAuth).(engine.Auth)
	db := ctx.Value(engine.ContextDatabase).(*data.DB)

	account, err := db.Users.GetDetail(keys.AccountID)
	if err != nil {
		engine.Respond(w, r, http.StatusBadRequest, err)
		return
	}

	var data BillingCardData
	if err := engine.ParseBody(r.Body, &data); err != nil {
		engine.Respond(w, r, http.StatusInternalServerError, err)
		return
	}

	newCard := &stripe.CardParams{
		Customer: stripe.String(account.StripeID),
		Name:     stripe.String(data.Name),
		Number:   stripe.String(data.Number),
		ExpMonth: stripe.String(data.Month),
		ExpYear:  stripe.String(data.Year),
		CVC:      stripe.String(data.CVC),
	}
	if c, err := card.New(newCard); err != nil {
		engine.Respond(w, r, http.StatusInternalServerError, err)
	} else {
		card := BillingCardData{
			ID:         c.ID,
			Number:     c.Last4,
			Expiration: fmt.Sprintf("%d / %d", c.ExpMonth, c.ExpYear),
			Brand:      string(c.Brand),
		}
		engine.Respond(w, r, http.StatusOK, card)
	}
}

func (b Billing) deleteCard(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	keys := ctx.Value(engine.ContextAuth).(engine.Auth)
	db := ctx.Value(engine.ContextDatabase).(*data.DB)

	account, err := db.Users.GetDetail(keys.AccountID)
	if err != nil {
		engine.Respond(w, r, http.StatusBadRequest, err)
		return
	}

	cardID := mux.Vars(r)["id"]

	if _, err := card.Del(cardID, &stripe.CardParams{Customer: stripe.String(account.StripeID)}); err != nil {
		engine.Respond(w, r, http.StatusInternalServerError, err)
	} else {
		engine.Respond(w, r, http.StatusOK, true)
	}
}

func (b Billing) invoices(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	keys := ctx.Value(engine.ContextAuth).(engine.Auth)
	db := ctx.Value(engine.ContextDatabase).(*data.DB)

	account, err := db.Users.GetDetail(keys.AccountID)
	if err != nil {
		engine.Respond(w, r, http.StatusInternalServerError, err)
		return
	}

	var invoices []*stripe.Invoice

	iter := invoice.List(&stripe.InvoiceListParams{Customer: stripe.String(account.StripeID)})
	for iter.Next() {
		invoices = append(invoices, iter.Invoice())
	}

	engine.Respond(w, r, http.StatusOK, invoices)
}

func (b Billing) getNextInvoice(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	keys := ctx.Value(engine.ContextAuth).(engine.Auth)
	db := ctx.Value(engine.ContextDatabase).(*data.DB)

	account, err := db.Users.GetDetail(keys.AccountID)
	if err != nil {
		engine.Respond(w, r, http.StatusInternalServerError, err)
		return
	}

	i, err := invoice.GetNext(&stripe.InvoiceParams{Customer: stripe.String(account.StripeID)})
	if err != nil {
		engine.Respond(w, r, http.StatusInternalServerError, err)
		return
	}
	engine.Respond(w, r, http.StatusOK, i)
}

// StripeWebhook is used to grab data sent by Stripe for a webhook
type StripeWebhook struct {
	Event stripe.Event `json:"event"`
}

// WebhookData used when stripe webhook is call
type WebhookData struct {
	ID   string            `json:"id"`
	Type string            `json:"type"`
	Data WebhookDataObject `json:"data"`
}

// WebhookDataObject is the container for the object received
type WebhookDataObject struct {
	Object WebhookDataObjectData `json:"object"`
}

// WebhookDataObjectData is the object being sent by stripe
type WebhookDataObjectData struct {
	ID           string `json:"id"`
	Customer     string `json:"customer"`
	Subscription string `json:"subscription"`
	Closed       bool   `json:"closed"`
}

func (b Billing) stripe(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	db := ctx.Value(engine.ContextDatabase).(*data.DB)

	// no matter what happen, Stripe wants us to send a 200
	defer w.Write([]byte("ok"))

	var data WebhookData
	if err := engine.ParseBody(r.Body, &data); err != nil {
		log.Println(err)
		return
	}

	if data.Type == "customer.subscription.deleted" {
		subID := data.Data.Object.ID
		if len(subID) == 0 {
			log.Println(fmt.Errorf("no subscription found to this customer.subscription.deleted %s", data.ID))
			return
		}

		stripeID := data.Data.Object.Customer
		if len(stripeID) == 0 {
			log.Println(fmt.Errorf("no customer found to this invoice.payment_succeeded %s", data.ID))
			return
		}

		// check if it's a failed payment_succeeded
		account, err := db.Users.GetByStripe(stripeID)
		if err != nil {
			log.Println(fmt.Errorf("no customer matches stripe id: %s", stripeID))
			return
		}

		if len(account.SubscriptionID) > 0 {
			//TODO: Send emails

			if err := db.Users.Cancel(account.ID); err != nil {
				log.Println(fmt.Errorf("unable to cancel this account: %v", account.ID))
				return
			}
		}
	}
}

func (b Billing) cancel(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	keys := ctx.Value(engine.ContextAuth).(engine.Auth)
	db := ctx.Value(engine.ContextDatabase).(*data.DB)

	var data = new(struct {
		Reason string `json:"reason"`
	})
	if err := engine.ParseBody(r.Body, &data); err != nil {
		engine.Respond(w, r, http.StatusBadRequest, err)
		return
	}

	// SendMail would be call here passing the reason

	account, err := db.Users.GetDetail(keys.AccountID)
	if err != nil {
		engine.Respond(w, r, http.StatusInternalServerError, err)
		return
	}

	if _, err := sub.Cancel(account.SubscriptionID, nil); err != nil {
		engine.Respond(w, r, http.StatusInternalServerError, err)
		return
	}

	if err := db.Users.Cancel(account.ID); err != nil {
		engine.Respond(w, r, http.StatusInternalServerError, err)
		return
	}

	engine.Respond(w, r, http.StatusOK, true)
}
