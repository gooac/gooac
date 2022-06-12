
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


# Quickstart
Add to your module with `go get github.com/gooac/gooac`

```go
package main
import "github.com/gooac/gooac"

func main() {
    g := gooa.NewGooa()
    g.Use(gooa.AttributeMiddleware())

    code, err := g.Compile([]byte(`
		print(123)

		function a(b=c)
			print(b)
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

# Syntax Examples
### Named Function Arguments
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

### C-Style Comments
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

### Continue Statement
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

### Function Attributes
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
