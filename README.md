# PWizard

Generate secure passwords in batches using predefined JSON templates.

Passwords are printed to stdout only; the utility does not perform any disk writes unless you redirect output.

## How to

1. Build

```bash
make
cp pwizard /usr/local/bin/pwizard
```

2. Configure passwords.json

### Configuration Template

| Key | Value type | Default value | Description |
| - | - | - | - |
| name | string | | password type or identifier |
| length | int | | password length (minimum: 4 symbols) |
| symbols | bool | false | enable/disable safe set of special symbols |
| symbols_extended | bool | false | enable/disable extended set of special symbols |
| symbols_dangerous | bool | false | enable/disable dangerous set of special symbols |
| custom_symbols_set | string | | specify custom set of symbols |
| digits | bool | true | enable/disable digits |
| letters | bool | true | enable/disable letters |
| uppercase_only | bool | false | enable/disable uppercase-only letters set |
| lowercase_only | bool | false | enable/disable lowercase-only letters set |

- Special symbols are disabled by default
- Digits are enabled by default
- Uppercase and lowercase letters are enabled by default

Password name and length specifications are required.

Symbol sets cannot be combined. Priority order:

1. custom_symbols_set
2. symbols_dangerous
3. symbols_extended
4. symbols

### Example Configuration

```json
[
  {
    "name": "database",
    "length": 16,
    "symbols": true
  },
  {
    "name": "api_key", 
    "length": 24
  }
]
```

3. Run `pwizard`

4. Or run it with custom config path

```bash
pwizard -config=$HOME/config.json
```

Default configuration path is `/etc/pwizard/passwords.json` 

## Output

Currently supported: human-readable string format

Planned: JSON, YAML

## Security notes

- Passwords are generated using cryptographically secure random generator
- Minimum password length is enforced (4 characters)
- No passwords are stored or logged
