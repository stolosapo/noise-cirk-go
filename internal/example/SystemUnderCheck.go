package example

var (
	systemHasProblem bool = false
)

func SystemHasProblem() bool {
	return systemHasProblem
}

func RestoreSystemProblem() {
	systemHasProblem = false
}

func MakeSystemProblematic() {
	systemHasProblem = true
}
