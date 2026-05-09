package main

import (
	"cloud-manage/service"
	"fmt"
	"os"
)

// CLI debug entry point for testing ECS queries without the GUI.
// Usage: go run cmd/cli/main.go <accessKeyId> <accessKeySecret> <region>
func main() {
	if len(os.Args) < 4 {
		fmt.Println("Usage: go run cmd/cli/main.go <accessKeyId> <accessKeySecret> <region>")
		os.Exit(1)
	}

	accessKeyId := os.Args[1]
	accessKeySecret := os.Args[2]
	region := os.Args[3]

	svc := service.NewECSService()
	result, err := svc.ListInstances(accessKeyId, accessKeySecret, region)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Total instances: %d\n", result.TotalCount)
	for _, inst := range result.Instances {
		fmt.Printf("  [%s] %s (%s) - %s | Public: %s | Private: %s | Created: %s\n",
			inst.Status, inst.InstanceId, inst.InstanceName,
			inst.ZoneId, inst.PublicIp, inst.PrivateIp, inst.CreationTime)
	}
}
