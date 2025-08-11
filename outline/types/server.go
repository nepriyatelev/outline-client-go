package types

// ServerInfo — ответ GET /server
type ServerInfoResponse struct {
	Name                 string `json:"name", yaml:"name", example:"My Server"`
	ServerID             string `json:"serverId", yaml:"serverId", example:"40f1b4a3-5c82-45f4-80a6-a25cf36734d3"`
	MetricsEnabled       bool   `json:"metricsEnabled", yaml:"metricsEnabled", example:"true"`
	CreatedTimestampMs   int64  `json:"createdTimestampMs", yaml:"createdTimestampMs", example:"1536613192052"`
	Version              string `json:"version", yaml:"version", example:"1.0.0"`
	PortForNewAccessKeys int    `json:"portForNewAccessKeys", yaml:"portForNewAccessKeys", example:"1234"`
	HostnameForAccessKeys string `json:"hostnameForAccessKeys", yaml:"hostnameForAccessKeys"`
}

// HostnameRequest — тело PUT /server/hostname-for-access-keys
type HostnameRequest struct {
	Hostname string `json:"hostname" yaml:"hostname" example:"myserver.example.com"`
}

// NameRequest — тело PUT /name
type NameRequest struct {
    Name string `json:"name" yaml:"name" example:"My Server"`
}

// MetricsEnabledRequest — тело PUT /metrics/enabled
type MetricsEnabledRequest struct {
    Enabled bool `json:"metricsEnabled" yaml:"metricsEnabled" example:"true"`
}

// MetricsEnabledResponse — ответ GET /metrics/enabled
type MetricsEnabledResponse struct {
    Enabled bool `json:"metricsEnabled" yaml:"metricsEnabled" example:"true"`
}