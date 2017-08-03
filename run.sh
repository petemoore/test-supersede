#!/bin/bash -e

cd "$(dirname "${0}")"
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
git push
sleep 2

go install
test-supersede
