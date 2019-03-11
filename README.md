Bottr go
========

Twitter Bot with sentiments

Features
====
* Stream tweets filtering for topics or users
* Compute Sentiment via MachineBox (Needs MachineBox Url)

In Action
=========
![Terminal](https://github.com/gabber12/bottrgo/raw/master/img/inaction.png)

How to Use
========

### Build App

```
make build
```

### Stream Filter 
```
 ./tweety -textFilter "#cats"       
```

### Stream Filter and classify with MachineBox
```
 ./tweety -textFilter "#cats" -classify -mbHost "http://localhost:80
90"         
```
