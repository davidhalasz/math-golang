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
	"gonum.org/v1/gonum/stat/distuv"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
	"gonum.org/v1/plot/vg/draw"
	"gonum.org/v1/plot/vg/vgimg"
)

type GnumPlotResponse struct {
	MeanPNG        []byte  `json:"mean_png"`
	Mean           float64 `json:"mean"`
	MedianPNG      []byte  `json:"median_png"`
	Median         float64 `json:"median"`
	StdDev         float64 `json:"std_dev"`
	Variance       float64 `json:"variance"`
	StdDevVarPNG   []byte  `json:"std_dev_var_png"`
	PDFPNG         []byte  `json:"pdf_png"`
	PMFPNG         []byte  `json:"pmf_png"`
	PoissonPNG     []byte  `json:"poisson_png"`
	Covariance1PNG []byte  `json:"covariance1_png"`
	Covariance1    float64 `json:"covariance1"`
	Covariance2PNG []byte  `json:"covariance2_png"`
	Covariance2    float64 `json:"covariance2"`
	Correlation    float64 `json:"correlation"`
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
	histogram, _ := plotter.NewHist(values, 50)

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
	histogram, _ := plotter.NewHist(values, 50)

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

func (app *application) StdVar(w http.ResponseWriter, r *http.Request) {
	seed := time.Now().UnixNano()
	localRand := rand.New(rand.NewSource(seed))

	n := 10000
	mean := 100.0
	stdDev := 100.0

	incomes := make([]float64, n)
	for i := range incomes {
		incomes[i] = localRand.NormFloat64()*stdDev + mean
	}

	stdDeviation := stat.StdDev(incomes, nil)

	variance := stat.Variance(incomes, nil)

	p := plot.New()

	values := make(plotter.Values, len(incomes))
	copy(values, incomes)
	histogram, _ := plotter.NewHist(values, 50)

	red := uint8(71)
	green := uint8(85)
	blue := uint8(105)

	histogram.FillColor = color.NRGBA{red, green, blue, 255}

	p.Add(histogram)

	imgCanvas := vgimg.New(vg.Points(800), vg.Points(400))
	dc := draw.New(imgCanvas)

	p.Draw(dc)

	var pngBuffer bytes.Buffer

	err := png.Encode(&pngBuffer, imgCanvas.Image())
	if err != nil {
		log.Printf("error encoding PNG: %v\n", err)
		return
	}

	svgResponse := GnumPlotResponse{StdDevVarPNG: pngBuffer.Bytes(), StdDev: stdDeviation, Variance: variance}
	jsonResponse, err := json.Marshal(svgResponse)
	if err != nil {
		fmt.Fprintln(w, "Error sending json:", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResponse)
}

func (app *application) PDF(w http.ResponseWriter, r *http.Request) {
	// Create a normal distribution with mean 0 and standard deviation 1
	dist := distuv.Normal{
		Mu:    0,
		Sigma: 1,
	}

	// Create a range of x values
	x := make([]float64, 6001)
	for i := range x {
		x[i] = -3 + float64(i)*0.001
	}

	// Create a plotter.XYs to hold the x, y values
	pts := make(plotter.XYs, len(x))
	for i, val := range x {
		pts[i].X = val
		pts[i].Y = dist.Prob(val)
	}

	red := uint8(71)
	green := uint8(85)
	blue := uint8(105)

	// Create a plot
	p := plot.New()
	p.Title.Text = "Normal Distribution"

	// Create a line plot
	line, err := plotter.NewLine(pts)
	if err != nil {
		panic(err)
	}
	line.Color = color.NRGBA{red, green, blue, 255}

	p.Add(line)

	imgCanvas := vgimg.New(vg.Points(800), vg.Points(400))
	dc := draw.New(imgCanvas)

	p.Draw(dc)

	var pngBuffer bytes.Buffer

	err = png.Encode(&pngBuffer, imgCanvas.Image())
	if err != nil {
		log.Printf("error encoding PNG: %v\n", err)
		return
	}

	svgResponse := GnumPlotResponse{PDFPNG: pngBuffer.Bytes()}
	jsonResponse, err := json.Marshal(svgResponse)
	if err != nil {
		fmt.Fprintln(w, "Error sending json:", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResponse)
}

func (app *application) Binomial(w http.ResponseWriter, r *http.Request) {
	n := 10.0
	p := 0.5

	// Define the binomial distribution
	dist := distuv.Binomial{
		N: n,
		P: p,
	}

	// Create a range of x values
	x := make([]float64, 10000)
	for i := range x {
		x[i] = float64(i) * 0.001
	}

	// Create a plotter.XYs to hold the x, y values
	pts := make(plotter.XYs, len(x))
	for i, val := range x {
		pts[i].X = val
		pts[i].Y = dist.Prob(val)
	}

	// Plot the PMF
	pmf := plot.New()

	pmf.Title.Text = "Binomial Distribution PMF"
	pmf.X.Label.Text = "X"
	pmf.Y.Label.Text = "Probability"

	// Set the X-axis range from 0 to 10
	pmf.X.Min = 0
	pmf.X.Max = 10

	// Create a line plot
	line, _, err := plotter.NewLinePoints(pts)
	if err != nil {
		panic(err)
	}

	red := uint8(71)
	green := uint8(85)
	blue := uint8(105)
	line.Color = color.NRGBA{red, green, blue, 255}
	pmf.Add(line)

	imgCanvas := vgimg.New(vg.Points(800), vg.Points(400))
	dc := draw.New(imgCanvas)

	pmf.Draw(dc)

	var pngBuffer bytes.Buffer

	err = png.Encode(&pngBuffer, imgCanvas.Image())
	if err != nil {
		log.Printf("error encoding PNG: %v\n", err)
		return
	}

	svgResponse := GnumPlotResponse{PMFPNG: pngBuffer.Bytes()}
	jsonResponse, err := json.Marshal(svgResponse)
	if err != nil {
		fmt.Fprintln(w, "Error sending json:", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResponse)
}

func (app *application) Poisson(w http.ResponseWriter, r *http.Request) {
	mu := 500.0
	x := make([]float64, 0)
	for i := 400.0; i < 600.0; i += 0.5 {
		x = append(x, i)
	}

	// Define the binomial distribution
	dist := distuv.Poisson{
		Lambda: mu,
	}

	// Create a plotter.XYs to hold the x, y values
	pts := make(plotter.XYs, len(x))
	for i, val := range x {
		pts[i].X = val
		pts[i].Y = dist.Prob(val)
	}

	// Plot the PMF
	poisson := plot.New()

	poisson.Title.Text = "Poisson Probability Mass Function"

	// Create a line plot
	line, err := plotter.NewLine(pts)
	if err != nil {
		panic(err)
	}

	red := uint8(71)
	green := uint8(85)
	blue := uint8(105)
	line.Color = color.NRGBA{red, green, blue, 255}
	poisson.Add(line)

	imgCanvas := vgimg.New(vg.Points(800), vg.Points(400))
	dc := draw.New(imgCanvas)

	poisson.Draw(dc)

	var pngBuffer bytes.Buffer

	err = png.Encode(&pngBuffer, imgCanvas.Image())
	if err != nil {
		log.Printf("error encoding PNG: %v\n", err)
		return
	}

	svgResponse := GnumPlotResponse{PoissonPNG: pngBuffer.Bytes()}
	jsonResponse, err := json.Marshal(svgResponse)
	if err != nil {
		fmt.Fprintln(w, "Error sending json:", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResponse)
}

func deMean(x []float64) []float64 {
	xmean := stat.Mean(x, nil)
	result := make([]float64, len(x))

	for i, xi := range x {
		result[i] = xi - xmean
	}
	return result
}

func covariance(x, y []float64) float64 {
	n := len(x)
	demeanX := deMean(x)
	demeanY := deMean(y)
	dotProduct := 0.0

	for i := 0; i < n; i++ {
		dotProduct += demeanX[i] * demeanY[i]
	}
	return dotProduct / float64(n-1)
}

func correlation(x, y []float64) float64 {
	stddevX := stat.StdDev(x, nil)
	stddevY := stat.StdDev(y, nil)
	return covariance(x, y) / (stddevX * stddevY)
}

func (app *application) CovCor(w http.ResponseWriter, r *http.Request) {
	var pageSpeeds, purchaseAmount1, purchaseAmount2 []float64
	for i := 0; i < 1000; i++ {
		pageSpeed := rand.NormFloat64()*1.0 + 3.0
		pageSpeeds = append(pageSpeeds, pageSpeed)
		purchase := rand.NormFloat64()*10.0 + 50.0

		purchaseAmount1 = append(purchaseAmount1, rand.NormFloat64()*10.0+50.0)
		purchaseAmount2 = append(purchaseAmount2, purchase/pageSpeed)
	}

	// Scatter plot 1
	red := uint8(71)
	green := uint8(85)
	blue := uint8(105)

	pts1 := make(plotter.XYs, len(pageSpeeds))
	for i := range pts1 {
		pts1[i].X = pageSpeeds[i]
		pts1[i].Y = purchaseAmount1[i]
	}

	p1 := plot.New()

	s1, err := plotter.NewScatter(pts1)
	if err != nil {
		fmt.Println("Error creating scatter plot:", err)
		return
	}

	p1.Add(s1)
	p1.Title.Text = "Scatter Plot"
	p1.X.Label.Text = "Page Speeds"
	p1.Y.Label.Text = "Purchase Amounts"

	s1.Color = color.NRGBA{red, green, blue, 255}

	imgCanvas1 := vgimg.New(vg.Points(800), vg.Points(400))
	dc1 := draw.New(imgCanvas1)

	p1.Draw(dc1)

	var pngBuffer1 bytes.Buffer

	err = png.Encode(&pngBuffer1, imgCanvas1.Image())
	if err != nil {
		log.Printf("error encoding PNG: %v\n", err)
		return
	}

	// Scatter plot 2
	pts2 := make(plotter.XYs, len(pageSpeeds))
	for i := range pts2 {
		pts2[i].X = pageSpeeds[i]
		pts2[i].Y = purchaseAmount2[i]
	}

	p2 := plot.New()

	s2, err := plotter.NewScatter(pts2)
	if err != nil {
		fmt.Println("Error creating scatter plot:", err)
		return
	}

	p2.Add(s2)
	p2.Title.Text = "Scatter Plot"
	p2.X.Label.Text = "Page Speeds"
	p2.Y.Label.Text = "Purchase Amounts"

	s2.Color = color.NRGBA{red, green, blue, 255}

	imgCanvas2 := vgimg.New(vg.Points(800), vg.Points(400))
	dc2 := draw.New(imgCanvas2)

	p2.Draw(dc2)

	var pngBuffer2 bytes.Buffer

	err = png.Encode(&pngBuffer2, imgCanvas2.Image())
	if err != nil {
		log.Printf("error encoding PNG: %v\n", err)
		return
	}

	covResult1 := covariance(pageSpeeds, purchaseAmount1)
	covResult2 := covariance(pageSpeeds, purchaseAmount2)
	correlation := correlation(pageSpeeds, purchaseAmount2)

	svgResponse := GnumPlotResponse{Covariance1PNG: pngBuffer1.Bytes(), Covariance1: covResult1, Covariance2PNG: pngBuffer2.Bytes(), Covariance2: covResult2, Correlation: correlation}
	jsonResponse, err := json.Marshal(svgResponse)
	if err != nil {
		fmt.Fprintln(w, "Error sending json:", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResponse)
}
