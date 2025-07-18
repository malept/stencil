---
order: 1013
---

<!-- Generated by tools/docgen. DO NOT EDIT. -->

# module.Call

Call executes a template function by name with the provided arguments.
The function must have been exported by the module that provides it
through the [Export] function.

Attempting to call a function that does not exist will return an error
outside of the first pass where it will return `nil` instead.

The template that is called must return a value using the `return`
template function, which is only available in this context.

In addition, all of the file, stencil and other functions are in the
context of the owning template, not the template calling the function.

`.` in a template function acts the same way as it does for [TplStencil.Include](#TplStencil.Include) (`stencil.Include`). Meaning, it points to [Values](#Values). The caller passed data is accessible on `.Data`.

Example:

```go
// module-a
{{- define "HelloWorld" }}
{{- return (printf "Hello, %s!" .Data) }}
{{- end }}
{{ module.Export "HelloWorld" }}

// module-b
{{ module.Call "github.com/rgst-io/module-a.HelloWorld" "Jared" }}
// Output: Hello, Jared
```
