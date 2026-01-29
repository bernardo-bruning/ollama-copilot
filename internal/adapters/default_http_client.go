package adapters

func NewDefaultHttpClient(baseURL string, token string) HttpClient {
	headers := map[string]string{
		"Content-Type": "application/json",
	}
	if token != "" {
		headers["Authorization"] = "Bearer " + token
	}
	return NewHttpClient(baseURL, headers)
}