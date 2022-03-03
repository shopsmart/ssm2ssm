#!/usr/bin/env bash

function ssm_setup() {
  :
}

function ssm_teardown() {
  :
}

function ssm_teardown_suite() {
  ssm_clean_params "$TEST_SSM_PREFIX"
}

function ssm_copy_params() {
  local input_path="$1"
  local output_path="$2"

  aws ssm get-parameters-by-path \
    --path "$input_path" \
    --with-decryption \
    --no-cli-pager \
  | jq -r '.Parameters[] | "\(.Name) \(.Value)"' \
  | while read -r name value; do
    name="${name/"$input_path"/"$output_path"}"
    aws ssm put-parameter \
      --name "$name" \
      --type SecureString \
      --value "$value" \
      --overwrite 2>&3
  done
}

function ssm_clean_params() {
  local path="$1"

  aws ssm get-parameters-by-path \
      --path "$path" \
      --no-cli-pager \
  | jq -r '.Parameters[].Name' \
  | while read -r param; do
    echo "[INFO ] Deleting parameter: $param" >&3
    aws ssm delete-parameter \
      --name "$param" 2>&3
  done

  echo "[INFO ] Sleeping for 30 seconds to allow SSM to propagate" >&3
  sleep 30
}

function ssm_pull_params() {
  local path="$1"
  local file="$2"

  echo "[INFO ] Pulling parameters from $path" >&3
  aws ssm get-parameters-by-path \
    --path "$path" \
    --with-decryption \
    --no-cli-pager \
  | jq -r '.Parameters // [] | .[] | "\(.Name)=\(.Value)"' \
  | sort >"$file"
}
