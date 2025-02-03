package server

import (
	"github.com/rebay1982/wsjtx-udp"
	"net"
)

type WSJTXServer struct {
	active bool
	conn   *net.UDPConn
	parser *wsjtxudp.WSJTXParser
}

// NewWSJTXServer creates a new WSJTXServer (UDP) server and listens on the speciified ip/port combination.
func NewWSJTXServer(ip string, port int) (*WSJTXServer, error) {
	server := &WSJTXServer{
		parser: &wsjtxudp.WSJTXParser{},
	}

	err := server.Listen(ip, port)

	return server, err
}

// Close closes the connection.
func (s *WSJTXServer) Close() error {
	err := s.conn.Close()
	if err == nil {
		s.active = false
	}

	return err
}

// Listen starts listening on the specified ip and port.
func (s *WSJTXServer) Listen(ip string, port int) error {
	if s.active {
		err := s.Close()
		if err != nil {
			return err
		}
	}

	addr := net.UDPAddr{
		Port: port,
		IP:   net.ParseIP(ip),
	}

	conn, err := net.ListenUDP("udp", &addr)
	if err == nil {
		s.conn = conn
		s.active = true
	}

	return err
}

// ReadFromUDP reads a WSJTX messages from the server.
func (s WSJTXServer) ReadFromUDP() (wsjtxudp.WSJTXMessage, error) {
	buff := make([]byte, 1024)

	_, _, err := s.conn.ReadFromUDP(buff)
	if err != nil {
		return wsjtxudp.WSJTXMessage{}, err
	}

	message, err := s.parser.Parse(buff)
	return message, err
}
