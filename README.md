# tcpBashServer

tcp Server that processes bash from a tcp client.
The tcp server manages commands to bash via pty terminal.


## tcpBashServer

## tcpBashServerV2

## tcpBashServerV3
Found that bash returns several char sequences.
In the test case it was 3.

## tcpBashServerV4

## tcpBashServerV5
Fixed io.copy. problem is io.copy return error code ELO

## tcpBashServerV6
goal is returning the bash output to the client

## tcpBashServerV7
goal is creating a loop


## todo
 - termios to manage echo and CRNL
 - loop for bash cli
 - stopping the loop with cntl-c
 - terminal size
 - making the server a service
 - logging


