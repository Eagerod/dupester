# Dupester

A little project I whipped together because I wanted to find similar documents among an unstructured collection of files I had.

Extracts text from provided files using an Apache Tika server, and saves them into Elasticsearch for comparrison purposes.
Presently, all it can really do is run the Elasticsearch `more_like_this` query, but it's pretty primitive as-is.

It seems to _technically_ function, but it has yet to actually be useful.

## Usage

```
docker-compose up -d
dupester run /a/file/path
```

Requires the Tika Server and ES Server to both be up.
Doesn't contain a configuration file yet to configure hostnames for the two servers either.

Should probably have a configuration file to set the hostnames of those two.
