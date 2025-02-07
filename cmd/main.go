package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

type LogResponse struct {
	Data []map[string]interface{} `json:"data"`
	Meta struct {
		Page struct {
			After string `json:"after"`
		} `json:"page"`
	} `json:"meta"`
}

func fetchLogs(storageTier, fromTime, toTime, query, apiUrl, output string) error {
	log.Printf("Starting to fetch logs from %s to %s\n", fromTime, toTime)

	apiKey := os.Getenv("DD_API_KEY")
	appKey := os.Getenv("DD_APP_KEY")

	if apiKey == "" || appKey == "" {
		return fmt.Errorf("please set the DD_API_KEY and DD_APP_KEY environment variables")
	}

	client := &http.Client{Timeout: 30 * time.Second}
	cursor := ""
	totalLogs := 0
	file, err := os.Create(output)
	if err != nil {
		return err
	}
	defer file.Close()

	bufWriter := bufio.NewWriter(file) // Usando um buffer para reduzir chamadas de write
	defer bufWriter.Flush()

	bufWriter.WriteString("[\n")
	firstEntry := true

	for {
		payload := fmt.Sprintf(`
{
    "filter": {
        "from": "%s",
        "to": "%s",
        "query": "%s",
        "storage_tier": "%s"
    },
    "sort": "timestamp:desc",
    "page": {
        "limit": 5000%s
    }
}`, fromTime, toTime, query, storageTier, func() string {
			if cursor != "" {
				return fmt.Sprintf(`,"cursor": "%s"`, cursor)
			}
			return ""
		}())

		req, err := http.NewRequest("POST", apiUrl, io.NopCloser(strings.NewReader(payload)))
		if err != nil {
			return err
		}

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("DD-API-KEY", apiKey)
		req.Header.Set("DD-APPLICATION-KEY", appKey)

		resp, err := client.Do(req)
		if err != nil {
			return err
		}

		defer resp.Body.Close()

		if resp.StatusCode == http.StatusTooManyRequests {
			log.Println("Rate limit reached, retrying in 30 seconds...")
			time.Sleep(30 * time.Second)
			continue // Re-tenta a mesma requisição
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return err
		}

		if resp.StatusCode != http.StatusOK {
			return fmt.Errorf("unexpected status code: %d - %s", resp.StatusCode, string(body))
		}

		var logResponse LogResponse
		err = json.Unmarshal(body, &logResponse)
		if err != nil {
			return err
		}

		batchSize := len(logResponse.Data)
		totalLogs += batchSize
		log.Printf("Fetched %d logs\n", totalLogs)

		for _, logEntry := range logResponse.Data {
			if !firstEntry {
				bufWriter.WriteString(",\n")
			}
			json.NewEncoder(bufWriter).Encode(logEntry)
			firstEntry = false
		}

		cursor = logResponse.Meta.Page.After
		if cursor == "" {
			break
		}
	}

	bufWriter.WriteString("\n]")
	log.Printf("✅ Fetched a total of %d logs\n", totalLogs)
	log.Printf("📁 Logs saved at: %s\n", output)
	return nil
}

func main() {
	apiUrlCollection := map[string]string{
		"us3":     "https://api.us3.datadoghq.com/api/v2/logs/events/search",
		"us5":     "https://api.us5.datadoghq.com/api/v2/logs/events/search",
		"us1":     "https://api.datadoghq.com/api/v2/logs/events/search",
		"ap1":     "https://api.ap1.datadoghq.com/api/v2/logs/events/search",
		"eu":      "https://api.datadoghq.eu/api/v2/logs/events/search",
		"us1-fed": "https://api.ddog-gov.com/api/v2/apicatalog/api",
	}

	storageTier := flag.String("storage_tier", "indexes", "Storage tier (ex: indexes,online-archives,flex) Default: indexes")
	fromTime := flag.String("from", "", "Start date (ex: 2024-12-25T00:00:00Z)")
	toTime := flag.String("to", "", "End date (ex: 2024-12-25T23:59:59Z)")
	query := flag.String("query", "", "Query (ex: source:auth0)")
	apiRegion := flag.String("api_region", "us1", "API region (ex: ap1, us1, us3, us5, eu, us1-fed) Default: us1")
	output := flag.String("output", "logs.json", "Output file (ex: my_logs.json) Default: logs.json")

	flag.Parse()

	if *fromTime == "" || *toTime == "" {
		fmt.Println("Error: '--from' and '--to' arguments are required.")
		flag.Usage()
		os.Exit(1)
	}

	apiUrl, ok := apiUrlCollection[*apiRegion]
	if !ok {
		fmt.Println("Error: Invalid API region.")
		flag.Usage()
		os.Exit(1)
	}

	err := fetchLogs(*storageTier, *fromTime, *toTime, *query, apiUrl, *output)
	if err != nil {
		log.Fatal(err)
	}
}
