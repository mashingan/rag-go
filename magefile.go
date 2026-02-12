//go:build mage
// +build mage

package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"

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
