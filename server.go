package main

import (
	"fmt"
	"go_grpc/services"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"google.golang.org/grpc"
	"log"
	"net/http"
	"strings"
)

func grpcHandlerFunc(grpcServer *grpc.Server, otherHandler http.Handler) http.Handler {
	return h2c.NewHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println(r)
		fmt.Println(r.Proto)
		fmt.Println(r.ProtoMajor)
		if r.ProtoMajor == 2 && strings.Contains(r.Header.Get("Content-Type"), "application/grpc") {
			grpcServer.ServeHTTP(w, r)
		} else {
			otherHandler.ServeHTTP(w, r)
		}
	}), &http2.Server{})
}

func main()  {
	rpcServer:=grpc.NewServer()
	services.RegisterProdServiceServer(rpcServer,&services.ProdService{})

	//lis,_:=net.Listen("tcp",":8081")
	//rpcServer.Serve(lis)

	mux:=http.NewServeMux()
	mux.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		//fmt.Println(request)
		//fmt.Println(request.Header)
		rpcServer.ServeHTTP(writer,request)
	})

	httsServer:=&http.Server{
		Addr:":8081",
		Handler:grpcHandlerFunc(rpcServer,mux),
	}
	err:=httsServer.ListenAndServe()
	if err!=nil{
		log.Fatal(err)
	}
}
