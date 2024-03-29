# Go Lade

[![Build Status](https://img.shields.io/github/actions/workflow/status/lade-io/go-lade/release.yml)](https://github.com/lade-io/go-lade/actions/workflows/release.yml)
[![Go.Dev Reference](https://img.shields.io/badge/go.dev-reference-blue.svg)](https://pkg.go.dev/github.com/lade-io/go-lade)
[![Release](https://img.shields.io/github/v/release/lade-io/go-lade.svg)](https://github.com/lade-io/go-lade/releases/latest)

Go Lade is a Go client library for the Lade V1 API.

## Installation

```sh
go get github.com/lade-io/go-lade
```

## Usage

Currently, the only authentication method is password credentials.
You can enter your username and password to create a new token:

```go
package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/lade-io/go-lade"
	"golang.org/x/crypto/ssh/terminal"
	"golang.org/x/oauth2"
)

func main() {
	conf := &oauth2.Config{
		ClientID: lade.DefaultClientID,
		Scopes:   lade.DefaultScopes,
		Endpoint: lade.Endpoint,
	}

	ctx := context.Background()
	username, password := getCredentials()

	token, err := conf.PasswordCredentialsToken(ctx, username, password)
	if err != nil {
		log.Fatal(err)
	}

	httpClient := conf.Client(ctx, token)
	client := lade.NewClient(httpClient)

	user, err := client.User.Me()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%+v\n", user)
}

func getCredentials() (string, string) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Username: ")
	username, _ := reader.ReadString('\n')

	fmt.Print("Password: ")
	bytePassword, _ := terminal.ReadPassword(0)
	password := string(bytePassword)
	fmt.Println()

	return strings.TrimSpace(username), strings.TrimSpace(password)
}
```
