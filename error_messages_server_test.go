package mysqlproto

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestErrorMessagesServerSequence(t *testing.T) {
	assert.Equal(t, ER_HASHCHK, uint16(1000))
	assert.Equal(t, ER_TABLE_NOT_LOCKED, uint16(1100))
	assert.Equal(t, ER_BAD_SLAVE, uint16(1200))
	assert.Equal(t, ER_INVALID_CHARACTER_STRING, uint16(1300))
	assert.Equal(t, ER_XAER_OUTSIDE, uint16(1400))
	assert.Equal(t, ER_SUBPARTITION_ERROR, uint16(1500))
	assert.Equal(t, ER_VIEW_INVALID_CREATION_CTX, uint16(1600))
	assert.Equal(t, ER_GRANT_PLUGIN_USER_EXISTS, uint16(1700))
	assert.Equal(t, ER_UNKNOWN_ALTER_ALGORITHM, uint16(1800))
	assert.Equal(t, ER_SLAVE_HAS_MORE_GTIDS_THAN_MASTER, uint16(1885))

	assert.Equal(t, ER_FILE_CORRUPT, uint16(3000))
	assert.Equal(t, ER_RUN_HOOK_ERROR, uint16(3100))
	assert.Equal(t, ER_CANNOT_CREATE_VIRTUAL_INDEX_CONSTRAINT, uint16(3175))
}
