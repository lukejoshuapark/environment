![](icon.png)

# environment

Automatically retrieves values from environment variables into a configuration
structure.

## Usage Example

```go
package main

import (
	"github.com/lukejoshuapark/environment"
)

type Config struct {
	Hostname string `environment:"APP_HOSTNAME,example.com"`
	Username string `environment:"APP_USERNAME"`
	Password string `environment:"APP_PASSWORD"`
}

func main() {
	cfg := &Config{}
	if err := environment.Populate(cfg); err != nil {
		// Failed to populate the provided structure.
		return
	}

	// cfg has now been populated from environment variables.
}
```

## Parsers

It is often useful to be able to represent richer types than `string` in our
configuration, even though we are forced to provide these values in `string`
form through environment variables.

Consider the scenario where we might want to source a port number from our
environment variables:

```go
type Config struct {
	Hostname string `environment:"APP_HOSTNAME,example.com"`
	Port     uint16 `environment:"APP_PORT,8080"` // Non-string type.
	Username string `environment:"APP_USERNAME"`
	Password string `environment:"APP_PASSWORD"`
}
```

By default, `environment.Populate` won't know how to handle the `uint16` in this
struct and will fail.  To show `environment.Populate` how to handle this type,
we can use a parser:

```go
func ParseUInt16(val string) (uint16, error) {
	n, err := strconv.ParseUint(val, 10, 16)
	if err != nil {
		return 0, err
	}

	return uint16(n), nil
}

environment.UseParser(ParseUInt16)
```

You can define a parser function for any type you like.  An example of a more
complex parser than can read base64 encoded Ed25519 private keys is provided
below:

```go
func ParseEd25519PrivateKey(val string) (ed25519.PrivateKey, error) {
	rawKey, err := base64.StdEncoding.DecodeString(val)
	if err != nil {
		return nil, fmt.Errorf("could not base64 decode private key: %w", err)
	}

	if len(rawKey) != ed25519.SeedSize {
		return nil, fmt.Errorf("expected a key length of %v bytes, but was %v bytes", ed25519.SeedSize, len(rawKey))
	}

	return ed25519.NewKeyFromSeed(rawKey), nil
}

environment.UseParser(ParseEd25519PrivateKey)
```

## Optional Environment Variables

By default, fields tagged with `environment` are required - if the
corresponding environment variable is not set, `Populate` will return a non-nil
error when processing.

If instead you'd like to have a field be optional, you can provide a default
value by suffixing a comma and the default string value for that field.

This can be seen above for the `Hostname` field.  If `APP_HOSTNAME` is set in
the environment, its value will be assigned to the `Hostname` field.  Otherwise,
the value `"example.com"` will be assigned.

Note that default values are assigned prior to parsing.  Because of this, for
non-string types, you must provide a default value in string form that can
successfully be parsed by the parser for that type.

---

Icons made by [Freepik](https://www.flaticon.com/authors/freepik) from
[www.flaticon.com](https://www.flaticon.com/).
