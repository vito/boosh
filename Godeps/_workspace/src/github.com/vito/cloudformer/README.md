# usage

bootstrapping a microbosh and using it to deploy drone:

```bash
pushd micro/
  spiff merge infrastructure-template.yml | cloudformer -o infrastructure.yml
  spiff merge deployment-template.yml infrastructure.yml > deployment.yml
  bosh micro deployment deployment.yml
  bosh micro deploy
popd

pushd drone/
  spiff merge infrastructure-template.yml | cloudformer -o infrastructure.yml
  spiff merge deployment-template.yml infrastructure.yml > deployment.yml
  bosh target <microbosh>
  bosh deployment deployment.yml
  # ...
popd
```
