package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"image/color"
	"image/png"
	"log"
	"math/rand"
	"net/http"
	"sort"
	"time"

	"gonum.org/v1/gonum/stat"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
	"gonum.org/v1/plot/vg/draw"
	"gonum.org/v1/plot/vg/vgimg"
)

type GnumPlotResponse struct {
	MeanPNG   []byte  `json:"mean_png"`
	Mean      float64 `json:"mean"`
	MedianPNG []byte  `json:"median_png"`
	Median    float64 `json:"median"`
}

func (app *application) HomePage(w http.ResponseWriter, r *http.Request) {
	if err := app.renderTemplate(w, r, "home", nil); err != nil {
		app.errorLog.Println(err)
	}
}

func (app *application) StatisticsPage(w http.ResponseWriter, r *http.Request) {
	if err := app.renderTemplate(w, r, "statistics", nil); err != nil {
		app.errorLog.Println(err)
	}
}

func (app *application) Mean(w http.ResponseWriter, r *http.Request) {
	// Use current time as a seed for random number generation
	seed := time.Now().UnixNano()
	localRand := rand.New(rand.NewSource(seed))

	// sample number
	n := 10000

	// mean and standard deviation
	mean := 27000.0
	stdDev := 15000.0

	// create normalized sample
	incomes := make([]float64, n)
	for i := range incomes {
		incomes[i] = localRand.NormFloat64()*stdDev + mean
	}

	// get average
	meanValue := stat.Mean(incomes, nil)

	// New Histogram
	p := plot.New()

	// add incomes to histogram
	values := make(plotter.Values, len(incomes))
	copy(values, incomes)
	histogram, _ := plotter.NewHist(values, 20)

	red := uint8(71)
	green := uint8(85)
	blue := uint8(105)

	histogram.FillColor = color.NRGBA{red, green, blue, 255}

	p.Add(histogram)

	// if zou would like to create png file
	//
	// wt, err := p.WriterTo(512, 512, "png")
	// if err != nil {
	// 	log.Fatalf("could not create writer: %v", err)
	// }

	// f, err := os.Create("out.png")
	// if err != nil {
	// 	log.Fatalf("could not create out.png: %v", err)
	// }
	// defer f.Close()

	// _, err = wt.WriteTo(f)
	// if err != nil {
	// 	log.Fatalf("could not write to out.png: %v", err)
	// }

	// if err := f.Close(); err != nil {
	// 	log.Fatalf("could not close out.png: %v", err)
	// }

	// Create an Image Canvas to draw the plot
	imgCanvas := vgimg.New(vg.Points(800), vg.Points(400))
	dc := draw.New(imgCanvas)

	// Set up the plot and draw it onto the canvas
	p.Draw(dc)

	// Create a buffer to store PNG data
	var pngBuffer bytes.Buffer

	// Encode the image to PNG and write to the buffer
	err := png.Encode(&pngBuffer, imgCanvas.Image())
	if err != nil {
		log.Printf("error encoding PNG: %v\n", err)
		return
	}

	// send JSON response
	svgResponse := GnumPlotResponse{MeanPNG: pngBuffer.Bytes(), Mean: meanValue}
	jsonResponse, err := json.Marshal(svgResponse)
	if err != nil {
		fmt.Fprintln(w, "Error sending json:", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResponse)
}

func (app *application) Median(w http.ResponseWriter, r *http.Request) {
	// Use current time as a seed for random number generation
	seed := time.Now().UnixNano()
	localRand := rand.New(rand.NewSource(seed))

	// sample number
	n := 10000

	// mean and standard deviation
	mean := 27000.0
	stdDev := 15000.0

	// create normalized sample
	incomes := make([]float64, n)
	for i := range incomes {
		incomes[i] = localRand.NormFloat64()*stdDev + mean
	}

	// Sort the incomes slice
	sort.Float64s(incomes)

	// Calculate the median
	medianValue := stat.Quantile(0.5, stat.Empirical, incomes, nil)

	// New Histogram
	p := plot.New()

	// add incomes to histogram
	values := make(plotter.Values, len(incomes))
	copy(values, incomes)
	histogram, _ := plotter.NewHist(values, 20)

	red := uint8(71)
	green := uint8(85)
	blue := uint8(105)

	histogram.FillColor = color.NRGBA{red, green, blue, 255}

	p.Add(histogram)

	// Create an Image Canvas to draw the plot
	imgCanvas := vgimg.New(vg.Points(800), vg.Points(400))
	dc := draw.New(imgCanvas)

	// Set up the plot and draw it onto the canvas
	p.Draw(dc)

	// Create a buffer to store PNG data
	var pngBuffer bytes.Buffer

	// Encode the image to PNG and write to the buffer
	err := png.Encode(&pngBuffer, imgCanvas.Image())
	if err != nil {
		log.Printf("error encoding PNG: %v\n", err)
		return
	}

	// send JSON response
	svgResponse := GnumPlotResponse{MedianPNG: pngBuffer.Bytes(), Median: medianValue}
	jsonResponse, err := json.Marshal(svgResponse)
	if err != nil {
		fmt.Fprintln(w, "Error sending json:", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResponse)
}
