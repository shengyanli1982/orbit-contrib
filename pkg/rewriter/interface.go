package rewriter

// Callback 是一个回调接口，用于定义路径重写后的回调函数
// Callback is a callback interface, used to define the callback function after path rewriting
type Callback interface {
	// OnPathRewrited 是一个路径重写回调函数，当路径被重写后，会调用这个函数
	// OnPathRewrited is a callback function for path rewriting. When the path is rewritten, this function will be called
	OnPathRewrited(new, old string)
}

// emptyCallback 是一个空回调函数的实现，不执行任何操作
// emptyCallback is an implementation of an empty callback function that does nothing
type emptyCallback struct{}

// OnPathRewrited 是 emptyCallback 的 OnPathRewrited 方法的实现，不执行任何操作
// OnPathRewrited is the implementation of the OnPathRewrited method of emptyCallback, it does nothing
func (e *emptyCallback) OnPathRewrited(new, old string) {}