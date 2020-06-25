package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"regexp"
	"runtime"
	"strings"
)

type ModInfoItem struct {
	Folder      string
	Name        string
	Description string
	Author      string
	Debug       string
}

var File *os.File
var WorkDir string
var ModsDir = "M:\\SteamLibrary\\steamapps\\common\\Don't Starve Together\\mods\\"
var EditorPath = "C:\\Program Files\\Sublime Text 3\\sublime_text.exe"
var ModsInfo []ModInfoItem
var ModInfoFile = "\\modinfo.lua"

func init() {
	ModsInfo = make([]ModInfoItem, 0)
	var err error
	File, err = os.OpenFile("modInfoTest.log",
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println(err)
	}
	WorkDir, err = os.Getwd()
	if err != nil {
		log.Println(err)
	}
	WorkDir, err = os.Getwd()
	if err != nil {
		log.Println(err)
	}
	fmt.Println("WorkDir: " + WorkDir)
	//defer f.Close()

	//TODO: create folders - data, logs, downloads
}

func appendToFile(text string) {
	if _, err := File.WriteString(text + "\n"); err != nil {
		log.Println(err)
	}
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func getModInfo(dir string) ModInfoItem {
	modInfoItem := ModInfoItem{}
	filePath := dir + ModInfoFile
	if !(fileExists(filePath)) {
		return modInfoItem
	}
	dat, err := ioutil.ReadFile(filePath)
	check(err)
	modInfoItem.Folder = dir
	modInfoItem.Author = regexp.MustCompile(`author =[\r\n]?[\r\n]?\s*(?:"|\[\[)([.\P{M}\p{M}]*?)(?:"|\]\])`).FindString(string(dat))
	modInfoItem.Name = regexp.MustCompile(`name =[\r\n]?[\r\n]?\s*(?:"|\[\[)([.\P{M}\p{M}]*?)(?:"|\]\])`).FindString(string(dat))
	modInfoItem.Description = regexp.MustCompile(`description =[\r\n]?[\r\n]?\s*(?:"|\[\[)([.\P{M}\p{M}]*?)(?:"|\]\])`).FindString(string(dat))
	modInfoItem.Debug = regexp.MustCompile(`name =[.\P{M}\p{M}]*?author\s?=`).FindString(string(dat))
	//add labels list
	ModsInfo = append(ModsInfo, modInfoItem)
	return modInfoItem
}

func main() {
	files, err := ioutil.ReadDir(ModsDir)
	if err != nil {
		log.Fatal(err)
	}
	for _, f := range files {
		if f.IsDir() {
			modInfo := getModInfo(ModsDir + f.Name())
			appendToFile(f.Name() + "\n" + modInfo.Name + "\n" + modInfo.Description + "\n\n")
		}
	}

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		exit := scanner.Text()
		if exit == "q" {
			println("Bye...")
			break
		} else {
			if fileExists(ModsDir + exit + ModInfoFile) {
				editInfo(ModsDir + exit + ModInfoFile)
			} else {
				//println("Doest exist")
				for _, modInfoItem := range ModsInfo {
					if strings.Index(modInfoItem.Name, exit) != -1 {
						println("Find in:" + modInfoItem.Name)
						println("Folder:" + modInfoItem.Folder)
						println("")
					}
				}
			}
		}
	}
}

func editInfo(filePath string) {
	cmd := exec.Command(EditorPath, filePath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if runtime.GOOS == "windows" {
		//cmd = exec.Command("tasklist")
	}

	err := cmd.Run()
	if err != nil {
		log.Fatalf("cmd.Run() failed with %s\n", err)
	} else {
		log.Println("OK")
	}
}

//features
// translate
// backups before translate
