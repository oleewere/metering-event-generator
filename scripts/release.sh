readlinkf(){
  # get real path on mac OSX
  perl -MCwd -e 'print Cwd::abs_path shift' "$1";
}

if [ "$(uname -s)" = 'Linux' ]; then
  METERINGP_SCRIPT_DIR="`dirname "$(readlink -f "$0")"`"
else
  METERINGP_SCRIPT_DIR="`dirname "$(readlinkf "$0")"`"
fi

METERINGP_ROOT_DIR="`dirname \"$METERINGP_SCRIPT_DIR\"`"

function print_help() {
  cat << EOF
   Usage: ./release.sh [additional options]
   --release-publish            run build and publish artifacts (on release tags)
   --release-major              create a tag with major version change (e.g.: 0.1.0 -> 1.0.0)
   --release-minor              create a tag with minor version change (e.g.: 0.1.0 -> 0.2.0)
   --release-patch              create a tag with patch version change (e.g.: 0.1.0 -> 0.1.1)
   -b, --release-build-only     create a dist from snapshot package
   -v, --version <version>      override the artifact versison
   -h, --help                   print help
EOF
}

function get_branch_name() {
  echo $(git rev-parse --abbrev-ref HEAD)
}

function get_last_release() {
  local last_rev=$(git rev-list --tags --max-count=1)
  if [[ -z "$last_rev" ]]; then
    echo "v0.0.0"
  else
    echo $(git describe --tags $last_rev)
  fi
}

function next_major_release() {
  echo $(get_version_number "$1") | awk '{split($0,a,"."); print "v"a[1]+1"."0"."0}'
}

function next_minor_release() {
  echo "$1" | awk '{split($0,a,"."); print a[1]"."a[2]+1"."0}'
}

function next_patch_release() {
  echo "$1" | awk '{split($0,a,"."); print a[1]"."a[2]"."a[3]+1}'
}

function new_branch_name() {
  echo "$1" | awk '{split($0,a,"."); print a[1]"."a[2]}'| sed -e 's/v/branch-/g'
}

function get_version_number() {
  local release_version="$1"
  echo "$release_version" | sed -e 's/v//g'
}

function get_head_tag() {
  echo $(git name-rev --tags --name-only $(git rev-parse HEAD))
}

function user_confirmation() {
  read -p "Continue (y/n)?" choice
  case "$choice" in
    y|Y|yes )
      echo "Release confirmed."
      ;;
    n|N|no )
      exit 0
      ;;
    * )
      echo "Invalid user anwser input"
      exit 1
      ;;
  esac
}

function release_and_push_new_branch() {
  local next_release="$1"
  local new_branch="$2"
  local version_number=$(get_version_number $next_release)
  git tag "$next_release"
  local release_result=$(run_release)
  git branch -D $new_branch
  git branch $new_branch $next_release
  git push origin master
  git push origin $new_branch
}

function release_and_push_actual_branch() {
  local next_release="$1"
  local actual_branch="$2"
  local version_number=$(get_version_number $next_release)
  git tag "$next_release" -m "$next_release (patch release)"
  local release_result=$(run_release)
  git push origin $actual_branch
}

function update_readme_version() {
  local next_release="$1"
  local readme_md_location="$METERINGP_ROOT_DIR/README.md"
  local version_number=$(get_version_number $next_release)
  sed -i.bak "s/METERINGP_VERSION=[[:digit:]]\.[[:digit:]]\.[[:digit:]]/METERINGP_VERSION=${version_number}/" "$readme_md_location"
  rm "$readme_md_location.bak"
  git add "$readme_md_location"
  git commit -m "Update README.md (for release version: $next_release)"
}

