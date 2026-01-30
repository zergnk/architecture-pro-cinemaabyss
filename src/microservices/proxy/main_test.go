package main

import (
	"testing"
)

const firstMax = 20
const secondMax = 30

func TestCreateBinaryCounter(t *testing.T) {

	var result = createBinaryCounter(firstMax, secondMax)
	if result.firstMax != firstMax {
		t.Errorf("firstMax = %d, ожидали %d", firstMax, result.firstMax)
	}

	if result.secondMax != secondMax {
		t.Errorf("secondMax = %d, ожидали %d", secondMax, result.secondMax)
	}

	if result.firstCounter != 0 {
		t.Errorf("firstCounter = %d, ожидали %d", 0, result.firstCounter)
	}

	if result.secondCounter != 0 {
		t.Errorf("secondCounter = %d, ожидали %d", 0, result.secondCounter)
	}
}

func TestNext(t *testing.T) {
	var result = createBinaryCounter(firstMax, secondMax)

	result.next()

	if result.firstCounter+result.secondCounter != firstMax+secondMax-1 {
		t.Errorf("%d + %d != %d + %d - 1", result.firstCounter, result.secondCounter, firstMax, secondMax)
	}

	var firstCounter = result.firstCounter
	var secondCounter = result.secondCounter

	for i := 0; i < firstMax+secondMax-1; i++ {
		next := result.next()

		if next == 1 {
			firstCounter--
		} else if next == 2 {
			secondCounter--
		} else {
			t.Errorf("Неожиданное значение: %d", next)
		}
	}

	if result.firstCounter != 0 {
		t.Errorf("result.firstCounter = %d, ожидали %d", 0, result.firstCounter)
	}

	if result.secondCounter != 0 {
		t.Errorf("result.secondCounter = %d, ожидали %d", 0, result.secondCounter)
	}

	if firstCounter != 0 {
		t.Errorf("firstCounter = %d, ожидали %d", 0, firstCounter)
	}

	if secondCounter != 0 {
		t.Errorf("secondCounter = %d, ожидали %d", 0, secondCounter)
	}
}
