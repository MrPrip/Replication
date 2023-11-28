# Replication

We suggest you run the program with four terminals: two clients and two servers (one being the backup)
Remember to cd into the folders. Run cd client for the client terminal and cd server for the server terminals

1. To start the main server run: go run .
2. The backup server must be started before five seconds after the main server is started. 
3. To start the backup sersver run: go run .
4. Start both clients and start bidding. To start a client run: go run . -name="Your name". Where "Your name" is replaced with some name. We assume names a re unique. 

After 30 seconds the main server dies and the backup server takes over if it was started in time.
Consult the programLog.txt file for all the print statments given in the terminals.

