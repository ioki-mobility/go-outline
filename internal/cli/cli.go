package cli

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/ioki-mobility/go-outline"
	"github.com/ioki-mobility/go-outline/internal/common"
	"github.com/spf13/cobra"
)

const (
	flagServerURL = "server"
	flagApiKey    = "key"
)

type config struct {
	serverUrl string
	apiKey    string
}

func Command() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:               "outcli",
		DisableAutoGenTag: true,
	}

	var cfg config
	rootCmd.PersistentFlags().StringVar(&cfg.serverUrl, flagServerURL, "", "The outline API server url")
	rootCmd.PersistentFlags().StringVar(&cfg.apiKey, flagApiKey, "", "The outline api key")

	rootCmd.AddCommand(&cobra.Command{
		Use: "version",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(common.Version)
		},
		Short: "Show app version",
	})

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

	collectionListCmd := &cobra.Command{
		Use:   "list",
		Short: "List all collections",
		Long:  "Get a list of all collections.",
		RunE: func(cmd *cobra.Command, args []string) error {
			return collectionList(cfg.serverUrl, cfg.apiKey)
		},
	}

	documentCmd := &cobra.Command{
		Use:   "document",
		Short: "Work with documents",
		Long:  `If you have to work with documents in any case, use this command`,
		Args:  cobra.MinimumNArgs(1),
	}

	documentCreateCmd := &cobra.Command{
		Use:   "create",
		Short: "Creates a document",
		Long:  "Creates a collection with the given name and collection id",
		Args:  cobra.MinimumNArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return documentCreate(cfg.serverUrl, cfg.apiKey, args[0], outline.CollectionID(args[1]))
		},
	}

	var isShareID bool
	documentGetCmd := &cobra.Command{
		Use:   "get",
		Short: "Get an existing document by its ID",
		Long:  "Get information about an existing document by specifying its document ID or a share ID",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return documentGet(cfg.serverUrl, cfg.apiKey, args[0], isShareID)
		},
	}
	documentGetCmd.Flags().BoolVar(&isShareID, "share", false, "Treat the argument as document share iD")

	rootCmd.AddCommand(collectionCmd)
	collectionCmd.AddCommand(collectionInfoCmd)
	collectionCmd.AddCommand(collectionCreateCmd)
	collectionCmd.AddCommand(collectionDocumentsCmd)
	collectionCmd.AddCommand(collectionListCmd)
	collectionCmd.AddCommand(collectionUpdate())
	rootCmd.AddCommand(documentCmd)
	documentCmd.AddCommand(documentCreateCmd)
	documentCmd.AddCommand(documentGetCmd)
	documentCmd.AddCommand(documentUpdate())

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

func collectionList(serverUrl string, apiKey string) error {
	oc := outline.New(serverUrl, &http.Client{}, apiKey)
	err := oc.Collections().List().Do(context.Background(), func(col *outline.Collection, err error) (bool, error) {
		if err != nil {
			return false, err
		}

		b, err := json.MarshalIndent(col, "", "  ")
		if err != nil {
			return false, fmt.Errorf("failed marshalling collection: %w", err)
		}
		fmt.Println(string(b))

		return true, nil
	})
	if err != nil {
		return fmt.Errorf("can't get list of collections: %w", err)
	}
	return nil
}

func collectionUpdate() *cobra.Command {
	var name, description, color string
	var permissionRead, permissionReadWrite bool

	cmd := &cobra.Command{
		Use:   "update",
		Short: "Update an existing collection",
		Long:  "Update an existing collection's name, description etc. properties",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(c *cobra.Command, args []string) error {
			id := args[0]
			errBase := fmt.Sprintf("failed updating collection with ID '%s'", id)

			// Extract value of global flags
			key, err := c.Flags().GetString(flagApiKey)
			if err != nil {
				return fmt.Errorf("%s: %w", errBase, err)
			}
			url, err := c.Flags().GetString(flagServerURL)
			if err != nil {
				return fmt.Errorf("%s: %w", errBase, err)
			}

			cl := outline.New(url, &http.Client{}, key).
				Collections().
				Update(outline.CollectionID(id)).
				Name(name).
				Description(description).
				Color(color)

			if permissionRead {
				cl.PermissionRead()
			}
			if permissionReadWrite {
				cl.PermissionReadWrite()
			}

			doc, err := cl.Do(context.Background())
			if err != nil {
				return fmt.Errorf("%s: %w", errBase, err)
			}

			b, err := json.MarshalIndent(doc, "", "  ")
			if err != nil {
				return fmt.Errorf("failed marshalling collection with ID '%s': %w", id, err)
			}
			fmt.Println(string(b))

			return nil
		},
	}

	cmd.Flags().StringVar(&name, "name", "", "The name of the collection")
	cmd.Flags().StringVar(&description, "description", "", "The description of the collection")
	cmd.Flags().StringVar(&color, "color", "", "The color of the collection. Should be in the format of #AABBCC")
	cmd.Flags().BoolVar(&permissionRead, "permission-read", false, "Change the permission to read only")
	cmd.Flags().BoolVar(&permissionReadWrite, "permission-read-write", false, "Change the permission to read write")
 cmd.MarkFlagsMutuallyExclusive("permission-read", "permission-read-write")

	return cmd
}

