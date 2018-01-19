package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"github.com/fsouza/go-dockerclient"
	"gopkg.in/yaml.v2"
	"os"
)

const (
	helpMsg = `docker_image_info - Docker Images Info tool

Usage: docker_image_info [options]

Options:
`
)

// The printHelp() function will use the default helpMsg
func printHelp() {
	fmt.Println(helpMsg)
	flag.PrintDefaults()
}

// We print the images cycling the slice containint APIImages structs
//type APIImages struct {
//ID          string            `json:"Id" yaml:"Id" toml:"Id"`
//RepoTags    []string          `json:"RepoTags,omitempty" yaml:"RepoTags,omitempty" toml:"RepoTags,omitempty"`
//Created     int64             `json:"Created,omitempty" yaml:"Created,omitempty" toml:"Created,omitempty"`
//Size        int64             `json:"Size,omitempty" yaml:"Size,omitempty" toml:"Size,omitempty"`
//VirtualSize int64             `json:"VirtualSize,omitempty" yaml:"VirtualSize,omitempty" toml:"VirtualSize,omitempty"`
//ParentID    string            `json:"ParentId,omitempty" yaml:"ParentId,omitempty" toml:"ParentId,omitempty"`
//RepoDigests []string          `json:"RepoDigests,omitempty" yaml:"RepoDigests,omitempty" toml:"RepoDigests,omitempty"`
//Labels      map[string]string `json:"Labels,omitempty" yaml:"Labels,omitempty" toml:"Labels,omitempty"`

func jsonEnc(imgs []docker.APIImages) (string, error) {
	jsonSer, err := json.Marshal(imgs)
	if err != nil {
		return "", err
	}
	result := fmt.Sprintf("%s", string(jsonSer))
	return result, nil
}

func yamlEnc(imgs []docker.APIImages) (string, error) {
	yamlser, err := yaml.Marshal(imgs)
	if err != nil {
		return "", err
	}
	result := fmt.Sprintf("%s", string(yamlser))
	return result, nil
}

func plainTextAll(imgs []docker.APIImages) error {
	if len(imgs) == 0 {
		return errors.New("Empty image list")
	}
	for _, img := range imgs {
		fmt.Println("ID: ", img.ID)
		fmt.Println("RepoTags: ", img.RepoTags)
		fmt.Println("Created: ", img.Created)
		fmt.Println("Size: ", img.Size)
		fmt.Println("VirtualSize: ", img.VirtualSize)
		fmt.Println("ParentID: ", img.ParentID)
		fmt.Println("RepoDigests: ", img.RepoDigests)
		fmt.Println("Labels: ", img.Labels)
		fmt.Printf("\n")
	}
	return nil
}

func checkMultiEncFlag(

	jsonON := flag.Bool("json", false, "Use JSON encoding")
	yamlON := flag.Bool("yaml", false, "Use YAML encondig")
	textON := flag.Bool("text", false, "Output the result in plain text")
	helpON := flag.Bool("help", false, "Print a more detailed help")

	flag.Parse()

	// Use an extended help function for a user friendly output
	if *helpON {
		printHelp()
		os.Exit(0)
	}

	// Let's create a client instance connecting to the local socket
	endpoint := "unix:///var/run/docker.sock"
	client, err := docker.NewClient(endpoint)
	if err != nil {
		panic(err)
	}

	// The ListImages method returns a slice of APIImages structs
	imgs, err := client.ListImages(docker.ListImagesOptions{All: false})
	if err != nil {
		panic(err)
	}

	if *jsonON && !*yamlON && !*textON {
		// Print the conent in json format
		res, err := jsonEnc(imgs)
		if err != nil {
			fmt.Printf("Json encoding error: %s\n", err)
			os.Exit(1)
		}
		fmt.Printf("%s", res)
		os.Exit(0)
	}

	if !*jsonON && *yamlON && !*textON {
		// Print the conent in json format
		res, err := yamlEnc(imgs)
		if err != nil {
			fmt.Println("Yaml encoding error: %s\n", err)
			os.Exit(1)
		}
		fmt.Printf("%s", res)
		os.Exit(0)
	}

	// We want text to be the default output format
	if (!*jsonON && !*yamlON && *textON) || (!*jsonON && !*yamlON && !*textON) {
		err := plainTextAll(imgs)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		os.Exit(0)
	}

	fmt.Println("Error: The program does not support more than one encoding flag")
	os.Exit(1)
}
