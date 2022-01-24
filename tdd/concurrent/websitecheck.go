package concurrent

type WebsiteChecker func(string) bool

type result struct {
	url string
	res bool
}

func CheckWebsites(chk WebsiteChecker, urls []string) map[string]bool {
	res := make(map[string]bool)
	resChan := make(chan result)
	for _, url := range urls {
		go func(u string) {
			resChan <- result{u, chk(u)}
		}(url)
	}

	for i := 0; i < len(urls); i++ {
		r := <-resChan
		res[r.url] = r.res
	}
	return res
}
