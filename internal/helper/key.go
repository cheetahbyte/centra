package helper

import (
	"fmt"
	"os"
)

func PrettyKey(pubKey string) {
	fmt.Println("------------------------------------------------")
	fmt.Println("Add this Deploy Key to your GitHub Repo:")
	b, _ := os.ReadFile(pubKey)
	fmt.Print(string(b))
	fmt.Println("------------------------------------------------")
}
