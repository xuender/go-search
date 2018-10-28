package search

const (
	// 文档Key前缀 dbKey
	_dbKey2docIDPrefix byte = 'K'
	// 文档ID前缀，保存文档及文档索引 docIndex
	_docIDPrefix byte = 'D'
	// docID -> docKey 前缀
	_docID2DocKeyPrefix byte = 'I'
	// 标题前缀
	_titleIDPrefix byte = 'T'
	// 标题前缀
	_summaryIDPrefix byte = 'S'
	// 内容前缀
	_contentIDPrefix byte = 'C'
)

var (
	// 关键词前缀
	_words = []byte("Words")
	// WordMaxLength 最大词长度
	WordMaxLength = 6
	// 索引前缀列表
	_IndexPrefixs = []byte{_titleIDPrefix, _summaryIDPrefix, _contentIDPrefix}
)
