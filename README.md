# OVH Service contact change

## help

run `./ovh-service-contact-change -h`

## first use - setup

run `./ovh-service-contact-change init`

setup your own API credentials and note them down

you can report them in your .env file (copied from .env.example file)

## batch accept the contact change requests

### IMAP

run `./ovh-service-contact-change imap --imap-server imap.domain.tld --imap-login me@domain.tld`

by default, TLS on port 993 will be used, and only mails received in the last day will be searched for.

run `./ovh-service-contact-change imap -h` for options

### local files

place your mails as text files in a folder (by default, subfolder "mails" in the current directory)

run `./ovh-service-contact-change file`

run `./ovh-service-contact-change file -h` for options