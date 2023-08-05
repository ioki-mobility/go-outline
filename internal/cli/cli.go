package cli

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/ioki-mobility/go-outline"
	"github.com/spf13/cobra"
)

const (
	flagNameBaseUrl = "baseUrl"
	flagNameApiKey  = "apiKey"
)

func Run() error {
	rootCmd := rootCmd()
	collectionCmd := collectionCmd(rootCmd)
	rootCmd.AddCommand(collectionCmd)

	return rootCmd.Execute()
}

func rootCmd() *cobra.Command {
	var rootCmd = &cobra.Command{Use: "outline"}

	var apiKeyFlag string
	var baseUrlFlag string
	rootCmd.PersistentFlags().StringVar(&baseUrlFlag, flagNameBaseUrl, "", "The base url to the outline API instance (should include /api/ as suffix)")
	rootCmd.PersistentFlags().StringVar(&apiKeyFlag, flagNameApiKey, "", "The outline apiKey")

	return rootCmd
}

func collectionCmd(rootCmd *cobra.Command) *cobra.Command {
	collectionCmd := &cobra.Command{
		Use:   "collection",
		Short: "Work with collections",
		Long:  `If you have to work with collection in any case, use this command`,
		Args:  cobra.MinimumNArgs(1),
	}

	fetchSubCmd := collectionCmdFetch(rootCmd)
	collectionCmd.AddCommand(fetchSubCmd)

	createSubCmd := collectionCmdCrate(rootCmd)
	collectionCmd.AddCommand(createSubCmd)

	return collectionCmd
}

func collectionCmdFetch(rootCmd *cobra.Command) *cobra.Command {
	return &cobra.Command{
		Use:   "fetch",
		Short: "Fetch collection details",
		Long:  "Fetch collection details for given ID and prints as json to stdout",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			baseUrl, err := rootCmd.Flags().GetString(flagNameBaseUrl)
			if err != nil {
				return fmt.Errorf("required flag '%s' not set: %w", flagNameBaseUrl, err)
			}
			apiKey, err := rootCmd.Flags().GetString(flagNameApiKey)
			if err != nil {
				return fmt.Errorf("required flag '%s' not set: %w", flagNameApiKey, err)
			}
			client := outline.New(baseUrl, &http.Client{}, apiKey)
			for _, colId := range args {
				col, err := client.Collections().Get(outline.CollectionID(colId)).Do(context.Background())
				if err != nil {
					return fmt.Errorf("can't get collection with id '%s': %w", colId, err)
				}

				b, err := json.Marshal(&col)
				if err != nil {
					return fmt.Errorf("failed marshalling collection with id '%s: %w", colId, err)
				}
				fmt.Println(string(b))
			}
			return nil
		},
	}
}

func collectionCmdCrate(rootCmd *cobra.Command) *cobra.Command {
	return &cobra.Command{
		Use:   "create",
		Short: "Creates a collection",
		Long:  "Creates a collection with the given name and prints the result as json to stdout",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			baseUrl, err := rootCmd.Flags().GetString(flagNameBaseUrl)
			if err != nil {
				return fmt.Errorf("required flag '%s' not set: %w", flagNameBaseUrl, err)
			}
			apiKey, err := rootCmd.Flags().GetString(flagNameApiKey)
			if err != nil {
				return fmt.Errorf("required flag '%s' not set: %w", flagNameApiKey, err)
			}
			client := outline.New(baseUrl, &http.Client{}, apiKey)

			name := args[0]
			col, err := client.Collections().Create(name).Do(context.Background())
			if err != nil {
				return fmt.Errorf("can't create collection with name '%s': %w", name, err)
			}

			b, err := json.Marshal(&col)
			if err != nil {
				return fmt.Errorf("failed marshalling collection with name '%s: %w", name, err)
			}
			fmt.Println(string(b))

			return nil
		},
	}
}
