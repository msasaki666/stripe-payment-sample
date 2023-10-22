package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/stripe/stripe-go/v76"
	"github.com/stripe/stripe-go/v76/webhook"
)

func handleWebhook(w http.ResponseWriter, req *http.Request) {
	const MaxBodyBytes = int64(65536)
	req.Body = http.MaxBytesReader(w, req.Body, MaxBodyBytes)

	body, err := io.ReadAll(req.Body)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading request body: %v\n", err)
		w.WriteHeader(http.StatusServiceUnavailable)
		return
	}

	// TODO: 変更する
	// This is your Stripe CLI webhook secret for testing your endpoint locally.
	endpointSecret := "sk_test_51NvebBAj9ehS6HaZoiSt7gY1MXsCOhKAPeM8LGbzL5BhL0PHPgE7XXaVob0ClnzMP12MNZ5e62RxdLr9bDMhP0Gk00s6HZZUAn"
	// Pass the request body and Stripe-Signature header to ConstructEvent, along
	// with the webhook signing key.
	event, err := webhook.ConstructEvent(
		body,
		req.Header.Get("Stripe-Signature"),
		endpointSecret,
	)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error verifying webhook signature: %v\n", err)
		w.WriteHeader(http.StatusBadRequest) // Return a 400 error on a bad signature
		return
	}

	// サブスクリプションでモニタリングが必要な最小限のイベントタイプ:
	// イベント名	説明
	// checkout.session.completed	Checkout で顧客が「支払う」または「登録」ボタンをクリックすると送信され、新しい購入が通知されます。
	// invoice.paid	請求期間ごとに、支払いが成功すると送信されます。
	// invoice.payment_failed	請求期間ごとに、顧客の支払い方法に問題がある場合に送信されます。
	// TODO: イベントの種類によって処理を変える
	// checkoutで購入完了した時のイベント
	// stripe.EventTypeCheckoutSessionCompleted
	// https://stripe.com/docs/payments/checkout/how-checkout-works#complete
	// checkoutでセッションの期限が切れた時のイベント。在庫確保する系なら、セッション作成後に在庫を確保し、これを機に在庫を戻す
	// stripe.EventTypeCheckoutSessionExpired
	switch event.Type {
	case stripe.EventTypeCheckoutSessionCompleted:
		var session stripe.CheckoutSession
		err := json.Unmarshal(event.Data.Raw, &session)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error parsing webhook JSON: %v\n", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		fmt.Println("OK")
		// params := &stripe.CheckoutSessionParams{}
		// params.AddExpand("line_items")
		// // Retrieve the session. If you require line items in the response, you may include them by expanding line_items.
		// session
		// sessionWithLineItems, _ := session.Get(session.ID, params)
		// lineItems := sessionWithLineItems.LineItems
		// // Fulfill the purchase...
		// FulfillOrder(lineItems)

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
