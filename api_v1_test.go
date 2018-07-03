package dynalist

import (
	"fmt"
	"io/ioutil"
	"testing"
	"time"
)

func TestAPI(t *testing.T) {
	api, err := New()
	if err != nil {
		t.Fatal(err)
	}
	var sampleID string
	t.Run("FileList", func(t *testing.T) {
		res, err := api.FileList()
		if err != nil {
			t.Error(err)
		}
		ioutil.WriteFile("testdata/file-list-"+time.Now().String(),
			[]byte(fmt.Sprintf("%+v\n", res)), 0644,
		)
		for _, file := range res.Files {
			if file.Type == TypeFolder {
				sampleID = file.ID
				break
			}
		}
	})
	t.Run("FileEdit", func(t *testing.T) {
		change := NewChange(ActionEdit)
		change.Type = TypeFolder
		change.FileID = sampleID
		change.Title = "test: " + time.Now().String()
		changes := []*Change{change}
		res, err := api.FileEdit(changes)
		if err != nil {
			t.Error(err)
		}
		ioutil.WriteFile("testdata/file-edit-"+time.Now().String(),
			[]byte(fmt.Sprintf("%+v\n", res)), 0644,
		)
		if res.Code != CodeOK {
			t.Log(res.Code)
			t.Error("response code is not Ok")
		}
	})
}
