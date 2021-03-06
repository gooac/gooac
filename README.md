
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
    - **Shorthand Syntax**: Adds a few different shorthands for certain things, also just a few fun ones for looks.


# Quickstart
Add to your module with `go get github.com/gooac/gooac`

```go
package main
import "github.com/gooac/gooac"

func main() {
    g := gooa.NewGooa() // Create a new instance of the tokenizer, parser and compiler
    g.Use(gooa.AttributeMiddleware()) // Add the middleware for function attributes

    // Compile the code
    code, err := g.Compile([]byte(`
		function a(fn)
			return (function(...)
                print("attribute called")
                fn(...)
            end )
		end

        $[a()]
        function some()
            print(123)
        end
	`))

    if err {
        print("Errored! ", err)
        return 
    }
    
    print("Success! \n")
    print(code)
}
```

# Examples
## Named Function Arguments
```lua
function test(a = 1, b = 2, c = somecall(1, 2, 3))
    print(a, b, c)
end
```
Compiles into
```lua
function test(a, b, c)
    a = a or 1
    b = b or 2
    c = c or somecall(1, 2, 3)
end
```

## C-Style Comments
```lua
// Valid Comment
-- Also Valid Comment

/*
    multiline comment
*/

--[==[
    Multiline Comment
]==]
```

## Continue Statement
> Warning: Continues in the base compilre use goto, dont be surprised if you break something with it, especially relating to whether the goto is actually visible to the scope continue is being used in.

```lua
for i=0, 10 do
    if i == 1 then
        continue
    end

    print(i)
end
```
Compiles Into
```lua
for i=0, 10 do
    if i == 1 then
        goto cont_0x0F00003010
    end

    print(i)
    ::cont_0x0F00003010::
end
```

## Function Attributes
```lua
$[route("/")]
local function index()
    print(123)
end
```
Compiles Into
```lua
local function index()
    print(123)
end
index = route(index, "/")
```

## Shorthand Syntax
- Call Arrow
    ```lua
    a()->b()
    ```
    Compiles into

    ```lua
    a():b()
    ```

- Short Function Declarations
    ```lua
    fn a
    end
    ```
    Compiles into
    ```
    function a()
    end
    ```

    This is based on 2 things, first, `fn` is a shorthand for `function`.  
    Second, function arguments can be omitted and are identical to `()`

- Elif!
    ```lua
    if a then

    elseif b then

    elif c then

    end
    ```
    Compiles into
    ```lua
    if a then
    elseif b then
    elseif c then
    end
    ```
