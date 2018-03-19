# TODO

A list of things to do, once the core language is implemented. In no particular order.

 - Booleang -> circuit diagram, using Graphviz
 - AST renderer with Graphviz
 - A web interface
 - Make a cool CPU example program
 - Make standard libraries
 - Extend language:
   - Constants
     - Like C macros
     - `#delay 5s`
	 - `#and (a & b)`
   - Number literals
     - `%five i64: 1205`
	 - `%num   u5: 57`
	 - `%flp  f32: 102.5`
   - String literals
     - UTF-8 encoded
     - `%name "hello 🌍"`
   - Macro init values
     - `%num (a:1, b:0, c: 1, d: 1)`