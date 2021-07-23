package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"regexp"
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
	var dir, yamlName, zipFileName, wDir, yamlPath, zipPath, mainName string
	var debug, final bool
	flag.StringVar(&dir, "d", "", "relative path directory to compress") // defaults to working dir
	flag.StringVar(&yamlName, "y", "", "yaml file name for file conversion (using .yml only)")
	flag.StringVar(&zipFileName, "zname", "", "zip file name for file compression (defaults to -y param)")
	flag.StringVar(&wDir, "wdir", "", "windows file directory path for path replacement")
	flag.BoolVar(&debug, "debug", false, "print important variables")
	flag.BoolVar(&final, "f", false, "final configuraiton version")
	flag.Parse()

	if !checkYaml(yamlName) {
		fmt.Fprintln(os.Stderr, "Must specify yaml configuraiton file")
		fmt.Fprintln(os.Stderr, "")
		flag.Usage()
		os.Exit(1)
	}
	mainName = yamlName[:len(yamlName)-4]
	if zipFileName == "" {
		zipFileName = mainName + ".zip"
	} else if !checkZip(zipFileName) {
		zipFileName += ".zip"
	}
	pwd = normalize(pwd)
	//
	if dir == "" {
		yamlPath = pwd + yamlName
		dir = yamlName[:len(yamlName)-4]
		zipPath = "./"
	} else {
		dir = normalize(dir)
		yamlPath = pwd + dir + yamlName
		zipPath = dir
	}
	// creating a copy of yaml file for file safety
	copyConfig(yamlName, yamlPath)
	var videoconfig VideoConfig
	getConfig(yamlPath, &videoconfig)
	if debug {
		fmt.Println("dir=", dir)
		fmt.Println("yamlName= ", yamlName)
		fmt.Println("yamlPath=", yamlPath)
		fmt.Println("zipPath=", zipPath)
		fmt.Println("final=", final)
		fmt.Println("debug=", debug)
		fmt.Printf("VideoConfig:\n\tTitle: %s\n\tDescription: %s\n\tTalentPath: %s\n\tProjectPath: %s\n\n", videoconfig.Title, videoconfig.Description, videoconfig.TalentPath, videoconfig.ProjectPath)
	}
	convertToWFS(wDir, dir, final, &videoconfig)
	if debug {
		fmt.Printf("VideoConfig after convert:\n\tTitle: %s\n\tDescriptio n: %s\n\tTalentPath: %s\n\tProjectPath: %s\n\n", videoconfig.Title, videoconfig.Description, videoconfig.TalentPath, videoconfig.ProjectPath)
	}

	if err := saveConfig(yamlPath, videoconfig); err != nil {
		log.Fatalln(err.Error())
		os.Exit(1)
	}
	cFiles, err := ioutil.ReadDir(zipPath)
	if err != nil {
		log.Fatalln(err.Error())
		os.Exit(1)
	}
	if debug {
		fmt.Printf("Compressing files into %s\n", zipFileName)
		for _, file := range cFiles {
			fmt.Printf("\t- %s\n", file.Name())
		}
	}

	if err := ZipFiles(zipFileName, zipPath, cFiles); err != nil {
		log.Fatalln(err.Error())
		os.Exit(1)
	}
	fmt.Printf("Success: %s created!\n", zipFileName)

}

// Reads the yaml configuration file and assigns file data to type
// VideoConfig variable
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

// Writes the yaml configuration file
func saveConfig(path string, config VideoConfig) error {
	data, err := yaml.Marshal(config)
	if err != nil {
		return err
	}
	// write out file
	ioutil.WriteFile(path, data, 0777)
	return nil
}

// Checks if the string provided contains a "/" in the last string index
// creates a new string with a trailing "/" if one is not trailing
func normalize(dir string) string {
	if dir[len(dir)-1] != '/' {
		return dir + "/"
	}
	return dir
}

// Converts Unix file notation to Windows
func convertToWFS(wDir, dir string, isFinal bool, config *VideoConfig) {
	dir = strings.Replace(dir, "/", "\\", 1)
	if wDir == "" {
		// default for personal use case
		config.ProjectPath = "F:\\videos\\" + dir + "\\"

	} else {
		config.ProjectPath = wDir + dir + "\\"
	}
	if isFinal {
		config.TalentPath = "final.webm"
	}

}

// Validates YAML file extension
func checkYaml(yFile string) bool {
	return regexp.MustCompile(`^.*\.(yml|YML)$`).MatchString(yFile)
}

// Validates Zip file extension
func checkZip(zFile string) bool {
	return regexp.MustCompile(`^.*\.(zip|ZIP)$`).MatchString(zFile)
}

// Creates a local copy of the configuration file
func copyConfig(name, yPath string) {
	//Read all the contents of the  original file
	bytesRead, err := ioutil.ReadFile(yPath)
	if err != nil {
		log.Fatal(err)
	}

	//Copy all the contents to the desitination file
	err = ioutil.WriteFile("copy-"+name, bytesRead, 0755)
	if err != nil {
		log.Fatal(err)
	}

}
