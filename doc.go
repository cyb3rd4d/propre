/*
Package propre helps you to implement a clean architecture for HTTP applications in Go.

The Clean Architecture has been theorized by Robert C. Martin (aka Uncle Bob) in 2012 in his [blog].
There is no exact or perfect way to apply it, but its layered nature makes it quite easy to
design some kind of "framework" at least to avoid writing some boilerplate code again and again.

Currently Propre (which means "clean" in French) exposes a [HTTPHandler] type which is generic over
2 types: Input and Output:
  - Input is the type used by a use case as an entrypoint. It contains the data needed by the business
    rules of your use case, or an error raised by the request decoder which is responsible of reading
    the incoming HTTP request.
  - Output is the type produced by the use case as a result. It also holds either the data of a
    successful scenario or an error.

This is not idiomatic for functions or methods in Go to not return an error type in case of failure.
Propre enforces the use of "monads" to make the mechanics easier between the application layers. If you're
not familiar with monads, you can check out the fantastic [samber/mo] project.

Please see [this repository] containing examples about how to implement Propre in a project.

[blog]: https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html
[samber/mo]: https://github.com/samber/mo
[this repository]: https://github.com/cyb3rd4d/propre-examples
*/
package propre
