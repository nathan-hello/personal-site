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

	sig := r.Header.Get("X-Hub-Signature-256")
    if !validateSignature(sig, payload, parsed.WEBHOOK_SECRET) {
		log.Printf("CICD: invalid signature %v", sig)
        w.WriteHeader(http.StatusUnauthorized)
        return
    }

    event, err := github.ParseWebHook(github.WebHookType(r), payload)
    if err != nil {
		log.Printf("CICD: parse error on: %v, %v", event, err)
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

func deploy() {
    if err := exec.Command("git", "pull").Run(); err != nil {
        log.Println("git pull:", err)
        return
    }

    exec.Command("make", "build").Run()

    syscall.Exec("./personal-site", []string{"--build", "--serve"}, os.Environ())
}
