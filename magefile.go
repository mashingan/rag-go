//go:build mage
// +build mage

package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"sync"

	"github.com/magefile/mage/mg" // mg contains helpful utility functions, like Deps
)

// Default target to run when none is specified
// If not set, running mage will list available targets
// var Default = Build

func Init() error {
	mg.Deps(InstallDeps)
	/*
		cmd := exec.Command("go", "version")
		outb, err := cmd.Output()
		if err != nil {
			log.Println(err)
			return err
		}

		var major, minor, patch int
		n, err := fmt.Sscanf(string(outb), "go version go%d.%d.%d", &major, &minor, &patch)
		if err != nil {
			log.Println(err)
			return err
		}
		if n != 3 {
			msg := fmt.Sprintf("only scanning %d items, expected 3", n)
			log.Println(msg)
			return fmt.Errorf(msg)
		}
	*/
	cmd := exec.Command("w2v", "train", "--train",
		"./deps/word2vec/doc/leo-tolstoy-war-and-peace-en.txt",
		"--output", "vectors.bin")
	return cmd.Run()
}

// A build step that requires additional params, or platform specific steps for example
func Build() error {
	cmd := exec.Command("go", "build", "-o", "MyApp", ".")
	return cmd.Run()
}

// A custom install step if you need your bin someplace other than go/bin
func Install() error {
	mg.Deps(Build)
	fmt.Println("Installing...")
	return os.Rename("./MyApp", "/usr/bin/MyApp")
}

// Manage your deps, or running package managers.
func InstallDeps() error {
	fmt.Println("Installing Deps...")
	cmd := exec.Command("git", "submodule", "update", "--init")
	return cmd.Run()
}

func InstallInit() error {
	fmt.Println("Init submodule...")
	cmd := exec.Command("git", "submodule", "init")
	return cmd.Run()
}

// Download model from hugging face
func DownloadModel() error {
	const (
		targetsPath    = "./deps/all-MiniLM-L6-v2"
		modelfname     = "model.safetensors"
		modelUrl       = "https://huggingface.co/sentence-transformers/all-MiniLM-L6-v2/resolve/main/model.safetensors?download=true"
		vocabfname     = "vocab.txt"
		vocabUrl       = "https://huggingface.co/sentence-transformers/all-MiniLM-L6-v2/resolve/main/vocab.txt?download=true"
		tokenizerfname = "tokenizer.json"
		tokenizerUrl   = "https://huggingface.co/sentence-transformers/all-MiniLM-L6-v2/resolve/main/tokenizer.json?download=true"
	)
	fileurl := map[string]string{
		modelfname:     modelUrl,
		vocabfname:     vocabUrl,
		tokenizerfname: tokenizerUrl,
	}
	for fname, _ := range fileurl {
		if f, err := os.Open(filepath.Join(targetsPath, fname)); err != nil && !os.IsNotExist(err) {
			continue
		} else {
			f.Close()
		}
		delete(fileurl, fname)
	}
	if len(fileurl) == 0 {
		fmt.Println("All model files already downloaded.")
		return nil
	}
	fmt.Println("This will download around 90 MB model...")
	if err := os.MkdirAll(targetsPath, os.ModeDir|0665); err != nil {
		return err
	}
	var (
		wg   sync.WaitGroup
		errc = make(chan error, 3)
	)

	for fname, url := range fileurl {
		fname := fname
		url := url
		wg.Go(func() {
			fmt.Printf("Downloading %s...\n", fname)
			resp, err := http.Get(url)
			if err != nil {
				errc <- err
				return
			}
			if resp == nil {
				errc <- fmt.Errorf("response is nil")
				return
			}
			if resp.StatusCode != http.StatusOK {
				errc <- fmt.Errorf("download file %s not http status ok: got %s",
					fname, resp.Status)
				return
			}
			defer resp.Body.Close()
			f, err := os.Create(filepath.Join(targetsPath, fname))
			if err != nil {
				errc <- err
				return
			}
			io.Copy(f, resp.Body)
			errc <- nil
			fmt.Printf("Done downloading %s to %s\n", fname, targetsPath)
		})
	}
	wg.Wait()
	close(errc)
	var errs error
	for err := range errc {
		errs = errors.Join(errs, err)
	}
	return errs
}

// Clean up after yourself
func Clean() {
	fmt.Println("Cleaning...")
	os.RemoveAll("MyApp")
}

func Example_w2v() error {
	cmd := exec.Command("go", "run", "scratch/w2v-pure-go/main.go")
	outp, err := cmd.StdoutPipe()
	if err != nil {
		log.Println(err)
		return err
	}
	if err := cmd.Start(); err != nil {
		log.Println(err)
		return err
	}
	go io.Copy(os.Stdout, outp)
	cmd.Wait()
	return nil
}
