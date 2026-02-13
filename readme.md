# rag-go

Exploration and readily discarded repo to explore things relative to RAG
and its adjacent.

<!--
This is using package github.com/kirill-scherba/word2vec
with branch of feature/native_libw2

so first we need to clone it:
git clone https://github.com/kirill-scherba/word2vec --branch feature/native_libw2 --depth 1

also this requires Go 1.25.7 hence if our version is lower:
go version # go1.25.2

we need to install it first:
go install golang.org/dl/go1.25.7@latest
go1.25.7 download

set our alias:
alias go=go1.25.7

or powershell:
set-alias go=go1.25.7

next build and install the cmd:
cd word2vec/cmd
go install

back to up and modify our go.mod:
cd ../../rag-go

edit it to have:
replace "github.com/kirill-scherba/word2vec => ../word2vec
-->

## Instructions

The repo we're using is requiring `go1.25.7` hence we need to ensure
our `go version` is higher than that. In case it's lower, we can do:

```bash
go install golang.org/dl/go1.25.7@latest
go1.25.7 download
# on unix
alias go=go1.25.7
# on powershell
set-alias go=go1.25.7
go version # confirming it's targeted version
```

We're also using [`mage`][mgweb] so either installing it by running

```bash
go install github.com/magefile/mage@latest
```

or with various other ways we can check on its [website][mgweb].
Assuming now we have `mage`. We can do:

```bash
git clone https://github.com/mashingan/rag-go
cd rag-go
mage init
mage example_w2v # to see our setup is complete
mage downloadModel
```

[mgweb]: https://magefile.org
