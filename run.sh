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

# allow either https or ssh repo url format:
#   git@github.com:<user>/<repo>
#   https://github.com/<user>/<repo>
githubrepo="$(echo "${giturl}" | sed -e 's/^git@github.com:/@/' -e 's/^https:\/\/github.com\//@/' -e 's/\.git$//' | sed -n 's/^@//p')"

if [ -n "${githubrepo}" ]; then
  test-supersede "https://github.com/${githubrepo}/tree/${commit}" "https://raw.githubusercontent.com/${githubrepo}/${commit}/supersedes.txt" "${author}"
else
  echo "run.sh: Git remote doesn't seem to be a github repo: ${giturl}" >&2
  exit 64
fi
