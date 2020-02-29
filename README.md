# CGMiner API for Go #

## Installation ##

    # install the library:
    go get github.com/gOsToFf/go-cgminer-api

    // Use in your .go code:
    import (
        "github.com/gOsToFf/go-cgminer-api"
    )

## API Documentation ##

Full godoc output from the latest code in master is available here:

https://godoc.org/github.com/gOsToFf/go-cgminer-api

## Quickstart ##

```go
package main

import (
    "github.com/gOsToFf/go-cgminer-api"
    "log"
)

func main() {
    miner := cgminer.New("localhost", 4028)
    summary, err := miner.Summary()
    if err != nil {
        log.Fatalln("Unable to connect to CGMiner: ", err)
    }

    // Convert hashrate into KH/s
    hashrate := summary.MHSav * 1000.0

    log.Printf("Average Hashrate: %f KH/s\n", hashrate)
}
```
