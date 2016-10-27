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
	listenAddress = flag.String("listen", "0.0.0.0:80", "Where the server listens for connections. [interface]:port")
	staticPath    = flag.String("static", "../static/", "Location of static files.")
	modulesPath   = flag.String("modules", "./modules/", "Location of modules.")
)

func init() {
	flag.Parse()
}

func main() {
	http.Handle("/", http.FileServer(http.Dir(*staticPath)))
	http.HandleFunc("/server/", func(w http.ResponseWriter, r *http.Request) {
		module := r.URL.Query().Get("module")
		if module == "" {
			http.Error(w, "No module specified, or requested module doesn't exist.", 406)
			return
		}

		// Execute the command
		cmdPath := fmt.Sprintf("%s/%s.sh", *modulesPath, module)
		cmd := exec.Command(cmdPath)
		var output bytes.Buffer
		cmd.Stdout = &output
		err := cmd.Run()
		if err != nil {
			fmt.Printf("Error executing '%s': %s\n\tScript output: %s\n", module, err.Error(), output.String())
			http.Error(w, "Unable to execute module.", http.StatusInternalServerError)
			return
		}

		w.Write(output.Bytes())
	})

	fmt.Println("Starting http server at:", *listenAddress)
	err := http.ListenAndServe(*listenAddress, nil)
	if err != nil {
		fmt.Println("Error starting http server:", err)
		os.Exit(1)
	}
}
