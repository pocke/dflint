package main

type Problem struct {
	Line    int    `json:"line"`
	Column  int    `json:"column"`
	Length  int    `json:"length"`
	Path    string `json:"path"`
	Type    string `json:"type"`
	Message string `json:"message"`
}
