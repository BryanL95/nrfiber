package nrfiber

import (
	"net/http"
	"net/url"

	"github.com/gofiber/fiber/v2"
	"github.com/newrelic/go-agent/v3/newrelic"
)

func Middleware(app *newrelic.Application) fiber.Handler {
	return func(c *fiber.Ctx) error {
		if app != nil {
			txn := app.StartTransaction(c.Method() + " " + c.Path())
			defer func() {
				err := recover()
				if err != nil {
					switch err := err.(type) {
					case error:
						txn.NoticeError(err)
					}
				} else {
					txn.SetWebRequest(newrelic.WebRequest{
						Header: http.Header{
							"Response-Headers": []string{
								c.Get("Date"),
							},
						},
						URL:       &url.URL{Path: c.Path()},
						Method:    c.Method(),
						Transport: newrelic.TransportHTTP,
					})
				}
				txn.End()
			}()
			c.Context().SetUserValue("__newrelic_transaction__", txn)
		}
		c.Next()
		return nil
	}
}
