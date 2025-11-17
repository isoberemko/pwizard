# PWizard

Generate secure passwords in batches using predefined JSON templates.

## How to

1. Make

```bash
make
cp pwizard /usr/local/bin/pwizard
```

2. Configure passwords.json

### Configuration Template

- `name` - password type or identifier
- `length` - password length (minimum: 4 characters)  
- `symbols` - enable/disable special characters (true/false, default is false)

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
