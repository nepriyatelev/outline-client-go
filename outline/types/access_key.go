package types

type AccessKey struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Password  string `json:"password"`
	Port      int    `json:"port"`
	Method    string `json:"method"`
	AccessURL string `json:"accessUrl"`
}

// CreateAccessKeyRequest - запрос для создания ключа доступа
type CreateAccessKey struct {
	/*
		Обязательно
		method (string) - Алгоритм шифрования
		Определяет криптографический метод для защиты трафика.
		Пример: "method": "aes-192-gcm"
	*/
	Method string `json:"method"`
	/*
		Опционально
		name (string) - Имя ключа доступа
		Человекочитаемое имя для идентификации этого ключа.
		Используется для организации и управления ключами.
		Пример: "name": "Work Laptop"
		Примечание: Если не указано, сервер может присвоить имя по умолчанию.
	*/
	Name string `json:"name,omitempty"`
	/*
		Опционально
		password (string) - Пароль для подключения
		Пароль, который будет использоваться при подключении клиентом.
		Используется вместе с методом шифрования.
		Пример: "password": "8iu8V8EeoFVpwQvQeS9wiD"
		Примечание: Если не указано, сервер сгенерирует безопасный пароль.
	*/
	Password string `json:"password,omitempty"`
	/*
		Опционально
		port (integer) - Порт для прослушивания
		TCP/UDP порт, на котором будет доступен этот ключ доступа.
		Пример: "port": 8388
		Примечание: Если не указано, используется portForNewAccessKeys из конфигурации сервера
	*/
	Port uint16 `json:"port,omitempty"`
	/*
		Опционально
		limit (object) - Лимит передачи данных
		Максимальное количество байт, которые могут быть переданы через этот ключ доступа.
		После достижения лимита трафик может быть заблокирован.
		Пример:
		"limit": {
			"bytes": 10000
		}

		bytes (integer) - Максимальное количество байт (0 означает отсутствие лимита)
	*/
	Limit *Limit `json:"limit,omitempty"`
}
