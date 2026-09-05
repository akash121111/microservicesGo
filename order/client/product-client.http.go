package client

// type ProductClient struct {
// 	client  *http.Client
// 	baseUrl string
// }

// func NewProductClient(client *http.Client, baseUrl string) *ProductClient {
// 	return &ProductClient{
// 		client:  client,
// 		baseUrl: baseUrl,
// 	}
// }

// func (c *ProductClient) GetProduct(
// 	ctx context.Context,
// 	id string,
// ) (*model.ProductModel, error) {

// 	url := c.baseUrl + "/product/" + id

// 	maxRetries := 3

// 	for attempt := 0; attempt <= maxRetries; attempt++ {

// 		req, err := http.NewRequestWithContext(
// 			ctx,
// 			http.MethodGet,
// 			url,
// 			nil,
// 		)

// 		if err != nil {
// 			return nil, err
// 		}

// 		resp, err := c.client.Do(req)

// 		// Network error
// 		if err != nil {

// 			if !isRetryableError(err) || attempt == maxRetries {
// 				return nil, err
// 			}

// 			if err := waitForRetry(ctx, attempt); err != nil {
// 				return nil, err
// 			}

// 			continue
// 		}

// 		// Success
// 		if resp.StatusCode == http.StatusOK {

// 			defer resp.Body.Close()

// 			var response model.APIResponse[model.ProductModel]

// 			if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
// 				return nil, err
// 			}

// 			return &response.Data, nil
// 		}

// 		// Retry 5xx
// 		if resp.StatusCode >= 500 {

// 			resp.Body.Close()

// 			if attempt == maxRetries {
// 				return nil, fmt.Errorf(
// 					"product service returned status %d",
// 					resp.StatusCode,
// 				)
// 			}

// 			if err := waitForRetry(ctx, attempt); err != nil {
// 				return nil, err
// 			}

// 			continue
// 		}

// 		// 4xx → don't retry
// 		resp.Body.Close()

// 		return nil, fmt.Errorf(
// 			"product service returned status %d",
// 			resp.StatusCode,
// 		)
// 	}

// 	return nil, fmt.Errorf("product service request failed")
// }

// func isRetryableError(err error) bool {
// 	return true
// }
// func waitForRetry(ctx context.Context, attempt int) error {
// 	delay := time.Duration(100*(1<<attempt)) * time.Millisecond
// 	timer := time.NewTimer(delay)
// 	defer timer.Stop()
// 	select {
// 	case <-timer.C:
// 		return nil
// 	case <-ctx.Done():
// 		return ctx.Err()
// 	}

// }

// func (c *ProductClient) BookStock(
// 	ctx context.Context,
// 	productID string,
// 	stock int,
// ) error {

// 	url := c.baseUrl + "/product/" + productID + "/stocks"

// 	body, err := json.Marshal(map[string]int{
// 		"stock": int(stock),
// 	})
// 	if err != nil {
// 		return err
// 	}

// 	req, err := http.NewRequestWithContext(
// 		ctx,
// 		http.MethodPatch,
// 		url,
// 		bytes.NewBuffer(body),
// 	)
// 	if err != nil {
// 		return err
// 	}
// 	req.Header.Set("Content-Type", "application/json")
// 	resp, err := c.client.Do(req)
// 	if err != nil {
// 		return err
// 	}

// 	defer resp.Body.Close()

// 	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
// 		return fmt.Errorf(
// 			"product service returned status %d",
// 			resp.StatusCode,
// 		)
// 	}

// 	return nil
// }

// func (c *ProductClient) ReleaseStock(
// 	ctx context.Context,
// 	productID string,
// 	stock int,
// ) error {

// 	url := c.baseUrl + "/product/" + productID + "/relstocks"
// 	log.Println(url)
// 	body, err := json.Marshal(map[string]int{
// 		"stock": int(stock),
// 	})
// 	if err != nil {
// 		return err
// 	}

// 	req, err := http.NewRequestWithContext(
// 		ctx,
// 		http.MethodPatch,
// 		url,
// 		bytes.NewBuffer(body),
// 	)
// 	if err != nil {
// 		return err
// 	}
// 	req.Header.Set("Content-Type", "application/json")
// 	resp, err := c.client.Do(req)
// 	if err != nil {
// 		return err
// 	}

// 	defer resp.Body.Close()

// 	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
// 		return fmt.Errorf(
// 			"product service returned status %d",
// 			resp.StatusCode,
// 		)
// 	}

// 	return nil
// }
