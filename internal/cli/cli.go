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
	flagServerURL = "server"
	flagApiKey    = "key"
)

func Run() error {
	return parseCmd().Execute()
}

type config struct {
	serverUrl string
	apiKey    string
}

func parseCmd() *cobra.Command {
	rootCmd := &cobra.Command{Use: "outline"}

	var cfg config
	rootCmd.PersistentFlags().StringVar(&cfg.serverUrl, flagServerURL, "", "The outline API server url")
	rootCmd.PersistentFlags().StringVar(&cfg.apiKey, flagApiKey, "", "The outline api key")

	collectionCmd := &cobra.Command{
		Use:   "collection",
		Short: "Work with collections",
		Long:  `If you have to work with collection in any case, use this command`,
		Args:  cobra.MinimumNArgs(1),
	}

	collectionInfoCmd := &cobra.Command{
		Use:   "info",
		Short: "Get collection info",
		Long:  "Get information about the collection",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return collectionInfo(cfg.serverUrl, cfg.apiKey, args)
		},
	}

	collectionCreateCmd := &cobra.Command{
		Use:   "create",
		Short: "Creates a collection",
		Long:  "Creates a collection with the given name and prints the result as json to stdout",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return collectionCreate(cfg.serverUrl, cfg.apiKey, args[0])
		},
	}

	collectionDocumentsCmd := &cobra.Command{
		Use:   "docs",
		Short: "Get document structure",
		Long:  "Get a summary of associated documents (and children)",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return collectionDocuments(cfg.serverUrl, cfg.apiKey, args)
		},
	}

	rootCmd.AddCommand(collectionCmd)
	collectionCmd.AddCommand(collectionInfoCmd)
	collectionCmd.AddCommand(collectionCreateCmd)
	collectionCmd.AddCommand(collectionDocumentsCmd)

	return rootCmd
}

func collectionInfo(serverUrl string, apiKey string, ids []string) error {
	oc := outline.New(serverUrl, &http.Client{}, apiKey)
	for _, id := range ids {
		col, err := oc.Collections().Get(outline.CollectionID(id)).Do(context.Background())
		if err != nil {
			return fmt.Errorf("can't get collection with id '%s': %w", id, err)
		}

		b, err := json.MarshalIndent(col, "", "  ")
		if err != nil {
			return fmt.Errorf("failed marshalling collection with id '%s: %w", id, err)
		}
		fmt.Println(string(b))
	}
	return nil
}

func collectionCreate(serverUrl string, apiKey string, name string) error {
	oc := outline.New(serverUrl, &http.Client{}, apiKey)
	col, err := oc.Collections().Create(name).Do(context.Background())
	if err != nil {
		return fmt.Errorf("can't create collection with name '%s': %w", name, err)
	}

	b, err := json.MarshalIndent(col, "", "  ")
	if err != nil {
		return fmt.Errorf("failed marshalling collection with name '%s: %w", name, err)
	}
	fmt.Println(string(b))

	return nil
}

func collectionDocuments(serverUrl string, apiKey string, ids []string) error {
	oc := outline.New(serverUrl, &http.Client{}, apiKey)
	for _, id := range ids {
		st, err := oc.Collections().DocumentStructure(outline.CollectionID(id)).Do(context.Background())
		if err != nil {
			return fmt.Errorf("can't get collection with id '%s': %w", id, err)
		}

		b, err := json.MarshalIndent(st, "", "  ")
		if err != nil {
			return fmt.Errorf("failed marshalling collection with id '%s: %w", id, err)
		}
		fmt.Println(string(b))
	}
	return nil
}
