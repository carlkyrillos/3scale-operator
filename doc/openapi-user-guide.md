# OpenAPI Custom Resource

The [OpenAPI CRD](openapi-reference.md) is used as the source of truth to reconcile
one [3scale Product custom resource](product-reference.md) and
one [3scale Backend custom resource](backend-reference.md).

## Table of contents

* [OpenAPI Custom Resource](#openapi-custom-resource)
   * [Table of contents](#table-of-contents)
   * [OpenAPI document sources](#openapi-document-sources)
      * [Secret OpenAPI spec source](#secret-openapi-spec-source)
      * [URL OpenAPI spec source](#url-openapi-spec-source)
      * [OpenID Connect and OAuth2 example](#openid-connect-and-oauth2-example)
   * [Supported OpenAPI spec version and limitations](#supported-openapi-spec-version-and-limitations)
     * [OpenAPI 3.0.2 limitations](#openapi-302-limitation)
     * [OpenIdConnect and OAuth2 limitations](#openidconnect-and-oauth2-limitations)
   * [OpenAPI importing rules](#openapi-importing-rules)
      * [Product name](#product-name)
      * [Private Base URL](#private-base-url)
      * [3scale Methods](#3scale-methods)
      * [3scale Mapping Rules](#3scale-mapping-rules)
      * [Authentication](#authentication)
      * [ActiveDocs](#activedocs)
      * [3scale Application Plans](#3scale-application-plans)
      * [3scale Product Policy Chain](#3scale-product-policy-chain)
      * [3scale Deployment Mode](#3scale-deployment-mode)
   * [Minimum required OAS doc](#minimum-required-oas-doc)
   * [Update behavior](#update-behavior)
   * [Link your OpenAPI spec to your 3scale tenant or provider account](#link-your-openapi-spec-to-your-3scale-tenant-or-provider-account)

Generated using [github-markdown-toc](https://github.com/ekalinin/github-markdown-toc)

## OpenAPI document sources

The OpenAPI document <OAS> can be read from different sources:
* Kubernetes secret
* URL. Supported schemes are `http` and `https`.

*Note*: Accepted OpenAPI spec document formats are `json` and `yaml`.

### Secret OpenAPI spec source

Create a secret with the OpenAPI spec document. The name of the secret object will be referenced in the OpenAPI CR.

The following example shows how to create a secret out of a file:

```yaml
$ cat myopenapi.yaml
---
openapi: "3.0.2"
info:
  title: "some title"
  description: "some description"
  version: "1.0.0"
paths:
  /pet:
    get:
      operationId: "getPet"
      responses:
        405:
          description: "invalid input"


$ oc create secret generic myopenapi --from-file myopenapi.yaml
secret/myopenapi created
```

**NOTE** The filename used as key inside the secret is not read by the operator. Only the content is read.

Then, create your OpenAPI CR providing reference to the secret holding the OpenAPI document.

```yaml
apiVersion: capabilities.3scale.net/v1beta1
kind: OpenAPI
metadata:
  name: openapi1
spec:
  openapiRef:
    secretRef:
      name: myopenapi
```

[OpenAPI CRD Reference](openapi-reference.md) for more info.

### URL OpenAPI spec source

```yaml
apiVersion: capabilities.3scale.net/v1beta1
kind: OpenAPI
metadata:
  name: openapi1
spec:
  openapiRef:
    url: "https://raw.githubusercontent.com/OAI/OpenAPI-Specification/master/examples/v3.0/petstore.yaml"
```

[OpenAPI CRD Reference](openapi-reference.md) for more info.

### OpenID Connect and OAuth2 example

3scale requires some additional information that is not included 
in the OpenAPI spec - this data needs to be provided in the OpenAPI CR.
Specifically:
- OpenID Connect Issuer Type: defaults to rest but can be overridden using the OpenAPI CR.
- OpenID Connect Issuer Endpoint Reference (Secret): 3scale requires that the issuer URL must include a client secret.
- Flows object: When the security scheme is `oauth2`, the flows are specified via the OpenAPI spec (OAS). However, for the `openIdConnect` security scheme, the OAS does not provide the flows. In that case, they can be specified using the OpenAPI CR.

#### OIDC Issuer Secret

- The below secret example assumes a previously set up Issuer Client (RHSSO/Keycloak realm and client)
- The `issuerEndpoint` format is `https://<CLIENT_ID>:<CLIENT_SECRET>@<HOST>:<PORT>/auth/realms/<REALM_NAME>"`. The format is described in `3scale Portal/Products page - AUTHENTICATION SETTINGS - OpenID Connect Issuer`.
- The `<CLIENT_SECRET>` value will be taken from the Issuer Client - `Realm/Clients/ClientID/Credentials/Secret`.

```yaml
kind: Secret
apiVersion: v1
metadata:
  name: my-secret
  namespace: 3scale-test
data:
  issuerEndpoint: https://3scale-zync:some-secret@keycloak-rhsso-test.example.com/auth/realms/petstore
type: Opaque
```

#### OpenAPI CR example for OIDC and oauth2
```yaml
apiVersion: capabilities.3scale.net/v1beta1
kind: OpenAPI
metadata:
  generation: 1
  name: openapi-example
spec:
  openapiRef:
    url: "https://example.com/petstore.yaml"
  privateAPISecretToken: "xxxx"
  oidc:
    issuerType: keycloak
    issuerEndpointRef:
      name: my-secret
    jwtClaimWithClientID: azp
    jwtClaimWithClientIDType: plain
    authenticationFlow:
      standardFlowEnabled: true
      implicitFlowEnabled: true
      serviceAccountsEnabled: true
      directAccessGrantsEnabled: true
    gatewayResponse:
      errorStatusAuthFailed: 403
```
- `oidc` is an optional field in the OpenAPI CR
- The table below describes the fields in the `oidc` block: 

| **Field**                | **Required** | **Description**                                                                                                                                                                                                                                                                                                                                                                                                                                                            |
| --- | --- |---|
| issuerType               | no           | Valid values: [keycloak, rest]. Defaults to `rest`.                                                                                                                                                                                                                                                                                                                                                                                                                        |
| issuerEndpoint           | no           | Issuer endpoint. It can be defined in `issuerEndpointRef` or as plain value (please see CR example and notes below). The format of this endpoint is determined on your OpenID Provider setup. For RHSSO:  `https://<client_id>:<client_secret>@<host>:<port>/auth/realms/<realm_name>`.                                                                                                                                                                                    |
| issuerEndpointRef        | no           | The secret that contains  `issuerEndpoint`.                                                                                                                                                                                                                                                                                                                                                                                                                                |
| jwtClaimWithClientID     | no           | JSON Web Token (JWT) Claim with ClientID that contains the clientID. Defaults to 'azp'.                                                                                                                                                                                                                                                                                                                                                                                    |
| jwtClaimWithClientIDType | no           | JwtClaimWithClientIDType sets to process the ClientID Token Claim value as a string or as a liquid template. Valid values: plain, liquid. Defaults to 'plain'.                                                                                                                                                                                                                                                                                                             |
| authenticationFlow       | no           | Flows object. When the sec scheme is `oauth2`, the flows are provided by the OpenAPI spec. However, for the `openIdConnect` security scheme, the OpenAPI spec does not provide the flows. In this case, the flows can be specified in the OpenAPI CR. There are 4 flows parameters (for OIDC only): `standardFlowEnabled`, `implicitFlowEnabled`, `serviceAccountsEnabled`, `directAccessGrantsEnabled`. See [3scale product reference](product-reference.md) for more info |
| gatewayResponse          | no | Specifies custom gateway response on errors. See `GatewayResponseSpec` in [3scale product reference ](product-reference.md) for more info.                                                                                                                                                                                                                                                                                                                                 |

- One of IssuerEndpointRef or IssuerEndpoint must be defined in OIDC Spec (both fields can be defined, see next note).
- If issuerEndpoint plain value is defined in CR - it will be used as precedence over issuerEndpointRef secret.
- The format of issuerEndpoint is determined on your OpenID Provider setup;
  see in 3scale portal - `Product/Integration/Settings/AUTHENTICATION SETTINGS/OpenID Connect Issuer`.  
- HostHeader and SecretToken should only be set at OpenApi CR .spec level instead of at the .spec.oidc level since the OIDC value is ignored.  OIDC Security is populated in the Product CR from the OpenApi CR PrivateAPISecretToken and PrivateAPIHostHeader parameters if one or both of them are defined in the OpenAPI CR. See OpenAPISpec in [openapi reference](openapi-reference.md), OIDC specification in [product-reference.md](product-reference.md).

OpenAPI CR example where issuerEndpoint defined both as plain value and in secret (plain value will be used):
```yaml
apiVersion: capabilities.3scale.net/v1beta1
kind: OpenAPI
metadata:
  generation: 1
  name: openapi-example
spec:
  openapiRef:
    url: "https://example.com/petstore.yaml"
  oidc:
    issuerType: keycloak
    issuerEndpoint: https://3scale-zync:some-secret@keycloak-rhsso-test.example.com/auth/realms/petstore
    issuerEndpointRef:
      name: my-secret
    jwtClaimWithClientID: azp
    jwtClaimWithClientIDType: plain
    authenticationFlow:
      standardFlowEnabled: true
      implicitFlowEnabled: true
      serviceAccountsEnabled: true
      directAccessGrantsEnabled: true
``` 

- If OpenAPI CR spec is OIDC but securitySchemes type in OAS is `oauth2` then the CR OIDC Authentication Flows parameters will be ignored, and Product OIDC Authentication Flows will be set to match oauth2 flows as defined in OAS:
    - StandardFlowEnabled = true if oauth2 AuthorizationCode is defined
    - ImplicitFlowEnabled = true if oauth2 Implicit is defined
    - DirectAccessGrantsEnabled = true if oauth2 Password is defined
    - ServiceAccountsEnabled = true if oauth2 ClientCredentials is defined

An example of OAS securitySchemes definition that allows selection of all Product OIDC Authentication Flows (OIDC should be defined in the OpenAPI CR)
```yaml
      securitySchemes:
        myOauth:
          description: This API uses OAuth 2 with the implicit grant flow. [More info](https://api.example.com/docs/auth)
          flows:
            password:
              scopes:
                read_pets: read your pets
                write_pets: modify pets in your account
              tokenUrl: https://api.example.com/oauth2/token
            implicit:
              authorizationUrl: https://example.com/api/oauth/dialog
              scopes:
                write_pets: modify pets in your account
                read_pets: read your pets
            authorizationCode:
              authorizationUrl: https://example.com/api/oauth/dialog
              tokenUrl: https://example.com/api/oauth/token
              scopes:
                write_pets: modify pets in your account
                read_pets: read your pets 
            clientCredentials:
              tokenUrl: https://example.com/api/oauth/token
              scopes:
                write_pets: modify pets in your account
                read_pets: read your pets           
          type: oauth2
```

## Supported OpenAPI spec version and limitations

### OpenAPI 3.0.2 limitation

* [OpenAPI __3.0.2__ specification](https://github.com/OAI/OpenAPI-Specification/blob/main/versions/3.0.2.md) with some limitations:
  * `info.title` field value must not exceed `253-38 = 215` character length. It will be used to create some openshift object names with some length [limitations](https://kubernetes.io/docs/concepts/overview/working-with-objects/names/).
  * Only first `servers[0].url` element in `servers` list parsed as *private base url*. As OpenAPI specification `basePath` property, `servers[0].url` URL's base path component will be used.
  * `servers` element in path item or operation items are not supported.
  * Just a single top level security requirement supported. Operation level security requirements not supported.
  * Supported security schemes: 
    * `apiKey`
    * `openIdConnect/oauth2`   

### OpenIdConnect and OAuth2 limitations

* As detailed in section [OpenID Connect and OAuth2 example](#openid-connect-and-oauth2-example)  - 
3scale requires some additional pieces of information that are not included in the OpenAPI spec and need to be provided in the OpenAPI CR:
  * OpenID Connect Issuer Type: defaults to rest but can be overridden using the OpenAPI CR.
  * OpenID Connect Issuer Endpoint Reference (Secret): 3scale requires that the issuer URL must include a client secret.
  * Flows object: When the security scheme is `oauth2`, the flows are specified via the OpenAPI spec (OAS). However, for the `openIdConnect` security scheme, the OAS does not provide the flows. In that case, they can be specified using the OpenAPI CR.


## OpenAPI importing rules

### Product name

The default product system name is taken from the `info.title` field in the OpenAPI definition.
However, you can override this product name using the `spec.productSystemName` field
of the [OpenAPI CRD](openapi-reference.md).

### Private Base URL

Private base URL is read from OpenAPI `servers[0].url` field.
You can override this using the `spec.privateBaseURL` field
of the [OpenAPI CRD](openapi-reference.md).

### 3scale Methods

Each OpenAPI defined operation will translate in one 3scale method at product level.
The method name is read from the [operationId](https://github.com/OAI/OpenAPI-Specification/blob/main/versions/3.0.2.md#operationObject) field of the operation object.

### 3scale Metrics

The default `Hits` metric will be assigned to the OAS-generated Product CR by default however custom metrics can be created and assigned using [OAS 3scale extensions](openapi-3scale-extensions.md#root-level-3scale-extension).

### 3scale Mapping Rules

Each OpenAPI defined operation will translate in one 3scale mapping rule at product level.
Previously existing mapping rules will be replaced by those imported from the OpenAPI.

OpenAPI [paths](https://github.com/OAI/OpenAPI-Specification/blob/main/versions/3.0.2.md#pathsObject) object provides mapping rules *Verb* and *Pattern* properties. 3scale methods will be associated accordingly to the [operationId](https://github.com/OAI/OpenAPI-Specification/blob/main/versions/3.0.2.md#operationObject).
The mapping rule can be associated with custom methods/metrics using [OAS 3scale extensions](openapi-3scale-extensions.md#operation-level-3scale-extension).

*Delta* value is hard-coded to `1`. This can be set to an arbitrary value using [OAS 3scale extensions](openapi-3scale-extensions.md#operation-level-3scale-extension)

By default, *Strict matching* policy is being configured.
Matching policy can be switched to **Prefix matching** using the `spec.PrefixMatching` field
of the [OpenAPI CRD](openapi-reference.md).

### Authentication

Just one top level security requirement supported.
Operation level security requirements are not supported.

Supported security schemes: `apiKey`.

For the `apiKey` security scheme type:
* *credentials location* will be read from the OpenAPI [in](https://github.com/OAI/OpenAPI-Specification/blob/main/versions/3.0.2.md#security-scheme-object) field of the security scheme object.
* *Auth user key* will be read from the OpenAPI [name](https://github.com/OAI/OpenAPI-Specification/blob/main/versions/3.0.2.md#security-scheme-object) field of the security scheme object.

Partial example of OpenAPI (3.0.2) with `apiKey` security requirement

```yaml
---
openapi: "3.0.2"
security:
  - petstore_api_key: []
components:
  securitySchemes:
    petstore_api_key:
      type: apiKey
      name: api_key
      in: header
```

When OpenAPI does not specify any security requirements:
* The product authentication will be configured for `apiKey`.
* *credentials location* will default to 3scale value `As query parameters (GET) or body parameters (POST/PUT/DELETE)`.
* *Auth user key* will default to 3scale value `user_key`

3scale *Authentication Security* can be set using the `spec.privateAPIHostHeader` and the `spec.privateAPISecretToken` fields of the [OpenAPI CR](openapi-reference.md).

### ActiveDocs

No 3scale ActiveDoc is created.

### 3scale Application Plans

Custom application plans can be added to the product using [OAS 3scale extensions](openapi-3scale-extensions.md#root-level-3scale-extension).

### 3scale Product Policy Chain

3scale policy chain will be the default one created by 3scale.
This can be overridden with a custom policy chain using [OAS 3scale extensions](openapi-3scale-extensions.md#root-level-3scale-extension).

### 3scale Deployment Mode

By default, the configured 3scale deployment mode will be `APIcast 3scale managed`.
However, when the `spec.productionPublicBaseURL` or the `spec.stagingPublicBaseURL` (or both)
fields are provided in the [OpenAPI custom resource](openapi-reference.md),
the product's deployment mode will be `APIcast self-managed`.

Example of a [OpenAPI custom resource](openapi-reference.md) with custom public base URL:

```yaml
apiVersion: capabilities.3scale.net/v1beta1
kind: OpenAPI
metadata:
  name: openapi1
spec:
  openapiRef:
    url: "https://raw.githubusercontent.com/OAI/OpenAPI-Specification/master/examples/v3.0/petstore.yaml"
  productionPublicBaseURL: "https://production.my-gateway.example.com"
  stagingPublicBaseURL: "https://staging.my-gateway.example.com"
```

## Minimum required OAS doc

In [OAS 3.0.2](https://github.com/OAI/OpenAPI-Specification/blob/main/versions/3.0.2.md#oasDocument),
the minimum **valid** OpenAPI document just contains `info` and `paths` fields.

For instance:

```yaml
---
openapi: "3.0.2"
info:
  title: "some title"
  description: "some description"
  version: "1.0.0"
paths:
  /pet:
    get:
      operationId: "getPet"
      responses:
        405:
          description: "invalid input"
```

However, with this OpenAPI document, there is critical 3scale configuration lacking and
it must be provided for a working 3scale configuration:
* `Private Base URL` filling the `spec.privateBaseURL` field of the [OpenAPI CRD](openapi-reference.md)

The minimum valid [OpenAPI custom resource](openapi-reference.md) for a working 3scale product is:

```yaml
apiVersion: capabilities.3scale.net/v1beta1
kind: OpenAPI
metadata:
  name: openapi1
spec:
  openapiRef:
    url: "https://raw.githubusercontent.com/OAI/OpenAPI-Specification/master/examples/v3.0/petstore.yaml"
```

*Note*: The referenced OpenAPI document should include the `servers[0].url` field. For instance:

```yaml
---
openapi: "3.0.2"
info:
  title: "some title"
  description: "some description"
  version: "1.0.0"
servers:
  - url: https://petstore.swagger.io/v1
paths:
  /pet:
    get:
      operationId: "getPet"
      responses:
        405:
          description: "invalid input"
```

*Note*: 3scale still requires creating the application key, but this is out of scope.

## Update behavior
Once an OpenAPI CR has been created and reconciled, a corresponding Product CR and Backend CR are created. Any changes made to those Product and Backend CRs will be reverted by the operator. This is to ensure that the source of truth is the OpenAPI Spec (OAS) and OpenAPI CR. 

In order to make changes to the Product and Backend, you should instead update the OAS source itself as described below.

### Updating secret ref
When the OpenAPI CR is using a Secret for its OAS source, any changes made to the secret will be automatically passed on to the Product and Backend CR and then to the Product and Backend in the 3scale UI.

### Updating URL ref
When the OpenAPI CR is using a URL for its OAS source, the operator will scrape the URL to check for changes every 5 minutes. If changes are detected they will be automatically passed on to the Product and Backend CR and then to the Product and Backend in the 3scale UI.

## Link your OpenAPI spec to your 3scale tenant or provider account

When some [OpenAPI custom resource](openapi-reference.md) is found by the 3scale operator,
the *LookupProviderAccount* process is started to figure out the tenant owning the resource.

The process will check the following tenant credential sources. If none is found, an error is raised.

* Read credentials from *providerAccountRef* resource attribute. This is a secret local reference, for instance `mytenant`

```
apiVersion: capabilities.3scale.net/v1beta1
kind: OpenAPI
metadata:
  name: openapi1
spec:
  openapiRef:
    url: "https://raw.githubusercontent.com/OAI/OpenAPI-Specification/master/examples/v3.0/petstore.yaml"
  providerAccountRef:
    name: mytenant
```

[OpenAPI CRD Reference](openapi-reference.md) for more info about fields.

The `mytenant` secret must have`adminURL` and `token` fields with tenant credentials. For example:

```
apiVersion: v1
kind: Secret
metadata:
  name: mytenant
type: Opaque
stringData:
  adminURL: https://my3scale-admin.example.com:443
  token: "XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX"
```

* Default `threescale-provider-account` secret

For example: `adminURL=https://3scale-admin.example.com` and `token=123456`.

```
oc create secret generic threescale-provider-account --from-literal=adminURL=https://3scale-admin.example.com --from-literal=token=123456
```

* Default provider account in the same namespace 3scale deployment

The operator will gather required credentials automatically for the default 3scale tenant (provider account) if 3scale installation is found in the same namespace as the custom resource.
