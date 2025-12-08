package whatsapp

import (
	"context"
	"fmt"
	"strings"

	"github.com/sirupsen/logrus"
	"go.mau.fi/whatsmeow"
	waProto "go.mau.fi/whatsmeow/binary/proto"
	"go.mau.fi/whatsmeow/types"
)

func isNumber(s string) bool {
	for _, r := range s {
		if r < '0' || r > '9' {
			return false
		}
	}
	return true
}

func FindGroupJIDByName(name string) (string, error) {
	clientMu.RLock()
	c := client
	clientMu.RUnlock()

	if c == nil {
		return "", fmt.Errorf("WA client not initialized")
	}

	groups, err := c.GetJoinedGroups(context.Background())
	if err != nil {
		return "", err
	}

	for _, g := range groups {
		if strings.EqualFold(g.Name, name) {
			return g.JID.String(), nil
		}
	}

	return "", fmt.Errorf("group not found")
}

func jidListToStrings(jids []types.JID) []string {
	out := make([]string, 0, len(jids))
	for _, j := range jids {
		out = append(out, j.String())
	}
	return out
}

func SendText(ctx context.Context, to string, message string) error {
	clientMu.RLock()
	c := client
	clientMu.RUnlock()

	if c == nil {
		return fmt.Errorf("WA client not initialized")
	}

	if err := ensureConnected(c); err != nil {
		return err
	}

	var jid types.JID

	if isNumber(to) {
		jid = types.NewJID(to, types.DefaultUserServer)

	} else if strings.Contains(to, "@") {
		parsed, err := types.ParseJID(to)
		if err != nil {
			return fmt.Errorf("invalid JID: %w", err)
		}
		jid = parsed

	} else {
		groupJID, err := FindGroupJIDByName(to)
		if err != nil {
			return err
		}

		parsed, err := types.ParseJID(groupJID)
		if err != nil {
			return fmt.Errorf("failed to parse group JID: %w", err)
		}
		jid = parsed
	}

	if jid.Server == types.GroupServer {
		return sendGroupWithMention(ctx, c, jid, message)
	}

	msg := &waProto.Message{
		Conversation: &message,
	}

	_, err := c.SendMessage(ctx, jid, msg)
	if err != nil {
		logrus.Warnf("Failed sending message: %v", err)
		return err
	}

	return nil
}

func sendGroupWithMention(ctx context.Context, c *whatsmeow.Client, groupJID types.JID, message string) error {
	info, err := c.GetGroupInfo(context.Background(), groupJID)
	if err != nil {
		return fmt.Errorf("failed to get group info: %w", err)
	}

	var mentions []types.JID

	for _, p := range info.Participants {
		userPart := p.JID.User
		if strings.Contains(strings.ToLower(userPart), "meta") {
			continue
		}
		if !isNumber(userPart) {
			continue
		}
		mentions = append(mentions, p.JID)
	}

	if len(mentions) == 0 {
		msg := &waProto.Message{
			ExtendedTextMessage: &waProto.ExtendedTextMessage{
				Text: &message,
			},
		}
		_, err = c.SendMessage(ctx, groupJID, msg)
		if err != nil {
			return fmt.Errorf("send without mentions failed: %w", err)
		}
		return nil
	}

	msg := &waProto.Message{
		ExtendedTextMessage: &waProto.ExtendedTextMessage{
			Text: &message,
			ContextInfo: &waProto.ContextInfo{
				MentionedJID: jidListToStrings(mentions),
			},
		},
	}
	_, err = c.SendMessage(ctx, groupJID, msg)
	if err != nil {
		return fmt.Errorf("failed to send group mention: %w", err)
	}
	return nil
}
