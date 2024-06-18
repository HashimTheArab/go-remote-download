package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"

	"golang.org/x/term"
)

const url = ""
const name = "gophertunnel.exe"

func main() {
	fmt.Println("Golang Remote Downloader")
	for {
		fmt.Print("Enter 'd' to download the latest version or 'r' to run it: ")
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
	fmt.Printf("Downloading %s...\n", name)

	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	file, err := os.Create(name)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	_, err = io.Copy(file, resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Downloaded %s successfully.\n", name)
}

func runFile() {
	fmt.Printf("Running %s...\n", name)

	cmd := exec.Command("cmd", "/C", "start", name)
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
