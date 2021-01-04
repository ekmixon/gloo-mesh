#!/usr/bin/env bash

###################################################################################
# This script generates a versioned docs website for Gloo Mesh and
# deploys to Firebase.
###################################################################################

set -ex

# Update this array with all supported versions of GlooMesh to include in the versioned docs website.
go run codegen/versions/versionsgen.go
readarray -t versions < <(jq -r '.supported_versions[]' docs/version.json)
latestVersion=$(jq -r '.latest_version' docs/version.json)

# Firebase configuration
firebaseJson=$(cat <<EOF
{
  "hosting": {
    "site": "gloo-mesh",
    "public": "public",
    "ignore": [
      "firebase.json",
      "**/.*",
      "**/node_modules/**"
    ],
    "rewrites": [
      {
        "source": "/",
        "destination": "/gloo-mesh/latest/index.html"
      },
      {
        "source": "/gloo-mesh",
        "destination": "/gloo-mesh/latest/index.html"
      }
    ]
  }
}
EOF
)

# This script assumes that the working directory is the root of the repo.
workingDir=$(pwd)
docsSiteDir=$workingDir/ci/docs-site
repoDir=$workingDir/ci/gloo-mesh-temp

mkdir "$docsSiteDir"
echo "$firebaseJson" > "$docsSiteDir/firebase.json"

git clone https://github.com/solo-io/gloo-mesh.git "$repoDir"

# install go tools to sub-repo
export PATH=$workingDir/_output/.bin:$PATH
make -C "$repoDir" install-go-tools

# Generates a data/Solo.yaml file with $1 being the specified version.
function generateHugoVersionsYaml() {
  yamlFile="$repoDir/docs/data/Solo.yaml"
  # Truncate file first.
  {
    echo "LatestVersion: $latestVersion";
    # /gloo-mesh prefix is needed because the site is hosted under a domain name with suffix /gloo-mesh
    echo "DocsVersion: /gloo-mesh/$1";
    echo "CodeVersion: $1";
    echo "DocsVersions:";
  } > "$yamlFile"

  for hugoVersion in "${versions[@]}"; do
    echo "  - $hugoVersion" >> "$yamlFile"
  done
}

for version in "${versions[@]}"; do
  echo "Generating site for version $version"
  cd "$repoDir"
  if [[ "$version" == "main" ]]; then
    git checkout main
  else
    git checkout "tags/$version"
  fi
  # Replace version with "latest" if it's the latest version. This enables URLs with "/latest/..."
  if [[ "$version" ==  "$latestVersion" ]]; then
    version="latest"
  fi
  go run codegen/docs/docsgen.go
  cd docs
  # Generate data/Solo.yaml file with version info populated.
  generateHugoVersionsYaml $version
  # Use nav bar as defined in main, not the checked out temp repo.
  cp -f "$workingDir/docs/layouts/partials/versionnavigation.html" layouts/partials/versionnavigation.html
  # Generate the versioned static site.
  make site-release
  # Copy over versioned static site to firebase content folder.
  mkdir -p "$docsSiteDir/public/gloo-mesh/$version"
  cp -a site-latest/. "$docsSiteDir/public/gloo-mesh/$version/"
  # Discard git changes, vendor_any, and the generated changelog for subsequent checkouts
  cd "$repoDir"
  git reset --hard
  rm -rf vendor_any docs/content/reference/changelog
done
