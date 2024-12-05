package shared

import (
	"net/http"
	"testing"
	"time"

	"github.com/caarlos0/log"
	"github.com/carlmjohnson/requests"
	"github.com/carlmjohnson/requests/reqtest"
)

func HTTPTransport() http.RoundTripper {
	baseTrans := http.DefaultClient.Transport
	if testing.Testing() {
		return reqtest.Record(baseTrans, "fixtures")
	}
	logger := func(req *http.Request, res *http.Response, err error, d time.Duration) {
		log.Infof("method=%q url=%q err=%v status=%q duration=%v\n",
			req.Method, req.URL, err, res.Status, d.Round(1*time.Second))
	}

	return requests.LogTransport(baseTrans, logger)
}
