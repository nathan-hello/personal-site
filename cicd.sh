#/bin/sh

dir=${0%/*}

chmod +x ./personal-site
./personal-site --build-only

doas -u reluekiss

rm -rf /var/www/reluekiss.com/public

cp -r ./dist/public /var/www/reluekiss.com/public

./personal-site
