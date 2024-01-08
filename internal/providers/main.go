package providers

type Provider interface {
	Upload() error
	Download() error
}
