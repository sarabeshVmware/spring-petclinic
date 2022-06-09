# TAP Meta-Package Readme #

## Top-Level Keys ##

Top Level Keys (TLK) allow us to cascade common values down into multiple packages after specifiying only once. These can be useful in order to reduce the size and complexity of the values files necessary for installation.

### Current Top-Level Keys ###

- `shared.ca_cert_data` - Optional: PEM Encoded certificate data to trust TLS connections with a private CA.


### Process to propose a new TLK ###

Create an MR with the following proposed changes:

1. OpenAPI definition of the TLK under the `shared` key-space. You'll need to do this in the [template for the TAP package](../packages/tap/template/package.yaml) in the section similar to the below where you see the `ca_cert_data` key contributed.

```yaml
  shared:
    type: object
    properties:
    ca_cert_data:
        type: string
        default: ""
        description: "Optional: PEM Encoded certificate data to trust TLS connections with a private CA."
    description: "Common properties shared across multiple components"
```

2. Add your requested key to the `values.yml` for tap-pkg [here](./config/values.yml) as show below:

```yaml
shared:
  ca_cert_data: ""
```

3. Submit your PR and leave it in a pending state until you're able to solicit feedback from TAP Technical Leads. Notify the #tap-unboxing-and-installation slack channel of your intention in order to solicit reviews.


### Inherit a TLK into your package ###

There are two examples of methods to inherit TLK data into your package. You can see the TBS team's implementation [here](./config/tbs.yaml). You can also see the Convention Controller's implementation [here](./config/convention-controller.yaml)
