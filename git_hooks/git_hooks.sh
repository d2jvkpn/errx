#!/bin/bash
set -eu -o pipefail # -x
_wd=$(pwd); _path=$(dirname $0 | xargs -i readlink -f {})

cd ${_path}

set -x

for f in $(ls * | grep -v git_hooks.sh); do
    ln -frs $f ${_wd}/.git/hooks/
done
