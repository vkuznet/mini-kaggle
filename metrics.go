package main

// metrics module
//
// Copyright (c) 2019 - Valentin Kuznetsov <vkuznet@gmail.com>
//
import (
	"errors"
	"fmt"
	"log"

	"github.com/gonum/integrate"
	"github.com/gonum/stat"
)

// helper function to calculate metric
// feel free to extend it further to other metric functions
func calcMetric(args ...interface{}) float64 {
	if Config.Metric == "auc" {
		values := args[0].([]float64)
		labels := args[1].([]bool)
		return auc(values, labels)
	} else {
		msg := fmt.Sprintf("Not implemented metric: %s", Config.Metric)
		log.Fatal(errors.New(msg))
	}
	return 0
}

// helper function to calculate AUC
func auc(y []float64, classes []bool) float64 {
	// find tpr, fpr values
	tpr, fpr := stat.ROC(0, y, classes, nil)
	if _verbose > 0 {
		log.Println("preds", y, len(y), "classes", classes, len(classes), "TPR", tpr, "FPR", fpr)
	}
	// compute Area Under Curve
	auc := integrate.Trapezoidal(fpr, tpr)
	return auc
}
