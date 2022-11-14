package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/google"

	worker "github.com/the-runtime/gdrive/workers"
)

func main() {

	key := "tabish_secret"
	maxAge := 86400 * 30
	isProd := false

	store := sessions.NewCookieStore([]byte(key))
	store.MaxAge(maxAge)
	store.Options.Path = "/"
	store.Options.Secure = isProd

	gothic.Store = store

	goth.UseProviders(
		google.New("882134345746-46gep2dg1epok994utqh8o6dmurspkb3.apps.googleusercontent.com", "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", "http://localhost/google/callback"),
	)
	r := mux.NewRouter()

	r.HandleFunc("/auth/google/callback", func(res http.ResponseWriter, req *http.Request) {
		_, err := gothic.CompleteUserAuth(res, req)
		if err != nil {
			fmt.Fprintln(res, err)
			return
		}

	})

	r.HandleFunc("/auth/google", func(res http.ResponseWriter, req *http.Request) {
		gothic.BeginAuthHandler(res, req)
	})

	dispatch := worker.NewDispatcher(3)
	dispatch.Run()
	worker.InitJobQueue()

	r.HandleFunc("/test/{url}", func(res http.ResponseWriter, req *http.Request) {
		vars := mux.Vars(req)
		url := vars["url"]
		job := worker.NewJob(url)

		worker.JobQueue <- job
	})

	r.HandleFunc("/ping", func(res http.ResponseWriter, req *http.Request) {
		res.WriteHeader(http.StatusOK)
		res.Write([]byte("Ping is working"))
	})

}
