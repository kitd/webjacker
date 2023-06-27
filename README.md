# Webjacker

`Webjacker` is a simple package (seriously, it's only one file) to convert your service and/or domain objects into HTTP resources, allowing them to respond to standard HTTP verbs, or your own custom events. This is done by: 

1. injecting a single `*webjacker.HttpResource` field into your struct 
2. specifying what should happen in response to incoming HTTP events
3. registering it with a `ServeMux`.

Originally designed to make backends for [Htmx](https://htmx.org/)-powered web pages, it can easily be repurposed for any REST or AJAX based client.
Example:

```
// An autocompletion service type
type AutoCompleter struct {
    *webjacker.HttpResource
}

// An example autocompleter instance for looking up names
nameSearch := AutoCompleter{
    // Resource will be served on '/name_search' in the ServeMux
    webjacker.NewHttpResource("name_search"), 
}

nameSearch.On(http.MethodGet, 
    func (rw http.ResponseWriter, r *http.Request, params url.Values) {
        prefix := params.Get(nameSearch.Id)
        results := searchCustomers(prefix)
        templates.ExecuteTemplate(rw, "name_results", results)
    })

webjacker.RegisterHttpResource(nameSearch.HttpResource, http.DefaultServeMux)
```

See [examples](./examples/autocompleter) for a fuller example using Htmx. 

The path that a resource appears on is available via the `Path()` method. You can use this to inject links back to this resource in any output rendered to the ResponseWriter. If you can't or don't want to be use the standard HTTP verbs, the path to trigger custom events is available via `EventPath("event_name")`. You then define the handler for it via `.On("event_name", handler)` similar to the example above.

You will no doubt be asking "Why not just implement `http.Handler` and create a `ServeHTTP()` function on your objects?", which is a reasonable question. There are a few reasons:

1. Just having the function is not enough. There's the business of hooking it up under a certain path to a `ServeMux`, extracting parameters, dealing with custom events, responding to different HTTP methods, etc. `Webjacker` makes dealing with all this much simpler.

2. The `ServeHTTP()` function will handle calls for all instances of that struct. This is probably OK for data objects, but not for service objects, which require different handling for different instances. Eg (using the above example), an autocompleter for customer names and one for product names would require different processing, a distinction that would have to be made in the single `ServeHTTP()` function. Again, `Webjacker` makes handling this much simpler by allowing you to hook up processing on a per instance or per type basis.

3. It provides separation of concerns. Let your domain objects and services deal with domain problems. Let `Webjacker` handle the HTTP exchanges for you.