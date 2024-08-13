# ClearGO
- This is my "Go" at building an interpreter for the first time. This follows Thorsten Ball's "Writing an Interpreter in Go" very closely, if not exactly.
- My end goal is to create a compiler, which will follow along with the sequel to this book.
- Project goals:
    - Create a working interpreter with the following functions (as per the book):
        - Variables
        - Data types:
            - Int
            - Bool
            - String
            - Array
            - Hash
        - Arithmetic expressions
        - Built-in library of functions
        - First-class and Higher-order functions
        - Closures
    - Make good, detailed comments so others can read and understand the structure of the code without understanding Go
    - Prepare myself for creating a compiler in Go, maybe C++ or C also
    - Expand my understanding of Go, more specifically its features beyond basic API and testing
    - 10X my developing skills

# Notes

## Comments
Comment structure:
- **File-level comments**
  - These exist at the top of the file to explain what purpose that file serves and the content present in the file.
- **Type comments**
  - Describe a specific type or interface.
  - Also describe the role of the type and sometimes examples of how it's used.
- **Method comments**
  - Describe the purpose of a method, its parameters, and its return values.
  - Might describe additional notes about how it operates, preconditions, postconditions, or examples of output.
- **Marker method comments**
  - Marker methods are used to satisfy a type's interface, like seen here numerous times in the code:

    ```go
    type Boolean struct {
        Token token.Token
        Value bool
    }

    func (b *Boolean) expressionNode()      {}
    func (b *Boolean) TokenLiteral() string { return b.Token.Literal }
    func (b *Boolean) String() string       { return b.Token.Literal }
    ```
