# httpmockserverexample

This is a full working example of implemeting the [git.sr.ht/~ewintr/go-kit/test](git.sr.ht/~ewintr/go-kit/test) package to setup a mock http server. It is based on the blog post: [https://erikwinter.nl/articles/2020/unit-test-outbound-http-requests-in-golang](https://erikwinter.nl/articles/2020/unit-test-outbound-http-requests-in-golang)

The methods in the example program are simple just for the purpose of creating simple requests. They are not examples of best practices in writing http requests.

This package has its own assert and equal but they can be replaced by other testing packages without affecting the mock server.

## Advantages
- This package records requests to the mock server so they can be validated in testing
- Outside services do not need to be running for tests
- No outside mocking services like postman are needed
- Tests run without internet access
- Tests run faster