package main

import (
	cfg "github.com/jyyds/filestore/config"
	"github.com/jyyds/filestore/route"
)

func main() {

	router := route.Router()
	router.Run(cfg.UploadServiceHost)
}