func documentCreate(serverUrl string, apiKey string, name string, collectionId outline.CollectionID) error {
	oc := outline.New(serverUrl, &http.Client{}, apiKey)
	doc, err := oc.Documents().Create(name, collectionId).Do(context.Background())
	if err != nil {
		return fmt.Errorf("can't create document with name '%s': %w", name, err)
	}

	b, err := json.MarshalIndent(doc, "", "  ")
	if err != nil {
		return fmt.Errorf("failed marshalling document with name '%s: %w", name, err)
	}
	fmt.Println(string(b))

	return nil
}

func documentGet(serverURL string, apiKey string, id string, idIsShareID bool) error {
	var err error
	var doc *outline.Document
	oc := outline.New(serverURL, &http.Client{}, apiKey)

	if idIsShareID {
		doc, err = oc.Documents().Get().ByShareID(outline.DocumentShareID(id)).Do(context.Background())
		if err != nil {
			return fmt.Errorf("can't get document with share id '%s': %w", id, err)
		}
	} else {
		doc, err = oc.Documents().Get().ByID(outline.DocumentID(id)).Do(context.Background())
		if err != nil {
			return fmt.Errorf("can't get document with id '%s': %w", id, err)
		}
	}

	b, err := json.MarshalIndent(doc, "", "  ")
	if err != nil {
		return fmt.Errorf("failed marshalling document with id '%s': %w", id, err)
	}
	fmt.Println(string(b))

	return nil
}

func documentUpdate() *cobra.Command {
	var title string
	var append, publish, readText bool

	cmd := &cobra.Command{
		Use:   "update",
		Short: "Update an existing document",
		Long:  "Update an existing document's title, text etc. properties",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(c *cobra.Command, args []string) error {
			id := args[0]
			errBase := fmt.Sprintf("failed updating document with ID '%s'", id)

			// Extract value of global flags
			key, err := c.Flags().GetString(flagApiKey)
			if err != nil {
				return fmt.Errorf("%s: %w", errBase, err)
			}
			url, err := c.Flags().GetString(flagServerURL)
			if err != nil {
				return fmt.Errorf("%s: %w", errBase, err)
			}

			cl := outline.New(url, &http.Client{}, key).Documents().Update(outline.DocumentID(id)).
				Append(append).Publish(publish).Title(title)

			if readText {
				b, err := io.ReadAll(os.Stdin)
				if err != nil {
					return fmt.Errorf("%s: %w", errBase, err)
				}
				cl.Text(string(b))
			}

			doc, err := cl.Do(context.Background())
			if err != nil {
				return fmt.Errorf("%s: %w", errBase, err)
			}

			b, err := json.MarshalIndent(doc, "", "  ")
			if err != nil {
				return fmt.Errorf("failed marshalling document with ID '%s': %w", id, err)
			}
			fmt.Println(string(b))

			return nil
		},
	}

	cmd.Flags().StringVar(&title, "title", "", "The title of the document")
	cmd.Flags().BoolVar(&readText, "text", false, "Read document text from stdin")
	cmd.Flags().BoolVar(&append, "append", false, "Append new text to existing rather than replacing it")
	cmd.Flags().BoolVar(&publish, "publish", false,
		"Whether this document should be published and made visible to other team members, if a draft",
	)

	return cmd
}
