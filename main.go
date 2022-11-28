package main

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/sirupsen/logrus"
)

type M map[string]interface{}

func main() {
	tmpl := template.Must(template.ParseGlob("./*.html"))

	e := echo.New()

	const CSRFTokenHeader = "X-CSRF-Token"
	const CSRFKey = "csrf"

	e.Use(middleware.CSRFWithConfig(middleware.CSRFConfig{
		TokenLookup:  "header:" + CSRFTokenHeader,
		ContextKey:   CSRFKey,
		CookieMaxAge: 30,
	}))

	e.GET("/index", func(c echo.Context) error {
		data := make(M)
		data[CSRFKey] = c.Get(CSRFKey)
		logrus.Info(data)
		return tmpl.Execute(c.Response(), data)
	})

	e.POST("/sayhello", func(c echo.Context) error {
		data := make(M)
		if err := c.Bind(&data); err != nil {
			logrus.Error(err.Error())
			return err
		}

		message := fmt.Sprintf("hello %s", data["name"])
		return c.JSON(http.StatusOK, message)
	})

	e.Logger.Fatal(e.Start(":9000"))
}
