package outline

import (
	"context"
	"time"

	"github.com/nepriyatelev/outline-client-go/outline/types"
)

type ClientOutline interface {
	// Server Methods
	// Server Information
	GetServerInfo(ctx context.Context) (*types.ServerInfoResponse, error)
	// Server Configuration
	UpdateServerHostname(ctx context.Context, hostnameOrIP string) error
	UpdatePortNewAccessKeys(ctx context.Context, port uint16) error
	UpdateServerName(ctx context.Context, name string) error
	GetMetricsEnabled(ctx context.Context) (*types.MetricsEnabled, error)
	UpdateMetricsEnabled(ctx context.Context, enabled bool) error
	// Data Limits (Server-wide)
	UpdateKeyLimitBytes(ctx context.Context, bytes uint64) error
	DeleteKeyLimitBytes(ctx context.Context) error

	// Access Key Methods
	// CRUD Operations
	CreateAccessKey(
		ctx context.Context, createAccessKey *types.CreateAccessKey) (*types.AccessKey, error)
	GetAccessKeys(ctx context.Context) ([]*types.AccessKey, error)
	GetAccessKey(ctx context.Context, accessKeyID string) (*types.AccessKey, error)
	UpdateAccessKey(
		ctx context.Context, accessKeyID string, updateAccessKey *types.AccessKey) (
		*types.AccessKey, error)
	DeleteAccessKey(ctx context.Context, accessKeyID string) error
	// Management Operations
	UpdateNameAccessKey(ctx context.Context, accessKeyID, newName string) error
	UpdateDataLimitAccessKey(ctx context.Context, accessKeyID string, bytes uint64) error
	DeleteDataLimitAccessKey(ctx context.Context, accessKeyID string) error

	// Transfer Metrics
	GetMetricsTransfer(ctx context.Context) (*types.MetricsTransfer, error)

	// Experimental Metrics
	GetExperimentalMetrics(
		ctx context.Context, since time.Duration) (*types.ExperimentalMetricsResponse, error)
}
