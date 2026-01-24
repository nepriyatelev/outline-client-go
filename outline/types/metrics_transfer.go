package types

// MetricsTransfer represents metrics for data transfer grouped by user ID.
type MetricsTransfer struct {
	BytesTransferredByUserID map[string]int64 `json:"bytesTransferredByUserId"` // BytesTransferredByUserID maps user IDs to the number of bytes transferred by each user.
}
