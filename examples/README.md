# Examples

First of all you need to build the provider and link it to the examples:

```shell script
make build
make prepare-examples
```

To test examples, you need a running Hydra instance. For demonstrational purposes,
we can run Hydra locally as a Docker container:

```shell script
docker container run -i -t --rm --name hydra \
  -p 4444:4444 \
  -p 4445:4445 \
  -e LOG_LEVEL=debug \
  -e DSN=memory \
  oryd/hydra:v1.4.10 serve all --dangerous-force-http
```

Then you can simply run examples as usual terraform project. Enter a particular example directory (e.g. `cd examples/oryhydra_oauth2_client`)
and run:

```shell script
terraform init

terraform plan

terraform apply
```
