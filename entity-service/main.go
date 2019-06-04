package main

import (
	"log"
	"os"

	_ "github.com/lib/pq"
	"github.com/lisin-andrey/msvc-demo/common/pkg/tools"
	"github.com/lisin-andrey/msvc-demo/entity-service/cmd"
)

const (
	appNameLogPrefix = "EntitySvc: "
)

func main() {
	log.SetOutput(os.Stdout)
	log.SetFlags(log.LstdFlags | log.Lshortfile | log.Lmicroseconds)
	log.SetPrefix(appNameLogPrefix)

	app, err := cmd.Prepare()
	if err != nil {
		tools.Fatalfln("Error occured during preparation stage. [%s]", err.Error())
	}
	err = app.Run()
	if err != nil {
		tools.Fatalfln("Error occured during running stage. [%s]", err.Error())
	}
}
