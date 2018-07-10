## Performance
- Each request opens a new connection -> ideally, connections could be persisted for multiple requests (no connection pooling)

## UX
- UI must look better. Add green/red icon to show whether tunnel is active or not

## Usability
- `ssh_config` read from custom location

## Security
- Verify host key from `known_hosts`

## Stability
- Travis CI