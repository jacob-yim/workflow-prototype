#!/bin/bash

set -o errexit
set -o nounset
set -o pipefail

case "$(uname -s)" in
    Linux*)     linkutil=readlink;;
    Darwin*)    linkutil=greadlink;;
    *)          machine="UNKNOWN:${unameOut}"
esac

NIRMATA_DIR=$(dirname ${BASH_SOURCE})/..
NIRMATA_ROOT=$(${linkutil} -f ${NIRMATA_DIR})

CODEGEN_PKG="${GOPATH}/src/k8s.io/code-generator"

NIRMATA_PKG=${NIRMATA_ROOT#"${GOPATH}/src/"}

${CODEGEN_PKG}/generate-groups.sh \
    "deepcopy,client" \
    ${NIRMATA_PKG}/pkg/client \
    ${NIRMATA_PKG}/pkg/api \
    "workflow:v1"