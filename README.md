# cli-application

A Cobra-based CLI application with commands for OAuth2 authentication and status check.

## Features

- **Login Command:** Initiates the OAuth2 authentication flow, opening the user's default browser for authentication.
- **Status Command:** Checks if the CLI application is currently logged in.
- **Secure OAuth2 Integration:** Uses Google OAuth2 for authentication, ensuring secure access to user information.
- **Customizable Configuration:** Supports configuration of OAuth2 client ID, client secret, redirect URL, and port via command-line flags.

## Prerequisites

- Go version 1.16+
- Access to Google Cloud Console to register OAuth2 credentials

## Installation

Clone the repository:

```sh
git clone git@github.com:Greckas/cli-application.git
cd cli-application
```

Build the application:

```sh
go build -o cli-application cmd/main.go
```

## Run the CLI Application

Run the application with the login command to start the authentication process:

- Use default(my own) credentials configured in the application:

```sh
./cli-application login
```

- Or, specify your OAuth2 credentials and parameters via flags:

```sh
./cli-application login --client-id YOUR_CLIENT_ID --client-secret YOUR_CLIENT_SECRET --redirect-url http://localhost:8080/callback --port 8080
```

  Replace YOUR_CLIENT_ID and YOUR_CLIENT_SECRET with your OAuth2 credentials obtained from Google Cloud Console.
  List of available flags:
  - client-id
  - client-secret
  - redirect-url
  - port
  - scopes

  Note: The application by default supports only redirect URIs configured in my Google Cloud Console: http://localhost:8080/callback, http://localhost:8081/callback, and http://localhost/callback.

## Authenticate via Browser

The CLI application will open your default browser to authenticate with Google. Follow the instructions in the browser to log in and authorize the application.

## Check Login Status

To check if the CLI application is logged in, use the status command:

```sh
./cli-application status
```

## Create Executable Files

To create executable files depending on your operating system, use the following commands:

- Windows:

  ```sh
  GOOS=windows GOARCH=amd64 go build -o dist/application-cli-windows-amd64.exe
  ```

- Linux:

  ```sh
  GOOS=linux GOARCH=amd64 go build -o dist/application-cli-linux-amd64
  ```

- macOS:

  ```sh
  GOOS=darwin GOARCH=amd64 go build -o dist/application-cli-darwin-amd64
  ```
