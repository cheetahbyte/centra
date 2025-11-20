package helper

import (
	"fmt"
	"os"

	"github.com/cheetahbyte/centra/internal/config"
)

func PrettyKey(pubKey string) {
	if config.GetPublicSSHKey() == "" {
		fmt.Println("------------------------------------------------")
		fmt.Println("Add this Deploy Key to your GitHub Repo:")
		b, _ := os.ReadFile(pubKey)
		fmt.Println(string(b))
		fmt.Println("------------------------------------------------")
	}
}
