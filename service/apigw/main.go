package main

import "github.com/jyyds/filestore/route"

func main() {
	r := route.Router()
	r.Run(":8080")
}
