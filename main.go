package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"db-connect-demo/lib"

	"github.com/gin-gonic/gin"
)

type BackendSpec struct {
	Driver string `json:"driver"`
	DSN    string `json:"dsn"`
}

func main() {
	port := flag.String("port", "8080", "http server port")
	backendsFile := flag.String("backends-file", "", "path to json file defining backends map[name]={driver,dsn}")
	backendsJSON := flag.String("backends", "", "json map of backends as inline string")
	flag.Parse()

	var specs map[string]BackendSpec
	if *backendsFile != "" {
		b, err := ioutil.ReadFile(*backendsFile)
		if err != nil {
			fmt.Fprintln(os.Stderr, "failed to read backends file:", err)
			os.Exit(2)
		}
		if err := json.Unmarshal(b, &specs); err != nil {
			fmt.Fprintln(os.Stderr, "failed to parse backends file:", err)
			os.Exit(2)
		}
	} else if *backendsJSON != "" {
		if err := json.Unmarshal([]byte(*backendsJSON), &specs); err != nil {
			fmt.Fprintln(os.Stderr, "failed to parse backends json:", err)
			os.Exit(2)
		}
	} else {
		fmt.Fprintln(os.Stderr, "no backends specified; use -backends-file or -backends")
		os.Exit(2)
	}

	// register backends
	for name, s := range specs {
		if err := lib.RegisterBackend(name, s.Driver, s.DSN); err != nil {
			fmt.Fprintf(os.Stderr, "failed to register backend %s: %v\n", name, err)
			os.Exit(1)
		}
	}
	defer lib.CloseAllBackends()

	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) {
		res := lib.HealthAll(context.Background())
		c.JSON(http.StatusOK, res)
	})

	r.POST("/query", func(c *gin.Context) {
		var req struct {
			Backend string `json:"backend"`
			Query   string `json:"query"`
		}
		if err := c.BindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		rows, err := lib.QueryBackend(context.Background(), req.Backend, req.Query)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"rows": rows})
	})

	addr := ":" + *port
	fmt.Println("starting server on", addr)
	if err := r.Run(addr); err != nil {
		fmt.Fprintln(os.Stderr, "server error:", err)
		os.Exit(1)
	}
}
