# LDAP server

The ldap server used when running tests is an OpenLDAP server.

Its configuration is built from the different LDIF scripts in the `bootstrap`
directory. Those scripts are ordered by name when executed.

We use comments to define how the LDIF scripts should be run:
 - `#type: XXX` defines the executable to use to run the script.
   `#type: add` will use `slapadd` and `#type: modify` will use slapmodify.
 - `#db: XXX` defines which database to use. Use `#db: cn=config` to edit the
   OpenLDAP configuration, and `#db: dc=localhost,dc=priv` to edit the main
   database.
