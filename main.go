package main

import (
	"taskflow/bootstrap"
)

func main() {
	app := bootstrap.Boot()

	app.Start()
}
