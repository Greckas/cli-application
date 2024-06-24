package cmd

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"log"
	"net/http"
	"os/exec"
	"runtime"
	"sync"
)

var (
	clientID     string
	clientSecret string
	redirectURL  string
	port         string
	scopes       []string
	oauth2Config *oauth2.Config
	token        *oauth2.Token
	stateStore   = sync.Map{} // will store state for validation
)

func init() {

	// In prod version this moves to config file
	clientID = "21362453618-eldhvcb6s4ska7baqvuv2npmeg7o0bpn.apps.googleusercontent.com"
	clientSecret = "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"
	redirectURL = "http://localhost:8080/callback"
	port = "8080"
	scopes = []string{"https://www.googleapis.com/auth/userinfo.email"}

	loginCmd.Flags().StringVar(&clientID, "client-id", clientID, "Google OAuth2 Client ID")
	loginCmd.Flags().StringVar(&clientSecret, "client-secret", clientSecret, "Google OAuth2 Client Secret")
	loginCmd.Flags().StringVar(&redirectURL, "redirect-url", redirectURL, "Redirect URL for OAuth2 callback")
	loginCmd.Flags().StringVar(&port, "port", port, "Port for the local server")
	loginCmd.Flags().StringSliceVar(&scopes, "scopes", scopes, "OAuth2 scopes")

	oauth2Config = &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Endpoint:     google.Endpoint,
		RedirectURL:  redirectURL,
		Scopes:       scopes,
	}

	rootCmd.AddCommand(loginCmd)
}

var done = make(chan struct{})
var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Perform authentication via browser",
	Run: func(cmd *cobra.Command, args []string) {

		clientID, _ := cmd.Flags().GetString("client-id")
		clientSecret, _ := cmd.Flags().GetString("client-secret")
		redirectURL, _ := cmd.Flags().GetString("redirect-url")
		port, _ := cmd.Flags().GetString("port")
		scopes, _ := cmd.Flags().GetStringSlice("scopes")

		oauth2Config := &oauth2.Config{
			ClientID:     clientID,
			ClientSecret: clientSecret,
			Endpoint:     google.Endpoint,
			RedirectURL:  redirectURL,
			Scopes:       scopes,
		}

		state := generateState() // secure random state

		http.HandleFunc("/callback", func(w http.ResponseWriter, r *http.Request) {
			queryState := r.URL.Query().Get("state")
			if queryState != state {
				http.Error(w, "State did not match", http.StatusBadRequest)
				return
			}

			oauth2Token, err := oauth2Config.Exchange(context.Background(), r.URL.Query().Get("code"))
			if err != nil {
				http.Error(w, "Failed to exchange token: "+err.Error(), http.StatusInternalServerError)
				return
			}

			// token save
			token = oauth2Token
			if err := saveToken(token); err != nil {
				http.Error(w, "Failed to save token: "+err.Error(), http.StatusInternalServerError)
				return
			}

			fmt.Fprintf(w, "Login successful")
			color.Green("Login successful")
			fmt.Println()

			// signal that login is complete
			close(done)
		})

		go func() {
			color.Yellow("Starting local server for OAuth callback...")
			if err := http.ListenAndServe(":"+port, nil); err != nil {
				log.Fatalf("Failed to start server: %v", err)
			}
		}()

		url := oauth2Config.AuthCodeURL(state, oauth2.AccessTypeOffline)
		color.Blue("Opening browser for login...")

		openBrowser(url)

		<-done
	},
}

// openBrowser opens the specified URL in the default browser of your OS
func openBrowser(url string) {
	var err error

	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}

	if err != nil {
		log.Fatalf("Failed to open browser: %v", err)
	}
}

func generateState() string {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		log.Fatalf("Failed to generate state: %v", err)
	}
	state := base64.URLEncoding.EncodeToString(b)
	return state
}
