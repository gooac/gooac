
# ![Gooa](assets/gooacompiler.png)

Gooa is a Lua Preprocessor written in Go *(very clearly)*   
This is just a compiler library. Feel free to use this in its current state but I am currently working on tools to wrap this.

# Features
- ### Middleware
    Middleware allows you to control the processing at any time, whether it be allowing you to omit AST Nodes directly at parse time or to run a regex scan post compilation, Middleware allows you to control the entire process easily.
- ### Builtin Syntax Features
    - **Named Function Arguments**: Allows defining a default value for a functions argument at definition, no more `a=a or b()`.
    - **C-Style Comments**: Allows you to use C-Style comments the same way you would normal comments.
    - **Continue Statement**: Adds the `continue` keyword to control loops in a much simpler way.
    - **Function Attributes**: Allows you to easily wrap functions through a series of calls, progressively allowing you to create functions that are more powerful than they need to be.
