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

#### Setup Environment Variables

Get your credentials to use twitter api at http://apps.twitter.com
```
export TWITTER_CONSUMER_KEY=
export TWITTER_CONSUMER_SECRET=
export TWITTER_ACCESS_TOKEN=
export TWITTER_ACCESS_TOKEN_SECRET=
```

## Usage

### Stream Filter 
```
 ./tweety -textFilter "#cats"       
```

### Stream Filter and classify with MachineBox
```
 ./tweety -textFilter "#cats" -classify -mbHost "http://localhost:80
90"         
```

### Complete list of Flags

```
 ./tweety -textFilter "#cats|#dogs" -locationFilter "-74,40,-73,41" -userFilter "userId" -classify -mbHost "http://localhost:80
90"  
```

## Notes
Checkout https://machinebox.io.