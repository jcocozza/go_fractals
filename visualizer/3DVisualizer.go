package visualizer

import (
	"fmt"
	"log"
	"runtime"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

var (
	program uint32
	vao     uint32
)

var (
    rotationX float32
    rotationY float32
    rotationZ float32
)

const width, height, scaler = 1000,1000,5

func normalizePointsList(pointsList [][]float64, scalingFactor float32) []float32 {
	var normalizedPointsList []float32
	for _, point := range pointsList {
		for _, coord := range point {
			normalizedPointsList = append(normalizedPointsList, (scalingFactor * float32(coord)))
		}
	}
	return normalizedPointsList
}

var keys = make(map[glfw.Key]bool)
// Key callback function
func keyCallback(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
    if action == glfw.Press {
        keys[key] = true
    } else if action == glfw.Release {
        keys[key] = false
    }
}
// Draw a set of 3D points
func Draw3D(pointsList [][]float64) {
	runtime.LockOSThread()

	points := normalizePointsList(pointsList, scaler)

	if err := glfw.Init(); err != nil {
		log.Fatal(err)
	}
	defer glfw.Terminate()

	glfw.WindowHint(glfw.OpenGLForwardCompatible,glfw.True);
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)

    glfw.WindowHint(glfw.Resizable, glfw.True)
    glfw.WindowHint(glfw.ContextVersionMajor, 3)
    glfw.WindowHint(glfw.ContextVersionMinor, 2)
	window, err := glfw.CreateWindow(width, height, "3D Plot", nil, nil)
	if err != nil {
		log.Fatal(err)
	}
	window.MakeContextCurrent()
	// Handle keyboard input
	window.SetKeyCallback(keyCallback)

	if err := gl.Init(); err != nil {
		log.Fatal(err)
	}
	// Vertex and fragment shader code
	vertexShaderSource := `
		#version 410
		layout(location = 0) in vec3 vp;
		uniform mat4 modelView; // Add a uniform for the modelView matrix
		void main() {
			gl_Position = modelView * vec4(vp, 3.0);
			gl_PointSize = 1.0; // Adjust the point size as needed
		}
	`
	fragmentShaderSource := `
		#version 410
		out vec4 frag_color;
		void main() {
			frag_color = vec4(0.0, 1.0, 0.0, 1.0); // Point color (green in this case)
		}
	`

	// Create shaders
	vertexShader, err := compileShader(vertexShaderSource, gl.VERTEX_SHADER)
	if err != nil {
		log.Fatal(err)
	}
	defer gl.DeleteShader(vertexShader)

	fragmentShader, err := compileShader(fragmentShaderSource, gl.FRAGMENT_SHADER)
	if err != nil {
		log.Fatal(err)
	}
	defer gl.DeleteShader(fragmentShader)

	// Create shader program
	program = gl.CreateProgram()
	gl.AttachShader(program, vertexShader)
	gl.AttachShader(program, fragmentShader)
	gl.LinkProgram(program)

	// Create vertex array object (VAO)
	gl.GenVertexArrays(1, &vao)
	gl.BindVertexArray(vao)

	// Main loop
	for !window.ShouldClose() {
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		// Define the vertex data for the axes
		axisVertices := []float32{
			// X-axis
			0.0, 0.0, 0.0, // Start point
			5.0, 0.0, 0.0, // End point

			// Y-axis
			0.0, 0.0, 0.0, // Start point
			0.0, 5.0, 0.0, // End point

			// Z-axis
			0.0, 0.0, 0.0, // Start point
			0.0, 0.0, 5.0, // End point
		}

		// Create a VBO for the axes
		var vboAxes uint32
		gl.GenBuffers(1, &vboAxes)
		gl.BindBuffer(gl.ARRAY_BUFFER, vboAxes)
		gl.BufferData(gl.ARRAY_BUFFER, len(axisVertices)*4, gl.Ptr(axisVertices), gl.STATIC_DRAW)

		// Set up vertex array object (VAO) for the axes
		var vaoAxes uint32
		gl.GenVertexArrays(1, &vaoAxes)
		gl.BindVertexArray(vaoAxes)

		// Specify the layout of the vertex data
		positionAttribAxes := uint32(gl.GetAttribLocation(program, gl.Str("vp\x00")))
		gl.EnableVertexAttribArray(positionAttribAxes)
		gl.BindBuffer(gl.ARRAY_BUFFER, vboAxes)
		gl.VertexAttribPointer(positionAttribAxes, 3, gl.FLOAT, false, 0, nil)

		gl.UseProgram(program)

		// Draw x-axis (red)
		gl.Uniform3f(gl.GetUniformLocation(program, gl.Str("color\x00")), 1.0, 0.0, 0.0)
		gl.BindVertexArray(vaoAxes)
		gl.DrawArrays(gl.LINES, 0, 2)

		// Draw y-axis (green)
		gl.Uniform3f(gl.GetUniformLocation(program, gl.Str("color\x00")), 0.0, 1.0, 0.0)
		gl.DrawArrays(gl.LINES, 2, 2)

		// Draw z-axis (blue)
		gl.Uniform3f(gl.GetUniformLocation(program, gl.Str("color\x00")), 0.0, 0.0, 1.0)
		gl.DrawArrays(gl.LINES, 4, 2)

		// Inside the main loop
		modelView := mgl32.Ident4()
		modelView = modelView.Mul4(mgl32.HomogRotate3D(mgl32.DegToRad(rotationX), mgl32.Vec3{1, 0, 0}))
		modelView = modelView.Mul4(mgl32.HomogRotate3D(mgl32.DegToRad(rotationY), mgl32.Vec3{0, 1, 0}))
		modelView = modelView.Mul4(mgl32.HomogRotate3D(mgl32.DegToRad(rotationZ), mgl32.Vec3{0, 0, 1}))

		// Set the modelView matrix uniform in your shader program
		modelViewUniform := gl.GetUniformLocation(program, gl.Str("modelView\x00"))
		gl.UniformMatrix4fv(modelViewUniform, 1, false, &modelView[0])

		// Inside the main loop
		if window.GetKey(glfw.KeyLeft) == glfw.Press {
			//fmt.Println("KEY PRESSED LEFT")g
			rotationY += 1.0 // Adjust the rotation angle as needed
		}
		if window.GetKey(glfw.KeyRight) == glfw.Press {
			//fmt.Println("KEY PRESSED RIGHT")
			rotationY -= 1.0 // Adjust the rotation angle as needed
		}
		if window.GetKey(glfw.KeyUp) == glfw.Press {
			//fmt.Println("KEY PRESSED UP")
			rotationX += 1.0 // Adjust the rotation angle as needed
		}
		if window.GetKey(glfw.KeyDown) == glfw.Press {
			//fmt.Println("KEY PRESSED DOWN")
			rotationX += -1.0 // Adjust the rotation angle as needed
		}
		// Add more conditions for other keys or axes

		// Create a buffer for the points
		var vbo uint32
		gl.GenBuffers(1, &vbo)
		gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
		gl.BufferData(gl.ARRAY_BUFFER, len(points)*4, gl.Ptr(points), gl.STATIC_DRAW)

		// Set the attribute pointer for the vertex shader
		positionAttrib := uint32(gl.GetAttribLocation(program, gl.Str("vp\x00")))
		gl.EnableVertexAttribArray(positionAttrib)
		gl.VertexAttribPointer(positionAttrib, 3, gl.FLOAT, false, 0, nil)

		// Draw the points
		gl.DrawArrays(gl.POINTS, 0, int32(len(points)/3))

		gl.DisableVertexAttribArray(positionAttrib)

		gl.DeleteBuffers(1, &vbo)

		window.SwapBuffers()
		glfw.PollEvents()
	}
}

func compileShader(source string, shaderType uint32) (uint32, error) {
	shader := gl.CreateShader(shaderType)
	csource, free := gl.Strs(source)
	gl.ShaderSource(shader, 1, csource, nil)
	free()
	gl.CompileShader(shader)

	var status int32
	gl.GetShaderiv(shader, gl.COMPILE_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetShaderiv(shader, gl.INFO_LOG_LENGTH, &logLength)
		logMessage := make([]byte, logLength)
		gl.GetShaderInfoLog(shader, logLength, nil, &logMessage[0])
		return 0, fmt.Errorf("compileShader: %v", string(logMessage))
	}

	return shader, nil
}
