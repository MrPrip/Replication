# Replication

Initial thoughts:

- Run a server with a port given in command line arguments. Check if the possible to dial port, if not run a server with port+1 to as backup server.
- Run clients with another port. Dial the main server port, which is run with the same port every time. If not possible to connect dial serverport+1.
