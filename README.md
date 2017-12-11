# Stymie

## Features

- No need to remember another password, `stymie` will use the default GPG key.
- Create any number of key fields beyond the default username and password.  For example:

```
username: 'derp',
password: '1234',
SSN: '123-45-6789',
securityAnswer: '1st Avenue'
```

- Since everything is stored in `.stymie.d/`, it's easy to port between systems.
- GPG end-to-end encryption allows `stymie` to be safely versioned.
- Generate passwords using [Diceware], [Sillypass] or enter a custom password.

## Security Features

- Uses GPG/PGP public-key cryptography to encrypt everything (even configs).
- Uses OS-level permissions-based access control so only the user can view and list any files created by the user.
- Encrypts using the `--hidden-recipient` flag so as to not include the recipient's key ID in the encrypted file.
- Optionally, asks to set `$HISTIGNORE` so `stymie` commands aren't stored in history [1].

[1] As an alternative to setting `$HISTIGNORE`, most shells by default allow for any command preceded by a `[[SPACE]]` to be ignored by history. Check the value of `$HISTCONTROL` for support.

## Installation

go get github.com/btoll/stymie-go

## Suggestions

- Use `gpg-agent` to save typing.

## Examples

- Create the `example.com` key:
```
stymie add example.com
```

- Edit the `example.com` key:
```
stymie edit example.com
```

- Get just the `username` field value from the `example.com` key:
```
stymie get example.com --field username
```

- Get just the `password` field value from the `example.com` key and copy it to the system clipboard:
```
stymie get example.com --field password | xsel -b
```

## Usage

    Command | Description
    ------- | --------
    add | Adds a new entry
    edit | Edits an entry
    generate | Generates a diceware passphrase
    get | Retrieves an entry
    has | Checks if the entry exists
    init | Installs the password file directory and config file
    list | List all entries
    rm | Deletes an entry

### Options

    Option | Description
    ------- | --------
    --field, -f | Gets the specified key value from an entry (only with `get` command)
    --shadow | Obfuscates a password entry
    -h, --help | Display help

## Ports

- [JavaScript][js]

## License

[GPLv3](COPYING)

## Author

Benjamin Toll

[Diceware]: https://github.com/btoll/diceware
[Sillypass]: https://github.com/btoll/sillypass-go
[js]: https://github.com/btoll/stymie

