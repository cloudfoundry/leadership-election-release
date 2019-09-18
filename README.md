**Archived** Now a part of https://github.com/cloudfoundry/system-metrics-release

Leadership Election Release
=================

The leadership election release deploys job that implements the Raft distributed consensus algorithm to the desired
 instance group.
 
### Usage

Once deployed, the `leadership-election` job will expose an endpoint on localhost with the following responses:
- 200 OK: This node is the leader
- 423 LOCKED: This node is not the leader

By default, the endpoint is exposed on port 8080

Example:

`curl localhost:8080/v1/leader`

