# Wait Notify Engine


## 1. Motivation
It is a very common use case in modern systems that you want to do some long-running task and on its completion depending upon the response want to perform some action
These tasks can happen on the same service or a completely different one. 

However to reliably trigger these callbacks we need a mechanism that can be fault tolerant across machines. This library is an attempt to tackle such requirements 

## 2. Functional Specification

The wait notify engine should allow users to perform following operations

1. User can queue a `callback` against a `topic` for a given set of `notifyIds`
2. User can notify the `response` for the `notifyId` using the API
3. When the responses for all the notifyIds are received the `notify` method of `callback` is invoked with appropriate `responses`
4. User can also specify a timeout, if the response of all the notifyIds are not received within the specified duration the `timeout` method of callback is invoked with the received responses

## 3. Tech Choices

