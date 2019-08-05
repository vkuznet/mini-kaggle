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
		// we assume here that our line in csv file contains at last
		// position our prediction value. We parse it as float64
		// as it represents generic case for predictions (ints or floats)
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
	var scores []interface{}
	for {
		line, err := scoreReader.Read()
		if err == io.EOF {
			break
		}
		// we assume that score file contains at first position
		// given score
		scores = append(scores, line[0])
	}

	// split scores into public and private sets
	publicScores := scores[0:split]
	privateScores := scores[split:len(scores)]

	// return our public and private metrics
	return calcMetric(publicPreds, publicScores), calcMetric(privatePreds, privateScores)
}
