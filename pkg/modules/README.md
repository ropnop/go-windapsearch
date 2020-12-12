# Windapsearch Modules
`windapsearch` has adopted a modular structure for extending and adding common LDAP enumeration techniques.
The core functionality of setting up the connection and processing results is handled by the `windapsearch` and `ldapsession`
packages, so modules can be fairly standalone.

The following modules have been implemented, with functionality copied from the existing Python `windapsearch` script:

 * [admin-objects](#admin-objects)
 * [computers](#computers)
 * [custom](#custom)
 * [domain-admins](#domain-admins)
 * [gpos](#gpos)
 * [groups](#groups)
 * [members](#members)
 * [metadata](#metadata)
 * [privileged-users](#privileged-users)
 * [search](#search)
 * [unconstrained](#unconstrained)
 * [user-spns](#user-spns)
 * [users](#users)

**Common Options**
Every module inherits/hones the following command line switches:
`--attrs`: custom comma separated attributes to display. Overrides per-module defaults
`--full`: display all attributes (`*`). Overrides defaults and `--attrs`
`--json/-j`: Convert entries to JSON and convert availble fields to friendly formats

Also, `dn` will always be included as an attribute by default since it is always returned in responses.

## admin-objects
**Description**: `Enumerate all objects with protected ACLs (i.e admins)`

**Default Attrs**: `cn`

**Base Filter**: `(adminCount=1)`

**Additional Options**: ``

This module searches for any and all LDAP entries that have `adminCount=1`, indicating that they have protected ACLs, which means they are highly privileged objects (e.g. Domain Admins, priveleged groups, etc)

**Example Usage**:
```
$ ./windapsearch -d lab.ropnop.com -u agreen@lab.ropnop.com -p $PASS -m admin-objects -j | jq '.[0]'
{
  "cn": "Backup Operators",
  "dn": "CN=Backup Operators,CN=Builtin,DC=lab,DC=ropnop,DC=com"
}
```

## computers
**Description**: `Enumerate AD Computers`

**Default Attrs**: `cn, dNSHostName, operatingSystem, operatingSystemVersion,operatingSystemServicePack`

**Base Filter**: `(objectClass=Computer)`

**Additional Options**: ``

This module searches for all AD joined computers, and displays LDAP information about the computers, including DNS name and OS version.

**Example Usage**:
```
$ ./windapsearch -d lab.ropnop.com -u agreen@lab.ropnop.com -p $PASS -m computers -j | jq '.[0]'
{
  "cn": "WS03WIN10",
  "dNSHostName": "ws03win10.lab.ropnop.com",
  "dn": "CN=WS03WIN10,OU=computers,OU=LAB,DC=lab,DC=ropnop,DC=com",
  "operatingSystem": "Windows 10 Pro",
  "operatingSystemVersion": "10.0 (17134)"
}
```

## custom
**Description**: `Run a custom LDAP syntax filter`

**Default Attrs**: `*`

**Base Filter**: `custom`

**Additional Options**: `--filter`

The module lets you specify a custom LDAP syntax filter to run, and returns all attributes by default. *Note: your filter must be valid LDAP filter syntax and wrapped in parantheses*

**Example Usage**: 
```
$ ./bin/windapsearch -d lab.ropnop.com -u agreen@lab.ropnop.com -p $PASS -m custom --filter "(sAMAccountName=thoffman)" --attrs pwdLastSet -j | jq .
[
  {
    "dn": "CN=Trevor Hoffman,OU=users,OU=LAB,DC=lab,DC=ropnop,DC=com",
    "pwdLastSet": "2019-05-16T18:39:48.1597266-05:00"
  }
]
```

## dns-names
**Description**: `Query AD integrated DNS for domain names`

**Default Attrs**: `name, dnsTombstoned`

**Base Filter**: `(objectClass=*)`

**Additional Options**: ``

The module queries the Active Directory integrated DNS and returns all objects. Unfortunately, there is no attribute for the complete FQDN, therefore the dn is returned, containing sufficient information to recover the actual FQDN.

**Example Usage**: 
```
$ ./bin/windapsearch -d lab.ropnop.com -u agreen@lab.ropnop.com -p $PASS -m dns-names -j | jq .
[
  {
    "dn": "DC=sharepoint,DC=lab.ropnop.com,CN=MicrosoftDNS,DC=DomainDnsZones,DC=lab,DC=ropnop,DC=com"
  },
  {
    "dn": "DC=dc,DC=lab.ropnop.com,CN=MicrosoftDNS,DC=DomainDnsZones,DC=lab,DC=ropnop,DC=com"
  },
  {
    "dn": "DC=app01,DC=lab.ropnop.com,CN=MicrosoftDNS,DC=DomainDnsZones,DC=lab,DC=ropnop,DC=com"
  },
  ...
]

```

## dns-zones
**Description**: `Query AD integrated DNS for registered zones`

**Default Attrs**: `dn, name`

**Base Filter**: `(&(objectClass=dnsZone)(!name=RootDNSServers)(!name=*.in-addr.arpa)(!name=_msdcs.*)(!name=..TrustAnchors))`

**Additional Options**: ``

The module queries the Active Directory integrated DNS and returns all DNS zones. Unfortunately, there is no attribute for the complete FQDN, therefore the dn is returned, containing sufficient information to recover the actual FQDN.

**Example Usage**: 
```
$ ./bin/windapsearch -d lab.ropnop.com -u agreen@lab.ropnop.com -p $PASS -m dns-zones -j | jq .
[
  {
    "dn": "DC=lab.ropnop.com,CN=MicrosoftDNS,DC=DomainDnsZones,DC=lab,DC=ropnop,DC=com",
    "name": "lab.ropnop.com"
  },
  {
    "dn": "DC=dev.ropnop.net,CN=MicrosoftDNS,DC=DomainDnsZones,DC=dev,DC=ropnop,DC=net",
    "name": "dev.ropnop.net"
  }
]
```

## domain-admins
**Description**: `Recursively list all users objects in Domain Admins group`

**Default Attrs**: `cn, sAMAccountName`

**Base Filter**: `(&(objectClass=user)(|(memberof:1.2.840.113556.1.4.1941:=CN=Domain Admins,CN=Users,DC=lab,DC=ropnop,DC=com)(memberof:1.2.840.113556.1.4.1941:=CN=Domain-Admins,CN=Users,DC=lab,DC=ropnop,DC=com)(memberof:1.2.840.113556.1.4.1941:=CN=Domain Administrators,CN=Users,DC=lab,DC=ropnop,DC=com)(memberof:1.2.840.113556.1.4.1941:=CN=Domain-Administrators,CN=Users,DC=lab,DC=ropnop,DC=com)(memberof:1.2.840.113556.1.4.1941:=CN=Domänen Admins,CN=Users,DC=lab,DC=ropnop,DC=com)(memberof:1.2.840.113556.1.4.1941:=CN=Domänen-Admins,CN=Users,DC=lab,DC=ropnop,DC=com)(memberof:1.2.840.113556.1.4.1941:=CN=Domain Admins,CN=Users,DC=lab,DC=ropnop,DC=com)(memberof:1.2.840.113556.1.4.1941:=CN=Domain-Admins,CN=Users,DC=lab,DC=ropnop,DC=com)(memberof:1.2.840.113556.1.4.1941:=CN=Domänen Administratoren,CN=Users,DC=lab,DC=ropnop,DC=com)(memberof:1.2.840.113556.1.4.1941:=CN=Domänen-Administratoren,CN=Users,DC=lab,DC=ropnop,DC=com)))`

**Additional Options**: ``

This module lists every user object that is a member of the Domain Admins group. It performs recursive lookups using the OID `LDAP_MATCHING_RULE_IN_CHAIN`, so it will also display any user that is transitively part of the Domain Admins group as well. The filter includes language variations of the `Domain Admins` group too.

**Example Usage**:
```
$ ./windapsearch -d lab.ropnop.com -u agreen@lab.ropnop.com -p $PASS -m domain-admins -j | jq '.[0]'
{
  "cn": "Edna Dominguez",
  "dn": "CN=Edna Dominguez,OU=US,OU=users,OU=LAB,DC=lab,DC=ropnop,DC=com",
  "sAMAccountName": "edominguez"
}
```

## gpos
**Description**: `Enumerate Group Policy Objects`

**Default Attrs**: `displayName, gPCFileSysPath`

**Base Filter**: `(objectClass=groupPolicyContainer)`

**Additional Options**: ``

This module lists Group Policy Objects found in LDAP. It will display the display name and SYSVOL path by default:

**Example Usage**:
```
$ ./windapsearch -d lab.ropnop.com -u agreen@lab.ropnop.com -p $PASS -m gpos -j | jq '.[0]'
{
  "displayName": "firewall_rules",
  "dn": "CN={24722667-432E-4508-A58C-15D3D42FEFF4},CN=Policies,CN=System,DC=lab,DC=ropnop,DC=com",
  "gPCFileSysPath": "\\\\lab.ropnop.com\\SysVol\\lab.ropnop.com\\Policies\\{24722667-432E-4508-A58C-15D3D42FEFF4}"
}
```

## groups
**Description**: `List all AD groups`

**Default Attrs**: `cn`

**Base Filter**: `(objectcategory=group)`

**Additional Options**: `-s / --search`

This module lists all group objects. By default it only displays the CN. Optionally, it takes a `search` option to narrow down groups.

```
$ ./windapsearch -d lab.ropnop.com -u agreen@lab.ropnop.com -p $PASS -m groups --search IT --attrs cn,member -j | jq '.[0]'
{
  "cn": "IT Admins",
  "dn": "CN=IT Admins,OU=groups,OU=LAB,DC=lab,DC=ropnop,DC=com",
  "member": [
    "CN=vulnscanner,OU=service-accounts,OU=LAB,DC=lab,DC=ropnop,DC=com",
    "CN=Desktop Support,OU=groups,OU=LAB,DC=lab,DC=ropnop,DC=com",
    "CN=Mark Murdock,OU=US,OU=users,OU=LAB,DC=lab,DC=ropnop,DC=com",
    "CN=Susan Hendrickson,OU=US,OU=users,OU=LAB,DC=lab,DC=ropnop,DC=com",
    "CN=Michael Timpson,OU=US,OU=users,OU=LAB,DC=lab,DC=ropnop,DC=com",
    "CN=Herbert Smith,OU=US,OU=users,OU=LAB,DC=lab,DC=ropnop,DC=com",
    "CN=Paul Rivera,OU=US,OU=users,OU=LAB,DC=lab,DC=ropnop,DC=com"
  ]
}
```

## members
**Description**: `Query for members of a group`

**Default Attrs**: `cn, sAMAccountName`

**Base Filter**: `(memberOf=<group_dn>)`

**Additional Options**: `-g / --group, -r / --recursive, -s / --search, --users`

This module lists members of a group. You must specify a group with `-g` by its full distinguished name, or perform a search with `-s`. If more than one match is found, the module will prompt you for which group you meant.

Optionally, you can perform a `--recursive` lookup to list transitive members as well, or limit the results to only user objects with `--users`

**Example Usage**:
```
$ ./windapsearch -d lab.ropnop.com -u agreen@lab.ropnop.com -p $PASS -m members -s remote -j |jq '.[0]'
What DN do you want to use?

1. CN=Remote Desktop Users,CN=Builtin,DC=lab,DC=ropnop,DC=com
2. CN=Remote Management Users,CN=Builtin,DC=lab,DC=ropnop,DC=com

Enter a number: 1

[+] Using group: CN=Remote Desktop Users,CN=Builtin,DC=lab,DC=ropnop,DC=com

{
  "cn": "Peter Harris",
  "dn": "CN=Peter Harris,OU=US,OU=users,OU=LAB,DC=lab,DC=ropnop,DC=com",
  "sAMAccountName": "pharris"
}
```

## metadata
**Description**: `Print LDAP server metadata`

**Default Attrs**: `defaultNamingContext, domainFunctionality, forestFunctionality, domainControllerFunctionality, dnsHostName`

**Base Filter**: `(objectClass=*)`

**Additional Options**: ``

This module queries the LDAP server for metadata. It does not require an authenticated bind. By default it returns functionality levels, base DN, and DNS info

**Example Usage**:
```
$ ./windapsearch -d lab.ropnop.com -m metadata -j | jq .
[
  {
    "defaultNamingContext": "DC=lab,DC=ropnop,DC=com",
    "dnsHostName": "pdc01.lab.ropnop.com",
    "domainControllerFunctionality": "2012 R2",
    "domainFunctionality": "2012 R2",
    "forestFunctionality": "2012 R2"
  }
]
```

## privileged-users
**Description**: `Recursively list members of all highly privileged groups`

**Default Attrs**: `cn, sAMAccountName`

**Base Filter**: `(&(objectClass=user)(|(memberof:1.2.840.113556.1.4.1941:=CN=Administrators,CN=Users,DC=lab,DC=ropnop,DC=com)(memberof:1.2.840.113556.1.4.1941:=CN=Enterprise Admins,CN=Users,DC=lab,DC=ropnop,DC=com)(memberof:1.2.840.113556.1.4.1941:=CN=Schema Admins,CN=Users,DC=lab,DC=ropnop,DC=com)(memberof:1.2.840.113556.1.4.1941:=CN=Account Operators,CN=Users,DC=lab,DC=ropnop,DC=com)(memberof:1.2.840.113556.1.4.1941:=CN=Backup Operators,CN=Users,DC=lab,DC=ropnop,DC=com)(memberof:1.2.840.113556.1.4.1941:=CN=Server Management,CN=Users,DC=lab,DC=ropnop,DC=com)(memberof:1.2.840.113556.1.4.1941:=CN=Konten-Operatoren,CN=Users,DC=lab,DC=ropnop,DC=com)(memberof:1.2.840.113556.1.4.1941:=CN=Sicherungs-Operatoren,CN=Users,DC=lab,DC=ropnop,DC=com)(memberof:1.2.840.113556.1.4.1941:=CN=Server-Operatoren,CN=Users,DC=lab,DC=ropnop,DC=com)(memberof:1.2.840.113556.1.4.1941:=CN=Schema-Admins,CN=Users,DC=lab,DC=ropnop,DC=com)(memberof:1.2.840.113556.1.4.1941:=CN=Domain Admins,CN=Users,DC=lab,DC=ropnop,DC=com)(memberof:1.2.840.113556.1.4.1941:=CN=Domain-Admins,CN=Users,DC=lab,DC=ropnop,DC=com)(memberof:1.2.840.113556.1.4.1941:=CN=Domain Administrators,CN=Users,DC=lab,DC=ropnop,DC=com)(memberof:1.2.840.113556.1.4.1941:=CN=Domain-Administrators,CN=Users,DC=lab,DC=ropnop,DC=com)(memberof:1.2.840.113556.1.4.1941:=CN=Domänen Admins,CN=Users,DC=lab,DC=ropnop,DC=com)(memberof:1.2.840.113556.1.4.1941:=CN=Domänen-Admins,CN=Users,DC=lab,DC=ropnop,DC=com)(memberof:1.2.840.113556.1.4.1941:=CN=Domain Admins,CN=Users,DC=lab,DC=ropnop,DC=com)(memberof:1.2.840.113556.1.4.1941:=CN=Domain-Admins,CN=Users,DC=lab,DC=ropnop,DC=com)(memberof:1.2.840.113556.1.4.1941:=CN=Domänen Administratoren,CN=Users,DC=lab,DC=ropnop,DC=com)(memberof:1.2.840.113556.1.4.1941:=CN=Domänen-Administratoren,CN=Users,DC=lab,DC=ropnop,DC=com)))`

**Additional Options**: ``

This module recursively lists all members of Domain Admins and every other highly privileged group (e.g. `Schema Admins`, `Backup Operators`, etc)

**Example Usage**:
```
$ ./windapsearch -d lab.ropnop.com -u agreen@lab.ropnop.com -p $PASS -m privileged-users -j | jq '.[0]'
{
  "cn": "Paul Rivera",
  "dn": "CN=Paul Rivera,OU=US,OU=users,OU=LAB,DC=lab,DC=ropnop,DC=com",
  "sAMAccountName": "privera"
}
```

## search
**Description**: `Perform an ANR Search and return the results`

**Default Attrs**: `*`

**Base Filter**: `(anr=<search_term>)`

**Additional Options**: `--all, -s / --search`

This module performs an [Ambiguous Name Resolution](https://docs.microsoft.com/en-us/openspecs/windows_protocols/ms-adts/1a9177f4-0272-4ab8-aa22-3c3eafd39e4b) search. If more than one match is found, it will prompt you for the entry you wish to retrieve. If `--all` is specified, it will not prompt and instead dump every matching entry.

**Example Usage**:
```
$ ./windapsearch -d lab.ropnop.com -u agreen@lab.ropnop.com -p $PASS -m search -s ronnie --attrs objectSid -j | jq '.[0]'
What DN do you want to use?

1. CN=Ronnie Weinberg,OU=US,OU=users,OU=LAB,DC=lab,DC=ropnop,DC=com
2. CN=Ronnie Cooper,OU=US,OU=users,OU=LAB,DC=lab,DC=ropnop,DC=com

Enter a number: 1

{
  "dn": "CN=Ronnie Weinberg,OU=US,OU=users,OU=LAB,DC=lab,DC=ropnop,DC=com",
  "objectSid": "S-1-5-21-1654090657-4040911344-3269124959-1715"
}
```

## unconstrained
**Description**: `Find objects that allow unconstrained delegation`

**Default Attrs**: `cn, sAMAccountName`

**Base Filter**: `(userAccountControl:1.2.840.113556.1.4.803:=524288)`

**Additional Options**: `--computers, --users`

This module will search for LDAP objects that allow for unconstrained delegation. By default it will list all objects, though you can limit it either computers or users by using `--computers` or `--users`, respectively.

**Example Usage**:
```
$ ./windapsearch -d lab.ropnop.com -u agreen@lab.ropnop.com -p $PASS -m unconstrained -j | jq '.[0]'
{
  "cn": "PDC01",
  "dn": "CN=PDC01,OU=Domain Controllers,DC=lab,DC=ropnop,DC=com",
  "sAMAccountName": "PDC01$"
}
```

## user-spns
**Description**: `Enumerate all users objects with Service Principal Names (for kerberoasting)`

**Default Attrs**: `cn, sAMAccountName, servicePrincipalName`

**Base Filter**: `(&(&(servicePrincipalName=*)(UserAccountControl:1.2.840.113556.1.4.803:=512))(!(UserAccountControl:1.2.840.113556.1.4.803:=2)))`

**Additional Options**: ``

This module will identify user objects with servicePrincipalNames defined, which can be used for kerberoasting.

**Example Usage**:
```
$ ./windapsearch -d lab.ropnop.com -u agreen@lab.ropnop.com -p $PASS -m user-spns -j | jq '.[0]'
{
  "cn": "vulnscanner",
  "dn": "CN=vulnscanner,OU=service-accounts,OU=LAB,DC=lab,DC=ropnop,DC=com",
  "sAMAccountName": "vulnscanner",
  "servicePrincipalName": [
    "HTTP/webdev.lab.ropnop.com"
  ]
}
```

## users
**Description**: `List all user objects`

**Default Attrs**: `cn, sAMAccountName, userPrincipalName`

**Base Filter**: `(objectcategory=user)`

**Additional Options**: `--filter, -s / --search`

This module lists every LDAP user object. Depending on the size of the domain, this can get very big. You can limit results by adding an additional LDAP syntax filter with `--filter`, or an ANR search term with `--search`.

**Example Usage**:
```
$ ./windapsearch -d lab.ropnop.com -u agreen@lab.ropnop.com -p $PASS -m users --full -j -o full_user_dump.json
[+] full_user_dump.json written

$ jq '.|length' full_user_dump.json
2758

$ jq '.[] | select(.badPwdCount > 4)|.userPrincipalName' full_user_dump.json
"baguirre@lab.ropnop.com"
"avenezia@lab.ropnop.com"
"aturner@lab.ropnop.com"
"avelasquez@lab.ropnop.com"
"awoodell@lab.ropnop.com"
"ayoho@lab.ropnop.com"
"awilliams@lab.ropnop.com"
"atressler@lab.ropnop.com"
"awhite@lab.ropnop.com"
"awoods@lab.ropnop.com"
"ayim@lab.ropnop.com"
"aweiss@lab.ropnop.com"
"ayunker@lab.ropnop.com"
"barndt@lab.ropnop.com"
```








