package main

import(
	"io"
	"log"
	"net/http"
	"os"
	"io/ioutil"
	"strings"
	"time"
	"flag"
	"strconv"
)

func main() {
	host := flag.String("h","127.0.0.1","host name or ip address")
	port := flag.Int("p", 8080, "port")
	name := flag.String("n", "/request", "request root name")
	flag.CommandLine.Parse(os.Args[1:])

	http.HandleFunc(*name, request)

	wd, err := os.Getwd()
	if err != nil{
		log.Fatal(err)
	}
	http.Handle("/", http.FileServer(http.Dir(wd))) // view static directory

	err = http.ListenAndServe(*host + ":" + strconv.Itoa(*port), nil)
	if err != nil {
		log.Fatal(err)
	}
}

func request(w http.ResponseWriter, r *http.Request){
	if err := r.ParseForm(); err != nil {
		log.Println(err.Error())
	}
	log.Println(r.FormValue("url"), r.Method)

	body := strings.NewReader(r.PostForm.Encode())
	req, err := http.NewRequest(r.Method, r.Form.Get("url"), body)
	if err != nil {
		io.WriteString(w, err.Error())
		return
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded;param=value")
	client := &http.Client{Timeout: 10*time.Second} // important!!!
	resp, err := client.Do(req)
	if err != nil {
		io.WriteString(w, err.Error())
	}else {
		result, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			io.WriteString(w, "Get data failed: " + err.Error())
		}else {
			io.WriteString(w, string(result))
		}
	}
}
