package outline

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

func maskSecretPath(raw, secret string) string {
	if secret == "" {
		return raw
	}

	return strings.ReplaceAll(raw, "/"+secret+"/", "/*****/")
}

func formatDuration(d time.Duration) string {
	// Определяем знак и работаем с абсолютным значением
	sign := ""
	if d < 0 {
		sign = "-"
		d = -d
	}

	// Часы
	h := int64(d.Hours())
	if h != 0 {
		return fmt.Sprintf("%s%dh", sign, h)
	}
	// Минуты
	m := int64(d.Minutes())
	if m != 0 {
		return fmt.Sprintf("%s%dm", sign, m)
	}
	// Секунды
	s := int64(d.Seconds())
	if s != 0 {
		return fmt.Sprintf("%s%ds", sign, s)
	}
	// Если всё равно ноль, игнорируем знак и возвращаем "0s"
	return "0s"
}

func setIDInPath(path string, id string) string {
	replacedPath := strings.Replace(path, "{id}", id, 1)
	return replacedPath
}

// unmarshalJSONWithError декодирует JSON структуру и возвращает указатель
// Используйте для одиночных структур
func unmarshalJSONWithError[T any](data []byte) (*T, error) {
	target := new(T)
	if err := unmarshalWithErrorInternal(data, target, fmt.Sprintf("%T", target)); err != nil {
		return nil, err
	}
	return target, nil
}

// unmarshalJSONSliceOfPointersWithError для слайса указателей → []*T
func unmarshalJSONSliceOfPointersWithError[T any](data []byte) ([]*T, error) {
	var target []*T
	if err := unmarshalWithErrorInternal(data, &target, fmt.Sprintf("%T", target)); err != nil {
		return nil, err
	}
	return target, nil
}

// unmarshalWithErrorInternal общая логика для всех unmarshal операций
func unmarshalWithErrorInternal(data []byte, target interface{}, typeStr string) error {
	if len(data) == 0 {
		return &UnmarshalError{
			Data: data,
			Type: typeStr,
			Err:  fmt.Errorf("empty body"),
		}
	}

	if err := json.Unmarshal(data, target); err != nil {
		return &UnmarshalError{
			Data: data,
			Type: typeStr,
			Err:  err,
		}
	}
	return nil
}
