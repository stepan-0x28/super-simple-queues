# Super simple queues

Super simple queuing system.

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

## Interacting with the system via TCP

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

### Message types and their structure

Three types of messages are used to interact with the system:

- Init
- Payload
- Confirm

Each message begins with a one-byte header that defines the message type.

#### Message type "Init"

<table>
<tr>
<td align="center">Type</td>
<td align="center">Operating mode</td>
<td align="center">Queue key length</td>
<td align="center">Queue key</td>
</tr>
<tr>
<td align="center">1 byte<br>(<code>uint8</code>)</td>
<td align="center">1 byte<br>(<code>uint8</code>)</td>
<td align="center">1 byte<br>(<code>uint8</code>)</td>
<td align="center">N bytes<br>(<code>utf8</code>)</td>
</tr>
</table>

Init message type is always `0x01`. The client's operating mode can be either `0x00` or `0x01`, where `0x00` is the
receiving mode and `0x01` is the sending mode.

Example:

<table>
<tr>
<td align="center">Type</td>
<td align="center">Operating mode</td>
<td align="center">Queue key length</td>
<td align="center">Queue key</td>
</tr>
<tr>
<td align="center"><code>0x01</code></td>
<td align="center"><code>0x01</code></td>
<td align="center"><code>0x04</code></td>
<td align="center"><code>test</code></td>
</tr>
</table>

#### Message type "Payload"

<table>
<tr>
<td align="center">Type</td>
<td align="center">Data length</td>
<td align="center">Data</td>
</tr>
<tr>
<td align="center">1 byte<br>(<code>uint8</code>)</td>
<td align="center">4 bytes<br>(<code>uint32</code>)</td>
<td align="center">N bytes<br>(<code>utf8</code>)</td>
</tr>
</table>

Payload message type is always `0x02`.

Example:

<table>
<tr>
<td align="center">Type</td>
<td align="center">Data length</td>
<td align="center">Data</td>
</tr>
<tr>
<td align="center"><code>0x02</code></td>
<td align="center"><code>0x00000009</code></td>
<td align="center"><code>some data</code></td>
</tr>
</table>

#### Message type "Confirm"

<table>
<tr>
<td align="center">Type</td>
</tr>
<tr>
<td align="center">1 byte<br>(<code>uint8</code>)</td>
</tr>
</table>

Confirm message type is always `0x03`.

Example:

<table>
<tr>
<td align="center">Type</td>
</tr>
<tr>
<td align="center"><code>0x03</code></td>
</tr>
</table>