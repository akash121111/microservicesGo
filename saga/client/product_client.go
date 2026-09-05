package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"saga/model"
)

type ProductClient struct {
	client  *http.Client
	baseUrl string
}

func NewProductClient(client *http.Client, baseUrl string) *ProductClient {
	return &ProductClient{
		client:  client,
		baseUrl: baseUrl,
	}
}

func (c *ProductClient) ReleaseStock(
	ctx context.Context,
	items []model.OrderEventItem,
) error {

	url := fmt.Sprintf(
		"%s/product/stock/release",
		c.baseUrl,
	)
	log.Println(url)

	body, err := json.Marshal(map[string][]model.OrderEventItem{
		"items": items,
	})
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPatch,
		url,
		bytes.NewBuffer(body),
	)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf(
			"product service returned status %d",
			resp.StatusCode,
		)
	}

	return nil
}
