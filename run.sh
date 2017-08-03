#!/bin/bash -e

cd "$(dirname "${0}")"
go install
go get github.com/taskcluster/slugid-go/slug

cat << EOF > supersedes.txt
{
    "supersedes": [
        "$(slug)",
        "$(slug)",
        "$(slug)",
        "$(slug)",
        "$(slug)"
    ]
}
EOF

git add supersedes.txt
git commit -m "Updates slugs in supersedes.txt"
git push origin master
sleep 2

test-supersede "https://raw.githubusercontent.com/petemoore/test-supersede/$(git rev-parse HEAD)/supersedes.txt"
