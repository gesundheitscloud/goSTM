## Performance
- Each request opens a new connection -> ideally, connections could be persisted for multiple requests (no connection pooling)

## UX
- UI must look better. Add green/red icon to show whether tunnel is active or not

## Usability
- allow github tokens for private repos

## Security
- Verify host key from `known_hosts`

## Stability
- Travis CI: do not only build but actually try to create a tunnel

## Possible Features
- Integration with vault: get ssh private key from vault
