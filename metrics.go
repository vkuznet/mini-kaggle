package main

// metrics module
//
// Copyright (c) 2019 - Valentin Kuznetsov <vkuznet@gmail.com>
//
import (
	"log"

	"github.com/gonum/integrate"
	"github.com/gonum/stat"
)

// Metrics structure to hold information about availabel metrics
type Metrics struct{}

// MetricsMap contains a map of metrics and their associative functions
func MetricsMap() map[string]string {
	metricsMap := make(map[string]string)
	metricsMap["auc"] = "auc"
	return metricsMap
}

// helper function to calculate AUC
func (Metrics) auc(y []float64, classes []bool) float64 {
	// find tpr, fpr values
	tpr, fpr := stat.ROC(0, y, classes, nil)
	if _verbose > 0 {
		log.Println("preds", y, len(y), "classes", classes, len(classes), "TPR", tpr, "FPR", fpr)
	}
	// compute Area Under Curve
	auc := integrate.Trapezoidal(fpr, tpr)
	return auc
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
