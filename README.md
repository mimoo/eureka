# ENCRYPT FILE WITH AES-GCM

This is a simple tool to encrypt and decrypt files.

## Install

[Get it here](https://github.com/mimoo/EncryptFileWithAES-GCM/releases) or build it yourself.

## Usage

`./EncryptFileWithAESGCM -encrypt -file [your-file]` will encrypt a file AND give you a one-time 256-bit AES key.

You're supposed to upload that file somewhere, and send the key to your recipient in a separate channel.

The recipient can use the key and the file like that:

`./EncryptFileWithAESGCM -decrypt -file [encrypted-file] -key [hex-key]`
