package main

import (
	"log"
	"net/http"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"runtime"
	"time"
)

const httpAddr = "127.0.0.1:3998"

func main() {
	root := getRoot()
	log.Println("Serving content from", root)
	registerStatic(root)

	go func() {
		url := "http://" + httpAddr
		if waitServer(url) && startBrowser(url) {
			log.Printf("A browser window should open. If not, please visit %s", url)
		} else {
			log.Printf("Please open your web browser and visit %s", url)
		}
	}()
	log.Fatal(http.ListenAndServe(httpAddr, nil))
}

func getRoot() string {
	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}
	exPath := filepath.Dir(ex)
	return path.Join(exPath, "html")
}

func registerStatic(root string) {
	static := FileServer(root)
	http.Handle("/", static)
}

// waitServer waits some time for the http Server to start
// serving url. The return value reports whether it starts.
func waitServer(url string) bool {
	tries := 20
	for tries > 0 {
		resp, err := http.Get(url)
		if err == nil {
			resp.Body.Close()
			return true
		}
		time.Sleep(100 * time.Millisecond)
		tries--
	}
	return false
}

// startBrowser tries to open the URL in a browser, and returns
// whether it succeed.
func startBrowser(url string) bool {
	// try to start the browser
	var args []string
	switch runtime.GOOS {
	case "darwin":
		args = []string{"open"}
	case "windows":
		args = []string{"cmd", "/c", "start"}
	default:
		args = []string{"xdg-open"}
	}
	cmd := exec.Command(args[0], append(args[1:], url)...)
	return cmd.Start() == nil
}
