package main

import (
	"flag"
	"log"
	"net/http"

	"ytsruh.com/saas/cache"
	"ytsruh.com/saas/controllers"
	"ytsruh.com/saas/data"
)

func main() {
	dn := flag.String("driver", "postgres", "name of the database driver to use, postgres or mongo are supported")
	ds := flag.String("datasource", "", "database connection string")
	q := flag.Bool("queue", false, "set as queue pub/sub subscriber and task executor")
	e := flag.String("env", "dev", "set the current environment [dev|staging|prod]")
	flag.Parse()

	if len(*dn) == 0 || len(*ds) == 0 {
		flag.Usage()
		return
	}

	api := controllers.NewAPI()

	// open the database connection
	db := &data.DB{}

	if err := db.Open(*dn, *ds); err != nil {
		log.Fatal("unable to connect to the database:", err)
	}

	api.DB = db

	isDev := false
	if *e == "dev" {
		isDev = true
	}

	// Set as Redis pub/sub subscriber for the queue executor if q is true
	cache.New(*q, isDev)

	if err := http.ListenAndServe(":8080", api); err != nil {
		log.Println(err)
	}
}
