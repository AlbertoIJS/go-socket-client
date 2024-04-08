package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
	"os/exec"
	"strings"
)

func main() {
	// Connect to the server
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		fmt.Println("Error ", err)
		return
	}
	defer conn.Close()

	location := "local"
	// Get input from console
	reader := bufio.NewReader(os.Stdin)

	// Read line each time the users presses enter
	for {
		line, _ := reader.ReadString('\n')
		line = strings.TrimSpace(line)
		arr := strings.Split(line, " ")
		command := arr[0]

		if location == "local" && command != "put" {
			menu(line)
		} else {
			// Send input to the server
			_, err = conn.Write([]byte(line))
			if err != nil {
				fmt.Println("Error sending input to the server: ", err)
				return
			}

			if command == "get" {
				file, err := os.Create(arr[1])
				if err != nil {
					fmt.Println("Error al guardar archivo: ", err)
					return
				}

				_, err = io.Copy(file, conn)
				file.Close()
			} else if command == "put" {
				file, err := os.Open(arr[1])
				if err != nil {
					fmt.Println("Error al abrir el archivo: ", err)
					return
				}

				_, err = io.Copy(conn, file)
				if err != nil {
					fmt.Println("Error al copiar el archivo: ", err)
					return
				}
				file.Close()
			} else {
				// Read response from the server
				buffer := make([]byte, 1024)
				n, err := conn.Read(buffer)
				if err != nil {
					fmt.Println("Error reading data from the server: ", err)
					return
				}

				fmt.Println(string(buffer[:n]))
			}
		}

		if line == "quit" {
			break
		}

		// Change between local and remote directories
		if strings.Contains(line, "cd ") {
			arr := strings.Split(line, " ")
			if arr[1] == "local" || arr[1] == "remote" {
				location = arr[1]
			}
		}
	}
}

func menu(line string) {
	arr := strings.Split(line, " ")
	command := arr[0]

	switch command {
	case "list":
		output, err := exec.Command("ls", "-a").Output()
		if err != nil {
			fmt.Println("Error listing directory: ", err)
			return
		}
		fmt.Println(string(output))
	case "mkdir":
		output, err := exec.Command("mkdir", arr[1]).Output()
		if err != nil {
			fmt.Println("Error creating directory: ", err)
			return
		}
		fmt.Println(string(output))
	case "rmdir":
		output, err := exec.Command("rm -rf", arr[2]).Output()
		if err != nil {
			fmt.Println("Error deleting directory: ", err)
			return
		}
		fmt.Println(string(output))
	case "get":
	}
}
