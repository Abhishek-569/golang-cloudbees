// main.go
package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/google/uuid"
	"github.com/open-feature/go-sdk/pkg/openfeature"

	"github.com/rollout/cloudbees-openfeature-provider-go/pkg/cloudbees"
)

var (
	APP_KEY = os.Getenv("APP_KEY")
)

func main() {

	// Register your feature flag provider

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// 	// Check for user_id value from cookies
		var value1 bool
		appKey := APP_KEY
		if provider, err := cloudbees.NewProvider(appKey); err == nil {
			openfeature.SetProvider(provider)
			client := openfeature.NewClient("app")
			value, err := client.BooleanValue(context.Background(), "enableTutorial", false, openfeature.EvaluationContext{})
			value1 = value
			fmt.Printf("flag value: %v, error: %v", value, err)
		} else {
			fmt.Printf("error creating client %v", err)
		}
		cookie, err := r.Cookie("user_id")
		var userID string

		if err != nil || cookie.Value == "" {
			// Create a new user_id using a UUID function
			userID = uuid.New().String()
			// Set the user_id cookie
			http.SetCookie(w, &http.Cookie{
				Name:  "user_id",
				Value: userID,
			})
		} else {
			userID = cookie.Value
			// Use the existing user_id
		}

		//Usecase:1
		if value1 {
			fmt.Fprintf(w, "Your flag is enabled!")
		} else {
			fmt.Fprintf(w, "Hello World!")
		}

		//Usecase:2
		// if value1 {
		// 	fmt.Fprintf(w, "There is an error in your code!!!")
		// } else {
		// 	fmt.Fprintf(w, "Your Code is working fine!")
		// }

		//Usecase:3
		// if value1 {
		// 	fmt.Fprintf(w, "This part of code is complete. \nThis part of Code is work in progress.")
		// } else {
		// 	fmt.Fprintf(w, "This part of code is complete")
		// }

		//Usecase:4
		// if value1 {
		// 	fmt.Fprintf(w, "This is an Early feature to users in Beta!")
		// } else {
		// 	fmt.Fprintf(w, "This is the Current feature.")
		// }
	})
	http.ListenAndServe(":8080", nil)
}
