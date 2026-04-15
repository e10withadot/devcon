package main

import (
    "os"
    "log"
)

func create(args []string) {
    // mkdir -p .devcontainer
    err := os.Mkdir(".devcontainer", 0750)
    if err != nil && !os.IsExist(err) {
        log.Fatal(err)
    }
    /* 
    * cp ~/.devcontainer.json .devcontainer/devcontainer.json
    * repo="$(curl -s \
    *     -H "Accept: application/vnd.github+json" \
    *     "https://api.github.com/repos/devcontainers/images/git/trees/HEAD?recursive=1" |
    * jq -r '.tree[].path' |
    * awk -F/ '$1=="src" && NF>=2 { print $2 }' |
    * sort -u |
    * sed 's/^https:\/\/mcr.microsoft.com\/devcontainers\///' |
    * xargs gum choose)"
    * version=$(curl -s "https://mcr.microsoft.com/v2/devcontainers/${repo}/tags/list" | jq -r '.tags[]' | fzf -q latest)
    * DC_RUN="$(gum input --placeholder="Enter initial command...")"
    * printf "FROM %s\nRUN %s\n" "mcr.micrsoft.com/devcontainers/$repo:$version" "$DC_RUN" > .devcontainer/Dockerfile
    */
}
