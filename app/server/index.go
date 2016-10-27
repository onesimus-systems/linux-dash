package main

import (
	"bytes"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/exec"
)

var (
	listenAddress string
	staticPath    string
	modulesPath   string
)

func init() {
	flag.StringVar(&listenAddress, "listen", "0.0.0.0:80", "Where the server listens for connections. [interface]:port")
	flag.StringVar(&staticPath, "static", "../static/", "Location of static files.")
	flag.StringVar(&modulesPath, "modules", "./modules/", "Location of modules.")
	flag.Parse()
}

func main() {
	http.Handle("/", http.FileServer(http.Dir(staticPath)))
	http.HandleFunc("/server/", func(w http.ResponseWriter, r *http.Request) {
		module := r.URL.Query().Get("module")
		if module == "" {
			http.Error(w, "No module specified, or requested module doesn't exist.", http.StatusNotAcceptable)
			return
		}

		// Execute the command
		cmdPath := fmt.Sprintf("%s/%s.sh", modulesPath, module)
		if !fileExists(cmdPath) {
			http.Error(w, "No module specified, or requested module doesn't exist.", http.StatusNotAcceptable)
			return
		}

		var output bytes.Buffer
		cmd := exec.Command(cmdPath)
		cmd.Stdout = &output
		if err := cmd.Run(); err != nil {
			fmt.Printf("Error executing '%s': %s\n\tScript output: %s\n", module, err.Error(), output.String())
			http.Error(w, "Unable to execute module.", http.StatusInternalServerError)
			return
		}

		w.Write(output.Bytes())
	})

	fmt.Println("Starting http server at:", listenAddress)
	if err := http.ListenAndServe(listenAddress, nil); err != nil {
		fmt.Println("Error starting http server:", err)
		os.Exit(1)
	}
}

func fileExists(file string) bool {
	_, err := os.Stat(file)
	return !os.IsNotExist(err)
}
