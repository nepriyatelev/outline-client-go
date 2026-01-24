package types

// ServerInfoResponse represents the response containing information about the Outline server.
type ServerInfoResponse struct {
	Name                  string  `json:"name"`                  // Name is the human-readable name of the server.
	ServerID              string  `json:"serverId"`              // ServerID is the unique identifier of the server.
	MetricsEnabled        bool    `json:"metricsEnabled"`        // MetricsEnabled indicates whether metrics collection is enabled.
	CreatedTimestampMs    float64 `json:"createdTimestampMs"`    // CreatedTimestampMs is the creation timestamp in milliseconds since epoch.
	Version               string  `json:"version"`               // Version is the version of the Outline server software.
	PortForNewAccessKeys  int     `json:"portForNewAccessKeys"`  // PortForNewAccessKeys is the default port for new access keys.
	HostnameForAccessKeys string  `json:"hostnameForAccessKeys"` // HostnameForAccessKeys is the hostname used for access keys.
}
