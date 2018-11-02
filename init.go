package search

const (
	// 文档ID前缀，保存文档及文档索引 docIndex
	_docIDPrefix byte = 'D'
	// 文档Key前缀 dbKey
	_dbKey2docIDPrefix byte = 'K'
	// docID -> docKey 前缀
	_docID2DocKeyPrefix byte = 'I'
	// 关键词前缀
	_wordIDPrefix byte = 'W'

	// 标题前缀
	_titleIDPrefix byte = 'T'
	// 标题前缀
	_summaryIDPrefix byte = 'S'
	// 内容前缀
	_contentIDPrefix byte = 'C'
)

var (
	// WordMaxLength 最大词长度
	WordMaxLength = 6
	// 索引前缀列表
	_IndexPrefixs = []byte{_titleIDPrefix, _summaryIDPrefix, _contentIDPrefix}
)
