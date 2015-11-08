package main
import ( "net/http"
"os"
//"fmt"
)
func main() {
	
    dir := os.Args[1]
    port := os.Args[2]
	//port := "8080"
	//fmt.Println(dir)
    http.Handle("/",	http.FileServer(http.Dir(dir)))
    http.ListenAndServe(":"+port,nil)
}
