# Wait Notify Engine

## 1. Functional Specification

The wait notify engine should allow users to perform following operations

1. User can queue a `callback` against a `topic` for a given set of `notifyIds`
2. User can notify the `response` for the `notifyId` using the API
3. When the responses for all the notifyIds are received the `notify` method of `callback` is invoked with appropriate `responses`
4. User can also specify a timeout, if the response of all the notifyIds are not received within the specified duration the `timeout` method of callback is invoked with the received responses



