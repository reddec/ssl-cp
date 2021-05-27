package service

import (
	"context"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/reddec/ssl-cp/api"
	"gopkg.in/yaml.v2"
)

// ImportFromDir imports all yaml or JSON files recursively from directory to destination api.
// If dir is empty string, function returns without error and scanning.
func ImportFromDir(ctx context.Context, dir string, destination api.API) error {
	if dir == "" {
		return nil
	}
	var extensions = []string{".yaml", ".yml", ".json"}
	return filepath.Walk(dir, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}

		var matched bool
		fileExt := strings.ToLower(filepath.Ext(info.Name()))
		for _, ext := range extensions {
			if fileExt == ext {
				matched = true
				break
			}
		}

		if !matched {
			return nil
		}

		f, err := os.Open(path)
		if err != nil {
			return fmt.Errorf("open file %s: %w", path, err)
		}
		defer f.Close()

		var batched []api.Batch
		dec := yaml.NewDecoder(f)

		for {
			var item api.Batch
			err := dec.Decode(&item)
			if errors.Is(err, io.EOF) {
				break
			}
			if err != nil {
				return fmt.Errorf("decode file %s: %w", path, err)
			}
			batched = append(batched, item)
		}

		_, err = destination.BatchCreateCertificate(ctx, batched)
		if err != nil {
			return fmt.Errorf("create batch from %s: %w", path, err)
		}
		return nil
	})
}

func nullId(id uint) *uint {
	if id == 0 {
		return nil
	}
	return &id
}
