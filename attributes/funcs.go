package attributes

import (
	"fmt"
	"math/rand"
)

func UniqueNumberSequence(numberOfNumbers, numberRange int) []uint16 {
	var sequence []uint16 = make([]uint16, numberOfNumbers)

	for i := range sequence {
		var nextNumber int = rand.Intn(9) + 1
		var looping bool = true
		var allUnique bool = false

		for looping && !allUnique {
			allUnique = true
			for j := range sequence {
				if uint16(nextNumber) == sequence[j] {
					nextNumber = rand.Intn(9) + 1
					allUnique = false
				}
			}
		}

		sequence[i] = uint16(nextNumber)
	}
	return sequence
}

func generateCluster(cluster *[][]uint16, clusterSize int) {
	var values []uint16 = UniqueNumberSequence(clusterSize*clusterSize, clusterSize*clusterSize)
	var indexes []uint16 = UniqueNumberSequence(clusterSize*clusterSize, clusterSize*clusterSize)

	for row := range clusterSize {
		for col := range clusterSize {
			(*cluster)[row][col] = 0
			for i, index := range indexes {
				if uint16(row*clusterSize+col) == index-1 {
					(*cluster)[row][col] = values[i]
				}
			}
		}
	}
}

func clusterValuesClash(
	clusterRow1, clusterCol1 int, cluster1 [][]uint16,
	clusterRow2, clusterCol2 int, cluster2 [][]uint16,
	clusterSize int,
) bool {
	for r1 := 0; r1 < clusterSize; r1++ {
		for c1 := 0; c1 < clusterSize; c1++ {
			val1 := cluster1[r1][c1]
			if val1 == 0 {
				continue
			}

			globalRow1 := clusterRow1*clusterSize + r1
			globalCol1 := clusterCol1*clusterSize + c1

			for r2 := 0; r2 < clusterSize; r2++ {
				for c2 := 0; c2 < clusterSize; c2++ {
					val2 := cluster2[r2][c2]
					if val2 == 0 || val1 != val2 {
						continue
					}

					globalRow2 := clusterRow2*clusterSize + r2
					globalCol2 := clusterCol2*clusterSize + c2

					if globalRow1 == globalRow2 || globalCol1 == globalCol2 {
						return true
					}
				}
			}
		}
	}
	return false
}

func clustersClashInRow(initialValues [][][][]uint16, clusterSize, clusterRow, clusterCol int) bool {
	for row := 0; row < clusterRow; row++ {
		if clusterValuesClash(
			clusterRow, clusterCol, initialValues[clusterRow][clusterCol],
			row, clusterCol, initialValues[row][clusterCol],
			clusterSize,
		) {
			return true
		}
	}
	return false
}

func clustersClashInCol(initialValues [][][][]uint16, clusterSize, clusterRow, clusterCol int) bool {
	for col := 0; col < clusterCol; col++ {
		if clusterValuesClash(
			clusterRow, clusterCol, initialValues[clusterRow][clusterCol],
			clusterRow, col, initialValues[clusterRow][col],
			clusterSize,
		) {
			return true
		}
	}
	return false
}

func regenerateClashingClusters(initialValues *[][][][]uint16, boardSize, clusterSize int) {
	for clusterRow := range boardSize {
		for clusterCol := range boardSize {
			fmt.Println(clusterRow, clusterCol)
			generateCluster(&(*initialValues)[clusterRow][clusterCol], clusterSize)
			for clustersClashInRow(*initialValues, clusterSize, clusterRow, clusterCol) || clustersClashInCol(*initialValues, clusterSize, clusterRow, clusterCol) {
				generateCluster(&(*initialValues)[clusterRow][clusterCol], clusterSize)
			}
		}
	}
}

func GenerateBoard(boardSize, clusterSize int) (initialValues [][][][]uint16) {
	initialValues = make([][][][]uint16, boardSize)

	for clusterRow := range boardSize {
		initialValues[clusterRow] = make([][][]uint16, boardSize)

		for clusterCol := range boardSize {
			initialValues[clusterRow][clusterCol] = make([][]uint16, clusterSize)

			for cellRow := range clusterSize {
				initialValues[clusterRow][clusterCol][cellRow] = make([]uint16, clusterSize)
			}
		}
	}

	regenerateClashingClusters(&initialValues, boardSize, clusterSize)

	return initialValues
}
