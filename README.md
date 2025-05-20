# Packer Plugin Crypt

This plugin provides a way of generating password hashes. 

```
data "crypt-mkpasswd" "default" {
  plaintext = "password"
  salt = "salt"
}
```

When no salt is provided, 16 bits of random data from crypt/rand is used.
