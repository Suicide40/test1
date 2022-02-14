# MD5 web hasher

The app makes http requests and prints the address of the request along with the
MD5 hash of the response.

The tool ables to limit the number of parallel requests, to prevent exhausting local resources. To do that just use a flag `-parallel` to indicate this limit. By default, it is 10 if the flag is not provided.

##Examples:

```
./test1 -parallel 3 www.ya.ru www.google.com adjust.com google.com facebook.com yahoo.com yandex.com twitter.com reddit.com/r/funny reddit.com/r/notfunny baroquemusiclibrary.com
www.google.com 5e5f73d62a1a2efc45724844062ab079
www.ya.ru 72ec8fab7120c824a9abf6708feff309
adjust.com 339bacbfff8e2f28efdd80c819bc833c
google.com d12bdbf3147f087a87db186b72ffdeba
facebook.com 7f90991bf36b8fbff29d43a43737ef3f
yandex.com 3f7e2243237d2e805cb9e43c5fb8fe3e
yahoo.com ff354183f281a93519b9d066491c5c52
reddit.com/r/notfunny 76eec8d707252a00f76dd9a3f198ece4
twitter.com d4c0449cb99120286fefa49d93bd5cb0
reddit.com/r/funny 729e6f3cc3c0cb8031fb53eac0e47774
baroquemusiclibrary.com 990a7d6bf2963d76e7e6b15635e21b11
```

```
./test1 pornhub.com
pornhub.com 335665c5da0716ad9336c75c309a3494
```