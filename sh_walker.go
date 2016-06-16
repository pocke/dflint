package main

import "github.com/mvdan/sh"

type ShWalker struct {
	onBinaryCmd func(*sh.BinaryCmd)
	onCallExpr  func(*sh.CallExpr)
}

func (w *ShWalker) Walk(stmt sh.Stmt) {
	switch s := stmt.Cmd.(type) {
	case *sh.BinaryCmd:
		if w.onBinaryCmd != nil {
			w.onBinaryCmd(s)
		}
		w.Walk(s.X)
		w.Walk(s.Y)
	case *sh.CallExpr:
		if w.onCallExpr != nil {
			w.onCallExpr(s)
		}
	}
}
