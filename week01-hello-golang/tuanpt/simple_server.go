package main

import (
	"fmt"
	"net/http"
	"os"
	)
	
// execute : simple_server port dirpath

func main() {
	port := os.Args[1]
	dir  := os.Args[2]
	fmt.Println(port," ",dir)
	http.Handle("/",http.FileServer(http.Dir(dir)))
	http.ListenAndServe(":" + port,nil)
}
