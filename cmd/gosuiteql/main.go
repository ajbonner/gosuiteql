package main

import (
	"flag"
	"fmt"
	"io"
	"os"

	"gosuiteql/internal"
)

func main() {
	// Define command-line flags
	query := flag.String("query", "", "The SQL query to execute")
	queryFile := flag.String("file", "", "File containing the SQL query to execute")
	limit := flag.Int("limit", 0, "Number of results to return")
	offset := flag.Int("offset", 0, "Number of results to skip")
	help := flag.Bool("help", false, "Display help information")
	flag.Parse()

	if *help {
		fmt.Println("gosuiteql - A CLI tool for executing SuiteQL queries")
		fmt.Println("\nUsage:")
		flag.PrintDefaults()
		fmt.Println("\nExamples:")
		fmt.Println("  gosuiteql -query \"SELECT * FROM transaction\"")
		fmt.Println("  gosuiteql -file query.sql")
		os.Exit(0)
	}

	// Get the query from either file or command line
	var queryStr string
	if *queryFile != "" {
		file, err := os.Open(*queryFile)
		if err != nil {
			fmt.Printf("Error opening file: %v\n", err)
			os.Exit(1)
		}
		defer file.Close()

		content, err := io.ReadAll(file)
		if err != nil {
			fmt.Printf("Error reading file: %v\n", err)
			os.Exit(1)
		}
		queryStr = string(content)
	} else if *query != "" {
		queryStr = *query
	} else {
		fmt.Println("Error: Either -query or -file must be specified")
		flag.PrintDefaults()
		os.Exit(1)
	}

	// Initialize SuiteQL client
	client, err := internal.NewSuiteQLClient()
	if err != nil {
		fmt.Printf("Error initializing client: %v\n", err)
		os.Exit(1)
	}

	// Convert limit and offset to pointers for optional parameters
	var limitPtr, offsetPtr *int
	if *limit > 0 {
		limitPtr = limit
	}
	if *offset > 0 {
		offsetPtr = offset
	}

	// Execute the query
	result, err := client.ExecuteQuery(queryStr, limitPtr, offsetPtr)
	if err != nil {
		fmt.Printf("Error executing query: %v\n", err)
		os.Exit(1)
	}

	// Print results
	fmt.Println(result)
}
