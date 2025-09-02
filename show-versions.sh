set -eu

# Header
printf "%-20s %-30s %-20s %-20s\n" "CONTAINER" "IMAGE" "DIGEST" "VERSION"

# Loop over running containers
for cid in $(docker ps -q); do
    name=$(docker inspect --format '{{.Name}}' "$cid" | sed 's,^/,,')
    image=$(docker inspect --format '{{.Config.Image}}' "$cid")
    digest=$(docker inspect --format '{{.Image}}' "$cid")
    version=$(docker inspect --format '{{ index .Config.Labels "org.opencontainers.image.version" }}' "$cid" 2>/dev/null || echo "-")

    # Clean empty version
    [ -z "$version" ] && version="-"

    printf "%-20s %-30s %-20s %-20s\n" "$name" "$image" "$digest" "$version"
done