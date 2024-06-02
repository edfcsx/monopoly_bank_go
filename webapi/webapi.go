package webapi

import (
	"encoding/json"
	"fmt"
	"monopoly_bank_go/accounts"
	"monopoly_bank_go/connection"
	"monopoly_bank_go/http"
)

var routes = map[string]func(*http.Request, *connection.Connection){
	"POST/login": Login,
}

var defaultHeaders = map[string]string{
	"Content-Type": "application/json",
}

// Handler is the function that handles the webapi requests
func Handler(c *connection.Connection) {
	go http.HandlerRequest(c.Socket, func(r *http.Request, err error) {
		if err != nil {
			c.SendAndClose(http.MakeResponse(http.BadRequest, nil, ""))
			return
		}

		// CORS
		if r.Method == "OPTIONS" {
			c.SendAndClose(http.MakeResponse(http.OK, http.CorsHeaders, ""))
			return
		}

		path := fmt.Sprintf("%s%s", r.Method, r.Path)

		if handler, ok := routes[path]; ok {
			handler(r, c)
			return
		}

		res := map[string]string{"error": "Route not found"}

		j, err := json.Marshal(res)

		if err != nil {
			c.SendAndClose(http.MakeResponse(http.InternalServerError, defaultHeaders, ""))
			return
		}

		c.SendAndClose(http.MakeResponse(http.InternalServerError, defaultHeaders, string(j)))
	})
}

func Login(r *http.Request, c *connection.Connection) {
	reqBody := map[string]string{}
	err := json.Unmarshal([]byte(r.Body), &reqBody)

	if err != nil {
		msg, _ := json.Marshal(map[string]string{"error": "json body is invalid"})
		c.SendAndClose(http.MakeResponse(http.BadRequest, defaultHeaders, string(msg)))
		return
	}

	if !ValidateObject(reqBody, []string{"username", "password"}) {
		msg, _ := json.Marshal(map[string]string{"error": "invalid request body"})
		c.SendAndClose(http.MakeResponse(http.BadRequest, defaultHeaders, string(msg)))
		return
	}

	acc := accounts.ExistsByName(reqBody["username"])

	if acc != nil {
		if acc.Password == reqBody["password"] {
			msg, _ := json.Marshal(map[string]string{"player_hash": acc.AuthenticationHash})
			c.SendAndClose(http.MakeResponse(http.OK, defaultHeaders, string(msg)))
		} else {
			msg, _ := json.Marshal(map[string]string{"error": "invalid credentials"})
			c.SendAndClose(http.MakeResponse(http.Unauthorized, defaultHeaders, string(msg)))
		}

		return
	}

	authenticationHash := accounts.Create(reqBody["username"], reqBody["password"])
	msg, _ := json.Marshal(map[string]string{"player_hash": authenticationHash})
	c.SendAndClose(http.MakeResponse(http.OK, defaultHeaders, string(msg)))
}
