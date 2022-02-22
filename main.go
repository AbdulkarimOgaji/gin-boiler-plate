package main

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"os/exec"
	"runtime"

	"github.com/urfave/cli"
)

const boilerPlate = `package main
	import(
		"github.com/gin-gonic/gin"
		"log"
	)

	func homePage(c *gin.Context) {
		c.HTML(200, "index.html", nil)
	}


	func main() {
		r := gin.Default()
		r.LoadHTMLGlob("templates/*")

		r.GET("/", homePage)

		err := r.Run(":%v")
		if err != nil {
			log.Fatal(err)
		}
	}
`
const index = `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta http-equiv="X-UA-Compatible" content="IE=edge">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Gin-App</title>
</head>
<body>
    <h1>Welcome to gin App</h1>
</body>
</html>`

func openbrowser(url string) {
	var err error

	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}
	if err != nil {
		log.Fatal(err)
	}

}

func enableReload() {
	log.Println("Live Reload enables")
}

func main() {
	app := cli.NewApp()
	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:  "live",
			Usage: "decide weather to enable live reload for the gin app",
		},
		cli.StringFlag{
			Name:  "port",
			Value: "8000",
			Usage: "specifies port for the server",
		},
		cli.BoolFlag{
			Name:  "browser",
			Usage: "decide weather you want to make a request to the server from browser on start up",
		},
	}

	app.Action = func(c *cli.Context) error {
		if c.NArg() < 1 {
			err := fmt.Errorf("you did not specify the directory for the project")
			log.Println(err)
			return fmt.Errorf("you did not specify the directory for the project")
		}
		mainDir := c.Args()[0]
		port := c.String("port")
		os.Mkdir(mainDir, fs.ModeDir)
		os.Chdir(mainDir)
		os.Mkdir("templates", fs.ModeDir)
		mainFile, err := os.Create("main.go")
		if err != nil {
			return err
		}
		mainFile.WriteString(fmt.Sprintf(boilerPlate, port))
		cmd := exec.Command("go", "mod", "init", mainDir)
		cmd.Run()
		cmd = exec.Command("go", "mod", "tidy")
		cmd.Run()
		os.Chdir("templates")
		mainFile, _ = os.Create("index.html")
		mainFile.WriteString(index)
		os.Chdir("..")
		cmd = exec.Command("go", "run", "main.go")
		cmd.Stdout = os.Stdout
		cmd.Start()
		if c.Bool("browser") {
			openbrowser(fmt.Sprintf("http://localhost:%v", port))
		}
		if c.Bool("live") {
			enableReload()
		}
		return nil
	}
	app.Run(os.Args)
}
