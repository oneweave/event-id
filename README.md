# event-id

A Go package for generating and working with event-id (Prefixed UUIDv7, formerly puidv7) identifiers. This repository is a fork of the original [puidv7-go](https://github.com/puidv7/puidv7-go).

## Fork & Refactoring Summary
This fork introduces major optimizations and design changes:
- **Performance Optimizations**: Removed all regex matching and compile overhead, resulting in 50x-200x faster encode/decode speeds and reducing memory allocations to nearly zero (using stack-allocated buffers).
- **Strict Parsing**: Removed the Crockford look-alike character mappings (e.g. mapping `i`/`l` to `1`, `o` to `0`) to enforce strict, machine-only canonical representations.
- **Robust UUID Handling**: Leveraged `github.com/google/uuid`'s native binary operations instead of custom string parsing, fixing formatting errors.


## What is an event-id?

It's a prefixed UUIDv7 which:

1. Is encoded using base32 crockford encoding
2. Is always lowercase
3. Does not contain any hyphens
4. Is prefixed with a 3 character alphabetic (a-z) prefix

## UUIDv7 <> event-id conversion example

Check out the online converter at <https://puidv7.dev> or step through the following examples...

With `acc` prefix:

- UUIDv7 = `01970a1c-e31e-7422-9cd5-e9651d11cc97`

- event-id = `acc06bgm7733st2576nx5jht4ecjw`

How to manually verify:

- Remove dashes from the UUIDv7 to get the HEX-encoded string

  e.g. `01970a1c-e31e-7422-9cd5-e9651d11cc97` becomes `01970a1ce31e74229cd5e9651d11cc97`

- Convert the HEX-encoded string to crockford base32
  e.g. Use <https://cryptii.com/pipes/crockford-base32> with Bytes (in Hexadecimal format), and Encode to Base32 (Crockford's Base32 variant).

  i.e. `01970a1ce31e74229cd5e9651d11cc97` becomes `06BGM7733ST2576NX5JHT4ECJW`

- Convert the value to lowercase and add the prefix.

  e.g. `06BGM7733ST2576NX5JHT4ECJW` becomes `06bgm7733st2576nx5jht4ecjw`

  then `06bgm7733st2576nx5jht4ecjw` becomes `acc06bgm7733st2576nx5jht4ecjw`

## Why does the world need event-id?

Because:

1. UUIDv7 is great for databases and distributed systems
2. UUIDv7 is not as great for end users/humans

Advantages of the human-friendly event-id format:

1. URL-safe and case insensitive
2. Shorter than UUIDv7 (29 characters vs 36)
3. Easier to copy & paste (no hyphens)
4. Types can be inferred (great for customer support, future-proofing APIs)

## Usage

```go
package main

import (
    "fmt"
    "github.com/oneweave/event-id"
)

func main() {
    // Generate a new event-id with prefix "acc"
    id, err := eventid.New("acc")
    if err != nil {
        panic(err)
    }
    fmt.Println(id) // e.g. "acc069rz3kw7dyyz2gj28t5cy4tqg"

    // Encode an existing UUID with a prefix
    encoded, err := eventid.Encode("01938f8e-7c3b-7def-8a12-123456789abc", "acc")
    if err != nil {
        panic(err)
    }
    fmt.Println(encoded) // "acc069rz3kw7dyyz2gj28t5cy4tqg"

    // Decode an event-id back to UUID (with optional prefix validation)
    uuid, err := eventid.Decode(encoded, "acc")
    if err != nil {
        panic(err)
    }
    fmt.Println(uuid) // "01938f8e-7c3b-7def-8a12-123456789abc"
}
```

If you would like to use event-id in Go with a
[github.com/go-playground/validator/v10](https://github.com/go-playground/validator/v10)
validator:

```go
func ValidateEventID(fl validator.FieldLevel) bool {
	_, err := eventid.Decode(fl.Field().String(), "")
	if err != nil {
		return false
	}
	return true
}
func Example() {
  validate := validator.New(validator.WithRequiredStructEnabled())
  if err := validate.RegisterValidation("eventid", ValidateEventID); err != nil {
    fmt.Printf("error registering eventid validator: %v\n", err)
  }
  if err := validate.Struct(struct {
    ID    string `validate:"required,eventid"`
  }{
    ID:    "abc06awcb4f5hzmfey7qwt7s8a6q4",
  }); err != nil {
    fmt.Printf("validation errors: %+s\n", err)
  }
}
```

## API

### `New(prefix string) (string, error)`

Generates a new UUIDv7 and encodes it with the given 3-character lowercase prefix.

### `Encode(uuid string, prefix string) (string, error)`

Encodes an existing UUID into an event-id string with the given prefix.

### `Decode(id string, prefix string) (string, error)`

Decodes an event-id string back to a UUID. If prefix is provided, validates that the ID starts with that prefix.

## Format

- **Prefix**: 1 to 5 lowercase letters (a-z)
- **Separator**: Underscore `_`
- **Encoded UUID**: 26 characters using Crockford Base32 alphabet
- **Total length**: 28 to 32 characters (depending on prefix length)

Example: `acc_069rz3kw7dyyz2gj28t5cy4tqg`

## Further Reading

- Stripe's [Designing APIs for humans: Object IDs](https://dev.to/stripe/designing-apis-for-humans-object-ids-3o5a) article.
- Buildkite's [Goodbye integers. Hello UUIDv7!](https://buildkite.com/resources/blog/goodbye-integers-hello-uuids/) article.
- Douglas [Crockford Base32](https://www.crockford.com/base32.html) page.
- [RFC 4122](https://datatracker.ietf.org/doc/html/rfc4122) - A Universally Unique Identifier (UUID) URN Namespace.

## License

event-id is licensed under the Apache License, Version 2.0.
See the [LICENSE](./LICENSE) file for details.
