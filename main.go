package main

import (
)

func main() {
	blkC := NewBlockChain()
	str := "首先是它的树的结构,默克尔树常见的结构是二叉"
	blkC.AddBlock(str,3, 0)
	blkC.GetAllBlockHash()
}


