# ddlogexporter

`ddlogexporter` is a tool to fetch logs from Datadog from a specified region and save the data in a JSON file.

## Features

- Allows fetching logs from different Datadog API regions (e.g., `us1`, `us3`, `us5`, `eu`).
- Supports different storage tiers (such as `indexes`, `online-archives`, `flex`).
- Queries logs based on a time range (`from` and `to`), and optionally, a `query` filter.
- Saves the logs in a formatted JSON file.

## Prerequisites

- You need to set up an API key (`DD_API_KEY`) and an application key (`DD_APP_KEY`) in your environment. Otherwise, the tool will not work.

## How to Use

### Building

To build the binary for your desired platform, you can use the provided Makefile:

```bash
make all
```

### Running 

Once you've built the binary, run the tool with the necessary parameters:

```bash
# Example usage

./bin/linux/ddlogexporter --from "2024-12-25T00:00:00Z" --to "2024-12-25T10:05:59Z" --query "source:auth0" --output "/tmp/my_logs.json" --api_region "us3" --storage_tier flex
```

- --from (required): Start date and time for the search (e.g., 2024-12-25T00:00:00Z).
- --to (required): End date and time for the search (e.g., 2024-12-25T23:59:59Z).
- --query (optional): Query filter for logs (e.g., source:auth0).
- --storage_tier (optional): Storage tier (e.g., indexes, online-archives, flex) (default: indexes).
- --api_region (optional): Datadog API region to use (e.g., us1, us3, us5, eu, us1-fed) (default: us1).
- --output (optional): Output file name where the logs will be saved (default: logs.json).

### Expected output

```bash
2025/02/07 12:15:35 Starting to fetch logs from 2024-12-25T00:00:00Z to 2024-12-25T10:05:59Z
2025/02/07 12:15:39 Fetched 5000 logs
2025/02/07 12:15:42 Fetched 10000 logs
2025/02/07 12:15:45 Fetched 13922 logs
2025/02/07 12:15:45 ✅ Fetched a total of 13922 logs
2025/02/07 12:15:45 📁 Logs saved at: /tmp/my_logs.json
```