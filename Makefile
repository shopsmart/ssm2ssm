#!/usr/bin/env make

test-fixtures: tests/regions.txt

.PHONY: tests/regions.txt # Forces rebuild
tests/regions.txt:
	aws ssm get-parameters-by-path \
		--path /aws/service/global-infrastructure/regions \
		--no-cli-pager \
		--query 'sort_by(Parameters, &Name)[].{key:Name,value:Value}' | \
		jq -r '.[] | "\(.key)=\(.value)"' \
	> tests/regions.txt
