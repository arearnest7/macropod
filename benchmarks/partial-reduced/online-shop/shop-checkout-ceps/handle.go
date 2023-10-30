package function

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"net"
	"os"
	"time"

	"cloud.google.com/go/profiler"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/GoogleCloudPlatform/microservices-demo/src/checkoutservice/genproto"
	money "github.com/GoogleCloudPlatform/microservices-demo/src/checkoutservice/money"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"

	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

func PlaceOrder(var req) {
	log.Infof("[PlaceOrder] user_id=%q user_currency=%q", req.UserId, req.UserCurrency)

	orderID, err := uuid.NewUUID()
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to generate order uuid")
	}

	prep, err := prepareOrderItemsAndShippingQuoteFromCart(req.UserId, req.UserCurrency, req.Address)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	total := pb.Money{CurrencyCode: req.UserCurrency,
		Units: 0,
		Nanos: 0}
	total = money.Must(money.Sum(total, *prep.shippingCostLocalized))
	for _, it := range prep.orderItems {
		multPrice := money.MultiplySlow(*it.Cost, uint32(it.GetItem().GetQuantity()))
		total = money.Must(money.Sum(total, multPrice))
	}

	txID, err := chargeCard(&total, req.CreditCard)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to charge card: %+v", err)
	}
	log.Infof("payment went through (transaction_id: %s)", txID)

	shippingTrackingID, err := shipOrder(req.Address, prep.cartItems)
	if err != nil {
		return nil, status.Errorf(codes.Unavailable, "shipping error: %+v", err)
	}

	_ = emptyUserCart(req.UserId)

	orderResult := &pb.OrderResult{
		OrderId:            orderID.String(),
		ShippingTrackingId: shippingTrackingID,
		ShippingCost:       prep.shippingCostLocalized,
		ShippingAddress:    req.Address,
		Items:              prep.orderItems,
	}

	if err := sendOrderConfirmation(req.Email, orderResult); err != nil {
		log.Warnf("failed to send order confirmation to %q: %+v", req.Email, err)
	} else {
		log.Infof("order confirmation email sent to %q", req.Email)
	}
	resp := &pb.PlaceOrderResponse{Order: orderResult}
	return resp, nil
}

func prepareOrderItemsAndShippingQuoteFromCart(var userID, var userCurrency, var address) {
	var out orderPrep
	cartItems, err := getUserCart(userID)
	if err != nil {
		return out, fmt.Errorf("cart failure: %+v", err)
	}
	orderItems, err := prepOrderItems(cartItems, userCurrency)
	if err != nil {
		return out, fmt.Errorf("failed to prepare order: %+v", err)
	}
	shippingUSD, err := quoteShipping(address, cartItems)
	if err != nil {
		return out, fmt.Errorf("shipping quote failure: %+v", err)
	}
	shippingPrice, err := convertCurrency(shippingUSD, userCurrency)
	if err != nil {
		return out, fmt.Errorf("failed to convert shipping cost to currency: %+v", err)
	}

	out.shippingCostLocalized = shippingPrice
	out.cartItems = cartItems
	out.orderItems = orderItems
	return out, nil
}

func quoteShipping(var address, var items) {
	resp := os.system("go run shipping " + json.dumps({requestType: "quote", Address: address, Items: items}))
	return resp.GetCostUsd(), nil
}

func getUserCart(var userID) {
        contents, err := ioutil.ReadFile("/etc/secret-volume/shop-frontend-arpc")
        if err != nil {
                log.Fatal(err)
        }
        shop-currency := string(contents)
        requestURL := shop-currency + ":80"
        req_url, err := http.NewRequest(http.MethodPost, requestURL, {request: "cart", requestType: "get", UserID: userID})
        resp, err := http.DefaultClient.Do(req_url)
        if err != nil {
                return nil, fmt.Errorf("failed to get user cart during checkout: %+v", err)
        }
        return resp.GetItems(), nil
}

func emptyUserCart(var userID) {
        contents, err := ioutil.ReadFile("/etc/secret-volume/shop-frontend-arpc")
        if err != nil {
                log.Fatal(err)
        }
        shop-cart := string(contents)
        requestURL := shop-cart + ":80"
        req_url, err := http.NewRequest(http.MethodPost, requestURL, {request: "cart", requestType: "empty", UserId: userID})
        resp, err := http.DefaultClient.Do(req_url)
        if err != nil {
                return fmt.Errorf("failed to empty user cart during checkout: %+v", err)
        }
        return nil
}

func prepOrderItems(var items) {
	out := make([]*pb.OrderItem, len(items))
        contents, err := ioutil.ReadFile("/etc/secret-volume/shop-frontend-arpc")
        if err != nil {
                log.Fatal(err)
        }
        shop-productcatalog := string(contents)
        requestURL := shop-productcatalog + ":80"

	for i, item := range items {
        	req_url, err := http.NewRequest(http.MethodPost, requestURL, {request: "productcatalog", &pb.GetProductRequest{I>
                product, err := http.DefaultClient.Do(req_url)
                if err != nil {
                        return nil, fmt.Errorf("failed to get product #%q", item.GetProductId())
                }
		price, err := convertCurrency(product.GetPriceUsd(), userCurrency)
		if err != nil {
			return nil, fmt.Errorf("failed to convert price of %q to %s", item.GetProductId(), userCurrency)
		}
		out[i] = &pb.OrderItem{
			Item: item,
			Cost: price}
	}
	return out, nil
}

func convertCurrency(var from, var toCurrency) {
	result, err := os.system("npm run currency.js " + json.dumps({From: from, ToCode: toCurrency}))
	return result, err
}

func chargeCard(var amount, var paymentInfo) {
	paymentResp, err := os.system("npm run payment.js " + json.dumps({Amount: amount, CreditCard: paymentInfo}))
	return paymentResp.GetTransactionId(), nil
}

func sendOrderConfirmation(var email, var order) {
	resp, err := os.system("go run shipping.go " + json.dumps({Email: email, Order: order}))
	return err
}

func shipOrder(var address, var order) {
	resp, err := os.system("go run shipping.go " + json.dumps({requestType: "order", Address: address, Items: items}))
	return resp.GetTrackingId(), nil
}

// Handle an HTTP Request.
func Handle(ctx context.Context, res http.ResponseWriter, req *http.Request) {
	if req.json.request == "checkout" {
        	fmt.Fprintf(res, PlaceOrder(req.json)) // echo to caller
	}
	else if req.json.request == "currency" {
		fmt.Fprintf(res, os.system("npm run currency.js " + req.json)) // echo to caller
	}
	else if req.json.request == "shipping" {
		fmt.Fprintf(res, os.system("go run shipping.go " + req.json)) // echo to caller
	}
	fmt.Fprintf(res, "") // echo to caller
}
