package utils

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"syscall"

	"github.com/google/go-github/v45/github"
)

func HookHandler(w http.ResponseWriter, r *http.Request) {
    payload, err := io.ReadAll(r.Body)
    if err != nil {
		log.Printf("CICD: git payload read error: %v", err)
        w.WriteHeader(http.StatusInternalServerError)
        return
    }
    defer r.Body.Close()
	log.Printf("CICD: payload: %v", payload)

	sig := r.Header.Get("X-Hub-Signature-256")
    if !validateSignature(sig, payload, parsed.WEBHOOK_SECRET) {
		log.Printf("CICD: invalid signature %v", sig)
        w.WriteHeader(http.StatusUnauthorized)
        return
    }

	hookType := github.WebHookType(r)
    event, err := github.ParseWebHook(hookType, payload)
    if err != nil {
		log.Printf("CICD: parse error on: %v, %v, %v", event, hookType, err)
        w.WriteHeader(http.StatusBadRequest)
        return
    }

    if push, ok := event.(*github.PushEvent); ok && push.GetRef() == "refs/heads/chat" {
        go deploy()
    }

    w.WriteHeader(http.StatusOK)
}

func validateSignature(sigHeader string, body []byte, secret string) bool {
    const prefix = "sha256="
    if !strings.HasPrefix(sigHeader, prefix) {
        return false
    }
    sig, err := hex.DecodeString(sigHeader[len(prefix):])
    if err != nil {
        return false
    }
    mac := hmac.New(sha256.New, []byte(secret))
    mac.Write(body)
    return hmac.Equal(mac.Sum(nil), sig)
}

// TODO: is this necerssary if watch files exists?
func deploy() {
    if err := exec.Command("git", "pull", "origin", "main").Run(); err != nil {
        log.Println("git pull:", err)
        return
    }

    exec.Command("bun", "run", "tailwindcss", "-i", "./public/css/tw-input.css", "-o", "./public/css/tw-output.css").Run()
    exec.Command("go", "run", "github.com/sqlc-dev/sqlc/cmd/sqlc@v1.27.0", "generate").Run()
    exec.Command("templ", "generate").Run()

    exe, err := os.Executable()
    if err != nil {
        log.Fatal(err)
    }
    if err := exec.Command("go", "build", "-o", exe, ".").Run(); err != nil {
        log.Fatal(err)
    }
    syscall.Exec(exe, append([]string{exe, "--serve"}, os.Args[1:]...), os.Environ())
}
