package mapper

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/accurics/terrascan/pkg/utils"
	"github.com/stretchr/testify/assert"
)

func TestCFTMapper(t *testing.T) {
	tests := []struct {
		name          string
		template      string
		expectedError bool
	}{
		{
			name:          "test-for-valid-json",
			template:      "cft/ecs-service/deploy.json",
			expectedError: false,
		},
	}
	m := NewMapper("cft")
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			doc, err := iacDocumentFromFile(test.template)
			if err != nil {
				t.Error(err)
			}
			d, _ := m.Map(doc)
			fmt.Println("+++++++++++++++++++++++++++")
			b, err := json.MarshalIndent(d, "", "    ")
			if err != nil {
				fmt.Println(err)
			}
			fmt.Println(string(b))
			fmt.Println("+++++++++++++++++++++++++++")
			if err != nil {
				assert.True(t, test.expectedError)
				return
			}
			assert.False(t, test.expectedError)
		})
	}
}

func iacDocumentFromFile(name string) (*utils.IacDocument, error) {
	data, err := readFile(name)
	if err != nil {
		return nil, err
	}
	return &utils.IacDocument{
		Type:      utils.JSONDoc,
		StartLine: 0,
		EndLine:   183,
		FilePath:  filepath.Join("test_data", name),
		Data:      data,
	}, nil
}

func readFile(name string) ([]byte, error) {
	const testData = "test_data"
	f, err := os.Open(filepath.Join(testData, name))
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return ioutil.ReadAll(f)
}
