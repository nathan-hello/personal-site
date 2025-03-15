#/bin/sh

dir=${0%/*}

chmod +x ./personal-site
./personal-site --build

doas -u reluekiss rm -rf /var/www/reluekiss.com/public

cp -r ./dist/public /var/www/reluekiss.com/public
chown -R reluekiss:reluekiss /var/www/reluekiss.com

./personal-site --serve
