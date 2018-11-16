# ATEM

## Type `Atem`

> Represents an ATEM device

**Usage:**
```go
// Create an Atem instance without debug output
at := Atem.Create("192.168.0.100", false)
```

## Properties

* **Ip** `string` - The IP address of the ATEM
* **State** `ConnectionState` - The state of the connection to the ATEM
* **Debug** `bool` - Whether to log information about device communication to the console
* **UID** `uint16` - Unique identifier for the connection
* **ProtocolVersion** `types.Version` - Version of Blackmagic Design ATEM protocol used by the device in X.Y format
* **ProductId** `types.NullTerminatedString` - Model of the ATEM device
* **Warn** `types.NullTerminatedString` - Warning communicated by the device
* **Topology** `types.Topology` - Features and capabilites of the device
* **MixEffectConfig** `types.MixEffectConfig` - Capabilities of the device's Mix Effect Block(s)
* **MediaPlayers** `types.MediaPlayers` - Number of stills and video clips the devices Media Players can hold
* **MultiViewCount** `uint8` - Number of multi-views supported by the device
* **AudioMixerConfig** `types.AudioMixerConfig` - Capabilities of the audio mixer
* **VideoMixerConfig** `VideoMixerConfig` - List of video modes supported by the device
* **MacroPool** `uint8` - Number of macros supported by the device
* **PowerStatus** `types.PowerStatus` - Power status of the device 
* **VideoMode** `*types.VideoMode` - Properties of the video standard used by the device
* **VideoSources** `*types.VideoSources` - List of the device's video sources
* **ProgramInput** `*types.ProgramInput` - Input currently on the program bus
* **PreviewInput** `*types.PreviewInput` - Input currently on the preview bus

## Static Methods

### **`Atem.Create(ipAddress, debug) *Atem`**

* **ipAddress** `string` - The IP address of the device with which to establish a connection
* **debug** `bool` - Whether to log additional output to the console regarding communication with the device

Retuns an `Atem` instance (not yet connected)

## Instance Methods

### **`atem.Connect()` error**

Attempts to the device at the specified IP address (`atem.Ip`)

Returns an error if there was a problem connecting to the device

### **`atem.Connected() bool`**

Returns a `bool` indicating whether the device is connected (`true`) or not (`false`)

### **`atem.Close()`**

Closes the connection with the device

### **`atem.On(eventName, callback)`**

* **eventName** `string` - Name of the event that triggers the callback function
* **callback** `func()` - Callback to execute when the event is emitted

Triggers a callback to be executed when an event by the name specified in `eventName` is emitted (see [events](#Instance_Events))

## Instance Events

Instances of the `Atem` type emit the following events:

**Event: "connected"**

Emitted when a connection with the ATEM device is established

**Event: "closed"**

Emitted when a connection with the ATEM device is closed