package acadbp

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/schollz/progressbar/v3"
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
}

// NewBatcher returns a new batcher that batch-processes drawing files with accoreconsole
func NewBatcher(accorepath string) *Batcher {
	return &Batcher{
		accore: accorepath,
		iflag:  "/i",
		sflag:  "/s",
		encode: unicode.UTF16(unicode.LittleEndian, unicode.IgnoreBOM),
		out:    os.Stdout,
	}
}

// SetOutput sets a writer to output result
func (b *Batcher) SetOutput(w io.Writer) {
	b.out = w
}

// Run runs batch processing
func (b *Batcher) Run(scrfile string) {
	inflog := log.New(b.out, "[INFO] ", log.Ldate|log.Ltime)
	errlog := log.New(b.out, "[ERROR] ", log.Ldate|log.Ltime)
	dec := b.encode.NewDecoder()

	bar := b.makePbar(1)
	out, err := exec.Command(b.accore, b.sflag, scrfile).CombinedOutput()
	if err != nil {
		errlog.Println(err)
	}
	bt, _, err := transform.Bytes(dec, out)
	if err != nil {
		errlog.Println(err)
	}
	inflog.Println("run script")
	fmt.Fprintln(b.out, string(bt))
	inflog.Println("end script")
	bar.Add(1)
}

// RunForEach runs batch processing for each input file
func (b *Batcher) RunForEach(scrfile string, files []string, ext string) {
	inflog := log.New(b.out, "[INFO] ", log.Ldate|log.Ltime)
	errlog := log.New(b.out, "[ERROR] ", log.Ldate|log.Ltime)
	dec := b.encode.NewDecoder()

	bar := b.makePbar(len(files))
	for i, f := range files {
		bar.Set(i)
		if err := createEmptyFile(f, ext); err != nil {
			errlog.Println(err)
			continue
		}
		out, err := exec.Command(b.accore, b.iflag, f, b.sflag, scrfile).CombinedOutput()
		if err != nil {
			errlog.Println(err)
			continue
		}
		bt, _, err := transform.Bytes(dec, out)
		if err != nil {
			errlog.Println(err)
			continue
		}
		inflog.Println("run script for " + filepath.Base(f))
		fmt.Fprintln(b.out, string(bt))
		inflog.Println("end script")
	}
	bar.Add(1)
}

func (b *Batcher) makePbar(max int) *progressbar.ProgressBar {
	return progressbar.NewOptions(max,
		progressbar.OptionUseANSICodes(true),
		progressbar.OptionSetVisibility(b.out != os.Stdout),
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
