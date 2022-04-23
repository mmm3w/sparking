package main

import (
	"context"
	"fmt"
	"net/http"

	"mmm3w/sparking/resb"
	"mmm3w/sparking/wsk"

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
	//resb
	resbComponent := viper.Get("component.resb")
	if resbComponent != nil {
		fmt.Println("component:resb--> ", resbComponent)
		res := viper.GetString("component.resb.res")
		if res != "" {
			r.HandleFunc(res, resb.Find)
		}

		upload := viper.GetString("component.resb.upload")
		if upload != "" {
			r.HandleFunc(upload, resb.Upload)
		}
	}

	//wsk
	wskComponent := viper.Get("component.wsk")
	if wskComponent != nil {
		wsk.Init()
		fmt.Println("component:wsk--> ", wskComponent)
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
