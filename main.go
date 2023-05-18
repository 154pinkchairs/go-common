package main

import (
	"bufio"
	"fmt"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"strings"

	. "github.com/154pinkchairs/go-common/validators"
)

func main() {
	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Println("\nSelect an option:")
		fmt.Println("1. Process a tiny string (32B), no goroutines")
		fmt.Println("2. Process a small string (1KB), no goroutines")
		fmt.Println("3. Process a mid string (1MB), no goroutines")
		fmt.Println("4. Process a large string (4MB), no goroutines")
		fmt.Println("5. Process a tiny string (32B), with goroutines")
		fmt.Println("6. Process a small string (1KB), with goroutines")
		fmt.Println("7. Process a mid string (1MB), with goroutines")
		fmt.Println("8. Process a large string (4MB), with goroutines")
		fmt.Println("q - Quit")
		fmt.Println("d - Exit the prompt loop but keep the program running")

		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Printf("Error reading input: %s", err.Error())
			continue
		}
		input = strings.TrimSpace(input)

		switch input {
		case "1":
			tiny := strings.Repeat("a", 1<<5) // ~32B string. One UTF-8 char is 1-4 bytes.
			EnsureUTF8(tiny)
		case "2":
			small := strings.Repeat("a", 1<<10) // ~1KB strings
			EnsureUTF8(small)
		case "3":
			mid := strings.Repeat("a", 1<<20) // ~1MB strings
			EnsureUTF8(mid)
		case "4":
			mid := strings.Repeat("a", 1<<20)  // ~1MB strings
			large := strings.Repeat(mid, 1<<2) // ~4MB strings
			EnsureUTF8(large)
		case "5":
			tiny := strings.Repeat("a", 1<<5)
			GEnsureUTF8(tiny)
		case "6":
			small := strings.Repeat("a", 1<<10)
			GEnsureUTF8(small)
		case "7":
			mid := strings.Repeat("a", 1<<20)
			GEnsureUTF8(mid)
		case "8":
			mid := strings.Repeat("a", 1<<20)
			large := strings.Repeat(mid, 1<<2)
			GEnsureUTF8(large)
		case "q":
			os.Exit(0)
		case "d":
			return
		default:
			fmt.Println("Invalid option")
		}
	}
}
