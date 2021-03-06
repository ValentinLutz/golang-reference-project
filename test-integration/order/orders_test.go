package order_test

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"os"
	"test-integration/order"
	"testing"
)

func initClient() order.Client {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	return order.Client{
		Server: "http://localhost:8080",
		Client: client,
	}
}

func TestGetOrder(t *testing.T) {
	// GIVEN
	client := initClient()

	// WHEN
	apiOrder, err := client.GetApiOrdersOrderId(context.Background(), "IsQah2TkaqS-NONE-DEV-JewgL0Ye73g")
	if err != nil {
		t.Fatal(err)
	}
	defer apiOrder.Body.Close()

	// THEN
	var actualResponse order.OrderResponse
	readToObject(t, apiOrder.Body, &actualResponse)
	var expectedResponse order.OrderResponse
	readToObject(t, readFile(t, "orderResponse.json"), &expectedResponse)

	assert.Equal(t, 200, apiOrder.StatusCode)
	assert.Equal(t, expectedResponse, actualResponse)
}

func TestGetOrderNotFound(t *testing.T) {
	// GIVEN
	client := initClient()

	// WHEN
	apiOrder, err := client.GetApiOrdersOrderId(context.Background(), "NOPE")
	if err != nil {
		t.Fatal(err)
	}
	defer apiOrder.Body.Close()

	// THEN
	var actualResponse order.ErrorResponse
	readToObject(t, apiOrder.Body, &actualResponse)
	var expectedResponse order.ErrorResponse
	readToObject(t, readFile(t, "orderNotFoundResponse.json"), &expectedResponse)
	expectedResponse.Timestamp = actualResponse.Timestamp

	assert.Equal(t, 404, apiOrder.StatusCode)
	assert.Equal(t, expectedResponse, actualResponse)
}

func TestGetOrders(t *testing.T) {
	// GIVEN
	client := initClient()

	// WHEN
	apiOrder, err := client.GetApiOrders(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	defer apiOrder.Body.Close()

	// THEN
	var actualResponse order.OrdersResponse
	readToObject(t, apiOrder.Body, &actualResponse)
	var expectedResponse order.OrdersResponse
	readToObject(t, readFile(t, "ordersResponse.json"), &expectedResponse)

	assert.Equal(t, 200, apiOrder.StatusCode)
	assert.Equal(t, expectedResponse, actualResponse)
}

func readToObject(t *testing.T, reader io.Reader, object interface{}) {
	decoder := json.NewDecoder(reader)
	err := decoder.Decode(object)
	if err != nil {
		t.Fatalf("Failed to decode input, %v", err)
	}
}

func readFile(t *testing.T, path string) *os.File {
	file, err := os.Open(path)
	if err != nil {
		t.Fatalf("Failed to read file from path %v, %v", path, err)
	}
	return file
}
