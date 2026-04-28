package main

import (
    "encoding/json"
    "log"
    "net/http"
    "os"
    "strings"

    huh "charm.land/huh/v2"
)

func create(args []string) {
    if len(args) < 1 {
        log.Fatal("Path required.")
    }
    pwd := args[0] + "/.devcontainer"

    def_repo := "mcr.microsoft.com/devcontainers/"
    urls := get_container_urls("https://api.github.com/repos/devcontainers/images/git/trees/HEAD?recursive=1", def_repo)

    var compose bool; var repo string
    form := huh.NewForm(
        huh.NewGroup(
            huh.NewSelect[bool]().
            Title("Which build file?").
            Options(
                huh.NewOption("Dockerfile", false),
                huh.NewOption("Docker Compose", true),
            ).
            Value(&compose),
        ),
        huh.NewGroup(
            huh.NewSelect[string]().
            Title("Which base devcontainer?").
            Options(urls...).
            Value(&repo),
        ),
    )
    if err := form.Run(); err != nil {
        log.Fatal(err)
    }

    repo_name := strings.TrimPrefix(repo, def_repo)
    def_tag, tag_opts := get_container_tags(repo_name, "https://mcr.microsoft.com/v2/devcontainers/")
 
    var use_def bool
    form = huh.NewForm(
        huh.NewGroup(
            huh.NewConfirm().
            Title("Use the default version? (" + def_tag + ")").
            Affirmative("Yes").
            Negative("No").
            Value(&use_def),
        ),
    )
    if err := form.Run(); err != nil {
        log.Fatal(err)
    }

    var version string
    if use_def { version = def_tag } else {
        form = huh.NewForm(
            huh.NewGroup(
                huh.NewSelect[string]().
                Title("Which version?").
                Options(tag_opts...).
                Value(&version),
            ),
        )
        if err := form.Run(); err != nil {
            log.Fatal(err)
        }
    }

    if err := os.Mkdir(pwd, 0750); err != nil && !os.IsExist(err) {
        log.Fatal(err)
    }

    var dc_link string
    if compose {
        os.WriteFile(pwd + "/compose.yml", []byte("services:\n    main:\n        image: " + repo + ":" + version + "\n        command: sleep infinity"), 0666)
        dc_link = "{ \"dockerComposeFile\": \"compose.yml\", \"service\": \"main\" }"
    } else {
        os.WriteFile(pwd + "/Dockerfile", []byte("FROM " + repo + ":" + version), 0666)
        dc_link = "{ \"build\": { \"dockerfile\": \"Dockerfile\" \"context\": \"..\" } }"
    }
    os.WriteFile(pwd + "/devcontainer.json", []byte(dc_link), 0666)
    // DC_RUN="$(gum input --placeholder="Enter initial command...")"
}

func get_container_urls(repo string, image_repo string) []huh.Option[string] {
    var err error
    var request *http.Request; var response *http.Response
    if request, err = http.NewRequest("GET", repo, nil); err != nil {
        log.Fatal(err)
    }
    client := &http.Client{}
    if response, err = client.Do(request); err != nil {
        log.Fatal(err)
    }
    defer response.Body.Close()

    var dc_dict struct {
        Tree []struct {
            Path string `json:"path"`
            Mode string `json:"mode"`
        } `json:"tree"`
    }
    if err = json.NewDecoder(response.Body).Decode(&dc_dict); err != nil {
        log.Fatal(err)
    }
    urls := make([]huh.Option[string], 0)
    for _, item := range dc_dict.Tree {
        if strings.HasPrefix(item.Path, "src/") && item.Mode == "040000" {
            var name string
            if name = strings.TrimPrefix(item.Path, "src/"); !strings.Contains(name, "/") {
                urls = append(urls, huh.NewOption(
                    name,
                    image_repo + name,
                ))
            }
        }
    }
    return urls
}

func get_container_tags(image string, repo string) (string, []huh.Option[string]) {
    var request *http.Request; var response *http.Response; var err error
    if request, err = http.NewRequest("GET", repo + image + "/tags/list", nil); err != nil {
        log.Fatal(err)
    }
    client := &http.Client{}
    if response, err = client.Do(request); err != nil {
        log.Fatal(err)
    }
    defer response.Body.Close()

    var tag_dict struct {
        Tags []string `json:"tags"`
    }
    if err = json.NewDecoder(response.Body).Decode(&tag_dict); err != nil {
        log.Fatal(err)
    }
    def_tag := tag_dict.Tags[0]
    tag_opts := make([]huh.Option[string], 0)
    for _, tag := range tag_dict.Tags {
        tag_opts = append(tag_opts, huh.NewOption(tag, tag))
    }
    return def_tag, tag_opts
}
