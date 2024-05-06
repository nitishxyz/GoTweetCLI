package twitter

import (
	"context"
	"fmt"
	"net/http"
)

func (c *Client) ListenAndServe(port string) error {
	mux := http.NewServeMux()
	mux.HandleFunc("/token", c.tokenHandler)

	c.server = &http.Server{
		Addr:    port,
		Handler: mux,
	}

	return c.server.ListenAndServe()
}

func (c *Client) StopServer() error {
	if c.server == nil {
		return fmt.Errorf("server is not running")
	}
	return c.server.Shutdown(context.Background())
}
