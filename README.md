# Getzemail

Receive only email inbox.

---

I have implemented an example public email inbox: [getzemail.com](http://getzemail.com/). 

If you would like to test it out, please go to [getzemail.com](http://getzemail.com/) and search for an inbox, when you search an inbox, if it doesn't exist, it's created. After searching for it, you can start sending emails to that address. Ex: search for "anything" and then send your emails to anything@getzemail.com.

Source code here: [https://github.com/KorayGocmen/getzemail.com](https://github.com/KorayGocmen/getzemail.com) 

## Implementation Details

![getzemail.com-overview.png](assets/getzemail.com-overview.png)

### Components

- SMTP: Implemented in Golang. Receives emails, uploads the email body as well as attachments to S3 and sends the details to the API. Email body is clipped to 255 characters and stored in the database for quick display but the actual contents are downloaded on client side via the pre-signed S3 URLs in order not to overwhelm the server with IO and offload it to clients.
- API: Implemented in Golang. Implements mail, mail_inbox and mail_message objects for storing the details of the mail inboxes and mail messages. Pre-signs S3 URLs for Web application to fetch directly from the S3 bucket. Provides mail inbox discovery to SMTP server. PostgreSQL database is used but supports Mysql and Sqlite.
- WEB: Implemented in React with a small Golang web server to serve React resources. Receives mail inbox details from the API, fetches mail contents and attachments from S3.

### Data types

![getzemail.com-datatypes.jpg](assets/getzemail.com-datatypes.jpg)

- `mail`: The overarching datatype that stores the host of the mail. For all getzemails, the `mail` type is the same and the host is `getzemail.com`, this allows us to use the API for multiple email hosts. For example we can use `[windscribe.com](http://windscribe.com)` as well as `[controld.com](http://controld.com)` for two different email services.
- `mail_inbox`: `mail` has multiple mail inboxes, koray@[windscribe.com](http://windscribe.com) will be linked to windscribe.com mail and koray@[controld.com](http://controld.com) will be linked to controld.com mail. `address` field for the first one would be `koray@windscribe.com` and vice versa for `koray@controld.com`.
- `mail_message`: All mail messages for a certain mail inbox will be stored under this table. The first 255 chars of the mail message text and html are stored for display purposes and the full path of the text and html is `$s3_bucket/$message_id/text` and `$s3_bucket/$message_id/html` under S3 respectively.
- `mail_message_relation`: Possible `type` fields are `to`, `cc` and `bcc`. Stores the address and the display name for the relation.
- `mail_message_file`: The S3 bucket details for `inline` and `attachment` disposition types. the full path of any file is `$s3_bucket/$disposition/$content_id` under S3.

### S3 Bucket Contents

![getzemail.com-bucketdetails.png](assets/getzemail.com-bucketdetails.png)

Emails are stored under the path of the `mail_message.message_id` in the bucket. Stores the Raw MIME data, as well as the text and html data. Any attachments including inline and attachment dispositions will also be stored under the same path. For attachments, the S3 key and url will be saved under `mail_message_files` in the database.

### Full cycle, step by step

- Email for koray@getzemail.com is received by the SMTP server. SMTP server checks Redis for a mail instance with the hostname "getzemail.com". 
    1. If mail instance not found in Redis, `GET /mails/getzemail.com` request to API.
    2. Save the mail instance to Redis.
- Check if mail inbox (ie. koray@getzemail.com) is a known mail inbox for the mail instance.
    1. If not, reject the email.
- Parse mime type and upload mail message and any attachments to S3.
- Send new mail message to API.
- API receives mail message, saves to database.
- User visits [getzemail.com](http://getzemail.com) and searches "koray" inbox.
- API receives `GET /mails/getzemail.com/inboxes/koray` from the Web.
- API returns 

```
{
  "mail_inbox": {
    "id": 1,
    "created_at": "2021-11-04T22:58:38.293324Z",
    "updated_at": "2021-11-04T22:58:38.293324Z",
    "deleted_at": null,
    "mail": 1,
    "address": "koray@getzemail.com",
    "display_name": "Koray",
    "mail_messages": [
      {
        "id": 1,
        "created_at": "2021-11-05T22:43:43.41723Z",
        "updated_at": "2021-11-05T22:43:43.41723Z",
        "deleted_at": null,
        "mail_inbox": 1,
        "message_id": "CAPFK=UQuSk4=rZQriRi_SpHkgn8WEN7fFzH3aQS AE0OM4QvZA@mail.gmail.com",
        "in_reply_to_id": "CAPFK=UTPR6n-nN3jFg08s0y7th=pOf a48fiDWN4523Q=EhFvA@mail.gmail.com",
        "subject": "Re: hey, test email",
        "text": "reply to the test email\r\n\r\nOn Fri, 5 Nov 2021 at 17:32, Koray Göçmen <gocmen.koray@gmail.com> wrote:\r\n>\r\n> hello\r\n",
        "html": "",
        "mail_message_relations": [
          {
            "id": 4,
            "created_at": "2021-11-05T22:43:43.418051Z",
            "updated_at": "2021-11-05T22:43:43.418051Z",
            "deleted_at": null,
            "mail_message": 4,
            "type": "to",
            "address": "koray@getzemail.com",
            "display_name": ""
          }
        ]
      }
    ]
  }
}
```

- User clicks on the first mail message.
- API receives `GET /mails/getzemail.com/messages/1`.
- API pre-signs the S3 urls for text, html and any mail message files and returns the mail message details to Web 

```
{
  "mail_message": {
    "id": 1,
    "created_at": "2021-11-05T22:43:43.41723Z",
    "updated_at": "2021-11-05T22:43:43.41723Z",
    "deleted_at": null,
    "mail_inbox": 1,
    "message_id": "CAPFK=UQuSk4=rZQriRi_SpHkgn8WEN7fFzH3aQS AE0OM4QvZA@mail.gmail.com",
    "in_reply_to_id": "CAPFK=UTPR6n-nN3jFg08s0y7th=pOf a48fiDWN4523Q=EhFvA@mail.gmail.com",
    "subject": "Re: hey, test email",
    "text": "reply to the test email\r\n\r\nOn Fri, 5 Nov 2021 at 17:32, Koray Göçmen <gocmen.koray@gmail.com> wrote:\r\n>\r\n> hello\r\n",
    "html": "",
    "mail_message_relations": [
      {
        "id": 1,
        "created_at": "2021-11-05T22:43:43.418051Z",
        "updated_at": "2021-11-05T22:43:43.418051Z",
        "deleted_at": null,
        "mail_message": 1,
        "type": "to",
        "address": "koray@getzemail.com",This proposal is for receiving only email inboxes for Windscribe and ControlD users. I believe this can be a key feature for converting free users to pro users. Receiving only email inboxes can be used as a throwaway email address for spammy services that require email addresses or it can be used for privacy reasons via using different aliases.
        "display_name": ""
      }
    ],
    "text_url": "https://drop-inbox.s3.amazonaws.com/CAPFK%3DUQuSk4%3DrZQriRi_SpHkgn8WEN7fFzH3aQS%20AE0OM4QvZA%40mail.gmail.com/text?X-Amz-Algorithm=AWS4-HMAC-SHA256&X-Amz-Credential=AKIAUOD6FWIM4OLQZLHO%2F20211107%2Fus-east-1%2Fs3%2Faws4_request&X-Amz-Date=20211107T192548Z&X-Amz-Expires=900&X-Amz-SignedHeaders=host&X-Amz-Signature=ee1b0e6c875ccf01adaca4c036a5b5ae44c0bddd90d44941c3ab7b5e3b2759f6",
    "html_url": "https://drop-inbox.s3.amazonaws.com/CAPFK%3DUQuSk4%3DrZQriRi_SpHkgn8WEN7fFzH3aQS%20AE0OM4QvZA%40mail.gmail.com/html?X-Amz-Algorithm=AWS4-HMAC-SHA256&X-Amz-Credential=AKIAUOD6FWIM4OLQZLHO%2F20211107%2Fus-east-1%2Fs3%2Faws4_request&X-Amz-Date=20211107T192548Z&X-Amz-Expires=900&X-Amz-SignedHeaders=host&X-Amz-Signature=717ccedc103250a12832e3dedeb36fbe9943628e5a30c86d56be9f0348663faf"
  },
  "success": true
}
```

- Web fetches the text, html and any attachments directly from S3 before rendering the message details.