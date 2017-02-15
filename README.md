# Mtr-go
Multi language translator api wrapper in Go, translate or compare strings or arrays of strings with language pairs supported by multiple services.

[![Build Status](https://travis-ci.org/untoreh/mtr-go.svg?branch=master)](https://travis-ci.org/untoreh/mtr-go)
[![codecov](https://codecov.io/gh/untoreh/mtr-go/branch/master/graph/badge.svg)](https://codecov.io/gh/untoreh/mtr-go)
[![CodeFactor](https://www.codefactor.io/repository/github/untoreh/mtr-go/badge)](https://www.codefactor.io/repository/github/untoreh/mtr-go)
[![Code Climate](https://codeclimate.com/github/untoreh/mtr-go/badges/gpa.svg)](https://codeclimate.com/github/untoreh/mtr-go)

## Install 
Download the bin from the release page.

Or get the package
```bash
go get github.com/untoreh/mtr-go
```
 
## Usage
The binary has no flags atm, so just run it
```bash
./mtr
```
## Package
For using directly into code

Pass source/target language and a string or array of strings
```golang
import "github.com/untoreh/mtr-go"
m := mtr_go.New(nil)

m.Tr("en", "fr", "the fox hides quickly", nil)
// returns : Le renard se cache rapidement

```
Map keys are __preserved__.

List of base usable language codes, the priority is to google codes which means if you want 
to translate chinese you should use `zh-TW` or `zh-CN`
```golang
m.SupLangs()
// returns : [ "en", "fr", ... ]
```

Choose which services to use

```golang
m.Tr("en", "fr", "the fox hides quickly", {"google", "bing"});
```

Add a weight to it to specify how many times a service should be chosen over the others
```golang
m.Tr("en", "fr", "the fox hides quickly", {"google" : 50, "bing" : 5})
```

Custom http options 
```golang
m := new(mtr_go.New( {"httpClient" : http.Client{}})
```

Api keys 
```golang
m = new Mtr({"systran_key" : key});
```

## Conventions
- It is _recommended_ to **not** rely on the translation to return _consistent punctuation_, 
therefore input text should be as __atomic__ as possible.
- Some services arbitrarily encode/decode html or even add html tags themselves, such 
aggressive services have active decoding before the output.

## Notes
- Requests are limited to `1000~` __chars__, strings and arrays get _split or merged_ up to this
size to try to make uniform requests. 
- All the parts of a request are run __concurrently__, _pools_ are not used (yet).
- Default services `weight` is `30` for _google, bing, yandex_ and 10 for the rest.
- Cached keys start with `mtr_`
- Because services may be fickle, they will be dropped as they go down or block access.
- Not all services supports all the languages, the group of services used is transparently trimmed to the ones that support the requested language pair.

## TODO/Improvements
In the [issues](https://github.com/untoreh/mtr-go/issues)

## Credits
- [Stichoza/google-translate-php](https://github.com/Stichoza/google-translate-php) - _google token generator_
- [leodido/langcode-conv](https://github.com/leodido/langcode-conv) - _base for the language code converter_

