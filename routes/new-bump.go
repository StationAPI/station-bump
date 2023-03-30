package routes

import (
	"net/http"
	"os"

	neon "github.com/stationapi/station-bump/db"
	"github.com/stationapi/station-bump/session"
	"github.com/stripe/stripe-go/v74"
	checkout "github.com/stripe/stripe-go/v74/checkout/session"
	"gorm.io/gorm"
)

type Request struct {
	Id string `json:"site_id"`
}

func NewBump(w http.ResponseWriter, r *http.Request, db gorm.DB) error {
	cookie, cookieErr := r.Cookie("station")		

	if cookieErr != nil {
		http.Error(w, "you are not authentictaed", http.StatusForbidden)

		return cookieErr
	}

	githubId, authErr := session.GetSession(cookie.Value) 

	if authErr != nil {
		http.Error(w, "you are not authenticated", http.StatusForbidden)

		return authErr
	}

	website := Request{}

	err := ProcessBody(r.Body, &website)

	if err != nil {
		http.Error(w, "there was an error processing the request body", http.StatusBadRequest)

		return err
	}

	_, user := neon.GetUser(githubId, db)

	if user.Bumps > 0 {
		neon.Bump(website.Id, db)		

		return nil
	}		

	stripe.Key = os.Getenv("STRIPE_KEY")

	stripeSession, err := checkout.New(
		&stripe.CheckoutSessionParams{
			LineItems: []*stripe.CheckoutSessionLineItemParams{
				{
					Price: stripe.String(os.Getenv("STRIPE_PRICE_ID")),
					Quantity: stripe.Int64(1),
				},
			},
			Mode: stripe.String("payment"),
			SuccessURL: stripe.String(os.Getenv("STRIPE_SUCCESS_URL")),
		},
	)

	if err != nil {
		http.Error(w, "there was an error generating the checkout link", http.StatusInternalServerError)

		return nil
	}

	return nil
}
