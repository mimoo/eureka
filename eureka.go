//
// Eureka
// ======
// Eureka is a handy utility to encrypt files and folders. It follows several principles:
//
// - I want to encrypt and send a file to someone or myself.
// - Eureka should be easy to install and share.
// - PGP is too cumbersome to use, I want something simple (right-click > encrypt).
// - I already share two separate and secure channels with the recipient (mail + signal for example).
//
// Here how the code is organized:
//
// - eureka.go 		// the main code to encrypt files
// - folders.go 	// the code to compress folders
// - ui_windows.go 	// the code to add right-click encrypt/decrypt on windows
// - ui_macOS.go 	// the code to add right-click encrypt/decrypt on macOS
// - ui_linux.go 	// the code to add right-click encrypt/decrypt on linux

package main

import (
	"bufio"
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	// cross-platform clipboard
	"github.com/atotto/clipboard"
	// cross-platform GUI
	"github.com/sqweek/dialog"
)

// open a link in your favorite browser
func openBrowser(url string) {
	var err error

	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}
	if err != nil {
		panic(err)
	}
}

// prompt for the key with a GUI
func promptKey() string {

	ok := dialog.Message("%s", "Do you want to use your clipboard as the key?").Title("Eureka").YesNo()
	if ok { // clipboard option
		key, err := clipboard.ReadAll()
		fmt.Println([]byte(key))
		if err != nil {
			dialog.Message("%s", "error: couldn't read the key from clipboard").Title("Eureka").Info()
			return ""
		}
		return key
	}

	// terminal option
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter 256-bit hexadecimal key: ")
	key, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("couldn't read the key")
		return ""
	}

	return key
}

func main() {
	var err error

	// parse flags
	encrypt := flag.Bool("encrypt", false, "to encrypt")
	about := flag.Bool("about", false, "to get redirected to github.com/mimoo/eureka")
	decrypt := flag.Bool("decrypt", false, "to decrypt")
	keyHex := flag.String("key", "", "256-bit key")
	inFile := flag.String("file", "", "file to encrypt or decrypt")
	setupUI := flag.String("setup-ui", "", "to setup interactivity in your OS, pass the directory where Eureka is installed as argument")

	flag.Parse()

	// redirect to github.com/mimoo/eureka
	if *about {
		openBrowser("https://www.github.com/mimoo/eureka")
		return
	}

	// setup user-interaction
	if *setupUI != "" {
		if *setupUI == "remove" {
			uninstall()
		} else {
			fmt.Println("setting up right-click context menus")
			install(strings.Trim(*setupUI, `"`))
		}
		return
	}

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

	// nonce = 1111...
	nonce := bytes.Repeat([]byte{1}, 12)

	// key = ?
	var key []byte

	if *encrypt { // generate random key if we're encrypting
		key = make([]byte, 32)
		if _, err = io.ReadFull(rand.Reader, key); err != nil {
			fmt.Println("error: randomness cannot be generated on your system")
			flag.Usage()
			return
		}
	} else { // get key from flag if we are decrypting
		if *keyHex == "" { // if flag key is empty, prompt the user
			*keyHex = promptKey()
		}
		// decode and check key
		key, err = hex.DecodeString(*keyHex)
		if err != nil || len(key) != 32 {
			fmt.Println(err)
			fmt.Println("error: the key has to be a 256-bit hexadecimal string")
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

	// encrypt or decrypt
	var contentAfter []byte

	if *encrypt {
		// compress file or folder
		var buf bytes.Buffer
		if err := compress(*inFile, &buf); err != nil {
			fmt.Println(err)
			return
		}
		// encrypt compressed content
		contentAfter = AESgcm.Seal(nil, nonce, buf.Bytes(), nil)
		// write file to disk
		_, outFile := filepath.Split(*inFile)
		outFile = outFile + ".encrypted"
		if err = ioutil.WriteFile(outFile, contentAfter, 0600); err != nil {
			fmt.Println(err)
			return
		}
		// place key in clipboard
		stringKey := fmt.Sprintf("%032x", key)
		// notification
		ok := dialog.Message("File encrypted at %s\nyour recipient will need Eureka to decrypt the file: https://github.com/mimoo/eureka\nIn a different secure channel, pass the following one-time key to your recipient.\n%s\nDo you want to use your clipboard as the key?", outFile, stringKey).Title("Eureka").YesNo()

		if ok { // copy to clipboard
			clipboard.WriteAll(stringKey)
		} else { // print to terminal and pause
			fmt.Println(stringKey)
			fmt.Scanln()
		}
	} else {
		// open file
		content, err := ioutil.ReadFile(*inFile)
		if err != nil {
			fmt.Println("error: cannot open input file")
			flag.Usage()
			return
		}
		// decrypt
		contentAfter, err = AESgcm.Open(nil, nonce, content, nil)
		if err != nil {
			fmt.Println("error: cannot decrypt. The key is not correct or someone tried to modify your file.")
			return
		}
		// create a decrypted folder
		if _, err := os.Stat("./decrypted"); err != nil {
			if err := os.MkdirAll("./decrypted", 0755); err != nil {
				fmt.Println("error: cannot create folder 'decrypted'")
				return
			}
		} else {
			fmt.Println("error: the folder 'decrypted' already exists. Decrypting the file could overwrite files.")
			return
		}
		// decompress it
		buf := bytes.NewReader(contentAfter)
		if err := decompress(buf, "./decrypted"); err != nil {
			fmt.Println(err)
			return
		}
		// notification
		dialog.Message("File decrypted at decrypted/\nCheers.").Title("Eureka").Info()
	}
}
