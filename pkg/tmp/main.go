package main

import (
	"context"
	"fmt"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
	"io/ioutil"
	"log"
	"net/http"
)

var (
	githubOauthConfig *oauth2.Config
	// TODO: randomize it
	//oauthStateString = "pseudo-random"
)

func init() {
	githubOauthConfig = &oauth2.Config{
		RedirectURL:  "http://localhost:8080/callback",
		ClientID:     "e727ab35809f450261ae",
		ClientSecret: "c3c98ff11b7dee45528dcc050f8bf016bc44cd39",
		//Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email"},
		Endpoint: github.Endpoint,
	}
}

func main() {
	http.HandleFunc("/", handleMain)
	http.HandleFunc("/login", handleGithubLogin)
	http.HandleFunc("/callback", handleGithubCallback)
	fmt.Println(http.ListenAndServe(":8080", nil))
}

func handleMain(w http.ResponseWriter, r *http.Request) {
	var htmlIndex = `<html>
<body>
	<a href="/login">Github Log In</a>
</body>
</html>`

	fmt.Fprintf(w, htmlIndex)
}

func handleGithubLogin(w http.ResponseWriter, r *http.Request) {
	url := githubOauthConfig.AuthCodeURL("")
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func handleGithubCallback(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	accessToken, err := githubOauthConfig.Exchange(context.TODO(), code)

	fmt.Println(accessToken)

	if err != nil {
		panic(err)
	}

	// Get request to a set URL
	req, reqerr := http.NewRequest(
		"GET",
		"https://api.github.com/user",
		nil,
	)
	if reqerr != nil {
		log.Panic("API Request creation failed")
	}

	// Set the Authorization header before sending the request
	// Authorization: token XXXXXXXXXXXXXXXXXXXXXXXXXXX
	authorizationHeaderValue := fmt.Sprintf("token %s", accessToken.AccessToken)
	req.Header.Set("Authorization", authorizationHeaderValue)

	// Make the request
	resp, resperr := http.DefaultClient.Do(req)
	if resperr != nil {
		log.Panic("Request failed")
	}

	// Read the response as a byte slice
	respbody, _ := ioutil.ReadAll(resp.Body)

	fmt.Println(string(respbody))
}
