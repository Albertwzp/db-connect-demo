package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	dbv1 "db-connect-demo/api/v1"
	"db-connect-demo/controllers"
	"db-connect-demo/lib"

	"github.com/gin-gonic/gin"
	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	_ "k8s.io/client-go/plugin/pkg/client/auth"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/metrics/server"
)

var (
	scheme = runtime.NewScheme()
)

func init() {
	utilruntime.Must(clientgoscheme.AddToScheme(scheme))
	utilruntime.Must(dbv1.AddToScheme(scheme))
}

func main() {
	port := flag.String("port", "8080", "http server port")
	metricsAddr := flag.String("metrics-bind-address", ":8081", "metrics bind address")
	healthProbeAddr := flag.String("health-probe-bind-address", ":8082", "health probe bind address")
	leaderElect := flag.Bool("leader-elect", false, "enable leader election for controller manager")
	flag.Parse()

	// Create manager
	mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), ctrl.Options{
		Scheme:                 scheme,
		Metrics:                server.Options{BindAddress: *metricsAddr},
		HealthProbeBindAddress: *healthProbeAddr,
		LeaderElection:         *leaderElect,
		LeaderElectionID:       "db-connect-demo.local",
	})
	if err != nil {
		fmt.Fprintln(os.Stderr, "unable to start manager:", err)
		os.Exit(1)
	}

	// Register Reconcilers
	if err := (&controllers.PostgreSQLConnectionReconciler{Client: mgr.GetClient(), Scheme: mgr.GetScheme()}).SetupWithManager(mgr); err != nil {
		fmt.Fprintln(os.Stderr, "unable to create PostgreSQL controller:", err)
		os.Exit(1)
	}
	if err := (&controllers.MySQLConnectionReconciler{Client: mgr.GetClient(), Scheme: mgr.GetScheme()}).SetupWithManager(mgr); err != nil {
		fmt.Fprintln(os.Stderr, "unable to create MySQL controller:", err)
		os.Exit(1)
	}
	if err := (&controllers.SQLiteConnectionReconciler{Client: mgr.GetClient(), Scheme: mgr.GetScheme()}).SetupWithManager(mgr); err != nil {
		fmt.Fprintln(os.Stderr, "unable to create SQLite controller:", err)
		os.Exit(1)
	}
	if err := (&controllers.KafkaConnectionReconciler{Client: mgr.GetClient(), Scheme: mgr.GetScheme()}).SetupWithManager(mgr); err != nil {
		fmt.Fprintln(os.Stderr, "unable to create Kafka controller:", err)
		os.Exit(1)
	}
	if err := (&controllers.SolaceConnectionReconciler{Client: mgr.GetClient(), Scheme: mgr.GetScheme()}).SetupWithManager(mgr); err != nil {
		fmt.Fprintln(os.Stderr, "unable to create Solace controller:", err)
		os.Exit(1)
	}

	// Start manager in background
	ctx := ctrl.SetupSignalHandler()
	go func() {
		if err := mgr.Start(ctx); err != nil {
			fmt.Fprintln(os.Stderr, "manager exited non-zero:", err)
			os.Exit(1)
		}
	}()

	// Start Gin API server
	r := gin.Default()

	// API routes: read backends from in-memory lib registry populated by Reconcilers
	r.GET("/ping", func(c *gin.Context) {
		res := lib.HealthAll(c.Request.Context())
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
		rows, err := lib.QueryBackend(c.Request.Context(), req.Backend, req.Query)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"rows": rows})
	})

	// Serve UI
	if _, err := os.Stat("frontend/dist/index.html"); err == nil {
		r.Static("/ui", "frontend/dist")
		r.NoRoute(func(c *gin.Context) {
			p := c.Request.URL.Path
			if strings.HasPrefix(p, "/ui") || p == "/" {
				c.File("frontend/dist/index.html")
				return
			}
			c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		})
	} else {
		r.StaticFile("/ui", "frontend.html")
		r.GET("/", func(c *gin.Context) { c.Redirect(http.StatusFound, "/ui") })
	}

	addr := ":" + *port
	fmt.Println("starting API server on", addr, "(UI at /ui)")

	// Run server (blocking)
	if err := r.Run(addr); err != nil {
		fmt.Fprintln(os.Stderr, "server error:", err)
		// allow manager goroutine to exit via signal
		// give it a moment to shutdown
		time.Sleep(1 * time.Second)
		os.Exit(1)
	}
}
