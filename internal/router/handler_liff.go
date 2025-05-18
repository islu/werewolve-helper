package router

import (
	"log"
	"net/http"
	"text/template"
	"werewolve-helper/internal"
)

func RegisterLIFF(config internal.BotConfig) {
	http.HandleFunc("/liff/setting-role", func(w http.ResponseWriter, r *http.Request) {
		t, err := template.ParseFiles("internal/router/liff/setting_role.html")
		if err != nil {
			log.Fatalln(err)
		}

		err = t.Execute(w, config)
		if err != nil {
			log.Println(err)
		}
	})
}
