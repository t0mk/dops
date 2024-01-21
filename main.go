package main

import (
	"context"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/olekukonko/tablewriter"
	"log"
	"os"
)

func main() {
	// Create a Docker client
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		log.Fatal(err)
	}

	// List containers
	containers, err := cli.ContainerList(context.Background(), types.ContainerListOptions{All: true})
	if err != nil {
		log.Fatal(err)
	}

	// Prepare table
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Container ID", "Container Name", "Image", "Status", "Open Ports"})

	// Fill in table with container information
	for _, container := range containers {
		// Get container details to retrieve port information
		containerDetails, err := cli.ContainerInspect(context.Background(), container.ID)
		if err != nil {
			log.Fatal(err)
		}

		// Prepare open ports string
		openPorts := ""
		for port, _ := range containerDetails.Config.ExposedPorts {
			openPorts += fmt.Sprintf("%s, ", port)
		}
		openPorts = trimCommaSpace(openPorts)

		// Add row to the table
		table.Append([]string{
			container.ID[:12],
			container.Names[0],
			container.Image,
			container.Status,
			openPorts,
		})
	}

	// Render the table
	table.Render()
}

// Helper function to trim the trailing comma and space
func trimCommaSpace(s string) string {
	if len(s) > 2 {
		return s[:len(s)-2]
	}
	return s
}