function release_major() {
  echo "Create major release ..."
  local branch_name=$(get_branch_name)
  local last_release=$(get_last_release)
  echo "Branch name: $branch_name"
  if [[ "$branch_name" == "master" ]]; then
    echo "Last release: $last_release"
    local next_release=$(next_major_release $last_release)
    echo "New release: $next_release"
    local new_branch=$(new_branch_name $next_release)
    echo "New branch: $new_branch"
    user_confirmation
    update_readme_version "$next_release"
    release_and_push_new_branch "$next_release" "$new_branch"
  else
    echo "Major release can be created only on master branch. Exiting ..."
    exit 0
  fi
}

function release_minor() {
  echo "Create minor release ..."
  local branch_name=$(get_branch_name)
  local last_release=$(get_last_release)
  echo "Branch name: $branch_name"
  if [[ "$branch_name" == "master" ]]; then
    echo "Last release: $last_release"
    local next_release=$(next_minor_release $last_release)
    echo "New release: $next_release"
    local new_branch=$(new_branch_name $next_release)
    echo "New branch: $new_branch"
    user_confirmation
    update_readme_version "$next_release"
    release_and_push_new_branch "$next_release" "$new_branch"
  else
    echo "Minor release can be created only on master branch. Exiting ..."
    exit 0
  fi
}

function release_patch() {
  echo "Create patch release ..."
  local branch_name=$(get_branch_name)
  local last_release=$(get_last_release)
  echo "Branch name: $branch_name"
  if [[ "$branch_name" != "master" ]]; then
    if [[ "$branch_name" != branch* ]]; then
      echo "Cannot create patch release on feature branch"
      exit 0
    fi
    echo "Last release: $last_release"
    local next_release=$(next_patch_release $last_release)
    echo "New release: $next_release"
    user_confirmation
    update_readme_version "$next_release"
    release_and_push_actual_branch "$next_release" "$branch_name"
  else
    echo "Patch release cannot be created on master branch. Exiting ..."
    exit 0
  fi
}

function run_release() {
  docker run -w /go/src/github.com/oleewere/meteringp -e GITHUB_TOKEN=$GITHUB_TOKEN --rm -v $METERINGP_ROOT_DIR/vendor/:/go/src/ -v $METERINGP_ROOT_DIR:/go/src/github.com/oleewere/meteringp goreleaser/goreleaser:latest --debug --rm-dist
}

function build_only() {
  docker run -w /go/src/github.com/oleewere/meteringp -e GITHUB_TOKEN=$GITHUB_TOKEN --rm -v $METERINGP_ROOT_DIR/vendor/:/go/src/ -v $METERINGP_ROOT_DIR:/go/src/github.com/oleewere/meteringp goreleaser/goreleaser:latest --snapshot --debug --rm-dist --skip-publish
}

function main() {

  local RELEASE_BUILD_ONLY="false"
  local RELEASE="false"

  while [[ $# -gt 0 ]]
    do
      key="$1"
      case $key in
        -b|--release-build-only)
          shift 1
          RELEASE_BUILD_ONLY="true"
        ;;
        --release-major)
          local RELEASE_MAJOR="true"
          shift 1
        ;;
        --release-minor)
          local RELEASE_MINOR="true"
          shift 1
        ;;
        --release-patch)
          local RELEASE_PATCH="true"
          shift 1
        ;;
        -h|--help)
          shift 1
          print_help
          exit 0
        ;;
        *)
          echo "Unknown option: $1"
          exit 1
        ;;
      esac
    done

  if [[ -z "$GITHUB_TOKEN" ]] ; then
    echo "Setting GITHUB_TOKEN variable is required."
    exit 1
  fi

  if [[ "$RELEASE_BUILD_ONLY" == "true" ]]; then
    build_only
  fi

  if [[ "$RELEASE_MAJOR" == "true" ]] ; then
    release_major
  fi

  if [[ "$RELEASE_MINOR" == "true" ]] ; then
    release_minor
  fi

  if [[ "$RELEASE_PATCH" == "true" ]] ; then
    release_patch
  fi
}

main ${1+"$@"}