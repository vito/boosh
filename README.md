usage:

```bash
# generate deployment
spiff merge infrastructure-template.yml > infrastructure.yml

# generate cloudformation template
go run *.go infrastructure.yml > cloudformation.json
```

and then deploy that cloudformation via the console

visualizing the infrastructure:

```bash
cat cloudformation.json | ./cfviz | dot -Tsvg -o drone.svg
```

you'll probably want to delete your stack after; this costs ~$40/mo.
