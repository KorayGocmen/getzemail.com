import React, {useEffect, useState} from "react";
import { withRouter, Link } from "react-router-dom";

import api from "../../api";

function Inbox(props) {
  const [inbox, setInbox] = useState(null);
  const [messages, setMessages] = useState([]);

  useEffect(() => {
    async function fetchInbox() {
      try {
        const response = await api.fetchInbox(props.match.params.address);
        setInbox(response.data.mail_inbox);
        setMessages(response.data.mail_inbox.mail_messages);
      } catch (e) {
        console.log(e);
      }
    }

    fetchInbox();
    const timer = window.setInterval(() => fetchInbox(), 5000);

    return () => {
      window.clearInterval(timer);
    }
  }, []);

  return (
    <div className="Inbox">
      <h2> { inbox ? inbox.address : "Inbox not found"} </h2>
      <ol className="Messages">
        {messages ? messages.map((message, index) => {
          return <li key={message.id}>
            <Link to={{pathname: `/messages/${message.id}`}}>
              <div>
              <b> { message.subject } </b>
              <p> { message.text || message.html } </p>
              </div>
            </Link>
          </li>
        }) : "No Messages" }
      </ol>
    </div>
  )
}

export default withRouter(Inbox);
