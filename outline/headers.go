package outline

type Headers map[string]string

func DefaultHeaders() Headers {
	return Headers{
		"Content-Type": "application/json",
		"Accept":       "application/json",
	}
}
