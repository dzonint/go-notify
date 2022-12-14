# go-notify

Package go-notify enables independent components of an application to observe notable events in a decoupled fashion. It eliminates the need for components to have intimate konwledge of each other (only names of the events are shared).

It contains 2 implementations:
- Go - Uses channels exclusively to send and receive data between producers and listeners.
- Kafka - Uses Kafka writer to send data and consumer to receive it and forward it to appropriate listeners, thus enabling inter-application real time communication.


Both implementations work on a principle of storing channel - event pairings in a map and using it to determine which channels need to receive which messages.

Example:

    // Producer of "my_event".
    go func() {
        for {
            time.Sleep(time.Duration(1) * time.Second):
            notify.Post("my_event", time.Now().Unix())
        }
    }()

    // Consumer of "my_event" (normally some independent component that
    // needs to be notified when "my_event" occurs).
    myEventChan := make(chan interface{})
    notify.Start("my_event", myEventChan)
    go func() {
        for {
            data := <-myEventChan
            log.Printf("MY_EVENT: %#v", data)
        }
    }()

### Functions

    func Post(event string, data interface{}) error
        Post a notification (arbitrary data) to the specified event

    func Start(event string, outputChan chan interface{})
        Start observing the specified event via provided output channel

    func Stop(event string, outputChan chan interface{}) error
        Stop observing the specified event on the provided output channel

    func Close()
        Performs any graceful shutdown required

## License

go-notify is licensed under the New BSD License. More information can be found in LICENSE.txt.