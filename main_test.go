package main

import (
	"bufio"
	"bytes"
	"os"
	"strings"
	"testing"
)

func TestWithNewOld(t *testing.T) {
	testFiles(t, "1")
}

func TestOnlyOld(t *testing.T) {
	testFiles(t, "2")
}

func TestOnlyNew(t *testing.T) {
	testFiles(t, "3")
}

func testFiles(t *testing.T, suf string) {
	buf := &bytes.Buffer{}

	fp := "testdata/" + suf + ".in"
	f, err := os.Open(fp)
	if err != nil {
		t.Fatalf("open file %q error: %v", fp, err)
	}
	worker := Worker{
		size: 2,
		r:    f,
		w:    buf,
	}

	worker.Do()

	exp := strings.Join(lines(t, "testdata/"+suf+".out"), "\n") + "\n"
	res := buf.String()

	if exp != res {
		t.Errorf("Incorrect result. Expected\n%s\n, got\n%s\n", exp, res)
	}
}

func lines(t *testing.T, fp string) []string {
	_, err := os.Stat(fp)
	if err != nil {
		t.Fatal(err)
	}

	f, err := os.Open(fp)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	res := []string{}
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		res = append(res, sc.Text())
	}
	return res
}
