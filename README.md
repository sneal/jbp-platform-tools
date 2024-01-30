# CF Platform Tools

Tools to help rationalize about your CF system.

Requires go 1.20+

## Usage
All the commands assume you're already signed into the CF CLI as an admin user or a user which has wide read only 
visibility into the CF v3 API.

### List all apps using a buildpack
```bash
$ go run main.go buildpack --name java_buildpack_offline
```

### List all apps using port 18000
```bash
$ go run main.go routing --port 18000
```

### List all apps using network policy port 61443
```bash
$ go run main.go network-policy --port 61443
```