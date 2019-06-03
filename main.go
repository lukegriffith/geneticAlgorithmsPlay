package main

import (
	"fmt"
	"math/rand"
	"sort"
	"strconv"
	"time"
)


type chromosome struct {
	Data string
}

func newChromosome(c string) *chromosome {
	return &chromosome{c}
}


func (c *chromosome) fitness() int  {

	i, err := strconv.ParseInt(c.Data, 2, 64)

	if err != nil {
		fmt.Println(err)
	}

	return int(i)

}

func (c *chromosome) mutate() *chromosome {

	rand.Seed(time.Now().UnixNano())
	mbit := rand.Intn(8)

	if c.Data[mbit] == '1' {
		c.Data= c.Data[:mbit] + "0" + c.Data[mbit+1:]
	} else {
		c.Data= c.Data[:mbit] + "1" + c.Data[mbit+1:]
	}

	return c
}

func (c *chromosome) crossover(x *chromosome) *chromosome {

	d := c.Data[:4] + x.Data[4:]

	return newChromosome(d)
}

func genPopulation(size int) []*chromosome {

	var min = 10
	var max = 256

	population := make([]*chromosome, size)

	for x := 0; x < size; x++ {
		rand.Seed(time.Now().UnixNano())
		y := rand.Intn((max - min) + min)
		population[x] = newChromosome(fmt.Sprintf("%08b", y))
	}

	return population
}

// eliteList contains ordered map index by fitness.
func breedPopulation(popSize int, eliteList []int, pop map[int]*chromosome) []*chromosome{

	var population []*chromosome
	population = append(population, pop[eliteList[0]].crossover(pop[eliteList[1]]))

	rand.Seed(time.Now().UnixNano())

	// Randomness, potentially 1 new chromosome will mutate in a generation.
	// Stops from getting trapped on a local optimum.
	shouldMutate := rand.Intn(popSize + 10)

	for x := 0; x < popSize -1; x++ {
		rand.Seed(time.Now().UnixNano())
		parnter := rand.Intn(popSize -1)

		if shouldMutate == x {
			population = append(population, pop[eliteList[x]].crossover(pop[eliteList[parnter]]).mutate())
		} else {
			population = append(population, pop[eliteList[x]].crossover(pop[eliteList[parnter]]))
		}
	}

	return population
}

func main() {

	var popSize int = 6
	var maxGenerations int = 100
	var generation int = 0
	var terminationCriterion bool = false
	var popHealth int

	genX := genPopulation(popSize)

	for (! terminationCriterion) && generation < maxGenerations {

		popHealth = 0

		orderedMap := make(map[int]*chromosome)
		var intList []int

		for _, x := range genX {

			fit := x.fitness()

			popHealth += fit
			orderedMap[fit] = x
			intList = append(intList, fit)


			if fit == 0 {
				terminationCriterion = true
			}
		}

		fmt.Println("Generation:", generation, "Average:", popHealth / popSize)

		sort.Ints(intList)

		genX = breedPopulation(popSize, intList, orderedMap)

		generation++
	}
}


