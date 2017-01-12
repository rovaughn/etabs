etabs
=====

Simple elastic tabs formatter.

It takes in a file and automatically aligns the elastic tabs.  Leading
whitespace is completely ignored.  I'm still unsure what the best way is to
mark an elastic tab in a file.  Currently it's marked as a TAB character or as
two or more spaces (excluding leading whitespace).

As an example, let's say the following is `file.c`.

	struct A {
		int        a;  // field 1
		int     b;   // field 2
		char      *c;      // field 3
	};

	struct X {
		struct D     x;  // field 1
		const char    *y;   // field 2
		char    *z;            // field 3
	};

When you run `etabs file.c` it becomes

	struct A {
		int   a;   // field 1
		int   b;   // field 2
		char  *c;  // field 3
	};

	struct X {
		struct D    x;   // field 1
		const char  *y;  // field 2
		char        *z;  // field 3
	};

in general it's hard to make an example that works for all syntaxes but this use
of elastic tabs probably covers 99% of use cases.  The downsides are for
instance in the above example you might want the variable names aligned
regardless of the asterisk.  Also elastic tabs are recorded as at least two
spaces, otherwise it might interpret the `D` in `struct D` or the `char` in
`const char` as needing to be aligned with the variable names.

The upside is it doesn't really require any editor support or weird unicode
characters to implement elastic tabstops in a way that works with any language.

TODO
----

- Tests
- Benchmarks.  There are some potential efficiency gains although I don't know
  if they'd be noticeable.

