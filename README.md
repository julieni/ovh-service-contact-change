# OVH Service contact change

## help

run `./ovh-service-contact-change -h`

## first use - setup

run `./ovh-service-contact-change init`

setup your own API credentials and note them down

you can report them in your .env file (copied from .env.example file)

## batch accept the contact change requests

place your mails as text files in a folder (by default, subfolder "mails" in the current directory)

run `./ovh-service-contact-change accept`