package cmd

import (
	"fmt"
	"strconv"
	"strings"

	"deployaja-cli/internal/ui"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(listCmd())
}

func listCmd() *cobra.Command {
	var (
		category  string
		page      int
		limit     int
		sortBy    string
		sortOrder string
		query     string
		author    string
	)

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List all apps in the marketplace",
		Long: `List all applications available in the marketplace.

Examples:
  aja list                              # List all apps (default)
  aja list --query blog                 # Search apps with 'blog'
  aja list --category productivity      # Filter by category
  aja list --author "John Doe"          # Filter by author
  aja list --sort-by downloads          # Sort by downloads
  aja list --page 2 --limit 5           # Pagination

For advanced usage and more examples, see the documentation.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := ensureAuthenticated(); err != nil {
				return err
			}

			// Validate pagination parameters
			if page < 1 {
				return fmt.Errorf("page must be greater than 0")
			}
			if limit < 1 || limit > 100 {
				return fmt.Errorf("limit must be between 1 and 100")
			}

			// Validate sort parameters
			if sortBy != "" {
				validSortFields := []string{"name", "downloads", "rating", "createdAt", "updatedAt"}
				isValid := false
				for _, field := range validSortFields {
					if sortBy == field {
						isValid = true
						break
					}
				}
				if !isValid {
					return fmt.Errorf("invalid sort-by field. Valid options: %s", strings.Join(validSortFields, ", "))
				}
			}

			if sortOrder != "" {
				if sortOrder != "asc" && sortOrder != "desc" {
					return fmt.Errorf("sort-order must be either 'asc' or 'desc'")
				}
			}

			// Build query parameters
			params := make(map[string]string)
			if query != "" {
				params["q"] = query
			}
			if author != "" {
				params["author"] = author
			}
			if category != "" {
				params["category"] = category
			}
			if page > 1 {
				params["page"] = strconv.Itoa(page)
			}
			if limit != 10 { // Default limit is 10
				params["limit"] = strconv.Itoa(limit)
			}
			if sortBy != "" {
				params["sortBy"] = sortBy
			}
			if sortOrder != "" {
				params["sortOrder"] = sortOrder
			}

			// Build display message
			var filters []string
			if query != "" {
				filters = append(filters, fmt.Sprintf("query: %s", query))
			}
			if author != "" {
				filters = append(filters, fmt.Sprintf("author: %s", author))
			}
			if category != "" {
				filters = append(filters, fmt.Sprintf("category: %s", category))
			}
			if page > 1 {
				filters = append(filters, fmt.Sprintf("page: %d", page))
			}
			if limit != 10 {
				filters = append(filters, fmt.Sprintf("limit: %d", limit))
			}
			if sortBy != "" {
				filters = append(filters, fmt.Sprintf("sort: %s %s", sortBy, sortOrder))
			}

			if len(filters) > 0 {
				fmt.Printf("%s Listing marketplace apps with filters: %s\n\n", ui.InfoPrint("📋"), strings.Join(filters, ", "))
			} else {
				fmt.Printf("%s Listing all marketplace apps\n\n", ui.InfoPrint("📋"))
			}

			// Fetch apps from API
			response, err := apiClient.ListMarketplaceApps(params)
			if err != nil {
				return fmt.Errorf("failed to list marketplace apps: %v", err)
			}

			if len(response.Apps) == 0 {
				fmt.Printf("%s No apps found", ui.WarningPrint("⚠️"))
				if len(filters) > 0 {
					fmt.Printf(" with the specified filters")
				}
				fmt.Println()
				return nil
			}

			total := response.Total
			if total == 0 {
				total = len(response.Apps)
			}
			fmt.Printf("%s Found %d apps", ui.SuccessPrint("✅"), total)
			if len(response.Apps) < total {
				fmt.Printf(" (showing %d)", len(response.Apps))
			}
			fmt.Println("\n")

			// Display results in a table format
			for i, app := range response.Apps {
				fmt.Printf("%s %s\n", color.New(color.Bold, color.FgCyan).Sprint(i+1), color.New(color.Bold).Sprint(app.Name))
				fmt.Printf("   %s\n", app.Description)
				fmt.Printf("   Category: %s\n", app.Category)
				fmt.Printf("   Author: %s\n", app.Author)
				fmt.Printf("   Version: %s\n", app.Version)
				fmt.Printf("   Downloads: %d\n", app.Downloads)
				fmt.Printf("   Rating: %.1f\n", app.Rating)

				if len(app.Tags) > 0 {
					fmt.Printf("   Tags: %s\n", strings.Join(app.Tags, ", "))
				}

				if app.Repository != "" {
					fmt.Printf("   Repository: %s\n", app.Repository)
				}

				fmt.Println()
			}

			// Show pagination info if applicable
			if response.Total > len(response.Apps) {
				totalPages := (response.Total + limit - 1) / limit
				fmt.Printf("%s Page %d of %d (showing %d-%d of %d total apps)\n",
					ui.InfoPrint("📄"), page, totalPages,
					(page-1)*limit+1, (page-1)*limit+len(response.Apps), response.Total)
				fmt.Println()
			}

			fmt.Printf("%s Use 'aja install <app-name>' to install an app\n", ui.InfoPrint("💡"))

			return nil
		},
	}

	// Add flags
	cmd.Flags().StringVarP(&query, "query", "q", "", "Search query (name, description, tags)")
	cmd.Flags().StringVar(&author, "author", "", "Filter apps by author")
	cmd.Flags().StringVarP(&category, "category", "c", "", "Filter apps by category")
	cmd.Flags().IntVarP(&page, "page", "p", 1, "Page number (default: 1)")
	cmd.Flags().IntVarP(&limit, "limit", "l", 10, "Number of apps per page (1-100, default: 10)")
	cmd.Flags().StringVar(&sortBy, "sort-by", "", "Sort by field (name, downloads, rating, createdAt, updatedAt)")
	cmd.Flags().StringVar(&sortOrder, "sort-order", "desc", "Sort order (asc, desc)")

	return cmd
}
