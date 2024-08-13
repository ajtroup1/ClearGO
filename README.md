# ClearGO

# Notes

## Comments
Comment structure:
    - File-level comments
        - These exists at the top of the file to explain what purpose that file serves and the content present in the file
    - Type comments
        - Describe a specific type / interface
        - Also describes the role of the type and sometimes examples of how it's used
    - Method comments
        - Describe the purpose of a method, its parameters, and its return values
        - Might describe additional notes about how it operates, preconditions, postconditions, or examples of output
    - Marker method comments
        - Marker methods are used to satisfy a type's inteface, like seen here numerous times in the code:
            - ```go
            type Boolean struct {
                Token token.Token
                Value bool
            }

            func (b *Boolean) expressionNode()      {}
            func (b *Boolean) TokenLiteral() string { return b.Token.Literal }
            func (b *Boolean) String() string       { return b.Token.Literal }
            ```
