package main

import "github.com/FlareZone/melon-backend/cmd"

func main() {
	//r := gin.Default()
	//
	//conf := &oauth2.Config{
	//	ClientID:     "1030441591409-i0eiesff2uj64mhe3bl66338ofcv8sar.apps.googleusercontent.com",
	//	ClientSecret: "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",
	//	RedirectURL:  "https://3fdfc921.r2.cpolar.top/callback",
	//	//RedirectURL: "http://localhost:8080/callback",
	//
	//	Scopes:   []string{"email", "profile"},
	//	Endpoint: google.Endpoint,
	//}
	//r.GET("/login", func(c *gin.Context) {
	//	// Redirect user to consent page to ask for permission
	//	// for the scopes specified above.
	//	url := conf.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	//	c.Redirect(http.StatusTemporaryRedirect, url)
	//})
	//r.GET("/callback", func(c *gin.Context) {
	//
	//	// Handle the exchange code to initiate a transport.
	//	code := c.Query("code")
	//	fmt.Println("code ==========> ")
	//	fmt.Println(code)
	//	tok, err := conf.Exchange(c, code)
	//	if err != nil {
	//		log.Fatal(err)
	//	} else {
	//		fmt.Println("token ==========> ")
	//		fmt.Println(tok)
	//	}
	//	client := conf.Client(c, tok)
	//	response, err := client.Get("https://www.googleapis.com/oauth2/v3/userinfo")
	//	if err != nil {
	//		log.Print(err)
	//		return
	//	}
	//	defer response.Body.Close()
	//	bodyBytes, err := io.ReadAll(response.Body)
	//	if err != nil {
	//		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to get userinfo", "error": err.Error()})
	//		return
	//	}
	//	log.Print("bodyBytes====>", bodyBytes)
	//})
	//r.Run(":8080") // listen and serve on 0.0.0.0:8080
	cmd.Execute()
}
