//takes in two executable paths
//starts up two of these dockers, with port 10000 in the VM bound to two random numbers on the host
//copies in the two executables to /code in both images
//makes a connection to localhost:whatever, which is docker:10000 on both. this makes relay.go start the provided /code, and forwards stdin/stdout over this socket
//runs iterated prisoner's dilemma, and finally returns number of rounds and final socre
