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
    pwd := args[0] + "/.devcontainer"

    // curl -s -H "Accept: application/vnd.github+json" "https://api.github.com/repos/devcontainers/images/git/trees/HEAD?recursive=1"
    var err error
    var request *http.Request
    client := &http.Client{}
    if request, err = http.NewRequest("GET", "https://api.github.com/repos/devcontainers/images/git/trees/HEAD?recursive=1", nil); err != nil {
        log.Fatal(err)
    }
    request.Header.Add("Accept", "application/vnd.github+json")
    var response *http.Response
    if response, err = client.Do(request); err != nil {
        log.Fatal(err)
    }
    defer response.Body.Close()

    // jq -r '.tree[].path'
    var dc_dict struct {
        Tree []struct {
            Path string `json:"path"`
            Mode string `json:"mode"`
        } `json:"tree"`
    }
    if err := json.NewDecoder(response.Body).Decode(&dc_dict); err != nil {
        log.Fatal(err)
    }
    urls := make([]huh.Option[string], 0)
    for _, item := range dc_dict.Tree {
        if strings.HasPrefix(item.Path, "src/") && item.Mode == "040000" {
            var name string
            if name = strings.TrimPrefix(item.Path, "src/"); !strings.Contains(name, "/") {
                urls = append(urls, huh.NewOption(
                    name,
                    "mcr.microsoft.com/devcontainers/" + name,
                ))
            }
        }
    }

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
    if err = form.Run(); err != nil {
        log.Fatal(err)
    }

    // curl -s "https://mcr.microsoft.com/v2/devcontainers/${repo}/tags/list"
    repo_name := strings.TrimPrefix(repo, "mcr.microsoft.com/devcontainers/")
    if request, err = http.NewRequest("GET", "https://mcr.microsoft.com/v2/devcontainers/" + repo_name + "/tags/list", nil); err != nil {
        log.Fatal(err)
    }
    request.Header.Add("Accept", "application/vnd.github+json")
    if response, err = client.Do(request); err != nil {
        log.Fatal(err)
    }
    defer response.Body.Close()

    // jq -r '.tags[]'
    var tag_dict struct {
        Tags []string `json:"tags"`
    }
    if err := json.NewDecoder(response.Body).Decode(&tag_dict); err != nil {
        log.Fatal(err)
    }
    def_tag := tag_dict.Tags[0]
    tag_opts := make([]huh.Option[string], 0)
    for _, tag := range tag_dict.Tags {
        tag_opts = append(tag_opts, huh.NewOption(tag, tag))
    }
 
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
    if err = form.Run(); err != nil {
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
        if err = form.Run(); err != nil {
            log.Fatal(err)
        }
    }

    // mkdir -p .devcontainer
    if err = os.Mkdir(pwd, 0750); err != nil && !os.IsExist(err) {
        log.Fatal(err)
    }

    // printf "FROM %s\nRUN %s\n" "mcr.micrsoft.com/devcontainers/$repo:$version" "$DC_RUN" > .devcontainer/Dockerfile
    var dc_link string
    if compose {
        os.WriteFile(pwd + "/compose.yml", []byte("services:\n\tmain:\n\t\timage: " + repo + ":" + version), 0666)
        dc_link = "{ \"dockerComposeFile\": \"compose.yml\"\"service\": \"main\" }"
    } else {
        os.WriteFile(pwd + "/Dockerfile", []byte("FROM " + repo + ":" + version), 0666)
        dc_link = "{ \"build\": { \"dockerfile\": \"Dockerfile\" \"context\": \"..\" } }"
    }
    os.WriteFile(pwd + "/devcontainer.json", []byte(dc_link), 0666)
    // DC_RUN="$(gum input --placeholder="Enter initial command...")"
}
