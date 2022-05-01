package plates

const ServerFile = `package main
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

const Indexhtml = `<!DOCTYPE html>
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
