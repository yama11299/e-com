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
	"github.com/yama11299/e-com/order/internal/app/bl"
	"github.com/yama11299/e-com/order/internal/app/bl/dl"
	"github.com/yama11299/e-com/order/internal/app/handler"
	"github.com/yama11299/e-com/order/internal/app/spec"
	productGRPC "github.com/yama11299/e-com/product/grpc"
	"google.golang.org/grpc"
)

func main() {
	log.Println("Initializing order service...")

	logFile, err := os.OpenFile("order.log", os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	if err != nil {
		return
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

	conn, err := grpc.Dial("0.0.0.0:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("failed to connect to product service !!!")
	}

	product := productGRPC.NewProductClient(conn)

	dl := dl.NewOrderDL(db)
	bl := bl.NewOrderBL(logger, dl, product)
	// setup router
	logger.Println("Setting up router")
	router := mux.NewRouter()

	router.HandleFunc(spec.CreateOrderPath, handler.CreateOrder(bl)).Methods(http.MethodPost)
	router.HandleFunc(spec.GetOrderPath, handler.GetOrder(bl)).Methods(http.MethodGet)

	// start server
	go func() {
		_ = http.ListenAndServe(":8081", router)
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
