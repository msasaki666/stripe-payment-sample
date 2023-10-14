package main

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/stripe/stripe-go/v75/webhook"
)

func handleWebhook(w http.ResponseWriter, req *http.Request) {
	const MaxBodyBytes = int64(65536)
	req.Body = http.MaxBytesReader(w, req.Body, MaxBodyBytes)
	payload, err := io.ReadAll(req.Body)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading request body: %v\n", err)
		w.WriteHeader(http.StatusServiceUnavailable)
		return
	}

	// TODO: 変更する
	// This is your Stripe CLI webhook secret for testing your endpoint locally.
	endpointSecret := "whsec_cee980c3fb2443c92af0cec2996e629b3c34d4890cd6920ca0e962f54356d639"
	// Pass the request body and Stripe-Signature header to ConstructEvent, along
	// with the webhook signing key.
	event, err := webhook.ConstructEvent(payload, req.Header.Get("Stripe-Signature"),
		endpointSecret)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error verifying webhook signature: %v\n", err)
		w.WriteHeader(http.StatusBadRequest) // Return a 400 error on a bad signature
		return
	}
	// TODO: イベントの種類によって処理を変える
	// checkoutで購入完了した時のイベント
	// stripe.EventTypeCheckoutSessionCompleted
	// https://stripe.com/docs/payments/checkout/how-checkout-works#complete
	// checkoutでセッションの期限が切れた時のイベント。在庫確保する系なら、セッション作成後に在庫を確保し、これを機に在庫を戻す
	// stripe.EventTypeCheckoutSessionExpired
	switch event.Type {
	// case stripe.EventTypePaymentIntentSucceeded:
	// 	var paymentIntent stripe.PaymentIntent
	// 	err := json.Unmarshal(event.Data.Raw, &paymentIntent)
	// 	if err != nil {
	// 		fmt.Fprintf(os.Stderr, "Error parsing webhook JSON: %v\n", err)
	// 		w.WriteHeader(http.StatusBadRequest)
	// 		return
	// 	}
	// 	// Then define and call a func to handle the successful payment intent.
	// 	// handlePaymentIntentSucceeded(paymentIntent)
	// case stripe.EventTypePaymentMethodAttached:
	// 	var paymentMethod stripe.PaymentMethod
	// 	err := json.Unmarshal(event.Data.Raw, &paymentMethod)
	// 	if err != nil {
	// 		fmt.Fprintf(os.Stderr, "Error parsing webhook JSON: %v\n", err)
	// 		w.WriteHeader(http.StatusBadRequest)
	// 		return
	// 	}
	// 	// Then define and call a func to handle the successful attachment of a PaymentMethod.
	// 	// handlePaymentMethodAttached(paymentMethod)
	// // ... handle other event types
	default:
		fmt.Fprintf(os.Stderr, "Unhandled event type: %s\n", event.Type)
	}

	w.WriteHeader(http.StatusOK)
}
