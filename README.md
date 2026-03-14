# Super simple queuing system

## An example of the system operation in the diagram

```mermaid
flowchart LR
    q1[(Queue)]
    q2[(Queue)]
    s1[Sender]
    s2[Sender]
    s3[Sender]
    s4[Sender]
    r1[Receiver]
    r2[Receiver]
    r3[Receiver]
    r4[Receiver]

    subgraph Super simple queues
        q1
        q2
    end

    s1 & s2 --> q1 --> r1 & r2
    s3 & s4 --> q2 --> r3 & r4
```

## Interaction with the system

### Sequence diagram

```mermaid
sequenceDiagram
    participant c as Client
    participant s as Server
    c ->>+ s: Init message
    s -->>- c: Confirm message
    loop Further communication
        alt Client in sending mode
            c ->>+ s: Payload message
            s -->>- c: Confirm message
        else Client in receiving mode
            s ->>+ c: Payload message
            c -->>- s: Confirm message
        end
    end
```