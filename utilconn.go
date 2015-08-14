package canopus

import (
	"log"
	"net"
	"time"
)

// Sends a 402 Error - Bad Option
func SendError402BadOption(messageId uint16, conn *net.UDPConn, addr *net.UDPAddr) {
	msg := NewMessage(TYPE_NONCONFIRMABLE, COAPCODE_501_NOT_IMPLEMENTED, messageId)
	msg.SetStringPayload("Bad Option: An unknown option of type critical was encountered")

	SendMessageTo(msg, NewCanopusUDPConnection(conn), addr)
}

// Sends a CoAP Message to UDP address
func SendMessageTo(msg *Message, conn CanopusConnection, addr *net.UDPAddr) (*Response, error) {
	if conn == nil {
		return nil, ERR_NIL_CONN
	}

	if msg == nil {
		return nil, ERR_NIL_MESSAGE
	}

	if addr == nil {
		return nil, ERR_NIL_ADDR
	}

	b, _ := MessageToBytes(msg)
	_, err := conn.WriteTo(b, addr)

	if err != nil {
		return nil, err
	}

	if msg.MessageType == TYPE_NONCONFIRMABLE {
		return nil, err
	} else {
		// conn.SetReadDeadline(time.Now().Add(time.Second * 2))
		buf, n, err := conn.Read()
		if err != nil {
			return nil, err
		}
		msg, err := BytesToMessage(buf[:n])
		resp := NewResponse(msg, err)

		return resp, err
	}
	return nil, nil
}

// Sends a CoAP Message to a UDP Connection
func SendMessage(msg *Message, conn CanopusConnection) (*Response, error) {
	if conn == nil {
		return nil, ERR_NIL_CONN
	}

	b, _ := MessageToBytes(msg)
	_, err := conn.Write(b)

	if err != nil {
		return nil, err
	}

	if msg.MessageType == TYPE_NONCONFIRMABLE {
		return nil, err
	} else {
		var buf []byte = make([]byte, 1500)
		conn.SetReadDeadline(time.Now().Add(time.Second * 2))
		buf, n, err := conn.Read()

		if err != nil {
			return nil, err
		}

		msg, err := BytesToMessage(buf[:n])

		resp := NewResponse(msg, err)

		return resp, err
	}
}

func SendAsyncMessage(msg *Message, conn *net.UDPConn, fn ResponseHandler) {
	b, _ := MessageToBytes(msg)
	_, err := conn.Write(b)

	if err != nil {
		log.Println(err)

		fn(nil, err)
	}

	if msg.MessageType == TYPE_NONCONFIRMABLE {
		fn(nil, err)
	} else {
		var buf []byte = make([]byte, 1500)
		conn.SetReadDeadline(time.Now().Add(time.Second * 2))
		n, _, err := conn.ReadFromUDP(buf)

		if err != nil {
			fn(nil, err)
		}

		msg, err := BytesToMessage(buf[:n])

		resp := NewResponse(msg, err)

		fn(resp, err)
	}
}
