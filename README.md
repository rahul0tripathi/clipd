## Clipd

> clipd is helps encrypting/decrypting text on the local clipboard/pasteboard

### FEATURES
- clipd uses the local clipboard/pasteboard service and can directly encrypt and decrypt into the clipboard
- to encrypt copy the password and press `Cmd+shift+E` which changes the copied password to a formatted string which can be saved any where
> [WHOAMI]!0a4a2860-5857-42f2-ac91-3e341b115180!331168d55b964ad6ee3841f0b04558bc549f55748206332c6e62f3cd527e049df8c915d1bdbd5542d794327f5f75361c9ce6e598d2
- to decrypt, copy the saved encrypted string and press `Cmd+shift+D` and the contents are decrypted into the clip board and can be pasted as original password
### SETUP
> to setup run `clipd init` and follow along the wizard
[![asciicast](https://asciinema.org/a/WQusPav0PftLUEAxBBGLRwkwk.svg)](https://asciinema.org/a/WQusPav0PftLUEAxBBGLRwkwk)
- Now clipd can be started using command and can keep running in the background

  $ clipd  server --salt randomsalt --keydir your_key_directory_given_in_setup

### Encrypting
[![asciicast](https://asciinema.org/a/XnQM0XH4gbc2XyFrAOUjlW9Wm.svg)](https://asciinema.org/a/XnQM0XH4gbc2XyFrAOUjlW9Wm)

### Decryption
[![asciicast](https://asciinema.org/a/7Hlj68zSAz9s1sxYUUtDr9BbW.svg)](https://asciinema.org/a/7Hlj68zSAz9s1sxYUUtDr9BbW)