package main

import (
	"encoding/json"
	"fmt"
	"time"
)

type Event struct {
	Name      string                 `json:"name"`
	Timestamp time.Time              `json:"timestamp"`
	Data      map[string]interface{} `json:"data"`
}

func LogEvent(name string, data map[string]interface{}) {
	event := Event{
		Name:      name,
		Timestamp: time.Now(),
		Data:      data,
	}

	jsonData, err := json.Marshal(event)
	if err != nil {
		fmt.Println("Error marshalling event:", err)
		return
	}

	fmt.Println(string(jsonData))
}
