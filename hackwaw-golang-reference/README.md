## Golang example

This is simple app that implements the HTTP server `/latest` endpoint
and tries to reach `https://www.google.com` through HTTP proxy.

The HTTP proxy is configured with `HTTP_PROXY` and `HTTPS_PROXY` environment variable.

To workaround the trusted certificate issue,
it currently it uses `InsecureSkipVerify` of `tls.Config`.
