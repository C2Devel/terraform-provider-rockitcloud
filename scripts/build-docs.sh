#!/bin/bash


# Path of the directory with the documentation
SOURCE_DIR="./site/"
WEBSITE_DOCS_DIR="website/docs"

# Function to show help menu
help () {
    echo "
    --help  Show this message and exit.
    --tools Install all required dependencies.
    --push  Build and push internal documentation to the bucket.
    --local Build and run internal documentation locally.
    --clean Remove temp files excluding the backup folder."
}

# Function to install required dependencies
tools () {
  pip3 install mkdocs mkdocs-material mkdocs-awesome-pages-plugin | pip install mkdocs mkdocs-material mkdocs-awesome-pages-plugin
}

# Function to build internal documentation
mkdocs_build () {
    # Copying of configuration for the documentatiion to the website folder
    cp -r docs/c2/mkdocs/assets docs/c2/mkdocs/images website/docs/
    cp docs/c2/mkdocs/mkdocs.yml website/

    # Generation of the documentation
    cd website
    echo "Generation of the documentation"
    mkdocs build -q --clean
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
    local dir="docs"

    for item in "$dir"/*; do
      base_item=$(basename "$item")

      if [ "$base_item" != "r" ] && [ "$base_item" != "d" ] && [ "$base_item" != "index.md" ]; then
        echo "Removing: webste/$item"
        rm -rf "$item"
      fi
    done
    echo "Removing of the temp files"
    rm -fr $dir/images $dir/assets mkdocs.yml site
    echo "Cleaning is complete!"
    echo "The backup folder is not removed! Remove it manually."
}

# Function to restore doc's sources from backup
restore_docs_source_cod () {
    cp -r ../backup/r docs/
    cp -r ../backup/d docs/
}

# Function to create folder's structure
process_files() {
  local src_dir="$1"
  local dest_dir_name="$2"

  find "$src_dir" -type f -name "*.md" | while read -r filepath; do
    # Taking of the subcategory from .md files
    subcategory=$(grep '^subcategory:' "$filepath" | cut -d ':' -f 2 | xargs)

    if [[ -n "$subcategory" ]]; then
      # Creating of the target folder if it doesn't exist
      dest_dir="$WEBSITE_DOCS_DIR/$subcategory/$dest_dir_name"
      mkdir -p "$dest_dir"

      # Removing of the type of the page
      if [[ $dest_dir_name == "Data Sources" ]]; then
        sed -i 's/Data Source://g' $filepath
      elif [[ $dest_dir_name == "Resources" ]]; then
        sed -i 's/Resource://g' $filepath
      fi

      # Moving of the pages to the target folder
      cp "$filepath" "$dest_dir/"
    fi
  done
}

#Function to make a backup of the r and d directories
backup () {
  echo "Backup of the $WEBSITE_DOCS_DIR/r and $WEBSITE_DOCS_DIR/d directories..."
  mkdir -p backup
  cp -r $WEBSITE_DOCS_DIR/r backup/
  cp -r $WEBSITE_DOCS_DIR/d backup/
}

if [[ "$1" == "--push" ]]; then
    if [[ -n $S3_DOCS_BUCKET_NAME ]]; then
        backup
        # Moving of the files from the r folder to the Resources folder
        echo "Moving of the .md files from $WEBSITE_DOCS_DIR/r to appropriate folders..."
        process_files "$WEBSITE_DOCS_DIR/r" "Resources"
        # Moving of the files from the d folder to the Data Sources folder
        echo "Moving of the .md files from $WEBSITE_DOCS_DIR/d to appropriate folders..."
        process_files "$WEBSITE_DOCS_DIR/d" "Data Sources"
        rm -fr $WEBSITE_DOCS_DIR/r
        rm -fr $WEBSITE_DOCS_DIR/d
        mkdocs_build
        upload_other_files
        modify_content_type css text/css
        modify_content_type js application/javascript
        restore_docs_source_cod
        remove_temp_files
    else
        echo "Define S3_DOCS_BUCKET_NAME environment variable."
        exit 1
    fi
elif [[ "$1" == "--local" ]]; then
    backup
    # Moving of the files from the r folder to the Resources folder
    echo "Moving of the .md files from $WEBSITE_DOCS_DIR/r to appropriate folders..."
    process_files "$WEBSITE_DOCS_DIR/r" "Resources"
    # Moving of the files from the d folder to the Data Sources folder
    echo "Moving of the .md files from $WEBSITE_DOCS_DIR/d to appropriate folders..."
    process_files "$WEBSITE_DOCS_DIR/d" "Data Sources"
    rm -fr $WEBSITE_DOCS_DIR/r
    rm -fr $WEBSITE_DOCS_DIR/d
    mkdocs_build
    mkdocs_run_locally
    restore_docs_source_cod
    remove_temp_files
elif [[ "$1" == "--clean" ]]; then
    remove_temp_files
elif [[ "$1" == "--tools" ]]; then
    tools
elif [[ "$1" == "--help" ]]; then
    help
else
    echo "Choose one of the options: --push, --local, --tools, --help or --clean"
    exit 1
fi