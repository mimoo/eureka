![EUREKA](https://i.imgur.com/qSscFjx.png)

Eureka is a simple tool to encrypt files and folders. It works on Windows, Linux and MacOS.

## Security and Status

Eureka is pretty simple, and well commented. Anyone is free to audit the software themselves.

## Install

There are several ways to install Eureka, with more on the way.

**Binary**.

[Get a binary here](https://github.com/mimoo/eureka/releases/tag/1.0).

**Go get**.

If you have [Golang](https://golang.org/) installed and `/usr/local/go/bin` is in your PATH, you should be able to simply get the binary by doing

```
go get github.com/mimoo/eureka
```

**Homebrew**.

If you are on MacOS, just use Homebrew:

```
brew tap mimoo/eureka && brew install eureka
```

## Usage

**1.** You are trying to send *Bob* the file `myfile.txt`.Start by encrypting the file via:

```
eureka -encrypt -file myfile.txt
```

which will return a one-time 256-bit AES key and create a new `myfile.txt.encrypted` file:

```
File encrypted at myfile.txt.encrypted
In a different secure channel, pass the following one-time key to your recipient.
613800fc6cf88f09aa6aeafab3eedd627503e6c5de28040c549efc2c6f80178d
```

**2.** Find a channel to send the encrypted file to *Bob*. It could be via email, or via dropbox, or via google drive, etc.

**3.** You then need to transmit the one-time key (`613800fc6cf88f09aa6aeafab3eedd627503e6c5de28040c549efc2c6f80178d`) to *Bob* in a **different channel**. For example, if you exchanged the file (or a link to the file) via email, then send this key to *Bob* via WhatsApp. 

**If you send both the encrypted file and the one-time key in the same channel, encryption is useless**.

**4.** Once *Bob* receives the file and the one-time key from two different channels, he can decrypt the file via this command:

```
eureka -decrypt -file myfile.txt.encrypted
```

which will create a new file `myfile.txt` under a `decrypted` folder containing the original content.
