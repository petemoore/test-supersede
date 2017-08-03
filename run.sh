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
git push
sleep 2
author="$(git log --pretty=format:%ae -1)"
commit="$(git rev-parse HEAD)"
giturl="$(git remote get-url origin)"
githubrepo="$(echo "${giturl}" | sed -n 's/^git@github.com://p' | sed 's/\.git$//')"

if [ -n "${githubrepo}" ]; then
  test-supersede "https://github.com/${githubrepo}" "https://raw.githubusercontent.com/${githubrepo}/${commit}/supersedes.txt" "${author}"
else
  echo "run.sh: Git remote doesn't seem to be a github repo: ${giturl}" >&2
  exit 64
fi
