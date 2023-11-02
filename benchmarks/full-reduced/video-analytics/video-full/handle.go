package function

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"flag"
	"fmt"
	"net"
	"os"
	"io/ioutil"
	"log"
	"os/exec"

	ctrdlog "github.com/containerd/containerd/log"
	log "github.com/sirupsen/logrus"
)

var (
	videoFragment []byte
	videoFile     *string
	AWS_S3_BUCKET = "vhive-video-bench"
)

const (
	INLINE = "INLINE"
	XDT    = "XDT"
	S3     = "S3"
)

// Handle an HTTP Request.
func Handle(ctx context.Context, res http.ResponseWriter, req *http.Request) {
	contents, err := ioutil.ReadFile("/etc/secret-volume/video-decoder")
	if err != nil {
      		log.Fatal(err)
   	}
	video-decoder := string(contents)

	videoFile = "reference/video.mp4"
	videoFragment, _ = os.ReadFile(*videoFile)
	log.Infof("read video fragment, size: %v\n", len(videoFragment))

	requestURL := video-decoder + ":80"
	req_url, err := http.NewRequest(http.MethodPost, requestURL, {"video": videoFragment})
	if err != nil {
		log.Fatal(err)
	}
	ret, err := exec.Command("python3 decoder.py " + req_url).Output()
	fmt.Fprintf(res, ret) // echo to caller
}
