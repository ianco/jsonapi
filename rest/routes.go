package rest

import "net/http"

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type Routes []Route

var routes = Routes{
	Route{
		"Index",
		"GET",
		"/",
		Index,
	},
	Route{
		"ConfigShow",
		"GET",
		"/config",
		ConfigShow,
	},
	Route{
		"ConfigUpdate",
		"POST",
		"/config",
		ConfigUpdate,
	},
	Route{
		"ConfigUpdate",
		"PUT",
		"/config",
		ConfigUpdate,
	},
}
