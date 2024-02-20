package cmd

// general
var fps int
var FpsDefault int = 10

var filePath string
var FilePathDefault string = "goFractalsOutput"

var width int
var WidthDefault int = 600

var height int
var HeightDefault int = 600

// IFS variables
var ifsPath string
var IfsPathDefault string = "path/to/my.ifs"

var numIterations int
var NumIterationsDefault int = 1

var algorithmProbabilistic bool
var AlgorithmProbabilisticDefault bool = false

var algorithmDeterministic bool
var AlgorithmDeterministicDefault bool = true

var probabilitiesList []float64
var ProbabilitiesListDefault []float64 = []float64{}

var random bool
var RandomDefault bool = false

var numTransforms int
var NumTransformsDefault int = 2

var numPoints int
var NumPointsDefault = 1

var numStacks int
var NumStacksDefault int = 1

var thickness float32
var ThicknessDefault float32 = 15

// julia & mandelbrot
var juliaEquation string
var JuliaEquationDefault string = "z*z - 1"

var colored bool
var ColoredDefault bool = false

var threeDimensional bool
var ThreeDimensionalDefault bool = false

var cInitString string
var CInitStringDefault = "0+0i"

var cIncrementString string
var CIncrementStringDefault string = "0.01+.01i"

var numIncrements int
var NumIncrementsDefault int = 10

var writeBinary bool
var WriteBinaryDefault bool = false

var solid bool
var SolidDefault bool = false

var mandelbrotEquation string
var MandelbrotEquationDefault string = "z*z + c"

var centerPointString string
var CenterPointStringDefault string = "0+0i"

var zoom float64
var ZoomDefault float64 = 4

var maxItr int
var MaxItrDefault int = 1000