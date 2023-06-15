package example

import "fmt"

type (
	UseCase interface {
		Enact() error
	}

	useCase struct {
	}
)

func NewUseCase() *useCase {
	return &useCase{}
}

func (u useCase) Enact() error {
	fmt.Println("Application Run Successfully..")
	return nil
}
