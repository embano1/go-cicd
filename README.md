## About

Little tool to trigger command, executable, shell script, etc. ("-e") when a file or directory ("-f") changes. See `examples` folder for inspiration.

Of course, make sure to install go-cicd first and puth it in your $PATH. You can pick a release or build your own (see below). 

```
   _____  ____     _____ _____       _____ _____
  / ____|/ __ \   / ____|_   _|     / ____|  __ \
 | |  __| |  | | | |      | |______| |    | |  | |
 | | |_ | |  | | | |      | |______| |    | |  | |
 | |__| | |__| | | |____ _| |_     | |____| |__| |
  \_____|\____/   \_____|_____|     \_____|_____/

        The World's most basic CI/CD Tool
```

## Build

### Binary (needs Go tools installed)

```
go build -o go-cicd cmd/main.go
# copy into your $PATH
```

### Docker Image

```
docker build -t <USER/IMAGE:TAG> .
```