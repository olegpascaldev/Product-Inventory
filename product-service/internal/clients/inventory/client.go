package inventory

import (
	"context"
	pb "product-service/pkg/protos/inventory"

	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Client struct {
	conn   *grpc.ClientConn
	client pb.InventoryServiceClient
}

// NewClient создает новый gRPC-клиент для службы инвентаризации.
func NewClient(addr string) (*Client, error) {
	conn, err := grpc.Dial(addr, grpc.WithInsecure()) // For production, use TLS
	if err != nil {
		return nil, err
	}
	return &Client{
		conn:   conn,
		client: pb.NewInventoryServiceClient(conn),
	}, nil
}

// Закрывает gRPC-соединение.
func (c *Client) Close() error {
	return c.conn.Close()
}

// Функция GetStock получает текущее количество товара на складе.
func (c *Client) GetStock(ctx context.Context, productID uuid.UUID) (int32, error) {
	resp, err := c.client.GetStock(ctx, &pb.GetStockRequest{ProductId: productID.String()})
	if err != nil {
		// Если товар не найден, считать его наличие нулевым (или вернуть конкретную ошибку)
		if status.Code(err) == codes.NotFound {
			return 0, nil
		}
		return 0, err
	}
	return resp.Quantity, nil
}
