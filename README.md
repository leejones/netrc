# netrc

Package netrc gets credentials for a given host from a netrc file.

## Example Usage

This example shows how to make an HTTP GET request using the credentials from
the user's `.netrc` file in their home directory:

```go
// note: some error handling is removed for brevity

url, _ := url.Parse("https://example.com/")
client := http.Client{Timeout: 5 * time.Second}
request, _ := http.NewRequest("GET", url.String(), nil)

credentials, _ := netrc.Get(url.Host)
if err != nil {
  fmt.Fprintf(os.Stderr, "Failed to load credentials from netrc file: %v\n", err)
  os.Exit(1)
}

request.SetBasicAuth(credentials.Login, credentials.Password)

response, _ := client.Do(request)
// handle err and response etc
```
