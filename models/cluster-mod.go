package models

import (
	"time"
)

type NodesData struct {
	Kind       string `json:"kind"`
	APIVersion string `json:"apiVersion"`
	Metadata   struct {
		SelfLink string `json:"selfLink"`
	} `json:"metadata"`
	Items []struct {
		Metadata struct {
			Name              string    `json:"name"`
			SelfLink          string    `json:"selfLink"`
			CreationTimestamp time.Time `json:"creationTimestamp"`
		} `json:"metadata"`
		Timestamp time.Time `json:"timestamp"`
		Window    string    `json:"window"`
		Usage     struct {
			CPU    string `json:"cpu"`
			Memory string `json:"memory"`
		} `json:"usage"`
	} `json:"items"`
}

type PodsData struct {
	Kind       string `json:"kind"`
	APIVersion string `json:"apiVersion"`
	Metadata   struct {
		Name              string    `json:"name"`
		SelfLink          string    `json:"selfLink"`
		CreationTimestamp time.Time `json:"creationTimestamp"`
	} `json:"metadata"`
	Timestamp time.Time `json:"timestamp"`
	Window    string    `json:"window"`
	Usage     struct {
		CPU    string `json:"cpu"`
		Memory string `json:"memory"`
	} `json:"usage"`
}