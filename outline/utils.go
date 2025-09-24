package outline

import (
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


