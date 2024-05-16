// Package graph ...
package graph

import (
	"fmt"
	"image/color"
	"io/fs"
	"math"
	"math/rand"
	"os"
	"path/filepath"
	"regexp"
	"strconv"

	"gonum.org/v1/gonum/stat/distuv"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
)

const (
	PNG  = ".png"
	JPG  = ".jpg"
	JPEG = ".jpeg"
	PDF  = ".pdf"
	svg  = ".svg"
)

func CreateRandomGraphicDefault() Settings {
	return CreateRandomGraphic(5, 300)
}

// Settings is config struct to get some point on graphic
// Range
type Settings struct {
	Range       [2]int
	DotsCount   int
	graphicName string
	graphicDots *plotter.XYs
	path        string
}

// GetGraphicName Возвращает имя графика с которым происходит работа
func (s Settings) GetGraphicName() string {
	return s.graphicName
}

// GetGraphicDots возвращает ссылку на срез точек
// на основании, которых строился график
func (s Settings) GetGraphicDots() *plotter.XYs {
	return s.graphicDots
}

// SetSettings Устанавливает параметры
// Значения Range должны быть в диапазоне 0-100
func (s *Settings) SetSettings(newRange [2]int, newDotsCount int) bool {
	for _, numb := range newRange {
		if numb > 100 || numb < 0 {
			return false
		}
	}
	if newRange[0] > newRange[1] {
		newRange[1], newRange[0] = newRange[0], newRange[1]
	}
	s.Range = newRange
	s.DotsCount = newDotsCount
	return true
}

// CreateDotsGraphics создаёт график с указанием точек
func (s Settings) CreateDotsGraphics() []plotter.XY {
	//указание границ с использованием s.Range
	min, max := len(*s.graphicDots)/100*s.Range[0], len(*s.graphicDots)/100*s.Range[1]
	//формула 95%?
	//mean + stdDev*3
	//mean - stdDev*3
	mean := max - min
	stdDev := mean / 6
	normalDist := distuv.Normal{
		Mu:    float64(mean),
		Sigma: float64(stdDev),
	}
	fmt.Println(len(*s.graphicDots))
	fmt.Println(s.Range[0])
	fmt.Println(mean)
	fmt.Println(stdDev)
	fmt.Println(max)
	fmt.Println(min)
	scatterData := make(plotter.XYs, len(*s.graphicDots)/2)
	for i := 0; i < len(scatterData)-1; i++ {
		rn := int(normalDist.Rand())
		scatterData[i].X = (*s.graphicDots)[rn].X
		scatterData[i].Y = (*s.graphicDots)[rn].Y
	}
	createGraphicWithDot(*s.graphicDots, scatterData, namingRandomGraphic("dots_", PNG), color.RGBA{R: 255, G: 0, B: 0, A: 255})

	return nil
}

// createGraphic создаёт график с заданными точками и путём сохранения графика
func createGraphicWithDot(pts, ptsDot plotter.XYs, filePath string, color color.RGBA) {
	// Создаем новый график
	p := plot.New()
	// Добавляем точки на график
	line, err := plotter.NewLine(pts)
	if err != nil {
		panic(err)
	}
	p.Add(line)

	// Создаем график точек
	scatter, err := plotter.NewScatter(ptsDot)
	if err != nil {
		panic(err)
	}
	// Устанавливаем цвет точек
	scatter.Color = color // Красный цвет
	p.Add(scatter)

	// Сохраняем график в файл
	if err := p.Save(4*vg.Inch, 4*vg.Inch, filePath); err != nil {
		panic(err)
	}
}

// CreateRandomGraphic return Settings type for created graphic
//
// sinCount is count of multiply sin()
//
//	/*good number for*/ sinCount is 5
//
// dots is number of points on the graph
//
//	/*good number for*/ dots is 300
func CreateRandomGraphic(sinCount int, dots int) Settings {
	if sinCount == 0 {
		sinCount = 5
	}
	// Создаем новый график
	p := plot.New()

	pts := make(plotter.XYs, dots)
	for i := range pts {
		x := float64(i) * 2 * math.Pi / float64(dots)
		y := math.Sin(x)
		for i := 0; i < int(sinCount); i++ {
			y *= math.Sin(x+float64(rand.Intn(sinCount)+1)) * math.Cos(x+float64(rand.Intn(sinCount)+1))
		}
		pts[i].X = x
		pts[i].Y = y
	}

	Settings := Settings{
		graphicName: namingRandomGraphic("sin_", PNG),
		Range:       [2]int{25, 75},
		DotsCount:   dots / 4,
		graphicDots: &pts,
	}
	if pth, err := os.Getwd(); err == nil {
		Settings.path = pth
	} else {
		Settings.path = filepath.Join(".", Settings.graphicName)
		fmt.Println(err)
	}

	// Добавляем точки на график
	line, err := plotter.NewLine(pts)
	if err != nil {
		panic(err)
	}
	p.Add(line)

	// Сохраняем график в файл
	if err := p.Save(4*vg.Inch, 4*vg.Inch, Settings.graphicName); err != nil {
		panic(err)
	}

	return Settings
}

func namingRandomGraphic(prefix, extension string) string {
	max := 1
	//re, err := regexp.Compile("^sin_.*.png$")
	re, err := regexp.Compile(fmt.Sprintf("^%s.*%s$", prefix, extension))
	if err != nil {
		return prefix + "err" + extension
	}
	reNum, err := regexp.Compile(fmt.Sprintf(`^%s(\d+)\%s$`, prefix, extension))
	if err != nil {
		return prefix + "err" + extension
	}
	err = filepath.WalkDir(".", func(path string, d fs.DirEntry, err error) error {
		if re.MatchString(d.Name()) && d != nil {
			matches := reNum.FindStringSubmatch(d.Name())
			//Лексиграфическое сравнение строк
			if len(matches) >= 2 {
				num, err := strconv.Atoi(matches[1])
				if err != nil {
					return nil
				}
				if num >= max {
					max = num + 1
				}
			}
		}
		return nil
	})
	if err != nil {
		return prefix + "err" + extension
	}
	return fmt.Sprintf("%s%d%s", prefix, max, extension)
}
