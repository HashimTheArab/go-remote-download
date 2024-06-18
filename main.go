package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"golang.org/x/term"
)

const url = "https://raw.githubusercontent.com/hashimthearab/go-remote-download/master/assets/gophertunnel.exe"

// target is the path to install the file to.
const target = "gophertunnel.exe"

func main() {
	fmt.Println("Golang Remote Downloader")
	for {
		fmt.Printf("Actions:\n  [d] Download %s\n  [r] Run %s\n  [u] Update\nRun: ", target, target)
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

	resp, err := http.Get(url)
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