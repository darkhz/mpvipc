// This code is based on the hub.go code of the Gorilla WebSocket:
//
// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license:
// Copyright (c) 2013 The Gorilla WebSocket Authors. All rights reserved.
// 
// Redistribution and use in source and binary forms, with or without
// modification, are permitted provided that the following conditions are met:
// 
//   Redistributions of source code must retain the above copyright notice, this
//   list of conditions and the following disclaimer.
// 
//   Redistributions in binary form must reproduce the above copyright notice,
//   this list of conditions and the following disclaimer in the documentation
//   and/or other materials provided with the distribution.
// 
// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS" AND
// ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED
// WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
// DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE
// FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL
// DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR
// SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER
// CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY,
// OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
// OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.

package mpvipc

// Hub maintains the set of active EventListeners and broadcasts 
// messages to them.
type Hub struct {
	// Registered EventListeners.
	listeners map[*EventListener]bool

	// Inbound Event messages
	broadcast chan *Event

	// Register requests from the EventListeners.
	register chan *EventListener

	// Unregister requests from EventListeners.
	unregister chan *EventListener
}

func newHub() *Hub {
	return &Hub{
		broadcast:  make(chan *Event),
		register:   make(chan *EventListener),
		unregister: make(chan *EventListener),
		listeners:  make(map[*EventListener]bool),
	}
}

func (h *Hub) run() {
	for {
		select {
		case listener := <-h.register:
			h.listeners[listener] = true
		case listener := <-h.unregister:
			if _, ok := h.listeners[listener]; ok {
				delete(h.listeners, listener)
				close(listener.send)
			}
		case message := <-h.broadcast:
			for listener := range h.listeners {
				listener.send <- message
			}
		}
	}
}
