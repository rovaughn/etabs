etabs
=====

Simple [elastic tabs](http://nickgravgaard.com/elastic-tabstops/) formatter.  A
one-size-fits-all syntax formatter.

	etabs file.c
	cat file.c | etabs - | less

etabs takes in a file and automatically aligns the elastic tabs.  An elastic
tabstop is, excluding leading whitespace, any sequence of whitespace that isn't
just a single space.

As an example, let's say the following is `file.c`.

	struct A {
		int        a;  // field 1
		int     b;   // field 2
		char*      c;      // field 3
	};

	struct X {
		struct D     x;  // field 1
		const char*    y;   // field 2
		char*    z;            // field 3
	};

When you run `etabs file.c` it becomes

	struct A {
		int    a;  // field 1
		int    b;  // field 2
		char*  c;  // field 3
	};

	struct X {
		struct D     x;  // field 1
		const char*  y;  // field 2
		char*        z;  // field 3
	};

The reason elastic tabs have to be recorded as two or more spaces are because
otherwise it might interpret the `D` in `struct D` or the `char` in `const char`
as needing to be aligned with the variable names.  Also the only time multiple
spaces are generally used is specifically when aligning something.  TAB
characters that aren't leading whitespace are also elastic tabs, allowing you to
just hit TAB to insert an elastic tab without worrying about alignment.

The upside is it doesn't really require any editor support or weird unicode
characters to implement elastic tabstops in a way that works with any language.

vim
---

You can have etabs automatically run when saving a file with

	autocmd BufWritePost * silent !etabs <afile>

You can replace `*` with a pattern to constrain it to only files with a certain
extension, e.g.

	autocmd BufWritePost *.c,*.h silent !etabs <afile>

to only have it run when saving a C file.

building
--------

Building this program depends on [Go](https://golang.org/doc/install) being
installed.  After it is installed the following can build it:

	export GOPATH=~/go/src
	mkdir -p ~/go/src
	cd ~/go/src
	git clone https://github.com/rovaughn/etabs
	cd etabs
	go install

This will put etabs at go/bin/etabs, so then you can just add go/bin to your
path or move the binary somewhere visible on your path.

TODO
----

- Possibly more efficiency gains.
- Maybe more checking for alignment issues.  E.g. trailing spaces, or when using
  spaces to indent, using an inconsistent multiple of spaces, or mixing
  tabs/spaces.  Also it could convert all indentation to a given format.
  Probably more suitable to the language-specific formatters though.

