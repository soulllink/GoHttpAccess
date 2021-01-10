package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const indexhtml = `<!DOCTYPE html>
<html lang="en">

<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Document</title>
</head>

<body>
    <h3>Files:</h3>
    <p>@@@</p>
</body>

</html>`

var fileslist []string

func start() {
	err := filepath.Walk(".",
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			fmt.Println(path, info.Size())
			fileslist = append(fileslist, path)
			return nil
		})
	if err != nil {
		log.Println(err)
	}
}

func main() {
	start()
	conn, err := net.Dial("tcp", "192.168.1.1:80")
	if err != nil {
		log.Print(err)
	}
	fmt.Println(conn.LocalAddr())
	http.HandleFunc("/files/", downloadHandler)
	http.HandleFunc("/", templateHandler)
	http.ListenAndServe(":3000", nil)
}

func templateHandler(w http.ResponseWriter, r *http.Request) {
	s := strings.Split(indexhtml, "@@@")
	out := strings.Join(s, makelistoflinks())
	io.WriteString(w, out)
}

func downloadHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	StoredAs := strings.Replace(r.URL.Path, "/files", "", 1) // file name
	fmt.Println(StoredAs)
	data, err := ioutil.ReadFile("." + StoredAs)
	if err != nil {
		fmt.Fprint(w, err)
	}
	http.ServeContent(w, r, StoredAs, time.Now(), bytes.NewReader(data))
}

func makelistoflinks() string {
	var s string
	var l []string

	for i := 0; i < len(fileslist); i++ {
		decorator := "<a href=\"./files/" + fileslist[i] + "\">" + fileslist[i] + "</a>"
		l = append(l, decorator)
	}
	s = strings.Join(l, " ")
	return s
}
