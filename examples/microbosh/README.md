# usage

```bash
# import keypair as 'bosh'
ssh-keygen -t rsa -f id_rsa_bosh

# configure AWS credentials
export AWS_ACCESS_KEY_ID=xxx
export AWS_SECRET_ACCESS_KEY=xxx
export AWS_DEFAULT_REGION=us-east-1

# generate cloudformation template
cat microbosh-infrastructure.yml | boosh generate > cloudformation.json

# deploy infrastructure
cat cloudformation.json | boosh deploy --name microbosh

# generate stub to be fed into deployment templates
boosh resources --name microbosh > microbosh-stub.yml

# generate microbosh deployment manifest
mkdir micro/
spiff merge microbosh-deployment.yml microbosh-stub.yml > micro/microbosh.yml
bosh micro deployment micro/microbosh.yml

# deploy microbosh!
bosh micro deploy ami-2bf3fb42
```
