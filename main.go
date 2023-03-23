package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"unicode"

	"github.com/alecthomas/chroma/quick"
	"github.com/atotto/clipboard"
)

const interfaceFileUrl = "https://raw.githubusercontent.com/SunWeb3Sec/DeFiHackLabs/main/src/test/interface.sol"

func getInterfaceFilePath() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		fmt.Println(err)
	}
	interfaceFilePath := dir + "/interface.sol"
	//fmt.Println("interfaceFilePath:", interfaceFilePath)
	return interfaceFilePath
}

func download() {
	interfaceFilePath := getInterfaceFilePath()

	err := os.MkdirAll(filepath.Dir(interfaceFilePath), os.ModePerm)
	if err != nil {
		fmt.Printf("Error creating directory: %s\n", err)
		os.Exit(1)
	}
	out, err := os.Create(interfaceFilePath)
	if err != nil {
		fmt.Printf("Error creating file: %s\n", err)
		os.Exit(1)
	}
	defer out.Close()
	resp, err := http.Get(interfaceFileUrl)
	if err != nil {
		fmt.Printf("Error downloading file: %s\n", err)
		os.Exit(1)
	}
	defer resp.Body.Close()
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		fmt.Printf("Error saving file: %s\n", err)
		os.Exit(1)
	}
	fmt.Println("contract file downloaded successfully")
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: ./interface.go <interface_name>")
		os.Exit(1)
	}

	interfaceFilePath := getInterfaceFilePath()

	// if contractFile not exist ,download it
	if _, err := os.Stat(interfaceFilePath); os.IsNotExist(err) {
		fmt.Println("contract file not exist,downloading...")
		download()
	}

	interfaceName := os.Args[1]
	rawData, err := ioutil.ReadFile(interfaceFilePath)
	if err != nil {
		fmt.Printf("Error reading file: %s\n", err)
		os.Exit(1)
	}

	if !strings.HasPrefix(interfaceName, "I") {
		interfaceName = "I" + interfaceName
	}

	re := regexp.MustCompile(`(?ms)(?i)interface ` + strings.Title(interfaceName) + ` \{[^}]+\}`)
	interfaceCode := re.FindString(string(rawData))

	if interfaceCode == "" {
		// fmt.Printf("Interface not found: %s\n", interfaceName)
		// 查找类似的接口
		pattern := "interface\\s+"
		for _, char := range interfaceName {
			pattern += fmt.Sprintf("[%c%c].*", unicode.ToUpper(char), unicode.ToLower(char))
		}
		re = regexp.MustCompile(pattern)
		//log.Println("re", re)
		interfaces := re.FindAllString(string(rawData), -1)
		if len(interfaces) == 0 {
			fmt.Printf("\033[31mInterface not found: %s\033[0m\n", interfaceName)
		} else {
			fmt.Printf("\033[31mInterface not found: %s\033[0m\n", interfaceName)
			fmt.Println("Did you mean:")
			fmt.Println("--------------------------")
			for _, i := range interfaces {

				fmt.Printf("\033[32m%s\033[0m\n", strings.TrimPrefix(strings.Split(i, "{")[0], "interface "))
				//quick.Highlight(os.Stdout, i, "solidity", "terminal256", "dracula")
			}
		}
	} else {
		// style 参考 https://pkg.go.dev/github.com/alecthomas/chroma@v0.10.0/styles
		quick.Highlight(os.Stdout, interfaceCode, "solidity", "terminal256", "dracula")

		err = clipboard.WriteAll(interfaceCode)
		if err != nil {
			fmt.Println("\n\nError copying to clipboard:", err)
		} else {
			fmt.Println("\n\nInterface code copied to clipboard.")
		}
	}
}
