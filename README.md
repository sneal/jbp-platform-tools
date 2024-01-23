# CF Java Buildpack Platform Tools

Tools to help rationalize Java buildpack upgrades.

Requires go 1.20+

## Usage
Assuming you're already signed into the CF CLI, list all apps using java_buildpack_offline:
```bash
$ go run main.go
```

List all apps using a specifically named buildpack
```bash
$ go run main.go --name binary_buildpack
```
