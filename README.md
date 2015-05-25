lclip
=====

cli interface for labeled clipboard.

	$ echo -e "Hello\nworld" | lclip -s hello
	$ lclip -g hello
	Hello
	world


Usage
-----

	$ lclip <operation> [...]
	operations:
		lclip {-h --help}                      # Show this help message
		lclip {-v --version}                   # Print the version
		lclip {-l --labels}                    # List labels
		lclip {-g --get}     <label>           # Paste text from label
		lclip {-s --set}     <label> [file(s)] # Copy text to label
		lclip {-d --delete}  <label(s)>        # Delete label(s)


Installation
------------

###compiled binary

See [releases](https://github.com/kusabashira/lclip/releases)


###go get

	go get github.com/kusabashira/lclip/cmd/lclip


Operations
----------

### -h, --help

Display a help message.


### -v, --version

Display the version of lclip.


### -l, --labels

List sorted labels.


### -g, --get *label*

Paste text from label.
If not exists label, paste a newline.


### -s, --set *label* *file(s)*

Copy text to label.
Read from file(s), or standard input, and nothing output.


### -d, --delete *label(s)*

Delete label(s) from the storage.


Other Specification
-------------------

- Storage path is `~/.lclip.db`
- If not exist exists storage, storage will be created.


License
-------

MIT License


Author
------

wara <kusabashira227@gmail.com>
