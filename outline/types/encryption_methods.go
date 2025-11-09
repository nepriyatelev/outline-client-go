package types

import "slices"

// Encryption methods (поддерживаемые Outline Server алгоритмы шифрования)
// Outline использует AEAD шифры (Authenticated Encryption with Associated Data)
// для обеспечения конфиденциальности, целостности и аутентичности данных
const (
	// MethodAES128GCM - AES-128 с аутентификацией (основной рекомендуемый)
	// Обязательный метод в современных реализациях Shadowsocks
	// Отличная производительность с аппаратным ускорением на современном оборудовании
	MethodAES128GCM = "aes-128-gcm"

	// MethodAES256GCM - AES-256 с аутентификацией (повышенная безопасность)
	// Максимальный уровень защиты, рекомендуется для критичных данных
	// Работает эффективнее при наличии аппаратного ускорения AES (AES-NI)
	MethodAES256GCM = "aes-256-gcm"

	// MethodChaCha20IETFPoly1305 - ChaCha20 с аутентификацией (универсальный)
	// Обязательный метод согласно спецификации Shadowsocks (SIP004)
	// Лучший выбор для систем без аппаратного ускорения AES
	// Используется как метод по умолчанию в Outline Server
	MethodChaCha20IETFPoly1305 = "chacha20-ietf-poly1305"
)

// ValidEncryptionMethods возвращает список всех поддерживаемых методов шифрования
// в порядке рекомендации (от более предпочтительного к менее предпочтительному)
var ValidEncryptionMethods = []string{
	MethodChaCha20IETFPoly1305, // Default в Outline Server
	MethodAES128GCM,            // Рекомендуемый для современного оборудования
	MethodAES256GCM,            // Для требующих максимальной безопасности
}

// IsValidEncryptionMethod проверяет, поддерживается ли переданный метод шифрования
// Возвращает true если метод есть в списке поддерживаемых методов
func IsValidEncryptionMethod(method string) bool {
	return slices.Contains(ValidEncryptionMethods, method)
}

// GetDefaultEncryptionMethod возвращает метод шифрования по умолчанию в Outline Server
func GetDefaultEncryptionMethod() string {
	return MethodChaCha20IETFPoly1305
}
