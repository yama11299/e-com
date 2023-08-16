package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gorilla/mux"

	_ "github.com/mattn/go-sqlite3"
	"github.com/yama11299/e-com/product/internal/app/bl"
	"github.com/yama11299/e-com/product/internal/app/bl/dl"
	"github.com/yama11299/e-com/product/internal/app/handler"
	"github.com/yama11299/e-com/product/internal/app/spec"
	"github.com/yama11299/e-com/product/pb"
)

func main() {
	log.Println("Initializing E-Com Server...")

	logFile, err := os.OpenFile("product.log", os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	if err != nil {
		os.Exit(1)
	}

	defer logFile.Close()
	w := io.MultiWriter(os.Stdout, logFile)
	logger := log.New(w, "", 0)

	// setup database
	logger.Println("Initializing and connecting database")
	db, err := dl.InitDB()
	if err != nil {
		return
	}

	dl := dl.NewProductDL(db)
	bl := bl.NewProductBL(logger, dl)
	// setup router
	logger.Println("Setting up router")
	router := mux.NewRouter()

	router.HandleFunc(spec.ListPath, handler.List(bl)).Methods(http.MethodGet)
	router.HandleFunc(spec.UpdateQuantityPath, handler.UpdateQuantity(bl)).Methods(http.MethodPatch)

	// start server
	go func() {
		_ = http.ListenAndServe(":8080", router)
	}()

	go func() {
		_ = pb.StartRPCServer(bl)
	}()

	errs := make(chan error)
	go func() {
		c := make(chan os.Signal, 2)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errs <- fmt.Errorf("%s", <-c)
	}()

	logger.Println("HTTP server up and running")
	logger.Println("Exiting server", "error", <-errs)

}
