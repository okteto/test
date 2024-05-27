#!/bin/sh

set -e

name=$1
namespace=$2
file=$3
deploy=$4
nocache=$5
variables=$6
timeout=$7
log_level=$8
tests=$9

if [ -n "$OKTETO_CA_CERT" ]; then
   echo "Custom certificate is provided"
   echo "$OKTETO_CA_CERT" > /usr/local/share/ca-certificates/okteto_ca_cert.crt
   update-ca-certificates
fi

params=""

if [ -n "$name" ]; then
   params="$params --name $name"
fi

if [ -n "$namespace" ]; then
   params="$params --namespace $namespace"
fi

if [ -n "$file" ]; then
   params="$params -f $file"
fi

if [ "$deploy" = "true" ]; then
      params="$params --deploy"
fi

if [ "$nocache" = "true" ]; then
      params="$params --no-cache"
fi

variable_params=""
if [ ! -z "${variables}" ]; then
  for ARG in $(echo "${variables}" | tr ',' '\n'); do
    variable_params="${variable_params} --var ${ARG}"
  done

  params="${params} ${variable_params}"
fi


if [ -n "$timeout" ]; then
   params="$params --timeout $timeout"
fi

params="$params $tests"

if [ ! -z "$log_level" ]; then
  if [ "$log_level" = "debug" ] || [ "$log_level" = "info" ] || [ "$log_level" = "warn" ] || [ "$log_level" = "error" ] ; then
    log_level="--log-level ${log_level}"
  else
    echo "unsupported log-level ${log_level}, supported options are: debug, info, warn, error"
    exit 1
  fi
fi

# https://docs.github.com/en/actions/monitoring-and-troubleshooting-workflows/enabling-debug-logging
# https://docs.github.com/en/actions/learn-github-actions/variables#default-environment-variables
if [ "${RUNNER_DEBUG}" = "1" ]; then
  log_level="--log-level debug"
fi

echo running: okteto test $log_level "$params"
# shellcheck disable=SC2086
okteto test $log_level $params