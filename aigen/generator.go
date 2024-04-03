package aigen

/**
Generator is a function specific to a AI Model for interacting with a text generator
**/

import "ai-manager/aimsg"

type Generator func(model, apikey string, conversation aimsg.Conversation, meta ...aimsg.Meta) (aimsg.Message, error)
