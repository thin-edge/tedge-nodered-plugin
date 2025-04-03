set dotenv-load
set export

DEVICE_ID := env_var_or_default("DEVICE_ID", "CI_" + file_name(home_directory()) + "_tedge-nodered-plugin" )
IMAGE := env_var_or_default("IMAGE", "debian-systemd-native")
IMAGE_SRC := env_var_or_default("IMAGE_SRC", "debian-systemd-native")

# Initialize a dotenv file for usage with a local debugger
# WARNING: It will override any previously generated dotenv file
init-dotenv:
  @echo "Recreating .env file..."
  @echo "DEVICE_ID=$DEVICE_ID" > .env
  @echo "IMAGE=$IMAGE" >> .env
  @echo "IMAGE_SRC=$IMAGE_SRC" >> .env
  @echo "C8Y_BASEURL=$C8Y_BASEURL" >> .env
  @echo "C8Y_USER=$C8Y_USER" >> .env
  @echo "C8Y_PASSWORD=$C8Y_PASSWORD" >> .env


# Release all artifacts
build *ARGS='':
    mkdir -p output
    go run main.go completion bash > output/completions.bash
    go run main.go completion zsh > output/completions.zsh
    go run main.go completion fish > output/completions.fish

    docker context use default
    goreleaser release --clean --auto-snapshot {{ARGS}}

# Build a release locally (for testing the release artifacts)
build-local:
    just -f "{{justfile()}}" build --snapshot

# Install python virtual environment
venv:
  [ -d .venv ] || python3 -m venv .venv
  ./.venv/bin/pip3 install -r tests/requirements.txt

# Build test images and test artifacts
build-test:
  docker buildx build --load -t {{IMAGE}} -f ./test-images/{{IMAGE_SRC}}/Dockerfile .

# Run tests
test *args='':
  ./.venv/bin/python3 -m robot.run --outputdir output {{args}} tests

# Download/update vendor packages
update-vendor:
  go mod vendor

# Print yocto licensing string
print-yocto-licenses:
  @echo 'LIC_FILES_CHKSUM = " \'
  @find . -name "LICENSE*" -exec bash -c 'echo -n "    file://src/\${GO_IMPORT}/{};md5=" && md5 -q {}' \; 2>/dev/null | grep -v "/\.venv/" | sed 's|$| \\|g' | sed 's|/\./|/|g'
  @echo '"'
