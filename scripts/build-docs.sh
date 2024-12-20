#!/bin/bash


# Path of the directory with the documentation
SOURCE_DIR="./site/"

# Function to build internal documentation
mkdocs_build () {
    # Copying of configuration for the documentatiion to the website folder
    cp -r docs/c2/mkdocs/assets docs/c2/mkdocs/images website/docs/
    cp docs/c2/mkdocs/mkdocs.yml website/

    # Generation of the documentation
    cd website
    echo "Generation of the documentation"
    mkdocs build --clean
}

# Funtion to upload the documentation files to the bucket
upload_other_files () {
    echo "Uploading of the documentation to the $S3_DOCS_BUCKET_NAME bucket"
    s3cmd sync "$SOURCE_DIR" "s3://$S3_DOCS_BUCKET_NAME" --acl-public
}

# Funtion to upload .css files to the bucket with specified Content-Type
modify_content_type () {
    local file_type=$1
    local content_type=$2
    echo "Setting of the content-type for .$file_type files in bucket"
    find "$SOURCE_DIR" -type f -name "*.$file_type" | while read -r file; do
        s3cmd modify "s3://$S3_DOCS_BUCKET_NAME/${file#$SOURCE_DIR}" --add-header="Content-Type:$content_type"
    done
}

# Function to run internal documentation locally
mkdocs_run_locally () {
    echo "Run internal documentation locally"
    mkdocs serve
}

# Function to remove temp files
remove_temp_files () {
    echo "Removing of the temp files"
    rm -fr docs/images docs/assets mkdocs.yml site
    echo Complete!
}

if [[ "$1" == "push" ]]; then
    if [[ -n $S3_DOCS_BUCKET_NAME ]]; then
        mkdocs_build
        upload_other_files
        modify_content_type css text/css
        modify_content_type js application/javascript
        remove_temp_files
    else
        echo "Define S3_DOCS_BUCKET_NAME environment variable."
        exit 1
    fi
elif [[ "$1" == "local" ]]; then
    mkdocs_build
    mkdocs_run_locally
    remove_temp_files
elif [[ "$1" == "clean" ]]; then
    remove_temp_files
else
    echo "Choose one of the options: push, local or clean"
    exit 1
fi
