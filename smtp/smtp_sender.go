package main

import (
	"fmt"
	"sort"

	"github.com/jhillyerd/enmime"
)

// smtpSend tries to send an smtp message to the provided upstreams.
func smtpSend(messageEncoded enmime.MailBuilder, upstreams []typeMailUpstream) error {
	// Sort the upstreams by priority.
	sort.Slice(upstreams, func(n1, n2 int) bool {
		return upstreams[n1].Priority < upstreams[n2].Priority
	})

	// Add the port number to the upstream if there is no port number.
	var sortedUpstreams []typeMailUpstream
	for _, upstream := range upstreams {
		upstream.Target = fmt.Sprintf("%s:%d", upstream.Target, config.Server.Port)
		sortedUpstreams = append(sortedUpstreams, upstream)
	}

	var err error
	for _, upstream := range sortedUpstreams {
		if err = messageEncoded.Send(upstream.Target, nil); err != nil {
			logger.Errorf("Failed to send message %v", err)
			continue
		}

		// If there is no error, relay succeded.
		return nil
	}

	return err
}
