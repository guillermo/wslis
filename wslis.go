/*
Package wslis allows a server to create connections to clients.

wslis allows a server to create connections to clients. For that, clients
previously stablish a connection to the server through websockets. Once the
connection arrives, the server and the clients, multiplex the stream
so several connections can go through the same http stream.

History and Design Considerations

wslis comes from the needs of delegatescreen.com that mainly display webpages
in office dashboards. The most common operations are: Set the web page to
display, and get a screenshot of the current page. This operations have to be
syncronous, and completes as fast as possible.

This lead us to the following desing considerations:

- Request-Response paradigm.

- It works over http, so edges, load balancers and other network
components work out of the box.


And this are the options considered.

- ZeroMq, Mqtt or most iot platforms are pub sub based. So even if most of them can work over http, we will need to abstract some kind of request/response ciclye.

- Signaling where the server sends a signal (kind of request) to the client,
and the client initiates an http request. For signaling we considered: longpush,
server-sent-events, http2 push . This option was discarded
to remove the need of designing a signaling mecanish, and the added
complication of matching signals to client/requests.

- connect the clients to vpns through wireguard. As clients can run on
windows/mac/linux this will imply the automation of all the configuraiton in
three ooss, and maintain that over time. Also a missconfiguration may brick
the computer.

- ssh over http. This will give the server control over the clients to run
any commands. Even if this might be flexible enough, a compromise server could
infect a client and access client network.


Usage

See Example for usage info.

*/
package wslis
