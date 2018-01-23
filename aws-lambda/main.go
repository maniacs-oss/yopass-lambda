package main

import (
	"net/http"
	"os"

	"github.com/akrylysov/algnhsa"
	"github.com/julienschmidt/httprouter"
	yopass "github.com/yopass/yopass-lambda"
)

func main() {
	router := httprouter.New()
	db := yopass.NewDynamo(os.Getenv("TABLE_NAME"))
	router.GET("/secret/:key", func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		yopass.GetSecret(w, r, p, db)
	})
	router.POST("/secret", func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		yopass.CreateSecret(w, r, p, db)
	})
	algnhsa.ListenAndServe(router, nil)
}
