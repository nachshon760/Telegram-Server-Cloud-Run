package main

import (
	"context"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gotd/td/internal/crypto"
	"github.com/gotd/td/tgtest"
)

// מפתח ה-RSA הפרטי של השרת
const privateKeyPEM = `-----BEGIN RSA PUBLIC KEY-----
MIIEowIBAAKCAQEAy3W81o1S7qH3n8XN9rU2hMh3yJv9XW18O2J6hX8n4vKjSdB7
m9Z6Y7t6X8n4vKjSdB7m9Z6Y7t6X8n4vKjSdB7m9Z6Y7t6X8n4vKjSdB7m9Z6Y7t
6X8n4vKjSdB7m9Z6Y7t6X8n4vKjSdB7m9Z6Y7t6X8n4vKjSdB7m9Z6Y7t6X8n4v
KjSdB7m9Z6Y7t6X8n4vKjSdB7m9Z6Y7t6X8n4vKjSdB7m9Z6Y7t6X8n4vKjSdB7
m9Z6Y7t6X8n4vKjSdB7m9Z6Y7t6X8n4vKjSdB7m9Z6Y7t6X8n4vKjSdB7m9Z6Y7
t6X8n4vKjSdB7m9Z6Y7t6X8n4vKjSdB7m9Z6Y7wIFwIDAQAB
AoIBADCCp1T4vV3E5Gv3XJnzX7n7Y9fB3n8XN9rU2hMh3yJv9XW18O2J6hX8n4vK
jSdB7m9Z6Y7t6X8n4vKjSdB7m9Z6Y7t6X8n4vKjSdB7m9Z6Y7t6X8n4vKjSdB7m
9Z6Y7t6X8n4vKjSdB7m9Z6Y7t6X8n4vKjSdB7m9Z6Y7t6X8n4vKjSdB7m9Z6Y7t
6X8n4vKjSdB7m9Z6Y7t6X8n4vKjSdB7m9Z6Y7t6X8n4vKjSdB7m9Z6Y7t6X8n4v
KjSdB7m9Z6Y7t6X8n4vKjSdB7m9Z6Y7t6X8n4vKjSdB7m9Z6Y7t6X8n4vKjSdB7
m9Z6Y7t6X8n4vKjSdB7m9Z6Y7t6X8n4vKjSdB7m9Z6Y7t6X8n4vKjSdB7m9ZCEC
gYEA8d7M1o1S7qH3n8XN9rU2hMh3yJv9XW18O2J6hX8n4vKjSdB7m9Z6Y7t6X8n
4vKjSdB7m9Z6Y7t6X8n4vKjSdB7m9Z6Y7t6X8n4vKjSdB7m9Z6Y7t6X8n4vKjSdB
7m9Z6Y7t6X8n4vKjSdB7m9Z6Y7t6X8n4vKjSdB7m9Z6Y7t6X8n4vKjSdB7m9ZCE
CgYEA1uW81o1S7qH3n8XN9rU2hMh3yJv9XW18O2J6hX8n4vKjSdB7m9Z6Y7t6X8
n4vKjSdB7m9Z6Y7t6X8n4vKjSdB7m9Z6Y7t6X8n4vKjSdB7m9Z6Y7t6X8n4vKjS
dB7m9Z6Y7t6X8n4vKjSdB7m9Z6Y7t6X8n4vKjSdB7m9Z6Y7t6X8n4vKjSdB7m9Z
CECgYEA3XW81o1S7qH3n8XN9rU2hMh3yJv9XW18O2J6hX8n4vKjSdB7m9Z6Y7t6
X8n4vKjSdB7m9Z6Y7t6X8n4vKjSdB7m9Z6Y7t6X8n4vKjSdB7m9Z6Y7t6X8n4vK
jSdB7m9Z6Y7t6X8n4vKjSdB7m9Z6Y7t6X8n4vKjSdB7m9Z6Y7t6X8n4vKjSdB7m
9ZCECgYEAyuW81o1S7qH3n8XN9rU2hMh3yJv9XW18O2J6hX8n4vKjSdB7m9Z6Y7
t6X8n4vKjSdB7m9Z6Y7t6X8n4vKjSdB7m9Z6Y7t6X8n4vKjSdB7m9Z6Y7t6X8n4
vKjSdB7m9Z6Y7t6X8n4vKjSdB7m9Z6Y7t6X8n4vKjSdB7m9Z6Y7t6X8n4vKjSdB
7m9ZCECgYEA6XW81o1S7qH3n8XN9rU2hMh3yJv9XW18O2J6hX8n4vKjSdB7m9Z6
Y7t6X8n4vKjSdB7m9Z6Y7t6X8n4vKjSdB7m9Z6Y7t6X8n4vKjSdB7m9Z6Y7t6X8
n4vKjSdB7m9Z6Y7t6X8n4vKjSdB7m9Z6Y7t6X8n4vKjSdB7m9Z6Y7t6X8n4vKjS
dB7m9ZCeQ==
-----END RSA PUBLIC KEY-----`

func main() {
	// 1. פענוח המפתח הפרטי של השרת
	block, _ := pem.Decode([]byte(privateKeyPEM))
	if block == nil {
		log.Fatal("failed to parse PEM block containing the private key")
	}
	privKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		log.Fatalf("failed to parse private key: %v", err)
	}

	key := crypto.RSAPrivateKey{
		PrivateKey: privKey,
	}

	// 2. יצירת סימולטור שרת טלגרם רשמי (שמזהה פקודות אימות ורושם משתמשים)
	ctx := context.Background()
	testServer := tgtest.NewServer(tgtest.Key(key), tgtest.WithLogger(nil))
	defer testServer.Close()

	// 3. קביעת פורט הריצה של גוגל
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// 4. הגדרת נתיב החיבור לאפליקציה
	http.HandleFunc("/tg", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("=== התקבלה בקשת חיבור MTProto חדשה מהאפליקציה ===")
		
		// העברת הבקשה ישירות למנוע ה-MTProto שמבין את כל השפה של טלגרם
		testServer.ServeHTTP(w, r)
	})

	fmt.Printf("MTProto Server is running on port %s...\n", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}
