package types

type ServerInfoResponse struct {
	Name                  string  `json:"name"`
	ServerID              string  `json:"serverId"`
	MetricsEnabled        bool    `json:"metricsEnabled"`
	CreatedTimestampMs    float64 `json:"createdTimestampMs"`
	Version               string  `json:"version"`
	PortForNewAccessKeys  int     `json:"portForNewAccessKeys"`
	HostnameForAccessKeys string  `json:"hostnameForAccessKeys"`
}
