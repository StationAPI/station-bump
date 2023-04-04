package routes

import (
	"errors"
	"net/http"

	neon "github.com/stationapi/station-bump/db"
	sess "github.com/stationapi/station-bump/session"
	"github.com/stripe/stripe-go/v74"
	"github.com/stripe/stripe-go/v74/checkout/session"
	"gorm.io/gorm"
)

func CheckoutSuccess(w http.ResponseWriter, r *http.Request, db gorm.DB) error {
	sessionId := r.URL.Query().Get("session_id")

	session, err := session.Get(sessionId, &stripe.CheckoutSessionParams{})

	if err != nil {
		http.Error(w, "there was an error fetching the session", http.StatusInternalServerError)

		return err
	}

	if session.PaymentIntent.Status != stripe.PaymentIntentStatusSucceeded {
		http.Error(w, "the checkout session was invalid", http.StatusForbidden)

		return errors.New("unsuccessful checkout session")
	}

	cookie, err := r.Cookie("station")

	if err != nil {
		http.Error(w, "you are not authenticated", http.StatusForbidden)

		return err
	}

	githubId, err := sess.GetSession(cookie.Value)

	if err != nil {
		http.Error(w, "you are not authenticated", http.StatusForbidden)

		return err
	}

	_, user := neon.GetUser(githubId, db)

	user.Bumps += 1

	db.Save(user)

	w.WriteHeader(200)
	w.Write([]byte("you have a new bump available"))

	return nil
}
