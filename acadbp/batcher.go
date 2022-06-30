package acadbp

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/schollz/progressbar/v3"
	"golang.org/x/term"
	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/unicode"
	"golang.org/x/text/transform"
)

type Batcher struct {
	accore string
	iflag  string
	sflag  string
	encode encoding.Encoding
	out    io.Writer
	numbat int
}

// NewBatcher returns a new batcher that batch-processes drawing files with accoreconsole
func NewBatcher(accorepath string) *Batcher {
	return &Batcher{
		accore: accorepath,
		iflag:  "/i",
		sflag:  "/s",
		encode: unicode.UTF16(unicode.LittleEndian, unicode.IgnoreBOM),
		out:    os.Stdout,
		numbat: runtime.NumCPU(),
	}
}

// SetOutput sets a writer to output result
func (b *Batcher) SetOutput(w io.Writer) {
	b.out = w
}

// Run runs batch processing
func (b *Batcher) Run(scrfile string) {
	dec := b.encode.NewDecoder()
	bar := b.makePbar(1)

	var buf bytes.Buffer
	defer func() {
		buf.WriteTo(b.out)
		bar.Add(1)
	}()
	ilog := log.New(&buf, "[INFO] ", log.Ldate|log.Ltime)
	elog := log.New(&buf, "[ERROR] ", log.Ldate|log.Ltime)
	ilog.Println("run script")
	out, err := exec.Command(b.accore, b.sflag, scrfile).CombinedOutput()
	if err != nil {
		elog.Println(err)
		return
	}
	dout, _, err := transform.Bytes(dec, out)
	if err != nil {
		elog.Println(err)
		return
	}
	fmt.Fprintln(&buf, string(dout))
	ilog.Println("end script")
}

// RunForEach runs batch processing for each input file
func (b *Batcher) RunForEach(scrfile string, files []string, ext string) {
	dec := b.encode.NewDecoder()
	bar := b.makePbar(len(files))

	sp := make(chan struct{}, b.numbat)
	var wg sync.WaitGroup
	for _, f := range files {
		sp <- struct{}{}
		wg.Add(1)
		go func(in string) {
			var buf bytes.Buffer
			defer func() {
				buf.WriteTo(b.out)
				bar.Add(1)
				<-sp
				wg.Done()
			}()
			ilog := log.New(&buf, "[INFO] ", log.Ldate|log.Ltime)
			elog := log.New(&buf, "[ERROR] ", log.Ldate|log.Ltime)
			ilog.Println("run script for " + filepath.Base(in))
			if err := createEmptyFile(in, ext); err != nil {
				elog.Println(err)
				return
			}
			out, err := exec.Command(b.accore, b.iflag, in, b.sflag, scrfile).CombinedOutput()
			if err != nil {
				elog.Println(err)
				return
			}
			dout, _, err := transform.Bytes(dec, out)
			if err != nil {
				elog.Println(err)
				return
			}
			fmt.Fprintln(&buf, string(dout))
			ilog.Println("end script")
		}(f)
		time.Sleep(50 * time.Millisecond)
	}
	wg.Wait()
}

func (b *Batcher) makePbar(max int) *progressbar.ProgressBar {
	return progressbar.NewOptions(max,
		progressbar.OptionUseANSICodes(true),
		progressbar.OptionSetVisibility(b.out != os.Stdout && term.IsTerminal(int(syscall.Stdout))),
		progressbar.OptionSetDescription("Running accoreconsole..."),
		progressbar.OptionSetRenderBlankState(true),
		progressbar.OptionOnCompletion(func() { time.Sleep(500 * time.Millisecond) }),
		progressbar.OptionShowCount(),
	)
}

// CheckAccore returns an error if accoreconsole path is invalid
func (b *Batcher) CheckAccore() error {
	if _, err := os.Stat(b.accore); err != nil {
		return errors.New("accoreconsole binary is not found")
	}
	return nil
}

func createEmptyFile(src string, ext string) error {
	if !strings.HasPrefix(ext, ".") {
		return nil
	}
	if _, err := os.Stat(src); err != nil {
		return err
	}
	dst := strings.TrimSuffix(src, filepath.Ext(src)) + ext
	if _, err := os.Stat(dst); err != nil {
		f, err := os.Create(dst)
		if err != nil {
			return err
		}
		f.Close()
	}
	return nil
}
