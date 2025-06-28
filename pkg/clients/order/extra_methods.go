package order

import (
    "context"
    "fmt"

    "go.uber.org/zap"

    orderv1 "github.com/leonvanderhaeghen/stockplatform/services/orderSvc/api/gen/go/proto/order/v1"
)

// AddTrackingCode adds a tracking code to an order
func (c *Client) AddTrackingCode(ctx context.Context, req *orderv1.AddTrackingCodeRequest) (*orderv1.AddTrackingCodeResponse, error) {
    c.logger.Debug("Adding tracking code", zap.String("order_id", req.OrderId))

    resp, err := c.client.AddTrackingCode(ctx, req)
    if err != nil {
        c.logger.Error("Failed to add tracking code", zap.Error(err))
        return nil, fmt.Errorf("failed to add tracking code: %w", err)
    }

    return resp, nil
}
