package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"

	wm "github.com/ynqa/wego/pkg/model"
	"github.com/ynqa/wego/pkg/model/modelutil/vector"
	w2v "github.com/ynqa/wego/pkg/model/word2vec"

	ann "github.com/rikonor/go-ann"
)

func newWvModel() wm.Model {
	model, err := w2v.New(
		w2v.Window(5),
		w2v.Model(w2v.Cbow),
		w2v.Optimizer(w2v.NegativeSampling),
		w2v.NegativeSampleSize(5),
		w2v.Verbose(),
	)
	if err != nil {
		log.Fatal(err)
	}
	return model
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	model := newWvModel()
	knowledge, err := os.Open("knowledge.txt")
	if err != nil {
		log.Fatal(err)
	}

	defer knowledge.Close()
	content, _ := io.ReadAll(knowledge)
	log.Println(string(content))
	knowledge.Seek(0, io.SeekStart)
	if err := model.Train(knowledge); err != nil {
		log.Fatal(err)
	}

	// model.Save(os.Stdin, vector.Agg)
	v := model.WordVector(vector.Agg)
	// fmt.Println(v)

	vc := [][]float64{}
	for i := 0; i < v.Col(); i++ {
		vv := v.Slice(i)
		fmt.Printf("%v\n", vv)
		vc = append(vc, vv)
	}

	query := "how is dog doing?"
	qm := newWvModel()
	if err := qm.Train(bytes.NewReader([]byte(query))); err != nil {
		log.Fatal(err)
	}

	// qmv := qm.WordVector(vector.Single)
	qmv := qm.WordVector(vector.Agg)
	qv := make([]float64, qmv.Col())
	log.Println("qmv.Col():", qmv.Col())
	log.Println("qmv.Row():", qmv.Row())
	for i := 0; i < qmv.Row(); i++ {
		qvv := qmv.Slice(i)
		qv[i] = average(qvv)
	}
	log.Println("qmv:", qmv)
	log.Println("qv:", qv)

	ex := ann.NewExhaustiveNNer(vc)
	nb := ex.ANN(qv, 4)
	log.Println("neighbors:", nb)
}

type theOrdered interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 |
		~float32 | ~float64
}

func average[T theOrdered](arr []T) float64 {
	var sum T
	length := len(arr)
	for _, a := range arr {
		sum += a
	}
	return float64(sum) / float64(length)
}
