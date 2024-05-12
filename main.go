package main

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

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
		return err.Error(), fmt.Errorf("Error al conectar al servidor SSH: %v", err)
	}
	defer client.Close()

	// Crear una nueva sesión SSH
	session, err := client.NewSession()
	if err != nil {
		return err.Error(), fmt.Errorf("Error al crear la sesión SSH: %v", err)
	}
	defer session.Close()

	// Create buffers to store output and error
	var stdoutBuf, stderrBuf bytes.Buffer

	// Redirect command output and error to buffers
	session.Stdout = &stdoutBuf
	session.Stderr = &stderrBuf

	// Ejecutar el comando en la sesión SSH
	err = session.Run(command)
	if err != nil {
		return err.Error(), fmt.Errorf("%v", stderrBuf.String())
	}

	// Retrieve output and error messages

	return stdoutBuf.String(), nil
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
		e.Router.GET("/functions/:server/dockercompose/:composeid/up", func(c echo.Context) error {
			server := c.PathParam("server")
			composeId := c.PathParam("composeid")

			record, err := app.Dao().FindRecordById("server", server)
			if err != nil {
				return c.JSON(http.StatusForbidden, string("[]"))
			}

			compose, err := app.Dao().FindRecordById("docker_compose", composeId)
			if err != nil {
				return c.JSON(http.StatusForbidden, string("[]"))
			}

			fmt.Printf("Test %v", compose.GetString("content"))
			serverAddress := fmt.Sprintf("%s:%s", record.GetString("ip"), strconv.Itoa(record.GetInt("port")))
			command := "mkdir docktopusrun > /dev/null 2> /dev/null; cd docktopusrun > /dev/null 2> /dev/null; echo '" + compose.GetString("content") + "' > docker-compose.yml 2> /dev/null; docker compose up -d"

			output, err := executarComanda(serverAddress, record.GetString("username"), record.GetString("pass"), command)
			if err != nil {
				fmt.Printf("Error al ejecutar el comando SSH: %v\n", err)
				var sortidaFinal = string(err.Error())
				if len(sortidaFinal) > 0 {
					sortidaFinal = sortidaFinal[:len(sortidaFinal)-1]
				}
				sortidaFinal = strings.ReplaceAll(sortidaFinal, `"`, `\"`)
				sortidaFinal = strings.ReplaceAll(sortidaFinal, "\n", " ")
				return c.JSON(http.StatusOK, string("{\"stat\": \"err\", \"resultat\": \""+sortidaFinal+"\"}"))
			}

			var sortidaFinal = string(output)
			if len(sortidaFinal) > 0 {
				sortidaFinal = sortidaFinal[:len(sortidaFinal)-2]
			}

			fmt.Println("Sortida final:")
			fmt.Println(sortidaFinal)

			return c.JSON(http.StatusOK, string("{\"stat\": \"ok\", \"resultat\": \""+sortidaFinal+"\"}"))
		} /* optional middlewares */)

		return nil
	})

	app.OnBeforeServe().Add(func(e *core.ServeEvent) error {
		e.Router.GET("/functions/:server/dockercompose/:composeid/down", func(c echo.Context) error {
			server := c.PathParam("server")
			composeId := c.PathParam("composeid")

			record, err := app.Dao().FindRecordById("server", server)
			if err != nil {
				return c.JSON(http.StatusForbidden, string("[]"))
			}

			compose, err := app.Dao().FindRecordById("docker_compose", composeId)
			if err != nil {
				return c.JSON(http.StatusForbidden, string("[]"))
			}

			fmt.Printf("Test %v", compose.GetString("content"))
			serverAddress := fmt.Sprintf("%s:%s", record.GetString("ip"), strconv.Itoa(record.GetInt("port")))
			command := "mkdir docktopusrun > /dev/null 2> /dev/null; cd docktopusrun > /dev/null 2> /dev/null; echo '" + compose.GetString("content") + "' > docker-compose.yml 2> /dev/null; docker compose down"

			output, err := executarComanda(serverAddress, record.GetString("username"), record.GetString("pass"), command)
			if err != nil {
				fmt.Printf("Error al ejecutar el comando SSH: %v\n", err)
				var sortidaFinal = string(err.Error())
				if len(sortidaFinal) > 0 {
					sortidaFinal = sortidaFinal[:len(sortidaFinal)-1]
				}
				sortidaFinal = strings.ReplaceAll(sortidaFinal, `"`, `\"`)
				sortidaFinal = strings.ReplaceAll(sortidaFinal, "\n", " ")
				return c.JSON(http.StatusOK, string("{\"stat\": \"err\", \"resultat\": \""+sortidaFinal+"\"}"))
			}

			var sortidaFinal = string(output)
			if len(sortidaFinal) > 0 {
				sortidaFinal = sortidaFinal[:len(sortidaFinal)-2]
			}

			fmt.Println("Sortida final:")
			fmt.Println(sortidaFinal)

			return c.JSON(http.StatusOK, string("{\"stat\": \"ok\", \"resultat\": \""+sortidaFinal+"\"}"))
		} /* optional middlewares */)

		return nil
	})

	// docker exec -t f42bf0fe4043 ls /bin
	app.OnBeforeServe().Add(func(e *core.ServeEvent) error {
		e.Router.GET("/functions/:server/container/:container/list/:directori", func(c echo.Context) error {
			server := c.PathParam("server")
			container := c.PathParam("container")
			directori := c.PathParam("directori")

			record, err := app.Dao().FindRecordById("server", server)
			if err != nil {
				return c.JSON(http.StatusForbidden, string("[]"))
			}

			serverAddress := fmt.Sprintf("%s:%s", record.GetString("ip"), strconv.Itoa(record.GetInt("port")))
			command := "docker exec -t " + container + " ls -la -p -d " + directori + "/*"

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

			return c.JSON(http.StatusOK, string(sortidaFinal))
		} /* optional middlewares */)

		return nil
	})

	app.OnBeforeServe().Add(func(e *core.ServeEvent) error {
		e.Router.GET("/functions/:server/container/:container/list", func(c echo.Context) error {
			server := c.PathParam("server")
			container := c.PathParam("container")

			record, err := app.Dao().FindRecordById("server", server)
			if err != nil {
				return c.JSON(http.StatusForbidden, string("[]"))
			}

			serverAddress := fmt.Sprintf("%s:%s", record.GetString("ip"), strconv.Itoa(record.GetInt("port")))
			command := "docker exec -t " + container + " find . -print -ls"

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

			return c.JSON(http.StatusOK, string(sortidaFinal))
		} /* optional middlewares */)

		return nil
	})

	app.OnBeforeServe().Add(func(e *core.ServeEvent) error {
		e.Router.GET("/functions/:server/image/:image/pull", func(c echo.Context) error {
			server := c.PathParam("server")
			image := c.PathParam("image")

			record, err := app.Dao().FindRecordById("server", server)
			if err != nil {
				return c.JSON(http.StatusForbidden, string("[]"))
			}

			serverAddress := fmt.Sprintf("%s:%s", record.GetString("ip"), strconv.Itoa(record.GetInt("port")))
			command := "docker pull " + image

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
		e.Router.GET("/functions/volumes/:server", func(c echo.Context) error {
			server := c.PathParam("server")

			record, err := app.Dao().FindRecordById("server", server)
			if err != nil {
				return c.JSON(http.StatusForbidden, string("[]"))
			}

			serverAddress := fmt.Sprintf("%s:%s", record.GetString("ip"), strconv.Itoa(record.GetInt("port")))

			command := "docker volume ls --format='{{json .}},'"

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

			return c.JSON(http.StatusOK, string("{\"volumes\": ["+sortidaFinal+"]}"))
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
