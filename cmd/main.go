package main

import (
	"flag"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"

	"github.com/yed-palettes/svg-to-graphml/converter"
)

type Cfg struct {
	inPath      string
	outPath     string
	noRecursive bool
}

func (c Cfg) validate() error {
	if c.inPath == "" {
		return fmt.Errorf("you must provide a path\n-path <file or directory>")
	}
	return nil
}

func loadNodesFromPath(loadSVGPath string, noRecurse bool) (converter.NodeMap, error) {
	info, err := os.Stat(loadSVGPath)
	if err != nil {
		return nil, err
	}

	nodesMap := converter.NodeMap{}

	// Single file
	if !info.IsDir() {
		node, err := converter.NewNode(loadSVGPath)
		if err != nil {
			return nil, err
		}
		parentDir := filepath.Base(filepath.Dir(loadSVGPath))
		nodesMap[parentDir] = converter.Nodes{node}
		return nodesMap, nil
	}

	// Directory
	filepath.WalkDir(loadSVGPath, func(path string, d fs.DirEntry, err error) error {
		if d.IsDir() {
			if noRecurse && path != loadSVGPath {
				return fs.SkipDir
			}
			return nil
		}

		if filepath.Ext(path) != ".svg" {
			log.Printf("Skipping %s | ext is: %s", path, filepath.Ext(path))
			return nil
		}

		parentDir := filepath.Base(filepath.Dir(path))
		node, err := converter.NewNode(path)
		if err != nil {
			log.Println("error creating node:", err, path)
			return nil
		}
		nodesMap[parentDir] = append(nodesMap[parentDir], node)
		return nil
	})

	return nodesMap, nil
}

func main() {
	inPath := flag.String("path", "", "path to search for svgs & folders of svgs")
	outDir := flag.String("out", "./out/", "output location")
	noRecursive := flag.Bool("no-recurse", false, "disable searching recursively")

	flag.Parse()

	cfg := Cfg{
		inPath:      *inPath,
		outPath:     *outDir,
		noRecursive: *noRecursive,
	}
	if err := cfg.validate(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	nodesMap, err := loadNodesFromPath(cfg.inPath, cfg.noRecursive)
	if err != nil {
		log.Fatalf("failed to load nodes: %v", err)
	}

	// Ensure output directory exists
	if err := os.MkdirAll(cfg.outPath, 0755); err != nil {
		log.Fatalf("failed to create output directory: %v", err)
	}

	for outfile, nodes := range nodesMap {
		// Anonymous function so we can close files right after we're done with them  (rather than end of main)
		err := func() error {
			dstFile, err := os.Create(filepath.Join(cfg.outPath, outfile+".graphml"))
			if err != nil {
				return err
			}
			defer dstFile.Close()

			err = nodes.RenderTemplate(dstFile)
			if err != nil {
				return err
			}
			return nil
		}()

		if err != nil {
			log.Printf("error processing %s: %v", outfile, err)
		}
	}

	fmt.Println("completed, results stored:", cfg.outPath)
}
