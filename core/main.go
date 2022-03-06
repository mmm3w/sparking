package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/fsnotify/fsnotify"
	"github.com/gorilla/mux"
	"github.com/spf13/viper"
)

func start() *http.Server {
	fmt.Println("start")

	r := mux.NewRouter()

	routerMapping := viper.GetStringMapString("router")
	fmt.Println("router:", routerMapping)
	for k, v := range routerMapping {
		// r.Handle(k, http.StripPrefix(k, http.FileServer(http.Dir(v))))
		r.PathPrefix(k).Handler(http.StripPrefix(k, http.FileServer(http.Dir(v))))
	}

	srv := &http.Server{
		Addr:    viper.GetString("port"),
		Handler: r,
	}

	go func() {
		err := srv.ListenAndServe()
		if err != nil {
			fmt.Println(err.Error())
		}
	}()

	return srv
}

func main() {
	viper.SetConfigFile("./config.yaml")
	viper.SetDefault("port", ":8088")
	viper.SetDefault("router", map[string]string{})

	err := viper.ReadInConfig()
	if err != nil {
		viper.WriteConfig()
		return
	}

	restart := make(chan int)
	var srv *http.Server

	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println("config changed", e.Name)
		restart <- 1
	})

	srv = start()
	for {
		<-restart
		srv.Shutdown(context.Background())
		srv = start()
	}
}
