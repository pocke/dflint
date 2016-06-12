package main

import "testing"

func TestDockerfileNodes(t *testing.T) {
	d, err := newDockerfile([]byte(`FROM busybox
RUN ls
RUN echo 'hoge'
ENV foo bar
RUN sl
`), "/dev/null")
	if err != nil {
		t.Error(err)
	}

	res := d.Nodes("run")
	if len(res) != 3 {
		t.Errorf("Len should be 3, but got %d", len(res))
	}
}
