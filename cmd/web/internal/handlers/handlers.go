package handlers

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

	"github.com/davidhalasz/gomath/cmd/web/internal/config"
	"github.com/davidhalasz/gomath/cmd/web/internal/render"
	"gonum.org/v1/gonum/stat"
	"gonum.org/v1/gonum/stat/distuv"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/plotutil"
	"gonum.org/v1/plot/vg"
	"gonum.org/v1/plot/vg/draw"
	"gonum.org/v1/plot/vg/vgimg"
)

var app *config.AppConfig

func NewHandlers(a *config.AppConfig) {
	app = a
}

type GnumPlotResponse struct {
	MeanPNG             []byte  `json:"mean_png"`
	Mean                float64 `json:"mean"`
	MedianPNG           []byte  `json:"median_png"`
	Median              float64 `json:"median"`
	StdDev              float64 `json:"std_dev"`
	Variance            float64 `json:"variance"`
	StdDevVarPNG        []byte  `json:"std_dev_var_png"`
	PDFPNG              []byte  `json:"pdf_png"`
	PMFPNG              []byte  `json:"pmf_png"`
	PoissonPNG          []byte  `json:"poisson_png"`
	Covariance1PNG      []byte  `json:"covariance1_png"`
	Covariance1         float64 `json:"covariance1"`
	Covariance2PNG      []byte  `json:"covariance2_png"`
	Covariance2         float64 `json:"covariance2"`
	Correlation         float64 `json:"correlation"`
	LinearRegressionR   float64 `json:"linearRegressionR"`
	LinearRegression    []byte  `json:"linearRegression"`
	PolynomalRegression []byte  `json:"polynomalRegression"`
}

var plotColor = map[string]uint8{
	"r": uint8(71),
	"g": uint8(85),
	"b": uint8(105),
}

func HomePage(w http.ResponseWriter, r *http.Request) {
	if err := render.Template(w, r, "home.page.gohtml", nil); err != nil {
		app.ErrorLog.Println(err)
	}
}

func StatisticsPage(w http.ResponseWriter, r *http.Request) {
	if err := render.Template(w, r, "statistics.page.gohtml", nil); err != nil {
		app.ErrorLog.Println(err)
	}
}

func Mean(w http.ResponseWriter, r *http.Request) {
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

	histogram.FillColor = color.NRGBA{plotColor["r"], plotColor["g"], plotColor["b"], 255}

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

func Median(w http.ResponseWriter, r *http.Request) {
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

	histogram.FillColor = color.NRGBA{plotColor["r"], plotColor["g"], plotColor["b"], 255}

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

func StdVar(w http.ResponseWriter, r *http.Request) {
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

	histogram.FillColor = color.NRGBA{plotColor["r"], plotColor["g"], plotColor["b"], 255}

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

func PDF(w http.ResponseWriter, r *http.Request) {
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

	// Create a plot
	p := plot.New()
	p.Title.Text = "Normal Distribution"

	// Create a line plot
	line, err := plotter.NewLine(pts)
	if err != nil {
		panic(err)
	}
	line.Color = color.NRGBA{plotColor["r"], plotColor["g"], plotColor["b"], 255}

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

func Binomial(w http.ResponseWriter, r *http.Request) {
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

	line.Color = color.NRGBA{plotColor["r"], plotColor["g"], plotColor["b"], 255}
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

func Poisson(w http.ResponseWriter, r *http.Request) {
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

	line.Color = color.NRGBA{plotColor["r"], plotColor["g"], plotColor["b"], 255}
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

func CovCor(w http.ResponseWriter, r *http.Request) {
	var pageSpeeds, purchaseAmount1, purchaseAmount2 []float64
	for i := 0; i < 1000; i++ {
		pageSpeed := rand.NormFloat64()*1.0 + 3.0
		pageSpeeds = append(pageSpeeds, pageSpeed)
		purchase := rand.NormFloat64()*10.0 + 50.0

		purchaseAmount1 = append(purchaseAmount1, rand.NormFloat64()*10.0+50.0)
		purchaseAmount2 = append(purchaseAmount2, purchase/pageSpeed)
	}

	// Scatter plot 1

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

	s1.Color = color.NRGBA{plotColor["r"], plotColor["g"], plotColor["b"], 255}

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

	s2.Color = color.NRGBA{plotColor["r"], plotColor["g"], plotColor["b"], 255}

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

func LinearRegression(w http.ResponseWriter, r *http.Request) {
	// Generate 1000 random points
	localRand := rand.New(rand.NewSource(0))
	n := 1000
	pageSpeeds := make([]float64, n)
	purchaseAmount := make([]float64, n)

	for i := 0; i < n; i++ {
		pageSpeeds[i] = localRand.NormFloat64()*1.0 + 3.0
		purchaseAmount[i] = 100.0 - (pageSpeeds[i]+localRand.NormFloat64()*0.1)*3.0
	}

	// Perform linear regression
	slope, intercept := stat.LinearRegression(pageSpeeds, purchaseAmount, nil, false)

	// Calculate R-squared value
	rSquared := stat.RSquared(pageSpeeds, purchaseAmount, nil, slope, intercept)
	fmt.Printf("R-squared: %f\n", rSquared)

	// Create a plot
	p := plot.New()

	// Create points for the scatter plot
	points := make(plotter.XYs, n)
	for i := range pageSpeeds {
		points[i].X = pageSpeeds[i]
		points[i].Y = purchaseAmount[i]
	}

	// Create a scatter plot
	scatter, err := plotter.NewScatter(points)
	if err != nil {
		log.Fatal(err)
	}
	scatter.Color = color.NRGBA{plotColor["r"], plotColor["g"], plotColor["b"], 255}

	// Add scatter plot to the plot
	p.Add(scatter)

	// Create points for the linear regression line
	line := plotter.NewFunction(func(x float64) float64 { return intercept*x + slope })
	line.LineStyle.Width = vg.Points(3)
	line.Color = plotutil.Color(0)

	// Add linear regression line to the plot
	p.Add(line)

	// Set plot title and labels
	p.Title.Text = "Linear Regression on Random Points"
	p.X.Label.Text = "X"
	p.Y.Label.Text = "Y"

	// // Save the plot to a PNG file
	// if err := p.Save(8*vg.Inch, 4*vg.Inch, "random_linear_regression.png"); err != nil {
	// 	log.Fatal(err)
	// }

	imgCanvas := vgimg.New(vg.Points(800), vg.Points(400))
	dc := draw.New(imgCanvas)

	p.Draw(dc)

	var pngBuffer bytes.Buffer

	err = png.Encode(&pngBuffer, imgCanvas.Image())
	if err != nil {
		log.Printf("error encoding PNG: %v\n", err)
		return
	}

	svgResponse := GnumPlotResponse{LinearRegression: pngBuffer.Bytes(), LinearRegressionR: rSquared}
	jsonResponse, err := json.Marshal(svgResponse)
	if err != nil {
		fmt.Fprintln(w, "Error sending json:", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResponse)
}
