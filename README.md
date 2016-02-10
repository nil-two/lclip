lclip
=====

[![Build Status](https://travis-ci.org/kusabashira/lclip.svg?branch=master)](https://travis-ci.org/kusabashira/lclip)

CLI interface for labeled clipboard.

```
$ echo -e "Hello\nworld" | lclip -s hello
$ lclip -g hello
Hello
world
```

Usage
-----

```
$ lclip <operation> [...]
operations:
  lclip {-h --help}                      # show this help message
  lclip {-v --version}                   # print the version
  lclip {-l --labels}                    # list labels
  lclip {-g --get}     <label>           # paste text from label
  lclip {-s --set}     <label> [file(s)] # copy text to label
  lclip {-d --delete}  <label(s)>        # delete label(s)
```

Installation
------------

### compiled binary

See [releases](https://github.com/kusabashira/lclip/releases)

### go get

```
go get github.com/kusabashira/lclip
```

Operations
----------

### -h, --help

Display the help message.

### -v, --version

Display the version of lclip.

### -l, --labels

List sorted labels.

### -g, --get *label*

Paste text from label.
If label doesn't exist, it paste a newline.

### -s, --set *label* *file(s)*

Copy text to label.
Read from file(s), or standard input.

### -d, --delete *label(s)*

Delete label(s) from the storage.

Environment Variables
---------------------

### LCLIP\_PATH

This variable specifies the storage path.
The default is `~/.lclip.json`.

Other Specification
-------------------

- Default storage path is `~/.lclip.json`.
  - You can change the storage path by setting `$LCLIP_PATH`.
- Create storage if it doesn't exist.

License
-------

MIT License

Author
------

kusabashira <kusabashira227@gmail.com>
