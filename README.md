# Chain
Execute a series of functions, passing the request to each defined callback.

Typically used as middleware for modifying request headers, or the request responses. In addition, you can short-ciruit requests before they're sent out, and use this to define easily adaptable transports for logging / instrumentation.

[Modeled off of okhttp](https://github.com/square/okhttp/blob/1d8233ddb7a0dfa490a340a06433909148f21610/okhttp/src/main/java/okhttp3/Interceptor.java)


TODO: Add examples and documentation :-D, please use tests as a reference for now.
