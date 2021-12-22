package board

import (
	"bufio"
	"fmt"
	"io"
	"math"
	"regexp"
	"strconv"
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

// New creates board with given SUBGRID width and height.
// The board is always square and number of rows / columns
// must be the same as number of fields in a subgrid, so
// total board size is not provided as an argument -
// it must be always subgridWidth*subgridHeight.
// For example, with subgrid 3x2 the board will be 6x6:
// +-------+-------+
// | 0 0 0 | 0 0 0 |
// | 0 0 0 | 0 0 0 |
// +-------+-------+
// | 0 0 0 | 0 0 0 |
// | 0 0 0 | 0 0 0 |
// +-------+-------+
// | 0 0 0 | 0 0 0 |
// | 0 0 0 | 0 0 0 |
// +-------+-------+
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

var filterLineRegex = regexp.MustCompile(`\d+`)

// NewFromSerializedFormat creates board from serialized format. It accepts
// data format produced by Serialize, i.e. the first line is
// subgrid width and height, following lines is board data. Example:
// 2 2
// +-----+-----+
// | 0 0 | 0 3 |
// | 0 1 | 0 4 |
// +-----+-----+
// | 4 2 | 3 1 |
// | 1 3 | 4 2 |
// +-----+-----+
// In fact, it is more "tolerant", i.e. numbers in each line can be
// separated by anything, for exammple this would also work:
// 2x2 :)
// 0 0 0 3
// 0 1 0 4
// 4 2 3 1
// 1 3 4 2
// some random comment not containing digits
func NewFromSerializedFormat(reader io.Reader) (*Board, error) {
	scanner := bufio.NewScanner(reader)
	if !scanner.Scan() {
		return nil, fmt.Errorf("error - no data")
	}
	firstLineNumbers := filterLineRegex.FindAll(scanner.Bytes(), 2)
	if len(firstLineNumbers) != 2 {
		return nil, fmt.Errorf("error parsing - expected 2 numbers, line: %s", scanner.Text())
	}
	subgridWidth, err := strconv.Atoi(string(firstLineNumbers[0]))
	if err != nil {
		return nil, fmt.Errorf("error parsing number %w in line: %s", err, scanner.Text())
	}
	subgridHeight, err := strconv.Atoi(string(firstLineNumbers[1]))
	if err != nil {
		return nil, fmt.Errorf("error parsing number %w in line: %s", err, scanner.Text())
	}
	board, err := New(subgridWidth, subgridHeight)
	if err != nil {
		return nil, fmt.Errorf("error creating board: %w", err)
	}
	y := 0
	for scanner.Scan() {
		numbers := filterLineRegex.FindAll(scanner.Bytes(), board.gridSize)
		if len(numbers) == 0 {
			continue
		}
		if y >= board.gridSize { // we check it as late as here, because there might be lines without numbers at the end
			return nil, fmt.Errorf("too many board lines, expected %d", board.gridSize)
		}
		if len(numbers) != board.gridSize {
			return nil, fmt.Errorf("expected %d numbers, got %d in line %s", board.gridSize, len(numbers), scanner.Text())
		}
		for x, numberBytes := range numbers {
			number, err := strconv.Atoi(string(numberBytes))
			if err != nil {
				return nil, fmt.Errorf("error parsing number %w in line: %s", err, scanner.Text())
			}
			if number < 0 || number > board.gridSize {
				return nil, fmt.Errorf("inalid number %d in line: %s", number, scanner.Text())
			}
			board.Set(x, y, uint16(number))
		}
		y++
	}
	if y != board.gridSize {
		return nil, fmt.Errorf("invalid number of board lines, expected %d, got %d", board.gridSize, y)
	}
	return board, nil
}

func (b *Board) Get(x, y int) uint16 {
	offset := y*b.gridSize + x
	return b.data[offset]
}

func (b *Board) Set(x, y int, value uint16) {
	if value > uint16(b.gridSize) {
		panic(fmt.Sprintf("cannot set value %d for grid with size %d", value, b.gridSize))
	}
	offset := y*b.gridSize + x
	b.data[offset] = value
}

func (b *Board) ForEachInRow(y int, operation func(x, y int)) {
	for x := 0; x < b.gridSize; x++ {
		operation(x, y)
	}
}

func (b *Board) ForEachInColumn(x int, operation func(x, y int)) {
	for y := 0; y < b.gridSize; y++ {
		operation(x, y)
	}
}

func (b *Board) ForEach(operation func(x, y int)) {
	for x := 0; x < b.gridSize; x++ {
		for y := 0; y < b.gridSize; y++ {
			operation(x, y)
		}
	}
}

func (b *Board) ForEachInSubgrid(x, y int, operation func(x, y int)) {
	gridBeginX := x - x%b.subgridWidth
	gridBeginY := y - y%b.subgridHeight
	for dy := 0; dy < b.subgridHeight; dy++ {
		for dx := 0; dx < b.subgridWidth; dx++ {
			operation(gridBeginX+dx, gridBeginY+dy)
		}
	}
}

func (b *Board) Serialize(writer io.Writer) error {
	if _, err := io.WriteString(writer, fmt.Sprintf("%d %d\n", b.subgridWidth, b.subgridHeight)); err != nil {
		return err
	}
	if _, err := io.WriteString(writer, b.String()); err != nil {
		return err
	}
	return nil
}

// String writes board in "ASCII art", for example:
// +-------+-------+
// | 1 2 3 | 4 5 6 |
// | 0 0 0 | 0 0 0 |
// +-------+-------+
// | 0 0 0 | 0 0 0 |
// | 0 0 0 | 0 0 0 |
// +-------+-------+
// | 0 0 0 | 0 0 0 |
// | 6 5 4 | 3 2 1 |
// +-------+-------+
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
