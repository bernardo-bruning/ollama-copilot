package adapters

type HttpClient interface {
	Post(url string, req any, resp any) error
}
