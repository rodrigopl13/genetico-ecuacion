package main

import (
	"fmt"
	"fyne.io/fyne"
	"fyne.io/fyne/app"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/widget"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/plotutil"
	"gonum.org/v1/plot/vg"
	"math"
	"solucionador-ecuacion/genetico"
	"time"
)

type application struct {
	solution       *canvas.Image
	window         fyne.Window
	generation     *genetico.Generation
	labelBest      *widget.Label
	bestSolution   float64
	plot           *plot.Plot
	xy             plotter.XYs
	bestChromosome []uint8
}

func main() {
	a := app.New()
	g := application{
		solution: &canvas.Image{
			FillMode: canvas.ImageFillOriginal,
		},
		window:       a.NewWindow("Solución ecuación"),
		generation:   &genetico.Generation{},
		labelBest:    widget.NewLabel(fmt.Sprintf("%.3f", 0.0)),
		bestSolution: math.MaxFloat64,
	}
	g.window.SetContent(
		widget.NewVBox(
			g.solution,
			fyne.NewContainerWithLayout(
				layout.NewGridLayout(3),
				widget.NewButton("Iniciar Applicacion", g.startApp),
				widget.NewLabel("Best Solution"),
				g.labelBest,
			),
		),
	)
	g.window.ShowAndRun()

}

func (a *application) startApp() {
	p, err := plot.New()
	if err != nil {
		panic(err)
	}

	a.plot = p

	p.Title.Text = "Ecuacion"
	p.Y.Label.Text = "Aptitude function"
	p.X.Label.Text = "Generation"
	a.solution.File = "image/graph.png"

	a.xy = plotter.XYs{}
	a.generation = genetico.NewGeneration(100, 3, 255, 0.5, aptFunc)
	fmt.Println(a.generation.Population)
	time.Sleep(1 * time.Second)
	a.createGraph(0)
	fmt.Println(a.generation.Population)
	for i := 0; i < 10; i++ {
		a.generation = genetico.NextGeneration(a.generation)
		fmt.Println(a.generation.Population)
		time.Sleep(1 * time.Second)
		a.createGraph(i + 1)
	}
}

func (a *application) createGraph(i int) {
	a.getBestChromosome()
	points := plotter.XY{}
	points.X = float64(i)
	points.Y = a.bestSolution

	a.xy = append(a.xy, points)
	err := plotutil.AddLinePoints(a.plot, a.xy)
	if err != nil {
		panic(err)
	}

	if err = a.plot.Save(7*vg.Inch, 7*vg.Inch, "image/graph.png"); err != nil {
		panic(err)
	}

	canvas.Refresh(a.solution)
	a.labelBest.SetText(fmt.Sprintf("%.3f", a.bestSolution))
}

func (a *application) getBestChromosome() {
	minSolution := math.MaxFloat64
	index := 0
	for i, v := range a.generation.Aptitud {
		if minSolution > v {
			index = i
			minSolution = v
		}
	}
	if a.bestSolution > minSolution {
		a.bestChromosome = a.generation.Population[index]
		a.bestSolution = minSolution
	}

	fmt.Println("BEST >>>>>>>>>>>>", a.generation.Population[index], "distance: "+fmt.Sprintf("%.3f", minSolution))
}

func aptFunc(c []uint8) float64 {
	//equation
	// A/25.5x^2+B/25.5x+C/25.5 = 13
	x := 2.0
	return math.Abs(((float64(c[0]) / 25.5) * math.Pow(x, 2)) +
		float64((float64(c[1])/25.5)*x) +
		(float64(c[2]) / 25.5) - 13)
}
