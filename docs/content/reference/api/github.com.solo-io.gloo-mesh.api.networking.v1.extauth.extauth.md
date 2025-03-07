
---

title: "extauth.proto"

---

## Package : `extauth.networking.mesh.gloo.solo.io`



<a name="top"></a>

<a name="API Reference for extauth.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## extauth.proto


## Table of Contents
  - [BufferSettings](#extauth.networking.mesh.gloo.solo.io.BufferSettings)
  - [CustomAuth](#extauth.networking.mesh.gloo.solo.io.CustomAuth)
  - [CustomAuth.ContextExtensionsEntry](#extauth.networking.mesh.gloo.solo.io.CustomAuth.ContextExtensionsEntry)
  - [GatewayExtauth](#extauth.networking.mesh.gloo.solo.io.GatewayExtauth)
  - [HttpService](#extauth.networking.mesh.gloo.solo.io.HttpService)
  - [HttpService.Request](#extauth.networking.mesh.gloo.solo.io.HttpService.Request)
  - [HttpService.Request.HeadersToAddEntry](#extauth.networking.mesh.gloo.solo.io.HttpService.Request.HeadersToAddEntry)
  - [HttpService.Response](#extauth.networking.mesh.gloo.solo.io.HttpService.Response)
  - [RouteExtauth](#extauth.networking.mesh.gloo.solo.io.RouteExtauth)

  - [GatewayExtauth.ApiVersion](#extauth.networking.mesh.gloo.solo.io.GatewayExtauth.ApiVersion)






<a name="extauth.networking.mesh.gloo.solo.io.BufferSettings"></a>

### BufferSettings
Configuration for buffering the request data.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| maxRequestBytes | uint32 |  | Sets the maximum size of a message body that the filter will hold in memory. Envoy will return *HTTP 413* and will *not* initiate the authorization process when buffer reaches the number set in this field. Note that this setting will have precedence over failure_mode_allow. Defaults to 4KB. |
  | allowPartialMessage | bool |  | When this field is true, Envoy will buffer the message until *max_request_bytes* is reached. The authorization request will be dispatched and no 413 HTTP error will be returned by the filter. |
  | packAsBytes | bool |  | When this field is true, Envoy will send the body sent to the external authorization service with raw bytes. |
  





<a name="extauth.networking.mesh.gloo.solo.io.CustomAuth"></a>

### CustomAuth
Gloo Mesh is not expected to configure the ext auth server in this case. This is used with custom auth servers.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| contextExtensions | [][extauth.networking.mesh.gloo.solo.io.CustomAuth.ContextExtensionsEntry]({{< versioned_link_path fromRoot="/reference/api/github.com.solo-io.gloo-mesh.api.networking.v1.extauth.extauth#extauth.networking.mesh.gloo.solo.io.CustomAuth.ContextExtensionsEntry" >}}) | repeated | When a request matches the route or traffic policy on which this configuration is defined, Gloo Mesh will add the given context_extensions to the request that is sent to the external authorization server. This allows the server to base the auth decision on metadata that you define on the source of the request.<br>This attribute is analogous to Envoy's config.filter.http.ext_authz.v2.CheckSettings. See the official [Envoy documentation](https://www.envoyproxy.io/docs/envoy/latest/api-v2/config/filter/http/ext_authz/v2/ext_authz.proto.html?highlight=ext_authz#config-filter-http-ext-authz-v2-checksettings) for more details. |
  





<a name="extauth.networking.mesh.gloo.solo.io.CustomAuth.ContextExtensionsEntry"></a>

### CustomAuth.ContextExtensionsEntry



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| key | string |  |  |
  | value | string |  |  |
  





<a name="extauth.networking.mesh.gloo.solo.io.GatewayExtauth"></a>

### GatewayExtauth
Enterprise-Only: Configure the Extauth Filter on a Gateway


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| extauthzRef | [core.skv2.solo.io.ObjectRef]({{< versioned_link_path fromRoot="/reference/api/github.com.solo-io.skv2.api.core.v1.core#core.skv2.solo.io.ObjectRef" >}}) |  | The destination ref for the envoy external authentication service to ask about auth decisions.<br>The fqdn used by envoy for extauth requests is derived from the kube service or is the first hostname in an external service destination. IMPORTANT: Envoy's extauth requests go to the port chosen from the destination that has the name or protocol of the extauth service protocol (i.e., grpc or http). In practice, this means the port handling extauth requests should be named grpc and/or have protocol GRPC (both case-insensitive), unless `http_service` is configured.<br>The provided ref will be used to search for a service of the given name/namespace on each cluster that a Gateway will be created.<br>If omitted, Gloo Mesh will search for a service with the name ext-auth-service in each namespace on each cluster that a Gateway will be created. |
  | httpService | [extauth.networking.mesh.gloo.solo.io.HttpService]({{< versioned_link_path fromRoot="/reference/api/github.com.solo-io.gloo-mesh.api.networking.v1.extauth.extauth#extauth.networking.mesh.gloo.solo.io.HttpService" >}}) |  | If this is set, communication to the upstream will be via HTTP and not GRPC. |
  | requestTimeout | [google.protobuf.Duration]({{< versioned_link_path fromRoot="/reference/api/github.com.solo-io.protoc-gen-ext.external.google.protobuf.duration#google.protobuf.Duration" >}}) |  | Timeout for the ext auth service to respond. Defaults to 2000ms. (OIDC requests to external IDPs can be slow) For latency critical applications, this value should be tuned much lower, as extauth is on the request path. |
  | failureModeAllow | bool |  | In case of a failure or timeout querying the auth server, normally a request is denied. if this is set to true, the request will be allowed. |
  | requestBody | [extauth.networking.mesh.gloo.solo.io.BufferSettings]({{< versioned_link_path fromRoot="/reference/api/github.com.solo-io.gloo-mesh.api.networking.v1.extauth.extauth#extauth.networking.mesh.gloo.solo.io.BufferSettings" >}}) |  | Set this if you also want to send the body of the request, and not just the headers. |
  | clearRouteCache | bool |  | Clears route cache in order to allow the external authorization service to correctly affect routing decisions. Filter clears all cached routes when:<br>1. The field is set to *true*.<br>2. The status returned from the authorization service is a HTTP 200 or gRPC 0.<br>3. At least one *authorization response header* is added to the client request, or is used for altering another client request header. |
  | statusOnError | uint32 |  | Sets the HTTP status that is returned to the client when there is a network error between the filter and the authorization server. The default status is HTTP 403 Forbidden. If set, this must be one of the following: - 100 - 200 201 202 203 204 205 206 207 208 226 - 300 301 302 303 304 305 307 308 - 400 401 402 403 404 405 406 407 408 409 410 411 412 413 414 415 416 417 421 422 423 424 426 428 429 431 - 500 501 502 503 504 505 506 507 508 510 511 |
  | transportApiVersion | [extauth.networking.mesh.gloo.solo.io.GatewayExtauth.ApiVersion]({{< versioned_link_path fromRoot="/reference/api/github.com.solo-io.gloo-mesh.api.networking.v1.extauth.extauth#extauth.networking.mesh.gloo.solo.io.GatewayExtauth.ApiVersion" >}}) |  | Determines the API version for the `ext_authz` transport protocol that will be used by Envoy to communicate with the auth server. Defaults to `V3`. For more info, see the `transport_api_version` field [here](https://www.envoyproxy.io/docs/envoy/latest/api-v3/extensions/filters/http/ext_authz/v3/ext_authz.proto#extensions-filters-http-ext-authz-v3-extauthz). |
  | statPrefix | string |  | Optional additional prefix to use when emitting statistics. This allows to distinguish emitted statistics between configured ext_authz filters in an HTTP filter chain. |
  





<a name="extauth.networking.mesh.gloo.solo.io.HttpService"></a>

### HttpService



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| pathPrefix | string |  | Sets a prefix to the value of authorization request header *Path*. |
  | request | [extauth.networking.mesh.gloo.solo.io.HttpService.Request]({{< versioned_link_path fromRoot="/reference/api/github.com.solo-io.gloo-mesh.api.networking.v1.extauth.extauth#extauth.networking.mesh.gloo.solo.io.HttpService.Request" >}}) |  |  |
  | response | [extauth.networking.mesh.gloo.solo.io.HttpService.Response]({{< versioned_link_path fromRoot="/reference/api/github.com.solo-io.gloo-mesh.api.networking.v1.extauth.extauth#extauth.networking.mesh.gloo.solo.io.HttpService.Response" >}}) |  |  |
  





<a name="extauth.networking.mesh.gloo.solo.io.HttpService.Request"></a>

### HttpService.Request



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| allowedHeaders | []string | repeated | These headers will be copied from the incoming request to the request going to the auth server. Note that in addition to the user's supplied matchers:<br>1. *Host*, *Method*, *Path* and *Content-Length* are automatically included to the list.<br>2. *Content-Length* will be set to 0 and the request to the authorization service will not have a message body. |
  | headersToAdd | [][extauth.networking.mesh.gloo.solo.io.HttpService.Request.HeadersToAddEntry]({{< versioned_link_path fromRoot="/reference/api/github.com.solo-io.gloo-mesh.api.networking.v1.extauth.extauth#extauth.networking.mesh.gloo.solo.io.HttpService.Request.HeadersToAddEntry" >}}) | repeated | These headers that will be included to the request to authorization service. Note that client request of the same key will be overridden. |
  





<a name="extauth.networking.mesh.gloo.solo.io.HttpService.Request.HeadersToAddEntry"></a>

### HttpService.Request.HeadersToAddEntry



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| key | string |  |  |
  | value | string |  |  |
  





<a name="extauth.networking.mesh.gloo.solo.io.HttpService.Response"></a>

### HttpService.Response



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| allowedUpstreamHeaders | []string | repeated | When this is set, authorization response headers that have a will be added to the original client request and sent to the upstream. Note that coexistent headers will be overridden. |
  | allowedClientHeaders | []string | repeated | When this. is set, authorization response headers that will be added to the client's response when auth request is denied. Note that when this list is *not* set, all the authorization response headers, except *Authority (Host)* will be in the response to the client. When a header is included in this list, *Path*, *Status*, *Content-Length*, *WWW-Authenticate* and *Location* are automatically added. |
  





<a name="extauth.networking.mesh.gloo.solo.io.RouteExtauth"></a>

### RouteExtauth
Extauth configuration for a Route or TrafficPolicy. Configures extauth for individual HTTP routes


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| disable | bool |  | Set to true to disable auth on the route. |
  | configRef | [core.skv2.solo.io.ObjectRef]({{< versioned_link_path fromRoot="/reference/api/github.com.solo-io.skv2.api.core.v1.core#core.skv2.solo.io.ObjectRef" >}}) |  | A reference to an AuthConfig. This is used to configure the mesh clients to identify themselves by matching their client identifier to the extauth server config for the same AuthConfig. |
  | customAuth | [extauth.networking.mesh.gloo.solo.io.CustomAuth]({{< versioned_link_path fromRoot="/reference/api/github.com.solo-io.gloo-mesh.api.networking.v1.extauth.extauth#extauth.networking.mesh.gloo.solo.io.CustomAuth" >}}) |  | Use this field if you are running your own custom extauth server. |
  




 <!-- end messages -->


<a name="extauth.networking.mesh.gloo.solo.io.GatewayExtauth.ApiVersion"></a>

### GatewayExtauth.ApiVersion
Describes the transport protocol version to use when envoy connects to the ext auth server.

| Name | Number | Description |
| ---- | ------ | ----------- |
| GLOO_MESH_AUTO | 0 | Use transport version that matches the version the default enterprise ext-auth-service uses. |
| ENVOY_AUTO | 1 | Use envoy's auto transport version. This will change as envoy cycles through transport api versions. |
| V3 | 2 | Use v3 API. |


 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->

