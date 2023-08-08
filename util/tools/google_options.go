package tools

type GoogleOptions func(s *Client)

func WithApiKey(key string) GoogleOptions {
	return func(s *Client) {
		s.apiKey = key
	}
}

func WithMaxResults(maxResults int) GoogleOptions {
	return func(s *Client) {
		s.maxResults = maxResults
	}
}
