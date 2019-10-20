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