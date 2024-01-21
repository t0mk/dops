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
	table.SetHeader([]string{"Container ID", "Container Name", "Image", "Status", "IP Address", "Open Ports"})

	// Fill in table with container information
	for _, container := range containers {
		// Get container details to retrieve port information
		containerDetails, err := cli.ContainerInspect(context.Background(), container.ID)
		if err != nil {
			log.Fatal(err)
		}

		// Get the container's IP address
		ipAddress := getContainerIPAddress(containerDetails)

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
			ipAddress,
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

// Helper function to get the container's IP address
func getContainerIPAddress(containerDetails types.ContainerJSON) string {
	// Assuming the container has a single network (common case)
	if len(containerDetails.NetworkSettings.Networks) > 0 {
		for _, network := range containerDetails.NetworkSettings.Networks {
			return network.IPAddress
		}
	}
	return "N/A"
}
