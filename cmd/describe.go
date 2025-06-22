package cmd

import (
	"fmt"

	"deployaja-cli/internal/ui"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(describeCmd())
}

func describeCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "describe DEPLOYMENT_NAME",
		Short: "Describe deployment pod details (status, containers, events, etc.)",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			deploymentName := args[0]
			if deploymentName == "" {
				return fmt.Errorf("deployment name is required")
			}
			if err := ensureAuthenticated(); err != nil {
				return err
			}

			fmt.Printf("%s Fetching pod details for %s...\n", ui.InfoPrint("🔍"), deploymentName)

			describeResp, err := apiClient.Describe(deploymentName)
			if err != nil {
				return err
			}

			// Print pod description
			fmt.Printf("%s Pod Description\n\n", ui.InfoPrint("📦"))
			printPodDescription(describeResp.Pod)

			// Print events if any
			if len(describeResp.Events) > 0 {
				fmt.Printf("\n%s Pod Events\n\n", ui.InfoPrint("📅"))
				for _, event := range describeResp.Events {
					fmt.Printf("- [%s] %s: %s (count: %v, first: %v, last: %v)\n",
						event["type"], event["reason"], event["message"],
						event["count"], event["firstTimestamp"], event["lastTimestamp"])
				}
			} else {
				fmt.Printf("%s No events found for this pod\n", ui.WarningPrint("⚠"))
			}

			return nil
		},
	}
	return cmd
}

func printPodDescription(pod map[string]interface{}) {
	fmt.Printf("Name:        %v\n", pod["name"])
	fmt.Printf("Namespace:   %v\n", pod["namespace"])
	fmt.Printf("Node:        %v\n", pod["nodeName"])
	fmt.Printf("Phase:       %v\n", pod["phase"])
	fmt.Printf("Pod IP:      %v\n", pod["podIP"])
	fmt.Printf("Host IP:     %v\n", pod["hostIP"])
	fmt.Printf("Start Time:  %v\n", pod["startTime"])

	// Conditions
	if conds, ok := pod["conditions"].([]interface{}); ok && len(conds) > 0 {
		fmt.Println("Conditions:")
		for _, c := range conds {
			if cond, ok := c.(map[string]interface{}); ok {
				fmt.Printf("  - Type: %v, Status: %v, Reason: %v, Message: %v\n",
					cond["type"], cond["status"], cond["reason"], cond["message"])
			}
		}
	}

	// Containers
	if containers, ok := pod["containers"].([]interface{}); ok && len(containers) > 0 {
		fmt.Println("Containers:")
		for _, cs := range containers {
			if c, ok := cs.(map[string]interface{}); ok {
				fmt.Printf("  - Name: %v\n", c["name"])
				fmt.Printf("    Image: %v\n", c["image"])
				fmt.Printf("    Ready: %v\n", c["ready"])
				fmt.Printf("    Restarts: %v\n", c["restartCount"])
				if state, ok := c["state"].(map[string]interface{}); ok {
					fmt.Printf("    State: %v\n", state)
				}

				// Print ports if available
				if ports, ok := c["ports"].([]interface{}); ok && len(ports) > 0 {
					fmt.Println("    Ports:")
					for _, p := range ports {
						if port, ok := p.(map[string]interface{}); ok {
							containerPort := getValueOrNil(port["containerPort"])
							protocol := getValueOrDefault(port["protocol"], "TCP")
							name := getValueOrNil(port["name"])
							hostPort := getValueOrNil(port["hostPort"])
							fmt.Printf("      - Container Port: %v, Protocol: %v, Name: %v, Host Port: %v\n",
								containerPort, protocol, name, hostPort)
						}
					}
				}

				// Print environment variables if available
				if env, ok := c["environment"].([]interface{}); ok && len(env) > 0 {
					fmt.Println("    Environment Variables:")
					for _, e := range env {
						if envVar, ok := e.(map[string]interface{}); ok {
							name := getValueOrNil(envVar["name"])
							value := getValueOrNil(envVar["value"])
							fmt.Printf("      - %v = %v\n", name, value)

							// Print valueFrom if available
							if valueFrom, ok := envVar["valueFrom"].(map[string]interface{}); ok {
								if configMapRef := getValueOrNil(valueFrom["configMapKeyRef"]); configMapRef != "<nil>" {
									fmt.Printf("        (from ConfigMap: %v)\n", configMapRef)
								}
								if secretRef := getValueOrNil(valueFrom["secretKeyRef"]); secretRef != "<nil>" {
									fmt.Printf("        (from Secret: %v)\n", secretRef)
								}
								if fieldRef := getValueOrNil(valueFrom["fieldRef"]); fieldRef != "<nil>" {
									fmt.Printf("        (from Field: %v)\n", fieldRef)
								}
								if resourceFieldRef := getValueOrNil(valueFrom["resourceFieldRef"]); resourceFieldRef != "<nil>" {
									fmt.Printf("        (from Resource Field: %v)\n", resourceFieldRef)
								}
							}
						}
					}
				}

				// Print volume mounts if available
				if mounts, ok := c["volumeMounts"].([]interface{}); ok && len(mounts) > 0 {
					fmt.Println("    Volume Mounts:")
					for _, m := range mounts {
						if mount, ok := m.(map[string]interface{}); ok {
							name := getValueOrNil(mount["name"])
							mountPath := getValueOrNil(mount["mountPath"])
							readOnly := getValueOrDefault(mount["readOnly"], false)
							fmt.Printf("      - Name: %v, Mount Path: %v, Read Only: %v\n",
								name, mountPath, readOnly)
						}
					}
				}
			}
		}
	}
}

// Helper function to get value or return "<nil>"
func getValueOrNil(value interface{}) string {
	if value == nil {
		return "<nil>"
	}
	return fmt.Sprintf("%v", value)
}

// Helper function to get value or return default
func getValueOrDefault(value interface{}, defaultValue interface{}) interface{} {
	if value == nil {
		return defaultValue
	}
	return value
}
