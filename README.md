# ENCRYPT A FILE WITH EASE

This is a simple tool to encrypt and decrypt files.

## Install

[Get a binary here](https://github.com/mimoo/eureka/releases).

If you have [Go]() installed and `/usr/local/go/bin` is in your PATH, you should be able to simply get the binary by doing

```
go get github.com/mimoo/eureka
```

## Usage

`./eureka -encrypt -file [your-file]` will encrypt a file AND give you a one-time 256-bit AES key.

You're supposed to upload that file somewhere, and send the key to your recipient in a separate channel.

The recipient can use the key and the file like that:

`./eureka -decrypt -file [encrypted-file] -key [hex-key]`
