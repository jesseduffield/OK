#!sh

set -e

cd ../site
npm run build
cd ../pulumi/site
pulumi up
