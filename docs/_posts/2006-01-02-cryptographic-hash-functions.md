---
layout: default
title: Cryptographic Hash Functions - Reference Manual - csvq
category: reference
---

# Cryptographic Hash Functions

| name | description |
| :- | :- |
| [MD5](#md5) | Generate a MD5 hash value |
| [SHA1](#sha1) | Generate a SHA-1 hash value |
| [SHA256](#sha256) | Generate a SHA-256 hash value |
| [SHA512](#sha512) | Generate a SHA-512 hash value |
| [MD5_HMAC](#md5_hmac) | Generate a MD5 keyed-hash value |
| [SHA1_HMAC](#sha1_hmac) | Generate a SHA-1 keyed-hash value |
| [SHA256_HMAC](#sha256_hmac) | Generate a SHA-256 keyed-hash value |
| [SHA512_HMAC](#sha512_hmac) | Generate a SHA-512 keyed-hash value |

## Definitions

### MD5
{: #md5}

```
MD5(data)
```

_data_
: [string]({{ '/reference/value.html#string' | relative_url }})

_return_
: [string]({{ '/reference/value.html#string' | relative_url }})

Generates a MD5 hash value.

### SHA1
{: #sha1}

```
SHA1(data)
```

_data_
: [string]({{ '/reference/value.html#string' | relative_url }})

_return_
: [string]({{ '/reference/value.html#string' | relative_url }})

Generates a SHA-1 hash value.

### SHA256
{: #sha256}

```
SHA256(data)
```

_data_
: [string]({{ '/reference/value.html#string' | relative_url }})

_return_
: [string]({{ '/reference/value.html#string' | relative_url }})

Generates a SHA-256 hash value.

### SHA512
{: #sha512}

```
SHA512(data)
```

_data_
: [string]({{ '/reference/value.html#string' | relative_url }})

_return_
: [string]({{ '/reference/value.html#string' | relative_url }})

Generates a SHA-512 hash value.

### MD5_HMAC
{: #md5_hmac}

```
MD5_HMAC(data, key)
```

_data_
: [string]({{ '/reference/value.html#string' | relative_url }})

_key_
: [string]({{ '/reference/value.html#string' | relative_url }})

_return_
: [string]({{ '/reference/value.html#string' | relative_url }})

Generates a MD5 keyed-hash value using the HMAC method.

### SHA1_HMAC
{: #sha1_hmac}

```
SHA1_HMAC(data, key)
```

_data_
: [string]({{ '/reference/value.html#string' | relative_url }})

_key_
: [string]({{ '/reference/value.html#string' | relative_url }})

_return_
: [string]({{ '/reference/value.html#string' | relative_url }})

Generates a SHA-1 keyed-hash value using the HMAC method.

### SHA256_HMAC
{: #sha256_hmac}

```
SHA256_HMAC(data, key)
```

_data_
: [string]({{ '/reference/value.html#string' | relative_url }})

_key_
: [string]({{ '/reference/value.html#string' | relative_url }})

_return_
: [string]({{ '/reference/value.html#string' | relative_url }})

Generates a SHA-256 keyed-hash value using the HMAC method.

### SHA512_HMAC
{: #sha512_hmac}

```
SHA512_HMAC(str, key)
```

_data_
: [string]({{ '/reference/value.html#string' | relative_url }})

_key_
: [string]({{ '/reference/value.html#string' | relative_url }})

_return_
: [string]({{ '/reference/value.html#string' | relative_url }})

Generates a SHA-512 keyed-hash value using the HMAC method.
