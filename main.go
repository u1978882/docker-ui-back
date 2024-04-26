package main

import (
	"encoding/json"
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

type ImageInfo struct {
	Repository string `json:"repository"`
	Tag        string `json:"tag"`
	ImageID    string `json:"image_id"`
	Created    string `json:"created"`
	Size       string `json:"size"`
}

func main() {
	app := pocketbase.New()

	app.OnBeforeServe().Add(func(e *core.ServeEvent) error {
		e.Router.GET("/functions/images/:server", func(c echo.Context) error {
			server := c.PathParam("server")

			record, err := app.Dao().FindRecordById("server", server)
			if err != nil {
				return c.JSON(http.StatusForbidden, string("[]"))
			}

			fmt.Println("password")
			fmt.Println(record.GetString("pass"))
			fmt.Println("kaladin")
			// Configuración de la conexión SSH
			sshConfig := &ssh.ClientConfig{
				User: record.GetString("username"),
				//User: "arnau",
				Auth: []ssh.AuthMethod{
					ssh.Password(record.GetString("pass")),
					//ssh.Password("kaladin"),
					// También puedes utilizar otros métodos de autenticación, como claves SSH, dependiendo de tu configuración
				},
				HostKeyCallback: ssh.InsecureIgnoreHostKey(),
				// Otras configuraciones, como HostKeyCallback, pueden ser necesarias dependiendo de tu entorno.
			}

			// Dirección del servidor SSH (host:port)
			serverAddress := fmt.Sprintf("%s:%s", record.GetString("ip"), strconv.Itoa(record.GetInt("port")))

			// Comando a ejecutar en el servidor remoto
			command := "docker images"
			//docker images --format '{"Repository": "{{.Repository}}", "Tag": "{{.Tag}}", "ID": "{{.ID}}", "Created": "{{.CREATED}}"}'

			// Realizar la conexión SSH
			client, err := ssh.Dial("tcp", serverAddress, sshConfig)
			if err != nil {
				fmt.Printf("Error al conectar al servidor SSH: %v\n", err)
				return c.JSON(http.StatusForbidden, string("[]"))
			}
			defer client.Close()

			// Crear una nueva sesión SSH
			session, err := client.NewSession()
			if err != nil {
				fmt.Printf("Error al crear la sesión SSH: %v\n", err)
				return c.JSON(http.StatusForbidden, string("[]"))
			}
			defer session.Close()

			// Ejecutar el comando en la sesión SSH
			output, err := session.CombinedOutput(command)
			if err != nil {
				fmt.Printf("Error al ejecutar el comando en el servidor remoto: %v\n", err)
				return c.JSON(http.StatusForbidden, string("[]"))
			}

			// Imprimir la salida del comando
			fmt.Println("Output del comando:")
			fmt.Println(string(output))

			lineArray := strings.Split(strings.TrimSpace(string(output)), "\n")

			if len(lineArray) > 0 {
				lineArray = lineArray[1:]
			}

			// Arreglo para almacenar los objetos ImageInfo
			var images []ImageInfo

			// Procesar cada línea y convertirla a un objeto ImageInfo
			for _, line := range lineArray {
				// Separar la línea en campos
				fields := strings.Fields(line)

				// Crear un objeto ImageInfo
				imageInfo := ImageInfo{
					Repository: fields[0],
					Tag:        fields[1],
					ImageID:    fields[2],
					Created:    fields[3],
					Size:       fields[4],
				}

				// Agregar el objeto al arreglo
				images = append(images, imageInfo)
			}

			// Convertir el arreglo a JSON
			jsonData, err := json.Marshal(images)
			if err != nil {
				fmt.Printf("Error al convertir a JSON: %v", err)
				return c.JSON(http.StatusForbidden, string("[]"))
			}

			return c.JSON(http.StatusOK, string(jsonData))
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
