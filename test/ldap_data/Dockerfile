FROM alpine:3.18

COPY ./entrypoint /entrypoint
RUN chmod 755 /entrypoint

# Install openldap and dependencies
RUN apk update
RUN apk add openldap openldap-back-mdb openldap-clients openldap-overlay-all openssl

# Create ldap files/directories
RUN rm -f /etc/openldap/slapd.conf
RUN mkdir /etc/openldap/slapd.d
RUN mkdir /var/lib/openldap/run

# File to keep track of which bootstrap scripts that have been run
RUN touch /var/lib/bootstrap.txt

# Create ldaps certificates
RUN mkdir /etc/openldap/certs
WORKDIR /etc/openldap/certs
RUN openssl req -x509 -new -nodes -newkey rsa:2048 -keyout ca.key -out ca.crt -batch -subj '/OU=tests/O=dory/' -days 3650
RUN openssl req -new -nodes -newkey rsa:2048 -keyout ldaps.key -out ldaps.csr -days 3650 -subj '/OU=tests/O=dory/CN=localhost.priv/'
RUN openssl x509 -req -in ldaps.csr -CA ca.crt -CAkey ca.key -CAcreateserial -out ldaps.crt -days 3650
RUN chmod -R 755 .
WORKDIR /

# LDIF scripts ran to boostrap the ldap server
VOLUME /bootstrap

EXPOSE 389 636
ENTRYPOINT ["/entrypoint", "-dstats"]
