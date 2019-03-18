# go-tests-kvstore
There are times when you're working with dynamic data that needs to be up to date, but, 
since it's not exactly at hand to retrieve it ever so often, not 100% up to date. Here comes my tool.

I wrote this initially while working on a web crawler's unit tests, which couldn't fetch the urls before each run,
but still needed fresh html to run against. 

My solution was this thing that downloads the html from a set set of urls, encodes the contents to base64 
(so the crawled site owners don't find it on google) and then writes the whole key-value set 
(`[](url => base64(contents))`) to a json file for later usage.

So at first there is this `fetch_urls` command which outputs the whole json to stdout. 
It requires a set of urls to be specified as arguments.

Then you can import the package into your tests file and use the `kvstore.Read` to easily access the decoded data.

```bash
url1='https://google.com'
url2='https://1337x.to/torrent/3651764/Love-Death-And-Robots-S01-COMPLETE-720p-WEB-x264-GalaxyTV/'

go run fetch_urls.go ${url1} ${url2} > data.json

# ofc you can also build a binary / go install it

# ofc you can do a one liner or integrate this into your CI pipeline.
```
this will result in a data.json file containing the following:

```json
[
  {
    "key": "https://google.com",
    "value": "UENGa2IyTjBlWEJsSUdoMGJX..."
  },
  {
    "key": "https://gmail.com",
    "value": "Q2p3aFJFOURWRmxRUlNCb2RH..."
  }
]
```

These are my sane defaults that generate json with base64 encoded contents, but you can plug your own - 
just go ahead and check out the `Write` and `Read` funcs in 
[kv.go](https://github.com/florinutz/go-tests-kvstore/blob/master/kv.go). 
The url fetching thing is also something that suited my workflow, 
but you could Write any kv store that you somehow retrieved previously.
