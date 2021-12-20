package board

import (
	"fmt"
	"math"
	"strings"
)

const MaxSize = math.MaxUint16

type Board struct {
	data           []uint16
	subgridWidth   int
	subgridHeight  int
	gridSize       int
	subgridsCountX int
	subgridsCountY int
}

func New(subgridWidth, subgridHeight int) (*Board, error) {
	if subgridWidth < 1 || subgridHeight < 1 {
		return nil, fmt.Errorf("invalid grid size, subgrid sizes must be at least 1, got %d, %d", subgridWidth, subgridHeight)
	}
	if subgridWidth > MaxSize || subgridHeight > MaxSize {
		return nil, fmt.Errorf("invalid grid size, subgrid sizes can be max %d, got %d, %d", MaxSize, subgridWidth, subgridHeight)
	}
	gridSize := subgridHeight * subgridWidth
	if gridSize > MaxSize {
		return nil, fmt.Errorf("grid size (%d) > max available grid size (%d)", gridSize, MaxSize)
	}
	subgridsCountX := subgridHeight // = gridSize / subgridWidth
	subgridsCountY := subgridWidth  // = gridSize / subgridHeight
	return &Board{
		data:           make([]uint16, gridSize*gridSize),
		subgridWidth:   subgridWidth,
		subgridHeight:  subgridHeight,
		gridSize:       gridSize,
		subgridsCountX: subgridsCountX,
		subgridsCountY: subgridsCountY,
	}, nil
}

func (b *Board) FillExampleData() {
	lastIndex := b.gridSize*b.gridSize - 1
	for i := 0; i < b.gridSize; i++ {
		number := uint16(i%b.gridSize + 1)
		b.data[i] = number
		b.data[lastIndex-i] = number
	}
}

func (b *Board) String() string {
	var s strings.Builder
	digitLen := len(fmt.Sprint(b.gridSize))
	charsPerSubgridX := b.subgridWidth + 1 + b.subgridWidth*digitLen
	s.Grow((b.gridSize + b.subgridsCountY + 1) * (charsPerSubgridX*b.subgridsCountX + b.subgridsCountX + 2))
	dataIndex := 0
	for y := 0; y < b.gridSize; y++ {
		if y%b.subgridHeight == 0 {
			b.writeBoardHeaderLine(charsPerSubgridX, &s)
		}
		for x := 0; x < b.gridSize; x++ {
			if x%b.subgridWidth == 0 {
				s.WriteString("| ")
			}
			s.WriteString(fmt.Sprintf("%*d ", digitLen, b.data[dataIndex]))
			dataIndex++
		}
		s.WriteString("|\n")
	}
	b.writeBoardHeaderLine(charsPerSubgridX, &s)
	return s.String()
}

func (b *Board) writeBoardHeaderLine(charsPerSubgridX int, s *strings.Builder) {
	for subgrid := 0; subgrid < b.subgridsCountX; subgrid++ {
		s.WriteString("+")
		for i := 0; i < charsPerSubgridX; i++ {
			s.WriteString("-")
		}
	}
	s.WriteString("+\n")
}
