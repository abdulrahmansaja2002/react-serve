package main

import (
	"echo-react-serve/config"
	"echo-react-serve/server"
)

func main() {
	e := server.New()
	e.Logger.Fatal(e.Start(":" + config.Envs.App.Port))
	defer server.Close()
	// temp.Show()
}
