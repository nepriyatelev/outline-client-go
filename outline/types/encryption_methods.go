package types

import "slices"

// EncryptionMethods defines the supported encryption algorithms for Outline Server.
// Outline uses AEAD ciphers (Authenticated Encryption with Associated Data)
// to ensure confidentiality, integrity, and authenticity of data.
const (
	// MethodAES128GCM is AES-128 with authentication (primary recommended).
	// Required method in modern Shadowsocks implementations.
	// Excellent performance with hardware acceleration on modern equipment.
	MethodAES128GCM = "aes-128-gcm"

	// MethodAES256GCM is AES-256 with authentication (enhanced security).
	// Maximum level of protection, recommended for critical data.
	// Works more efficiently with hardware AES acceleration (AES-NI).
	MethodAES256GCM = "aes-256-gcm"

	// MethodChaCha20IETFPoly1305 is ChaCha20 with authentication (universal).
	// Required method according to Shadowsocks specification (SIP004).
	// Best choice for systems without hardware AES acceleration.
	// Used as the default method in Outline Server.
	MethodChaCha20IETFPoly1305 = "chacha20-ietf-poly1305"
)

// ValidEncryptionMethods lists all supported encryption methods
// in order of recommendation (from most preferred to least preferred).
var ValidEncryptionMethods = []string{
	MethodChaCha20IETFPoly1305, // Default in Outline Server
	MethodAES128GCM,            // Recommended for modern equipment
	MethodAES256GCM,            // For those requiring maximum security
}

// IsValidEncryptionMethod reports whether the given encryption method is supported.
// It returns true if the method is in the list of supported methods.
func IsValidEncryptionMethod(method string) bool {
	return slices.Contains(ValidEncryptionMethods, method)
}

// GetDefaultEncryptionMethod returns the default encryption method used in Outline Server.
func GetDefaultEncryptionMethod() string {
	return MethodChaCha20IETFPoly1305
}
