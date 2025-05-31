#!/usr/bin/env bash

set -e

VERSION=$1
IMAGE="quay.io/cbennett/todo-api:${VERSION}"
REPO="https://github.com/col1985/todo-app.git"
PVC="https://raw.githubusercontent.com/openshift/pipelines-tutorial/master/01_pipeline/03_persistent_volume_claim.yaml"

tkn pipeline start build-and-push \
  --prefix-name build-push-pipelinerun \
  -w name="shared-workspace,volumeClaimTemplateFile=${PVC}" \
  -p build_name="todo-api-${VERSION}" \
  -p context_dir=./api \
  -p repo_url="${REPO}" \
  -p IMAGE_REPO=$IMAGE \
  -p dockerfile=./Containerfile \
  -p revision=main

