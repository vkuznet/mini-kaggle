package main

// metrics module
//
// Copyright (c) 2019 - Valentin Kuznetsov <vkuznet@gmail.com>
//
import (
	"errors"
	"fmt"
	"log"

	"github.com/bsm/mlmetrics"
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
	} else if Config.Metric == "accuracy" {
		yTrue := args[0].([]int)
		yPred := args[1].([]int)
		return accuracy(yTrue, yPred)
	} else if Config.Metric == "logloss" {
		predictions := args[0].([]map[string]float64)
		categories := args[1].([]string)
		return logloss(predictions, categories)
	} else if Config.Metric == "mae" {
		yTrue := args[0].([]float64)
		yPred := args[1].([]float64)
		return mae(yTrue, yPred)
	} else if Config.Metric == "mse" {
		yTrue := args[0].([]float64)
		yPred := args[1].([]float64)
		return mse(yTrue, yPred)
	} else if Config.Metric == "rmse" {
		yTrue := args[0].([]float64)
		yPred := args[1].([]float64)
		return rmse(yTrue, yPred)
	} else if Config.Metric == "msle" {
		yTrue := args[0].([]float64)
		yPred := args[1].([]float64)
		return msle(yTrue, yPred)
	} else if Config.Metric == "rmsle" {
		yTrue := args[0].([]float64)
		yPred := args[1].([]float64)
		return rmsle(yTrue, yPred)
	} else if Config.Metric == "r2" {
		yTrue := args[0].([]float64)
		yPred := args[1].([]float64)
		return r2(yTrue, yPred)
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

func confusionMatrix(yTrue []int, yPred []int) *mlmetrics.ConfusionMatrix {
	mat := mlmetrics.NewConfusionMatrix()
	for i := range yTrue {
		mat.Observe(yTrue[i], yPred[i])
	}
	return mat
}
func accuracy(yTrue []int, yPred []int) float64 {
	mat := confusionMatrix(yTrue, yPred)
	return mat.Accuracy()
}
func logloss(preds []map[string]float64, categories []string) float64 {
	metric := mlmetrics.NewLogLoss()
	for i, actual := range categories {
		probability := preds[i][actual]
		metric.Observe(probability)
	}
	return metric.Score()
}
func regression(yTrue, yPred []float64) *mlmetrics.Regression {
	metric := mlmetrics.NewRegression()
	for i := range yTrue {
		metric.Observe(yTrue[i], yPred[i])
	}
	return metric
}
func mae(yTrue, yPred []float64) float64 {
	metric := regression(yTrue, yPred)
	return metric.MAE()
}
func mse(yTrue, yPred []float64) float64 {
	metric := regression(yTrue, yPred)
	return metric.MSE()
}
func rmse(yTrue, yPred []float64) float64 {
	metric := regression(yTrue, yPred)
	return metric.RMSE()
}
func msle(yTrue, yPred []float64) float64 {
	metric := regression(yTrue, yPred)
	return metric.MSLE()
}
func rmsle(yTrue, yPred []float64) float64 {
	metric := regression(yTrue, yPred)
	return metric.RMSLE()
}
func r2(yTrue, yPred []float64) float64 {
	metric := regression(yTrue, yPred)
	return metric.R2()
}
