package router

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var SetFruitAccesses = promauto.NewCounter(prometheus.CounterOpts{
	Name: "set_fruit_accesses_total",
	Help: "The total number of accesses to SetFruit with any method",
})

var GetTokenAccesses = promauto.NewCounter(prometheus.CounterOpts{
	Name: "get_token_accesses_total",
	Help: "The total number of accesses to GetToken with any method",
})

var IndexAccesses = promauto.NewCounter(prometheus.CounterOpts{
	Name: "index_accesses_total",
	Help: "The total number of accesses to index page with any method",
})

var LoginUserAccesses = promauto.NewCounter(prometheus.CounterOpts{
	Name: "login_user_accesses_total",
	Help: "The total number of times user tried to log in regardless of success",
})

var LogoutUserAccesses = promauto.NewCounter(prometheus.CounterOpts{
	Name: "logout_user_accesses_total",
	Help: "The total number of times user logged out",
})
