package recounter

type Client struct {
	output chan string
}

func newClient(output chan string) *Client {
	return &Client{
		output: output,
	}
}

func (c *Client) TryRecount(releaseID string) {
	c.output <- releaseID
}
