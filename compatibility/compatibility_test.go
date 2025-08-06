package compatibility

import (
	"fmt"
	"testing"
)

func TestRequestScore(t *testing.T) {
	content := []string{
		"i love food! it tastes so wonderful. mmmm",
		"food is not amazing, it disappointed me although i liked the alarm clock.",
		"this food is absolutely terrible. everyone slapped me and i cried.",
	}
	for _, c := range content {
		fmt.Println(c)
		m, err := queryTextClassificationApi(c)
		if err != nil {
			t.Fatal(err)
		}
		fmt.Println(m)
	}
	m, err := queryTextGenerationApi(content[0])
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(m)
}
