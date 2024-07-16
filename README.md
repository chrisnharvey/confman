# Confman
Linux/macOS configuration/dotfile manager

## Disclaimer
- ⚠️ This project is under active development.
- ⚠️ Do not rely on this unless you know what you are doing and have backups.
- ⚠️ The CLI and the repository format will change
- ⚠️ The code is (currenly) a mess, I just wanted to get something working. I'll tidy it up at some point.
- ⚠️ There are no docs (yet).
- ⚠️ You have been warned!

## How it works

Confman works from a directory on your system (this is called the "repository"). Config/dot files
are moved into this directory when you run `confman add /etc/myconffile.conf myconffile`. A symlink
is then added where the conf file was which points to the file in the Confman repository. An entry
is also added to the `.confman.yaml` file which maps the file in the Confman repository to the destination
on your system.

You can commit the Confman repository a VCS repo, or sync it to a cloud storage provider.

On a new system, you can run `confman link` to setup all of the symlinks from your Confman repo to
their destinations on your system.

## Confman Repository

The confman repository is where your config files will be kept, aloing with a `.confman.yaml` yaml file which keeps track of the mappings between files in the repository and their destination on your system.

Example:
```yaml
paths:
    /etc/pam.d/kde-fingerprint: kde-fingerprint
    /home/chris/code/confman/hi: foo2
    /home/chris/code/confman/hi2: foo3
```

## CLI

```
Available Commands:
  add         Add a configuration file to the Confman repository [source] [destination]
  link        Creates symlinks for all configuration files managed by Confman
  list        Lists all configuration files and their status
```
