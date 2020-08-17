package concurrency

//WebsiteChecker is a the type of a function that checks a website url
type WebsiteChecker func(string) bool

//Result gives the result to a request which receives a string and returns a bool
type result struct {
	string
	bool
}

//CheckWebsites takes in a WebsiteChecker and []string of urls. It checks each url using the checker and returns a map of the results
func CheckWebsites(wc WebsiteChecker, urls []string) map[string]bool {
	results := make(map[string]bool)
	resultChannel := make(chan result)

	for _, url := range urls {
		go func(url string) {
			resultChannel <- result{url, wc(url)}
		}(url)
	}

	for i := 0; i < len(urls); i++ {
		res := <-resultChannel
		results[res.string] = res.bool
	}

	return results
}
