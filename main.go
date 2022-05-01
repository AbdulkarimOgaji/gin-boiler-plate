package main

import (
	"create-gin-app/plates"
	"fmt"
	"io/fs"
	"log"
	"os"
	"os/exec"
	"runtime"

	"github.com/urfave/cli"
)

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
		log.Println(c.Args())
		if c.NArg() < 1 {
			err := fmt.Errorf("you did not specify the directory for the project")
			log.Println(err)
			return fmt.Errorf("you did not specify the directory for the project")
		}
		mainDir := c.Args()[0]
		port := c.String("port")

		os.Mkdir(mainDir, fs.ModeDir)
		os.Chdir(mainDir)

		//main.go file
		mainFile, err := os.Create("main.go")
		if err != nil {
			return err
		}
		mainFile.WriteString(fmt.Sprintf(plates.ServerFile, port))

		// live reload
		if c.Bool("live") {
			cmd := exec.Command("go", "install", "github.com/codegangsta/gin")
			cmd.Run()
			port = "3001"
		}

		// go mod init and tidy
		cmd := exec.Command("go", "mod", "init", mainDir)
		cmd.Run()
		cmd = exec.Command("go", "mod", "tidy")
		cmd.Run()

		// make html file
		os.Mkdir("templates", fs.ModeDir)
		os.Chdir("templates")
		indexHtml, err := os.Create("index.html")
		indexHtml.WriteString(plates.Indexhtml)

		// go back to main dir
		os.Chdir("..")

		// start process in current stdout
		if c.Bool("live") {
			cmd = exec.Command("gin", "-p", "3001", "run", "main.go")
		} else {
			cmd = exec.Command("go", "run", "main.go")
		}

		cmd.Stdout = os.Stdout
		cmd.Start()

		// open browser if instructed
		if c.Bool("browser") {
			openbrowser(fmt.Sprintf("http://localhost:%v", port))
		}

		return nil
	}
	app.Run(os.Args)
}
