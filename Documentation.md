This file contains documentation related to the development
of LaME. It is intended to be a broad overview of the
project's architecture.

## Package Stories
**engine** depends on **target** so that it can
  make a program blueprint that a code generator
  can understand.

**engine** depends on **model** so that it can
  process models to make decisions about the
  program blueprint it's going to generate.

**generators** depend on **target** so that they can
  understand a universal program blueprint definition,
  as well as provide an interface for invocation that
  the engine can understand.

**generators** depend on **support** so that they can
  implement common code generation logic with minimal
  code.

## How to handle scoping

First a disclaimer. We all have experience with different
programming languages and may have slightly different
definitions for even some of the most fundamendal terms
in computer science. With this in mind, if any definition
here seems inaccurate please address this by adding an
issue to the repository. With that said, it is my beleif
that the concepts described regarding the implementation
of LaME will result in a system capable of generating code
in all but the most esoteric of programming languages.

### Block Scope vs Function Scope
Firstly, block scope vs function scope doesn't matter.
LaME will internally be based on function scope with
the C89 rule for variable placement, which greatly
simplifies code generation. However, a LisPI syntax
frontend could easily implement block scope under these
constraints by manipulating variable names or adding
logic for variable access control. Similarly, a code
generator could optimize variable declarations based
on what is most efficient for the target language, albiet
this may be more challenging; this is irrelevant to
functionality however, and the algorithms used in a
program will always have precedence in determining the
program's performance, so performance optimizations with
variable scope should be implemented on an as-needed basis.

### Lexical Scope vs Dynamic Scope
Writing this is on the TODO, but the short version is
LaME will be implemented with lexical scope. For languages
with dynamic scope, certain variables preceeding a
function call will be considered that function's arguments
and the code generator will need to add logic to replace
names that may collide.

## How to handle implementation variables
It may be necessary for a code generator to create its own
variables in order to build a program that functions in the
way LaME expects it to. To avoid collisions, single
underscores in modelled variables can be escaped as
double underscores. This should be handled at code
generation in case language quirks result. For instance,
in a language that only supports pure-alpha identifiers,
underscores could be escaped as zu and the letter 'z'
could be escaped as zz. Access to properties on objects
outside of LaME would require special handling.