package bender

import (
	"gopkg.in/yaml.v2"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type Plans map[string][]Plan

// Top level of test case structure
type Plan struct {
	Description     string
	Sets            []string
	Profile         string
	Host            string
	AdditionalKargs string
	Provision       string // pxe | libvirt | ISO
}

type Set struct {
	// Description string
	Libs     []string
	Filename string
	Timeout  string
}

type TestScripts struct {
	Libs    SimpleSet
	Scripts [][2]string
}

func (p Plans) GetAllPlans() {
	currentPath, _ := os.Getwd()

	err1 := os.Chdir(PlanPath)
	defer os.Chdir(currentPath)

	if err1 != nil {
		log.Fatalln(err1)
	}

	fl, _ := filepath.Glob("*.yaml")
	for _, fn := range fl {
		var plans []Plan
		v := SplitMultiYamlToSingle(filepath.Join(PlanPath, fn))
		for _, i := range v {
			var plan Plan
			yaml.Unmarshal(i, &plan)
			plans = append(plans, plan)
		}
		p[strings.Replace(fn, ".yaml", "", 1)] = plans
	}
}

func (p Plan) ParseAllSets() *TestScripts {
	testScripts := TestScripts{
		Libs:    make(SimpleSet),
		Scripts: make([][2]string, 0, 100),
	}
	for _, set := range p.Sets {
		log.Printf("Start parsing set %s", set)
		v := SplitMultiYamlToSingle(filepath.Join(SetPath, set)+".yaml", 1)

		for idx, i := range v {
			var set Set
			yaml.Unmarshal(i, &set)
			if idx == 0 {
				for _, lib := range set.Libs {
					testScripts.Libs.Add(lib)
				}
			} else {
				tmp := [2]string{}
				tmp[0], tmp[1] = set.Filename, set.Timeout
				testScripts.Scripts = append(testScripts.Scripts, tmp)
			}
		}
	}
	return &testScripts
}
