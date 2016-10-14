Cache Server
=====================

![](http://dl2.joxi.net/drive/2016/10/13/0011/0529/758289/89/f84998ed49.jpg)

Theory
=====================

Hash tables
--------

For imetentation Hash tables i used chaining scheme, in this scheme, we allocated slice size of "m"(cardinality of a hash function), each element of
this slice is chain that has link to the next element.

Evaluation:
- Memory consumption in O(n+m), where n is number of objects currently stored in the hash table and m is the cardinality of the hash function
- Operations work in time O(c+1), where c is the lenght of the longest chain

Main problem: how to make both m and c small?







  


