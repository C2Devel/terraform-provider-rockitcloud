#!/bin/bash

# Path of the directory with the documentation
SOURCE_DIR="site/"
S3_CMD_CFG_LOCATION=${S3_CMD_CFG_LOCATION:-"$HOME/.s3cfg"}
S3_CMD_CFG_LOCATION_ROCKIT=${S3_CMD_CFG_LOCATION_ROCKIT:-"$HOME/.s3cfg-rockit"}
S3_DOCS_BUCKET_NAMES=("docs.tf.k2.cloud")
if [[ -n "${S3_DOCS_BUCKET_NAME:-}" ]]; then
    S3_DOCS_BUCKET_NAMES=("$S3_DOCS_BUCKET_NAME")
elif [[ -f "$S3_CMD_CFG_LOCATION_ROCKIT" && -r "$S3_CMD_CFG_LOCATION_ROCKIT" && -s "$S3_CMD_CFG_LOCATION_ROCKIT" ]]; then
    S3_DOCS_BUCKET_NAMES+=("docs.tf.rockitcloud.ru")
fi

help () {
    echo "
    --help  Show this message and exit.
    --tools Install all required dependencies.
    --push  Build and push internal documentation to the configured buckets (or S3_DOCS_BUCKET_NAME).
    --local Build and run internal documentation locally."
}

tools () {
  python3 -m venv .venv
  source .venv/bin/activate
  pip3 install mkdocs mkdocs-material mkdocs-awesome-pages-plugin || pip install mkdocs mkdocs-material mkdocs-awesome-pages-plugin
  deactivate
}

# Copying of configuration for the documentatiion to the website folder
copy () {
  cp -r mkdocs/images docs/mkdocs-images
  cp -r mkdocs/assets docs/mkdocs-assets
  cp mkdocs/mkdocs.yml ./
}

# Generation of the documentation
build () {
  source .venv/bin/activate
  echo "Generation of the documentation"
  mkdocs build -f ./mkdocs.yml -d $SOURCE_DIR --clean
  local build_status=$?
  deactivate
  return "$build_status"
}

# Run the documentation locally
run_local () {
  source .venv/bin/activate
  echo "Run the documentation locally"
  mkdocs serve -f ./mkdocs.yml
  deactivate
}

s3_config_for_bucket () {
    if [[ "$1" == "docs.tf.rockitcloud.ru" ]]; then
        printf '%s\n' "$S3_CMD_CFG_LOCATION_ROCKIT"
    else
        printf '%s\n' "$S3_CMD_CFG_LOCATION"
    fi
}

# Function to upload the documentation files to the bucket
upload_other_files () {
    local bucket="$1"
    local config="$2"
    s3cmd sync --config="$config" "$SOURCE_DIR" "s3://$bucket" --acl-public
}

# Function to upload .css files to the bucket with specified Content-Type
upload_css_files () {
    local bucket="$1"
    local config="$2"
    find "$SOURCE_DIR" -type f -name "*.css" | while read -r file; do
    s3cmd modify --config="$config" "s3://$bucket/${file#$SOURCE_DIR}" --add-header='Content-Type:text/css' || return 1
done
}

# Function to upload .js files to the bucket with specified Content-Type
upload_js_files () {
    local bucket="$1"
    local config="$2"
    find "$SOURCE_DIR" -type f -name "*.js" | while read -r file; do
    s3cmd modify --config="$config" "s3://$bucket/${file#$SOURCE_DIR}" --add-header='Content-Type:application/javascript' || return 1
done
}

# Removing of the temp files
cleanup () {
  echo "Removing of the temp files"
  rm -fr docs/mkdocs-images docs/mkdocs-assets mkdocs.yml site
  echo Complete!
}

# Uploading of the documentation to the S3 buckets

if [[ "$1" == "--push" ]]; then
    for bucket in "${S3_DOCS_BUCKET_NAMES[@]}"; do
        config=$(s3_config_for_bucket "$bucket")
        if [[ ! -f "$config" || ! -r "$config" || ! -s "$config" ]]; then
            echo "S3 config is missing, unreadable or empty for $bucket: $config" >&2
            exit 1
        fi
    done
    if ! copy || ! build; then
        cleanup
        exit 1
    fi
    upload_status=0
    for bucket in "${S3_DOCS_BUCKET_NAMES[@]}"; do
        config=$(s3_config_for_bucket "$bucket")
        upload_other_files "$bucket" "$config" && \
            upload_css_files "$bucket" "$config" && \
            upload_js_files "$bucket" "$config" || upload_status=1
    done
    cleanup
    exit "$upload_status"
elif [[ "$1" == "--local" ]]; then
    copy
    run_local
    cleanup
elif [[ "$1" == "--tools" ]]; then
    tools
elif [[ "$1" == "--help" ]]; then
    help
else
    echo "Choose one of the options: --push, --local, --tools or --help"
    exit 1
fi
