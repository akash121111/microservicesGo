package client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"saga/model"

	"github.com/google/uuid"
)

type OrderClient struct {
	client  *http.Client
	baseUrl string
}

func NewOrderClient(client *http.Client, baseUrl string) *OrderClient {
	return &OrderClient{
		client:  client,
		baseUrl: baseUrl,
	}
}

func (c *OrderClient) CancelOrder(
	ctx context.Context,
	orderID uuid.UUID,
) error {

	url := fmt.Sprintf(
		"%s/order/cancel/%s",
		c.baseUrl,
		orderID,
	)

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPatch,
		url,
		nil,
	)

	if err != nil {
		return err
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		return fmt.Errorf(
			"cancel order failed: status=%d",
			resp.StatusCode,
		)
	}

	return nil
}

func (c *OrderClient) UpdateOrderStatus(
	ctx context.Context,
	orderID uuid.UUID,
	status string,
) error {

	url := fmt.Sprintf(
		"%s/order/%s/status/%s",
		c.baseUrl,
		orderID,
		status,
	)

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPatch,
		url,
		nil,
	)

	if err != nil {
		return err
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		return fmt.Errorf(
			"cancel order failed: status=%d",
			resp.StatusCode,
		)
	}

	return nil
}

func (c *OrderClient) GetOrderById(ctx context.Context, orderId uuid.UUID) (model.OrderDataWithItems, error) {
	url := fmt.Sprintf("%s/order/%s", c.baseUrl, orderId)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)

	if err != nil {
		return model.OrderDataWithItems{}, err
	}
	resp, err := c.client.Do(req)
	if err != nil {
		return model.OrderDataWithItems{}, err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		return model.OrderDataWithItems{}, fmt.Errorf(
			"get order failed: status=%d",
			resp.StatusCode,
		)
	}
	var responce model.APIResponse[model.OrderDataWithItems]
	if err := json.NewDecoder(resp.Body).Decode(&responce); err != nil {
		return model.OrderDataWithItems{}, fmt.Errorf("decode order response: %w", err)
	}
	return responce.Data, nil

}
