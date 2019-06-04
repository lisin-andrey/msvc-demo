package cmd

import (
	"bytes"
	"path"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/lisin-andrey/msvc-demo/entity-service/pkg/model"
	"github.com/stretchr/testify/assert"
)

func prepare(t *testing.T) *webData {
	wd := &webData{}

	const relTemplateDir = "./templates/"

	_, file, _, _ := runtime.Caller(0)
	dir := filepath.Dir(file)
	templateDir := path.Join(path.Dir(dir), relTemplateDir)
	t.Logf("\ntemplateDir: %s\n", templateDir)

	err := wd.initTemplatesFromDir(templateDir)
	// t.Logf("\n%#v\n", *wd)
	assert.NoError(t, err, "init templates")
	return wd
}

func TestTemplateErr(t *testing.T) {
	wd := prepare(t)
	data := struct{ Message string }{Message: "Some err message"}
	var tpl bytes.Buffer
	err := wd.tmpltErr.Execute(&tpl, data)
	// err := wd.tmpltErr.Execute(os.Stdout, data)
	assert.NoError(t, err, "template tmpltErr")
	// t.Log("\n", tpl.String())
}

func TestTemplateTest(t *testing.T) {
	wd := prepare(t)
	data := struct{ Message string }{Message: "Some err message"}
	var tpl bytes.Buffer
	err := wd.tmpltTest.Execute(&tpl, data)
	assert.NoError(t, err, "template tmpltTest")
	// t.Log("\n", tpl.String())
}

func TestTemplateCreate(t *testing.T) {
	wd := prepare(t)
	data := struct{}{}
	var tpl bytes.Buffer
	err := wd.tmpltCreate.Execute(&tpl, data)
	assert.NoError(t, err, "template tmpltCreate")
	// t.Log("\n", tpl.String())
}

func TestTemplateUpdate(t *testing.T) {
	wd := prepare(t)
	data := model.EntityQuery{
		ID:           1,
		Name:         "$Name1$",
		Descr:        "$Descr1$",
		Created:      time.Now(),
		LastUpdated:  time.Now(),
		LastOperator: "$User1$",
	}

	var tpl bytes.Buffer
	err := wd.tmpltUpdate.Execute(&tpl, data)
	assert.NoError(t, err, "template tmpltUpdate")
	content := tpl.String()

	//t.Log("\n", content)
	assert.True(t, strings.Contains(content, "$Name1$"), "Name is not found")
	assert.True(t, strings.Contains(content, "$Descr1$"), "Descr is not found")
	assert.True(t, strings.Contains(content, "$User1$"), "LastOperator is not found")
}

func TestTemplateEntity(t *testing.T) {
	wd := prepare(t)
	data := model.EntityQuery{
		ID:           1,
		Name:         "$Name1$",
		Descr:        "$Descr1$",
		Created:      time.Now(),
		LastUpdated:  time.Now(),
		LastOperator: "$User1$",
	}

	var tpl bytes.Buffer
	err := wd.tmpltEntity.Execute(&tpl, data)
	assert.NoError(t, err, "template tmpltEntity")
	content := tpl.String()
	// t.Log("\n", content)
	assert.True(t, strings.Contains(content, "$Name1$"), "Name is not found")
	assert.True(t, strings.Contains(content, "$Descr1$"), "Descr is not found")
	assert.True(t, strings.Contains(content, "$User1$"), "LastOperator is not found")
}

// Test entities: [2]EntityQuery, msgs: [0]string,
func TestTemplateEntitiesV1(t *testing.T) {
	wd := prepare(t)
	entities := []model.EntityQuery{
		model.EntityQuery{
			ID:           1,
			Name:         "$Name1$",
			Descr:        "$Descr1$",
			Created:      time.Now(),
			LastUpdated:  time.Now(),
			LastOperator: "$User1$",
		},
		model.EntityQuery{
			ID:           1,
			Name:         "$Name2$",
			Descr:        "$Descr2$",
			Created:      time.Now(),
			LastUpdated:  time.Now(),
			LastOperator: "$User2$",
		},
	}
	msgs := []string{}

	data := struct {
		Entities []model.EntityQuery
		Messages []string
	}{
		Entities: entities,
		Messages: msgs,
	}

	var tpl bytes.Buffer
	err := wd.tmpltEntities.Execute(&tpl, data)
	assert.NoError(t, err, "template tmpltEntities with 2 EntityQueries")
	content := tpl.String()
	//t.Log("\n", content)
	assert.True(t, strings.Contains(content, "$Name1$"), "Name1 is not found")
	assert.True(t, strings.Contains(content, "$Descr1$"), "Descr1 is not found")
	assert.True(t, strings.Contains(content, "$User1$"), "LastOperator1 is not found")
	assert.True(t, strings.Contains(content, "$Name2$"), "Name2 is not found")
	assert.True(t, strings.Contains(content, "$Descr2$"), "Descr2 is not found")
	assert.True(t, strings.Contains(content, "$User2$"), "LastOperator2 is not found")
}

// Test entities: [0]EntityQuery, msgs: [2]string,
func TestTemplateEntitiesV2(t *testing.T) {
	wd := prepare(t)
	entities := []model.EntityQuery{}
	msgs := []string{"$msg1$", "$msg2$"}

	data := struct {
		Entities []model.EntityQuery
		Messages []string
	}{
		Entities: entities,
		Messages: msgs,
	}

	var tpl bytes.Buffer
	err := wd.tmpltEntities.Execute(&tpl, data)
	assert.NoError(t, err, "template tmpltEntities with 2 messages")
	content := tpl.String()
	//t.Log("\n", content)
	assert.True(t, strings.Contains(content, "$msg1$"), "messge 1 is not found")
	assert.True(t, strings.Contains(content, "$msg2$"), "message 2 is not found")
}
