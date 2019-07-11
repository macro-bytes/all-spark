package template_reader

import (
	"template"
	"testing"
)

func TestDeserialize(t *testing.T) {
	var template template.DockerTemplate

	err := Deserialize("does-not-exist", &template)
	if err == nil {
		t.Error("Expected non-nil error")
	}

	err = Deserialize("../../../sample_templates/docker.json", &template)
	if err != nil {
		t.Error(err)
	}
}
