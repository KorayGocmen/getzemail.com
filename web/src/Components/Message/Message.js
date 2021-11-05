import axios from "axios";
import React, {useEffect, useState} from "react";
import { withRouter } from "react-router-dom";

import api from "../../api";

function Message(props) {
  const [message, setMessage] = useState(null);

  useEffect(()=>{
    (async () => {
      try {
        const response = await api.fetchMessage(props.match.params.messageId);
        const mailMessage = response.data.mail_message;
        
        const textdata = await axios.get(mailMessage.text_url);
        mailMessage.text = textdata.data;

        const htmldata = await axios.get(mailMessage.html_url);
        mailMessage.html = htmldata.data;

        setMessage(mailMessage);
      } catch (e) {
        console.log(e);
      }
    })();
  }, []);

  function Relations(props) {
    if (!message.mail_message_relations) return null

    const relations = message.mail_message_relations
      .filter(m => m.type === props.type)
      .map(r => 
        <span key={r.id}> {r.display_name} &#60;{r.address}&#x3e; </span>
      )

    if (relations.length === 0) {
      return null
    }

    return (
      <div> 
        {props.type}: {relations}
      </div>
    )
  }

  function Files(props) {
    if (!message.mail_message_files) return null

    const files = message.mail_message_files
      .filter(f => f.type === props.type)
      .map(f => 
        <a key={f.id} href={f.url}> {f.file_name} </a>
      )

    if (files.length === 0) {
      return null
    }

    return (
      <div> 
        <b> Attachments </b>
        {files}
      </div>
    )
  }

  return (
    message ? 
      <div className="Message">
      <h2> { message.subject } </h2>
      <Relations type={"to"} />
      <Relations type={"cc"} />
      <Relations type={"bcc"} />
      <p> { message.text || message.html } </p>
      <Files />
    </div> : 

    <h2> Message not found </h2>
  )
}

export default withRouter(Message);
