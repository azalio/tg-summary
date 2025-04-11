package storage

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

const testDBPath = "test_storage.db"

func setupTestStorage(t *testing.T) *GormStorage {
	_ = os.Remove(testDBPath)
	st, err := NewGormStorage(testDBPath)
	require.NoError(t, err)
	err = st.Init(context.Background())
	require.NoError(t, err)
	return st
}

func teardownTestStorage() {
	_ = os.Remove(testDBPath)
}

func TestGormStorage_CRUD(t *testing.T) {
	st := setupTestStorage(t)
	defer teardownTestStorage()

	ctx := context.Background()

	chat := &Chat{ID: 123, Title: "Test Group", Type: "group"}
	user := &User{ID: 456, Username: "testuser", DisplayName: "Test User"}
	msg := &Message{
		ChatID:    chat.ID,
		MessageID: 1,
		AuthorID:  user.ID,
		Text:      "Hello, world!",
		Timestamp: time.Now().Unix(),
	}

	// Save chat and user
	require.NoError(t, st.SaveChat(ctx, chat))
	require.NoError(t, st.SaveUser(ctx, user))

	// Save message
	require.NoError(t, st.SaveMessage(ctx, msg))

	// Save duplicate message (should not error, DoNothing)
	require.NoError(t, st.SaveMessage(ctx, msg))

	// Get last message timestamp
	ts, err := st.GetLastMessageTimestamp(ctx, chat.ID)
	require.NoError(t, err)
	require.Equal(t, msg.Timestamp, ts)

	// Get messages after timestamp-1 (should return 1)
	msgs, err := st.GetMessagesAfter(ctx, chat.ID, msg.Timestamp-1)
	require.NoError(t, err)
	require.Len(t, msgs, 1)
	require.Equal(t, msg.Text, msgs[0].Text)
	require.Equal(t, msg.AuthorID, msgs[0].AuthorID)
}