package genetico

import (
	"fmt"
	"math"
	"math/rand"
	"sync"
	"time"
)

type Generation struct {
	Population [][]uint8
	Aptitud    []float64
}

type aptFunction func(c []uint8) float64

var (
	aptFunc           aptFunction
	maxNumberAlelo    int
	competePercentaje float32
	sizeChromosome    int
)

func NewGenetic(
	population,
	sizeChromosome,
	maxNumberAlelo int,
	competePercentaje float32,
	af aptFunction,
) *Generation {
	ng := &Generation{
		Population: make([][]uint8, population),
		Aptitud:    make([]float64, population),
	}
	aptFunc = af
	maxNumberAlelo = maxNumberAlelo
	competePercentaje = competePercentaje
	sizeChromosome = sizeChromosome

	ng.generateRandom()
	ng.calculateAptitud()

	return ng
}

func (g *Generation) NextGeneration() *Generation {
	population := len(g.Population)
	ng := &Generation{
		Population: make([][]uint8, population),
		Aptitud:    make([]float64, population),
	}

	g.competeParents(ng)

	//g.generateOperations(ng)

	ng.calculateAptitud()
	return ng
}

//
func (g *Generation) generateRandom() {
	var wg sync.WaitGroup
	for i := range g.Population {
		wg.Add(1)
		go g.randomChromosome(i, &wg)
	}
	wg.Wait()
}

func (g *Generation) randomChromosome(chromosome int, wg *sync.WaitGroup) {
	c := make(map[int]bool)
	var random int
	g.Population[chromosome] = make([]uint8, sizeChromosome)
	for i := 0; i < sizeChromosome; i++ {
		rand.Seed(time.Now().UnixNano())
		random = rand.Intn(maxNumberAlelo) + 1
		for c[random] {
			rand.Seed(time.Now().UnixNano())
			random = rand.Intn(maxNumberAlelo) + 1
		}
		c[random] = true
		g.Population[chromosome][i] = uint8(random)
	}
	wg.Done()
}

func (g *Generation) competeSingle(currentGeneration, newGeneration *Generation) {
	population := len(currentGeneration.Population)
	var wg sync.WaitGroup

	for i := 0; i < population; i++ {
		wg.Add(1)
		go g.reproduceSingleChromosome(newGeneration, i, &wg)
	}

	wg.Wait()
}

func (g *Generation) reproduceSingleChromosome(
	newGeneration *Generation,
	position int,
	wg *sync.WaitGroup,
) {
	population := len(g.Population)
	p := float32(population) * competePercentaje
	bestApt := math.MaxFloat64
	var randomIndex, minIndex int
	for i := 0; i < int(p); i++ {
		rand.Seed(time.Now().UnixNano())
		randomIndex = rand.Intn(population)
		if g.Aptitud[randomIndex] < bestApt {
			bestApt = g.Aptitud[randomIndex]
			minIndex = randomIndex
		}
	}
	sizeChromosome := len(g.Population[minIndex])
	newChromosome := make([]uint8, sizeChromosome)
	copy(newChromosome, g.Population[minIndex])
	newGeneration.Population[position] = newChromosome
	wg.Done()
}

func (g *Generation) competeParents(newGeneration *Generation) {
	population := len(g.Population)
	var wg sync.WaitGroup

	for i := 0; i < population; i += 2 {
		wg.Add(1)
		go g.reproduceChildsChromosome(newGeneration, i, &wg)
	}

	wg.Wait()
}

func (g *Generation) reproduceChildsChromosome(
	newGeneration *Generation,
	position int,
	wg *sync.WaitGroup,
) {
	population := len(g.Population)
	p := float32(population) * competePercentaje
	bestApt1 := math.MaxFloat64
	bestApt2 := math.MaxFloat64
	var randomIndex, min1, min2 int
	for i := 0; i < int(p); i++ {
		rand.Seed(time.Now().UnixNano())
		randomIndex = rand.Intn(population)
		if g.Aptitud[randomIndex] < bestApt1 {
			bestApt1 = g.Aptitud[randomIndex]
			min1 = randomIndex
		}
	}

	for i := 0; i < int(p); i++ {
		rand.Seed(time.Now().UnixNano())
		randomIndex = rand.Intn(population)
		if g.Aptitud[randomIndex] < bestApt2 {
			bestApt2 = g.Aptitud[randomIndex]
			min2 = randomIndex
		}
	}
	fmt.Println("============", min1 == min2)
	r1, r2 := Cruza(g.Population[min1], g.Population[min2])
	fmt.Println("=============", r1, r2)
	newGeneration.Population[position] = r1
	newGeneration.Population[position+1] = r2
	wg.Done()
}

func (g *Generation) generateOperations() {
	var wg sync.WaitGroup
	countInversion := 0
	var r int
	for i := range g.Population {
		wg.Add(1)
		rand.Seed(time.Now().UnixNano())
		r = rand.Intn(2)
		if r == 0 && countInversion <= 50 {
			go Inversion(g.Population[i], &wg)
			countInversion++
		} else {
			go Intercambio(g.Population[i], &wg)
		}
	}
	wg.Wait()
}

func (g *Generation) calculateAptitud() {
	for i := range g.Population {
		g.Aptitud[i] = aptFunc(g.Population[i])
	}
}
