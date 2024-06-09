package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/spf13/cobra"
)

func printTimeCmd() *cobra.Command {
	return &cobra.Command{
		Use: "curtime",
		RunE: func(cmd *cobra.Command, args []string) error {
			now := time.Now()
			prettyTime := now.Format(time.RubyDate)
			cmd.Println("Hey! The current time is ", prettyTime)
			return nil //Tells cobra that no errors happened
		},
	}
}

func getBooksCmd() *cobra.Command {
	return &cobra.Command{
		Use: "books",
		RunE: func(cmd *cobra.Command, args []string) error {
			resp, err := http.Get("http://localhost:8080/books")
			if err != nil {
				log.Fatal(err)
			}
			defer resp.Body.Close()
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println(string(body))
			return nil //Tells cobra that no errors happened
		},
	}
}

func getBookCmd() *cobra.Command {
	return &cobra.Command{
		Use: "book",
		RunE: func(cmd *cobra.Command, args []string) error {
			resp, err := http.Get("http://localhost:8080/books/" + args[0])
			if err != nil {
				log.Fatal(err)
			}
			defer resp.Body.Close()
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println(string(body))
			return nil //Tells cobra that no errors happened
		},
	}
}

func addBookCmd() *cobra.Command {
	return &cobra.Command{
		Use: "addbook",
		RunE: func(cmd *cobra.Command, args []string) error {
			payload, err := json.Marshal(map[string]interface{}{
				"id":       4,
				"title":    "Colors for adults",
				"author":   "Adeline",
				"quantity": 2,
			})
			if err != nil {
				log.Fatal(err)
			}
			client := &http.Client{}
			url := "http://localhost:8080/books"
			req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(payload))
			req.Header.Set("Content-Type", "application/json")
			if err != nil {
				log.Fatal(err)
			}
			resp, err := client.Do(req)
			if err != nil {
				log.Fatal(err)
			}
			defer resp.Body.Close()
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println(string(body))
			return nil
		},
	}
}

func checkoutBookCmd() *cobra.Command {
	return &cobra.Command{
		Use: "checkout",
		RunE: func(cmd *cobra.Command, args []string) error {
			client := &http.Client{}
			url := "http://localhost:8080/checkout?id=" + args[0]
			req, err := http.NewRequest(http.MethodPatch, url, nil)
			req.Header.Set("Content-Type", "application/json")
			if err != nil {
				log.Fatal(err)
			}
			resp, err := client.Do(req)
			if err != nil {
				log.Fatal(err)
			}
			defer resp.Body.Close()
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println(string(body))
			return nil
		},
	}
}

func returnBookCmd() *cobra.Command {
	return &cobra.Command{
		Use: "return",
		RunE: func(cmd *cobra.Command, args []string) error {
			client := &http.Client{}
			url := "http://localhost:8080/return?id=" + args[0]
			req, err := http.NewRequest(http.MethodPatch, url, nil)
			req.Header.Set("Content-Type", "application/json")
			if err != nil {
				log.Fatal(err)
			}
			resp, err := client.Do(req)
			if err != nil {
				log.Fatal(err)
			}
			defer resp.Body.Close()
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println(string(body))
			return nil
		},
	}
}

func main() {
	cmd := &cobra.Command{
		Use:   "gifm",
		Short: "Welcome to the Client!",
	}

	cmd.AddCommand(printTimeCmd())
	cmd.AddCommand(getBooksCmd())
	cmd.AddCommand(getBookCmd())
	cmd.AddCommand(addBookCmd())
	cmd.AddCommand(checkoutBookCmd())
	cmd.AddCommand(returnBookCmd())

	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
