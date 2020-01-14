package actions

import (
	"fmt"
	"log"
)

func (a *ActionTemplate) ActionStorage() bool {

	if a.data.Query.Storage.XMLName.Space == "storage:metacontacts" {
		return a.ActionGetMetaContacts()
	}
	if a.data.Query.Storage.XMLName.Space == "storage:bookmarks" {
		return a.ActionGetBookmarks()
	}
	if a.data.Query.Storage.XMLName.Space == "storage:rosternotes" {
		return a.ActionGetRosterNotes()
	}
	if a.data.Query.Storage.XMLName.Space == "roster:delimiter" {
		return a.ActionGetRosterDelimiter()
	}
	return true
}

func (a *ActionTemplate) ActionGetMetaContacts() bool {
	//todo implement contactlist
	log.Println("DO ContactList")
	DoRespond(a.conn,
		fmt.Sprintf("<iq type=\"result\" xmlns=\"jabber:client\" id=\"%s\" to=\"%s\">"+
			"<query xmlns=\"jabber:iq:private\">"+
			"<storage xmlns=\"storage:metacontacts\" >"+
			"<meta jid='11@zzz.club' tag='1' order='1'/>"+
			"<meta jid='12@zzz.club' tag='2' order='1'/>"+
			"<meta jid='13@zzz.club' tag='3' order='2'/>"+
			"<meta jid='14@zzz.club' tag='4' order='2'/>"+
			"</storage>"+
			"</query>"+
			"</iq>",
			a.data.Id,
			a.user.FullAddr,
		), a.data.Id)
	return true
}

func (a *ActionTemplate) ActionGetRosterDelimiter() bool {
	log.Println("DO Roster Delimiter")
	DoRespond(a.conn,
		fmt.Sprintf("<iq type=\"result\" xmlns=\"jabber:client\" id=\"%s\" to=\"%s\">"+
			"<query xmlns=\"jabber:iq:private\">"+
			"<storage xmlns=\"roster:delimiter\" />"+
			"</query>"+
			"</iq>",
			a.data.Id,
			a.user.FullAddr,
		), a.data.Id)
	return true
}

func (a *ActionTemplate) ActionGetRosterNotes() bool {
	//todo implement roster notes XEP-0145
	log.Println("DO Roster Notes")
	DoRespond(a.conn,

		fmt.Sprintf("<iq type='result' id='%s'>"+
			"<query xmlns='jabber:iq:private'>"+
			"<storage xmlns='storage:rosternotes'>"+
			"<note jid='hamlet@shakespeare.lit' "+
			" cdate='2004-09-24T15:23:21Z' "+
			" mdate='2004-09-24T15:23:21Z'>Seems to be a good writer</note>"+
			"<note jid='juliet@capulet.com'"+
			" cdate='2004-09-27T17:23:14Z'"+
			" mdate='2004-09-28T12:43:12Z'>Oh my sweetest love ...</note>"+
			"</storage>"+
			"</query>"+
			"</iq>",
			a.data.Id,
		), a.data.Id)
	return true
}
