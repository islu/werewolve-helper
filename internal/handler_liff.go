package internal

import (
	"log"
	"net/http"
	"text/template"
)

func RegisterLIFF(config BotConfig) {

	http.HandleFunc("/liff/setting-role", func(w http.ResponseWriter, r *http.Request) {

		t, err := template.ParseFiles("internal/liff/setting_role.html")
		if err != nil {
			log.Fatalln(err)
		}

		err = t.Execute(w, config)
		if err != nil {
			log.Println(err)
		}
	})

}
