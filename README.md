# go-windapsearch
[![CircleCI](https://circleci.com/gh/ropnop/go-windapsearch.svg?style=svg)](https://circleci.com/gh/ropnop/go-windapsearch)

`windapsearch` is a tool to assist in Active Directory Domain enumeration through LDAP queries. It contains several modules to enumerate users, groups, computers, as well as perform searching and unauthenticated information gathering.

For usage examples of each of the modules, view the [modules README](pkg/modules/README.md)

In addition to performing common LDAP searches, `windapsearch` now also has the option to convert LDAP results to JSON format for easy parsing. When performing JSON encoding, `windapsearch` will automatically convert certain LDAP attributes to a more human friendly format as well (e.g. timestamps, GUIDs, enumerations, etc)

This is a complete re-write of my earlier [Python implementation](https://github.com/ropnop/windapsearch). For some more background/explanation on how I'm using Go, and more advanced usage examples see [this blog post](TODO).

## Installation
You can download pre-compiled binaries for amd64 Linux/Mac/Windows from the [latest releases](https://github.com/ropnop/go-windapsearch/releases)

To build from source, I use [mage](https://github.com/magefile/mage), a Make like tool written in Go. Install `mage` then run the mage targets:

```
$ git clone https://github.com/ropnop/go-windapsearch.git && cd go-windapsearch
$ go get github.com/magefile/mage
$ mage
Targets:
  build    Compile windapsearch for current OS and ARCH
  clean    Delete bin and dist dirs
  dist     Cross-compile for Windows, Linux, Mac x64 and put in ./dist
$ mage build
$ ./windapsearch --version
```

# Usage
`windapsearch` is a standalone binary with multiple modules for various common LDAP queries

```
windapsearch: a tool to perform Windows domain enumeration through LDAP queries
Version: dev (131fd6d) | Built: 06/10/20 (go1.14.3) | Ronnie Flathers @ropnop

Usage: ./windapsearch [options] -m [module] [module options]

Options:
  -d, --domain string     The FQDN of the domain (e.g. 'lab.example.com'). Only needed if dc not provided
      --dc string         The Domain Controller to query against
  -u, --username string   The full username with domain to bind with (e.g. 'ropnop@lab.example.com' or 'LAB\ropnop')
                           If not specified, will attempt anonymous bind
  -p, --password string   Password to use. If not specified, will be prompted for
      --hash string       NTLM Hash to use instead of password (i.e. pass-the-hash)
      --ntlm              Use NTLM auth (automatic if hash is set)
      --port int          Port to connect to (if non standard)
      --secure            Use LDAPS. This will not verify TLS certs, however. (default: false)
      --full              Output all attributes from LDAP
  -o, --output string     Save results to file
  -j, --json              Convert LDAP output to JSON
      --page-size int     LDAP page size to use (default 1000)
      --version           Show version info and exit
  -v, --verbose           Show info logs
      --debug             Show debug logs
  -h, --help              Show this help
  -m, --module string     Module to use

Available modules:
    admin-objects       Enumerate all objects with protected ACLs (i.e admins)
    computers           Enumerate AD Computers
    custom              Run a custom LDAP syntax filter
    domain-admins       Recursively list all users objects in Domain Admins group
    gpos                Enumerate Group Policy Objects
    groups              List all AD groups
    members             Query for members of a group
    metadata            Print LDAP server metadata
    privileged-users    Recursively list members of all highly privileged groups
    search              Perform an ANR Search and return the results
    unconstrained       Find objects that allow unconstrained delegation
    user-spns           Enumerate all users objects with Service Principal Names (for kerberoasting)
    users               List all user objects
```

## Selecting a Module
Select a module to use with the `-m` option. Some modules have additional options which can be seen by specifying a module when running `-h`:

```
$ ./windapsearch -m users -h
<...>
Options for "users" module:
      --attrs strings   Comma separated custom atrributes to display (default [cn,sAMAccountName])
      --filter string   Extra LDAP syntax filter to use
  -s, --search string   Search term to filter on
```

Each module defines a default set of attributes to return. These can always be overriden by the comma separated `--attrs` option, or by specifying `--full`, which will always return every attribute.


## Output Formats
With no other options specified, `windapsearch` will display output to the terminal in the same text based format used by `ldapsearch`. Output can also be written to a file by specifying the `-o` option.

When specifying the `-j` option, the tool will convert LDAP responses to JSON format and outut a JSON array of LDAP entries. During the marshalling, `windapsearch` will also convert binary values to human-readable formats, and perform some enumeration substitution with string values.

For example, when looking a single user, these are the normal "text" attributes:
```
whenCreated: 20170806185838.0Z
objectSid: AQUAAAAAAAUVAAAAoWuXYvBp2/Bf49rCUgQAAA==
lastLogonTimestamp: 132340658159483754
userAccountControl: 66048
```

But in JSON format, they are converted:
```json
    "whenCreated": "2017-08-06T18:58:38Z",
    "objectSid": "S-1-5-21-1654090657-4040911344-3269124959-1106",
    "lastLogonTimestamp": "2020-05-15T20:23:35.9483754-05:00",
    "userAccountControl": [
      "DONT_EXPIRE_PASSWORD",
      "NORMAL_ACCOUNT"
    ],
```

*Note: I have not implemented full mapping/pretty printing of every LDAP attribute. If you see one that should be converted to something else and isn't, please open an Issue - or better yet a PR ;)*

## Logging
To see more information, including the full LDAP queries that are being sent, use the `--verbose` option, which will display helpful information.

If you are experiencing issues, please use the `--debug` option for much more detailed log information, including every entry being parsed

# Credits
 - The authors of [go-ldap](https://github.com/go-ldap/ldap) for the LDAP client that powers all of this
 - [audibleblink](https://twitter.com/4lex) for the [idea](https://twitter.com/4lex/status/1254037754842931200?s=20) and the [package](github.com/audibleblink/msldapuac) to parse UserAccountControl from LDAP

