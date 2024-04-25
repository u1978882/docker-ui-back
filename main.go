package main

import (
    "log"
    "os"
	"net/http"

	"github.com/labstack/echo/v5"
    "github.com/pocketbase/pocketbase"
    "github.com/pocketbase/pocketbase/apis"
    "github.com/pocketbase/pocketbase/core"
)

func main() {
    app := pocketbase.New()

    // serves static files from the provided public dir (if exists)
    app.OnBeforeServe().Add(func(e *core.ServeEvent) error {
        e.Router.GET("/*", apis.StaticDirectoryHandler(os.DirFS("./pb_public"), false))
        return nil
    })

    if err := app.Start(); err != nil {
        log.Fatal(err)
    }

	app.OnBeforeServe().Add(func(e *core.ServeEvent) error {
		e.Router.GET("/api/images", func(c echo.Context) error {
	
			//[{name: "alpine linux"}, {name: "Ubuntu 16.04"}, {name: "Minecraft Server"}, {name: "Ftp server"}]
			return c.JSON(http.StatusOK, map[string]string{"images": "test"})
		}, /* optional middlewares */)
	
		return nil
	})

}