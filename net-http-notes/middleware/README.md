## Middleware

**Running Code Before or After Handler Code in a HTTP Request Lifecycle**

Any `type` can be a `Handler` so long as it implements `ServeHTTP`.

so construct our own Handler(middleware handler) that wraps original handler.
