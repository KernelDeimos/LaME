package engine

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"

	"github.com/KernelDeimos/LaME/lamego/model"
)

type FileWalkModelProducer struct {
	sourceDirectory__ string
	modelReader__     ModelReader
	doneWalk__        bool
}

func (producer *FileWalkModelProducer) SetSourceDirectory(s string) {
	producer.sourceDirectory__ = s
}

func (producer *FileWalkModelProducer) SetModelReader(
	reader ModelReader,
) {
	producer.modelReader__ = reader
}

func (producer *FileWalkModelProducer) InvokeModels() {
	if producer.doneWalk__ {
		return
	}
	callback := func(path string, info os.FileInfo, err error) error {
		if filepath.Ext(path) == ".yaml" {
			models := []model.Model{}
			b, err := ioutil.ReadFile(path)
			if err != nil {
				return err
			}
			err = yaml.Unmarshal(b, &models)
			if err != nil {
				return err
			}
			for _, m := range models {
				producer.modelReader__.AddModel(m)
			}
		}
		return nil
	}
	err := filepath.Walk(
		producer.sourceDirectory__, callback)
	producer.doneWalk__ = true
	if err != nil {
		// TODO: ErrorProducer interface
		panic(err)
	}
}
