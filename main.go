package main

import (
	"flag"
	"github.com/BurntSushi/toml"
	"github.com/mitchellh/go-homedir"
	"github.com/xanzy/go-gitlab"
	"log"
	"os"
	"path/filepath"
)

type Config struct {
	ApiToken          string
	DefaultVisibility string
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	visibilityFlag := flag.String("visibility", "public", "Repo visibility")
	flag.Parse()

	home, err := homedir.Dir()
	configLocation := home + "/.config/gitlab-create-repo.toml"
	check(err)

	var conf Config
	if _, err := toml.DecodeFile(configLocation, &conf); err != nil {
		log.Print("The config file couldn't be decoded. Please make sure it is valid TOML syntax.")
		check(err)
	}

	visibility := gitlab.PublicVisibility
	if conf.DefaultVisibility == "Private" {
		visibility = gitlab.PrivateVisibility
	}

	if *visibilityFlag == "private" {
		visibility = gitlab.PrivateVisibility
	}

	git := gitlab.NewClient(nil, conf.ApiToken)

	projectName, err := os.Getwd()
	if err != nil {
		check(err)
	}
	projectName = filepath.Base(projectName)

	p := &gitlab.CreateProjectOptions{
		Name:            gitlab.String(projectName),
		VisibilityLevel: gitlab.VisibilityLevel(visibility),
	}
	project, _, err := git.Projects.CreateProject(p)
	if err != nil {
		log.Fatal(err)
	}

	log.Print("Created ", project.Name)
	log.Print("Set your git remote to: ", project.SSHURLToRepo)
}
