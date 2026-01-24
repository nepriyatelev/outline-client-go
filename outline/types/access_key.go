// Package types defines data structures for Outline VPN access keys,
// server information, metrics, and related API requests and responses.
package types

// AccessKey represents an access key for VPN connection.
type AccessKey struct {
	ID        string `json:"id"`        // ID is the unique identifier of the access key.
	Name      string `json:"name"`      // Name is the human-readable name of the access key.
	Password  string `json:"password"`  // Password is the password used for client connection.
	Port      int    `json:"port"`      // Port is the TCP/UDP port on which the access key is available.
	Method    string `json:"method"`    // Method is the encryption method used.
	AccessURL string `json:"accessUrl"` // AccessURL is the URL for accessing the key.
}

// CreateAccessKey represents a request to create a new access key.
type CreateAccessKey struct {
	Method   string `json:"method"`             // Method is the required encryption algorithm that defines the cryptographic method for protecting traffic. Example: "aes-192-gcm".
	Name     string `json:"name,omitempty"`     // Name is the optional human-readable name for identifying this key, used for organization and management. Example: "Work Laptop". If not specified, the server may assign a default name.
	Password string `json:"password,omitempty"` // Password is the optional password used for client connection, used together with the encryption method. Example: "8iu8V8EeoFVpwQvQeS9wiD". If not specified, the server will generate a secure password.
	Port     uint16 `json:"port,omitempty"`     // Port is the optional TCP/UDP port on which this access key will be available. Example: 8388. If not specified, uses portForNewAccessKeys from server configuration.
	Limit    *Limit `json:"limit,omitempty"`    // Limit is the optional data transfer limit specifying the maximum number of bytes that can be transferred through this access key. After reaching the limit, traffic may be blocked. Example: {"bytes": 10000} where bytes is the maximum number of bytes (0 means no limit).
}
