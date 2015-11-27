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

## Writing Guidelines

* Always give categories a title on the first line with `#`.
* Keep the card front as succinct as possible to increase the likelihood that
  it can be matched against on subsequent imports. For example, instead of
  "what is the definition of Mach?" just put "Mach definition?".
* Wrap inline equations in the Anki latex markers `[$]` and `[/$]`.
