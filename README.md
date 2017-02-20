# ENCRYPT FILE WITH AES-GCM

I couldn't encrypt a file with AES-GCM using the OpenSSL command line tool, so I made this.

Use it at your own risks. Inb4 someone accuses me of killing innocent people.

## Usage

`go run main.go -encrypt -file [your-file]` will encrypt a file AND give you a one-time 256-bit AES key.

You're suppose to upload that file somewhere, and send the key to your recipient in a separate channel.

The recipient can use the key and the file like that:

`go run main.go -decrypt -file [encrypted-file] -key [hex-key]`
