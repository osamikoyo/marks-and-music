package fetcher

type FetcherClient struct{
	queryChan chan string
}

func newFetcherClient(qchan chan string) *FetcherClient {
	return &FetcherClient{
		queryChan: qchan,
	}
}

func (fc *FetcherClient) AsyncFetch(query string) {
	fc.queryChan <- query
}