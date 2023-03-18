package main

import (
	"flag"
	"monolith/router"
	"monolith/users"
)

func main() {
	var create = flag.String("create", "", "Create a new user")
	var password = flag.String("password", "", "Set user password")
	var super = flag.Int("super", 0, "Set user as super")
	var unsuper = flag.Int("unsuper", 0, "Set user as not super")
	var serve = flag.Bool("serve", false, "Serve the web app")

	flag.Parse()

	if *create != "" && *password != "" {
		err := users.CreateUser(*create, *password)
		if err != nil {
			panic(err)
		}
		return
	}

	if *super != 0 {
		err := users.UpdateUserSuperStatus(*super, 1)
		if err != nil {
			panic(err)
		}
		return
	}

	if *unsuper != 0 {
		err := users.UpdateUserSuperStatus(*unsuper, 0)
		if err != nil {
			panic(err)
		}
		return
	}

	if *serve {
		err := router.Serve(8080)
		if err != nil {
			panic(err)
		}
		return
	}

	flag.PrintDefaults()
}
