package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

type VideoConfig struct {
	Title       string `yaml:"title"`
	Description string `yaml:"description"`
	TalentPath  string `yaml:"talentPath"`
	ProjectPath string `yaml:"projectPath"`
	Items       yaml.Node
}

func main() {
	pwd, err := os.Getwd()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	var dir, yamlName, zipFileName string
	flag.StringVar(&dir, "d", "", "relative path directory to compress") // defaults to working dir
	flag.StringVar(&yamlName, "y", "", "yaml file name for file conversion")
	flag.StringVar(&zipFileName, "zname", "", "zip file name for file compression")
	flag.Parse()

	// video-zip -d folder1/ -y main
	var yamlPath, zipPath string
	if yamlName == "" {
		fmt.Fprintln(os.Stderr, "Must specify yaml configuraiton file name")
		fmt.Fprintln(os.Stderr, "")
		flag.Usage()
		os.Exit(1)
	}
	if zipFileName == "" {
		fmt.Fprintln(os.Stderr, "Must specify zip file name")
		fmt.Fprintln(os.Stderr, "")
		flag.Usage()
		os.Exit(1)
	}
	pwd = normalize(pwd)

	if dir == "" {
		yamlPath = pwd + yamlName
		zipPath = pwd
	} else {
		dir = normalize(dir)
		yamlPath = pwd + dir + yamlName
		zipPath = pwd + dir

	}
	fmt.Println("flag dir: ", dir)
	fmt.Println("flag yamlName: ", yamlName)
	fmt.Println("yamlPath: ", yamlPath)
	fmt.Println("zipPath: ", zipPath)

	var videoconfig VideoConfig
	getConfig(yamlPath, &videoconfig)
	fmt.Printf("VideoConfig:\nTitle: %s\nDescription: %s\nTalentPath: %s\nProjectPath: %s\nItems: %v\n", videoconfig.Title, videoconfig.Description, videoconfig.TalentPath, videoconfig.ProjectPath, videoconfig.Items)
	convertToWFS("", dir, &videoconfig)
	fmt.Printf("after: %s\n", videoconfig.ProjectPath)
	err = saveConfig(yamlPath, videoconfig)
	if err != nil {
		log.Fatalln(err.Error())
		os.Exit(1)
	}
	cFiles, err := ioutil.ReadDir(zipPath)
	if err != nil {
		log.Fatalln(err.Error())
		os.Exit(1)
	}
	fmt.Println("Files to be zipped...")
	for _, file := range cFiles {
		fmt.Println(file.Name())
	}
	// for testing
	err = ZipFiles(zipFileName, cFiles)
	if err != nil {
		log.Fatalln(err.Error())
		os.Exit(1)
	}

}

func getConfig(file string, data *VideoConfig) {
	yfile, err := ioutil.ReadFile(file)
	if err != nil {
		log.Fatalln("Error: ", err.Error())
	}
	derr := yaml.Unmarshal(yfile, &data)
	if derr != nil {
		log.Fatalln("Failed: ", derr.Error())
		os.Exit(1)
	}
}

func saveConfig(path string, config VideoConfig) error {
	data, err := yaml.Marshal(config)
	if err != nil {
		return err
	}
	// write out file
	ioutil.WriteFile(path, data, 0777)
	return nil
}

func normalize(dir string) string {
	if dir[:len(dir)-1] != "/" {
		return dir + "/"
	}
	return dir
}

func convertToWFS(wDir, dir string, config *VideoConfig) {
	dir = strings.Replace(dir, "/", "\\", 1)
	if wDir == "" {
		config.ProjectPath = "F:\\videos\\" + dir

	} else {
		config.ProjectPath = wDir + dir

	}
}
