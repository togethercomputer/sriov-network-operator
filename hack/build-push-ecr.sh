#!/usr/bin/env bash
#
# build-push-ecr.sh — multi-arch (linux/amd64 + linux/arm64) build & publish of
# the Together sriov-network-operator fork images to AWS ECR. (TCL-9012)
#
# Builds all three images and pushes them to the public ECR gallery:
#   public.ecr.aws/k6t4m3l7/sriov-network-operator:<TAG>
#   public.ecr.aws/k6t4m3l7/sriov-network-operator-config-daemon:<TAG>
#   public.ecr.aws/k6t4m3l7/sriov-network-operator-webhook:<TAG>
# then copies the operator image into the private substrate-apps repository
# (registry-to-registry via `docker buildx imagetools create`, no rebuild):
#   651706779278.dkr.ecr.us-west-2.amazonaws.com/substrate-apps/sriov-network-operator:<PRIVATE_TAG>
# and finally verifies both architectures are present in every pushed manifest.
#
# PREREQUISITES:
#   1. AWS SSO session in the prod account:
#        export AWS_PROFILE=TC-Prod-PowerUserAccess-651706779278
#        aws sso login
#   2. The public ECR repository for the OPERATOR image may not exist yet
#      (config-daemon and webhook repos already exist in the k6t4m3l7 gallery).
#      If the first push fails with "repository does not exist", create it once:
#        aws ecr-public create-repository --repository-name sriov-network-operator --region us-east-1
#   3. Docker Buildx with a docker-container driver builder for multi-arch
#      (QEMU/binfmt for arm64 when on an amd64 host):
#        docker buildx create --use
#
# USAGE:
#   ./hack/build-push-ecr.sh            # build, push, copy, verify
#   DRY_RUN=1 ./hack/build-push-ecr.sh  # print every command, execute nothing
#
# OVERRIDES (environment):
#   AWS_PROFILE      (default: TC-Prod-PowerUserAccess-651706779278)
#   TAG              (default: git rev-parse --short=8 HEAD)
#   PRIVATE_TAG      (default: v1.4.0-0.0.2 — bump per release;
#                     last prod private pin was v1.4.0-0.0.1)
#   PUBLIC_REGISTRY  (default: public.ecr.aws/k6t4m3l7)
#   PRIVATE_REGISTRY (default: 651706779278.dkr.ecr.us-west-2.amazonaws.com)
#
# NOTE: run from anywhere; the script cd's to the repo root itself.
# Rollout after publish = pin the new tags in tcloud-infra (versions.env +
# per-region cluster.env + sriov-network-operator app values) — separate PR.

set -euo pipefail

cd "$(dirname "$0")/.."

AWS_PROFILE="${AWS_PROFILE:-TC-Prod-PowerUserAccess-651706779278}"
export AWS_PROFILE
TAG="${TAG:-$(git rev-parse --short=8 HEAD)}"
PRIVATE_TAG="${PRIVATE_TAG:-v1.4.0-0.0.2}"
PUBLIC_REGISTRY="${PUBLIC_REGISTRY:-public.ecr.aws/k6t4m3l7}"
PRIVATE_REGISTRY="${PRIVATE_REGISTRY:-651706779278.dkr.ecr.us-west-2.amazonaws.com}"
PRIVATE_REPO="${PRIVATE_REGISTRY}/substrate-apps"
DRY_RUN="${DRY_RUN:-0}"

# Dockerfile:image-name pairs; image names are the public ECR repository names.
SPECS=(
  "Dockerfile:sriov-network-operator"
  "Dockerfile.sriov-network-config-daemon:sriov-network-operator-config-daemon"
  "Dockerfile.webhook:sriov-network-operator-webhook"
)

run() {
  echo "+ $*"
  if [[ "${DRY_RUN}" != "1" ]]; then
    "$@"
  fi
}

verify_arches() {
  local ref="$1"
  echo "+ docker buildx imagetools inspect ${ref}"
  if [[ "${DRY_RUN}" == "1" ]]; then
    return 0
  fi
  local out platform
  out="$(docker buildx imagetools inspect "${ref}")"
  echo "${out}"
  for platform in linux/amd64 linux/arm64; do
    if ! grep -q "Platform:.*${platform}" <<<"${out}"; then
      echo "ERROR: ${ref} is missing platform ${platform}" >&2
      return 1
    fi
  done
  echo "OK: ${ref} has both linux/amd64 and linux/arm64"
}

echo "==> Logging in to ECR (AWS_PROFILE=${AWS_PROFILE})"
if [[ "${DRY_RUN}" == "1" ]]; then
  echo "+ aws ecr-public get-login-password --region us-east-1 | docker login --username AWS --password-stdin public.ecr.aws"
  echo "+ aws ecr get-login-password --region us-west-2 | docker login --username AWS --password-stdin ${PRIVATE_REGISTRY}"
else
  aws ecr-public get-login-password --region us-east-1 \
    | docker login --username AWS --password-stdin public.ecr.aws
  aws ecr get-login-password --region us-west-2 \
    | docker login --username AWS --password-stdin "${PRIVATE_REGISTRY}"
fi

for SPEC in "${SPECS[@]}"; do
  DF="${SPEC%%:*}"
  IMG="${SPEC##*:}"
  if [[ ! -f "${DF}" ]]; then
    echo "ERROR: ${DF} not found (run from the repo root)" >&2
    exit 1
  fi
  echo "==> Building ${PUBLIC_REGISTRY}/${IMG}:${TAG} from ${DF} for linux/amd64,linux/arm64"
  # Multi-arch is explicit and unconditional here: this script always --push'es,
  # so building both platforms is correct. (A local `docker build` without
  # --push cannot load multi-platform results — `make image` therefore stays
  # single-arch; see IMAGE_BUILD_OPTS in the Makefile.)
  run docker buildx build --platform linux/amd64,linux/arm64 -f "${DF}" \
    -t "${PUBLIC_REGISTRY}/${IMG}:${TAG}" --push .
done

echo "==> Copying operator image to private substrate-apps as ${PRIVATE_TAG}"
run docker buildx imagetools create \
  -t "${PRIVATE_REPO}/sriov-network-operator:${PRIVATE_TAG}" \
  "${PUBLIC_REGISTRY}/sriov-network-operator:${TAG}"

echo "==> Verifying both architectures are present in the pushed manifests"
for SPEC in "${SPECS[@]}"; do
  IMG="${SPEC##*:}"
  verify_arches "${PUBLIC_REGISTRY}/${IMG}:${TAG}"
done
verify_arches "${PRIVATE_REPO}/sriov-network-operator:${PRIVATE_TAG}"

echo "==> Done. Next: pin the new tags in tcloud-infra (versions.env / cluster.env / app values) — separate PR."
