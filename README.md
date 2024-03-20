# healthyy

A healthcheck tool written in golang with 0 dependencies.

# Configuration files

Specifiy your checks using the following configuration format:

https://github.com : 15s 

Or directly from cli:

./healthyy https://github.com 15s https://google.com 3s
