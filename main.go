package main

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"path/filepath"
	"strings"
)

func main() {
	// stupid vars
	var err error

	// parse flags
	encrypt := flag.Bool("encrypt", false, "to encrypt")
	decrypt := flag.Bool("decrypt", false, "to decrypt")
	keyHex := flag.String("key", "", "256-bit key")
	inFile := flag.String("file", "", "file to encrypt or decrypt")

	flag.Parse()

	if *encrypt == false && *decrypt == false {
		fmt.Println("===================ᕙ(⇀‸↼‶)ᕗ===================")
		fmt.Println(" Eureka is a tool to help you encrypt/decrypt a file")
		fmt.Println(" to encrypt:")
		fmt.Println("     eureka -encrypt -file [your-file]")
		fmt.Println(" to decrypt:")
		fmt.Println("     eureka -decrypt -file [encrypted-file] -key [hex-key]")
		fmt.Println("===================ᕙ(⇀‸↼‶)ᕗ===================")
		flag.Usage()
		return
	}

	// nonce = 0
	nonce := bytes.Repeat([]byte{0}, 12)

	// key = ?
	var key []byte

	if *encrypt {
		// create key
		key = make([]byte, 32)
		if _, err = io.ReadFull(rand.Reader, key); err != nil {
			fmt.Println("Can't create key")
			flag.Usage()
			return
		}
	} else {
		// get hex key from flag
		key, err = hex.DecodeString(*keyHex)
		if err != nil || len(key) != 32 {
			fmt.Println("This is not a 256-bit key")
			flag.Usage()
			return
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

	// open input file
	content, err := ioutil.ReadFile(*inFile)
	if err != nil {
		fmt.Println("cannot open input file")
		flag.Usage()
		return
	}

	// encrypt or decrypt
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

	// write output file
	var outFile string

	if *decrypt {
		var path string
		path, outFile = filepath.Split(*inFile)
		if filepath.Ext(outFile) == ".encrypted" && len(outFile) > 10 {
			fmt.Println(outFile)
			outFile = strings.TrimSuffix(outFile, ".encrypted")
			fmt.Println(outFile)
		} else {
			outFile = outFile + ".decrypted"
		}
		outFile = path + outFile
	} else {
		_, outFile = filepath.Split(*inFile)
		outFile = outFile + ".encrypted"
	}

	//	now := time.Now()
	//	outFile += now.Format("_2006-01-02_03-04-05")

	err = ioutil.WriteFile(outFile, content_after, 0644)

	if err != nil {
		if *decrypt {
			fmt.Printf("Can't write file at %s\n", outFile)
		} else {
			fmt.Printf("Can't write file at %s\n", outFile)
		}
		return
	}

	if *decrypt {
		fmt.Printf("File decrypted at %s\n", outFile)
		fmt.Println("Cheers.")
	} else {
		fmt.Printf("File encrypted at %s\n", outFile)
		fmt.Println("In a different secure channel, pass the following one-time key to your recipient.")
		fmt.Printf("%032x\n", key)
	}

	//
}
