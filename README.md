# Webjacker

`Webjacker` is a simple package (seriously, it's only one file) to convert your service and/or domain objects into HTTP resources, allowing them to respond to standard HTTP verbs, or your own custom events. This is done by: 

1. injecting a single `*webjacker.HttpResource` field into your struct 
2. specifying what should happen in response to incoming HTTP events
3. registering it with a `ServeMux`.

Originally designed to make backends for [Htmx](https://htmx.org/)-powered web pages, it can easily be repurposed for any REST or AJAX based client.
Example:

```
// An example autocompletion service
type AutoCompleter struct {
	*webjacker.HttpResource
}

func (ac *AutoCompleter) Handle(rw http.ResponseWriter, r *http.Request, params url.Values) {
	prefix := params.Get("name")
	results := searchCustomers(prefix)
	templates.ExecuteTemplate(rw, "name_results", results)
}

// An example autocompleter instance for looking up names
nameSearch := AutoCompleter{
    // Resource served on '/name_search' in the ServeMux
    webjacker.NewHttpResource("name_search"), 
}

nameSearch.On(http.MethodGet, nameSearch.Handle)

webjacker.RegisterHttpResource(nameSearch.HttpResource, http.DefaultServeMux)
```

See [examples](./example/) for a fuller example. 

The path that a resource appears on is available via the `Path()` method. You can use this to inject calls back to the resource in any output rendedred to the ResponseWriter. If you can't or don't want to be constrained to the standard HTTP verbs, the path to trigger custom events is available via `PathForEvent(my_event)`. You then define the handler for it via `.On("my_event" handler)` similar to the example above.
