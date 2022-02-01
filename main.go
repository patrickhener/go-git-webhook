package main

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/gorilla/mux"
)

var (
	secret string = os.Getenv("WEBHOOK_SECRET")
	ip     string = os.Getenv("WEBHOOK_IP")
	port   string = os.Getenv("WEBHOOK_PORT")
	cmd    string = os.Getenv("WEBHOOK_CMD")
)

func main() {

	mux := mux.NewRouter()

	mux.PathPrefix("/").HandlerFunc(webhookhandler)
	// Make dynamic later on
	addr := fmt.Sprintf("%+v:%+v", ip, port)

	server := http.Server{
		Addr:    addr,
		Handler: mux,
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 120 * time.Second,
		ReadTimeout:  120 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	log.Panic(server.ListenAndServe())

}

func webhookhandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		ok, err := checksig(r)
		if err != nil {
			log.Printf("There was an error checking signature: %+v", err)
			return
		}
		if !ok {
			log.Print("Signature was wrong")
			return
		}

		log.Println("Signature match")

		if err := docmd(); err != nil {
			log.Printf("Error when running command: %+v", err)
		}

		return
	}
}

func docmd() error {
	c := exec.Command(cmd)
	var out bytes.Buffer
	c.Stdout = &out
	err := c.Run()

	if err != nil {
		return err
	}

	log.Printf("Command says: %+v", out.String())

	return nil
}

func checksig(r *http.Request) (bool, error) {
	key := []byte(secret)
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return false, err
	}

	sig := hmac.New(sha256.New, key)
	sig.Write(body)

	signature := hex.EncodeToString(sig.Sum(nil))
	// log.Printf("Signature is: %+v", signature)
	if len(r.Header["X-Hub-Signature-256"]) > 0 {
		headerSignature := r.Header["X-Hub-Signature-256"][0]
		signaturePart := strings.Split(headerSignature, "=")[1]
		// log.Printf("Header Signature says: %+v", signaturePart)

		if headerSignature != "" && signaturePart == signature {
			return true, nil
		}
	}

	return false, nil
}
