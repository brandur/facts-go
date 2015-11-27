# facts-go

[![build status](https://secure.travis-ci.org/brandur/facts-go.png)](https://travis-ci.org/brandur/facts-go)

A program that exports the contents of `categories/` (Markdown files containing
lists of facts) to a tab-separated file (TSV) suitable for ingestion into Anki.

## Run

```
go get -u github.com/brandur/facts-go
facts-go
```

Then look for the output file `facts.tsv`.

## Test

```
go test
```
