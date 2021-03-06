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

var findNumbersRegex = regexp.MustCompile(`\d+`)

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
	firstLineNumbers := findNumbersRegex.FindAll(scanner.Bytes(), 3) // 3 instead of 2 to find if there are too many numbers
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
	y, lineNumber := 0, 1 // y is number of saved rows in the board, lineNumber is number of lines read (only for error reporting)
	for scanner.Scan() {
		lineNumber++
		numbers := findNumbersRegex.FindAll(scanner.Bytes(), board.gridSize+1) // +1 to find if there are too many numbers
		if len(numbers) == 0 {
			continue
		}
		if y >= board.gridSize { // we check it as late as here, because there might be lines without numbers at the end
			return nil, fmt.Errorf("too many board lines, expected %d", board.gridSize)
		}
		if len(numbers) != board.gridSize {
			return nil, fmt.Errorf("expected %d numbers, got %d in line %d: %s", board.gridSize, len(numbers), lineNumber, scanner.Text())
		}
		for x, numberBytes := range numbers {
			number, err := strconv.Atoi(string(numberBytes))
			if err != nil {
				return nil, fmt.Errorf("error parsing number %w in line %d: %s", err, lineNumber, scanner.Text())
			}
			if number < 0 || number > board.gridSize {
				return nil, fmt.Errorf("inalid number %d in line %d: %s", number, lineNumber, scanner.Text())
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

func (b *Board) Copy() *Board {
	dataCopy := make([]uint16, len(b.data))
	copy(dataCopy, b.data)
	return &Board{
		data:           dataCopy,
		subgridWidth:   b.subgridWidth,
		subgridHeight:  b.subgridHeight,
		gridSize:       b.gridSize,
		subgridsCountX: b.subgridsCountX,
		subgridsCountY: b.subgridsCountY,
	}
}

// Size returns total width/height of the board.
// For example, for standard 9x9 sudoku with 3x3 subgrids, this will return 9.
func (b *Board) Size() int {
	return b.gridSize
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

func (b *Board) Equal(b2 *Board) bool {
	if b.subgridWidth != b2.subgridWidth ||
		b.subgridHeight != b2.subgridHeight ||
		b.gridSize != b2.gridSize ||
		b.subgridsCountX != b2.subgridsCountX ||
		b.subgridsCountY != b2.subgridsCountY ||
		len(b.data) != len(b2.data) {
		return false
	}
	for i := range b.data {
		if b.data[i] != b2.data[i] {
			return false
		}
	}
	return true
}

func (b *Board) ForEach(operation func(x, y int, n uint16)) {
	for y := 0; y < b.gridSize; y++ {
		for x := 0; x < b.gridSize; x++ {
			operation(x, y, b.Get(x, y))
		}
	}
}

func (b *Board) ForEachNeighbour(x0, y0 int, operation func(x, y int)) {
	gridBeginX := x0 - x0%b.subgridWidth
	gridBeginY := y0 - y0%b.subgridHeight
	gridEndX := gridBeginX + b.subgridWidth
	gridEndY := gridBeginY + b.subgridHeight

	// vertical line above subgrid
	for y := 0; y < gridBeginY; y++ {
		operation(x0, y)
	}

	// subgrid above point
	for y := gridBeginY; y < y0; y++ {
		for x := gridBeginX; x < gridEndX; x++ {
			operation(x, y)
		}
	}

	// point's row, including subgrid, excluding point itself
	for x := 0; x < x0; x++ {
		operation(x, y0)
	}
	for x := x0 + 1; x < b.gridSize; x++ {
		operation(x, y0)
	}

	// subgrid below point
	for y := y0 + 1; y < gridEndY; y++ {
		for x := gridBeginX; x < gridEndX; x++ {
			operation(x, y)
		}
	}

	// vertical line below subgrid
	for y := gridEndY; y < b.gridSize; y++ {
		operation(x0, y)
	}
}

func (b *Board) Validate(validate func(x, y int, n uint16) error, nextFieldGroup func()) error {
	// for each row
	for y := 0; y < b.gridSize; y++ {
		for x := 0; x < b.gridSize; x++ {
			if err := validate(x, y, b.Get(x, y)); err != nil {
				return err
			}
		}
		nextFieldGroup()
	}

	// for each column
	for x := 0; x < b.gridSize; x++ {
		for y := 0; y < b.gridSize; y++ {
			if err := validate(x, y, b.Get(x, y)); err != nil {
				return err
			}
		}
		nextFieldGroup()
	}

	// for each subgrid
	for y0 := 0; y0 < b.gridSize; y0 += b.subgridHeight {
		for x0 := 0; x0 < b.gridSize; x0 += b.subgridWidth {
			for y := y0; y < y0+b.subgridHeight; y++ {
				for x := x0; x < x0+b.subgridWidth; x++ {
					if err := validate(x, y, b.Get(x, y)); err != nil {
						return err
					}
				}
			}
			nextFieldGroup()
		}
	}

	return nil
}

func (b *Board) HaveCommonSubgrid(x1, y1, x2, y2 int) bool {
	gridBeginX1 := x1 - x1%b.subgridWidth
	gridBeginX2 := x2 - x2%b.subgridWidth
	if gridBeginX1 != gridBeginX2 {
		return false
	}
	gridBeginY1 := y1 - y1%b.subgridHeight
	gridBeginY2 := y2 - y2%b.subgridHeight
	return gridBeginY1 == gridBeginY2
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
