#/bin/sh

# https://stackoverflow.com/questions/2870992/automatic-exit-from-bash-Shell-script-on-error
# -e exits on error
# -u errors on undefined variables
# -x prints commands before execution,
# and -o (for option) pipefail exits on command pipe failures
set -euxo pipefail

dir=${0%/*}

chmod +x ./personal-site
./personal-site --build

doas -u reluekiss rm -rf /var/www/reluekiss.com/public

cp -r ./dist/public /var/www/reluekiss.com/public
chown -R reluekiss:reluekiss /var/www/reluekiss.com

./personal-site --serve
