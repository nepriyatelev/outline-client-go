package types

// ExperimentalMetricsResponse represents the response containing experimental metrics
// for the server and all access keys.
type ExperimentalMetricsResponse struct {
	Server     ServerMetrics      `json:"server"`     // Server contains metrics for the Outline server.
	AccessKeys []AccessKeyMetrics `json:"accessKeys"` // AccessKeys contains metrics for each access key.
}

// ServerMetrics represents metrics collected for the Outline server.
type ServerMetrics struct {
	Locations []LocationMetrics `json:"locations"` // Locations contains metrics grouped by geographic location.
}

// TimeMetric represents a time duration in seconds.
type TimeMetric struct {
	Seconds float64 `json:"seconds"` // Seconds is the duration in seconds.
}

// DataMetric represents an amount of data in bytes.
type DataMetric struct {
	Bytes float64 `json:"bytes"` // Bytes is the amount of data in bytes.
}

// BandwidthMetrics represents bandwidth usage metrics including current and peak values.
type BandwidthMetrics struct {
	Current BandwidthPoint `json:"current"` // Current is the current bandwidth usage at the time of measurement.
	Peak    BandwidthPoint `json:"peak"`    // Peak is the highest bandwidth usage recorded.
}

// BandwidthPoint represents a bandwidth measurement at a specific timestamp.
type BandwidthPoint struct {
	Data      DataMetric `json:"data"`      // Data is the amount of data transferred in this measurement.
	Timestamp int64      `json:"timestamp"` // Timestamp is the Unix timestamp when the measurement was taken.
}

// LocationMetrics represents metrics for a specific geographic location.
type LocationMetrics struct {
	Location        string     `json:"location"`        // Location is the geographic location identifier.
	ASN             *int64     `json:"asn"`             // ASN is the Autonomous System Number, if available.
	ASOrg           *string    `json:"asOrg"`           // ASOrg is the Autonomous System organization name, if available.
	DataTransferred DataMetric `json:"dataTransferred"` // DataTransferred is the amount of data transferred from this location.
	TunnelTime      TimeMetric `json:"tunnelTime"`      // TunnelTime is the total tunnel time for connections from this location.
}

// AccessKeyMetrics represents metrics for a specific access key.
type AccessKeyMetrics struct {
	AccessKeyID     int64             `json:"accessKeyId"`     // AccessKeyID is the unique identifier of the access key.
	TunnelTime      TimeMetric        `json:"tunnelTime"`      // TunnelTime is the total time the access key has been used for tunneling.
	DataTransferred DataMetric        `json:"dataTransferred"` // DataTransferred is the total amount of data transferred using this access key.
	Connection      ConnectionMetrics `json:"connection"`      // Connection contains connection-related metrics for this access key.
}

// ConnectionMetrics represents connection-related metrics for an access key.
type ConnectionMetrics struct {
	LastTrafficSeen int64           `json:"lastTrafficSeen"` // LastTrafficSeen is the Unix timestamp of the last traffic seen for this access key.
	PeakDeviceCount PeakDeviceCount `json:"peakDeviceCount"` // PeakDeviceCount is the peak number of devices connected simultaneously.
}

// PeakDeviceCount represents the peak number of devices connected at a specific time.
type PeakDeviceCount struct {
	Data      int64 `json:"data"`      // Data is the number of devices connected.
	Timestamp int64 `json:"timestamp"` // Timestamp is the Unix timestamp when this peak was recorded.
}
