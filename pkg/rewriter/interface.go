package rewriter

// Callback 是一个回调接口
// Callback is a callback interface
type Callback interface {
	// OnPathRewrited 是一个路径重写回调函数
	// OnPathRewrited is a callback function to rewrite the path
	OnPathRewrited(new, old string)
}

// emptyCallback 是一个空回调函数，不执行任何操作
// emptyCallback is an empty callback function that does nothing
type emptyCallback struct{}

// OnLimited 是空回调函数，不执行任何操作
// OnLimited is an empty callback function that does nothing
func (e *emptyCallback) OnPathRewrited(new, old string) {}
