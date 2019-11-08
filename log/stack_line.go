package log

import (
	"fmt"
	"runtime"
	"bytes"
	"strings"
	"io/ioutil"
)

// The maximum number of stackframes on any error.
var MaxStackDepth = 50

type StackFrame struct {
	// The path to the file containing this ProgramCounter
	File string
	// The LineNumber in that file
	LineNumber int
	// The Name of the function that contains this ProgramCounter
	Name string
	// The Package that contains this function
	Package string
	// The underlying ProgramCounter
	ProgramCounter uintptr
}

func GetLineDetail() string {
	stackMax := make([]uintptr, MaxStackDepth)
	length := runtime.Callers(2, stackMax[:])
	stack := stackMax [:length]

	frames := make([]StackFrame, len(stack))
	for i, pc := range stack {
		frames[i] = NewStackFrame(pc)
	}

	buf := bytes.Buffer{}
	for i := len(frames)-1; i>=0 ; i= i-1 {
		buf.WriteString(frames[i].String())
	}

	return string(buf.Bytes())
}

func packageAndName(fn *runtime.Func) (string, string) {
	name := fn.Name()
	pkg := ""
	// we first remove the path prefix if there is one.
	if lastslash := strings.LastIndex(name, "/"); lastslash >= 0 {
		pkg += name[:lastslash] + "/"
		name = name[lastslash+1:]
	}
	if period := strings.Index(name, "."); period >= 0 {
		pkg += name[:period]
		name = name[period+1:]
	}

	name = strings.Replace(name, "·", ".", -1)
	return pkg, name
}

// Func returns the function that contained this frame.
func (frame *StackFrame) Func() *runtime.Func {
	if frame.ProgramCounter == 0 {
		return nil
	}
	return runtime.FuncForPC(frame.ProgramCounter)
}

// SourceLine gets the line of code (from File and Line) of the original source if possible.
func (frame *StackFrame) SourceLine() (string, error) {
	data, err := ioutil.ReadFile(frame.File)

	if err != nil {
		return "", err
	}

	lines := bytes.Split(data, []byte{'\n'})
	if frame.LineNumber <= 0 || frame.LineNumber >= len(lines) {
		return "???", nil
	}
	// -1 because line-numbers are 1 based, but our array is 0 based
	return string(bytes.Trim(lines[frame.LineNumber-1], " \t")), nil
}

// String returns the stackframe formatted in the same way as go does
// in runtime/debug.Stack()
func (frame *StackFrame) String() string {
	str := fmt.Sprintf("%s:%d \t", frame.File, frame.LineNumber)

	source, err := frame.SourceLine()
	if err != nil {
		return str + "\n"
	}

	return str + fmt.Sprintf("\t%s: %s\n", frame.Name, source)
}

// NewStackFrame popoulates a stack frame object from the program counter.
func NewStackFrame(pc uintptr) (frame StackFrame) {
	frame = StackFrame{ProgramCounter: pc}
	if frame.Func() == nil {
		return
	}
	frame.Package, frame.Name = packageAndName(frame.Func())

	// pc -1 because the program counters we use are usually return addresses,
	// and we want to show the line that corresponds to the function call
	frame.File, frame.LineNumber = frame.Func().FileLine(pc - 1)
	return

}