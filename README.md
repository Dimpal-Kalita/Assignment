## Problem 1:
Build a webserver in Go to receive data external and route through messaging systems.

### Solution:
The solution is a simple web server that listens on port $8080$, with redis running on port $6379$ and kafka on $29092$. It accepts **POST** requests on the `/push` endpoint and **GET** requests on the  `/get` endpoint.

Request Format for push: 

```sh
curl --location 'http://localhost:8080/push' \
--header 'Content-Type: application/json' \
--data '{
    "key": 26,
    "value": 20
}'
```
Request Format for get: 

```sh
curl --location 'http://localhost:8080/get?key=26'
```

The server will store the data in redis and send the data to kafka. The server will also fetch the data from redis and return it to the client.


Redis is used to store the latest query and also all the queries that are made to the server. 



## Problem 2:
Write code in C to build a tree of arbitrary depth where the number of branches per level
equals the integer sequence in the expansion of the irrational number Pi.

### Solution:
The depth depends on the digits provided in the value of the pi. (default: $pi= 3.141592653589793238462643$)
The levels is the depth of the tree and the number of branches as the value of the pi.

The Tree is designed like this:

```
             1     
          / /\ \
         2 3  4 5
        /
        6
    / / | \ \
   7 8  9 10 11
   .
   .
   .
```

After building the tree the program will print the tree in the following format:

```
BFS Traversal : 1 2 3 4 5 6 7 8 9 10 11 12 13 14 15 16 17 18 19 20 21 22 23 24 25 26 27 28 29 30 31 32 33 34 35 36 37 38 39 40 41 42 43 44 45 46 47 48 49 50 51 52 53 54 55 56 57 58 59 60 61 62 63 64 65 66 67 68 69 70 71 72 73 74 75 76 77 78 79 80 81 82 83 84 85 86 87 88 89 90 91 92 93 94 95 96 97 98 99 100 101 102 103 104 105 106 107 108 109 110 111 112 113 114 115 

DFS Traversal : 1 2 6 7 12 21 23 29 34 37 42 50 59 66 75 78 80 83 91 95 101 103 109 113 114 115 110 111 112 104 105 106 107 108 102 96 97 98 99 100 92 93 94 84 85 86 87 88 89 90 81 82 79 76 77 67 68 69 70 71 72 73 74 60 61 62 63 64 65 51 52 53 54 55 56 57 58 43 44 45 46 47 48 49 38 39 40 41 35 36 30 31 32 33 24 25 26 27 28 22 13 14 15 16 17 18 19 20 8 9 10 11 3 4 5 
```
