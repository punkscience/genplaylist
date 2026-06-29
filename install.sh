#!/usr/bin/env bash
# genplaylist one-liner installer.
#
#   curl -fsSL https://punkscience.github.io/genplaylist/install.sh | bash
set -euo pipefail

GREEN='\033[0;32m'; BOLD='\033[1m'; NC='\033[0m'

if [ "$(uname -s)" != "Linux" ]; then
  echo "This installer is for Debian/Ubuntu (APT). On other platforms install genplaylist another way." >&2
  exit 1
fi

# Require root for package operations. Re-exec via sudo if needed. When piped
# (curl|bash or bash <(curl)), $0 is /dev/fd/N which sudo cannot read, so
# re-fetch the script from its canonical URL before elevating.
if [ "$(id -u)" -ne 0 ] && [ -z "${SUDO_USER:-}" ]; then
  echo -e "${BOLD}📦 genplaylist installer — Linux (APT)${NC}"
  echo "  Adds the genplaylist APT repo + signing key, then installs genplaylist (needs root)."
  if [ ! -f "$0" ] || [ "${0#/dev/fd/}" != "$0" ]; then
    TMP="$(mktemp)"
    curl -fsSL "https://punkscience.github.io/genplaylist/install.sh" -o "$TMP"
    exec sudo bash "$TMP"
  fi
  exec sudo bash "$0"
fi

echo "  → Installing signing key…"
curl -fsSL "https://punkscience.github.io/genplaylist/apt/genplaylist-archive-keyring.gpg" \
  -o /usr/share/keyrings/genplaylist-archive-keyring.gpg

echo "  → Adding APT source…"
cat > /etc/apt/sources.list.d/genplaylist.list <<SOURCELIST
deb [signed-by=/usr/share/keyrings/genplaylist-archive-keyring.gpg] https://punkscience.github.io/genplaylist/apt/ stable main
SOURCELIST

echo "  → Updating package lists…"
apt-get update -qq

echo "  → Installing genplaylist…"
apt-get install -y genplaylist

echo -e "${GREEN}✓ genplaylist installed. Run: genplaylist --help${NC}"