package main

// utils module
//
// Copyright (c) 2019 - Valentin Kuznetsov <vkuznet@gmail.com>
//
import (
	"bufio"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"

	"github.com/gonum/stat"
)

// helper function to get the scores
// it should be extended when new metric is implemented
func getMetric(values []float64) float64 {
	// read _scoreFile
	csvFile, err := os.Open(_scoreFile)
	if err != nil {
		log.Println("error", err, "file", _scoreFile)
		return 0
	}
	defer csvFile.Close()
	reader := csv.NewReader(bufio.NewReader(csvFile))
	var labels []bool
	for {
		line, err := reader.Read()
		if err == io.EOF {
			break
		}
		// true is 0, false is 1 to make gonum be aligned with scikit-earn
		// https://scikit-learn.org/stable/modules/generated/sklearn.metrics.roc_auc_score.html
		// https://godoc.org/github.com/gonum/stat#ROC
		if line[0] == "0" || line[0] == "true" {
			labels = append(labels, true)
		} else {
			labels = append(labels, false)
		}
	}
	// sort our values and labels
	stat.SortWeightedLabeled(values, labels, nil)
	if Config.Metric == "auc" {
		return auc(values, labels)
	} else {
		msg := fmt.Sprintf("Not implemented metric: %s", Config.Metric)
		log.Fatal(errors.New(msg))
	}
	return 0
}

func getScore(file string) float64 {
	csvFile, err := os.Open(file)
	if err != nil {
		log.Println("error", err, "file", file)
		return 0
	}
	defer csvFile.Close()
	reader := csv.NewReader(bufio.NewReader(csvFile))
	var preds []float64
	for {
		line, err := reader.Read()
		if err == io.EOF {
			break
		}
		p, err := strconv.ParseFloat(line[1], 64)
		if err != nil {
			if _verbose > 0 {
				log.Println("findScore: parse error", err, "value", line)
			}
			continue
		}
		preds = append(preds, p)
	}
	return getMetric(preds)
}
