module rag-go

go 1.25.5

require (
	github.com/kirill-scherba/word2vec v1.0.7
	github.com/magefile/mage v1.15.0
	github.com/rikonor/go-ann v0.0.0-20180822120657-22df548f1e42
	github.com/ynqa/wego v0.0.0-20230402162916-bce06112d2fe
)

replace github.com/kirill-scherba/word2vec => ./deps/word2vec

require (
	github.com/gonum/blas v0.0.0-20181208220705-f22b278b28ac // indirect
	github.com/gonum/floats v0.0.0-20181209220543-c233463c7e82 // indirect
	github.com/gonum/internal v0.0.0-20181124074243-f884aa714029 // indirect
	github.com/gonum/lapack v0.0.0-20181123203213-e4cdc5a0bff9 // indirect
	github.com/gonum/matrix v0.0.0-20181209220409-c518dec07be9 // indirect
	github.com/inconshreveable/mousetrap v1.0.0 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/spf13/cobra v1.1.1 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	golang.org/x/sync v0.0.0-20201020160332-67f06af15bc9 // indirect
)
