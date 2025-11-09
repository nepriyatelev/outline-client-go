package types

type ExperimentalMetricsResponse struct {
	Server     ServerMetrics      `json:"server"`
	AccessKeys []AccessKeyMetrics `json:"accessKeys"`
}

type ServerMetrics struct {
	TunnelTime      TimeMetric        `json:"tunnelTime"`
	DataTransferred DataMetric        `json:"dataTransferred"`
	Bandwidth       BandwidthMetrics  `json:"bandwidth"`
	Locations       []LocationMetrics `json:"locations"`
}

type TimeMetric struct {
	Seconds float64 `json:"seconds"`
}

type DataMetric struct {
	Bytes float64 `json:"bytes"`
}

type BandwidthMetrics struct {
	Current BandwidthPoint `json:"current"`
	Peak    BandwidthPoint `json:"peak"`
}

type BandwidthPoint struct {
	Data      DataMetric `json:"data"`
	Timestamp int64      `json:"timestamp"`
}

type LocationMetrics struct {
	Location        string     `json:"location"`
	ASN             *int64     `json:"asn"`
	ASOrg           *string    `json:"asOrg"`
	DataTransferred DataMetric `json:"dataTransferred"`
	TunnelTime      TimeMetric `json:"tunnelTime"`
}

type AccessKeyMetrics struct {
	AccessKeyID     int64             `json:"accessKeyId"`
	TunnelTime      TimeMetric        `json:"tunnelTime"`
	DataTransferred DataMetric        `json:"dataTransferred"`
	Connection      ConnectionMetrics `json:"connection"`
}

type ConnectionMetrics struct {
	LastTrafficSeen int64           `json:"lastTrafficSeen"`
	PeakDeviceCount PeakDeviceCount `json:"peakDeviceCount"`
}

type PeakDeviceCount struct {
	Data      int64 `json:"data"`
	Timestamp int64 `json:"timestamp"`
}
