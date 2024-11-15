#!/bin/bash

# Load new versions of specified terraform provider from official terraform registry 
# to private s3 registry.

set -euo pipefail

TF_REGISTRY_URL=${TF_REGISTRY_URL:-"https://registry.terraform.io/"}
S3_REGISTRY_URL=${S3_REGISTRY_URL:-}
S3_BUCKET_NAME=${S3_BUCKET_NAME:-}
PROVIDER_NAME=${PROVIDER_NAME:-"c2devel/rockitcloud"}

S3_BACKUP_DIR=${S3_BACKUP_DIR:-}
TMP_DIR="/tmp/"

TF_VERSIONS_FILE="${TMP_DIR}/tf-versions.json"
S3_VERSIONS_FILE="${TMP_DIR}/s3-versions.json"

function trim_slashes() {
  # $1 - url part

  echo "${1}" | sed 's:^/*::;s:/*$::'
}

function to_lower() {
  # $1 - str

  echo "${1}" | tr '[:upper:]' '[:lower:]'
}

function curl_and_check() {
  # $1 - svc name
  # $2 - url
  # $3 - output

  if [[ -z "${3:-}" ]]; then
    http_code=$(curl -k -s -o /dev/null -w '%{http_code}' "${2}")
  else
    http_code=$(curl -k -s -o "${3}" -w '%{http_code}' "${2}")
  fi

  if [[ "${http_code}" == "200" ]]; then
    echo "  ${1}... OK (endpoint: ${2}, http code: ${http_code})"
  else
    echo "  ${1}... FAIL (endpoint: ${2}, http code: ${http_code})"
    exit 1
  fi
}


echo "Start updating s3 registry"


echo "Check env variables"

if [[ -z "${S3_REGISTRY_URL}" || -z "${S3_BUCKET_NAME}" ]]; then
  echo "  S3_REGISTRY_URL and S3_BUCKET_NAME must not be empty"
  exit 1
fi

echo "  TF_REGISTRY_URL = ${TF_REGISTRY_URL}"
echo "  S3_REGISTRY_URL = ${S3_REGISTRY_URL}"
echo "  S3_BUCKET_NAME = ${S3_BUCKET_NAME}"
echo "  PROVIDER_NAME = ${PROVIDER_NAME}"
echo "  S3_BACKUP_DIR = ${S3_BACKUP_DIR}"


TF_REGISTRY_URL=$(trim_slashes "${TF_REGISTRY_URL}")
S3_REGISTRY_URL=$(trim_slashes "${S3_REGISTRY_URL}")

PROVIDER_NAME=$(to_lower "${PROVIDER_NAME}")


echo "Check availability of registries:"

curl_and_check \
  "terraform registry" \
  "${TF_REGISTRY_URL}/.well-known/terraform.json"

curl_and_check \
  "s3 registry" \
  "${S3_REGISTRY_URL}/.well-known/terraform.json"


echo "Get providers url prefix:"

tf_provider_prefix=$(curl -k -s "${TF_REGISTRY_URL}/.well-known/terraform.json" | jq -r '."providers.v1"')
s3_provider_prefix=$(curl -k -s "${S3_REGISTRY_URL}/.well-known/terraform.json" | jq -r '."providers.v1"')

echo "  terraform registry: ${tf_provider_prefix}"
echo "  s3 registry: ${s3_provider_prefix}"

tf_provider_prefix=$(trim_slashes "${tf_provider_prefix}")
s3_provider_prefix=$(trim_slashes "${s3_provider_prefix}")


echo "Get versions for provider '${PROVIDER_NAME}':"

curl_and_check \
  "tf registry versions" \
  "${TF_REGISTRY_URL}/${tf_provider_prefix}/${PROVIDER_NAME}/versions" \
  "${TF_VERSIONS_FILE}"

curl_and_check \
  "s3 registry versions" \
  "${S3_REGISTRY_URL}/${s3_provider_prefix}/${PROVIDER_NAME}/versions" \
  "${S3_VERSIONS_FILE}"

tf_provider_versions=$(< "${TF_VERSIONS_FILE}" jq -r '.versions[] | .version')
s3_provider_versions=$(< "${S3_VERSIONS_FILE}" jq -r '.versions[] | .version')

echo "  terraform registry:" $tf_provider_versions
echo "  s3 registry:" $s3_provider_versions


if [[ -n "${S3_BACKUP_DIR}" ]]; then
  timestamp=$(date +%Y%m%d-%H%M%S)

  echo "Backup s3 bucket to ${S3_BACKUP_DIR}/${S3_BUCKET_NAME}-${timestamp}/"

  mkdir -p "${S3_BACKUP_DIR}"

  s3cmd sync --config="./.s3cfg" --no-preserve --quiet "s3://${S3_BUCKET_NAME}/" "${S3_BACKUP_DIR}/${S3_BUCKET_NAME}-${timestamp}/"

  echo "Finish backup"
fi


echo "Find new provider versions in terraform registry"

new_versions_count=0
for version in $tf_provider_versions; do
  if [[ "${s3_provider_versions[*]}" =~ ${version} ]]; then
    continue
  fi

  echo "  Add new version '${version}' to s3 registry"

  tf_version_platforms=$(< "${TF_VERSIONS_FILE}" \
    jq --arg version "$version" '.versions[] | select(.version == $version) | .platforms')

  platforms_len=$(echo "${tf_version_platforms}" | jq '. | length')

  i=0
  while [[ $i -lt $platforms_len ]]; do
    os=$(echo "${tf_version_platforms}" | jq --argjson i $i -r '.[$i] | .os')
    arch=$(echo "${tf_version_platforms}" | jq --argjson i $i -r '.[$i] | .arch')

    curl_and_check \
      "  ${os}/${arch}" \
      "${TF_REGISTRY_URL}/${tf_provider_prefix}/${PROVIDER_NAME}/${version}/download/${os}/${arch}" \
      "${TMP_DIR}/${version}_${os}_${arch}.json"

    s3cmd put --dry-run --config="./.s3cfg" --quiet --acl-public --content-type=application/json "${TMP_DIR}/${version}_${os}_${arch}.json" \
      "s3://${S3_BUCKET_NAME}/${s3_provider_prefix}/${PROVIDER_NAME}/${version}/download/${os}/${arch}/index.json"

    rm -f "${TMP_DIR}/${version}_${os}_${arch}.json"

    ((i+=1))
  done

  ((new_versions_count+=1))
  echo "  Finish adding new version '${version}'"
done

if [[ $new_versions_count -gt 0 ]]; then
  echo "${new_versions_count} version(s) were added"

  echo "Update versions meta in s3 registry"

  s3cmd put --dry-run --config="./.s3cfg" --quiet --acl-public --content-type=application/json "${TF_VERSIONS_FILE}" \
    "s3://${S3_BUCKET_NAME}/${s3_provider_prefix}/${PROVIDER_NAME}/versions/index.json"

  echo "Finish versions meta update"
else
  echo "No new versions were found"
fi

rm -f "${TF_VERSIONS_FILE:?}"
rm -f "${S3_VERSIONS_FILE:?}"

echo "Done!"
