//
// Copyright (c) 2022-2023 Markku Rossi
//
// All rights reserved.
//

package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"strings"
)

// Assets define static assets from a directory tree.
type Assets struct {
	root  string
	files map[string]os.DirEntry
}

// NewAssets creates a new assets object for the specified directory
// tree root.
func NewAssets(root string) *Assets {
	return &Assets{
		root:  root,
		files: make(map[string]os.DirEntry),
	}
}

// Dir returns the assets root directory
func (assets *Assets) Dir() string {
	return assets.root
}

// Add adds a file to the assets. The file must be located under the
// assets object's root directory.
func (assets *Assets) Add(file string, asset os.DirEntry) error {
	if !strings.HasPrefix(file, assets.root) {
		return fmt.Errorf("file '%s' is not under root directory '%s'",
			file, assets.root)
	}
	assets.files[file] = asset
	return nil
}

// AddDir recursively adds the assets directory. The directory must be
// located under the assets object's root directory.
func (assets *Assets) AddDir(dir string) error {
	if !strings.HasPrefix(dir, assets.root) {
		return fmt.Errorf("directory '%s' is not subdirectory of '%s'",
			dir, assets.root)
	}
	f, err := os.Open(dir)
	if err != nil {
		return err
	}
	defer f.Close()

	files, err := f.ReadDir(0)
	if err != nil {
		return err
	}
	for _, file := range files {
		name := path.Join(dir, file.Name())
		if file.IsDir() {
			err = assets.AddDir(name)
			if err != nil {
				return err
			}
		} else {
			assets.files[name] = file
		}
	}

	return nil
}

// Copy copies the assets to the argument directory.
func (assets *Assets) Copy(dir string) error {
	dir = path.Clean(dir)
	for asset, assetEntry := range assets.files {

		assetInfo, err := assetEntry.Info()
		if err != nil {
			return err
		}
		output := path.Join(dir, asset[len(assets.root):])

		if isValid(assetInfo, output) {
			continue
		}

		err = os.MkdirAll(path.Dir(output), 0777)
		if err != nil {
			return err
		}

		src, err := os.Open(asset)
		if err != nil {
			return err
		}
		dst, err := os.OpenFile(output, os.O_RDWR|os.O_CREATE|os.O_TRUNC,
			assetInfo.Mode())
		if err != nil {
			src.Close()
			return err
		}
		_, err = io.Copy(dst, src)
		src.Close()
		dst.Close()
		if err != nil {
			return err
		}

		log.Printf("%s\t=> %s\n", asset, output)
	}
	return nil
}
