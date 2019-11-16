package lldb

//最大KEY长度
const LLDB_KEY_LEN_MAX = 255
//scan搜索的最大值
const LLDB_SCAN_LIMIT = 1000

const (
	dt_kv = 'k'

	//hash
	dt_hash = 'h'; // hashmap(sorted by key)
	dt_hsize = 'H'; // key = size


	//queue
	dt_queue = 'q';
	dt_qsize = 'Q';

	//sset
	dt_sset = 's' //key|vlaue => ""
	dt_ssize = 'S' // key => size

	//zset
	dt_zset = 'x'; // key|value => score
	dt_zscore = 'z'; // key|score|value =>
	dt_zsize = 'Z'; // key = size
)

