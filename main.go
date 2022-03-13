package main

import (
	"context"
	"fmt"
	"net/http"

	"mmm3w/sparking/resb"

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
		r.PathPrefix(k).Handler(http.StripPrefix(k, http.FileServer(http.Dir(v))))
	}

	attachComponent(r)

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

func attachComponent(r *mux.Router) {
	componentMapping := viper.GetStringMapString("component")
	fmt.Println("component:", componentMapping)

	//resb
	if componentMapping["resb"] != "" {
		r.HandleFunc(componentMapping["resb"], resb.Find)
	}
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
