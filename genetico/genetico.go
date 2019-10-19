package genetico

import (
	"math"
	"math/rand"
	"sync"
	"time"
)

type Generation struct {
	Population [][]uint8
	Aptitud    []float64
	maxNumber  int
	percentaje float32
}

type aptFunction func(c []uint8) float64

var aptFunc aptFunction

func NewGeneration(population, sizeChromosome, maxNumber int, percentaje float32, af aptFunction) *Generation {
	ng := &Generation{
		Population: make([][]uint8, population),
		Aptitud:    make([]float64, population),
		maxNumber:  maxNumber,
		percentaje: percentaje,
	}
	aptFunc = af
	generateRandom(ng, sizeChromosome)
	calculateAptitud(ng)
	return ng
}

func NextGeneration(currentGeneration *Generation) *Generation {
	population := len(currentGeneration.Population)
	ng := &Generation{
		Population: make([][]uint8, population),
		Aptitud:    make([]float64, population),
		maxNumber:  currentGeneration.maxNumber,
	}

	competeParents(currentGeneration, ng)

	//generateOperations(ng)

	calculateAptitud(ng)
	return ng
}

//
func generateRandom(g *Generation, sizeChromosome int) {
	var wg sync.WaitGroup
	for i := range g.Population {
		wg.Add(1)
		go randomChromosome(g.Population, i, sizeChromosome, g.maxNumber, &wg)
	}
	wg.Wait()
}

func randomChromosome(p [][]uint8, chromosome, sizeChromosome, maxNumber int, wg *sync.WaitGroup) {
	c := make(map[int]bool)
	var random int
	p[chromosome] = make([]uint8, sizeChromosome)
	for i := 0; i < sizeChromosome; i++ {
		rand.Seed(time.Now().UnixNano())
		random = rand.Intn(maxNumber) + 1
		for c[random] {
			rand.Seed(time.Now().UnixNano())
			random = rand.Intn(maxNumber) + 1
		}
		c[random] = true
		p[chromosome][i] = uint8(random)
	}
	wg.Done()
}

func competeSingle(currentGeneration, newGeneration *Generation) {
	population := len(currentGeneration.Population)
	var wg sync.WaitGroup

	for i := 0; i < population; i++ {
		wg.Add(1)
		go reproduceSingleChromosome(currentGeneration, newGeneration, i, currentGeneration.percentaje, &wg)
	}

	wg.Wait()
}

func reproduceSingleChromosome(
	currentGeneration,
	newGeneration *Generation,
	position int,
	percentaje float32,
	wg *sync.WaitGroup,
) {
	population := len(currentGeneration.Population)
	p := float32(population) * percentaje
	bestApt := math.MaxFloat64
	var randomIndex, minIndex int
	for i := 0; i < int(p); i++ {
		rand.Seed(time.Now().UnixNano())
		randomIndex = rand.Intn(population)
		if currentGeneration.Aptitud[randomIndex] < bestApt {
			bestApt = currentGeneration.Aptitud[randomIndex]
			minIndex = randomIndex
		}
	}
	sizeChromosome := len(currentGeneration.Population[minIndex])
	newChromosome := make([]uint8, sizeChromosome)
	copy(newChromosome, currentGeneration.Population[minIndex])
	newGeneration.Population[position] = newChromosome
	wg.Done()
}

func competeParents(currentGeneration, newGeneration *Generation) {
	population := len(currentGeneration.Population)
	var wg sync.WaitGroup

	for i := 0; i < population; i += 2 {
		wg.Add(1)
		go reproduceChildsChromosome(currentGeneration, newGeneration, i, currentGeneration.percentaje, &wg)
	}

	wg.Wait()
}

func reproduceChildsChromosome(
	currentGeneration,
	newGeneration *Generation,
	position int,
	percentaje float32,
	wg *sync.WaitGroup,
) {
	population := len(currentGeneration.Population)
	p := float32(population) * percentaje
	bestApt1 := math.MaxFloat64
	bestApt2 := math.MaxFloat64
	var randomIndex, min1, min2 int
	for i := 0; i < int(p); i++ {
		rand.Seed(time.Now().UnixNano())
		randomIndex = rand.Intn(population)
		if currentGeneration.Aptitud[randomIndex] < bestApt1 {
			bestApt1 = currentGeneration.Aptitud[randomIndex]
			min1 = randomIndex
		}
	}

	for i := 0; i < int(p); i++ {
		rand.Seed(time.Now().UnixNano())
		randomIndex = rand.Intn(population)
		if currentGeneration.Aptitud[randomIndex] < bestApt2 {
			bestApt2 = currentGeneration.Aptitud[randomIndex]
			min2 = randomIndex
		}
	}
	r1, r2 := Cruza(currentGeneration.Population[min1], currentGeneration.Population[min2])
	//fmt.Println("=============", r1, r2)
	newGeneration.Population[position] = r1
	newGeneration.Population[position+1] = r2
	wg.Done()
}

func generateOperations(g *Generation) {
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

func calculateAptitud(g *Generation) {
	for i := range g.Population {
		g.Aptitud[i] = aptFunc(g.Population[i])
	}
}
