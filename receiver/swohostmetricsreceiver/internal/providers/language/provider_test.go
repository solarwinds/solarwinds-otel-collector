package language

import (
	"fmt"
	"testing"
)

func Test_Functional(t *testing.T) {
	t.Skip("This test should be run manually")

	sut := CreateLanguageProvider()
	result := <-sut.Provide()
	fmt.Printf("Result: %+v\n", result)
}
