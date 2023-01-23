package votifier

// TokenProvider provides a token for a given vote service.
// The vote service is the name of the service the user is voting from
// and sent by the vote service itself (e.g. "minecraft-serverlist.net").
type TokenProvider interface {
	// Token returns the token for a service.
	Token(service string) string
}

// TokenProviderFunc is a function that implements TokenProvider.
type TokenProviderFunc func(service string) string

// Token implements TokenProvider.
func (f TokenProviderFunc) Token(service string) string {
	return f(service)
}

// StaticTokenProvider returns the same token for every request.
func StaticTokenProvider(token string) TokenProvider {
	return TokenProviderFunc(func(service string) string {
		return token
	})
}
