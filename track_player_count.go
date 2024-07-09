package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

const (
	apiURL         = "https://api.mcsrvstat.us/3/mc.micr0.dev"
	outputFilePath = "./minecraft_player_counts.json"
)

type APIResponse struct {
	Online  bool `json:"online"`
	Players struct {
		Online int `json:"online"`
	} `json:"players"`
}

type PlayerCount struct {
	Timestamp int64 `json:"timestamp"`
	Count     int   `json:"count"`
}

func main() {
	for {
		recordPlayerCount()
		time.Sleep(1 * time.Minute) // Adjust the interval as needed
	}
}

func recordPlayerCount() {
	resp, err := http.Get(apiURL)
	if err != nil {
		fmt.Println("Error fetching player count:", err)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return
	}

	var data APIResponse
	err = json.Unmarshal(body, &data)
	if err != nil {
		fmt.Println("Error unmarshalling JSON:", err)
		return
	}

	if !data.Online {
		fmt.Println("Server is offline")
		return
	}

	playerCount := PlayerCount{
		Timestamp: time.Now().Unix(),
		Count:     data.Players.Online,
	}

	writePlayerCountToFile(playerCount)
}

func writePlayerCountToFile(playerCount PlayerCount) {
	var currentData []PlayerCount

	file, err := os.OpenFile(outputFilePath, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	fileContents, err := io.ReadAll(file)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return
	}

	if len(fileContents) > 0 {
		err = json.Unmarshal(fileContents, &currentData)
		if err != nil {
			fmt.Println("Error unmarshalling existing data:", err)
			return
		}
	}

	currentData = append(currentData, playerCount)

	newData, err := json.MarshalIndent(currentData, "", "    ")
	if err != nil {
		fmt.Println("Error marshalling new data:", err)
		return
	}

	err = os.WriteFile(outputFilePath, newData, 0644)
	if err != nil {
		fmt.Println("Error writing to file:", err)
		return
	}
}
