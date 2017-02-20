package main

import(
	"crypto/cipher"
	"crypto/aes"
	"crypto/rand"
	"fmt"
	"io"
	"bytes"
	"io/ioutil"
	"flag"
	"encoding/hex"
	"time"
)

func main(){
	var err error
	
	// parse
	encrypt := flag.Bool("encrypt", false, "to encrypt")
	decrypt := flag.Bool("decrypt", false, "to decrypt")

	keyHex := flag.String("key", "", "256-bit key")
	inFile := flag.String("file", "", "file to encrypt or decrypt")

	flag.Parse()

	// validate encrypt or decrypt
	if *encrypt == false && *decrypt == false {
		fmt.Println("do you want to encrypt or decrypt? Use -h to get help.")
		return
	}

	//
	var key []byte

	nonce := bytes.Repeat([]byte{0}, 12)

	if *encrypt {
		// create key
		key = make([]byte, 32)
		if _, err = io.ReadFull(rand.Reader, key); err != nil {
			fmt.Println("Can't create key")
			return
		}
	} else {
		key, err = hex.DecodeString(*keyHex)
		if err != nil || len(key) != 32 {
			fmt.Println("This is not a 256-bit key")
		}
	}
	
	// create AES-GCM instance
	cipherAES, err := aes.NewCipher(key)

	if err != nil {
		fmt.Println("Can't instantiate AES")
		return
	}

	AESgcm, err := cipher.NewGCM(cipherAES)

	if err != nil {
		fmt.Println("Can't instantiate GCM")
		return
	}

	// open file
	content, err := ioutil.ReadFile(*inFile)
	if err != nil {
		fmt.Println("cannot open input file")
		return
	}

	// encrypt / decrypt
	var content_after []byte
	
	if *encrypt {
		content_after = AESgcm.Seal(nil, nonce, content, nil)
	} else { // decrypt
		content_after, err = AESgcm.Open(nil, nonce, content, nil)
		if err != nil {
			fmt.Println("Cannot decrypt. The key is not correct or someone tried to modify your file.")
			return
		}
	}

	// write file
	var outFile string
	if *decrypt {
		outFile = "DECRYPTED_FILE"
	} else {
		outFile = "ENCRYPTED_FILE"
	}

	now := time.Now()
	outFile += now.Format("_2006-01-02_03-04-05")

	err = ioutil.WriteFile(outFile, content_after, 0644)
	if err != nil {
		if *decrypt {
			fmt.Printf("can't write file at %s\n", outFile)
		} else {
			fmt.Printf("can't write file at %s\n", outFile)
		}
		return
	}

	if *decrypt {
		fmt.Printf("file decrypted at %s\n", outFile)
	} else {
		fmt.Printf("file encrypted at %s\n", outFile)
		fmt.Println("In a different secure channel, pass the following one-time key to your recipient.")
		fmt.Printf("%032x\n", key)
	}
	
	//
}
