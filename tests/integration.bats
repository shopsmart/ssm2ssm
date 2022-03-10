#!/usr/bin/env bats

load helpers/ssm.bash

function setup() {
  export TEST_SSM2SSM_EXECUTABLE="${TEST_SSM2SSM_EXECUTABLE:-$BATS_TEST_DIRNAME/../bin/ssm2ssm}"
  ssm_setup

  export REGIONS_FILE="$BATS_TEST_DIRNAME/regions.txt"
}

function teardown() {
  ssm_teardown
}

function teardown_suite() {
  ssm_teardown_suite
}

@test "it should print help if no args are provided" {
  run "$TEST_SSM2SSM_EXECUTABLE"

  echo "$output"

  [ $status -eq 0 ]
}

@test "it should print help if only one arg is provided" {
  run "$TEST_SSM2SSM_EXECUTABLE" /aws/service/global-infrastructure/regions

  echo "$output"

  [ $status -eq 0 ]
}

@test "it should print the version" {
  run "$TEST_SSM2SSM_EXECUTABLE" --version

  echo "$output"

  [ $status -eq 0 ]
}

@test "it should copy the ssm params from one path to another" {
  # Delete any existing parameters
  ssm_clean_params "$TEST_SSM_PREFIX"

  run "$TEST_SSM2SSM_EXECUTABLE" /aws/service/global-infrastructure/regions "$TEST_SSM_PREFIX"

  [ $status -eq 0 ]

  echo "$output"

  ssm_pull_params "$TEST_SSM_PREFIX" "$BATS_TEST_TMPDIR/actual.txt"
  sed 's|/aws/service/global-infrastructure/regions|'"$TEST_SSM_PREFIX"'|' "$REGIONS_FILE" | sort > "$BATS_TEST_TMPDIR/expected.txt"

  diff -y "$BATS_TEST_TMPDIR/expected.txt" "$BATS_TEST_TMPDIR/actual.txt"
}

@test "it should cowardly refuse to copy the ssm params from one path to another" {
  # Copy over the parameters ahead of time
  ssm_copy_params /aws/service/global-infrastructure/regions "$TEST_SSM_PREFIX"

  run "$TEST_SSM2SSM_EXECUTABLE" /aws/service/global-infrastructure/regions "$TEST_SSM_PREFIX"

  [ $status -ne 0 ]

  echo "$output"
}

@test "it should copy the ssm params from one path to another overwriting the existing params" {
  # Delete any existing parameters
  ssm_copy_params /aws/service/global-infrastructure/regions "$TEST_SSM_PREFIX"

  # Change one of the values to know we overwrote its value
  aws ssm put-parameter \
    --name "$TEST_SSM_PREFIX/us-east-1" \
    --value not-us-east-1 \
    --type SecureString \
    --overwrite

  run "$TEST_SSM2SSM_EXECUTABLE" /aws/service/global-infrastructure/regions "$TEST_SSM_PREFIX" --overwrite

  [ $status -eq 0 ]

  echo "$output"

  ssm_pull_params "$TEST_SSM_PREFIX" "$BATS_TEST_TMPDIR/actual.txt"
  sed 's|/aws/service/global-infrastructure/regions|'"$TEST_SSM_PREFIX"'|' "$REGIONS_FILE" | sort > "$BATS_TEST_TMPDIR/expected.txt"

  diff -y "$BATS_TEST_TMPDIR/expected.txt" "$BATS_TEST_TMPDIR/actual.txt"
}
