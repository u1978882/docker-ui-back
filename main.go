package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
	"golang.org/x/crypto/ssh"
)

func executarComanda(serverAddress, username, password, command string) (string, error) {
	// Configuración de la conexión SSH
	sshConfig := &ssh.ClientConfig{
		User: username,
		Auth: []ssh.AuthMethod{
			ssh.Password(password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	// Realizar la conexión SSH
	client, err := ssh.Dial("tcp", serverAddress, sshConfig)
	if err != nil {
		return "", fmt.Errorf("Error al conectar al servidor SSH: %v", err)
	}
	defer client.Close()

	// Crear una nueva sesión SSH
	session, err := client.NewSession()
	if err != nil {
		return "", fmt.Errorf("Error al crear la sesión SSH: %v", err)
	}
	defer session.Close()

	// Ejecutar el comando en la sesión SSH
	output, err := session.CombinedOutput(command)
	if err != nil {
		return "", fmt.Errorf("Error al ejecutar el comando en el servidor remoto: %v", err)
	}

	return string(output), nil
}

func main() {
	app := pocketbase.New()

	app.OnBeforeServe().Add(func(e *core.ServeEvent) error {
		e.Router.GET("/functions/:server/container/:containerId/start", func(c echo.Context) error {
			server := c.PathParam("server")
			container := c.PathParam("containerId")

			record, err := app.Dao().FindRecordById("server", server)
			if err != nil {
				return c.JSON(http.StatusForbidden, string("[]"))
			}

			serverAddress := fmt.Sprintf("%s:%s", record.GetString("ip"), strconv.Itoa(record.GetInt("port")))
			command := "docker start " + container

			output, err := executarComanda(serverAddress, record.GetString("username"), record.GetString("pass"), command)
			if err != nil {
				fmt.Printf("Error al ejecutar el comando SSH: %v\n", err)
				return c.JSON(http.StatusForbidden, "[]")
			}

			var sortidaFinal = string(output)
			if len(sortidaFinal) > 0 {
				sortidaFinal = sortidaFinal[:len(sortidaFinal)-2]
			}

			fmt.Println("Sortida final:")
			fmt.Println(sortidaFinal)

			return c.JSON(http.StatusOK, string("{\"resultat\": "+sortidaFinal+"}"))
		} /* optional middlewares */)

		return nil
	})

	app.OnBeforeServe().Add(func(e *core.ServeEvent) error {
		e.Router.GET("/functions/:server/container/:containerId/stop", func(c echo.Context) error {
			server := c.PathParam("server")
			container := c.PathParam("containerId")

			record, err := app.Dao().FindRecordById("server", server)
			if err != nil {
				return c.JSON(http.StatusForbidden, string("[]"))
			}

			serverAddress := fmt.Sprintf("%s:%s", record.GetString("ip"), strconv.Itoa(record.GetInt("port")))

			command := "docker stop " + container

			output, err := executarComanda(serverAddress, record.GetString("username"), record.GetString("pass"), command)
			if err != nil {
				fmt.Printf("Error al ejecutar el comando SSH: %v\n", err)
				return c.JSON(http.StatusForbidden, "[]")
			}

			var sortidaFinal = string(output)
			if len(sortidaFinal) > 0 {
				sortidaFinal = sortidaFinal[:len(sortidaFinal)-2]
			}

			fmt.Println("Sortida final:")
			fmt.Println(sortidaFinal)

			return c.JSON(http.StatusOK, string("{\"resultat\": "+sortidaFinal+"}"))
		} /* optional middlewares */)

		return nil
	})

	app.OnBeforeServe().Add(func(e *core.ServeEvent) error {
		e.Router.GET("/functions/images/:server", func(c echo.Context) error {
			server := c.PathParam("server")

			record, err := app.Dao().FindRecordById("server", server)
			if err != nil {
				return c.JSON(http.StatusForbidden, string("[]"))
			}

			serverAddress := fmt.Sprintf("%s:%s", record.GetString("ip"), strconv.Itoa(record.GetInt("port")))

			command := "docker images --format='{{json .}},'"

			output, err := executarComanda(serverAddress, record.GetString("username"), record.GetString("pass"), command)
			if err != nil {
				fmt.Printf("Error al ejecutar el comando SSH: %v\n", err)
				return c.JSON(http.StatusForbidden, "[]")
			}

			var sortidaFinal = string(output)
			if len(sortidaFinal) > 0 {
				sortidaFinal = sortidaFinal[:len(sortidaFinal)-2]
			}

			fmt.Println("Sortida final:")
			fmt.Println(sortidaFinal)

			return c.JSON(http.StatusOK, string("{\"images\": ["+sortidaFinal+"]}"))
		} /* optional middlewares */)

		return nil
	})

	app.OnBeforeServe().Add(func(e *core.ServeEvent) error {
		e.Router.GET("/functions/containers/:server", func(c echo.Context) error {
			server := c.PathParam("server")

			record, err := app.Dao().FindRecordById("server", server)
			if err != nil {
				return c.JSON(http.StatusForbidden, string("[]"))
			}

			serverAddress := fmt.Sprintf("%s:%s", record.GetString("ip"), strconv.Itoa(record.GetInt("port")))

			command := "docker ps -a --format='{{json .}},'"

			output, err := executarComanda(serverAddress, record.GetString("username"), record.GetString("pass"), command)
			if err != nil {
				fmt.Printf("Error al ejecutar el comando SSH: %v\n", err)
				return c.JSON(http.StatusForbidden, "[]")
			}

			var sortidaFinal = string(output)
			if len(sortidaFinal) > 0 {
				sortidaFinal = sortidaFinal[:len(sortidaFinal)-2]
			}

			fmt.Println("Sortida final:")
			fmt.Println(sortidaFinal)

			return c.JSON(http.StatusOK, string("{\"containers\": ["+sortidaFinal+"]}"))
		} /* optional middlewares */)

		return nil
	})

	// serves static files from the provided public dir (if exists)
	app.OnBeforeServe().Add(func(e *core.ServeEvent) error {
		e.Router.GET("/*", apis.StaticDirectoryHandler(os.DirFS("./pb_public"), false))
		return nil
	})

	if err := app.Start(); err != nil {
		log.Fatal(err)
	}

}
