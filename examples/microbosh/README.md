# usage

```bash
# configure AWS
export AWS_ACCESS_KEY_ID=xxx
export AWS_SECRET_ACCESS_KEY=xxx
export AWS_DEFAULT_REGION=us-east-1

# generate a keypair and import it as 'bosh'
mkdir -p micro/
ssh-keygen -t rsa -f micro/id_rsa_bosh -N ''

# deploy a MicroBOSH
./deploy

# feel free to make changes to the template and simply re-run ./deploy.
# everything is idempotent and will converge.
```
