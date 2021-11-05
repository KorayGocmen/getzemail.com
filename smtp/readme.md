# EmailReceiver

EmailReceiver is a receive-only smtp server that works with an API to offer a complete email receiver solution. It does not support IMAP/POP3 and only works with a web based email client which will pull the data from the API / S3.

This is a stripped down version of the Email Firewall we use for [Violetnorth.com](https://www.violetnorth.com).

---

## Overview

EmailReceiver works alongside an API and uses Redis as a cache and S3 as the main file storage.

```
Incoming Email → EmailReceiver → Redis
                                   ↓
                              EmailReceiver -> API
                                   ↓
                                  S3
```