package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(publishCmd())
}

func publishCmd() *cobra.Command {
	var (
		name        string
		description string
		category    string
		author      string
		version     string
		repository  string
		image       string
		tags        string
		file        string
	)

	cmd := &cobra.Command{
		Use:   "publish",
		Short: "Publish your app to the marketplace (experimental)",
		Long: `Publish your app to the Aja marketplace.
You can specify metadata via flags or interactively.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := ensureAuthenticated(); err != nil {
				return err
			}

			// If not provided, try to infer from deployaja.yaml or prompt user
			if name == "" {
				fmt.Print("App name: ")
				fmt.Scanln(&name)
			}
			if version == "" {
				fmt.Print("Version: ")
				fmt.Scanln(&version)
			}
			if description == "" {
				fmt.Print("Description: ")
				fmt.Scanln(&description)
			}
			if category == "" {
				fmt.Print("Category: ")
				fmt.Scanln(&category)
			}
			if author == "" {
				fmt.Print("Author: ")
				fmt.Scanln(&author)
			}
			if repository == "" {
				fmt.Print("Repository (optional): ")
				fmt.Scanln(&repository)
			}
			if image == "" {
				fmt.Print("Image (optional): ")
				fmt.Scanln(&image)
			}
			if tags == "" {
				fmt.Print("Tags (comma separated, optional): ")
				fmt.Scanln(&tags)
			}

			tagList := []string{}
			if tags != "" {
				for _, t := range strings.Split(tags, ",") {
					trimmed := strings.TrimSpace(t)
					if trimmed != "" {
						tagList = append(tagList, trimmed)
					}
				}
			}

			// Check config file existence
			configFile := file
			if configFile == "" {
				configFile = "deployaja.yaml"
			}
			if _, err := os.Stat(configFile); err != nil {
				return fmt.Errorf("config file '%s' not found: %v", configFile, err)
			}

			fmt.Printf("📤 Publishing app '%s' (version: %s)...\n", name, version)

			resp, err := apiClient.PublishApp(
				name,
				description,
				category,
				author,
				version,
				repository,
				image,
				tagList,
				configFile,
			)
			if err != nil {
				return fmt.Errorf("failed to publish app: %v", err)
			}

			fmt.Printf("✅ App published! ID: %s\n", resp.ID)
			fmt.Printf("Status: %s\n", resp.Status)
			if resp.PublishedAt != "" {
				fmt.Printf("Published at: %s\n", resp.PublishedAt)
			}
			if resp.Message != "" {
				fmt.Printf("Message: %s\n", resp.Message)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&name, "name", "n", "", "App name")
	cmd.Flags().StringVarP(&version, "version", "v", "", "App version")
	cmd.Flags().StringVarP(&description, "description", "d", "", "App description")
	cmd.Flags().StringVarP(&category, "category", "c", "", "App category")
	cmd.Flags().StringVarP(&author, "author", "a", "", "Author name")
	cmd.Flags().StringVar(&repository, "repository", "", "Repository URL")
	cmd.Flags().StringVar(&image, "image", "", "App image URL")
	cmd.Flags().StringVar(&tags, "tags", "", "Comma separated tags")
	cmd.Flags().StringVarP(&file, "file", "f", "", "Path to deployaja.yaml (default: deployaja.yaml)")

	return cmd
}
