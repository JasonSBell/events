package events

import "encoding/json"

type GenericEvent Event[json.RawMessage]
