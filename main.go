package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"golang.org/x/term"
)

const windowsUrl = "https://raw.githubusercontent.com/hashimthearab/go-remote-download/master/assets/gophertunnel.exe"
const darwinUrl = "https://raw.githubusercontent.com/hashimthearab/go-remote-download/master/assets/gophertunnel_darwin"

// target is the path to install the file to. Leave it blank to use the current directory and file name.
var target = ""

func main() {
	fmt.Println("Golang Remote Downloader")

	if target == "" {
		parsedUrl, err := url.Parse(getUrl())
		if err != nil {
			log.Fatal(err)
		}
		parsedUrl.RawQuery = ""
		target = filepath.Base(parsedUrl.Path)
	}

	for {
		fmt.Printf("[d] Download %s\n[r] Run %s\n[e] Exit\nExecute: ", target, target)
		char, err := readSingleChar()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println() // new line

		switch strings.ToLower(string(char)) {
		case "d":
			downloadFile()
		case "r":
			runFile()
		case "e":
			return
		default:
			fmt.Println("Invalid input, try again.")
		}
	}
}

func readSingleChar() (byte, error) {
	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		return 0, err
	}
	defer term.Restore(int(os.Stdin.Fd()), oldState)

	var buf [1]byte
	_, err = os.Stdin.Read(buf[:])
	if err != nil {
		return 0, err
	}
	return buf[0], nil
}

func downloadFile() {
	fmt.Printf("Downloading %s...\n", target)

	resp, err := http.Get(getUrl())
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	file, err := os.Create(target)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	_, err = io.Copy(file, resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Downloaded %s successfully.\n", target)
}

func runFile() {
	fmt.Printf("Running %s...\n", target)

	if _, err := os.Stat(target); os.IsNotExist(err) {
		fmt.Printf("%s does not exist, download it first.\n", target)
		return
	}

	var cmd *exec.Cmd
	if runtime.GOOS == "linux" || runtime.GOOS == "darwin" {
		err := os.Chmod(target, 0755)
		if err != nil {
			log.Fatal(err)
		}
	}
	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("cmd", "/C", "start", target)
	case "darwin":
		cmd = exec.Command("open", target)
	case "linux":
		cmd = exec.Command("xdg-open", target)
	default:
		fmt.Println("Unsupported OS.")
		return
	}
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Start()
	if err != nil {
		log.Fatal(err)
	}

	err = cmd.Wait()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("File execution completed.")
}

func getUrl() string {
	var target = windowsUrl
	if runtime.GOOS == "darwin" {
		target = darwinUrl
	}

	if strings.Contains(target, "raw.githubusercontent.com") {
		// Bypass github cache
		target += "?cb=" + fmt.Sprintf("%d", time.Now().Unix())
	}

	return target
}
