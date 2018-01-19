package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"github.com/fsouza/go-dockerclient"
	"gopkg.in/yaml.v2"
	"os"
	"regexp"
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

// We print the images cycling the slice []docker.APIImages containing APIImages structs
//type APIImages struct {
//ID          string            `json:"Id" yaml:"Id" toml:"Id"`
//RepoTags    []string          `json:"RepoTags,omitempty" yaml:"RepoTags,omitempty" toml:"RepoTags,omitempty"`
//Created     int64             `json:"Created,omitempty" yaml:"Created,omitempty" toml:"Created,omitempty"`
//Size        int64             `json:"Size,omitempty" yaml:"Size,omitempty" toml:"Size,omitempty"`
//VirtualSize int64             `json:"VirtualSize,omitempty" yaml:"VirtualSize,omitempty" toml:"VirtualSize,omitempty"`
//ParentID    string            `json:"ParentId,omitempty" yaml:"ParentId,omitempty" toml:"ParentId,omitempty"`
//RepoDigests []string          `json:"RepoDigests,omitempty" yaml:"RepoDigests,omitempty" toml:"RepoDigests,omitempty"`
//Labels      map[string]string `json:"Labels,omitempty" yaml:"Labels,omitempty" toml:"Labels,omitempty"`

// Function jsonEnc() translates []docker.APIImages slice to json string
func jsonEnc(imgs []docker.APIImages) (string, error) {
	jsonSer, err := json.Marshal(imgs)
	if err != nil {
		return "", err
	}
	result := fmt.Sprintf("%s", string(jsonSer))
	return result, nil
}

// Function yamlEnc() translates []docker.APIImages slice to yaml string
func yamlEnc(imgs []docker.APIImages) (string, error) {
	yamlser, err := yaml.Marshal(imgs)
	if err != nil {
		return "", err
	}
	result := fmt.Sprintf("%s", string(yamlser))
	return result, nil
}

// Function jsonAge() translates map[string]int64 map to json string
func jsonAge(m map[string]int64) (string, error) {
	jsonAgeMarshal, err := json.Marshal(m)
	if err != nil {
		return "", err
	}
	result := fmt.Sprintf("%s", string(jsonAgeMarshal))
	return result, nil
}

// Function yamlAge() translates map[string]int64 map to yaml string
func yamlAge(m map[string]int64) (string, error) {
	yamlAgeMarshal, err := yaml.Marshal(m)
	if err != nil {
		return "", err
	}
	result := fmt.Sprintf("%s", string(yamlAgeMarshal))
	return result, nil
}

//TODO: Convert to a map[string]string output and pass it
// to the calling function who will print it by itself.
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

// Function imageTimeStamp() can be useful to collect the age of old images on
// node caches.
func imageTimeStamp(imgs []docker.APIImages) map[string]int64 {
	imgTSMap := make(map[string]int64)
	// We don't want the hash algorithm info here
	pattern, _ := regexp.Compile("sha256:")
	for _, img := range imgs {
		cleanID := pattern.ReplaceAllString(img.ID, "")
		imgTSMap[cleanID] = img.Created
	}
	return imgTSMap
}

// Function checkMultiEncFlag() checks if the user uses more than one encoding flag
func checkMultiEncFlag(jsonFlag bool, yamlFlag bool, textFlag bool) bool {
	if (jsonFlag && !yamlFlag && !textFlag) ||
		(!jsonFlag && yamlFlag && !textFlag) ||
		(!jsonFlag && !yamlFlag && textFlag) ||
		(!jsonFlag && !yamlFlag && !textFlag) {
		return false
	} else {
		return true
	}
}

func main() {

	jsonON := flag.Bool("json", false, "Use JSON encoding")
	yamlON := flag.Bool("yaml", false, "Use YAML encondig")
	textON := flag.Bool("text", false, "Output the result in plain text")
	helpON := flag.Bool("help", false, "Print a more detailed help")
	imgAgeON := flag.Bool("age", false, "Print age of images in Unix epoch (text/json/yaml format)")

	flag.Parse()

	// Use an extended help function for a user friendly output
	if *helpON {
		printHelp()
		os.Exit(0)
	}

	if checkMultiEncFlag(*jsonON, *yamlON, *textON) {
		fmt.Println("Error: The program does not support more than one encoding flag")
		os.Exit(1)
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

	if *jsonON && !*imgAgeON {
		// Print the conent in json format
		res, err := jsonEnc(imgs)
		if err != nil {
			fmt.Printf("Json encoding error: %s\n", err)
			os.Exit(1)
		}
		fmt.Printf("%s", res)
		os.Exit(0)
	}

	if *yamlON && !*imgAgeON {
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
	if *textON && !*imgAgeON || (!*jsonON && !*yamlON && !*textON && !*imgAgeON) {
		err := plainTextAll(imgs)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		os.Exit(0)
	}

	if *imgAgeON {
		imgAgeMap := imageTimeStamp(imgs)
		if *textON || (!*jsonON && !*yamlON && !*textON) {
			fmt.Printf("IMAGE ID\t\t\t\t\t\t\t  AGE\n")
			for id, epoch := range imgAgeMap {
				fmt.Printf("%s  %d\n", id, epoch)
			}
		} else if *yamlON {
			res, err := yamlAge(imgAgeMap)
			if err != nil {
				fmt.Println("Yaml encoding error: %s\n", err)
				os.Exit(1)
			}
			fmt.Printf("%s", res)
		} else if *jsonON {
			res, err := jsonAge(imgAgeMap)
			if err != nil {
				fmt.Println("Json encoding error: %s\n", err)
				os.Exit(1)
			}
			fmt.Printf("%s", res)
		}
	}
}
