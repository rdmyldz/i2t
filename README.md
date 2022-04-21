- building for windows

```
GOOS=windows CGO_ENABLED=1 CC=x86_64-w64-mingw32-gcc go build
```

- run this in order to copy all necessary dlls

```
mingw-copy-deps.sh /usr/x86_64-w64-mingw32/sys-root cli.exe
```
