## Testing

A sample test is included to help you verify that your provider correctly implements the [libdns contract](https://github.com/libdns/libdns/blob/master/libdns.go).

To use it, copy the file into your provider directory, rename it to `provider_test.go`, and update the `newProvider` method to return an instance of your provider.

After that, you should be able to run your tests like this:

```bash

API_KEY=.... go test -v 
```