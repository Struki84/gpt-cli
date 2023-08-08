package metaphor

type MetaphorOptions func(*MetaphorSearch)

func WithNumResults(numResults int) MetaphorOptions {
	return func(metaphor *MetaphorSearch) {
		metaphor.client.RequestBody.NumResults = numResults
	}
}

func WithIncludeDomains(includeDomains []string) MetaphorOptions {
	return func(metaphor *MetaphorSearch) {
		metaphor.client.RequestBody.IncludeDomains = includeDomains
	}
}

func WithExcludeDomains(excludeDomains []string) MetaphorOptions {
	return func(metaphor *MetaphorSearch) {
		metaphor.client.RequestBody.ExcludeDomains = excludeDomains
	}
}

func WithStartCrawlDate(startCrawlDate string) MetaphorOptions {
	return func(metaphor *MetaphorSearch) {
		metaphor.client.RequestBody.StartCrawlDate = startCrawlDate
	}
}

func WithEndCrawlDate(endCrawlDate string) MetaphorOptions {
	return func(metaphor *MetaphorSearch) {
		metaphor.client.RequestBody.EndCrawlDate = endCrawlDate
	}
}

func WithStartPublishedDate(startPublishedDate string) MetaphorOptions {
	return func(metaphor *MetaphorSearch) {
		metaphor.client.RequestBody.StartPublishedDate = startPublishedDate
	}
}

func WithEndPublishedDate(endPublishedDate string) MetaphorOptions {
	return func(metaphor *MetaphorSearch) {
		metaphor.client.RequestBody.EndPublishedDate = endPublishedDate
	}
}

func WithAutoprompt(useAutoprompt bool) MetaphorOptions {
	return func(metaphor *MetaphorSearch) {
		metaphor.client.RequestBody.UseAutoprompt = useAutoprompt
	}
}

func WithType(type_ string) MetaphorOptions {
	return func(metaphor *MetaphorSearch) {
		metaphor.client.RequestBody.Type = type_
	}
}
