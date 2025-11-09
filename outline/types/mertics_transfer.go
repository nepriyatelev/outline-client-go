package types

type MetricsTransfer struct {
	BytesTransferredByUserID map[string]int64 `json:"bytesTransferredByUserId"`
}