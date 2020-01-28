package actions

func MessageHelloReply() []byte {
	var t = "<?xml version='1.0'?>" +
		"<stream:stream xmlns='jabber:client' xmlns:stream='http://etherx.jabber.org/streams' id='1' from='jabber.ru' version='1.0' xml:lang='en'>" +
		"<stream:features>" +
		"<starttls xmlns='urn:ietf:params:xml:ns:xmpp-tls'/>" +
		//"<compression xmlns='http://jabber.org/features/compress'><method>zlib</method></compression>" +
		"<mechanisms xmlns='urn:ietf:params:xml:ns:xmpp-sasl'><mechanism>PLAIN</mechanism></mechanisms>" +
		"</stream:features>"
	return []byte(t)
}

func MessageAfterLogged() string {
	var t = "<?xml version='1.0'?>" +
		"<stream:stream xmlns='jabber:client' xmlns:stream='http://etherx.jabber.org/streams' id='1' from='jabber.ru' version='1.0' xml:lang='en'>" +
		"<stream:features>" +
		"<bind xmlns='urn:ietf:params:xml:ns:xmpp-bind'/><session xmlns='urn:ietf:params:xml:ns:xmpp-session'/>" +
		"</stream:features>"
	return t
}
