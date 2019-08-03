package main

// utils module
//
// Copyright (c) 2019 - Valentin Kuznetsov <vkuznet@gmail.com>
//
import (
	"bufio"
	"encoding/csv"
	"io"
	"log"
	"os"
	"strconv"

	"github.com/gonum/stat"
)

// helper function to get scores from provided file
func getScores(file string) (float64, float64) {
	csvFile, err := os.Open(file)
	if err != nil {
		log.Println("error", err, "file", file)
		return 0, 0
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
	// split predictions according to configuration score split
	// one set represents public scores and another private one
	split := int(float64(len(preds)) * Config.ScoreSplit)
	publicPreds := preds[0:split]
	privatePreds := preds[split:len(preds)]

	// read scores from _scoreFile
	csvScoreFile, err := os.Open(_scoreFile)
	if err != nil {
		log.Println("error", err, "file", _scoreFile)
		return 0, 0
	}
	defer csvScoreFile.Close()
	scoreReader := csv.NewReader(bufio.NewReader(csvScoreFile))
	var labels []bool
	for {
		line, err := scoreReader.Read()
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
	publicLabels := labels[0:split]
	privateLabels := labels[split:len(labels)]

	// sort our values and labels
	stat.SortWeightedLabeled(publicPreds, publicLabels, nil)
	stat.SortWeightedLabeled(privatePreds, privateLabels, nil)

	// return our public and private scores
	return calcMetric(publicPreds, publicLabels), calcMetric(privatePreds, privateLabels)
}
